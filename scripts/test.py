import os
import subprocess
import sys

# --- Configuration ---
# The root of the script is the project root.
SCRIPT_ROOT = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))

# Coverage reports will be saved here.
REPORTS_DIR = os.path.join(SCRIPT_ROOT, "reports")

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

def run_tests_with_coverage():
    """Runs Go unit tests and generates a coverage profile in the reports directory."""
    print("--- Running Unit Tests with Coverage ---")
    os.makedirs(REPORTS_DIR, exist_ok=True)
    
    coverage_file = os.path.join(REPORTS_DIR, "coverage.out")
    print(f"Coverage profile will be saved to: {coverage_file}")
    
    cmd = ["go", "test", "-v", f"-coverprofile={coverage_file}", "./..."]
    run_command(cmd)
    print("--- Unit Tests Passed ---")

def main():
    """Main function to run the tests."""
    run_tests_with_coverage()

if __name__ == "__main__":
    main()