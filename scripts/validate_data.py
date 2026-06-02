#!/usr/bin/env python3
"""
scripts/validate_data.py

Apply validation rules to ensure data integrity for Meta Earth Phase I.
Covers high-value trust boundaries: Bridge, Sequencer, Governance, Reward,
DID/KYC, and Validator logic.

Usage:
    python scripts/validate_data.py --input raw_data.json --output validated_data.json
"""

import argparse
import json
import logging
import sys
import hashlib
import hmac
import time
import asyncio
from dataclasses import dataclass, field, asdict
from datetime import datetime, timezone
from enum import Enum
from pathlib import Path
from typing import Any, Dict, List, Optional, Set, Tuple, Union, TypeVar, Generic, Callable, AsyncIterator
from collections import defaultdict
from functools import lru_cache, wraps
import re
import os
import signal
from contextlib import contextmanager, asynccontextmanager
from concurrent.futures import ThreadPoolExecutor, ProcessPoolExecutor
from threading import Lock
import traceback
import uuid
from abc import ABC, abstractmethod

# ---------------------------------------------------------------------------
# Logging Configuration
# ---------------------------------------------------------------------------
class LoggingConfig:
    """Centralized logging configuration with rotation support and structured logging."""
    
    _instance: Optional[logging.Logger] = None
    _lock: Lock = Lock()
    
    @classmethod
    def get_logger(cls) -> logging.Logger:
        """Get or create the singleton logger instance."""
        if cls._instance is None:
            with cls._lock:
                if cls._instance is None:
                    cls._instance = cls._setup_logging()
        return cls._instance
    
    @classmethod
    def _setup_logging(cls, log_level: str = "INFO", log_file: Optional[Path] = None) -> logging.Logger:
        """Configure logging with proper formatting and optional file output."""
        logger = logging.getLogger("validate_data")
        logger.setLevel(getattr(logging, log_level.upper(), logging.INFO))
        
        # Structured formatter for better log parsing
        formatter = logging.Formatter(
            '{"timestamp": "%(asctime)s", "level": "%(levelname)s", "logger": "%(name)s", '
            '"module": "%(module)s", "function": "%(funcName)s", "line": %(lineno)d, '
            '"message": "%(message)s"}',
            datefmt="%Y-%m-%dT%H:%M:%S.%fZ"
        )
        
        # Console handler with stderr for separation from stdout
        console_handler = logging.StreamHandler(sys.stderr)
        console_handler.setFormatter(formatter)
        logger.addHandler(console_handler)
        
        # File handler if specified with rotation support
        if log_file:
            from logging.handlers import RotatingFileHandler
            file_handler = RotatingFileHandler(
                log_file,
                maxBytes=10*1024*1024,  # 10MB
                backupCount=5
            )
            file_handler.setFormatter(formatter)
            logger.addHandler(file_handler)
        
        return logger

logger = LoggingConfig.get_logger()


# ---------------------------------------------------------------------------
# Custom Exceptions
# ---------------------------------------------------------------------------
class ValidationError(Exception):
    """Base exception for validation errors."""
    def __init__(self, message: str, code: str = "VALIDATION_ERROR", details: Optional[Dict[str, Any]] = None):
        self.code = code
        self.details = details or {}
        super().__init__(message)


class DataIntegrityError(ValidationError):
    """Raised when data integrity checks fail."""
    def __init__(self, message: str, details: Optional[Dict[str, Any]] = None):
        super().__init__(message, code="DATA_INTEGRITY_ERROR", details=details)


class SecurityValidationError(ValidationError):
    """Raised when security-related validation fails."""
    def __init__(self, message: str, details: Optional[Dict[str, Any]] = None):
        super().__init__(message, code="SECURITY_VALIDATION_ERROR", details=details)


class ConfigurationError(ValidationError):
    """Raised when configuration is invalid."""
    def __init__(self, message: str, details: Optional[Dict[str, Any]] = None):
        super().__init__(message, code="CONFIGURATION_ERROR", details=details)


class ResourceExhaustedError(ValidationError):
    """Raised when system resources are exhausted."""
    def __init__(self, message: str, details: Optional[Dict[str, Any]] = None):
        super().__init__(message, code="RESOURCE_EXHAUSTED", details=details)


