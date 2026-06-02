#!/usr/bin/env python3
"""
scripts/deliver_results.py

Package final deliverables into JSON, CSV, and markdown formats.
Part of the BountyScout data pipeline delivery stage.

This module provides production-grade functionality for transforming
analysis data into multiple output formats with comprehensive error
handling, logging, validation, and performance optimization.
"""

import json
import csv
import os
import sys
import logging
from datetime import datetime, timezone
from pathlib import Path
from typing import Any, Dict, List, Optional, Union, Sequence, Mapping, Set, Tuple
from dataclasses import dataclass, field, asdict
from enum import Enum, auto
from contextlib import contextmanager
import traceback
from functools import wraps
import hashlib
import io
from concurrent.futures import ThreadPoolExecutor, as_completed
import threading
from collections import defaultdict

# ---------------------------------------------------------------------------
# Constants & Configuration
# ---------------------------------------------------------------------------

DEFAULT_OUTPUT_DIR = Path("output/deliverables")
DEFAULT_ANALYSIS_FILE = Path("output/analysis/analysis_output.json")
DEFAULT_ENCODING = "utf-8"
DEFAULT_JSON_INDENT = 2
MAX_FILE_SIZE_BYTES = 100 * 1024 * 1024  # 100MB
MAX_BATCH_SIZE = 1000
THREAD_POOL_SIZE = 4
CACHE_TTL_SECONDS = 300

SCOPE_MAP_COMPONENTS: List[str] = [
    "Bridge",
    "Sequencer logic",
    "Governance",
    "Reward logic",
    "DID/KYC",
    "Validator logic",
]

RISK_NOTES_TEMPLATE: Dict[str, str] = {
    "Bridge": "Fund theft",
    "Sequencer logic": "Fake state submissions",
    "Governance": "Unauthorized control",
    "Reward logic": "Infinite rewards",
    "DID/KYC": "Identity bypass",
    "Validator logic": "Consensus manipulation",
}

VALID_PIPELINE_STAGES: Set[str] = {"delivery", "analysis", "collection", "processing"}
VALID_OUTPUT_FORMATS: Set[str] = {"json", "csv", "markdown"}

# ---------------------------------------------------------------------------
# Custom Exceptions
# ---------------------------------------------------------------------------

class DeliverableError(Exception):
    """Base exception for deliverable-related errors."""
    pass

class ValidationError(DeliverableError):
    """Raised when input data validation fails."""
    pass

class FormattingError(DeliverableError):
    """Raised when output formatting fails."""
    pass

class FileWriteError(DeliverableError):
    """Raised when file writing operations fail."""
    pass

class DataParsingError(DeliverableError):
    """Raised when data parsing fails."""
    pass

class SecurityError(DeliverableError):
    """Raised when security validation fails."""
    pass

class ResourceLimitError(DeliverableError):
    """Raised when resource limits are exceeded."""
    pass

# ---------------------------------------------------------------------------
# Enums
# ---------------------------------------------------------------------------

class OutputFormat(Enum):
    """Supported output formats."""
    JSON = auto()
    CSV = auto()
    MARKDOWN = auto()

class LogLevel(Enum):
    """Logging levels for configuration."""
    DEBUG = logging.DEBUG
    INFO = logging.INFO
    WARNING = logging.WARNING
    ERROR = logging.ERROR
    CRITICAL = logging.CRITICAL

class ValidationSeverity(Enum):
    """Severity levels for validation issues."""
    CRITICAL = auto()
    HIGH = auto()
    MEDIUM = auto()
    LOW = auto()
    INFO = auto()

# ---------------------------------------------------------------------------
# Logging Setup
# ---------------------------------------------------------------------------

class StructuredFormatter(logging.Formatter):
    """Custom formatter with structured output."""
    
    def format(self, record: logging.LogRecord) -> str:
        """Format log record with additional context."""
        record.component = getattr(record, 'component', 'unknown')
        record.correlation_id = getattr(record, 'correlation_id', 'N/A')
        return super().format(record)

