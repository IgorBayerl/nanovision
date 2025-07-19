#!/usr/bin/env python3
"""
Universal Coverage Reporter Script

This script automates the process of generating code coverage reports for both C# (.NET)
and Go projects. It performs the following main steps:

1.  **For C# Projects**:
    *   Runs `dotnet test` to generate Cobertura XML coverage files.
    *   Uses a custom Go-based report generator (`go_report_generator`) to create reports
        from the Cobertura XML.
    *   Uses the standard .NET `ReportGenerator` tool to create reports from the
        Cobertura XML.

2.  **For Go Projects**:
    *   Runs `go test` to generate native Go coverage profiles (`coverage.out`).
    *   Uses the custom Go tool to generate a report directly from `coverage.out`.
    *   Converts `coverage.out` to Cobertura XML format using `gocover-cobertura`.
    *   Uses the custom Go tool to generate a report from the converted Cobertura XML.
    *   Uses the standard .NET `ReportGenerator` tool to create reports from the Cobertura XML.

3.  **For Merged Reports**:
    *   Takes the C# Cobertura XML and the Go native coverage file as inputs.
    *   Uses the custom Go tool to generate a single, unified report combining the
        coverage data from both projects.

The script allows specifying desired report types by modifying the
`SELECTED_REPORT_TYPES_CONFIG_STRING` variable in the `main` function.

Prerequisites:
    *   .NET SDK (for `dotnet test` and `ReportGenerator.dll`)
    *   Go (for `go test`, `go run`, and `gocover-cobertura`)
    *   `gocover-cobertura` executable in your PATH or accessible.
        (Install via: go install github.com/t-yuki/gocover-cobertura@latest)
    *   Node.js and npm to build the Angular SPA for the HTML reports. Before running,
        ensure you have run `npm install` and `npm run build` inside the
        `go_report_generator/angular_frontend_spa/` directory.

Usage Example:
    1.  Edit `SELECTED_REPORT_TYPES_CONFIG_STRING` in `main()` to define report types.
    2.  Run the script from the `Testprojects` directory:
        ```bash
        python generate_reports.py
        ```
"""
import subprocess
import sys
import pathlib
import os

# --- Constants for Paths ---
SCRIPT_ROOT = pathlib.Path(__file__).resolve().parent

# --- C# Project Specific Paths ---
CSHARP_TEST_PROJECT_PATH = SCRIPT_ROOT / "CSharp/Project_DotNetCore/UnitTests/UnitTests.csproj"
CSHARP_COVERAGE_OUTPUT_DIR = SCRIPT_ROOT / "CSharp/Reports"
CSHARP_COBERTURA_XML_PATH = CSHARP_COVERAGE_OUTPUT_DIR / "coverage.cobertura.xml"
CSHARP_REPORTS_FROM_GO_TOOL_DIR = SCRIPT_ROOT.parent / "reports/csharp_project_go_tool_report"
CSHARP_REPORTS_FROM_DOTNET_TOOL_DIR = SCRIPT_ROOT.parent / "reports/csharp_project_dotnet_tool_report"

# --- Go Project Specific Paths ---
GO_PROJECT_TO_TEST_PATH = SCRIPT_ROOT / "Go"
GO_PROJECT_NATIVE_COVERAGE_FILE = GO_PROJECT_TO_TEST_PATH / "coverage.out"
GO_PROJECT_COBERTURA_XML_FILE = GO_PROJECT_TO_TEST_PATH / "coverage.cobertura.xml"
GO_PROJECT_REPORTS_FROM_GO_TOOL_NATIVE_DIR = SCRIPT_ROOT.parent / "reports/go_project_go_tool_native_report"
GO_PROJECT_REPORTS_FROM_GO_TOOL_COBERTURA_DIR = SCRIPT_ROOT.parent / "reports/go_project_go_tool_cobertura_report"
GO_PROJECT_REPORTS_FROM_DOTNET_TOOL_DIR = SCRIPT_ROOT.parent / "reports/go_project_dotnet_tool_report"

# --- Merged Project Paths ---
MERGED_REPORT_DIR = SCRIPT_ROOT.parent / "reports/merged_csharp_go_report"

# --- Common Tool Paths ---
GO_REPORT_GENERATOR_CMD_PATH = SCRIPT_ROOT.parent / "go_report_generator/cmd"
DOTNET_REPORT_GENERATOR_DLL_PATH = SCRIPT_ROOT.parent / "src/ReportGenerator.Console.NetCore/bin/Debug/net8.0/ReportGenerator.dll"

