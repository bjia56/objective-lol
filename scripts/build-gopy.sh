#!/bin/bash

set -e

# Platform-specific executable suffix
if [[ "$OSTYPE" == "msys" || "$OSTYPE" == "win32" ]]; then
    SUFFIX=".exe"
else
    SUFFIX=""
fi

# Paths
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(realpath "$SCRIPT_DIR/..")"
GO_TOOLCHAIN_DIR="$SCRIPT_DIR/.go"
GO_BIN_DIR="$GO_TOOLCHAIN_DIR/go/bin"
GO_EXE="$GO_BIN_DIR/go$SUFFIX"
GOBIN_DIR="$GO_TOOLCHAIN_DIR/gobin"
PYTHON_DIR="$SCRIPT_DIR/.python"

# Function to detect platform and architecture
detect_platform() {
    case "$OSTYPE" in
        linux*)
            case "$(uname -m)" in
                x86_64) echo "linux-x86_64" ;;
                aarch64) echo "linux-aarch64" ;;
                *) echo "linux-x86_64" ;;  # Default fallback
            esac
            ;;
        darwin*)
            echo "darwin-universal2"
            ;;
        msys*|win32*|cygwin*)
            case "$(uname -m)" in
                x86_64) echo "windows-x86_64" ;;
                aarch64) echo "windows-aarch64" ;;
                *) echo "windows-x86_64" ;;  # Default fallback
            esac
            ;;
        *)
            echo "linux-x86_64"  # Default fallback
            ;;
    esac
}

# Function to download and extract Python
download_python() {
    local platform="$1"
    local python_version="3.13.5"
    local build_type="headless"  # Use headless build for smaller size
    local archive_name="python-${build_type}-${python_version}-${platform}.zip"
    local download_url="https://github.com/bjia56/portable-python/releases/download/cpython-v${python_version}-build.1/${archive_name}"

    # Create python directory
    mkdir -p "$PYTHON_DIR"

    # Download Python if not already present
    if [[ ! -f "$PYTHON_DIR/$archive_name" ]]; then
        echo "Downloading Python ${python_version} for ${platform}..." >&2
        if command -v curl >/dev/null 2>&1; then
            curl -L -o "$PYTHON_DIR/$archive_name" "$download_url"
        elif command -v wget >/dev/null 2>&1; then
            wget -O "$PYTHON_DIR/$archive_name" "$download_url"
        else
            echo "Error: Neither curl nor wget found. Please install one of them." >&2
            exit 1
        fi
    else
        echo "Python archive already downloaded." >&2
    fi

    # Extract Python if not already extracted
    local extract_dir="$PYTHON_DIR/python-${python_version}"
    if [[ ! -d "$extract_dir" ]]; then
        echo "Extracting Python..." >&2
        if command -v unzip >/dev/null 2>&1; then
            unzip -q "$PYTHON_DIR/$archive_name" -d "$PYTHON_DIR"
        elif command -v 7z.exe >/dev/null 2>&1; then
            7z.exe x "$PYTHON_DIR/$archive_name" -o"$PYTHON_DIR" -y
        elif command -v 7z >/dev/null 2>&1; then
            7z x "$PYTHON_DIR/$archive_name" -o"$PYTHON_DIR" -y
        else
            echo "Error: Neither unzip nor 7z found. Please install one of them." >&2
            exit 1
        fi

        # The archive extracts to a directory named after the archive (without .zip)
        local extracted_name="${archive_name%.zip}"
        if [[ -d "$PYTHON_DIR/$extracted_name" ]]; then
            mv "$PYTHON_DIR/$extracted_name" "$extract_dir"
        fi
    fi

    echo "$extract_dir"
}