def setup_logging(
    level: LogLevel = LogLevel.INFO,
    log_file: Optional[Path] = None,
    correlation_id: Optional[str] = None
) -> logging.Logger:
    """
    Configure logging with structured format and optional file output.
    
    Args:
        level: Logging level
        log_file: Optional path to log file
        correlation_id: Optional correlation ID for request tracing
        
    Returns:
        Configured logger instance
        
    Raises:
        ValidationError: If log_file path is invalid
    """
    logger = logging.getLogger("deliver_results")
    logger.setLevel(level.value)
    
    # Prevent duplicate handlers
    if logger.handlers:
        return logger
    
    formatter = StructuredFormatter(
        fmt="%(asctime)s [%(levelname)s] [%(correlation_id)s] %(name)s (%(filename)s:%(lineno)d): %(message)s",
        datefmt="%Y-%m-%dT%H:%M:%S%z"
    )
    
    # Console handler
    console_handler = logging.StreamHandler(sys.stdout)
    console_handler.setFormatter(formatter)
    logger.addHandler(console_handler)
    
    # File handler (optional)
    if log_file:
        try:
            log_file = Path(log_file)
            log_file.parent.mkdir(parents=True, exist_ok=True)
            file_handler = logging.FileHandler(log_file, encoding=DEFAULT_ENCODING)
            file_handler.setFormatter(formatter)
            logger.addHandler(file_handler)
        except (OSError, PermissionError) as e:
            logger.warning(f"Failed to create log file {log_file}: {e}")
    
    return logger

logger = setup_logging()

# ---------------------------------------------------------------------------
# Thread-safe Cache
# ---------------------------------------------------------------------------

class ThreadSafeCache:
    """Thread-safe cache with TTL support."""
    
    def __init__(self, ttl_seconds: int = CACHE_TTL_SECONDS):
        """
        Initialize cache.
        
        Args:
            ttl_seconds: Time-to-live in seconds for cache entries
        """
        self._cache: Dict[str, Tuple[Any, float]] = {}
        self._lock = threading.RLock()
        self._ttl_seconds = ttl_seconds
    
    def get(self, key: str) -> Optional[Any]:
        """
        Get value from cache.
        
        Args:
            key: Cache key
            
        Returns:
            Cached value or None if not found/expired
        """
        with self._lock:
            if key not in self._cache:
                return None
            
            value, timestamp = self._cache[key]
            if time.time() - timestamp > self._ttl_seconds:
                del self._cache[key]
                return None
            
            return value
    
    def set(self, key: str, value: Any) -> None:
        """
        Set value in cache.
        
        Args:
            key: Cache key
            value: Value to cache
        """
        with self._lock:
            self._cache[key] = (value, time.time())
    
    def clear(self) -> None:
        """Clear all cache entries."""
        with self._lock:
            self._cache.clear()
    
    def invalidate(self, key: str) -> None:
        """
        Invalidate specific cache entry.
        
        Args:
            key: Cache key to invalidate
        """
        with self._lock:
            self._cache.pop(key, None)

_cache = ThreadSafeCache()

# ---------------------------------------------------------------------------
# Decorators
# ---------------------------------------------------------------------------

def log_exceptions(func):
    """Decorator to log exceptions with traceback."""
    @wraps(func)
    def wrapper(*args, **kwargs):
        try:
            return func(*args, **kwargs)
        except Exception as e:
            logger.error(
                "Exception in %s: %s\n%s",
                func.__name__,
                str(e),
                traceback.format_exc(),
                extra={'component': func.__name__}
            )
            raise
    return wrapper

def validate_file_path(func):
    """Decorator to validate file paths."""
    @wraps(func)
    def wrapper(*args, **kwargs):
        # Extract filepath from args or kwargs
        filepath = None
        for arg in args:
            if isinstance(arg, Path):
                filepath = arg
                break
        if not filepath and 'filepath' in kwargs:
            filepath = kwargs['filepath']
        
        if filepath:
            if not isinstance(filepath, Path):
                raise ValidationError(f"File path must be Path instance, got {type(filepath)}")
            if filepath.exists() and not filepath.is_file():
                raise ValidationError(f"Path exists but is not a file: {filepath}")
            if filepath.exists() and filepath.stat().st_size > MAX_FILE_SIZE_BYTES:
                raise ResourceLimitError(f"File exceeds maximum size: {filepath}")
        
        return func(*args, **kwargs)
    return wrapper

