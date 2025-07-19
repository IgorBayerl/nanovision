#!/usr/bin/env python3
"""
Generate coverage reports from existing data or regenerate everything with --force.
"""
import subprocess
import sys
import pathlib
import argparse

# Paths
SCRIPT_ROOT = pathlib.Path(__file__).resolve().parent
CSHARP_TEST_PROJECT = SCRIPT_ROOT / "CSharp/Project_DotNetCore/UnitTests/UnitTests.csproj"
CSHARP_COVERAGE_DIR = SCRIPT_ROOT / "CSharp/Reports"
CSHARP_COBERTURA_XML = CSHARP_COVERAGE_DIR / "coverage.cobertura.xml"
GO_PROJECT_DIR = SCRIPT_ROOT / "Go"
GO_COVERAGE_FILE = GO_PROJECT_DIR / "coverage.out"
GO_COBERTURA_XML = GO_PROJECT_DIR / "coverage.cobertura.xml"
GO_TOOL_CMD_DIR = SCRIPT_ROOT.parent / "cmd"

def run_command(cmd, working_dir=None, show_output=False):
    """Run command and exit on failure."""
    try:
        if show_output:
            # Stream output in real-time
            process = subprocess.Popen(
                cmd, 
                cwd=working_dir,
                stdout=subprocess.PIPE,
                stderr=subprocess.STDOUT,
                text=True,
                bufsize=1,
                universal_newlines=True
            )
            
            # Print output line by line as it comes
            for line in process.stdout:
                print(line.rstrip())
            
            process.wait()
            if process.returncode != 0:
                print(f"Command failed with return code {process.returncode}")
                sys.exit(1)
        else:
            # Capture output
            result = subprocess.run(cmd, cwd=working_dir, check=True, 
                                  capture_output=True, text=True)
        return None
    except FileNotFoundError:
        print(f"Command not found: {cmd[0]}")
        sys.exit(1)

def ensure_dir(path):
    """Create directory if it doesn't exist."""
    path.mkdir(parents=True, exist_ok=True)

def check_existing_data():
    """Check if coverage data already exists."""
    csharp_exists = CSHARP_COBERTURA_XML.exists() and CSHARP_COBERTURA_XML.stat().st_size > 0
    go_exists = GO_COVERAGE_FILE.exists() and GO_COVERAGE_FILE.stat().st_size > 0
    
    return csharp_exists, go_exists

def generate_csharp_coverage():
    """Generate C# coverage data."""
    print("Generating C# coverage...")
    ensure_dir(CSHARP_COVERAGE_DIR)
    
    cmd = [
        "dotnet", "test", str(CSHARP_TEST_PROJECT),
        "--configuration", "Release",
        "/p:CollectCoverage=true", 
        "/p:CoverletOutputFormat=cobertura",
        f"/p:CoverletOutput={CSHARP_COBERTURA_XML.resolve()}"
    ]
    run_command(cmd)
    
    if not CSHARP_COBERTURA_XML.exists():
        print("Failed to generate C# coverage file")
        sys.exit(1)

def generate_go_coverage():
    """Generate Go coverage data."""
    print("Generating Go coverage...")
    
    # Clean old files
    GO_COVERAGE_FILE.unlink(missing_ok=True)
    GO_COBERTURA_XML.unlink(missing_ok=True)
    
    # Generate native coverage
    cmd = ["go", "test", f"-coverprofile={GO_COVERAGE_FILE.name}", "./..."]
    run_command(cmd, working_dir=GO_PROJECT_DIR)
    
    if not GO_COVERAGE_FILE.exists():
        print("Failed to generate Go coverage file")
        sys.exit(1)
    
    # Convert to Cobertura XML
    cmd_str = f'gocover-cobertura < "{GO_COVERAGE_FILE.name}" > "{GO_COBERTURA_XML.name}"'
    subprocess.run(cmd_str, cwd=GO_PROJECT_DIR, shell=True, check=True)

def run_report_tool(report_types="Html,TextSummary,Lcov"):
    """Run the main report generation tool."""
    print("Running report tool...")
    
    if not GO_TOOL_CMD_DIR.exists():
        print(f"Report tool not found at {GO_TOOL_CMD_DIR}")
        sys.exit(1)
    
    # Reports output directories
    reports_base = SCRIPT_ROOT.parent / "reports"
    csharp_reports = reports_base / "csharp_reports"
    go_reports = reports_base / "go_reports" 
    merged_reports = reports_base / "merged_reports"
    
    for report_dir in [csharp_reports, go_reports, merged_reports]:
        ensure_dir(report_dir)
    
    # Generate C# report
    if CSHARP_COBERTURA_XML.exists():
        print("Generating C# report...")
        cmd = [
            "go", "run", ".",
            f"--report={CSHARP_COBERTURA_XML.resolve()}",
            f"--output={csharp_reports.resolve()}",
            f"--reporttypes={report_types}"
            f"--verbose"
        ]
        run_command(cmd, working_dir=GO_TOOL_CMD_DIR, show_output=True)
    else:
        print("C# coverage file not found, skipping C# report")
    
    # Generate Go report
    if GO_COVERAGE_FILE.exists():
        print("Generating Go report...")
        cmd = [
            "go", "run", ".",
            f"--report={GO_COVERAGE_FILE.resolve()}",
            f"--output={go_reports.resolve()}",
            f"--reporttypes={report_types}",
            f"--sourcedirs={GO_PROJECT_DIR.resolve()}"
            f"--verbose"
        ]
        run_command(cmd, working_dir=GO_TOOL_CMD_DIR, show_output=True)
    else:
        print("Go coverage file not found, skipping Go report")
    
    # Generate merged report
    if CSHARP_COBERTURA_XML.exists() and GO_COVERAGE_FILE.exists():
        print("Generating merged report...")
        merged_input = f"{CSHARP_COBERTURA_XML.resolve()};{GO_COVERAGE_FILE.resolve()}"
        cmd = [
            "go", "run", ".",
            f"--report={merged_input}",
            f"--output={merged_reports.resolve()}",
            f"--reporttypes={report_types}",
            f"--sourcedirs={GO_PROJECT_DIR.resolve()}"
            f"--verbose"
        ]
        run_command(cmd, working_dir=GO_TOOL_CMD_DIR, show_output=True)
    else:
        print("Missing coverage files, skipping merged report")

def main():
    """Main function with force flag support."""
    parser = argparse.ArgumentParser(description="Generate coverage reports")
    parser.add_argument("--force", action="store_true", 
                       help="Force regeneration of all coverage data")
    parser.add_argument("--report-types", default="Html,TextSummary,Lcov",
                       help="Comma-separated list of report types")
    
    args = parser.parse_args()
    
    if args.force:
        print("Force mode: regenerating all coverage data")
        generate_csharp_coverage()
        generate_go_coverage()
    else:
        # Check existing data
        csharp_exists, go_exists = check_existing_data()
        
        if not csharp_exists and not go_exists:
            print("No existing coverage data found, generating all...")
            generate_csharp_coverage()
            generate_go_coverage()
        elif not csharp_exists:
            print("Missing C# coverage, generating...")
            generate_csharp_coverage()
        elif not go_exists:
            print("Missing Go coverage, generating...")
            generate_go_coverage()
        else:
            print("Using existing coverage data")
    
    run_report_tool(args.report_types)
    print("Report generation complete")

if __name__ == "__main__":
    main()