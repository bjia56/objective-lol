from inspect import signature
from functools import partial
import json

from .vm import *

definedFunctions = {}

# gopy does not support passing complex types directly,
# so we wrap arguments and return values as JSON strings.
# Additionally, using closures seems to result in a segfault
# at https://github.com/python/cpython/blob/v3.13.5/Python/generated_cases.c.h#L2462
# so we use a global dictionary to store the actual functions.
def gopy_wrapper(name: str, json_args: str):
    args = json.loads(json_args)
    try:
        result = definedFunctions[name](*args)
        return json.dumps({"result": result, "error": None}).encode('utf-8')
    except Exception as e:
        return json.dumps({"result": None, "error": str(e)}).encode('utf-8')

class ObjectiveLOLVM:
    def __init__(self):
        # todo: figure out how to bridge stdout/stdin
        self.vm = NewVM(DefaultConfig())

    def declare_variable(self, name: str, value, constant: bool = False) -> None:
        if isinstance(value, int):
            var_type = "INTEGR"
            goValue = WrapInt(value)
        elif isinstance(value, float):
            var_type = "DUBBLE"
            goValue = WrapFloat(value)
        elif isinstance(value, str):
            var_type = "STRIN"
            goValue = WrapString(value)
        elif isinstance(value, bool):
            var_type = "BOOL"
            goValue = WrapBool(value)
        else:
            raise ValueError("Unsupported variable type %s" % type(value))

        self.vm.DefineVariable(name, var_type, goValue, constant)

    def declare_function(self, name: str, function) -> None:
        argc = len(signature(function).parameters)
        definedFunctions[name] = function
        self.vm.DefineFunctionMaxCompat(name, argc, gopy_wrapper)

    def execute(self, code: str) -> None:
        return self.vm.Execute(code)

