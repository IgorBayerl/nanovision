<#
    Builds the project, runs the tests, and creates a Cobertura‑style XML
    coverage report via gcovr.  Run this after setup.ps1 has prepared
    MinGW/CMake/Python.
#>

Write-Host "--- Build & Cobertura XML Report ---" -ForegroundColor Green

# Paths
$ScriptDir     = Split-Path -Parent $MyInvocation.MyCommand.Path
$ProjectDir    = Join-Path $ScriptDir "project"
$BuildDir      = Join-Path $ProjectDir "build"
$ReportDir     = Join-Path $ScriptDir "report"
$CoberturaDir  = Join-Path $ReportDir "cobertura"
$XmlReport     = Join-Path $CoberturaDir "cobertura.xml"

# Fresh build folder
if (Test-Path $BuildDir) { Remove-Item $BuildDir -Recurse -Force }
New-Item $BuildDir -ItemType Directory | Out-Null

# Ensure report/cobertura exists
if (Test-Path $CoberturaDir) { Remove-Item $CoberturaDir -Recurse -Force }
New-Item $CoberturaDir -ItemType Directory | Out-Null

# Configure & build
cmake -S $ProjectDir -B $BuildDir -G "MinGW Makefiles"
cmake --build $BuildDir

# Run tests (produces .gcda)
Push-Location $BuildDir
& ".\run_tests.exe"
Pop-Location

# Generate Cobertura XML with gcovr
Write-Host "Generating Cobertura XML…"
gcovr -r $ProjectDir --xml-pretty -o $XmlReport

Write-Host "Report saved to $XmlReport" -ForegroundColor Green