def retry_on_failure(max_retries: int = 3, delay: float = 1.0):
    """
    Decorator to retry function on failure.
    
    Args:
        max_retries: Maximum number of retries
        delay: Delay between retries in seconds
    """
    def decorator(func):
        @wraps(func)
        def wrapper(*args, **kwargs):
            last_exception = None
            for attempt in range(max_retries):
                try:
                    return func(*args, **kwargs)
                except (IOError, OSError) as e:
                    last_exception = e
                    if attempt < max_retries - 1:
                        logger.warning(
                            f"Retry {attempt + 1}/{max_retries} for {func.__name__}: {e}"
                        )
                        time.sleep(delay * (attempt + 1))
                    else:
                        raise
            raise last_exception
        return wrapper
    return decorator

def measure_performance(func):
    """Decorator to measure function performance."""
    @wraps(func)
    def wrapper(*args, **kwargs):
        start_time = time.perf_counter()
        try:
            result = func(*args, **kwargs)
            elapsed = time.perf_counter() - start_time
            logger.debug(
                f"Performance: {func.__name__} took {elapsed:.4f}s",
                extra={'component': 'performance'}
            )
            return result
        except Exception as e:
            elapsed = time.perf_counter() - start_time
            logger.error(
                f"Performance: {func.__name__} failed after {elapsed:.4f}s: {e}",
                extra={'component': 'performance'}
            )
            raise
    return wrapper

# ---------------------------------------------------------------------------
# Data Models
# ---------------------------------------------------------------------------

@dataclass
class ScopeMapEntry:
    """Represents a single scope map entry."""
    component: str
    risk: str
    
    def __post_init__(self) -> None:
        """Validate fields after initialization."""
        self._validate()
    
    def _validate(self) -> None:
        """Internal validation method."""
        if not self.component or not isinstance(self.component, str):
            raise ValidationError(f"Invalid component: {self.component}")
        if not self.risk or not isinstance(self.risk, str):
            raise ValidationError(f"Invalid risk: {self.risk}")
        if len(self.component) > 100:
            raise ValidationError(f"Component name too long: {len(self.component)} chars")
        if len(self.risk) > 500:
            raise ValidationError(f"Risk description too long: {len(self.risk)} chars")
    
    def to_dict(self) -> Dict[str, str]:
        """Convert to dictionary."""
        return {"component": self.component, "risk": self.risk}

@dataclass
class RiskNoteEntry:
    """Represents a single risk note entry."""
    component: str
    risk: str
    severity: ValidationSeverity = ValidationSeverity.MEDIUM
    details: Optional[str] = None
    
    def __post_init__(self) -> None:
        """Validate fields after initialization."""
        self._validate()
    
    def _validate(self) -> None:
        """Internal validation method."""
        if not self.component or not isinstance(self.component, str):
            raise ValidationError(f"Invalid component: {self.component}")
        if not self.risk or not isinstance(self.risk, str):
            raise ValidationError(f"Invalid risk: {self.risk}")
        if len(self.component) > 100:
            raise ValidationError(f"Component name too long: {len(self.component)} chars")
        if len(self.risk) > 500:
            raise ValidationError(f"Risk description too long: {len(self.risk)} chars")
    
    def to_dict(self) -> Dict[str, Any]:
        """Convert to dictionary."""
        result = {"component": self.component, "risk": self.risk}
        if self.details:
            result["details"] = self.details
        return result

@dataclass
class Metadata:
    """Metadata for analysis output."""
    generated_at: str = field(default_factory=lambda: datetime.now(timezone.utc).isoformat())
    source: str = "unknown"
    version: str = "1.0.0"
    pipeline_stage: str = "delivery"
    correlation_id: Optional[str] = None
    
    def __post_init__(self) -> None:
        """Validate metadata fields."""
        self._validate()
    
    def _validate(self) -> None:
        """Internal validation method."""
        if self.pipeline_stage not in VALID_PIPELINE_STAGES:
            raise ValidationError(f"Invalid pipeline stage: {self.pipeline_stage}")
        if not self.source or not isinstance(self.source, str):
            raise ValidationError(f"Invalid source: {self.source}")
        if not self.version or not isinstance(self.version, str):
            raise ValidationError(f"Invalid version: {self.version}")
    
    def to_dict(self) -> Dict[str, Any]:
        """Convert to dictionary."""
        return {
            "generated_at": self.generated_at,
            "source": self.source,
            "version": self.version,
            "pipeline_stage": self.pipeline_stage,
            "correlation_id": self.correlation_id
        }

