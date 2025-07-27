# gen-csharp-lcov.ps1 - Runs C# tests and generates an lcov format coverage report.

Write-Host "--- Generate C# Lcov Report ---" -ForegroundColor Green

# Path Setup
$ScriptDir = Split-Path -Parent $MyInvocation.MyCommand.Path
$TestProjectFile = Join-Path $ScriptDir "project\UnitTests\UnitTests.csproj"
$ReportBaseDir   = Join-Path $ScriptDir "report"
$LcovDir         = Join-Path $ReportBaseDir "lcov_csharp"
$LcovReportFile  = Join-Path $LcovDir "coverage.lcov"

if (-not (Test-Path $TestProjectFile)) {
    Write-Error "Test project not found at: $TestProjectFile"; Read-Host; exit 1
}

# Cleaning
Write-Host "Cleaning up old report directory..."
if (Test-Path $LcovDir) { Remove-Item -Path $LcovDir -Recurse -Force }
New-Item -Path $LcovDir -ItemType Directory | Out-Null

# Run Tests and Generate Report
Write-Host "`n--- Running 'dotnet test' with lcov coverage ---" -ForegroundColor Green

dotnet test $TestProjectFile --configuration Release /p:CollectCoverage=true /p:CoverletOutputFormat=lcov "/p:CoverletOutput=$LcovReportFile"

if ($LASTEXITCODE -ne 0) {
    Write-Error "'dotnet test' command failed."; Read-Host; exit 1
}
if (-not (Test-Path $LcovReportFile)) {
    Write-Error "Coverage report was not generated."; Read-Host; exit 1
}

# Final Message
Write-Host "`nC# lcov report generated successfully!" -ForegroundColor Green
Write-Host "Report saved to: $LcovReportFile"