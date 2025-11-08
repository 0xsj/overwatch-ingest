"""Recovery interceptor for gRPC panic recovery."""

from typing import Callable, Any

import grpc

from scout_common.observability.logger import Logger


class RecoveryInterceptor(grpc.ServerInterceptor):
    """Recovers from exceptions in gRPC handlers."""
    
    def __init__(self, logger: Logger) -> None:
        """
        Initialize recovery interceptor.
        
        Args:
            logger: Logger instance
        """
        self._logger = logger
    
    def intercept_service(
        self,
        continuation: Callable,
        handler_call_details: grpc.HandlerCallDetails,
    ) -> grpc.RpcMethodHandler:
        """
        Intercept service calls to add recovery.
        
        Args:
            continuation: Function to invoke the next handler
            handler_call_details: Details about the RPC
            
        Returns:
            RPC method handler
        """
        method = handler_call_details.method
        
        def wrapper(behavior: Callable) -> Callable:
            def wrapped(request: Any, context: grpc.ServicerContext) -> Any:
                try:
                    return behavior(request, context)
                except Exception as e:
                    self._logger.error(
                        "grpc handler error recovered",
                        method=method,
                        error=str(e),
                        error_type=type(e).__name__,
                    )
                    
                    context.set_code(grpc.StatusCode.INTERNAL)
                    context.set_details("Internal server error")
                    raise
            
            return wrapped
        
        handler = continuation(handler_call_details)
        
        if handler and handler.unary_unary:
            return grpc.unary_unary_rpc_method_handler(
                wrapper(handler.unary_unary),
                request_deserializer=handler.request_deserializer,
                response_serializer=handler.response_serializer,
            )
        
        return handler