#!/usr/bin/env python3
"""
scripts/analyze_attack_surface.py

Generate scope map, risk notes, and attack surface summary for Meta Earth Phase I.
This script analyzes blockchain infrastructure components and produces structured
outputs for security assessment planning.

Usage:
    python scripts/analyze_attack_surface.py [--output-dir <path>] [--log-level <level>]

Outputs:
    - scope_map.json: High-value trust boundaries and components
    - risk_notes.json: Component-risk mapping with attack vectors
    - attack_surface_summary.md: Human-readable markdown report
"""

import argparse
import json
import logging
import os
import sys
import time
from dataclasses import dataclass, field, asdict
from datetime import datetime, timezone
from pathlib import Path
from typing import Dict, List, Optional, Any, Union, Final, Set, Tuple
from enum import Enum, auto
from functools import lru_cache
from collections import OrderedDict
import hashlib
import traceback
from contextlib import contextmanager

# ---------------------------------------------------------------------------
# Constants
# ---------------------------------------------------------------------------

SCRIPT_VERSION: Final[str] = "1.0.0"
DEFAULT_OUTPUT_DIR: Final[str] = "output"
MAX_RETRY_ATTEMPTS: Final[int] = 3
RETRY_DELAY_SECONDS: Final[float] = 1.0
VALID_SEVERITY_LEVELS: Final[Set[str]] = {"critical", "high", "medium", "low"}
VALID_PRIORITY_LEVELS: Final[Set[str]] = {"critical", "high", "medium", "low"}
MAX_COMPONENT_NAME_LENGTH: Final[int] = 100
MAX_DESCRIPTION_LENGTH: Final[int] = 1000
MAX_ATTACK_VECTORS: Final[int] = 20
MAX_COMPONENTS: Final[int] = 50
OUTPUT_FILE_PERMISSIONS: Final[int] = 0o644
LOG_FILE_PERMISSIONS: Final[int] = 0o644
MAX_LOG_FILE_SIZE: Final[int] = 10 * 1024 * 1024  # 10 MB
MAX_LOG_BACKUP_COUNT: Final[int] = 5
ENCODING: Final[str] = "utf-8"

# ---------------------------------------------------------------------------
# Logging Configuration
# ---------------------------------------------------------------------------

class LogLevel(Enum):
    """Enum for log levels to ensure type safety."""
    DEBUG = logging.DEBUG
    INFO = logging.INFO
    WARNING = logging.WARNING
    ERROR = logging.ERROR
    CRITICAL = logging.CRITICAL

class CustomRotatingFileHandler(logging.handlers.RotatingFileHandler):
    """Custom file handler with size-based rotation and proper permissions."""
    
    def __init__(self, filename: str, max_bytes: int = MAX_LOG_FILE_SIZE, 
                 backup_count: int = MAX_LOG_BACKUP_COUNT):
        super().__init__(filename, maxBytes=max_bytes, backupCount=backup_count)
        self._set_file_permissions()
    
    def _set_file_permissions(self) -> None:
        """Set proper file permissions for log files."""
        try:
            if os.path.exists(self.baseFilename):
                os.chmod(self.baseFilename, LOG_FILE_PERMISSIONS)
        except OSError as e:
            logging.getLogger(__name__).warning(f"Could not set log file permissions: {e}")

