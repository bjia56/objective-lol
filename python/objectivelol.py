import asyncio
import concurrent.futures
import inspect
import json
import threading
import uuid

from .vm import *

definedFunctions = {}

# gopy does not support passing complex types directly,
# so we wrap arguments and return values as JSON strings.
# Additionally, using closures seems to result in a segfault
# at https://github.com/python/cpython/blob/v3.13.5/Python/generated_cases.c.h#L2462
# so we use a global dictionary to store the actual functions.
def gopy_wrapper(id: str, json_args: str):
    args = json.loads(json_args)
    try:
        result = definedFunctions[id](*args)
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

    async def declare_variable_async(self, name: str, value, constant: bool = False) -> None:
        self.declare_variable(name, value, constant)

    def declare_function(self, name: str, function, argc: int = None) -> None:
        argc = argc is None and len(inspect.signature(function).parameters) or argc
        unique_id = str(uuid.uuid4())
        definedFunctions[unique_id] = function
        self.vm.DefineFunctionMaxCompat(unique_id, name, argc, gopy_wrapper)

    async def declare_function_async(self, name: str, function) -> None:
        self.declare_function(name, function)

    async def declare_coroutine_async(self, name: str, function) -> None:
        loop = asyncio.get_event_loop()
        argc = len(inspect.signature(function).parameters)

        def wrapper(*args):
            fut = concurrent.futures.Future()
            def do():
                try:
                    result = asyncio.run_coroutine_threadsafe(function(*args), loop).result()
                    fut.set_result(result)
                except Exception as e:
                    fut.set_exception(e)
            threading.Thread(target=do).start()
            return fut.result()

        self.declare_function(name, wrapper, argc)

    def execute(self, code: str) -> None:
        return self.vm.Execute(code)

    async def execute_async(self, code: str) -> None:
        fut = concurrent.futures.Future()
        def do():
            try:
                result = self.vm.Execute(code)
                fut.set_result(result)
            except Exception as e:
                fut.set_exception(e)
        threading.Thread(target=do).start()
        return await asyncio.wrap_future(fut)