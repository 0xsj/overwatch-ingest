# platform/pylib/scout_common/events/__init__.py
"""Event infrastructure for Scout platform.

Provides provider-agnostic event publishing and subscription.
Implementations: NATS, Kafka (future), RabbitMQ (future)

Example usage:
    from scout_common.events.nats import NATSConfig, NATSPublisher, NATSSubscriber
    
    # Publisher
    config = NATSConfig(url="nats://localhost:4224")
    async with NATSPublisher(config) as publisher:
        await publisher.publish_json("agent.created", {
            "agent_id": "123",
            "name": "my-agent"
        })
    
    # Subscriber
    async with NATSSubscriber(config) as subscriber:
        async def handle_event(msg):
            event = msg.to_dict()
            print(f"Received: {event}")
        
        await subscriber.subscribe("agent.created", handle_event)
"""

from .publisher import Publisher
from .subscriber import Subscriber, Subscription, MessageHandler
from .message import Message
from .config import Config

__all__ = [
    "Publisher",
    "Subscriber",
    "Subscription",
    "MessageHandler",
    "Message",
    "Config",
]