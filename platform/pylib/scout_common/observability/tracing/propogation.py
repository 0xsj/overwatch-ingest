# platform/pylib/scout_common/observability/tracing/propagation.py
"""Propagation helpers for request IDs across different transports."""

from typing import Dict, Optional, Callable, Awaitable

from .request_id import RequestID, REQUEST_ID_HEADER
from .context import get_request_id, set_request_id


def inject_http(headers: Dict[str, str], request_id: Optional[RequestID] = None) -> None:
    """Inject request ID into HTTP headers.
    
    Used by HTTP clients to propagate request IDs.
    
    Args:
        headers: HTTP headers dictionary
        request_id: Request ID to inject (or None to use from context)
    """
    if request_id is None:
        request_id = get_request_id()
    
    if request_id is None:
        return
    
    headers[REQUEST_ID_HEADER] = str(request_id)


def extract_http(headers: Dict[str, str]) -> Optional[RequestID]:
    """Extract request ID from HTTP headers.
    
    Used by HTTP servers to receive request IDs.
    
    Args:
        headers: HTTP headers dictionary
        
    Returns:
        Request ID if found, None otherwise
    """
    # Try case-insensitive lookup
    for key in [REQUEST_ID_HEADER, REQUEST_ID_HEADER.lower()]:
        value = headers.get(key)
        if value:
            return RequestID.parse(value)
    
    return None


def inject_event_metadata(
    metadata: Dict[str, str],
    request_id: Optional[RequestID] = None
) -> None:
    """Inject request ID into event metadata.
    
    Used when publishing events to NATS/RabbitMQ.
    
    Args:
        metadata: Event metadata dictionary
        request_id: Request ID to inject (or None to use from context)
    """
    if request_id is None:
        request_id = get_request_id()
    
    if request_id is None:
        return
    
    metadata[REQUEST_ID_HEADER] = str(request_id)


def extract_event_metadata(metadata: Dict[str, str]) -> Optional[RequestID]:
    """Extract request ID from event metadata.
    
    Used when consuming events from NATS/RabbitMQ.
    
    Args:
        metadata: Event metadata dictionary
        
    Returns:
        Request ID if found, None otherwise
    """
    value = metadata.get(REQUEST_ID_HEADER)
    if value:
        return RequestID.parse(value)
    
    return None


def propagate_to_event(metadata: Dict[str, str]) -> RequestID:
    """Extract request ID from context and inject into event metadata.
    
    Args:
        metadata: Event metadata dictionary
        
    Returns:
        Request ID (from context or newly generated)
    """
    request_id = get_request_id()
    if request_id is None:
        request_id = RequestID.generate()
        set_request_id(request_id)
    
    inject_event_metadata(metadata, request_id)
    return request_id


def extract_from_event(metadata: Dict[str, str]) -> RequestID:
    """Extract request ID from event metadata and set in context.
    
    Args:
        metadata: Event metadata dictionary
        
    Returns:
        Request ID (extracted or newly generated)
    """
    request_id = extract_event_metadata(metadata)
    if request_id is None:
        request_id = RequestID.generate()
    
    set_request_id(request_id)
    return request_id


# HTTP Middleware helper
async def http_middleware(
    request,
    call_next: Callable[[any], Awaitable[any]]
):
    """HTTP middleware for FastAPI/Starlette.
    
    Extracts or generates request ID and adds to context.
    
    Args:
        request: HTTP request
        call_next: Next middleware/handler
        
    Returns:
        HTTP response with request ID header
    """
    # Extract from headers or generate
    request_id = extract_http(dict(request.headers))
    if request_id is None:
        request_id = RequestID.generate()
    
    # Set in context
    set_request_id(request_id)
    
    # Process request
    response = await call_next(request)
    
    # Add to response headers
    response.headers[REQUEST_ID_HEADER] = str(request_id)
    
    return response