# --- Report File Names (for verification) ---
TEXT_SUMMARY_FILE_NAME = "Summary.txt"
HTML_REPORT_INDEX_FILE_NAME = "index.html"
LCOV_REPORT_FILE_NAME = "lcov.info"

# --- End of Constants ---

def run_command(command_args_or_string, working_dir=None, command_name="Command", shell=False):
    """Runs a shell command, prints output, and exits on error."""
    is_string_command = isinstance(command_args_or_string, str)

    if shell and not is_string_command:
        print(f"Error: If shell=True, command must be a string. Got: {command_args_or_string}", file=sys.stderr)
        sys.exit(1)

    cmd_display_str = command_args_or_string if is_string_command else ' '.join(map(str, command_args_or_string))
    print(f"Executing {command_name}: {cmd_display_str[:120]}{'...' if len(cmd_display_str) > 120 else ''}")
    if working_dir:
        print(f"  (in {working_dir})")
    try:
        process = subprocess.run(
            command_args_or_string,
            capture_output=True,
            text=True,
            cwd=working_dir,
            check=False,
            shell=shell,
            env=os.environ.copy()
        )

        if process.stdout:
            print(f"  Stdout from {command_name}:\n{process.stdout.strip()}")
        if process.stderr:
            print(f"  Stderr from {command_name}:\n{process.stderr.strip()}", file=sys.stderr)

        if process.returncode != 0:
            print(f"Error executing {command_name} (Return code: {process.returncode})", file=sys.stderr)
            sys.exit(1)

        return process
    except FileNotFoundError:
        executable = command_args_or_string if shell else command_args_or_string[0]
        print(f"Error: Command not found - {executable}. Ensure it's in your PATH or correctly specified.", file=sys.stderr)
        sys.exit(1)
    except Exception as e:
        print(f"An unexpected error occurred during {command_name}: {e}", file=sys.stderr)
        sys.exit(1)

def ensure_dir(dir_path: pathlib.Path):
    """Creates a directory if it doesn't exist."""
    try:
        dir_path.mkdir(parents=True, exist_ok=True)
        print(f"Directory ensured: {dir_path}")
    except OSError as e:
        print(f"Error creating directory {dir_path}: {e}", file=sys.stderr)
        sys.exit(1)

def check_generated_files(output_dir: pathlib.Path, report_types: list[str], tool_name: str) -> bool:
    """Checks if expected report files were generated based on report_types."""
    all_ok = True
    checked_something = False

    print(f"Verifying generated reports in {output_dir} for {tool_name} with types: {', '.join(report_types)}...")

    if not output_dir.is_dir():
        print(f"Error: Output directory {output_dir} for {tool_name} does not exist.", file=sys.stderr)
        return False

    if "TextSummary" in report_types:
        checked_something = True
        summary_file = output_dir / TEXT_SUMMARY_FILE_NAME
        if not summary_file.exists() or summary_file.stat().st_size == 0:
            print(f"Error: {tool_name} TextSummary report not generated or is empty at {summary_file}", file=sys.stderr)
            all_ok = False
        else:
            print(f"  OK: {tool_name} TextSummary report generated: {summary_file}")

    html_report_types_keywords = {"Html", "HtmlInline", "HtmlChart", "HtmlSummary", "Html_Dark"}
    if any(rt_keyword in rt for rt_keyword in html_report_types_keywords for rt in report_types):
        checked_something = True
        index_html_file = output_dir / HTML_REPORT_INDEX_FILE_NAME
        if not index_html_file.exists() or index_html_file.stat().st_size == 0:
            print(f"Error: {tool_name} HTML report (index.html) not generated or is empty at {index_html_file}", file=sys.stderr)
            all_ok = False
        else:
            print(f"  OK: {tool_name} HTML report (index.html) generated: {index_html_file}")

    if "lcov" in report_types:
        checked_something = True
        lcov_file = output_dir / LCOV_REPORT_FILE_NAME
        if not lcov_file.exists() or lcov_file.stat().st_size == 0:
            print(f"Error: {tool_name} LCOV report not generated or is empty at {lcov_file}", file=sys.stderr)
            all_ok = False
        else:
            print(f"  OK: {tool_name} LCOV report generated: {lcov_file}")

    if not checked_something and report_types:
        print(f"Warning: No specific file checks implemented for configured report types: {', '.join(report_types)} by {tool_name}. Assuming success if command ran.", file=sys.stderr)
        return True

    return all_ok


