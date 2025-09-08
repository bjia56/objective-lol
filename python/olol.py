import asyncio
import concurrent.futures
import functools
import inspect
import json
import threading
import uuid

from .api import (
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
    ClassDefinition,
    ClassVariable,
    ClassMethod,
    NewClassDefinition,
    GoValueIDKey
)


defined_functions = {}
object_instances = {}


def serialize_go_value(vm: VM, go_value: GoValue):
    if isinstance(go_value, GoValue):
        converted = convert_from_go_value(go_value)
        if isinstance(converted, GoValue):
            if converted.Type() == "NOTHIN":
                return None
            return {GoValueIDKey: converted.ID()}
        else:
            return converted
    else:
        return go_value


# gopy does not support passing complex types directly,
# so we wrap arguments and return values as JSON strings.
# Additionally, using closures seems to result in a segfault
# at https://github.com/python/cpython/blob/v3.13.5/Python/generated_cases.c.h#L2462
# so we use a global dictionary to store the actual functions.
def gopy_wrapper(id: str, json_args: str):
    args = json.loads(json_args)
    try:
        vm, fn = defined_functions[id]
        result = fn(*args)
        return json.dumps({"result": convert_to_go_value(vm, result), "error": None}, default=functools.partial(serialize_go_value, vm)).encode('utf-8')
    except Exception as e:
        return json.dumps({"result": None, "error": str(e)}).encode('utf-8')


def convert_to_go_value(vm: 'ObjectiveLOLVM', value):
    if value is None:
        return GoValue()
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
            slice.append(convert_to_go_value(vm, v))
        return slice
    elif isinstance(value, dict):
        map = Map_string_api_GoValue()
        for k, v in value.items():
            map[k] = convert_to_go_value(vm, v)
        return map
    else:
        vm.define_class(type(value))
        instance = vm._ObjectiveLOLVM__vm.NewObjectInstance(type(value).__name__)
        object_instances[instance.ID()] = value
        return instance


def convert_from_go_value(go_value: GoValue):
    if not isinstance(go_value, GoValue):
        return go_value
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
    __loop: asyncio.AbstractEventLoop

    def __init__(self):
        # todo: figure out how to bridge stdout/stdin
        self.__vm = NewVM(DefaultConfig())
        self.__compat = self.__vm.GetCompatibilityShim()
        self.__loop = asyncio.get_event_loop()

    def define_variable(self, name: str, value, constant: bool = False) -> None:
        goValue = convert_to_go_value(self, value)
        self.__vm.DefineVariable(name, goValue, constant)

    def define_function(self, name: str, function, argc: int = None) -> None:
        argc = argc is None and len(inspect.signature(function).parameters) or argc
        unique_id = str(uuid.uuid4())
        defined_functions[unique_id] = (self, function)
        self.__compat.DefineFunction(unique_id, name, argc, gopy_wrapper)

    def define_coroutine(self, name: str, function) -> None:
        argc = len(inspect.signature(function).parameters)

        def wrapper(*args):
            fut = concurrent.futures.Future()
            def do():
                try:
                    result = asyncio.run_coroutine_threadsafe(function(*args), self.__loop).result()
                    fut.set_result(result)
                except Exception as e:
                    fut.set_exception(e)
            threading.Thread(target=do).start()
            return fut.result()

        self.define_function(name, wrapper, argc)

    def define_class(self, python_class: type) -> None:
        class_name = python_class.__name__
        if self.__compat.IsClassDefined(class_name):
            pass

        # Use class builder to introspect and build the class definition
        builder = ClassBuilder(self)
        builder.set_name(class_name)
        builder.add_constructor(python_class)

        # Add class attributes as variables with getters/setters
        for attr_name in dir(python_class):
            if not attr_name.startswith('_') and not callable(getattr(python_class, attr_name)):
                builder.add_public_variable(
                    attr_name,
                    getter=lambda self: getattr(self, attr_name),
                    setter=lambda self, value: setattr(self, attr_name, value)
                )

        # Add methods
        for method_name in dir(python_class):
            if not method_name.startswith('_') and callable(getattr(python_class, method_name)):
                method = getattr(python_class, method_name)
                if not method_name == python_class.__name__:  # Skip constructor
                    if inspect.iscoroutinefunction(method):
                        builder.add_public_coroutine(method_name, method)
                    else:
                        builder.add_public_method(method_name, method)

        class_def = builder.get()
        self.__vm.DefineClass(class_def)

    def call(self, name: str, *args):
        goArgs = convert_to_go_value(self, args)
        result = self.__vm.Call(name, goArgs)
        return convert_from_go_value(result)

    async def call_async(self, name: str, *args):
        goArgs = convert_to_go_value(self, args)
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
        goArgs = convert_to_go_value(self, args)
        result = self.__vm.CallMethod(receiver, name, goArgs)
        return convert_from_go_value(result)

    async def call_method_async(self, receiver: GoValue, name: str, *args):
        goArgs = convert_to_go_value(self, args)
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


