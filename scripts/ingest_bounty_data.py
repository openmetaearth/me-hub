#!/usr/bin/env python3
"""
scripts/ingest_bounty_data.py

Ingest raw bug bounty data from GitHub issues and CSV templates.
Part of the Meta Earth Phase I data pipeline for security analysis.

Data Flow: ingestion -> validation -> transformation -> analysis -> delivery

Security Classification: INTERNAL - Contains vulnerability assessment data
"""

import csv
import json
import logging
import logging.handlers
import os
import re
import sys
import time
from dataclasses import dataclass, field, asdict
from datetime import datetime, timezone
from pathlib import Path
from typing import Any, Dict, List, Optional, Tuple, Generator, Set, Union
from urllib.parse import urlparse

import requests
from requests.adapters import HTTPAdapter
from requests.exceptions import (
    RequestException,
    HTTPError,
    ConnectionError,
    Timeout,
    RetryError,
)
from urllib3.util.retry import Retry

# ---------------------------------------------------------------------------
# Configuration & Constants
# ---------------------------------------------------------------------------

# Risk scoring for high-value trust boundaries
SCOPE_MAP: Dict[str, str] = {
    "Bridge": "Fund theft",
    "Sequencer logic": "Fake state submissions",
    "Governance": "Unauthorized control",
    "Reward logic": "Infinite rewards",
    "DID/KYC modules": "Identity bypass",
    "Validator logic": "Consensus manipulation",
    "RollApp settlement": "Settlement manipulation",
    "Cross-chain bridging": "Replay / forged messages",
    "Token minting": "Unlimited mint",
    "Gas logic": "Fee manipulation",
}

# Attack patterns to detect in descriptions
ATTACK_PATTERNS: Dict[str, List[str]] = {
    "replay": ["replay", "double-spend", "double spend", "nonce reuse"],
    "forged": ["forged", "fake proof", "invalid proof", "spoof"],
    "access_control": ["access control", "unauthorized", "privilege escalation", "admin"],
    "accounting": ["accounting", "overflow", "underflow", "rounding", "precision"],
    "governance": ["vote inflation", "snapshot", "proposal replay", "quorum bypass"],
    "bridge": ["bridge", "cross-chain", "relayer", "message passing"],
}

# Validation constants
VALID_SEVERITIES: Set[str] = {"critical", "high", "medium", "low", "unknown"}
VALID_STATUSES: Set[str] = {"open", "closed", "pending", "resolved", "wontfix"}
MAX_DESCRIPTION_LENGTH: int = 100000
MAX_TITLE_LENGTH: int = 500
MAX_LABELS: int = 50
MAX_RETRIES: int = 3
REQUEST_TIMEOUT: int = 30
PAGE_SIZE: int = 100

# File paths
DEFAULT_OUTPUT_DIR: Path = Path("output")
DEFAULT_LOG_DIR: Path = Path("logs")

# ---------------------------------------------------------------------------
# Custom Exceptions
# ---------------------------------------------------------------------------

class BountyIngestionError(Exception):
    """Base exception for bounty ingestion errors."""
    pass

class GitHubAPIError(BountyIngestionError):
    """Raised when GitHub API returns an error."""
    def __init__(self, message: str, status_code: Optional[int] = None, response: Optional[requests.Response] = None):
        self.status_code = status_code
        self.response = response
        super().__init__(message)

class ValidationError(BountyIngestionError):
    """Raised when data validation fails."""
    def __init__(self, message: str, errors: List[str]):
        self.errors = errors
        super().__init__(f"{message}: {'; '.join(errors)}")

class ConfigurationError(BountyIngestionError):
    """Raised when configuration is invalid."""
    pass

# ---------------------------------------------------------------------------
# Logging Setup
# ---------------------------------------------------------------------------

