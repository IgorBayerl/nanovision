<#
    Builds the project, runs the tests, and generates an interactive HTML
    coverage site (with per‑file & per‑line detail pages) using gcovr.
#>

Write-Host "--- Build & HTML Coverage Report ---" -ForegroundColor Green

# Paths
$ScriptDir = Split-Path -Parent $MyInvocation.MyCommand.Path
$ProjectDir = Join-Path $ScriptDir "project"
$BuildDir = Join-Path $ProjectDir "build"
$HtmlDir  = Join-Path $ScriptDir "report\html"

# Prepare folders
if (Test-Path $BuildDir) { Remove-Item $BuildDir -Recurse -Force }
New-Item $BuildDir -ItemType Directory | Out-Null
if (Test-Path $HtmlDir) { Remove-Item $HtmlDir -Recurse -Force }
New-Item $HtmlDir -ItemType Directory | Out-Null

# Configure & build
cmake -S $ProjectDir -B $BuildDir -G "MinGW Makefiles"
cmake --build $BuildDir

# Run tests (produces .gcda)
Push-Location $BuildDir
& ".\run_tests.exe"
Pop-Location

# Generate HTML dashboard
Write-Host "Generating HTML coverage site…"
$IndexFile = Join-Path $HtmlDir "index.html"
gcovr -r $ProjectDir --html --html-details -o $IndexFile

Write-Host "Open $IndexFile in your browser to view the report." -ForegroundColor Green
