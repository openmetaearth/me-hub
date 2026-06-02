#!/usr/bin/env python3
"""
scripts/transform_data.py

Normalize and structure data for scope mapping and risk analysis
in the Meta Earth bug bounty pipeline.

This module transforms validated raw data into structured scope maps,
risk notes, and attack surface summaries for high-value trust boundaries
including bridges, sequencers, governance, reward logic, DID/KYC, and validators.

Security Classification: INTERNAL
Version: 2.1.0
"""

import json
import csv
import logging
import os
from pathlib import Path
from typing import Any, Dict, List, Optional, Union, Set, Tuple, Iterator, Generator
from dataclasses import dataclass, field, asdict
from datetime import datetime, timezone
from enum import Enum, auto
from hashlib import sha256, sha3_256
from functools import lru_cache, wraps
from contextlib import contextmanager
import re
import sys
import time
import traceback
from collections import defaultdict, OrderedDict
from threading import Lock
import mmap
import io

# Configure logging with structured format and rotation
logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(name)s - %(levelname)s - %(message)s - %(filename)s:%(lineno)d - %(funcName)s',
    handlers=[
        logging.StreamHandler(sys.stdout),
        logging.handlers.RotatingFileHandler(
            'transform_data.log',
            maxBytes=10*1024*1024,  # 10MB
            backupCount=5
        )
    ]
)
logger = logging.getLogger(__name__)

# Performance monitoring decorator
def monitor_performance(func):
    """Decorator to monitor function performance and log metrics."""
    @wraps(func)
    def wrapper(*args, **kwargs):
        start_time = time.perf_counter()
        try:
            result = func(*args, **kwargs)
            elapsed = time.perf_counter() - start_time
            if elapsed > 1.0:  # Log slow operations
                logger.warning(f"Slow operation detected: {func.__name__} took {elapsed:.3f}s")
            return result
        except Exception as e:
            elapsed = time.perf_counter() - start_time
            logger.error(f"Function {func.__name__} failed after {elapsed:.3f}s: {str(e)}")
            raise
    return wrapper


# ---------------------------------------------------------------------------
# Enums and Constants
# ---------------------------------------------------------------------------

class SeverityLevel(str, Enum):
    """Severity levels for risk assessment."""
    CRITICAL = "Critical"
    HIGH = "High"
    MEDIUM = "Medium"
    LOW = "Low"
    INFO = "Info"

    @classmethod
    def from_string(cls, value: str) -> 'SeverityLevel':
        """Create SeverityLevel from string with validation."""
        try:
            return cls(value.capitalize())
        except ValueError:
            raise ValidationError(f"Invalid severity level: {value}")


class ComponentCategory(str, Enum):
    """Categories for scope components."""
    CROSS_CHAIN = "Cross-chain"
    CONSENSUS_EXECUTION = "Consensus/Execution"
    PROTOCOL_MANAGEMENT = "Protocol Management"
    ECONOMICS = "Economics"
    IDENTITY = "Identity"
    VALIDATOR = "Validator"

    @classmethod
    def from_string(cls, value: str) -> 'ComponentCategory':
        """Create ComponentCategory from string with validation."""
        try:
            return cls(value)
        except ValueError:
            raise ValidationError(f"Invalid component category: {value}")


class AttackLikelihood(str, Enum):
    """Likelihood ratings for attack vectors."""
    VERY_HIGH = "Very High"
    HIGH = "High"
    MEDIUM = "Medium"
    LOW = "Low"
    VERY_LOW = "Very Low"

    @classmethod
    def from_string(cls, value: str) -> 'AttackLikelihood':
        """Create AttackLikelihood from string with validation."""
        try:
            return cls(value.title())
        except ValueError:
            raise ValidationError(f"Invalid attack likelihood: {value}")


class AttackImpact(str, Enum):
    """Impact ratings for attack vectors."""
    CRITICAL = "Critical"
    HIGH = "High"
    MEDIUM = "Medium"
    LOW = "Low"
    NEGLIGIBLE = "Negligible"

    @classmethod
    def from_string(cls, value: str) -> 'AttackImpact':
        """Create AttackImpact from string with validation."""
        try:
            return cls(value.capitalize())
        except ValueError:
            raise ValidationError(f"Invalid attack impact: {value}")


