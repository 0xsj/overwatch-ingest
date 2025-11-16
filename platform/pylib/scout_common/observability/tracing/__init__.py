# platform/pylib/scout_common/observability/tracing/__init__.py
"""Request ID tracing for distributed systems.

Provides request ID generation, context storage, and propagation across:
- gRPC calls (via metadata)
- HTTP requests (via headers)
- Events (via metadata)

Example usage:
    from scout_common.observability.tracing import (
        RequestID,
        set_request_id,
        get_request_id,
        extract_or_generate_request_id,
    )
    
    # Generate and set
    request_id = RequestID.generate()
    set_request_id(request_id)
    
    # Get from context
    current_id = get_request_id()
    
    # gRPC server
    request_id = extract_or_generate_request_id(grpc_context)
    
    # Event publishing
    from scout_common.observability.tracing import propagate_to_event
    metadata = {}
    request_id = propagate_to_event(metadata)
    # metadata now has request ID
"""

from .request_id import RequestID, REQUEST_ID_HEADER, REQUEST_ID_PREFIX
from .context import (
    set_request_id,
    get_request_id,
    get_request_id_or_generate,
    has_request_id,
    clear_request_id,
    with_request_id,
    with_new_request_id,
)
from .metadata import (
    inject_request_id,
    extract_request_id,
    extract_or_generate_request_id,
    propagate_request_id,
)
from .propagation import (
    inject_http,
    extract_http,
    inject_event_metadata,
    extract_event_metadata,
    propagate_to_event,
    extract_from_event,
    http_middleware,
)

__all__ = [
    # Request ID
    "RequestID",
    "REQUEST_ID_HEADER",
    "REQUEST_ID_PREFIX",
    # Context
    "set_request_id",
    "get_request_id",
    "get_request_id_or_generate",
    "has_request_id",
    "clear_request_id",
    "with_request_id",
    "with_new_request_id",
    # gRPC Metadata
    "inject_request_id",
    "extract_request_id",
    "extract_or_generate_request_id",
    "propagate_request_id",
    # Propagation
    "inject_http",
    "extract_http",
    "inject_event_metadata",
    "extract_event_metadata",
    "propagate_to_event",
    "extract_from_event",
    "http_middleware",
]