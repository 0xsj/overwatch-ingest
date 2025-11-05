"""Scout Tools - ML tools service entry point."""

import asyncio
import logging
import signal
import sys
from typing import NoReturn

logging.basicConfig(
    level=logging.INFO,
    format="%(asctime)s - %(name)s - %(levelname)s - %(message)s"
)

logger = logging.getLogger(__name__)


async def process_tick() -> None:
    """Simulate processing tick."""
    logger.info("Tools service: Processing tick...")
    # TODO: Load ML models, process requests, etc.


async def run_worker() -> NoReturn:
    """Run the worker with periodic ticks."""
    logger.info("Tools worker running, processing every 5 seconds...")
    
    while True:
        try:
            await process_tick()
            await asyncio.sleep(5)
        except asyncio.CancelledError:
            logger.info("Tools worker stopped")
            break
        except Exception as e:
            logger.error(f"Error in worker: {e}", exc_info=True)
            await asyncio.sleep(5)


async def main() -> None:
    """Main entry point."""
    logger.info("Tools service starting...")
    
    # TODO: Initialize dependencies (config, logger, ML models, etc.)
    # TODO: Set up gRPC server
    # TODO: Register service handlers
    
    # Handle graceful shutdown
    loop = asyncio.get_running_loop()
    stop_event = asyncio.Event()
    
    def signal_handler() -> None:
        logger.info("Shutdown signal received, stopping service...")
        stop_event.set()
    
    for sig in (signal.SIGINT, signal.SIGTERM):
        loop.add_signal_handler(sig, signal_handler)
    
    # Run worker
    worker_task = asyncio.create_task(run_worker())
    
    # Wait for shutdown signal
    await stop_event.wait()
    
    # Cancel worker
    worker_task.cancel()
    await worker_task


if __name__ == "__main__":
    try:
        asyncio.run(main())
    except KeyboardInterrupt:
        logger.info("Service interrupted")
        sys.exit(0)