# Function to download and extract Go
download_go() {
    local platform="$1"
    local go_version="1.21.13"
    local go_dir="$GO_TOOLCHAIN_DIR/go"

    # Map platform to Go's naming convention
    local go_platform
    case "$platform" in
        linux-x86_64) go_platform="linux-amd64" ;;
        linux-aarch64) go_platform="linux-arm64" ;;
        darwin-universal2)
            # For macOS, we'll use amd64 as it works on both Intel and Apple Silicon
            go_platform="darwin-amd64" ;;
        windows-x86_64) go_platform="windows-amd64" ;;
        windows-aarch64) go_platform="windows-arm64" ;;
        *) go_platform="linux-amd64" ;;  # Default fallback
    esac

    local archive_name="go${go_version}.${go_platform}.tar.gz"
    local download_url="https://go.dev/dl/${archive_name}"

    # For Windows, use zip instead of tar.gz
    if [[ "$platform" == windows-* ]]; then
        archive_name="go${go_version}.${go_platform}.zip"
        download_url="https://go.dev/dl/${archive_name}"
    fi

    echo "Checking for Go ${go_version}..." >&2

    # Check if Go is already installed and correct version
    if [[ -f "$GO_EXE" ]]; then
        local current_version=$("$GO_EXE" version 2>/dev/null | grep -o 'go[0-9]\+\.[0-9]\+\.[0-9]\+' | head -1)
        if [[ "$current_version" == "go${go_version}" ]]; then
            echo "Go ${go_version} already installed." >&2
            echo "$GO_TOOLCHAIN_DIR"
            return 0
        fi
    fi

    # Create go directory
    mkdir -p "$GO_TOOLCHAIN_DIR"

    # Download Go if not already present
    if [[ ! -f "$GO_TOOLCHAIN_DIR/$archive_name" ]]; then
        echo "Downloading Go ${go_version} for ${go_platform}..." >&2
        if command -v curl >/dev/null 2>&1; then
            curl -L -o "$GO_TOOLCHAIN_DIR/$archive_name" "$download_url"
        elif command -v wget >/dev/null 2>&1; then
            wget -O "$GO_TOOLCHAIN_DIR/$archive_name" "$download_url"
        else
            echo "Error: Neither curl nor wget found. Please install one of them." >&2
            exit 1
        fi
    else
        echo "Go archive already downloaded." >&2
    fi

    # Remove existing Go installation if present
    if [[ -d "$go_dir" ]]; then
        rm -rf "$go_dir"
    fi

    # Extract Go
    echo "Extracting Go..." >&2
    if [[ "$platform" == windows-* ]]; then
        # Windows zip extraction
        if command -v unzip >/dev/null 2>&1; then
            unzip -q "$GO_TOOLCHAIN_DIR/$archive_name" -d "$GO_TOOLCHAIN_DIR"
        elif command -v 7z.exe >/dev/null 2>&1; then
            7z.exe x "$GO_TOOLCHAIN_DIR/$archive_name" -o"$GO_TOOLCHAIN_DIR" -y
        elif command -v 7z >/dev/null 2>&1; then
            7z x "$GO_TOOLCHAIN_DIR/$archive_name" -o"$GO_TOOLCHAIN_DIR" -y
        else
            echo "Error: No extraction tool found for zip files." >&2
            exit 1
        fi
    else
        # Unix tar.gz extraction
        tar -xzf "$GO_TOOLCHAIN_DIR/$archive_name" -C "$GO_TOOLCHAIN_DIR"
    fi

    echo "$GO_TOOLCHAIN_DIR"
}

# Detect platform and download Go and Python
PLATFORM=$(detect_platform)
download_go "$PLATFORM" >/dev/null
PYTHON_ROOT=$(download_python "$PLATFORM")

# Set Python executable path
if [[ "$OSTYPE" == "msys" || "$OSTYPE" == "win32" || "$OSTYPE" == "cygwin" ]]; then
    PYTHON_EXE="$PYTHON_ROOT/bin/python.exe"
else
    PYTHON_EXE="$PYTHON_ROOT/bin/python3"
fi