@dataclass
class AnalysisData:
    """Container for validated analysis output data."""
    
    metadata: Metadata = field(default_factory=Metadata)
    scope_map: List[ScopeMapEntry] = field(default_factory=list)
    risk_notes: List[RiskNoteEntry] = field(default_factory=list)
    attack_surface: Dict[str, Any] = field(default_factory=dict)
    _hash: Optional[str] = field(default=None, repr=False)
    
    def __post_init__(self) -> None:
        """Validate after initialization."""
        self._validate()
        self._compute_hash()
    
    def _validate(self) -> None:
        """Internal validation method."""
        if not isinstance(self.scope_map, list):
            raise ValidationError("scope_map must be a list")
        if not isinstance(self.risk_notes, list):
            raise ValidationError("risk_notes must be a list")
        if not isinstance(self.attack_surface, dict):
            raise ValidationError("attack_surface must be a dict")
        
        # Validate entries
        for entry in self.scope_map:
            if not isinstance(entry, ScopeMapEntry):
                raise ValidationError(f"Invalid scope map entry: {entry}")
        for entry in self.risk_notes:
            if not isinstance(entry, RiskNoteEntry):
                raise ValidationError(f"Invalid risk note entry: {entry}")
    
    def _compute_hash(self) -> None:
        """Compute hash of data for integrity verification."""
        data_str = json.dumps(self.to_dict(), sort_keys=True, default=str)
        self._hash = hashlib.sha256(data_str.encode()).hexdigest()
    
    def verify_integrity(self) -> bool:
        """Verify data integrity using hash."""
        old_hash = self._hash
        self._compute_hash()
        return self._hash == old_hash
    
    @classmethod
    def from_raw(cls, raw: Dict[str, Any]) -> 'AnalysisData':
        """
        Create AnalysisData from raw dictionary input.
        
        Args:
            raw: Raw analysis data dictionary
            
        Returns:
            Validated AnalysisData instance
            
        Raises:
            ValidationError: If input data is invalid
            DataParsingError: If parsing fails
        """
        if not isinstance(raw, dict):
            raise ValidationError(f"Input must be a dictionary, got {type(raw)}")
        
        try:
            # Parse metadata
            metadata_raw = raw.get("metadata", {})
            if not isinstance(metadata_raw, dict):
                raise ValidationError("metadata must be a dictionary")
            metadata = Metadata(**metadata_raw)
            
            # Parse scope map
            scope_map_raw = raw.get("scope_map", [])
            if not isinstance(scope_map_raw, list):
                raise ValidationError("scope_map must be a list")
            scope_map = []
            for entry in scope_map_raw:
                if isinstance(entry, dict):
                    scope_map.append(ScopeMapEntry(**entry))
                elif isinstance(entry, ScopeMapEntry):
                    scope_map.append(entry)
                else:
                    raise ValidationError(f"Invalid scope map entry type: {type(entry)}")
            
            # Parse risk notes
            risk_notes_raw = raw.get("risk_notes", [])
            if not isinstance(risk_notes_raw, list):
                raise ValidationError("risk_notes must be a list")
            risk_notes = []
            for entry in risk_notes_raw:
                if isinstance(entry, dict):
                    risk_notes.append(RiskNoteEntry(**entry))
                elif isinstance(entry, RiskNoteEntry):
                    risk_notes.append(entry)
                else:
                    raise ValidationError(f"Invalid risk note entry type: {type(entry)}")
            
            # Parse attack surface
            attack_surface = raw.get("attack_surface", {})
            if not isinstance(attack_surface, dict):
                raise ValidationError("attack_surface must be a dictionary")
            
            return cls(
                metadata=metadata,
                scope_map=scope_map,
                risk_notes=risk_notes,
                attack_surface=attack_surface
            )
            
        except TypeError as e:
            raise DataParsingError(f"Failed to parse data: {e}")
        except KeyError as e:
            raise DataParsingError(f"Missing required field: {e}")
    
    def to_dict(self) -> Dict[str, Any]:
        """Convert to dictionary."""
        return {
            "metadata": self.metadata.to_dict(),
            "scope_map": [entry.to_dict() for entry in self.scope_map],
            "risk_notes": [entry.to_dict() for entry in self.risk_notes],
            "attack_surface": self.attack_surface,
            "hash": self._hash
        }

#