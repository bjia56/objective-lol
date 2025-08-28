#!/bin/bash

# Test runner for Objective-LOL files
# Runs all .olol files in the tests directory
# Usage: ./run_tests.sh [-v|--verbose]

set -e

TESTS_DIR="tests"
INTERPRETER="./olol"
PASSED=0
FAILED=0
TOTAL=0
VERBOSE=false

# Parse command line arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        -v|--verbose)
            VERBOSE=true
            shift
            ;;
        -h|--help)
            echo "Usage: $0 [-v|--verbose] [-h|--help]"
            echo "  -v, --verbose    Show test output"
            echo "  -h, --help       Show this help message"
            exit 0
            ;;
        *)
            echo "Unknown option: $1"
            echo "Use -h or --help for usage information"
            exit 1
            ;;
    esac
done

if [[ "$VERBOSE" == true ]]; then
    echo "Running Objective-LOL tests (verbose mode)..."
    echo "=============================================="
else
    echo "Running Objective-LOL tests..."
    echo "================================"
fi

# Check if interpreter exists
if [[ ! -x "$INTERPRETER" ]]; then
    echo "Error: Interpreter '$INTERPRETER' not found or not executable"
    exit 1
fi

# Check if tests directory exists
if [[ ! -d "$TESTS_DIR" ]]; then
    echo "Error: Tests directory '$TESTS_DIR' not found"
    exit 1
fi

# Find and run all .olol test files
for test_file in "$TESTS_DIR"/*.olol; do
    if [[ -f "$test_file" ]]; then
        TOTAL=$((TOTAL + 1))
        
        if [[ "$VERBOSE" == true ]]; then
            echo ""
            echo "=== Running $(basename "$test_file") ==="
            
            if "$INTERPRETER" "$test_file" 2>&1; then
                echo ">>> RESULT: PASS"
                PASSED=$((PASSED + 1))
            else
                echo ">>> RESULT: FAIL"
                FAILED=$((FAILED + 1))
            fi
        else
            echo -n "Running $(basename "$test_file")... "
            
            if "$INTERPRETER" "$test_file" > /dev/null 2>&1; then
                echo "PASS"
                PASSED=$((PASSED + 1))
            else
                echo "FAIL"
                FAILED=$((FAILED + 1))
                echo "  Error running $test_file:"
                "$INTERPRETER" "$test_file" 2>&1 | sed 's/^/    /'
            fi
        fi
    fi
done

if [[ "$VERBOSE" == true ]]; then
    echo ""
    echo "=============================================="
else
    echo "================================"
fi
echo "Results: $PASSED passed, $FAILED failed, $TOTAL total"

if [[ $FAILED -eq 0 ]]; then
    echo "All tests passed!"
    exit 0
else
    echo "Some tests failed."
    exit 1
fi