from .olol import ObjectiveLOLVM

def run(source_code: str):
    vm = ObjectiveLOLVM()
    vm.execute(source_code)