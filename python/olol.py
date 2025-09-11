import asyncio
import concurrent.futures
import functools
import inspect
import json
import threading
from typing import Any, Callable, Dict, Tuple, Type
from unittest import result
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
    GoValueIDKey,
    ForeignModuleNamespace,
)


defined_functions: Dict[str, Tuple['ObjectiveLOLVM', Callable]] = {}
defined_classes: Dict[str, Type] = {}
object_instances: Dict[str, Any] = {}


# gopy does not support passing complex types directly,
# so we wrap arguments and return values as JSON strings.
# Additionally, using closures seems to result in a segfault
# at https://github.com/python/cpython/blob/v3.13.5/Python/generated_cases.c.h#L2462
# so we use a global dictionary to store the actual functions.
def gopy_wrapper(id: str, json_args: str) -> bytes:
    args = json.loads(json_args)
    try:
        vm, fn = defined_functions[id]
        converted_args = [vm.convert_from_go_value(arg) for arg in args]
        result = fn(*converted_args)
        return json.dumps({"result": vm.convert_to_go_value(result), "error": None}, default=vm.serialize_go_value).encode('utf-8')
    except Exception as e:
        return json.dumps({"result": None, "error": str(e)}).encode('utf-8')


def convert_to_simple_mro(mro: list[str]) -> list[str]:
    simple_mro = []
    for cls_name in mro:
        if cls_name.startswith(ForeignModuleNamespace):
            simple_mro.append(cls_name[len(ForeignModuleNamespace)+1:])
    return simple_mro


def generate_case_permutations(fname):
    """Generate all possible case combinations of fname"""
    if not fname:
        return ['']

    result = []
    first_char = fname[0]
    rest_permutations = generate_case_permutations(fname[1:])

    for rest in rest_permutations:
        if first_char.isalpha():
            result.append(first_char.lower() + rest)
            result.append(first_char.upper() + rest)
        else:
            result.append(first_char + rest)

    return result


class ProxyMeta(type):
    def __new__(mcs: Type, name: str, bases: tuple, attrs: dict, go_value: GoValue = None):
        cls = super().__new__(mcs, name, bases, attrs)
        cls._go_value = go_value
        return cls

    def __call__(cls, *args, **kwargs):
        instance = super().__call__(*args, **kwargs)
        instance._go_value = cls._go_value
        return instance