def run_csharp_workflow(report_types_list: list[str]):
    """Generates C# coverage and creates reports using both Go and .NET tools."""
    print("\n--- Starting C# Project Workflow ---")
    for dir_path in [CSHARP_COVERAGE_OUTPUT_DIR, CSHARP_REPORTS_FROM_GO_TOOL_DIR, CSHARP_REPORTS_FROM_DOTNET_TOOL_DIR]:
        ensure_dir(dir_path)

    print("\n--- Generating Cobertura XML for C# project ---")
    dotnet_test_command = [
        "dotnet", "test", str(CSHARP_TEST_PROJECT_PATH),
        "--configuration", "Release", "--verbosity", "minimal",
        "/p:CollectCoverage=true", "/p:CoverletOutputFormat=cobertura",
        f"/p:CoverletOutput={CSHARP_COBERTURA_XML_PATH.resolve()}"
    ]
    run_command(dotnet_test_command, command_name="dotnet test (C#)")
    if not (CSHARP_COBERTURA_XML_PATH.exists() and CSHARP_COBERTURA_XML_PATH.stat().st_size > 0):
        print(f"Error: C# Cobertura XML not generated or is empty at {CSHARP_COBERTURA_XML_PATH}", file=sys.stderr)
        sys.exit(1)
    print(f"C# Cobertura XML generated: {CSHARP_COBERTURA_XML_PATH}")

    if GO_REPORT_GENERATOR_CMD_PATH.is_dir():
        print("\n--- Generating C# report with Go tool ---")
        go_tool_report_types_arg = ",".join(report_types_list)
        go_report_command_csharp = [
            "go", "run", ".",
            f"-report={CSHARP_COBERTURA_XML_PATH.resolve()}",
            f"-output={CSHARP_REPORTS_FROM_GO_TOOL_DIR.resolve()}",
            f"-reporttypes={go_tool_report_types_arg}"
        ]
        run_command(go_report_command_csharp, working_dir=GO_REPORT_GENERATOR_CMD_PATH, command_name="Go Report Generator (for C#)")
        if not check_generated_files(CSHARP_REPORTS_FROM_GO_TOOL_DIR, report_types_list, "C# Go-tool"):
            sys.exit("Error: C# workflow Go tool report verification failed.")
    else:
        print("Skipping Go Report Generator for C#: Tool not found.")

    if DOTNET_REPORT_GENERATOR_DLL_PATH.exists():
        print("\n--- Generating C# report with .NET tool ---")
        dotnet_tool_report_types_arg = ";".join(report_types_list)
        dotnet_rg_command_csharp = [
            "dotnet", str(DOTNET_REPORT_GENERATOR_DLL_PATH.resolve()),
            f"-reports:{CSHARP_COBERTURA_XML_PATH.resolve()}",
            f"-targetdir:{CSHARP_REPORTS_FROM_DOTNET_TOOL_DIR.resolve()}",
            f"-reporttypes:{dotnet_tool_report_types_arg}"
        ]
        run_command(dotnet_rg_command_csharp, command_name=".NET ReportGenerator (for C#)")
        if not check_generated_files(CSHARP_REPORTS_FROM_DOTNET_TOOL_DIR, report_types_list, "C# .NET-tool"):
            sys.exit("Error: C# workflow .NET tool report verification failed.")
    else:
        print("Skipping .NET ReportGenerator for C#: Tool not found.")

    print("--- C# Project Workflow Finished Successfully ---")


