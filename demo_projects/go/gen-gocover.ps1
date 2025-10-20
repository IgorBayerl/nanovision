# gen-gocover.ps1 - Runs Go tests and generates a raw coverage profile (coverage.out).

Write-Host "--- Generate Go Coverage Profile ---" -ForegroundColor Green

# Path Setup
$ScriptDir = Split-Path -Parent $MyInvocation.MyCommand.Path
$ProjectDir = Join-Path $ScriptDir "project"
$ReportBaseDir = Join-Path $ScriptDir "report"
$GoCoverDir = Join-Path $ReportBaseDir "gocover"

# Define the full path for the output profile file
$CoverageProfile = Join-Path $GoCoverDir "coverage.out"

# Cleaning up previous reports
Write-Host "Cleaning up old report directory..."
if (Test-Path $GoCoverDir) {
    Remove-Item -Path $GoCoverDir -Recurse -Force
}
New-Item -Path $GoCoverDir -ItemType Directory | Out-Null

# Run Tests and Generate Coverage Profile
Write-Host "`n--- Running Go tests for all packages ---" -ForegroundColor Green

# Temporarily change to the project directory to run Go commands
Push-Location $ProjectDir

# The './...' argument tells Go to run tests in the current directory and all subdirectories.
# The -coverprofile flag generates the raw coverage data.
go test -v -cover -coverprofile "$CoverageProfile" ./...

# Check if the tests failed
if ($LASTEXITCODE -ne 0) {
    Write-Error "Go tests failed. Please check the output above for errors."
    Pop-Location
    Read-Host "Press Enter to exit"; exit 1
}

# Return to the original directory
Pop-Location

# Final Message
Write-Host "`nGo coverage profile generated successfully!" -ForegroundColor Green
Write-Host "Raw coverage data saved to:"
Write-Host $CoverageProfile