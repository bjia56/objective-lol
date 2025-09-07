import asyncio
import concurrent.futures
import inspect
import json
import threading
import uuid

from .vm import (
    VM,
    VMCompatibilityShim,
    NewVM,
    DefaultConfig,
    WrapInt,
    WrapFloat,
    WrapString,
    WrapBool,
    GoValue,
    Slice_api_GoValue,
    Map_string_api_GoValue,
)


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


def convert_to_go_value(value):
    if isinstance(value, int):
        return WrapInt(value)
    elif isinstance(value, float):
        return WrapFloat(value)
    elif isinstance(value, str):
        return WrapString(value)
    elif isinstance(value, bool):
        return WrapBool(value)
    elif isinstance(value, GoValue):
        # object handle, pass through
        return value
    elif isinstance(value, (list, tuple)):
        slice = Slice_api_GoValue()
        for v in value:
            slice.append(convert_to_go_value(v))
        return slice
    elif isinstance(value, dict):
        map = Map_string_api_GoValue()
        for k, v in value.items():
            map[k] = convert_to_go_value(v)
        return map
    else:
        raise ValueError("Unsupported argument type %s" % type(value))


def convert_from_go_value(go_value: GoValue):
    typ = go_value.Type()
    if typ == "INTEGR":
        return go_value.Int()
    elif typ == "DUBBLE":
        return go_value.Float()
    elif typ == "STRIN":
        return go_value.String()
    elif typ == "BOOL":
        return go_value.Bool()
    elif typ == "BUKKIT":
        return [convert_from_go_value(v) for v in go_value.Slice()]
    elif typ == "BASKIT":
        return {k: convert_from_go_value(v) for k, v in go_value.Map().items()}
    else:
        # object handle
        return go_value


class ObjectiveLOLVM:
    __vm: VM
    __compat: VMCompatibilityShim

    def __init__(self):
        # todo: figure out how to bridge stdout/stdin
        self.__vm = NewVM(DefaultConfig())
        self.__compat = self.__vm.GetCompatibilityShim()

    def define_variable(self, name: str, value, constant: bool = False) -> None:
        goValue = convert_to_go_value(value)
        self.__vm.DefineVariable(name, goValue, constant)

    async def define_variable_async(self, name: str, value, constant: bool = False) -> None:
        self.define_variable(name, value, constant)

    def define_function(self, name: str, function, argc: int = None) -> None:
        argc = argc is None and len(inspect.signature(function).parameters) or argc
        unique_id = str(uuid.uuid4())
        definedFunctions[unique_id] = function
        self.__compat.DefineFunction(unique_id, name, argc, gopy_wrapper)

    async def define_function_async(self, name: str, function) -> None:
        self.define_function(name, function)

    async def define_coroutine_async(self, name: str, function) -> None:
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

        self.define_function(name, wrapper, argc)

    def call(self, name: str, *args):
        goArgs = convert_to_go_value(args)
        result = self.__vm.Call(name, goArgs)
        return convert_from_go_value(result)

    async def call_async(self, name: str, *args):
        goArgs = convert_to_go_value(args)
        fut = concurrent.futures.Future()
        def do():
            try:
                result = self.__vm.Call(name, goArgs)
                fut.set_result(convert_from_go_value(result))
            except Exception as e:
                fut.set_exception(e)
        threading.Thread(target=do).start()
        return await asyncio.wrap_future(fut)

    def call_method(self, receiver: GoValue, name: str, *args):
        goArgs = convert_to_go_value(args)
        result = self.__vm.CallMethod(receiver, name, goArgs)
        return convert_from_go_value(result)

    async def call_method_async(self, receiver: GoValue, name: str, *args):
        goArgs = convert_to_go_value(args)
        fut = concurrent.futures.Future()
        def do():
            try:
                result = self.__vm.CallMethod(receiver, name, goArgs)
                fut.set_result(convert_from_go_value(result))
            except Exception as e:
                fut.set_exception(e)
        threading.Thread(target=do).start()
        return await asyncio.wrap_future(fut)

    def execute(self, code: str) -> None:
        return self.__vm.Execute(code)

    async def execute_async(self, code: str) -> None:
        fut = concurrent.futures.Future()
        def do():
            try:
                result = self.__vm.Execute(code)
                fut.set_result(result)
            except Exception as e:
                fut.set_exception(e)
        threading.Thread(target=do).start()
        return await asyncio.wrap_future(fut)