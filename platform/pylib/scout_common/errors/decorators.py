"""Decorators for error handling in Scout platform."""

import functools
import asyncio
from typing import TypeVar, Callable, ParamSpec, Awaitable, Any
from collections.abc import Coroutine

from .result import Result, Ok, Err
from .base import Error
from .constructors import (
    internal,
    timeout as timeout_error,
    validation,
    wrap,
)

from .types import ErrorType
from .codes import code

