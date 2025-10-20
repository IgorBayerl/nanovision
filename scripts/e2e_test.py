#!/usr/bin/env python3
"""
End-to-End (E2E) Test Script for the AdlerCov Project.

This script validates the compiled AdlerCov CLI tool by running it against a
variety of demo projects and coverage report formats.

Key Features:
- Builds the tool automatically for the host OS.
- Runs a predefined suite of test cases for C#, Go, and C++ projects.
- Tests both individual and merged report generation scenarios.
- Verifies command exit codes and the creation of expected output files.
- Provides a detailed summary of PASS/FAIL for all tests.
- Includes a '--self-cover' mode to generate a coverage report for the
  AdlerCov tool itself, combining unit and E2E test coverage.
- Use the '-v' or '--verbose' flag to see the live output from the CLI tool.
"""
import argparse
import os
import shutil
import subprocess
import sys
import tempfile
from dataclasses import dataclass, field

# ==============================================================================
#  Configuration: Paths and Constants
# ==============================================================================

# Core Project Paths
SCRIPT_ROOT = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))
SCRIPTS_DIR = os.path.join(SCRIPT_ROOT, "scripts")
REPORTS_OUTPUT_DIR = os.path.join(SCRIPT_ROOT, "reports") # For E2E report outputs
BINARY_DIR = os.path.join(SCRIPT_ROOT, "bin")
BUILD_SCRIPT_PATH = os.path.join(SCRIPTS_DIR, "build.py")

# Demo Project Paths
DEMO_PROJECTS_ROOT = os.path.join(SCRIPT_ROOT, "demo_projects")
CPP_DIR = os.path.join(DEMO_PROJECTS_ROOT, "cpp")
CPP_PROJECT_DIR = os.path.join(CPP_DIR, "project")
CPP_COBERTURA_XML = os.path.join(CPP_DIR, "report", "cobertura", "cobertura.xml")
CPP_GCOV_DIR = os.path.join(CPP_DIR, "report", "gcov", "branch-probabilities")
CPP_GCOV_PATTERN = os.path.join(CPP_GCOV_DIR, "*.gcov")

CSHARP_DIR = os.path.join(DEMO_PROJECTS_ROOT, "csharp")
CSHARP_PROJECT_DIR = os.path.join(CSHARP_DIR, "project")
CSHARP_COBERTURA_XML = os.path.join(CSHARP_DIR, "report", "cobertura", "cobertura.xml")

GO_DIR = os.path.join(DEMO_PROJECTS_ROOT, "go")
GO_PROJECT_DIR = os.path.join(GO_DIR, "project")
GO_COVERAGE_OUT = os.path.join(GO_DIR, "report", "gocover", "coverage.out")

# Self-Coverage Paths (for the --self-cover flag)
SELF_COVERAGE_DIR = os.path.join(REPORTS_OUTPUT_DIR, "adlercov_self_coverage")
UNIT_TEST_COVERAGE_OUT = os.path.join(SELF_COVERAGE_DIR, "coverage-unit.out")
INTEGRATION_TEST_COVERAGE_OUT = os.path.join(SELF_COVERAGE_DIR, "coverage-integration.out")


# ==============================================================================
#  Test Case Definitions
# ==============================================================================

@dataclass
class TestCase:
    """Defines a single E2E test case."""
    name: str
    output_dir_name: str
    args: list[str]
    expect_success: bool = True
    # List of files expected in the output dir (relative paths)
    output_files: list[str] = field(default_factory=lambda: ["index.html"])

# --- Primary E2E Test Cases ---
DEMO_PROJECT_TESTS = [
    # Individual Reports
    TestCase(
        name="C# Project Only (from Cobertura)",
        output_dir_name="csharp_cobertura_only",
        args=[
            f"-report={CSHARP_COBERTURA_XML}",
            f"-sourcedirs={CSHARP_PROJECT_DIR}",
        ],
    ),
    TestCase(
        name="Go Project Only (from gocover)",
        output_dir_name="go_gocover_only",
        args=[
            f"-report={GO_COVERAGE_OUT}",
            f"-sourcedirs={GO_PROJECT_DIR}",
        ],
    ),
    TestCase(
        name="C++ Project Only (from gcov)",
        output_dir_name="cpp_gcov_only",
        args=[
            f"-report={CPP_GCOV_PATTERN}",
            f"-sourcedirs={CPP_PROJECT_DIR}",
        ],
    ),
    TestCase(
        name="C++ Project Only (from Cobertura)",
        output_dir_name="cpp_cobertura_only",
        args=[
            f"-report={CPP_COBERTURA_XML}",
            f"-sourcedirs={CPP_PROJECT_DIR}",
        ],
    ),
    # Merged Reports
    TestCase(
        name="Merged - All Cobertura Reports",
        output_dir_name="merged_all_cobertura",
        args=[
            f"-report={CSHARP_COBERTURA_XML};{CPP_COBERTURA_XML}",
            f"-sourcedirs={CSHARP_PROJECT_DIR};{CPP_PROJECT_DIR}",
        ],
    ),
    TestCase(
        name="Merged - All C++ Reports",
        output_dir_name="merged_all_cpp",
        args=[
            f"-report={CPP_GCOV_PATTERN};{CPP_COBERTURA_XML}",
            f"-sourcedirs={CPP_PROJECT_DIR};{CPP_PROJECT_DIR}",
        ],
    ),
    TestCase(
        name="Merged - All Projects (Mixed Input Types)",
        output_dir_name="merged_all_projects_mixed",
        args=[
            f"-report={CSHARP_COBERTURA_XML};{GO_COVERAGE_OUT};{CPP_GCOV_PATTERN}",
            f"-sourcedirs={CSHARP_PROJECT_DIR};{GO_PROJECT_DIR};{CPP_PROJECT_DIR}",
        ],
    ),
    # Failure Case
    TestCase(
        name="Failure - Missing Report Argument",
        output_dir_name="failure_missing_report_arg",
        args=["-sourcedirs=."],
        expect_success=False,
        output_files=[] # No output expected on failure
    )
]