def run_go_project_workflow(report_types_list: list[str]):
    """Generates Go coverage and creates reports using multiple methods for comparison."""
    print("\n--- Starting Go Project Workflow ---")
    for dir_path in [
        GO_PROJECT_REPORTS_FROM_GO_TOOL_NATIVE_DIR,
        GO_PROJECT_REPORTS_FROM_GO_TOOL_COBERTURA_DIR,
        GO_PROJECT_REPORTS_FROM_DOTNET_TOOL_DIR
    ]:
        ensure_dir(dir_path)

    if not GO_PROJECT_TO_TEST_PATH.is_dir():
        sys.exit(f"Error: Go project to test not found at {GO_PROJECT_TO_TEST_PATH}")
    print(f"Go project directory found: {GO_PROJECT_TO_TEST_PATH}")

    GO_PROJECT_NATIVE_COVERAGE_FILE.unlink(missing_ok=True)
    GO_PROJECT_COBERTURA_XML_FILE.unlink(missing_ok=True)
    print("Old Go coverage files removed.")

    print("\n--- Step 1: Generating native Go coverage ---")
    go_test_command = ["go", "test", f"-coverprofile={GO_PROJECT_NATIVE_COVERAGE_FILE.name}", "./..."]
    run_command(go_test_command, working_dir=GO_PROJECT_TO_TEST_PATH, command_name="go test (Go project)")
    if not (GO_PROJECT_NATIVE_COVERAGE_FILE.exists() and GO_PROJECT_NATIVE_COVERAGE_FILE.stat().st_size > 0):
        sys.exit(f"Error: Go native coverage not generated or is empty at {GO_PROJECT_NATIVE_COVERAGE_FILE}")
    print(f"Go native coverage generated: {GO_PROJECT_NATIVE_COVERAGE_FILE}")

    go_tool_report_types_arg = ",".join(report_types_list)

    if GO_REPORT_GENERATOR_CMD_PATH.is_dir():
        print("\n--- Step 2: Generating report with Go tool (from native .out file) ---")
        go_report_command_go_native = [
            "go", "run", ".",
            f"-report={GO_PROJECT_NATIVE_COVERAGE_FILE.resolve()}",
            f"-output={GO_PROJECT_REPORTS_FROM_GO_TOOL_NATIVE_DIR.resolve()}",
            f"-reporttypes={go_tool_report_types_arg}",
            f"-sourcedirs={GO_PROJECT_TO_TEST_PATH.resolve()}"
        ]
        run_command(go_report_command_go_native, working_dir=GO_REPORT_GENERATOR_CMD_PATH, command_name="Go Report Generator (for Go native)")
        if not check_generated_files(GO_PROJECT_REPORTS_FROM_GO_TOOL_NATIVE_DIR, report_types_list, "Go-project Go-tool (Native)"):
            sys.exit("Error: Go project workflow native Go tool report verification failed.")

        print("\n--- Step 3: Converting Go native coverage to Cobertura XML ---")
        gocover_cobertura_command_str = f'gocover-cobertura < "{GO_PROJECT_NATIVE_COVERAGE_FILE.name}" > "{GO_PROJECT_COBERTURA_XML_FILE.name}"'
        run_command(gocover_cobertura_command_str, working_dir=GO_PROJECT_TO_TEST_PATH, command_name="gocover-cobertura", shell=True)
        if not (GO_PROJECT_COBERTURA_XML_FILE.exists() and GO_PROJECT_COBERTURA_XML_FILE.stat().st_size > 0):
            sys.exit(f"Error: Go Cobertura XML not generated or is empty at {GO_PROJECT_COBERTURA_XML_FILE}")
        print(f"Go project Cobertura XML generated: {GO_PROJECT_COBERTURA_XML_FILE}")

        print("\n--- Step 4: Generating report with Go tool (from Cobertura XML) ---")
        go_report_command_go_cobertura = [
            "go", "run", ".",
            f"-report={GO_PROJECT_COBERTURA_XML_FILE.resolve()}",
            f"-output={GO_PROJECT_REPORTS_FROM_GO_TOOL_COBERTURA_DIR.resolve()}",
            f"-reporttypes={go_tool_report_types_arg}"
        ]
        run_command(go_report_command_go_cobertura, working_dir=GO_REPORT_GENERATOR_CMD_PATH, command_name="Go Report Generator (for Go Cobertura)")
        if not check_generated_files(GO_PROJECT_REPORTS_FROM_GO_TOOL_COBERTURA_DIR, report_types_list, "Go-project Go-tool (Cobertura)"):
            sys.exit("Error: Go project workflow Cobertura Go tool report verification failed.")
    else:
        print("Skipping Go Report Generator for Go Project: Tool not found.")

    if DOTNET_REPORT_GENERATOR_DLL_PATH.exists():
        print("\n--- Step 5: Generating report with .NET tool (from Cobertura XML) ---")
        dotnet_tool_report_types_arg = ";".join(report_types_list)
        dotnet_rg_command_go_proj = [
            "dotnet", str(DOTNET_REPORT_GENERATOR_DLL_PATH.resolve()),
            f"-reports:{GO_PROJECT_COBERTURA_XML_FILE.resolve()}",
            f"-targetdir:{GO_PROJECT_REPORTS_FROM_DOTNET_TOOL_DIR.resolve()}",
            f"-reporttypes:{dotnet_tool_report_types_arg}"
        ]
        run_command(dotnet_rg_command_go_proj, command_name=".NET ReportGenerator (for Go project)")
        if not check_generated_files(GO_PROJECT_REPORTS_FROM_DOTNET_TOOL_DIR, report_types_list, "Go-project .NET-tool"):
            sys.exit("Error: Go project workflow .NET tool report verification failed.")
    else:
        print("Skipping .NET ReportGenerator for Go Project: Tool not found.")

    print("--- Go Project Workflow Finished Successfully ---")