# ---------------------------------------------------------------------------
# Enums & Constants
# ---------------------------------------------------------------------------
class RiskLevel(Enum):
    CRITICAL = "critical"
    HIGH = "high"
    MEDIUM = "medium"
    LOW = "low"
    INFO = "info"
    
    @classmethod
    def from_string(cls, value: str) -> "RiskLevel":
        """Convert string to RiskLevel with validation."""
        try:
            return cls(value.lower())
        except ValueError:
            raise ConfigurationError(f"Invalid risk level: {value}")


class Component(Enum):
    BRIDGE = "bridge"
    SEQUENCER = "sequencer"
    GOVERNANCE = "governance"
    REWARD = "reward"
    DID_KYC = "did_kyc"
    VALIDATOR = "validator"
    ROLLAPP = "rollapp"
    TOKEN_MINT = "token_mint"
    GAS = "gas"
    
    @classmethod
    def from_string(cls, value: str) -> "Component":
        """Convert string to Component with validation."""
        try:
            return cls(value.lower())
        except ValueError:
            raise ConfigurationError(f"Invalid component: {value}")


# Mapping of component to its risk notes (from spec)
COMPONENT_RISK_NOTES: Dict[Component, str] = {
    Component.BRIDGE: "Fund theft — replay attacks, forged messages, fake proofs, nonce reuse, chain ID confusion, partial verification",
    Component.SEQUENCER: "Fake state submissions — malicious sequencer can submit invalid state roots",
    Component.GOVERNANCE: "Unauthorized control — vote inflation, snapshot manipulation, proposal replay, timing attacks, quorum bypass",
    Component.REWARD: "Infinite rewards — reward accounting bugs, incorrect distribution logic",
    Component.DID_KYC: "Identity bypass — forged identity proofs, replay of KYC attestations",
    Component.VALIDATOR: "Consensus manipulation — validator set tampering, double-signing, liveness attacks",
    Component.ROLLAPP: "Settlement fraud — invalid state transitions, incorrect execution results",
    Component.TOKEN_MINT: "Unlimited minting — missing access control on mint functions",
    Component.GAS: "Gas manipulation — underpriced operations, gas token inflation",
}

# Validation constants with security considerations
MAX_NONCE_LENGTH = 256
MAX_MESSAGE_LENGTH = 1024 * 1024  # 1MB
MAX_SIGNATURE_LENGTH = 1024
MAX_BATCH_SIZE = 10000
MAX_RECURSION_DEPTH = 100
TIMEOUT_SECONDS = 30
RATE_LIMIT_CALLS = 100
RATE_LIMIT_WINDOW = 60  # seconds

# Secure field patterns
REQUIRED_FIELD_PATTERNS: Dict[str, re.Pattern] = {
    "address": re.compile(r"^0x[a-fA-F0-9]{40}$"),
    "hash": re.compile(r"^0x[a-fA-F0-9]{64}$"),
    "signature": re.compile(r"^0x[a-fA-F0-9]{130}$"),
    "public_key": re.compile(r"^0x[a-fA-F0-9]{66}$"),
    "private_key": re.compile(r"^0x[a-fA-F0-9]{64}$"),
    "mnemonic": re.compile(r"^[a-z]+( [a-z]+){11,23}$"),
}

# Sensitive fields that should never be logged
SENSITIVE_FIELDS: Set[str] = {
    "private_key", "mnemonic", "seed", "password", "secret", "token", "api_key"
}


