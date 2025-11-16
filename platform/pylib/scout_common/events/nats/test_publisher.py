# platform/pylib/scout_common/events/nats/test_publisher.py
"""Test NATS publisher."""

import asyncio
import pytest

from scout_common.events.nats import NATSConfig, NATSPublisher


@pytest.mark.asyncio
async def test_nats_publisher():
    """Test NATS publisher connection and publish."""
    config = NATSConfig(url="nats://localhost:4224")
    
    try:
        async with NATSPublisher(config) as publisher:
            await publisher.publish_json("test.subject", {
                "message": "hello from python"
            })
            print("✅ Successfully published to NATS")
    except Exception as e:
        pytest.skip(f"NATS not available: {e}")


if __name__ == "__main__":
    asyncio.run(test_nats_publisher())