def setup_logging(log_level: str = "INFO") -> logging.Logger:
    """
    Configure logging with proper levels, formatting, and rotation.
    
    Args:
        log_level: String representation of log level (DEBUG, INFO, WARNING, ERROR, CRITICAL)
    
    Returns:
        Configured logger instance
    
    Raises:
        ValueError: If invalid log level is provided
        OSError: If log directory cannot be created
    """
    try:
        level = LogLevel[log_level.upper()].value
    except KeyError:
        raise ValueError(f"Invalid log level: {log_level}. Must be one of {[e.name for e in LogLevel]}")
    
    logger = logging.getLogger("analyze_attack_surface")
    logger.setLevel(level)
    
    # Clear existing handlers to avoid duplication
    logger.handlers.clear()
    
    # Console handler with proper formatting
    console_handler = logging.StreamHandler(sys.stdout)
    console_handler.setLevel(level)
    console_formatter = logging.Formatter(
        "%(asctime)s - %(name)s - %(levelname)s - %(message)s",
        datefmt="%Y-%m-%d %H:%M:%S"
    )
    console_handler.setFormatter(console_formatter)
    logger.addHandler(console_handler)
    
    # File handler with rotation for persistent logging
    try:
        log_dir = Path("logs")
        log_dir.mkdir(exist_ok=True, mode=0o755)
        
        log_file = log_dir / f"attack_surface_{datetime.now(timezone.utc).strftime('%Y%m%d_%H%M%S')}.log"
        file_handler = CustomRotatingFileHandler(
            str(log_file),
            max_bytes=MAX_LOG_FILE_SIZE,
            backup_count=MAX_LOG_BACKUP_COUNT
        )
        file_handler.setLevel(logging.DEBUG)  # Always log debug to file
        file_formatter = logging.Formatter(
            "%(asctime)s - %(name)s - %(levelname)s - %(filename)s:%(lineno)d - %(funcName)s - %(message)s"
        )
        file_handler.setFormatter(file_formatter)
        logger.addHandler(file_handler)
        
        logger.info(f"Logging to file: {log_file}")
    except OSError as e:
        logger.warning(f"Could not create file handler: {e}")
        logger.warning("Continuing with console logging only")
    
    return logger

logger = setup_logging()

# ---------------------------------------------------------------------------
# Custom Exceptions
# ---------------------------------------------------------------------------

class AttackSurfaceError(Exception):
    """Base exception for attack surface analysis errors."""
    def __init__(self, message: str, error_code: Optional[str] = None):
        super().__init__(message)
        self.error_code = error_code
        self.timestamp = datetime.now(timezone.utc)

class ComponentValidationError(AttackSurfaceError):
    """Raised when component validation fails."""
    def __init__(self, message: str, component_name: Optional[str] = None):
        super().__init__(message, error_code="COMP_VAL_ERR")
        self.component_name = component_name

class FileOperationError(AttackSurfaceError):
    """Raised when file operations fail."""
    def __init__(self, message: str, file_path: Optional[str] = None):
        super().__init__(message, error_code="FILE_OP_ERR")
        self.file_path = file_path

class DataSerializationError(AttackSurfaceError):
    """Raised when data serialization/deserialization fails."""
    def __init__(self, message: str, data_type: Optional[str] = None):
        super().__init__(message, error_code="DATA_SER_ERR")
        self.data_type = data_type

class ConfigurationError(AttackSurfaceError):
    """Raised when configuration is invalid."""
    def __init__(self, message: str, config_key: Optional[str] = None):
        super().__init__(message, error_code="CONFIG_ERR")
        self.config_key = config_key

# ---------------------------------------------------------------------------
# Data Models
# ---------------------------------------------------------------------------

@dataclass(frozen=True)
class AttackVector:
    """
    Represents a specific attack vector for a component.
    
    Attributes:
        name: Attack vector name
        description: Detailed description of the attack
        severity: Severity level (critical, high, medium, low)
        likelihood: Likelihood of exploitation (high, medium, low)
        mitigation: Recommended mitigation strategies
    """
    name: str
    description: str
    severity: str = "high"
    likelihood: str = "medium"
    mitigation: str = ""
    
    def __post_init__(self) -> None:
        """Validate attack vector data after initialization."""
        self._validate()
    
    def _validate(self) -> None:
        """
        Validate attack vector fields.
        
        Raises:
            ComponentValidationError: If validation fails
        """
        errors: List[str] = []
        
        if not self.name or not self.name.strip():
            errors.append("Attack vector name cannot be empty")
        elif len(self.name) > MAX_COMPONENT_NAME_LENGTH:
            errors.append(f"Attack vector name exceeds maximum length of {MAX_COMPONENT_NAME_LENGTH}")
        
        if not self.description or not self.description.strip():
            errors.append("Attack vector description cannot be empty")
        elif len(self.description) > MAX_DESCRIPTION_LENGTH:
            errors.append(f"Attack vector description exceeds maximum length of {MAX_DESCRIPTION_LENGTH}")
        
        if self.severity not in VALID_SEVERITY_LEVELS:
            errors.append(f"Invalid severity: {self.severity}. Must be one of {VALID_SEVERITY_LEVELS}")
        
        if self.likelihood not in {"high", "medium", "low"}:
            errors.append(f"Invalid likelihood: {self.likelihood}. Must be one of high, medium, low")
        
        if errors:
            raise ComponentValidationError(
                f"Attack vector validation failed: {'; '.join(errors)}",
                component_name=self.name
            )


