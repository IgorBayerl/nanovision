param (
    [string]$InstallDir = "$env:LOCALAPPDATA\nanovision",
    [switch]$Force
)

$Repo = "IgorBayerl/nanovision"
$ApiUrl = "https://api.github.com/repos/$Repo/releases/latest"

Write-Host "Checking for nanovision updates..." -ForegroundColor Cyan

try {
    $LatestRelease = Invoke-RestMethod -Uri $ApiUrl -ErrorAction Stop
}
catch {
    Write-Error "Failed to connect to GitHub to check for updates."
    exit 1
}

$LatestVersion = $LatestRelease.tag_name
$CurrentVersion = "none"

# Check if installed
if (Get-Command nanovision -ErrorAction SilentlyContinue) {
    $VerStr = nanovision --version
    # Extract v0.X.X from output string
    if ($VerStr -match "(v[\d\.]+)") {
        $CurrentVersion = $matches[1]
    }
}

if ($CurrentVersion -eq $LatestVersion -and -not $Force) {
    Write-Host "You are already on the latest version ($LatestVersion)." -ForegroundColor Green
    exit 0
}

if ($CurrentVersion -ne "none") {
    Write-Host "Update available: $CurrentVersion -> $LatestVersion" -ForegroundColor Yellow
    $Choice = Read-Host "Do you want to read the Changelog? (y/N)"
    if ($Choice -eq 'y') {
        Write-Host "`n--- Changelog ---" -ForegroundColor Gray
        Write-Host $LatestRelease.body
        Write-Host "-----------------`n" -ForegroundColor Gray
    }
    $Choice = Read-Host "Update now? (Y/n)"
    if ($Choice -eq 'n') { exit 0 }
}

# Download & Install
$ZipName = "nanovision_${LatestVersion}_windows_amd64.zip"
$Asset = $LatestRelease.assets | Where-Object { $_.name -eq $ZipName }

if (-not $Asset) {
    Write-Error "Release asset $ZipName not found."
    exit 1
}

$ZipPath = "$env:TEMP\$ZipName"
Write-Host "Downloading..." -ForegroundColor Cyan
Invoke-WebRequest -Uri $Asset.browser_download_url -OutFile $ZipPath

if (-not (Test-Path $InstallDir)) { New-Item -ItemType Directory -Path $InstallDir | Out-Null }

# Unzip
Expand-Archive -Path $ZipPath -DestinationPath $InstallDir -Force

# === BYPASS UNKNOWN PUBLISHER ===
# Unblock the exe so SmartScreen doesn't block it
Unblock-File -Path "$InstallDir\nanovision.exe"
# ================================

Remove-Item $ZipPath

# Add to Path if missing
$UserPath = [Environment]::GetEnvironmentVariable("Path", "User")
if ($UserPath -notlike "*$InstallDir*") {
    [Environment]::SetEnvironmentVariable("Path", "$UserPath;$InstallDir", "User")
    $env:Path += ";$InstallDir"
    Write-Host "Added to PATH. You may need to restart your terminal." -ForegroundColor Yellow
}

Write-Host "nanovision $LatestVersion installed!" -ForegroundColor Green