class DataFormat(str, Enum):
    """Supported data formats for input/output."""
    JSON = "json"
    CSV = "csv"
    YAML = "yaml"
    TOML = "toml"


class ValidationMode(str, Enum):
    """Validation modes for data processing."""
    STRICT = "strict"
    LENIENT = "lenient"
    SILENT = "silent"


# ---------------------------------------------------------------------------
# Custom Exceptions
# ---------------------------------------------------------------------------

class TransformationError(Exception):
    """Base exception for transformation errors."""
    def __init__(self, message: str, error_code: Optional[str] = None):
        self.error_code = error_code or "TRANSFORM_ERROR"
        super().__init__(message)


class ValidationError(TransformationError):
    """Exception raised for validation failures."""
    def __init__(self, message: str, field: Optional[str] = None):
        self.field = field
        super().__init__(message, error_code="VALIDATION_ERROR")


class ConfigurationError(TransformationError):
    """Exception raised for configuration issues."""
    def __init__(self, message: str, config_key: Optional[str] = None):
        self.config_key = config_key
        super().__init__(message, error_code="CONFIG_ERROR")


class DataIntegrityError(TransformationError):
    """Exception raised for data integrity violations."""
    def __init__(self, message: str, expected_hash: Optional[str] = None, actual_hash: Optional[str] = None):
        self.expected_hash = expected_hash
        self.actual_hash = actual_hash
        super().__init__(message, error_code="INTEGRITY_ERROR")


class ResourceExhaustionError(TransformationError):
    """Exception raised when system resources are exhausted."""
    def __init__(self, message: str, resource_type: str):
        self.resource_type = resource_type
        super().__init__(message, error_code="RESOURCE_EXHAUSTION")


# ---------------------------------------------------------------------------
# Data Models
# ---------------------------------------------------------------------------

@dataclass(frozen=True)
class RiskNote:
    """Represents a risk note for a specific component.
    
    Attributes:
        component: The component name
        risk: Description of the risk
        severity: Severity level of the risk
        description: Detailed description of the risk
        attack_vectors: List of potential attack vectors
        mitigation: Recommended mitigation strategies
        created_at: ISO format timestamp of creation
        risk_id: Unique identifier for the risk note
        version: Version of the risk note
        tags: Tags for categorization
    """
    component: str
    risk: str
    severity: SeverityLevel = SeverityLevel.HIGH
    description: str = ""
    attack_vectors: Tuple[str, ...] = field(default_factory=tuple)
    mitigation: str = ""
    created_at: str = field(default_factory=lambda: datetime.now(timezone.utc).isoformat())
    risk_id: str = field(default_factory=lambda: sha3_256(str(datetime.now(timezone.utc).timestamp()).encode()).hexdigest()[:16])
    version: str = "1.0.0"
    tags: Tuple[str, ...] = field(default_factory=tuple)

    def __post_init__(self) -> None:
        """Validate risk note data after initialization."""
        if not self.component or not self.component.strip():
            raise ValidationError("Component name cannot be empty", field="component")
        if not self.risk or not self.risk.strip():
            raise ValidationError("Risk description cannot be empty", field="risk")
        if not isinstance(self.severity, SeverityLevel):
            raise ValidationError(f"Invalid severity level: {self.severity}", field="severity")
        if not re.match(r'^\d+\.\d+\.\d+$', self.version):
            raise ValidationError(f"Invalid version format: {self.version}", field="version")

    def to_dict(self) -> Dict[str, Any]:
        """Convert to dictionary with proper serialization."""
        return {
            'component': self.component,
            'risk': self.risk,
            'severity': self.severity.value,
            'description': self.description,
            'attack_vectors': list(self.attack_vectors),
            'mitigation': self.mitigation,
            'created_at': self.created_at,
            'risk_id': self.risk_id,
            'version': self.version,
            'tags': list(self.tags)
        }

    @classmethod
    def from_dict(cls, data: Dict[str, Any]) -> 'RiskNote':
        """Create RiskNote from dictionary with validation."""
        try:
            return cls(
                component=data['component'],
                risk=data['risk'],
                severity=SeverityLevel.from_string(data.get('severity', 'High')),
                description=data.get('description', ''),
                attack_vectors=tuple(data.get('attack_vectors', [])),
                mitigation=data.get('mitigation', ''),
                created_at=data.get('created_at', datetime.now(timezone.utc).isoformat()),
                risk_id=data.get('risk_id', ''),
                version=data.get('version', '1.0.0'),
                tags=tuple(data.get('tags', []))
            )
        except KeyError as e:
            raise ValidationError(f"Missing required field: {e}", field=str(e))