# Verify Python executable exists
if [[ ! -f "$PYTHON_EXE" ]]; then
    echo "Error: Python executable not found at $PYTHON_EXE" >&2
    exit 1
fi

echo "Using Python: $PYTHON_EXE"

$PYTHON_EXE -m pip install pybindgen

# Get go env and set environment variables
GO_ENV_OUTPUT=$("$GO_EXE" env)
declare -A GO_ENV_DICT

# Parse go env output
while IFS='=' read -r key value; do
    # Remove 'set' prefix if present
    key=$(echo "$key" | sed 's/^set//' | xargs)
    # Remove quotes
    value=$(echo "$value" | sed 's/['"'"'"]//g')
    if [[ -n "$key" ]]; then
        GO_ENV_DICT["$key"]="$value"
    fi
done <<< "$GO_ENV_OUTPUT"

# Set GOBIN
export GOBIN="$GOBIN_DIR"

# Install goimports
echo "Installing goimports..."
"$GO_EXE" install golang.org/x/tools/cmd/goimports@v0.17.0

# Install gopy
echo "Installing gopy..."
"$GO_EXE" install github.com/go-python/gopy@v0.4.10

# Define paths
TARGET_DIR="$PROJECT_ROOT/out"
SRC_DIR="$PROJECT_ROOT/pkg/api"

# Set environment variables for gopy
export GOWORK="off"
export CGO_ENABLED="1"
export CGO_LDFLAGS_ALLOW=".*"
export PATH="$GOBIN_DIR:$GO_BIN_DIR:$PATH"

# Platform-specific CGO flags
case "$OSTYPE" in
    darwin*)
        MACOSX_DEPLOYMENT_TARGET="${MACOSX_DEPLOYMENT_TARGET:-10.15}"
        export MACOSX_DEPLOYMENT_TARGET="$MACOSX_DEPLOYMENT_TARGET"
        export CGO_LDFLAGS="-mmacosx-version-min=$MACOSX_DEPLOYMENT_TARGET"
        export CGO_CFLAGS="-mmacosx-version-min=$MACOSX_DEPLOYMENT_TARGET"
        ;;
    msys*|win32*|cygwin*)
        # Use downloaded Python paths
        PYTHON_INCLUDE="$PYTHON_ROOT/include"
        PYTHON_LIB="$PYTHON_ROOT/libs"
        export CGO_CFLAGS="-I$PYTHON_INCLUDE"
        export CGO_LDFLAGS="-L$PYTHON_LIB -l:python313.lib"
        export GOPY_INCLUDE="$PYTHON_INCLUDE"
        export GOPY_LIBDIR="$PYTHON_LIB"
        export GOPY_PYLIB=":python313.lib"
        ;;
esac

# Print environment for debugging
echo "Environment variables:"
env | grep -E "^(GO|CGO|GOPY|MACOSX)" | sort

# Create target directory if it doesn't exist
mkdir -p "$TARGET_DIR"

# Change to src directory
cd "$SRC_DIR"

# Run gopy build
echo "Running gopy build..."
"$GOBIN_DIR/gopy$SUFFIX" build \
    -no-make \
    -dynamic-link=True \
    -symbols=False \
    -output "$TARGET_DIR" \
    --vm "$PYTHON_EXE" \
    "."

# Check if gopy succeeded
if [ $? -eq 0 ]; then
    echo "gopy build completed successfully."

    cd "$TARGET_DIR"

    # Write __init__.py file
    echo "from .api import *" > __init__.py

    exit 0

    # Remove unused files
    rm -f api.h \
          api.go \
          api.c \
          go.mod \
          go.sum \
          build.py \
          _api.cpython-313-x86_64-linux-gnu.h \
          _api.cpython-313-aarch64-linux-gnu.h \
          _api.cpython-313-darwin.h \
          _api.cp313-win_amd64.h

    echo "Build process completed successfully!"
else
    echo "gopy build failed!" >&2
    exit 1
fi