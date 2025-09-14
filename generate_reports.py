#!/usr/bin/env python3
"""
Runs the AdlerCov CLI tool to generate many coverage reports at once

It attempts to run all enabled report tasks, provides a final summary of
which tasks succeeded and which failed along with any error details.

It assumes that coverage data files (e.g., Cobertura XML, .gcov, .out)
have already been generated inside the 'demo_projects' directory.

If the coverage output files are not yet generated, each demo_project has a
script to generate the coverage reports.
"""
import subprocess
import sys
import pathlib
import argparse
import platform
import os
import shutil
import tempfile

# ==============================================================================
#  Configuration: Paths and Tool Information
# ==============================================================================

# The root of the script is the project root.
SCRIPT_ROOT = pathlib.Path(__file__).resolve().parent
REPORTS_OUTPUT_BASE = SCRIPT_ROOT / "reports"
DEMO_PROJECTS_ROOT = SCRIPT_ROOT / "demo_projects"

# Path for the tool's own coverage report, generated when using --cover
SELF_COVERAGE_OUT = REPORTS_OUTPUT_BASE / "adlercov_coverage" / "coverage.out"

# AdlerCov Tool Location
GO_TOOL_SRC_DIR = SCRIPT_ROOT
def get_binary_name():
    """Returns the platform-specific name for the AdlerCov executable."""
    return "adlercov.exe" if platform.system() == "Windows" else "adlercov"
BINARY_NAME = get_binary_name()
BINARY_PATH = SCRIPT_ROOT / BINARY_NAME

# C++ Project Paths (inside demo_projects)
CPP_DIR = DEMO_PROJECTS_ROOT / "cpp"
CPP_PROJECT_DIR = CPP_DIR / "project"
CPP_COBERTURA_XML = CPP_DIR / "report" / "cobertura" / "cobertura.xml"
CPP_GCOV_DIR = CPP_DIR / "report" / "gcov" / "branch-probabilities"
CPP_GCOV_PATTERN = str(CPP_GCOV_DIR.resolve() / "*.gcov")

# C# Project Paths (inside demo_projects)
CSHARP_DIR = DEMO_PROJECTS_ROOT / "csharp"
CSHARP_PROJECT_DIR = CSHARP_DIR / "project"
CSHARP_COBERTURA_XML = CSHARP_DIR / "report" / "cobertura" / "cobertura.xml"

# Go Project Paths (inside demo_projects)
GO_DIR = DEMO_PROJECTS_ROOT / "go"
GO_PROJECT_DIR = GO_DIR / "project"
GO_COVERAGE_OUT = GO_DIR / "report" / "gocover" / "coverage.out"

# ==============================================================================
# Report Generation Tasks
# ==============================================================================

REPORT_TASKS = [
    # Individual Reports
    {
        "name": "C# Project Only (from Cobertura)",
        "inputs": [CSHARP_COBERTURA_XML],
        "source_dirs": [CSHARP_PROJECT_DIR],
        "output_dir_suffix": "csharp_cobertura_only",
        "enabled": True,
    },
    {
        "name": "Go Project Only (from gocover)",
        "inputs": [GO_COVERAGE_OUT],
        "source_dirs": [GO_PROJECT_DIR],
        "output_dir_suffix": "go_gocover_only",
        "enabled": True,
    },
    {
        "name": "C++ Project Only (from gcov)",
        "inputs": [CPP_GCOV_PATTERN],
        "source_dirs": [CPP_PROJECT_DIR],
        "output_dir_suffix": "cpp_gcov_only",
        "enabled": True,
    },
    {
        "name": "C++ Project Only (from Cobertura)",
        "inputs": [CPP_COBERTURA_XML],
        "source_dirs": [CPP_PROJECT_DIR],
        "output_dir_suffix": "cpp_cobertura_only",
        "enabled": True,
    },
    # Merged Reports
    {
        "name": "Merged - All Cobertura Reports",
        "inputs": [CSHARP_COBERTURA_XML, CPP_COBERTURA_XML],
        "source_dirs": [CSHARP_PROJECT_DIR, CPP_PROJECT_DIR],
        "output_dir_suffix": "merged_all_cobertura",
        "enabled": True,
    },
    {
        "name": "Merged - All C++ Reports",
        "inputs": [CPP_GCOV_PATTERN, CPP_COBERTURA_XML],
        "source_dirs": [CPP_PROJECT_DIR, CPP_PROJECT_DIR],
        "output_dir_suffix": "merged_all_cpp",
        "enabled": True,
    },
    {
        "name": "Merged - All Projects (Mixed Input Types)",
        "inputs": [CSHARP_COBERTURA_XML, GO_COVERAGE_OUT, CPP_GCOV_PATTERN],
        "source_dirs": [CSHARP_PROJECT_DIR, GO_PROJECT_DIR, CPP_PROJECT_DIR],
        "output_dir_suffix": "merged_all_projects_mixed",
        "enabled": True,
    }
]