class ObjectiveLOLVM:
    class ClassBuilder:
        _vm: 'ObjectiveLOLVM'
        _class: ClassDefinition

        def __init__(self, vm: 'ObjectiveLOLVM'):
            self._vm = vm
            self._class = NewClassDefinition()

        def get(self) -> ClassDefinition:
            return self._class

        def set_name(self, name: str) -> 'ObjectiveLOLVM.ClassBuilder':
            self._class.Name = name
            return self

        def __build_variable(self, name: str, value, locked: bool, getter=None, setter=None) -> ClassVariable:
            class_variable = ClassVariable()
            class_variable.Name = name
            class_variable.Value = self._vm.convert_to_go_value(value)
            class_variable.Locked = locked
            if getter is not None:
                unique_id = str(uuid.uuid4())

                def wrapper(this_id):
                    return self._vm.convert_to_go_value(getter(object_instances[this_id]))

                defined_functions[unique_id] = (self._vm, wrapper)
                self._vm._compat.BuildNewClassVariableWithGetter(class_variable, unique_id, gopy_wrapper)
            if setter is not None:
                unique_id = str(uuid.uuid4())

                def wrapper(this_id, value):
                    setter(object_instances[this_id], self._vm.convert_from_go_value(value))

                defined_functions[unique_id] = (self._vm, wrapper)
                self._vm._compat.BuildNewClassVariableWithSetter(class_variable, unique_id, gopy_wrapper)
            return class_variable

        def add_public_variable(self, name: str, value = None, locked: bool = False, getter=None, setter=None) -> 'ObjectiveLOLVM.ClassBuilder':
            variable = self.__build_variable(name, value, locked, getter, setter)
            self._class.PublicVariables[name] = variable
            return self

        def add_private_variable(self, name: str, value = None, locked: bool = False, getter=None, setter=None) -> 'ObjectiveLOLVM.ClassBuilder':
            variable = self.__build_variable(name, value, locked, getter, setter)
            self._class.PrivateVariables[name] = variable
            return self

        def add_shared_variable(self, name: str, value = None, locked: bool = False, getter=None, setter=None) -> 'ObjectiveLOLVM.ClassBuilder':
            variable = self.__build_variable(name, value, locked, getter, setter)
            self._class.SharedVariables[name] = variable
            return self

        def __build_method(self, name: str, function, argc: int = None) -> ClassMethod:
            argc = len(inspect.signature(function).parameters) - 1 if argc is None else argc
            unique_id = str(uuid.uuid4())

            def wrapper(this_id, *args):
                return self._vm.convert_to_go_value(function(object_instances[this_id], *args))

            defined_functions[unique_id] = (self._vm, wrapper)
            class_method = ClassMethod()
            class_method.Name = name
            class_method.Argc = argc

            self._vm._compat.BuildNewClassMethod(class_method, unique_id, gopy_wrapper)
            return class_method

        def add_constructor(self, typ: type) -> 'ObjectiveLOLVM.ClassBuilder':
            # get init function
            init_function = typ.__init__
            argc = len(inspect.signature(init_function).parameters) - 1

            # ignore args and kwargs
            for param in inspect.signature(init_function).parameters.values():
                if param.kind in (param.VAR_POSITIONAL, param.VAR_KEYWORD):
                    argc = argc - 1

            unique_id = str(uuid.uuid4())

            def ctor_wrapper(this_id, *args):
                mro = self._vm._compat.GetObjectMRO(this_id)
                simple_mro = convert_to_simple_mro(mro)

                instance_class = typ
                if len(mro) > 1:
                    go_value = self._vm._compat.LookupObject(this_id)
                    instance_class = self._vm.create_proxy_class([defined_classes[cls_name.upper()] for cls_name in simple_mro if cls_name.upper() in defined_classes], go_value)

                instance = instance_class(*args)
                object_instances[this_id] = instance

            defined_functions[unique_id] = (self._vm, ctor_wrapper)
            class_method = ClassMethod()
            class_method.Name = typ.__name__
            class_method.Argc = argc

            self._vm._compat.BuildNewClassMethod(class_method, unique_id, gopy_wrapper)
            self._class.PublicMethods[typ.__name__] = class_method
            return self

        def add_public_method(self, name: str, function, argc: int = None) -> 'ObjectiveLOLVM.ClassBuilder':
            method = self.__build_method(name, function, argc)
            self._class.PublicMethods[name] = method
            return self

        def add_public_coroutine(self, name: str, function) -> 'ObjectiveLOLVM.ClassBuilder':
            argc = len(inspect.signature(function).parameters) - 1

            def wrapper(this, *args):
                fut = concurrent.futures.Future()
                def do():
                    try:
                        result = asyncio.run_coroutine_threadsafe(function(this, *args), self._vm._loop).result()
                        fut.set_result(result)
                    except Exception as e:
                        fut.set_exception(e)
                threading.Thread(target=do).start()
                return fut.result()

            method = self.__build_method(name, wrapper, argc)
            self._class.PublicMethods[name] = method
            return self

        def add_private_method(self, name: str, function, argc: int = None) -> 'ObjectiveLOLVM.ClassBuilder':
            method = self.__build_method(name, function, argc)
            self._class.PrivateMethods[name] = method
            return self

        def add_private_coroutine(self, name: str, function) -> 'ObjectiveLOLVM.ClassBuilder':
            argc = len(inspect.signature(function).parameters) - 1

            def wrapper(this, *args):
                fut = concurrent.futures.Future()
                def do():
                    try:
                        result = asyncio.run_coroutine_threadsafe(function(this, *args), self._vm._loop).result()
                        fut.set_result(result)
                    except Exception as e:
                        fut.set_exception(e)
                threading.Thread(target=do).start()
                return fut.result()

            method = self.__build_method(name, wrapper, argc)
            self._class.PrivateMethods[name] = method
            return self

        def add_unknown_function_handler(self, function) -> 'ObjectiveLOLVM.ClassBuilder':
            def handler(this_id: str, fname: str, from_context: str, *args):
                return function(object_instances[this_id], fname, from_context, *args)

            unique_id = str(uuid.uuid4())
            defined_functions[unique_id] = (self._vm, handler)
            self._class.UnknownFunctionHandler = self._vm._compat.BuildNewUnknownFunctionHandler(unique_id, gopy_wrapper)
            return self

        def add_unknown_coroutine_handler(self, function) -> 'ObjectiveLOLVM.ClassBuilder':
            def handler(this_id: str, fname: str, from_context: str, *args):
                fut = concurrent.futures.Future()
                def do():
                    try:
                        result = asyncio.run_coroutine_threadsafe(function(object_instances[this_id], fname, from_context, *args), self._vm._loop).result()
                        fut.set_result(result)
                    except Exception as e:
                        fut.set_exception(e)
                threading.Thread(target=do).start()
                return fut.result()

            unique_id = str(uuid.uuid4())
            defined_functions[unique_id] = (self._vm, handler)
            self._class.UnknownFunctionHandler = self._vm._compat.BuildNewUnknownFunctionHandler(unique_id, gopy_wrapper)
            return self

    _vm: VM
    _compat: VMCompatibilityShim
    _loop: asyncio.AbstractEventLoop
    _prefer_async_loop: bool

    def __init__(self, prefer_async_loop: bool = True):
        # todo: figure out how to bridge stdout/stdin
        self._vm = NewVM(DefaultConfig())
        self._compat = self._vm.GetCompatibilityShim()
        self._loop = asyncio.get_event_loop()
        self._prefer_async_loop = prefer_async_loop

    def convert_from_go_value(self, go_value: GoValue):
        if not isinstance(go_value, GoValue):
            if go_value and GoValueIDKey in go_value:
                go_value = self._compat.LookupObject(go_value[GoValueIDKey])
                return self.convert_from_go_value(go_value)
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
        elif typ == "NOTHIN":
            return None
        elif typ == "BUKKIT":
            return [self.convert_from_go_value(v) for v in go_value.Slice()]
        elif typ == "BASKIT":
            return {k: self.convert_from_go_value(v) for k, v in go_value.Map().items()}
        else:
            # object handle
            if go_value.ID() in object_instances:
                return object_instances[go_value.ID()]

            mro = self._compat.GetObjectMRO(go_value.ID())
            simple_mro = convert_to_simple_mro(mro)
            instance_class = self.create_proxy_class([defined_classes[cls_name.upper()] for cls_name in simple_mro if cls_name.upper() in defined_classes], go_value)
            instance = instance_class()
            object_instances[go_value.ID()] = instance
            return instance

    def convert_to_go_value(self, value):
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
                slice.append(self.convert_to_go_value(v))
            return slice
        elif isinstance(value, dict):
            map = Map_string_api_GoValue()
            for k, v in value.items():
                map[k] = self.convert_to_go_value(v)
            return map
        elif isinstance(type(value), ProxyMeta):
            return value._go_value
        else:
            self.define_class(type(value), fully_qualified=True)
            instance = self._vm.NewObjectInstance("{}.{}".format(type(value).__module__, type(value).__name__))

            # for attributes added at runtime (e.g. in __init__),
            # inject getters and setters for them
            for a in dir(value):
                if hasattr(type(value), a):
                    continue

                class_variable = ClassVariable()
                class_variable.Name = a.upper()

                getter = lambda obj, attr=a: getattr(obj, attr)
                unique_id = str(uuid.uuid4())
                def wrapper(this_id, getter=getter):
                    return self.convert_to_go_value(getter(object_instances[this_id]))
                defined_functions[unique_id] = (self, wrapper)
                self._compat.BuildNewClassVariableWithGetter(class_variable, unique_id, gopy_wrapper)

                setter = lambda obj, val, attr=a: setattr(obj, attr, val)
                unique_id = str(uuid.uuid4())
                def wrapper(this_id, value, setter=setter):
                    setter(object_instances[this_id], self.convert_from_go_value(value))
                defined_functions[unique_id] = (self, wrapper)
                self._compat.BuildNewClassVariableWithSetter(class_variable, unique_id, gopy_wrapper)

                self._compat.AddVariableToObject(instance.ID(), class_variable)

            object_instances[instance.ID()] = value
            return instance

    def serialize_go_value(self, go_value: GoValue):
        if isinstance(go_value, GoValue):
            if go_value.ID() != "":
                return {GoValueIDKey: go_value.ID()}
            return self.convert_from_go_value(go_value)
        else:
            return go_value

    def create_proxy_class(self, mro: list[type], go_value: GoValue) -> type:
        instance_immediate_functions = self._compat.GetObjectImmediateFunctions(go_value.ID())
        superself = self

        class Proxy(*mro, metaclass=ProxyMeta, go_value=go_value):
            def __getattribute__(self, name):
                # Handles basic object attributes to avoid infinite recursion
                if name in ('_go_value', '_create_proxy_method', '__class__', '__dict__'):
                    return super().__getattribute__(name)

                # Check if this method should be proxied to the VM
                if name.upper() in instance_immediate_functions:
                    # This method belongs to the immediate class, proxy to VM
                    return self._create_proxy_method(name)

                # For everything else, get it normally
                return super().__getattribute__(name)

            def _create_proxy_method(self, method_name):
                is_async = False
                try:
                    method = super().__getattribute__(method_name)
                    if callable(method):
                        if inspect.iscoroutinefunction(method):
                            is_async = True
                except:
                    pass

                if is_async:
                    return functools.partial(superself.call_method_async, self._go_value, method_name)
                else:
                    return functools.partial(superself.call_method, self._go_value, method_name)

        return Proxy

    def define_variable(self, name: str, value, constant: bool = False) -> None:
        goValue = self.convert_to_go_value(value)
        self._vm.DefineVariable(name, goValue, constant)

    def define_function(self, name: str, function, argc: int = None) -> None:
        argc = len(inspect.signature(function).parameters) if argc is None else argc
        unique_id = str(uuid.uuid4())
        defined_functions[unique_id] = (self, function)
        self._compat.DefineFunction(unique_id, name, argc, gopy_wrapper)

    def define_coroutine(self, name: str, function) -> None:
        argc = len(inspect.signature(function).parameters)

        def wrapper(*args):
            fut = concurrent.futures.Future()
            def do():
                try:
                    result = asyncio.run_coroutine_threadsafe(function(*args), self._loop).result()
                    fut.set_result(result)
                except Exception as e:
                    fut.set_exception(e)
            threading.Thread(target=do).start()
            return fut.result()

        self.define_function(name, wrapper, argc)

    def define_class(self, python_class: type, fully_qualified: bool = False) -> None:
        class_name = f"{python_class.__module__}.{python_class.__name__}" if fully_qualified else python_class.__name__
        if self._compat.IsClassDefined(class_name):
            return

        # Use class builder to introspect and build the class definition
        builder = ObjectiveLOLVM.ClassBuilder(self)
        builder.set_name(class_name)
        builder.add_constructor(python_class)

        # Add class attributes as variables with getters/setters
        for attr_name in dir(python_class):
            if not attr_name.startswith('_') and not callable(getattr(python_class, attr_name)):
                builder.add_public_variable(
                    attr_name,
                    getter=lambda self, attr=attr_name: getattr(self, attr),
                    setter=lambda self, value, attr=attr_name: setattr(self, attr, value)
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

        # Dynamic handler for unknown function calls
        # We assume all unknown functions are publicly available, so
        # no context checking is needed. However, since Objective-LOL
        # is case-insensitive for method names, we need to try all
        # case permutations to find a match.

        # Cache for case-insensitive method name lookups
        _method_name_cache = {}
        if self._prefer_async_loop:
            async def async_handler(this, fname: str, from_context: str, *args):
                # Check cache first
                cache_key = (id(this), fname.upper())
                if cache_key in _method_name_cache:
                    actual_method_name = _method_name_cache[cache_key]
                    if actual_method_name:
                        method = getattr(this, actual_method_name)
                        return await method(*args)
                    else:
                        # Cached as not found
                        raise AttributeError(f"'{type(this).__name__}' object has no attribute '{fname}'")

                # Try all case permutations
                for candidate in generate_case_permutations(fname):
                    try:
                        if hasattr(this, candidate):
                            method = getattr(this, candidate)
                            if callable(method):
                                _method_name_cache[cache_key] = candidate
                                return await method(*args)
                    except Exception as e:
                        continue

                # Cache as not found
                _method_name_cache[cache_key] = None
                raise AttributeError(f"'{type(this).__name__}' object has no attribute '{fname}'")
            builder.add_unknown_coroutine_handler(async_handler)
        else:
            def handler(this, fname: str, from_context: str, *args):
                # Check cache first
                cache_key = (id(this), fname.upper())
                if cache_key in _method_name_cache:
                    actual_method_name = _method_name_cache[cache_key]
                    if actual_method_name:
                        method = getattr(this, actual_method_name)
                        return method(*args)
                    else:
                        # Cached as not found
                        raise AttributeError(f"'{type(this).__name__}' object has no attribute '{fname}'")

                # Try all case permutations
                for candidate in generate_case_permutations(fname):
                    if hasattr(this, candidate):
                        method = getattr(this, candidate)
                        if callable(method):
                            _method_name_cache[cache_key] = candidate
                            return method(*args)

                # Cache as not found
                _method_name_cache[cache_key] = None
                raise AttributeError(f"'{type(this).__name__}' object has no attribute '{fname}'")
            builder.add_unknown_function_handler(handler)

        class_def = builder.get()
        self._vm.DefineClass(class_def)

        defined_classes[class_name.upper()] = python_class

    def call(self, name: str, *args):
        goArgs = self.convert_to_go_value(args)
        result = self._vm.Call(name, goArgs)
        return self.convert_from_go_value(result)

    async def call_async(self, name: str, *args):
        goArgs = self.convert_to_go_value(args)
        fut = concurrent.futures.Future()
        def do():
            try:
                result = self._vm.Call(name, goArgs)
                fut.set_result(self.convert_from_go_value(result))
            except Exception as e:
                fut.set_exception(e)
        threading.Thread(target=do).start()
        return await asyncio.wrap_future(fut)

    def call_method(self, receiver: GoValue, name: str, *args):
        goArgs = self.convert_to_go_value(args)
        result = self._vm.CallMethod(receiver, name, goArgs)
        return self.convert_from_go_value(result)

    async def call_method_async(self, receiver: GoValue, name: str, *args):
        goArgs = self.convert_to_go_value(args)
        fut = concurrent.futures.Future()
        def do():
            try:
                result = self._vm.CallMethod(receiver, name, goArgs)
                fut.set_result(self.convert_from_go_value(result))
            except Exception as e:
                fut.set_exception(e)
        threading.Thread(target=do).start()
        result = await asyncio.wrap_future(fut)
        return result

    def execute(self, code: str) -> None:
        return self._vm.Execute(code)

    async def execute_async(self, code: str) -> None:
        fut = concurrent.futures.Future()
        def do():
            try:
                result = self._vm.Execute(code)
                fut.set_result(result)
            except Exception as e:
                fut.set_exception(e)
        threading.Thread(target=do).start()
        return await asyncio.wrap_future(fut)

