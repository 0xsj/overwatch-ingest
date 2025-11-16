# platform/pylib/scout_common/events/nats/__init__.py
"""NATS implementation of event infrastructure.

Example usage:
    from scout_common.events.nats import NATSConfig, NATSPublisher, NATSSubscriber
    
    config = NATSConfig(
        url="nats://localhost:4224",
        max_reconnects=10,
        reconnect_wait=2.0
    )
    
    # Publish events
    async with NATSPublisher(config) as publisher:
        await publisher.publish_json("my.subject", {"key": "value"})
    
    # Subscribe to events
    async with NATSSubscriber(config) as subscriber:
        async def handler(msg):
            data = msg.to_dict()
            print(f"Received: {data}")
        
        await subscriber.subscribe("my.subject", handler)
        await asyncio.Event().wait()  # Keep listening
"""

from .config import NATSConfig
from .publisher import NATSPublisher
from .subscriber import NATSSubscriber

__all__ = [
    "NATSConfig",
    "NATSPublisher",
    "NATSSubscriber",
]