def setup_logging(
    verbose: bool = False,
    log_file: Optional[Path] = None,
    component: str = "bounty_ingest"
) -> logging.Logger:
    """
    Configure structured logging for the ingestion pipeline.
    
    Args:
        verbose: Enable debug logging
        log_file: Optional file path for log output
        component: Logger component name
        
    Returns:
        Configured logger instance
        
    Raises:
        ConfigurationError: If log directory cannot be created
    """
    logger = logging.getLogger(component)
    logger.setLevel(logging.DEBUG if verbose else logging.INFO)
    
    # Prevent duplicate handlers
    if logger.handlers:
        return logger
    
    formatter = logging.Formatter(
        "%(asctime)s | %(levelname)-8s | %(name)s | %(filename)s:%(lineno)d | %(message)s",
        datefmt="%Y-%m-%d %H:%M:%S",
    )
    
    # Console handler
    console_handler = logging.StreamHandler(sys.stdout)
    console_handler.setFormatter(formatter)
    logger.addHandler(console_handler)
    
    # File handler (if specified)
    if log_file:
        try:
            log_file.parent.mkdir(parents=True, exist_ok=True)
            file_handler = logging.handlers.RotatingFileHandler(
                log_file,
                maxBytes=10_000_000,  # 10MB
                backupCount=5,
            )
            file_handler.setFormatter(formatter)
            logger.addHandler(file_handler)
        except OSError as e:
            raise ConfigurationError(f"Cannot create log directory: {e}") from e
    
    return logger

logger = setup_logging()

# ---------------------------------------------------------------------------
# Data Models
# ---------------------------------------------------------------------------

@dataclass
class BountyIssue:
    """
    Represents a single bug bounty issue from GitHub or CSV.
    
    Attributes:
        source: Data source ("github" or "csv")
        source_id: Unique identifier from source
        title: Issue title
        description: Issue description/body
        component: Affected component
        risk: Risk assessment
        severity: Issue severity level
        status: Current status
        created_at: Creation timestamp
        updated_at: Last update timestamp
        labels: Issue labels
        raw_data: Original source data
    """
    source: str
    source_id: str
    title: str
    description: str
    component: str
    risk: str
    severity: str = "unknown"
    status: str = "open"
    created_at: str = ""
    updated_at: str = ""
    labels: List[str] = field(default_factory=list)
    raw_data: Dict[str, Any] = field(default_factory=dict)
    
    def __post_init__(self) -> None:
        """Validate and sanitize fields after initialization."""
        self.title = self.title.strip()[:MAX_TITLE_LENGTH]
        self.description = self.description.strip()[:MAX_DESCRIPTION_LENGTH]
        self.labels = self.labels[:MAX_LABELS]
        self.severity = self.severity.lower()
        self.status = self.status.lower()
    
    def to_dict(self) -> Dict[str, Any]:
        """Convert to dictionary, excluding raw data by default."""
        result = asdict(self)
        result.pop("raw_data", None)
        return result
    
    def to_json(self, include_raw: bool = False) -> str:
        """Serialize to JSON string."""
        data = asdict(self)
        if not include_raw:
            data.pop("raw_data", None)
        return json.dumps(data, default=str, indent=2)
    
    def validate(self) -> Tuple[bool, List[str]]:
        """
        Validate required fields and basic constraints.
        
        Returns:
            Tuple of (is_valid, list_of_errors)
        """
        errors: List[str] = []
        
        if not self.title:
            errors.append("Missing title")
        elif len(self.title) > MAX_TITLE_LENGTH:
            errors.append(f"Title exceeds {MAX_TITLE_LENGTH} characters")
            
        if not self.description:
            errors.append("Missing description")
        elif len(self.description) > MAX_DESCRIPTION_LENGTH:
            errors.append(f"Description exceeds {MAX_DESCRIPTION_LENGTH} characters")
            
        if not self.component:
            errors.append("Missing component")
        elif self.component not in SCOPE_MAP:
            errors.append(f"Unknown component: {self.component}")
            
        if not self.source:
            errors.append("Missing source")
        elif self.source not in {"github", "csv"}:
            errors.append(f"Invalid source: {self.source}")
            
        if not self.source_id:
            errors.append("Missing source_id")
            
        if self.severity not in VALID_SEVERITIES:
            errors.append(f"Invalid severity: {self.severity}")
            
        if self.status not in VALID_STATUSES:
            errors.append(f"Invalid status: {self.status}")
            
        return (len(errors) == 0, errors)

# ---------------------------------------------------------------------------
# HTTP Client with Retry Logic
# ---------------------------------------------------------------------------

