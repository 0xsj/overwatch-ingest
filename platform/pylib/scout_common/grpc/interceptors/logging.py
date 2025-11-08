"""Logging interceptor for gRPC."""

import time
from typing import Callable, Any

import grpc

from scout_common.observability.logger import Logger


class LoggingInterceptor(grpc.ServerInterceptor):
    """Logs gRPC requests and responses."""
    
    def __init__(self, logger: Logger) -> None:
        """
        Initialize logging interceptor.
        
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
        Intercept service calls to add logging.
        
        Args:
            continuation: Function to invoke the next handler
            handler_call_details: Details about the RPC
            
        Returns:
            RPC method handler
        """
        method = handler_call_details.method
        
        def wrapper(behavior: Callable) -> Callable:
            def wrapped(request: Any, context: grpc.ServicerContext) -> Any:
                start_time = time.time()
                
                self._logger.info("grpc request started", method=method)
                
                try:
                    response = behavior(request, context)
                    
                    duration_ms = (time.time() - start_time) * 1000
                    
                    self._logger.info(
                        "grpc request completed",
                        method=method,
                        duration_ms=round(duration_ms, 2),
                    )
                    
                    return response
                    
                except Exception as e:
                    duration_ms = (time.time() - start_time) * 1000
                    
                    self._logger.error(
                        "grpc request failed",
                        method=method,
                        duration_ms=round(duration_ms, 2),
                        error=str(e),
                    )
                    
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