def run_merged_workflow(report_types_list: list[str]):
    """Generates a single merged report from C# and Go coverage files."""
    print("\n--- Starting Merged Project Workflow ---")
    ensure_dir(MERGED_REPORT_DIR)

    # Check if input files from previous workflows exist
    if not CSHARP_COBERTURA_XML_PATH.exists():
        print(f"Skipping merged report: C# Cobertura XML not found at {CSHARP_COBERTURA_XML_PATH}", file=sys.stderr)
        return
    if not GO_PROJECT_NATIVE_COVERAGE_FILE.exists():
        print(f"Skipping merged report: Go native coverage not found at {GO_PROJECT_NATIVE_COVERAGE_FILE}", file=sys.stderr)
        return
    print("Source coverage files for merged report found.")

    print("\n--- Generating merged report from C# Cobertura and Go native files ---")
    
    # The Go tool accepts multiple reports separated by a semicolon
    merged_reports_arg = f"{CSHARP_COBERTURA_XML_PATH.resolve()};{GO_PROJECT_NATIVE_COVERAGE_FILE.resolve()}"
    go_tool_report_types_arg = ",".join(report_types_list)

    go_report_command_merged = [
        "go", "run", ".",
        f"-report={merged_reports_arg}",
        f"-output={MERGED_REPORT_DIR.resolve()}",
        f"-reporttypes={go_tool_report_types_arg}",
        # Provide source directory for the Go part of the report.
        # The Cobertura report should have its source paths embedded.
        f"-sourcedirs={GO_PROJECT_TO_TEST_PATH.resolve()}"
    ]
    
    run_command(go_report_command_merged, working_dir=GO_REPORT_GENERATOR_CMD_PATH, command_name="Go Report Generator (for Merged Report)")
    
    if not check_generated_files(MERGED_REPORT_DIR, report_types_list, "Merged Report Go-tool"):
        sys.exit("Error: Merged report verification failed.")

    print("--- Merged Project Workflow Finished Successfully ---")

def main():
    """Main function to orchestrate C#, Go, and merged project coverage and reporting."""
    print("Python script for C# and Go project coverage and reporting.")

    # --- Define desired report types here ---
    SELECTED_REPORT_TYPES_CONFIG_STRING = "Html,TextSummary,Lcov"

    if not SELECTED_REPORT_TYPES_CONFIG_STRING or SELECTED_REPORT_TYPES_CONFIG_STRING.isspace():
        sys.exit("Error: SELECTED_REPORT_TYPES_CONFIG_STRING is empty.")

    active_report_types = [rt.strip() for rt in SELECTED_REPORT_TYPES_CONFIG_STRING.split(',') if rt.strip()]
    if not active_report_types:
        sys.exit(f"Error: No valid report types parsed from '{SELECTED_REPORT_TYPES_CONFIG_STRING}'.")

    print(f"Target report types: {', '.join(active_report_types)}")

    # --- Pre-flight checks for tools ---
    if not (GO_REPORT_GENERATOR_CMD_PATH.is_dir() and (GO_REPORT_GENERATOR_CMD_PATH / "main.go").is_file()):
        print(f"Warning: Go Report Generator not found in {GO_REPORT_GENERATOR_CMD_PATH}. Go tool reporting will be skipped.", file=sys.stderr)
    else:
        print(f"Go Report Generator found at: {GO_REPORT_GENERATOR_CMD_PATH}")

    if not DOTNET_REPORT_GENERATOR_DLL_PATH.exists():
        print(f"Warning: .NET ReportGenerator.dll not found: {DOTNET_REPORT_GENERATOR_DLL_PATH}. .NET tool reporting will be skipped.", file=sys.stderr)
    else:
        print(f".NET ReportGenerator.dll found at: {DOTNET_REPORT_GENERATOR_DLL_PATH}")

    # --- Execute Workflows ---
    try:
        run_csharp_workflow(report_types_list=active_report_types)
        run_go_project_workflow(report_types_list=active_report_types)
        run_merged_workflow(report_types_list=active_report_types)
        print("\nAll workflows completed successfully. Reports generated!")
    except SystemExit as e:
        print(f"\nScript terminated prematurely with exit code {e.code}.", file=sys.stderr)
    except Exception as e:
        print(f"\nAn unexpected error occurred in main: {e}", file=sys.stderr)
        sys.exit(1)


if __name__ == "__main__":
    main()