def create_http_session(
    retries: int = MAX_RETRIES,
    backoff_factor: float = 0.5,
    timeout: int = REQUEST_TIMEOUT
) -> requests.Session:
    """
    Create a requests session with retry logic and connection pooling.
    
    Args:
        retries: Number of retry attempts
        backoff_factor: Exponential backoff factor
        timeout: Request timeout in seconds
        
    Returns:
        Configured requests session
    """
    session = requests.Session()
    
    retry_strategy = Retry(
        total=retries,
        backoff_factor=backoff_factor,
        status_forcelist=[429, 500, 502, 503, 504],
        allowed_methods=["GET", "POST", "PUT", "DELETE"],
    )
    
    adapter = HTTPAdapter(
        max_retries=retry_strategy,
        pool_connections=10,
        pool_maxsize=20,
    )
    
    session.mount("https://", adapter)
    session.mount("http://", adapter)
    session.timeout = timeout
    
    return session

# ---------------------------------------------------------------------------
# GitHub API Client
# ---------------------------------------------------------------------------

class GitHubClient:
    """
    Client for interacting with GitHub Issues API.
    
    Handles authentication, pagination, rate limiting, and error handling.
    """
    
    def __init__(
        self,
        token: Optional[str] = None,
        base_url: str = "https://api.github.com",
        session: Optional[requests.Session] = None,
    ):
        """
        Initialize GitHub client.
        
        Args:
            token: GitHub personal access token
            base_url: GitHub API base URL
            session: Optional pre-configured requests session
        """
        self.base_url = base_url.rstrip("/")
        self.session = session or create_http_session()
        self.headers: Dict[str, str] = {
            "Accept": "application/vnd.github.v3+json",
            "User-Agent": "MetaEarth-Bounty-Ingestion/1.0",
        }
        
        if token:
            self.headers["Authorization"] = f"token {token}"
            logger.debug("GitHub authentication configured")
        else:
            logger.warning("No GitHub token provided - rate limits will be restricted")
    
    def fetch_issues(
        self,
        repo: str,
        state: str = "open",
        labels: Optional[List[str]] = None,
        since: Optional[str] = None,
        max_pages: int = 10,
    ) -> Generator[Dict[str, Any], None, None]:
        """
        Fetch issues from a GitHub repository with pagination.
        
        Args:
            repo: Repository in format "owner/repo"
            state: Issue state (open, closed, all)
            labels: Filter by labels
            since: ISO 8601 timestamp to fetch issues updated after
            max_pages: Maximum number of pages to fetch
            
        Yields:
            Issue data dictionaries
            
        Raises:
            GitHubAPIError: On API errors
            ValueError: On invalid parameters
        """
        if not repo or "/" not in repo:
            raise ValueError(f"Invalid repository format: {repo}")
        
        params: Dict[str, Any] = {
            "state": state,
            "per_page": min(PAGE_SIZE, 100),
            "page": 1,
        }
        
        if labels:
            params["labels"] = ",".join(labels)
        
        if since:
            params["since"] = since
        
        url = f"{self.base_url}/repos/{repo}/issues"
        pages_fetched = 0
        
        while pages_fetched < max_pages:
            try:
                logger.debug(f"Fetching issues page {params['page']} from {repo}")
                
                response = self.session.get(
                    url,
                    headers=self.headers,
                    params=params,
                    timeout=REQUEST_TIMEOUT,
                )
                
                # Handle rate limiting
                if response.status_code == 403 and "rate limit" in response.text.lower():
                    retry_after = int(response.headers.get("Retry-After", 60))
                    logger.warning(f"Rate limited. Waiting {retry_after} seconds...")
                    time.sleep(retry_after)
                    continue
                
                response.raise_for_status()
                
                issues = response.json()
                if not issues:
                    break
                
                for issue in issues:
                    # Skip pull requests
                    if "pull_request" in issue:
                        continue
                    yield issue
                
                # Check for next page
                if "next" not in response.links:
                    break
                
                params["page"] += 1
                pages_fetched += 1
                
            except HTTPError as e:
                raise GitHubAPIError(
                    f"GitHub API error: {e}",
                    status_code=e.response.status_code if e.response else None,
                    response=e.response,
                )
            except (ConnectionError, Timeout) as e:
                logger.error(f"Network error fetching issues: {e}")
                raise GitHubAPIError(f"Network error: {e}")
            except RequestException as e:
                logger.error(f"Request error fetching issues: {e}")
                raise GitHubAPIError(f"Request error: {e}")
    
    def fetch_issue_detail(self, repo: str, issue_number: int) -> Dict[str, Any]:
        """
        Fetch a single issue by number.
        
        Args:
            repo: Repository in format "owner/repo"
            issue_number: Issue number
            
        Returns:
            Issue data dictionary
            
        Raises:
            GitHubAPIError: On API errors
        """
        url = f"{self.base_url}/repos/{repo}/issues/{issue_number}"
        
        try:
            response = self.session.get(url, headers=self.headers, timeout=REQUEST_TIMEOUT)
            response.raise_for_status()
            return response.json()
        except HTTPError as e:
            raise GitHubAPIError(
                f"Failed to fetch issue #{issue_number}: {e}",
                status_code=e.response.status_code if e.response else None,
                response=e.response,
            )
        except RequestException as e:
            raise GitHubAPIError(f"Request error: {e}")

