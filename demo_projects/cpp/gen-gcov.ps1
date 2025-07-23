# gen-gcov.ps1 - Builds the project, runs tests, and generates all gcov text report formats.

Write-Host "--- Build and Generate Gcov Reports ---" -ForegroundColor Green

# Path Setup and Cleaning
$ScriptDir = Split-Path -Parent $MyInvocation.MyCommand.Path
$ProjectDir = Join-Path $ScriptDir "project"
$BuildDir = Join-Path $ProjectDir "build"
$GcovReportDir = Join-Path $ScriptDir "report\gcov"

# Clean up previous build and report files for a fresh run
Write-Host "Cleaning up old build and report directories..."
if (Test-Path $BuildDir) { Remove-Item -Path $BuildDir -Recurse -Force }
if (Test-Path $GcovReportDir) { Remove-Item -Path $GcovReportDir -Recurse -Force }
New-Item -Path $BuildDir -ItemType Directory | Out-Null
New-Item -Path $GcovReportDir -ItemType Directory | Out-Null

# Build Project with Coverage Flags
Write-Host "`n--- Building the C++ Project with Coverage Flags ---" -ForegroundColor Green

Write-Host "Running CMake to configure the project..."
cmake -S $ProjectDir -B $BuildDir -G "MinGW Makefiles"
if ($LASTEXITCODE -ne 0) { Write-Error "CMake configuration failed."; Read-Host "Press Enter to exit"; exit 1 }

Write-Host "Running build..."
cmake --build $BuildDir
if ($LASTEXITCODE -ne 0) { Write-Error "Build failed."; Read-Host "Press Enter to exit"; exit 1 }

# Generate Coverage Data and Reports
Push-Location $BuildDir

Write-Host "`n--- Generating Coverage Data ---" -ForegroundColor Green
Write-Host "Running tests to create .gcda files..."
& ".\run_tests.exe"
if ($LASTEXITCODE -ne 0) { Write-Error "Test execution failed."; Pop-Location; Read-Host "Press Enter to exit"; exit 1 }

# Find all the instrumented object files for our library
$ObjectFiles = Get-ChildItem -Path ".\CMakeFiles\app_lib.dir" -Recurse -Filter "*.cpp.obj"

if ($ObjectFiles.Count -eq 0) {
    Write-Error "No .cpp.obj files found for 'app_lib'. Check build output."
    Pop-Location; Read-Host "Press Enter to exit"; exit 1
}

Write-Host "`n--- Generating All Gcov Report Formats ---" -ForegroundColor Green
Write-Host "Found $($ObjectFiles.Count) source files to process."

# Define the report types and their corresponding gcov flags
$reportTypes = @{
    "basic"                  = "";
    "branch-probabilities"   = "-b";
    "branch-counts"          = "-b -c";
    "unconditional-branches" = "-b -c -u";
}

foreach ($entry in $reportTypes.GetEnumerator()) {
    $dirName = $entry.Name
    $flags = $entry.Value
    $destinationDir = Join-Path $GcovReportDir $dirName
    New-Item -Path $destinationDir -ItemType Directory | Out-Null
    Write-Host "Generating report for: '$dirName'"

    # Process each object file found earlier
    foreach ($objFile in $ObjectFiles) {
        $gcovArgs = $flags.Split(' ', [System.StringSplitOptions]::RemoveEmptyEntries) + $objFile.FullName
        gcov.exe $gcovArgs | Out-Null

        # Move the resulting .gcov file to the correct report directory
        $gcovFile = (Split-Path -Leaf $objFile.Name) -replace '\.obj$', '.gcov'
        Move-Item -Path ".\$gcovFile" -Destination $destinationDir -Force
    }
}

Pop-Location

Write-Host "`nAll C++ gcov reports generated successfully in '$GcovReportDir'." -ForegroundColor Green