class ClassBuilder:
    __vm: ObjectiveLOLVM
    __compat: VMCompatibilityShim
    __loop: asyncio.AbstractEventLoop
    __class: ClassDefinition

    def __init__(self, vm: ObjectiveLOLVM):
        self.__vm = vm
        self.__compat = vm._ObjectiveLOLVM__compat
        self.__loop = vm._ObjectiveLOLVM__loop
        self.__class = NewClassDefinition()

    def get(self) -> ClassDefinition:
        return self.__class

    def set_name(self, name: str) -> 'ClassBuilder':
        self.__class.Name = name
        return self

    def __build_variable(self, name: str, value, locked: bool, getter=None, setter=None) -> ClassVariable:
        class_variable = ClassVariable()
        class_variable.Name = name
        class_variable.Value = convert_to_go_value(self.__vm, value)
        class_variable.Locked = locked
        if getter is not None:
            unique_id = str(uuid.uuid4())

            def wrapper(this_id):
                return convert_to_go_value(self.__vm, getter(object_instances[this_id]))

            defined_functions[unique_id] = (self.__vm, wrapper)
            self.__compat.BuildNewClassVariableWithGetter(class_variable, unique_id, gopy_wrapper)
        if setter is not None:
            unique_id = str(uuid.uuid4())

            def wrapper(this_id, value):
                setter(object_instances[this_id], convert_from_go_value(value))

            defined_functions[unique_id] = (self.__vm, wrapper)
            self.__compat.BuildNewClassVariableWithSetter(class_variable, unique_id, gopy_wrapper)
        return class_variable

    def add_public_variable(self, name: str, value = None, locked: bool = False, getter=None, setter=None) -> 'ClassBuilder':
        variable = self.__build_variable(name, value, locked, getter, setter)
        self.__class.PublicVariables[name] = variable
        return self

    def add_private_variable(self, name: str, value = None, locked: bool = False, getter=None, setter=None) -> 'ClassBuilder':
        variable = self.__build_variable(name, value, locked, getter, setter)
        self.__class.PrivateVariables[name] = variable
        return self

    def add_shared_variable(self, name: str, value = None, locked: bool = False, getter=None, setter=None) -> 'ClassBuilder':
        variable = self.__build_variable(name, value, locked, getter, setter)
        self.__class.SharedVariables[name] = variable
        return self

    def __build_method(self, name: str, function, argc: int = None) -> ClassMethod:
        argc = argc is None and len(inspect.signature(function).parameters) - 1 or argc
        unique_id = str(uuid.uuid4())

        def wrapper(this_id, *args):
            return convert_to_go_value(self.__vm, function(object_instances[this_id], *args))

        defined_functions[unique_id] = (self.__vm, wrapper)
        class_method = ClassMethod()
        class_method.Name = name
        class_method.Argc = argc

        self.__compat.BuildNewClassMethod(class_method, unique_id, gopy_wrapper)
        return class_method

    def add_constructor(self, typ) -> 'ClassBuilder':
        # get init function
        init_function = typ.__init__
        argc = len(inspect.signature(init_function).parameters) - 1

        # ignore args and kwargs
        for param in inspect.signature(init_function).parameters.values():
            if param.kind in (param.VAR_POSITIONAL, param.VAR_KEYWORD):
                argc = argc - 1

        unique_id = str(uuid.uuid4())

        def ctor_wrapper(this_id, *args):
            instance = typ(*args)
            object_instances[this_id] = instance

        defined_functions[unique_id] = (self.__vm, ctor_wrapper)
        class_method = ClassMethod()
        class_method.Name = typ.__name__
        class_method.Argc = argc

        self.__compat.BuildNewClassMethod(class_method, unique_id, gopy_wrapper)
        self.__class.PublicMethods[typ.__name__] = class_method
        return self

    def add_public_method(self, name: str, function, argc: int = None) -> 'ClassBuilder':
        method = self.__build_method(name, function, argc)
        self.__class.PublicMethods[name] = method
        return self

    def add_public_coroutine(self, name: str, function) -> 'ClassBuilder':
        argc = len(inspect.signature(function).parameters) - 1

        def wrapper(this, *args):
            fut = concurrent.futures.Future()
            def do():
                try:
                    result = asyncio.run_coroutine_threadsafe(function(this, *args), self.__loop).result()
                    fut.set_result(result)
                except Exception as e:
                    fut.set_exception(e)
            threading.Thread(target=do).start()
            return fut.result()

        method = self.__build_method(name, wrapper, argc)
        self.__class.PublicMethods[name] = method
        return self

    def add_private_method(self, name: str, function, argc: int = None) -> 'ClassBuilder':
        method = self.__build_method(name, function, argc)
        self.__class.PrivateMethods[name] = method
        return self

    def add_private_coroutine(self, name: str, function) -> 'ClassBuilder':
        argc = len(inspect.signature(function).parameters) - 1

        def wrapper(this, *args):
            fut = concurrent.futures.Future()
            def do():
                try:
                    result = asyncio.run_coroutine_threadsafe(function(this, *args), self.__loop).result()
                    fut.set_result(result)
                except Exception as e:
                    fut.set_exception(e)
            threading.Thread(target=do).start()
            return fut.result()

        method = self.__build_method(name, wrapper, argc)
        self.__class.PrivateMethods[name] = method
        return self