# ---------------------------------------------------------------------------
# CSV Parser
# ---------------------------------------------------------------------------

class CSVParser:
    """
    Parser for CSV bounty data with validation and error handling.
    """
    
    REQUIRED_COLUMNS: Set[str] = {"title", "description", "component", "source_id"}
    OPTIONAL_COLUMNS: Set[str] = {"severity", "status", "labels", "created_at", "updated_at"}
    
    def __init__(self, file_path: Path, delimiter: str = ","):
        """
        Initialize CSV parser.
        
        Args:
            file_path: Path to CSV file
            delimiter: CSV delimiter character
            
        Raises:
            FileNotFoundError: If file doesn't exist
            ValueError: If file is not a valid CSV
        """
        if not file_path.exists():
            raise FileNotFoundError(f"CSV file not found: {file_path}")
        
        if file_path.suffix.lower() not in {".csv", ".tsv"}:
            raise ValueError(f"Invalid file extension: {file_path.suffix}")
        
        self.file_path = file_path
        self.delimiter = delimiter if file_path.suffix.lower() == ".csv" else "\t"
    
    def parse(self) -> Generator[Dict[str, str], None, None]:
        """
        Parse CSV file and yield rows.
        
        Yields:
            Dictionary of column name to value
            
        Raises:
            ValidationError: On missing required columns or invalid data
        """
        try:
            with open(self.file_path, "r", encoding="utf-8-sig") as f:
                reader = csv.DictReader(f, delimiter=self.delimiter)
                
                if not reader.fieldnames:
                    raise ValidationError("Empty CSV file", ["No columns found"])
                
                # Validate required columns
                missing_columns = self.REQUIRED_COLUMNS - set(reader.fieldnames)
                if missing_columns:
                    raise ValidationError(
                        "Missing required columns",
                        [f"Missing: {', '.join(missing_columns)}"]
                    )
                
                # Validate no unknown columns (warn only)
                unknown_columns = set(reader.fieldnames) - self.REQUIRED_COLUMNS - self.OPTIONAL_COLUMNS
                if unknown_columns:
                    logger.warning(f"Unknown columns in CSV: {', '.join(unknown_columns)}")
                
                row_num = 0
                for row_num, row in enumerate(reader, start=2):  # Start at 2 for header
                    try:
                        # Validate required fields are not empty
                        errors = []
                        for col in self.REQUIRED_COLUMNS:
                            if not row.get(col, "").strip():
                                errors.append(f"Row {row_num}: Missing {col}")
                        
                        if errors:
                            logger.warning(f"Validation errors in CSV: {'; '.join(errors)}")
                            continue
                        
                        yield row
                        
                    except Exception as e:
                        logger.error(f"Error processing row {row_num}: {e}")
                        continue
                
                if row_num == 0:
                    logger.warning("CSV file contains no data rows")
                    
        except csv.Error as e:
            raise ValidationError(f"CSV parsing error: {e}", [str(e)])
        except IOError as e:
            raise ValidationError(f"File read error: {e}", [str(e)])

# ---------------------------------------------------------------------------
# Data Transformer
# ---------------------------------------------------------------------------

class BountyTransformer:
    """
    Transform raw data from various sources into standardized BountyIssue objects.
    """
    
    @staticmethod
    def from_github_issue(issue: Dict[str, Any]) -> BountyIssue:
        """
        Transform a GitHub issue API response into a BountyIssue.
        
        Args:
            issue: GitHub issue data dictionary
            
        Returns:
            BountyIssue instance
        """
        # Extract component from