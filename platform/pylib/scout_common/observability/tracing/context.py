# platform/pylib/scout_common/observability/tracing/context.py
"""Context storage for request IDs using contextvars."""

from contextvars import ContextVar
from typing import Optional

from .request_id import RequestID


# Thread-safe context variable for request ID
_request_id_var: ContextVar[Optional[RequestID]] = ContextVar(
    "request_id", default=None
)


def set_request_id(request_id: RequestID) -> None:
    """Set the request ID in the current context.
    
    Args:
        request_id: Request ID to set
    """
    _request_id_var.set(request_id)


def get_request_id() -> Optional[RequestID]:
    """Get the request ID from the current context.
    
    Returns:
        Request ID if set, None otherwise
    """
    return _request_id_var.get()


def get_request_id_or_generate() -> RequestID:
    """Get the request ID from context or generate a new one.
    
    Returns:
        Request ID from context or newly generated
    """
    rid = get_request_id()
    if rid is None:
        rid = RequestID.generate()
        set_request_id(rid)
    return rid


def has_request_id() -> bool:
    """Check if a request ID is set in the current context.
    
    Returns:
        True if request ID is set, False otherwise
    """
    return get_request_id() is not None


def clear_request_id() -> None:
    """Clear the request ID from the current context."""
    _request_id_var.set(None)


def with_request_id(request_id: RequestID):
    """Context manager to temporarily set a request ID.
    
    Args:
        request_id: Request ID to set
        
    Example:
        with with_request_id(RequestID.generate()):
            # Code here has request ID in context
            pass
    """
    token = _request_id_var.set(request_id)
    try:
        yield
    finally:
        _request_id_var.reset(token)


def with_new_request_id():
    """Context manager to temporarily set a new request ID.
    
    Example:
        with with_new_request_id():
            # Code here has new request ID in context
            pass
    """
    return with_request_id(RequestID.generate())