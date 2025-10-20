# gen-csharp-cobertura.ps1 - Runs C# tests and generates a Cobertura XML coverage report.

Write-Host "--- Generate C# Cobertura Report ---" -ForegroundColor Green

$ScriptDir = Split-Path -Parent $MyInvocation.MyCommand.Path

$TestProjectFile = Join-Path $ScriptDir "project\UnitTests\UnitTests.csproj"
$ReportBaseDir   = Join-Path $ScriptDir "report"
$CSharpReportDir = Join-Path $ReportBaseDir "cobertura"

$XmlReportFile = Join-Path $CSharpReportDir "cobertura.xml"

if (-not (Test-Path $TestProjectFile)) {
    Write-Error "Test project not found at: $TestProjectFile"
    Read-Host "Press Enter to exit"; exit 1
}

Write-Host "Cleaning up old report directory..."
if (Test-Path $CSharpReportDir) {
    Remove-Item -Path $CSharpReportDir -Recurse -Force
}
New-Item -Path $CSharpReportDir -ItemType Directory | Out-Null

Write-Host "`n--- Running 'dotnet test' with coverage ---" -ForegroundColor Green

dotnet test $TestProjectFile --configuration Release /p:CollectCoverage=true /p:CoverletOutputFormat=cobertura "/p:CoverletOutput=$XmlReportFile"

if ($LASTEXITCODE -ne 0) {
    Write-Error "'dotnet test' command failed. Please check the output for errors."
    Read-Host "Press Enter to exit"; exit 1
}

if (-not (Test-Path $XmlReportFile)) {
    Write-Error "Coverage report was not generated. Check your project setup and Coverlet configuration."
    Read-Host "Press Enter to exit"; exit 1
}

Write-Host "`nC# Cobertura report generated successfully!" -ForegroundColor Green
Write-Host "Report saved to:"
Write-Host $XmlReportFile