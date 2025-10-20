# gen-csharp-opencover.ps1 - Runs C# tests and generates an OpenCover XML coverage report.

Write-Host "--- Generate C# OpenCover Report ---" -ForegroundColor Green

# Path Setup
$ScriptDir = Split-Path -Parent $MyInvocation.MyCommand.Path
$TestProjectFile = Join-Path $ScriptDir "project\UnitTests\UnitTests.csproj"
$ReportBaseDir   = Join-Path $ScriptDir "report"
$OpenCoverDir    = Join-Path $ReportBaseDir "opencover"
$XmlReportFile   = Join-Path $OpenCoverDir "coverage.opencover.xml"

if (-not (Test-Path $TestProjectFile)) {
    Write-Error "Test project not found at: $TestProjectFile"; Read-Host; exit 1
}

# Cleaning
Write-Host "Cleaning up old report directory..."
if (Test-Path $OpenCoverDir) { Remove-Item -Path $OpenCoverDir -Recurse -Force }
New-Item -Path $OpenCoverDir -ItemType Directory | Out-Null

# Run Tests and Generate Report
Write-Host "`n--- Running 'dotnet test' with OpenCover coverage ---" -ForegroundColor Green

dotnet test $TestProjectFile --configuration Release /p:CollectCoverage=true /p:CoverletOutputFormat=opencover "/p:CoverletOutput=$XmlReportFile"

if ($LASTEXITCODE -ne 0) {
    Write-Error "'dotnet test' command failed."; Read-Host; exit 1
}
if (-not (Test-Path $XmlReportFile)) {
    Write-Error "Coverage report was not generated."; Read-Host; exit 1
}

# Final Message
Write-Host "`nC# OpenCover report generated successfully!" -ForegroundColor Green
Write-Host "Report saved to: $XmlReportFile"