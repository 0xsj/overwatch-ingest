"""Main entry point for Tools service."""

import signal
import sys

from scout_common.observability.logger import noop, new, Level
from scout_common.errors import Error

from config import load


def main() -> None:
    """Run the Tools service."""
    try:
        # Load configuration
        cfg = load(use_prefix=False)
        
        # Initialize logger
        if cfg.environment == "development":
            # Development: colorized console
            logger = noop(
                service="tools",
                environment=cfg.environment,
                version="0.1.0",
            )
        else:
            # Production: structured JSON logging
            level = Level.parse(cfg.log_level)
            logger = new(level=level, development=False)
            logger = logger.with_fields(
                service="tools",
                environment=cfg.environment,
                version="0.1.0",
            )
        
        # Log startup
        logger.info(
            "tools service starting",
            port=cfg.port,
            log_level=cfg.log_level,
            postgres_host=cfg.postgres.host,
            redis_host=cfg.redis.host,
            nats_url=cfg.nats.url,
        )
        
        # TODO: Initialize dependencies (database, cache, message bus)
        # TODO: Set up HTTP/gRPC server
        # TODO: Register handlers
        # TODO: Start server
        
        logger.info(
            "tools service running",
            port=cfg.port,
            address=f"http://localhost:{cfg.port}",
        )
        
        # Wait for interrupt signal
        def signal_handler(signum, frame):
            logger.info("shutdown signal received", signal=signal.Signals(signum).name)
            logger.info("tools service stopped")
            sys.exit(0)
        
        signal.signal(signal.SIGINT, signal_handler)
        signal.signal(signal.SIGTERM, signal_handler)
        
        # Keep the service running
        signal.pause()
        
    except Exception as e:
        # Catch all exceptions (not just Error)
        print(f"Failed to start service: {e}", file=sys.stderr)
        import traceback
        traceback.print_exc()
        sys.exit(1)


if __name__ == "__main__":
    main()