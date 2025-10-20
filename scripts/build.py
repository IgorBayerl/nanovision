import argparse
import os
import subprocess
import sys
import shutil

# --- Configuration ---
# The root of the script is the project root.
SCRIPT_ROOT = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))

# Final binaries will be placed in this directory.
OUTPUT_DIR = os.path.join(SCRIPT_ROOT, "bin")

# The single entry point for the Go application.
MAIN_PACKAGE_PATH = "cmd/main.go"

# Platform-specific binary names.
BINARY_NAME_LINUX = "adlercov"
BINARY_NAME_WINDOWS = "adlercov.exe"


def run_command(command, env=None, working_dir=SCRIPT_ROOT):
    """Executes a command, streams output, and exits if it fails."""
    print(f"--- Running Command: {' '.join(command)}")
    try:
        process_env = os.environ.copy()
        if env:
            process_env.update(env)
        subprocess.run(
            command, check=True, env=process_env, cwd=working_dir,
            stdout=sys.stdout, stderr=sys.stderr
        )
    except subprocess.CalledProcessError as e:
        print(f"--- Command failed with exit code {e.returncode}", file=sys.stderr)
        sys.exit(e.returncode)
    except FileNotFoundError:
        print(f"--- Error: Command not found: {command[0]}", file=sys.stderr)
        sys.exit(1)

def clean_directory(path):
    """Removes a directory and all its contents, then recreates it."""
    if os.path.exists(path):
        print(f"--- Cleaning directory: {path} ---")
        shutil.rmtree(path)
    print(f"--- Creating directory: {path} ---")
    os.makedirs(path)

def build_linux():
    """Builds the Go project for Linux."""
    print("--- Building for Linux (amd64) ---")
    env = {"GOOS": "linux", "GOARCH": "amd64"}
    output_path = os.path.join(OUTPUT_DIR, BINARY_NAME_LINUX)
    cmd = ["go", "build", "-mod=vendor", "-o", output_path, MAIN_PACKAGE_PATH]
    run_command(cmd, env=env)
    print(f"--- Linux build complete: {output_path} ---")

def build_windows():
    """Builds the Go project for Windows, enabling CGo for cross-compilation."""
    print("--- Building for Windows (amd64) ---")
    env = {
        "GOOS": "windows", "GOARCH": "amd64", "CGO_ENABLED": "1",
        "CC": "x86_64-w64-mingw32-gcc", "CXX": "x86_64-w64-mingw32-g++"
    }
    output_path = os.path.join(OUTPUT_DIR, BINARY_NAME_WINDOWS)
    cmd = ["go", "build", "-mod=vendor", "-o", output_path, MAIN_PACKAGE_PATH]
    run_command(cmd, env=env)
    print(f"--- Windows build complete: {output_path} ---")

def main():
    """Main function to parse arguments and execute the corresponding action."""
    parser = argparse.ArgumentParser(description="Build script for the AdlerCov project.")
    parser.add_argument(
        "--target",
        default='all',
        choices=['all', 'linux', 'windows'],
        help="The target platform to build for: 'all' (default), 'linux', or 'windows'."
    )
    args = parser.parse_args()

    clean_directory(OUTPUT_DIR)

    if args.target == 'all':
        build_linux()
        build_windows()
    elif args.target == 'linux':
        build_linux()
    elif args.target == 'windows':
        build_windows()

    print("\n--- All builds completed successfully! ---")
    print(f"Final binaries are located in: {OUTPUT_DIR}")


if __name__ == "__main__":
    main()