@dataclass(frozen=True)
class Component:
    """
    Represents a blockchain infrastructure component.
    
    Attributes:
        name: Component name (must be unique and non-empty)
        description: Detailed component description
        risk: Primary risk associated with the component
        attack_vectors: List of potential attack vectors
        trust_boundary: Whether this component is a trust boundary
        priority: Priority level for security assessment
        category: Component category for organization
        dependencies: List of component dependencies
    """
    name: str
    description: str
    risk: str
    attack_vectors: List[AttackVector] = field(default_factory=list)
    trust_boundary: bool = True
    priority: str = "high"
    category: str = "infrastructure"
    dependencies: List[str] = field(default_factory=list)
    
    def __post_init__(self) -> None:
        """Validate component data after initialization."""
        self._validate()
    
    def _validate(self) -> None:
        """
        Validate component fields.
        
        Raises:
            ComponentValidationError: If validation fails
        """
        errors: List[str] = []
        
        if not self.name or not self.name.strip():
            errors.append("Component name cannot be empty")
        elif len(self.name) > MAX_COMPONENT_NAME_LENGTH:
            errors.append(f"Component name exceeds maximum length of {MAX_COMPONENT_NAME_LENGTH}")
        
        if not self.description or not self.description.strip():
            errors.append("Component description cannot be empty")
        elif len(self.description) > MAX_DESCRIPTION_LENGTH:
            errors.append(f"Component description exceeds maximum length of {MAX_DESCRIPTION_LENGTH}")
        
        if not self.risk or not self.risk.strip():
            errors.append("Component risk cannot be empty")
        
        if self.priority not in VALID_PRIORITY_LEVELS:
            errors.append(f"Invalid priority: {self.priority}. Must be one of {VALID_PRIORITY_LEVELS}")
        
        if len(self.attack_vectors) > MAX_ATTACK_VECTORS:
            errors.append(f"Too many attack vectors: {len(self.attack_vectors)}. Maximum is {MAX_ATTACK_VECTORS}")
        
        if errors:
            raise ComponentValidationError(
                f"Component validation failed: {'; '.join(errors)}",
                component_name=self.name
            )
    
    def add_attack_vector(self, vector: AttackVector) -> 'Component':
        """
        Add an attack vector to the component (returns new instance).
        
        Args:
            vector: Attack vector to add
        
        Returns:
            New Component instance with added attack vector
        
        Raises:
            ComponentValidationError: If maximum attack vectors exceeded
        """
        if len(self.attack_vectors) >= MAX_ATTACK_VECTORS:
            raise ComponentValidationError(
                f"Maximum attack vectors ({MAX_ATTACK_VECTORS}) exceeded for component '{self.name}'",
                component_name=self.name
            )
        
        new_vectors = list(self.attack_vectors) + [vector]
        return Component(
            name=self.name,
            description=self.description,
            risk=self.risk,
            attack_vectors=new_vectors,
            trust_boundary=self.trust_boundary,
            priority=self.priority,
            category=self.category,
            dependencies=self.dependencies
        )


@dataclass(frozen=True)
class RiskNote:
    """
    Represents a risk note for a component-risk pair.
    
    Attributes:
        component: Component name
        risk: Risk description
        severity: Severity level
        attack_vectors: Specific attack vectors
        mitigation_strategies: Recommended mitigations
        references: External references or links
    """
    component: str
    risk: str
    severity: str
    attack_vectors: List[str] = field(default_factory=list)
    mitigation_strategies: List[str] = field(default_factory=list)
    references: List[str] = field(default_factory=list)
    
    def __post_init__(self) -> None:
        """Validate risk note data after initialization."""
        self._validate()
    
    def _validate(self) -> None:
        """
        Validate risk note fields.
        
        Raises:
            ComponentValidationError: If validation fails
        """
        errors: List[str] = []
        
        if not self.component or not self.component.strip():
            errors.append("Component name cannot be empty")
        
        if not self.risk or not self.risk.strip():
            errors.append("Risk description cannot be empty")
        
        if self.severity not in VALID_SEVERITY_LEVELS:
            errors.append(f"Invalid severity: {self.severity}. Must be one of {VALID_SEVERITY_LEVELS}")
        
        if errors:
            raise ComponentValidationError(
                f"Risk note validation failed: {'; '.join(errors)}",
                component_name=self.component
            )


