# platform/pylib/scout_common/events/message.py
"""Message wrapper for event bus messages."""

import json
from datetime import datetime
from typing import Any, Dict, TypeVar, Type
from dataclasses import dataclass, field


T = TypeVar("T")


@dataclass
class Message:
    """Represents an event message from the event bus."""

    subject: str
    """The topic/subject the message was published to."""

    data: bytes
    """The raw message payload."""

    metadata: Dict[str, str] = field(default_factory=dict)
    """Additional message information."""

    timestamp: datetime = field(default_factory=datetime.utcnow)
    """When the message was published."""

    def unmarshal_json(self, event_type: Type[T]) -> T:
        """Unmarshal the message data into a dataclass or pydantic model.
        
        Args:
            event_type: The type to deserialize into
            
        Returns:
            Deserialized event object
            
        Raises:
            json.JSONDecodeError: If data is not valid JSON
            TypeError: If event_type cannot be constructed from dict
        """
        data_dict = json.loads(self.data)
        
        # Support both dataclasses and pydantic models
        if hasattr(event_type, "model_validate"):
            # Pydantic v2
            return event_type.model_validate(data_dict)
        elif hasattr(event_type, "parse_obj"):
            # Pydantic v1
            return event_type.parse_obj(data_dict)
        else:
            # Dataclass or regular class
            return event_type(**data_dict)

    def to_dict(self) -> Dict[str, Any]:
        """Convert message data to dictionary.
        
        Returns:
            Dictionary representation of the message data
            
        Raises:
            json.JSONDecodeError: If data is not valid JSON
        """
        return json.loads(self.data)