# ---------------------------------------------------------------------------
# Security Utilities
# ---------------------------------------------------------------------------
class SecurityUtils:
    """Security utilities for validation operations."""
    
    @staticmethod
    def sanitize_for_logging(data: Dict[str, Any]) -> Dict[str, Any]:
        """Remove sensitive fields from data for logging purposes."""
        sanitized = {}
        for key, value in data.items():
            if key.lower() in SENSITIVE_FIELDS:
                sanitized[key] = "***REDACTED***"
            elif isinstance(value, dict):
                sanitized[key] = SecurityUtils.sanitize_for_logging(value)
            elif isinstance(value, list):
                sanitized[key] = [
                    SecurityUtils.sanitize_for_logging(item) if isinstance(item, dict) else item
                    for item in value
                ]
            else:
                sanitized[key] = value
        return sanitized
    
    @staticmethod
    def verify_signature(data: bytes, signature: bytes, public_key: bytes) -> bool:
        """Verify a cryptographic signature."""
        try:
            # Placeholder for actual signature verification
            # In production, use proper cryptographic libraries
            return hmac.compare_digest(
                hashlib.sha256(data).digest(),
                hashlib.sha256(signature).digest()
            )
        except Exception as e:
            logger.error(f"Signature verification failed: {e}")
            return False
    
    @staticmethod
    def hash_data(data: Any) -> str:
        """Create a secure hash of data."""
        try:
            if isinstance(data, str):
                data = data.encode('utf-8')
            elif isinstance(data, dict):
                data = json.dumps(data, sort_keys=True).encode('utf-8')
            return hashlib.sha256(data).hexdigest()
        except Exception as e:
            logger.error(f"Data hashing failed: {e}")
            raise SecurityValidationError("Failed to hash data", {"error": str(e)})


# ---------------------------------------------------------------------------
# Rate Limiter
# ---------------------------------------------------------------------------
class RateLimiter:
    """Thread-safe rate limiter for API calls."""
    
    def __init__(self, max_calls: int = RATE_LIMIT_CALLS, window: int = RATE_LIMIT_WINDOW):
        self.max_calls = max_calls
        self.window = window
        self._calls: List[float] = []
        self._lock = Lock()
    
    def acquire(self) -> bool:
        """Try to acquire a rate limit slot. Returns True if allowed."""
        with self._lock:
            now = time.time()
            # Remove old calls
            self._calls = [t for t in self._calls if now - t < self.window]
            
            if len(self._calls) >= self.max_calls:
                return False
            
            self._calls.append(now)
            return True
    
    def wait_and_acquire(self, timeout: float = 10.0) -> bool:
        """Wait for a rate limit slot with timeout."""
        start = time.time()
        while time.time() - start < timeout:
            if self.acquire():
                return True
            time.sleep(0.1)
        return False


# ---------------------------------------------------------------------------
# Circuit Breaker
# ---------------------------------------------------------------------------
class CircuitBreaker:
    """Circuit breaker pattern for fault tolerance."""
    
    STATE_CLOSED = "closed"
    STATE_OPEN = "open"
    STATE_HALF_OPEN = "half_open"
    
    def __init__(self, failure_threshold: int = 5, recovery_timeout: float = 30.0):
        self.failure_threshold = failure_threshold
        self.recovery_timeout = recovery_timeout
        self._state = self.STATE_CLOSED
        self._failure_count = 0
        self._last_failure_time = 0.0
        self._lock = Lock()
    
    @property
    def state(self) -> str:
        return self._state
    
    def call(self, func: Callable, *args, **kwargs) -> Any:
        """Execute a function with circuit breaker protection."""
        with self._lock:
            if self._state == self.STATE_OPEN:
                if time.time() - self._last_failure_time >= self.recovery_timeout:
                    self._state = self.STATE_HALF_OPEN
                else:
                    raise ResourceExhaustedError("Circuit breaker is open")
        
        try:
            result = func(*args, **kwargs)
            with self._lock:
                if self._state == self.STATE_HALF_OPEN:
                    self._state = self.STATE_CLOSED
                self._failure_count = 0
            return result
        except Exception as e:
            with self._lock:
                self._failure_count += 1
                self._last_failure_time = time.time()
                if self._failure_count >= self.failure_threshold:
                    self._state = self.STATE_OPEN
            raise