@dataclass(frozen=True)
class ScopeMap:
    """
    High-value trust boundaries and components for Phase I.
    
    Attributes:
        components: Dictionary of components by name
        generated_at: ISO 8601 timestamp of generation
        version: Schema version
        checksum: SHA-256 checksum of the data
    """
    components: Dict[str, Component] = field(default_factory=dict)
    generated_at: str = field(default_factory=lambda: datetime.now(timezone.utc).isoformat())
    version: str = SCRIPT_VERSION
    checksum: str = ""
    
    def __post_init__(self) -> None:
        """Generate checksum after initialization."""
        if not self.checksum:
            object.__setattr__(self, 'checksum', self._generate_checksum())
    
    def _generate_checksum(self) -> str:
        """Generate SHA-256 checksum of the scope map data."""
        data = json.dumps(self.to_dict(), sort_keys=True, default=str)
        return hashlib.sha256(data.encode(ENCODING)).hexdigest()
    
    def add_component(self, component: Component) -> 'ScopeMap':
        """
        Add a component to the scope map (returns new instance).
        
        Args:
            component: Component to add
        
        Returns:
            New ScopeMap instance with added component
        
        Raises:
            ComponentValidationError: If component with same name already exists
            ComponentValidationError: If maximum components exceeded
        """
        if component.name in self.components:
            raise ComponentValidationError(
                f"Component '{component.name}' already exists in scope map",
                component_name=component.name
            )
        
        if len(self.components) >= MAX_COMPONENTS:
            raise ComponentValidationError(
                f"Maximum components ({MAX_COMPONENTS}) exceeded",
                component_name=component.name
            )
        
        new_components = dict(self.components)
        new_components[component.name] = component
        
        return ScopeMap(
            components=new_components,
            generated_at=self.generated_at,
            version=self.version
        )
    
    def get_component(self, name: str) -> Optional[Component]:
        """
        Get a component by name.
        
        Args:
            name: Component name to look up
        
        Returns:
            Component if found, None otherwise
        """
        return self.components.get(name)
    
    def get_components_by_priority(self, priority: str) -> List[Component]:
        """
        Get all components with a specific priority.
        
        Args:
            priority: Priority level to filter by
        
        Returns:
            List of components with the specified priority
        """
        if priority not in VALID_PRIORITY_LEVELS:
            raise ValueError(f"Invalid priority: {priority}")
        return [comp for comp in self.components.values() if comp.priority == priority]
    
    def get_components_by_category(self, category: str) -> List[Component]:
        """
        Get all components in a specific category.
        
        Args:
            category: Category to filter by
        
        Returns:
            List of components in the specified category
        """
        return [comp for comp in self.components.values() if comp.category == category]
    
    def to_dict(self) -> Dict[str, Any]:
        """
        Convert scope map to dictionary for serialization.
        
        Returns:
            Dictionary representation of scope map
        """
        return {
            "generated_at": self.generated_at,
            "version": self.version,
            "checksum": self.checksum,
            "components": {
                name: {
                    "name": comp.name,
                    "description": comp.description,
                    "risk": comp.risk,
                    "attack_vectors": [
                        {
                            "name": av.name,
                            "description": av.description,
                            "severity": av.severity,
                            "likelihood": av.likelihood,
                            "mitigation": av.mitigation
                        }
                        for av in comp.attack_vectors
                    ],
                    "trust_boundary": comp.trust_boundary,
                    "priority": comp.priority,
                    "category": comp.category,
                    "dependencies": comp.dependencies
                }
                for name, comp in self.components.items()
            }
        }
    
    @classmethod
    def from_dict(cls, data: Dict[str, Any]) -> 'ScopeMap':
        """
        Create ScopeMap from dictionary.
        
        Args:
            data: Dictionary representation of scope map
        
        Returns:
            ScopeMap instance
        
        Raises:
            DataSerializationError: If data is invalid
        """
        try:
            components = {}
            for name, comp_data in data.get("components", {}).items