# Without __init__.py (verbose, ugly)
from scout_common.errors.types import ErrorType
from scout_common.errors.codes import Code, code
from scout_common.errors.base import Error, error
from scout_common.errors.result import Ok, Err, Result, unwrap, map_value, and_then
from scout_common.errors.constructors import not_found, validation, internal