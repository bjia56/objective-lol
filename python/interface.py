from .vm import *

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

    def execute(self, code: str) -> None:
        return self.vm.Execute(code)

