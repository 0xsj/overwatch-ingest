# platform/pylib/scout_common/observability/tracing/metadata.py
"""gRPC metadata helpers for request ID propagation."""

from typing import Optional, Tuple

import grpc

from .request_id import RequestID, REQUEST_ID_HEADER
from .context import get_request_id, set_request_id


def inject_request_id(
    metadata: Optional[list] = None,
    request_id: Optional[RequestID] = None
) -> list:
    """Inject request ID into gRPC outgoing metadata.
    
    Used by gRPC clients to propagate request IDs to downstream services.
    
    Args:
        metadata: Existing metadata list (or None for new)
        request_id: Request ID to inject (or None to use from context)
        
    Returns:
        Metadata list with request ID injected
    """
    if metadata is None:
        metadata = []
    
    if request_id is None:
        request_id = get_request_id()
    
    if request_id is None:
        return metadata
    
    # Remove existing request ID if present
    metadata = [
        (k, v) for k, v in metadata 
        if k.lower() != REQUEST_ID_HEADER.lower()
    ]
    
    # Add request ID
    metadata.append((REQUEST_ID_HEADER, str(request_id)))
    return metadata


def extract_request_id(context: grpc.ServicerContext) -> Optional[RequestID]:
    """Extract request ID from gRPC incoming metadata.
    
    Used by gRPC servers to receive request IDs from upstream services.
    
    Args:
        context: gRPC servicer context
        
    Returns:
        Request ID if found, None otherwise
    """
    metadata = dict(context.invocation_metadata())
    
    # Try case-insensitive lookup
    for key in [REQUEST_ID_HEADER, REQUEST_ID_HEADER.lower()]:
        value = metadata.get(key)
        if value:
            return RequestID.parse(value)
    
    return None


def extract_or_generate_request_id(
    context: grpc.ServicerContext
) -> RequestID:
    """Extract request ID from metadata or generate a new one.
    
    Also sets it in the context for downstream use.
    
    Args:
        context: gRPC servicer context
        
    Returns:
        Request ID (extracted or newly generated)
    """
    request_id = extract_request_id(context)
    if request_id is None:
        request_id = RequestID.generate()
    
    set_request_id(request_id)
    return request_id


def propagate_request_id(
    metadata: Optional[list] = None
) -> list:
    """Propagate request ID from context to outgoing metadata.
    
    Useful for client interceptors to automatically propagate request IDs.
    
    Args:
        metadata: Existing metadata list (or None for new)
        
    Returns:
        Metadata list with request ID propagated
    """
    request_id = get_request_id()
    if request_id is None:
        request_id = RequestID.generate()
        set_request_id(request_id)
    
    return inject_request_id(metadata, request_id)