@dataclass(frozen=True)
class ScopeComponent:
    """Represents a component in the scope map.
    
    Attributes:
        name: Component name
        category: Component category
        risk_level: Risk level assessment
        description: Component description
        sub_components: List of sub-components
        trust_boundary: Whether this is a trust boundary
        component_id: Unique identifier for the component
        dependencies: List of component dependencies
        version: Version of the component
    """
    name: str
    category: ComponentCategory
    risk_level: SeverityLevel
    description: str = ""
    sub_components: Tuple[str, ...] = field(default_factory=tuple)
    trust_boundary: bool = True
    component_id: str = field(default_factory=lambda: sha3_256(str(datetime.now(timezone.utc).timestamp()).encode()).hexdigest()[:16])
    dependencies: Tuple[str, ...] = field(default_factory=tuple)
    version: str = "1.0.0"

    def __post_init__(self) -> None:
        """Validate scope component data after initialization."""
        if not self.name or not self.name.strip():
            raise ValidationError("Component name cannot be empty", field="name")
        if not isinstance(self.category, ComponentCategory):
            raise ValidationError(f"Invalid category: {self.category}", field="category")
        if not isinstance(self.risk_level, SeverityLevel):
            raise ValidationError(f"Invalid risk level: {self.risk_level}", field="risk_level")
        if not re.match(r'^\d+\.\d+\.\d+$', self.version):
            raise ValidationError(f"Invalid version format: {self.version}", field="version")

    def to_dict(self) -> Dict[str, Any]:
        """Convert to dictionary with proper serialization."""
        return {
            'name': self.name,
            'category': self.category.value,
            'risk_level': self.risk_level.value,
            'description': self.description,
            'sub_components': list(self.sub_components),
            'trust_boundary': self.trust_boundary,
            'component_id': self.component_id,
            'dependencies': list(self.dependencies),
            'version': self.version
        }

    @classmethod
    def from_dict(cls, data: Dict[str, Any]) -> 'ScopeComponent':
        """Create ScopeComponent from dictionary with validation."""
        try:
            return cls(
                name=data['name'],
                category=ComponentCategory.from_string(data['category']),
                risk_level=SeverityLevel.from_string(data.get('risk_level', 'High')),
                description=data.get('description', ''),
                sub_components=tuple(data.get('sub_components', [])),
                trust_boundary=data.get('trust_boundary', True),
                component_id=data.get('component_id', ''),
                dependencies=tuple(data.get('dependencies', [])),
                version=data.get('version', '1.0.0')
            )
        except KeyError as e:
            raise ValidationError(f"Missing required field: {e}", field=str(e))