# --- Test cases for --self-cover mode ---
SELF_COVERAGE_TESTS = [
    TestCase(
        name="AdlerCov Self-Coverage (Unit + Integration Merged)",
        output_dir_name="adlercov_self_coverage_full",
        args=[
            f"-report={UNIT_TEST_COVERAGE_OUT};{INTEGRATION_TEST_COVERAGE_OUT}",
            f"-sourcedirs={SCRIPT_ROOT};{SCRIPT_ROOT}",
            "-filefilters=-**/*_test.go;-vendor/**;-tools/**"
        ],
    )
]


# ==============================================================================
#  Helper Functions
# ==============================================================================

def run_command(command, working_dir=SCRIPT_ROOT, suppress_output=True, critical=False):
    """Executes a command, returns its exit code, and optionally exits on failure."""
    print(f"--- Running Command: {' '.join(command)}")
    stdout = subprocess.DEVNULL if suppress_output else sys.stdout
    stderr = subprocess.DEVNULL if suppress_output else sys.stderr
    try:
        process = subprocess.run(
            command, cwd=working_dir, check=False, stdout=stdout, stderr=stderr
        )
        if critical and process.returncode != 0:
            print(f"--- CRITICAL COMMAND FAILED (Code: {process.returncode}). Aborting. ---", file=sys.stderr)
            sys.exit(1)
        return process.returncode
    except FileNotFoundError:
        print(f"--- Error: Command not found: {command[0]}", file=sys.stderr)
        if critical: sys.exit(1)
        return -1

def clean_directory(path):
    """Removes a directory and its contents, then recreates it."""
    if os.path.exists(path):
        shutil.rmtree(path)
    os.makedirs(path, exist_ok=True)

def get_platform_and_binary_name():
    """Determines the current OS and the corresponding binary name."""
    if sys.platform.startswith("linux"):
        return "linux", "adlercov"
    elif sys.platform == "win32":
        return "windows", "adlercov.exe"
    else:
        print(f"--- Unsupported platform: {sys.platform}", file=sys.stderr)
        sys.exit(1)

def print_summary_report(results):
    """Prints a formatted summary of all attempted tasks."""
    print("\n" + "="*80)
    print(" Final Report Generation Summary ")
    print("="*80)

    success_count, failure_count, skipped_count = 0, 0, 0

    for result in results:
        print(f"\nTask  : {result['name']}")
        print(f"Status: {result['status']}")
        print(f"Details: {result['details']}")

        if "FAILED" in result["status"]:
            failure_count += 1
        elif "SUCCESS" in result["status"]:
            success_count += 1
        else: # Skipped
            skipped_count += 1

    print("\n" + "-"*80)
    print(f"Summary: ✅ {success_count} succeeded, ❌ {failure_count} failed, ⚪ {skipped_count} skipped.")
    print("="*80)


# ==============================================================================
#  Core E2E and Self-Coverage Functions
# ==============================================================================

def build_tool(platform, cover=False):
    """Builds the Go application, optionally with coverage instrumentation."""
    action = "Building binary with coverage" if cover else "Building binary"
    print(f"\n--- {action} for {platform} ---")
    binary_path = os.path.join(BINARY_DIR, get_platform_and_binary_name()[1])
    build_cmd = ["go", "build", "-mod=vendor"]
    if cover:
        build_cmd.append("-cover")
    build_cmd.extend(["-o", binary_path, os.path.join(SCRIPT_ROOT, "cmd/main.go")])
    run_command(build_cmd, critical=True, suppress_output=False)
    print("--- Build successful ---")