# ---------------------------------------------------------------------------
# Data Models
# ---------------------------------------------------------------------------
@dataclass(frozen=True)
class ValidationResult:
    """Immutable result of a single validation check."""
    rule_name: str
    component: Component
    risk_level: RiskLevel
    passed: bool
    message: str
    details: Dict[str, Any] = field(default_factory=dict)
    timestamp: str = field(default_factory=lambda: datetime.now(timezone.utc).isoformat())
    correlation_id: str = field(default_factory=lambda: str(uuid.uuid4()))
    
    def __post_init__(self) -> None:
        """Validate the result after initialization."""
        if not self.rule_name or not self.rule_name.strip():
            raise ValueError("rule_name cannot be empty")
        if not self.message or not self.message.strip():
            raise ValueError("message cannot be empty")
        if not isinstance(self.component, Component):
            raise TypeError("component must be a Component enum")
        if not isinstance(self.risk_level, RiskLevel):
            raise TypeError("risk_level must be a RiskLevel enum")
    
    def to_dict(self) -> Dict[str, Any]:
        """Convert to dictionary for serialization."""
        return {
            "rule_name": self.rule_name,
            "component": self.component.value,
            "risk_level": self.risk_level.value,
            "passed": self.passed,
            "message": self.message,
            "details": self.details,
            "timestamp": self.timestamp,
            "correlation_id": self.correlation_id
        }


@dataclass
class ValidatedRecord:
    """A single validated data record with comprehensive tracking."""
    record_id: str
    source: str
    raw_data: Dict[str, Any]
    validation_results: List[ValidationResult] = field(default_factory=list)
    is_valid: bool = True
    errors: List[str] = field(default_factory=list)
    warnings: List[str] = field(default_factory=list)
    metadata: Dict[str, Any] = field(default_factory=dict)
    created_at: str = field(default_factory=lambda: datetime.now(timezone.utc).isoformat())
    
    def __post_init__(self) -> None:
        """Validate the record after initialization."""
        if not self.record_id or not self.record_id.strip():
            raise ValueError("record_id cannot be empty")
        if not self.source or not self.source.strip():
            raise ValueError("source cannot be empty")
        if not isinstance(self.raw_data, dict):
            raise TypeError("raw_data must be a dictionary")
    
    def add_validation_result(self, result: ValidationResult) -> None:
        """Add a validation result and update record status."""
        if not isinstance(result, ValidationResult):
            raise TypeError("result must be a ValidationResult")
        
        self.validation_results.append(result)
        if not result.passed:
            if result.risk_level in (RiskLevel.CRITICAL, RiskLevel.HIGH):
                self.is_valid = False
                self.errors.append(f"{result.rule_name}: {result.message}")
            else:
                self.warnings.append(f"{result.rule_name}: {result.message}")
    
    def to_dict(self) -> Dict[str, Any]:
        """Convert to dictionary for serialization."""
        return {
            "record_id": self.record_id,
            "source": self.source,
            "raw_data": self.raw_data,
            "validation_results": [r.to_dict() for r in self.validation_results],
            "is_valid": self.is_valid,
            "errors": self.errors,
            "warnings": self.warnings,
            "metadata": self.metadata,
            "created_at": self.created_at
        }


@dataclass
class ValidationReport:
    """Aggregate validation report with comprehensive statistics."""
    total_records: int = 0
    valid_records: int = 0
    invalid_records: int = 0
    critical_issues: int = 0
    high_issues: int = 0
    medium_issues: int = 0
    low_issues: int = 0
    info_issues: int = 0
    component_breakdown: Dict[str, Dict[str, int]] = field(default_factory=lambda: defaultdict(lambda: defaultdict(int)))
    records: List[ValidatedRecord] = field(default_factory=list)
    generated_at: str = field(default_factory=lambda: datetime.now(timezone.utc).isoformat())
    execution_time_ms: float = 0.0
    report_id: str = field(default_factory=lambda: str(uuid.uuid4()))
    
    def update_statistics(self) -> None:
        """Update aggregate statistics from records."""
        self.total_records = len(self.records)
        self.valid_records = sum(1 for r in self.records if r.is_valid)
        self.invalid_records = self.total_records - self.valid_records
        
        # Reset counters
        self.critical_issues = 0
        self.high_issues = 0
        self.medium_issues = 0
        self.low_issues = 0
        self.info_issues = 0
        self.component_breakdown.clear()
        
        for record in self.records:
            for result in record.validation_results:
                component_name = result.component.value
                if not result.passed:
                    if result.risk_level == RiskLevel.CRITICAL:
                        self.critical_issues += 1
                    elif result.risk_level == RiskLevel.HIGH:
                        self.high_issues += 1
                    elif result.risk_level == RiskLevel.MEDIUM:
                        self.medium_issues += 1
                    elif result.risk_level == Risk