import argparse
import os
import subprocess
import sys
import shutil
import datetime
import platform

# Configuration
SCRIPT_ROOT = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))
OUTPUT_DIR = os.path.join(SCRIPT_ROOT, "bin")
MAIN_PACKAGE = "cmd/main.go"
BINARY_NAME = "nanovision"

def get_version_info():
    """Returns the commit hash and current UTC timestamp."""
    try:
        commit = subprocess.check_output(
            ["git", "rev-parse", "--short", "HEAD"], cwd=SCRIPT_ROOT
        ).decode().strip()
    except Exception:
        commit = "unknown"
    
    date = datetime.datetime.now(datetime.timezone.utc).strftime("%Y-%m-%dT%H:%M:%SZ")
    return commit, date

def run_go_build(cmd, env):
    """Executes the Go build command. Returns (success, output_message)."""
    print(f"   > {' '.join(cmd)}")
    try:
        # Copy current env and merge with our build env
        full_env = os.environ.copy()
        full_env.update(env)

        result = subprocess.run(
            cmd, 
            cwd=SCRIPT_ROOT, 
            env=full_env, 
            capture_output=True, 
            text=True, 
            check=True
        )
        return True, result.stdout
    except subprocess.CalledProcessError as e:
        # Return the error output if build fails
        return False, e.stderr

def build_target(target_os, version):
    """Builds and archives the binary for a specific OS."""
    print(f"\nBuilding for {target_os.upper()} (Version: {version})")

    # Setup Paths
    commit, date = get_version_info()
    ext = ".exe" if target_os == "windows" else ""
    
    # We build into bin/windows/nanovision.exe or bin/linux/nanovision
    target_bin_dir = os.path.join(OUTPUT_DIR, target_os)
    os.makedirs(target_bin_dir, exist_ok=True)
    
    output_binary = os.path.join(target_bin_dir, f"{BINARY_NAME}{ext}")

    # Setup Linker Flags (Inject Version Info)
    ldflags = (
        f"-s -w "
        f"-X main.version={version} "
        f"-X main.commit={commit} "
        f"-X main.date={date}"
    )

    # Setup Build Environment
    env = {
        "GOOS": target_os,
        "GOARCH": "amd64",
        "CGO_ENABLED": "1"
    }

    # Detect Cross-Compilation to configure CGO and zig cc
    host_os = platform.system().lower()
    if target_os != host_os:
        if target_os == "windows":
            env["CC"] = "zig cc -target x86_64-windows-gnu"
            env["CXX"] = "zig c++ -target x86_64-windows-gnu"
        elif target_os == "linux":
            env["CC"] = "zig cc -target x86_64-linux-gnu"
            env["CXX"] = "zig c++ -target x86_64-linux-gnu"
        elif target_os == "darwin":
            env["CC"] = "zig cc -target x86_64-macos"
            env["CXX"] = "zig c++ -target x86_64-macos"

    # Run Build
    cmd = ["go", "build", "-mod=vendor", "-ldflags", ldflags, "-o", output_binary, MAIN_PACKAGE]
    success, output = run_go_build(cmd, env)

    if not success:
        print("Build compilation failed.")
        return False, output.strip()

    # Create Archive
    archive_format = "zip" if target_os == "windows" else "gztar"
    archive_name = f"nanovision_{version}_{target_os}_amd64"
    archive_path_root = os.path.join(OUTPUT_DIR, archive_name)

    try:
        final_archive = shutil.make_archive(archive_path_root, archive_format, target_bin_dir)
        print(f"Archive created: {os.path.basename(final_archive)}")
        return True, final_archive
    except Exception as e:
        return False, f"Archiving failed: {str(e)}"

def print_summary(results):
    """Prints a simple summary table of the build results."""
    print("\n" + "="*90)
    print("BUILD SUMMARY".center(90))
    print("="*90)
    
    exit_code = 0
    for res in results:
        status = "SUCCESS" if res['success'] else "FAILED "
        # If successful, show the relative path to the archive. If failed, show the error.
        if res['success']:
            details = os.path.relpath(res['msg'], SCRIPT_ROOT)
        else:
            details = res['msg'].split('\n')[0] # Show first line of error
            exit_code = 1
            
        print(f" {res['os'].ljust(10)} : {status} -> {details}")
    
    print("="*90 + "\n")
    return exit_code

def main():
    parser = argparse.ArgumentParser(description="Build nanovision binaries.")
    parser.add_argument("--version", default="dev", help="Version string to inject")
    parser.add_argument("--target", default='all', choices=['all', 'linux', 'windows'])
    args = parser.parse_args()

    # Clear previous builds
    if os.path.exists(OUTPUT_DIR):
        shutil.rmtree(OUTPUT_DIR)

    targets = ["linux", "windows"] if args.target == 'all' else [args.target]
    results = []

    for t in targets:
        success, msg = build_target(t, args.version)
        results.append({"os": t, "success": success, "msg": msg})

    sys.exit(print_summary(results))

if __name__ == "__main__":
    main()