def run_test_suite(test_cases, binary_path, global_args, title_prefix="E2E", verbose=False):
    """Runs a suite of test cases and returns a list of detailed result dicts."""
    results = []
    for case in test_cases:
        print(f"\n--- Running {title_prefix} Test Case: {case.name} ---")

        case_output_dir = os.path.join(REPORTS_OUTPUT_DIR, case.output_dir_name)
        os.makedirs(case_output_dir, exist_ok=True)

        command = [binary_path] + case.args + global_args + [f"-output={case_output_dir}"]
        return_code = run_command(command, suppress_output=not verbose)

        actual_success = (return_code == 0)
        test_passed = (actual_success == case.expect_success)

        if test_passed and case.expect_success:
            for file_path in case.output_files:
                if not os.path.exists(os.path.join(case_output_dir, file_path)):
                    test_passed = False
                    break

        if test_passed:
            results.append({
                "name": case.name,
                "status": "✅ SUCCESS",
                "details": f"Reports saved to '{case_output_dir}'"
            })
        else:
            results.append({
                "name": case.name,
                "status": "❌ FAILED",
                "details": f"Expected success: {case.expect_success}, but got exit code: {return_code}"
            })

    return results


def run_self_coverage_workflow(binary_path, global_args, verbose=False):
    """Handles the entire self-coverage report generation process."""
    print("\n" + "="*80)
    print("--- Starting AdlerCov Self-Coverage Workflow ---")
    print("="*80)

    clean_directory(SELF_COVERAGE_DIR)

    # 1. Run unit tests to generate the first coverage file.
    print("\n--- Step 1: Running unit tests for coverage ---")
    unit_test_cmd = ["go", "test", "-v", f"-coverprofile={UNIT_TEST_COVERAGE_OUT}", "./..."]
    run_command(unit_test_cmd, critical=True, suppress_output=not verbose)
    print(f"--- Unit test coverage saved to {UNIT_TEST_COVERAGE_OUT} ---")

    # 2. Process the raw integration coverage data.
    print("\n--- Step 2: Processing integration test coverage data ---")
    raw_cover_dir = os.environ.get("GOCOVERDIR")
    if not raw_cover_dir or not os.path.exists(raw_cover_dir):
        print("--- WARNING: GOCOVERDIR is not set or empty. ---", file=sys.stderr)
        return [{"name": "Self-Coverage Data Processing", "status": "❌ FAILED", "details": "GOCOVERDIR was not found or empty."}]

    convert_cmd = ["go", "tool", "covdata", "textfmt", f"-i={raw_cover_dir}", f"-o={INTEGRATION_TEST_COVERAGE_OUT}"]
    run_command(convert_cmd, critical=True, suppress_output=not verbose)
    print(f"--- Integration coverage saved to {INTEGRATION_TEST_COVERAGE_OUT} ---")

    # 3. Run the tool on its own coverage files.
    print("\n--- Step 3: Generating self-coverage reports ---")
    return run_test_suite(SELF_COVERAGE_TESTS, binary_path, global_args, title_prefix="Self-Cover", verbose=verbose)


def main():
    """Main function to set up, build, and run all tests."""
    parser = argparse.ArgumentParser(description="E2E test script for AdlerCov.")
    parser.add_argument("-v", "--verbose", action="store_true", help="Stream the live output from the AdlerCov tool during tests.")
    parser.add_argument("-sc", "--self-cover", action="store_true", help="Build with coverage and generate a coverage report for the tool itself.")
    parser.add_argument("--report-types", default="Html,TextSummary", help="Comma-separated list of report types to generate.")
    args = parser.parse_args()

    platform, binary_name = get_platform_and_binary_name()
    binary_path = os.path.join(BINARY_DIR, binary_name)
    global_cli_args = [f"-reporttypes={args.report_types}"]
    temp_cover_dir = None
    all_results = []

    try:
        if args.self_cover:
            temp_cover_dir = tempfile.mkdtemp(prefix="adlercov_raw_")
            os.environ["GOCOVERDIR"] = temp_cover_dir
            print(f"--- Raw coverage data will be collected in: {temp_cover_dir} ---")

        print("--- Setting up E2E test environment ---")
        clean_directory(REPORTS_OUTPUT_DIR)
        clean_directory(BINARY_DIR)
        build_tool(platform, cover=args.self_cover)

        print("\n" + "="*80)
        print("--- Running Primary E2E Tests ---")
        print("="*80)
        e2e_results = run_test_suite(DEMO_PROJECT_TESTS, binary_path, global_cli_args, verbose=args.verbose)
        all_results.extend(e2e_results)

        primary_tests_failed = any("FAILED" in r["status"] for r in e2e_results)

        if args.self_cover:
            if primary_tests_failed:
                print("\n--- SKIPPING self-coverage workflow because primary E2E tests failed. ---", file=sys.stderr)
                all_results.append({"name": "Self-Coverage Workflow", "status": "⚪ SKIPPED", "details": "Skipped due to failures in primary E2E tests."})
            else:
                self_cover_results = run_self_coverage_workflow(binary_path, global_cli_args, verbose=args.verbose)
                all_results.extend(self_cover_results)

    finally:
        if temp_cover_dir and os.path.exists(temp_cover_dir):
            print(f"--- Cleaning up temporary coverage directory: {temp_cover_dir} ---")
            shutil.rmtree(temp_cover_dir)

    print_summary_report(all_results)

    if any("FAILED" in r["status"] for r in all_results):
        print("\nOne or more tasks failed. Exiting with error status.")
        sys.exit(1)
    else:
        print("\nAll tasks completed successfully.")
        sys.exit(0)

if __name__ == "__main__":
    main()