@dataclass(frozen=True)
class AttackSurface:
    """Represents an attack surface entry.
    
    Attributes:
        component: Affected component
        attack_type: Type of attack
        likelihood: Likelihood of attack
        impact: Potential impact
        examples: Example attack scenarios
        notes: Additional notes
        surface_id: Unique identifier for the attack surface
        cwe_id: CWE identifier if applicable
        remediation: Remediation steps
    """
    component: str
    attack_type: str
    likelihood: AttackLikelihood
    impact: AttackImpact
    examples: Tuple[str, ...] = field(default_factory=tuple)
    notes: str = ""
    surface_id: str = field(default_factory=lambda: sha3_256(str(datetime.now(timezone.utc).timestamp()).encode()).hexdigest()[:16])
    cwe_id: Optional[str] = None
    remediation: str = ""

    def __post_init__(self) -> None:
        """Validate attack surface data after initialization."""
        if not self.component or not self.component.strip():
            raise ValidationError("Component name cannot be empty", field="component")
        if not self.attack_type or not self.attack_type.strip():
            raise ValidationError("Attack type cannot be empty", field="attack_type")
        if not isinstance(self.likelihood, AttackLikelihood):
            raise ValidationError(f"Invalid likelihood: {self.likelihood}", field="likelihood")
        if not isinstance(self.impact, AttackImpact):
            raise ValidationError(f"Invalid impact: {self.impact}", field="impact")
        if self.cwe_id and not re.match(r'^CWE-\d+$', self.cwe_id):
            raise ValidationError(f"Invalid CWE ID format: {self.cwe_id}", field="cwe_id")

    def to_dict(self) -> Dict[str, Any]:
        """Convert to dictionary with proper serialization."""
        return {
            'component': self.component,
            'attack_type': self.attack_type,
            'likelihood': self.likelihood.value,
            'impact': self.impact.value,
            'examples': list(self.examples),
            'notes': self.notes,
            'surface_id': self.surface_id,
            'cwe_id': self.cwe_id,
            'remediation': self.remediation
        }

    @classmethod
    def from_dict(cls, data: Dict[str, Any]) -> 'AttackSurface':
        """Create AttackSurface from dictionary with validation."""
        try:
            return cls(
                component=data['component'],
                attack_type=data['attack_type'],
                likelihood=AttackLikelihood.from_string(data.get('likelihood', 'Medium')),
                impact=AttackImpact.from_string(data.get('impact', 'Medium')),
                examples=tuple(data.get('examples', [])),
                notes=data.get('notes', ''),
                surface_id=data.get('surface_id', ''),
                cwe_id=data.get('cwe_id'),
                remediation=data.get('remediation', '')
            )
        except KeyError as e:
            raise ValidationError(f"Missing required field: {e}", field=str(e))


@dataclass
class TransformedData:
    """Container for all transformed data outputs.
    
    Attributes:
        scope_map: List of scope components
        risk_notes: List of risk notes
        attack_surfaces: List of attack surfaces
        metadata: Dictionary of metadata
        data_hash: Hash of the transformed data for integrity verification
        version: Version of the transformed data
        created_at: Creation timestamp
    """
    scope_map: List[ScopeComponent] = field(default_factory=list)
    risk_notes: List[RiskNote] = field(default_factory=list)
    attack_surfaces: List[AttackSurface] = field(default_factory=list)
    metadata: Dict[str, Any] = field(default_factory=dict)
    data_hash: str = field(default_factory=str)
    version: str = "2.1.0"
    created_at: str = field(default_factory=lambda: datetime.now(timezone.utc).isoformat())

    def __post_init__(self) -> None:
        """Initialize data hash after creation."""
        if not self.data_hash:
            self.data_hash = self.compute_hash()

    def compute_hash(self) -> str:
        """Compute hash of the transformed data for integrity verification."""
        data_string = json.dumps({
            'scope_map': [component.to_dict() for component in self.scope_map],
            'risk_notes': [note.to_dict() for note in self.risk_notes],
            'attack_surfaces': [surface.to_dict() for surface in self.attack_surfaces],
            'metadata': self.metadata,
            'version': self.version
        }, sort_keys=True)
        return sha3_256(data_string.encode()).hexdigest()

    def verify_integrity(self) -> bool:
        """Verify data integrity by comparing hashes."""
        current_hash = self.compute_hash()
        if current_hash != self.data_hash:
            raise DataIntegrityError(
                "Data integrity check failed",
                expected_hash=self.data_hash,
                actual_hash=current_hash
            )
        return True

    def to_dict(self) -> Dict[str, Any]:
        """Convert to dictionary with proper serialization."""
        return {
            'scope_map': [component.to_dict() for component in self.scope_map],
            'risk_notes': [note.to_dict() for note in self.risk_notes],
            'attack_surfaces': [surface.to_dict() for surface in self.attack_surfaces],
            'metadata': self.metadata,
            'data_hash': self.data_hash,
            'version': self.version,
            'created_at': self.created_at
        }

    @classmethod
    def from_dict(cls, data: Dict[str, Any]) -> 'TransformedData':
        """Create TransformedData from dictionary with validation."""
        try:
            instance = cls(
                scope_map=[ScopeComponent.from_dict(item) for item in data.get('scope_map', [])],
                risk_notes=[RiskNote.from_dict(item) for