def run_command(cmd, working_dir=None, critical=False):
    """
    Executes a command, streams its output in real-time, and returns the result.
    If 'critical' is True, the script will exit immediately on failure.
    """
    print(f"\n>>> Running Command: {' '.join(map(str, cmd))}")
    full_output = []

    try:
        process = subprocess.Popen(
            cmd,
            cwd=working_dir,
            stdout=subprocess.PIPE,
            stderr=subprocess.STDOUT,
            text=True,
            bufsize=1,
            universal_newlines=True,
            encoding='utf-8',
            errors='replace'
        )

        # Read and print output line-by-line in real-time
        for line in process.stdout:
            print(line, end='')
            full_output.append(line)

        process.wait()
        returncode = process.returncode
        output_str = "".join(full_output).strip()

        if returncode != 0:
            if critical:
                print(f"\n❌ Critical command failed with return code {returncode}.")
                sys.exit(1)
            return {"success": False, "output": output_str}

        return {"success": True, "output": output_str}

    except FileNotFoundError:
        message = f"❌ Command not found: {cmd[0]}"
        print(message)
        if critical: sys.exit(1)
        return {"success": False, "output": message}
    except Exception as e:
        message = f"❌ An unexpected error occurred: {e}"
        print(message)
        if critical: sys.exit(1)
        return {"success": False, "output": message}


def build_adlercov_binary(cover=False):
    """
    Builds the AdlerCov Go binary.
    If 'cover' is True, builds with coverage instrumentation.
    """
    if cover:
        print(" Building AdlerCov binary with coverage instrumentation ")
    else:
        print(" Building AdlerCov binary ")

    build_cmd = ["go", "build", "-mod=vendor"]
    if cover:
        build_cmd.append("-cover")

    build_cmd.extend(["-o", str(BINARY_PATH), "cmd/main.go"])

    run_command(build_cmd, working_dir=SCRIPT_ROOT, critical=True)
    print(f"✅ Successfully built '{BINARY_NAME}'")

def convert_go_coverage(raw_data_dir, output_file):
    """
    Merges raw Go coverage data files into a single coverage.out file.
    Returns True on success, False on failure.
    """
    print("\n" + "-"*80)
    print(" Converting Go Coverage Data for AdlerCov ")
    print(f"  Raw data source: '{raw_data_dir}'")
    print(f"  Output file: '{output_file}'")

    output_file.parent.mkdir(parents=True, exist_ok=True)

    cmd = [
        "go",
        "tool",
        "covdata",
        "textfmt",
        f"-i={raw_data_dir}",
        f"-o={str(output_file.resolve())}",
    ]

    result = run_command(cmd, working_dir=SCRIPT_ROOT)

    if result["success"]:
        print(f"✅ Successfully converted coverage data to '{output_file}'")
    else:
        print("❌ Failed to convert coverage data.")
    return result["success"]

def generate_reports(tasks_to_run, report_types):
    """Iterates through tasks, runs them, and collects the results."""
    print("\n" + "="*80)
    print(" Starting Report Generation ")
    print("="*80)
    results = []

    for task in tasks_to_run:
        task_name = task['name']
        if not task.get("enabled", True):
            results.append({"name": task_name, "status": "⚪ SKIPPED", "details": "Task was disabled in the script."})
            continue

        print(f"\n Processing Task: {task_name} ")

        report_patterns = [str(p) for p in task["inputs"]]

        output_dir = REPORTS_OUTPUT_BASE / task["output_dir_suffix"]
        output_dir.mkdir(parents=True, exist_ok=True)

        cmd = [
            str(BINARY_PATH),
            f"--report={';'.join(report_patterns)}",
            f"--output={str(output_dir.resolve())}",
            f"--reporttypes={report_types}",
            "--verbose"
        ]
        if task["source_dirs"]:
            source_paths = [str(p.resolve()) for p in task["source_dirs"]]
            cmd.append(f"--sourcedirs={';'.join(source_paths)}")

        result = run_command(cmd)
        if result["success"]:
            results.append({"name": task_name, "status": "✅ SUCCESS", "details": f"Reports saved to '{output_dir}'"})
        else:
            results.append({"name": task_name, "status": "❌ FAILED", "details": result["output"]})

    return results

