import argparse
import os
import platform
import subprocess
import sys

# Define color codes for UI
class Colors:
    GREEN = '\033[92m'
    RED = '\033[91m'
    YELLOW = '\033[93m'
    BLUE = '\033[94m'
    RESET = '\033[0m'
    BOLD = '\033[1m'

TOOLS = {
    "go": {
        "name": "Go",
        "description": "Go Compiler (>= 1.25)",
        "check_cmd": ["go", "version"],
        "install_win": ["winget", "install", "-e", "--id", "GoLang.Go"],
        "install_mac": ["brew", "install", "go"],
        "install_linux": ["sudo", "snap", "install", "go", "--classic"]
    },
    "zig": {
        "name": "Zig",
        "description": "Zig Toolchain (>= 0.11.0) for CGO cross-compilation",
        "check_cmd": ["zig", "version"],
        "install_win": ["winget", "install", "-e", "--id", "zig.zig"],
        "install_mac": ["brew", "install", "zig"],
        "install_linux": ["sudo", "snap", "install", "zig", "--classic"]
    },
    "node": {
        "name": "Node.js",
        "description": "Node.js (>= 20.x) for UI development",
        "check_cmd": ["node", "--version"],
        "install_win": ["winget", "install", "-e", "--id", "OpenJS.NodeJS"],
        "install_mac": ["brew", "install", "node"],
        "install_linux": ["sudo", "snap", "install", "node", "--classic"]
    },
    "goreleaser": {
        "name": "GoReleaser",
        "description": "GoReleaser for building and publishing artifacts",
        "check_cmd": ["goreleaser", "--version"],
        "install_win": ["winget", "install", "-e", "--id", "GoReleaser.GoReleaser"],
        "install_mac": ["brew", "install", "goreleaser"],
        "install_linux": ["sudo", "snap", "install", "goreleaser", "--classic"]
    }
}

def print_status(name, is_installed, details=""):
    symbol = f"{Colors.GREEN}[✓]{Colors.RESET}" if is_installed else f"{Colors.RED}[✗]{Colors.RESET}"
    text_color = Colors.GREEN if is_installed else Colors.RED
    detail_str = f" - {Colors.BLUE}{details}{Colors.RESET}" if details else ""
    print(f"  {symbol} {text_color}{name}{Colors.RESET}{detail_str}")

def check_tool(tool_key):
    try:
        cmd = TOOLS[tool_key]["check_cmd"]
        # Use shell=True on windows if it's a built-in or batch file, but list args are better without
        # Sometimes go/node/etc exist but don't respond well. Capturing output is best.
        result = subprocess.run(cmd, capture_output=True, text=True, check=True)
        # return the first line of output as version info
        output = result.stdout.strip() if result.stdout else result.stderr.strip()
        version = output.split('\n')[0]
        return True, version
    except (subprocess.CalledProcessError, FileNotFoundError):
        return False, "Not found"

def get_install_cmd(tool_key):
    os_name = platform.system().lower()
    tool = TOOLS[tool_key]
    
    if os_name == "windows":
        return tool.get("install_win")
    elif os_name == "darwin":
        return tool.get("install_mac")
    elif os_name == "linux":
        return tool.get("install_linux")
    return None

def install_tool(tool_key):
    cmd = get_install_cmd(tool_key)
    if not cmd:
        print(f"{Colors.RED}  No automated installation available for {TOOLS[tool_key]['name']} on your OS.{Colors.RESET}")
        return False
        
    print(f"{Colors.YELLOW}  Installing {TOOLS[tool_key]['name']}... ({' '.join(cmd)}){Colors.RESET}")
    try:
        subprocess.run(cmd, check=True)
        return True
    except subprocess.CalledProcessError as e:
        print(f"{Colors.RED}  Installation failed for {TOOLS[tool_key]['name']}.{Colors.RESET}")
        return False
    except FileNotFoundError:
        print(f"{Colors.RED}  Package manager ({cmd[0]}) not found. Please install manually.{Colors.RESET}")
        return False

def main():
    parser = argparse.ArgumentParser(description="Nanovision Local Environment Setup (`doctor`)")
    parser.add_argument("-y", "--yes", action="store_true", help="Automatically install missing dependencies without asking")
    args = parser.parse_args()

    # Enable VT100 escapes on Windows to display colors
    if platform.system().lower() == "windows":
        os.system("color") 

    print(f"\n{Colors.BOLD}Nanovision Status Doctor{Colors.RESET}")
    print("=" * 40)

    missing_tools = []

    for key, tool in TOOLS.items():
        is_installed, version_info = check_tool(key)
        print_status(tool["name"], is_installed, version_info if is_installed else tool["description"])
        if not is_installed:
            missing_tools.append(key)

    print("=" * 40)
    
    if not missing_tools:
        print(f"{Colors.GREEN}{Colors.BOLD}Your development environment is ready!{Colors.RESET}\n")
        sys.exit(0)

    print(f"\n{Colors.YELLOW}Some required tools are missing.{Colors.RESET}")
    
    if not args.yes:
        response = input(f"\nDo you want me to try and install the missing tools? [y/N]: ").strip().lower()
        if response != 'y':
            print("Setup aborted by user. Please install the missing tools manually.")
            sys.exit(1)

    print("\nStarting installation sequence...")
    all_success = True
    for key in missing_tools:
        success = install_tool(key)
        if not success:
            all_success = False

    print("\n" + "=" * 40)
    if all_success:
        print(f"{Colors.GREEN}{Colors.BOLD}Successfully installed everything! Please restart your terminal/IDE for PATH changes to take effect.{Colors.RESET}\n")
    else:
        print(f"{Colors.RED}{Colors.BOLD}Some installations failed. You may need to install them manually.{Colors.RESET}\n")

if __name__ == "__main__":
    main()
