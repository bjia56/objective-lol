import json
import os
import platform
import shutil
import subprocess
import sys
import tarfile
import zipfile
from urllib.request import urlopen, urlretrieve, Request
from urllib.error import URLError
try:
    from distutils.core import Extension
except:
    from setuptools import Extension

import setuptools
from setuptools.command.build_ext import build_ext


PACKAGE_PATH="pkg/api"
PACKAGE_NAME="olol"


def get_python_version():
    """Get the current Python interpreter version"""
    return f"{sys.version_info.major}.{sys.version_info.minor}.{sys.version_info.micro}"


def find_matching_python_release(target_version, detected_platform):
    """Find the best matching Python release from GitHub API with pagination support"""
    base_api_url = "https://api.github.com/repos/bjia56/portable-python/releases"

    try:
        # Parse target version
        target_parts = [int(x) for x in target_version.split('.')]

        # Find releases that match major.minor version
        matching_releases = []
        page = 1
        per_page = 100  # GitHub API max per page

        while True:
            api_url = f"{base_api_url}?per_page={per_page}&page={page}"

            # Create request with GitHub Actions token if available
            request = Request(api_url)
            github_token = os.environ.get('GITHUB_TOKEN')
            if github_token:
                request.add_header('Authorization', f'Bearer {github_token}')
                request.add_header('X-GitHub-Api-Version', '2022-11-28')

            with urlopen(request) as response:
                releases = json.loads(response.read())

            # If no releases returned, we've reached the end
            if not releases:
                break

            # Process releases on this page
            for release in releases:
                tag = release['tag_name']
                if tag.startswith('cpython-v'):
                    # Extract version and build suffix from tag (e.g., cpython-v3.13.5-build.1)
                    tag_suffix = tag.replace('cpython-v', '')
                    version_str = tag_suffix.split('-')[0]
                    build_suffix = tag_suffix[len(version_str)+1:]  # Extract build.X part

                    try:
                        version_parts = [int(x) for x in version_str.split('.')]
                        if len(version_parts) >= 3 and version_parts[0] == target_parts[0] and version_parts[1] == target_parts[1]:
                            # Check if assets exist for our platform
                            build_type = "headless"
                            expected_asset = f"python-{build_type}-{version_str}-{detected_platform}.zip"
                            if any(asset['name'] == expected_asset for asset in release['assets']):
                                matching_releases.append((version_str, version_parts, build_suffix))
                    except ValueError:
                        continue

            # If we found matching releases or got less than per_page results, we can stop
            if matching_releases or len(releases) < per_page:
                break

            page += 1

        if not matching_releases:
            print(f"Warning: No matching Python {target_parts[0]}.{target_parts[1]}.x release found for {detected_platform}")
            return None, None

        # Sort by version (highest patch version first)
        matching_releases.sort(key=lambda x: x[1], reverse=True)
        best_match = matching_releases[0]
        version_str, _, build_suffix = best_match

        print(f"Found Python {version_str} with build suffix {build_suffix} matching interpreter version {target_version}")
        return version_str, build_suffix

    except (URLError, json.JSONDecodeError, KeyError) as e:
        print(f"Error querying GitHub API: {e}")
        return None, None

def detect_platform():
    """Detect platform and architecture for downloading binaries"""
    os_name = platform.system().lower()
    arch = platform.machine().lower()

    if os_name == "linux":
        if arch in ["x86_64", "amd64"]:
            return "linux-x86_64"
        elif arch in ["aarch64", "arm64"]:
            return "linux-aarch64"
        else:
            raise RuntimeError(f"Unsupported architecture: {arch} on Linux")
    elif os_name == "darwin":
        return "darwin-universal2"
    elif os_name == "windows":
        if arch in ["x86_64", "amd64"]:
            return "windows-x86_64"
        elif arch in ["aarch64", "arm64"]:
            return "windows-aarch64"
        else:
            raise RuntimeError(f"Unsupported architecture: {arch} on Windows")
    else:
        raise RuntimeError(f"Unsupported OS: {os_name}")