def print_summary_report(results):
    """Prints a formatted summary of all attempted tasks."""
    print("\n" + "="*80)
    print(" Final Report Generation Summary ")
    print("="*80)

    success_count, failure_count, skipped_count = 0, 0, 0

    for result in results:
        print(f"\nTask  : {result['name']}")
        print(f"Status: {result['status']}")

        if "FAILED" in result["status"]:
            failure_count += 1
            details = "  " + result["details"].replace("\n", "\n  ")
            print(f"Details:\n{details}")
        elif "SUCCESS" in result["status"]:
            success_count += 1
        else: # Skipped
            skipped_count += 1
            print(f"Reason: {result['details']}")

    print("\n" + "-"*80)
    print(f"Summary: ✅ {success_count} succeeded, ❌ {failure_count} failed, 🟡 {skipped_count} skipped.")
    print("="*80)

def main():
    """Main execution function."""
    parser = argparse.ArgumentParser(description="Generate coverage reports using the AdlerCov tool.")
    parser.add_argument("--build", action="store_true", help="Force a rebuild of the AdlerCov binary.")
    parser.add_argument("--report-types", default="Html,TextSummary,Lcov,RawJson", help="Comma-separated report types.")
    parser.add_argument("--cover", action="store_true", help="Build with coverage instrumentation and generate a coverage report for the tool itself.")
    args = parser.parse_args()

    # Setup for coverage run
    temp_cover_dir = None
    if args.cover:
        temp_cover_dir = tempfile.mkdtemp(prefix="adlercov_raw_")
        os.environ["GOCOVERDIR"] = temp_cover_dir
        print(f"Coverage data will be collected in: {temp_cover_dir}")

    try:
        if args.build or not BINARY_PATH.exists() or args.cover:
            build_adlercov_binary(cover=args.cover)
        else:
            print(f"Using existing binary: {BINARY_PATH}")

        results = []
        tasks_failed = False
        try:
            results = generate_reports(REPORT_TASKS, args.report_types)
        finally:
            tasks_failed = any("FAILED" in r["status"] for r in results)

            # If coverage was enabled, process the tool's own coverage report
            if args.cover:
                if not tasks_failed:
                    conversion_ok = convert_go_coverage(temp_cover_dir, SELF_COVERAGE_OUT)

                    if conversion_ok:
                        # Define the self-coverage task and run it
                        self_coverage_task = {
                            "name": "AdlerCov Self-Coverage Report (from .out)",
                            "inputs": [SELF_COVERAGE_OUT],
                            "source_dirs": [GO_TOOL_SRC_DIR],
                            "output_dir_suffix": "adlercov_self_coverage",
                            "enabled": True,
                        }
                        self_results = generate_reports([self_coverage_task], args.report_types)
                        results.extend(self_results)
                    else:
                        # Append a failure result if conversion failed
                        results.append({
                            "name": "AdlerCov Self-Coverage Report",
                            "status": "❌ FAILED",
                            "details": "Failed to convert raw Go coverage data. See logs."
                        })
                else:
                    print("\n⚪ SKIPPING self-coverage report generation because one or more primary tasks failed.")
                    results.append({
                        "name": "AdlerCov Self-Coverage Report",
                        "status": "⚪ SKIPPED",
                        "details": "Skipped because one or more primary tasks failed."
                    })

            # Re-evaluate failure status before final summary
            tasks_failed = any("FAILED" in r["status"] for r in results)
            print_summary_report(results)
            if tasks_failed:
                print("\nOne or more tasks failed. Exiting with error status.")
                sys.exit(1)
    finally:
        # Cleanup the temporary directory for coverage data
        if temp_cover_dir:
            print(f"Cleaning up temporary coverage directory: {temp_cover_dir}")
            shutil.rmtree(temp_cover_dir)


if __name__ == "__main__":
    main()