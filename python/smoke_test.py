"""Quick check that a built objective_lol wheel imports and works.

Run by cibuildwheel (CIBW_TEST_COMMAND) against each wheel. Exercises the
Go <-> Python boundary in both directions, since callbacks from Go into Python
are the part that differs most between interpreters.
"""
import sys

import objective_lol as olol


def main():
    seen = []

    def add(a, b):
        return a + b

    def record(value):
        seen.append(value)
        return value

    class Calculator:
        def mul(self, x, y):
            return x * y

    vm = olol.ObjectiveLOLVM()
    vm.define_function("add", add)
    vm.define_function("record", record)
    vm.define_class(Calculator)

    vm.execute(
        """
HAI ME TEH FUNCSHUN MAIN
    I HAS A VARIABLE SUM TEH INTEGR ITZ record WIT (add WIT 10 AN WIT 20)
    I HAS A VARIABLE CALC TEH Calculator ITZ NEW Calculator
    I HAS A VARIABLE PRODUCT TEH INTEGR ITZ record WIT (CALC DO mul WIT 6 AN WIT 7)
KTHXBAI
"""
    )

    assert seen == [30, 42], f"unexpected callback results: {seen!r}"
    print(f"objective_lol smoke test passed on {sys.implementation.name} {sys.version.split()[0]}")


if __name__ == "__main__":
    main()