def download_and_extract_go(toolchain_dir, detected_platform):
    """Download and extract Go toolchain"""
    go_version = "1.21.13"
    go_dir = os.path.join(toolchain_dir, "go")

    # Map platform to Go's naming convention
    # For macOS, detect the actual architecture since Go doesn't have universal binaries
    if detected_platform == "darwin-universal2":
        # Detect actual architecture on macOS
        arch = platform.machine().lower()
        if arch in ["aarch64", "arm64"]:
            go_platform = "darwin-arm64"
        else:
            go_platform = "darwin-amd64"
    else:
        go_platform_map = {
            "linux-x86_64": "linux-amd64",
            "linux-aarch64": "linux-arm64",
            "windows-x86_64": "windows-amd64",
            "windows-aarch64": "windows-arm64"
        }
        go_platform = go_platform_map.get(detected_platform, "linux-amd64")

    # Determine archive format
    if detected_platform.startswith("windows"):
        archive_name = f"go{go_version}.{go_platform}.zip"
    else:
        archive_name = f"go{go_version}.{go_platform}.tar.gz"

    download_url = f"https://go.dev/dl/{archive_name}"
    archive_path = os.path.join(toolchain_dir, archive_name)

    # Check if Go is already installed
    go_exe = os.path.join(go_dir, "bin", "go")
    if detected_platform.startswith("windows"):
        go_exe += ".exe"

    if os.path.exists(go_exe):
        try:
            result = subprocess.run([go_exe, "version"], capture_output=True, text=True)
            if f"go{go_version}" in result.stdout:
                print(f"Go {go_version} already installed")
                return go_dir
        except:
            pass

    # Create directory and download
    os.makedirs(toolchain_dir, exist_ok=True)

    if not os.path.exists(archive_path):
        print(f"Downloading Go {go_version} for {go_platform}...")
        urlretrieve(download_url, archive_path)

    # Remove existing installation
    if os.path.exists(go_dir):
        shutil.rmtree(go_dir)

    # Extract
    print("Extracting Go...")
    if detected_platform.startswith("windows"):
        with zipfile.ZipFile(archive_path, 'r') as zip_ref:
            zip_ref.extractall(toolchain_dir)
    else:
        with tarfile.open(archive_path, 'r:gz') as tar_ref:
            tar_ref.extractall(toolchain_dir)

    return go_dir


def download_and_extract_python(toolchain_dir, detected_platform):
    """Download and extract Python"""
    # Get current Python version and find matching release
    current_version = get_python_version()
    python_version, build_suffix = find_matching_python_release(current_version, detected_platform)

    if python_version is None:
        raise RuntimeError(f"No compatible Python release found for version {current_version} and platform {detected_platform}")

    print(f"Using Python version {python_version} with build suffix {build_suffix}")
    build_type = "headless"
    archive_name = f"python-{build_type}-{python_version}-{detected_platform}.zip"
    download_url = f"https://github.com/bjia56/portable-python/releases/download/cpython-v{python_version}-{build_suffix}/{archive_name}"

    python_dir = os.path.join(toolchain_dir, "python")
    archive_path = os.path.join(toolchain_dir, archive_name)
    extract_dir = os.path.join(python_dir, f"python-{python_version}")

    # Check if already extracted
    if os.path.exists(extract_dir):
        print(f"Python {python_version} already installed")
        return extract_dir

    # Create directory and download
    os.makedirs(python_dir, exist_ok=True)

    if not os.path.exists(archive_path):
        print(f"Downloading Python {python_version} for {detected_platform}...")
        urlretrieve(download_url, archive_path)

    # Extract
    print("Extracting Python...")
    with zipfile.ZipFile(archive_path, 'r') as zip_ref:
        zip_ref.extractall(python_dir)

    # Handle directory naming - the archive might extract to different name
    extracted_name = archive_name.replace(".zip", "")
    extracted_path = os.path.join(python_dir, extracted_name)
    if os.path.exists(extracted_path) and not os.path.exists(extract_dir):
        os.rename(extracted_path, extract_dir)

    return extract_dir


def install_go_dependencies(go_exe, gobin_dir):
    """Install Go dependencies (goimports and gopy)"""
    env = os.environ.copy()
    env["GOBIN"] = gobin_dir

    print("Installing goimports...")
    subprocess.check_call([go_exe, "install", "golang.org/x/tools/cmd/goimports@v0.17.0"], env=env)

    print("Installing gopy...")
    subprocess.check_call([go_exe, "install", "github.com/go-python/gopy@v0.4.10"], env=env)


class CustomBuildExt(build_ext):
    def build_extension(self, ext: Extension):
        # Detect platform
        detected_platform = detect_platform()

        # Setup toolchain directories
        script_dir = os.path.dirname(os.path.abspath(__file__))
        toolchain_dir = os.path.join(script_dir, ".toolchain")

        # Download and setup Go
        go_dir = download_and_extract_go(toolchain_dir, detected_platform)
        go_exe = os.path.join(go_dir, "bin", "go")
        if detected_platform.startswith("windows"):
            go_exe += ".exe"

        # Download and setup Python
        python_dir = download_and_extract_python(toolchain_dir, detected_platform)
        if detected_platform.startswith("windows"):
            python_exe = os.path.join(python_dir, "bin", "python.exe")
        else:
            python_exe = os.path.join(python_dir, "bin", "python3")
            os.chmod(python_exe, 0o755)

        # Setup GOBIN directory and install dependencies
        gobin_dir = os.path.join(toolchain_dir, "gobin")
        os.makedirs(gobin_dir, exist_ok=True)
        install_go_dependencies(go_exe, gobin_dir)

        # Install pybindgen for the downloaded Python
        subprocess.check_call([python_exe, "-m", "pip", "install", "pybindgen"])

        # Get Go environment
        go_env_result = subprocess.check_output([go_exe, "env", "-json"], text=True)
        go_env = json.loads(go_env_result)

        # Setup destination directory
        destination = os.path.join(os.path.dirname(os.path.abspath(self.get_ext_fullpath(ext.name))), PACKAGE_NAME)
        if os.path.isdir(destination):
            shutil.rmtree(destination)

        # Setup environment variables
        # Use platform-appropriate path separator
        path_sep = ";" if sys.platform == "win32" else ":"
        new_path = f"{gobin_dir}{path_sep}{go_dir}/bin{path_sep}{os.environ.get('PATH', '')}"

        env = {
            "PATH": new_path,
            **go_env,
            "CGO_LDFLAGS_ALLOW": ".*",
            "GOWORK": "off",
            "CGO_ENABLED": "1",
        }

        # Platform-specific environment setup
        if sys.platform == "win32":
            env["SYSTEMROOT"] = os.environ.get("SYSTEMROOT", "")
            # Use downloaded Python paths
            python_include = os.path.join(python_dir, "include")
            python_lib = os.path.join(python_dir, "libs")
            # Use running Python version for lib name
            python_lib_name = f"python{sys.version_info.major}{sys.version_info.minor}.lib"
            env["CGO_CFLAGS"] = f"-I{python_include}"
            env["CGO_LDFLAGS"] = f"-L{python_lib} -l:{python_lib_name}"
            env["GOPY_INCLUDE"] = python_include
            env["GOPY_LIBDIR"] = python_lib
            env["GOPY_PYLIB"] = f":{python_lib_name}"
        elif sys.platform == "darwin":
            min_ver = os.environ.get("MACOSX_DEPLOYMENT_TARGET", "10.15")
            env["MACOSX_DEPLOYMENT_TARGET"] = min_ver
            env["CGO_LDFLAGS"] = f"-mmacosx-version-min={min_ver}"
            env["CGO_CFLAGS"] = f"-mmacosx-version-min={min_ver}"

        # Change to the source directory
        original_cwd = os.getcwd()
        source_dir = os.path.join(os.path.dirname(script_dir), PACKAGE_PATH)
        os.chdir(source_dir)

        try:
            # Run gopy build
            gopy_exe = os.path.join(gobin_dir, "gopy")
            if detected_platform.startswith("windows"):
                gopy_exe += ".exe"

            subprocess.check_call([
                gopy_exe,
                "build",
                "-no-make",
                "-dynamic-link=True",
                "-symbols=False",
                "-output", destination,
                "--vm", python_exe,
                "."
            ], env=env)
        finally:
            for file in ["__init__.py", "olol.py"]:
                source_file = os.path.join(script_dir, file)
                shutil.copy(source_file, os.path.join(destination, file))
            os.chdir(original_cwd)


setuptools.setup(
    name=PACKAGE_NAME,
    version="0.0.1",
    author="Brett Jia",
    author_email="dev.bjia56@gmail.com",
    description="Python bindings for Objective-LOL",
    url="https://github.com/bjia56/objective-lol",
    classifiers=[
        "Programming Language :: Python :: 3",
        "License :: OSI Approved :: MIT License",
    ],
    include_package_data=True,
    cmdclass={
        "build_ext": CustomBuildExt,
    },
    ext_modules=[
        Extension(PACKAGE_NAME, [PACKAGE_PATH]),
    ],
)
