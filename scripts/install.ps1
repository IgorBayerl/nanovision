param (
    [string]$InstallDir = "$env:LOCALAPPDATA\nanovision",
    [switch]$Force
)

$Repo = "IgorBayerl/nanovision"
$ApiUrl = "https://api.github.com/repos/$Repo/releases/latest"

Write-Host "Checking for nanovision updates..." -ForegroundColor Cyan

try {
    $LatestRelease = Invoke-RestMethod -Uri $ApiUrl -Headers @{ "User-Agent" = "nanovision-installer" } -ErrorAction Stop
} catch {
    Write-Error "Failed to connect to GitHub API."
    Write-Error "Details: $($_.Exception.Message)"
    exit 1
}

if (-not $LatestRelease.tag_name) {
    Write-Error "Invalid response from GitHub API."
    exit 1
}

$LatestVersion = $LatestRelease.tag_name
$CurrentVersion = "none"

if (Get-Command nanovision -ErrorAction SilentlyContinue) {
    try {
        # Pipe to Out-String to ensure we match against a single string, not an array
        $VerStr = nanovision --version | Out-String
        if ($VerStr -match "(v[\d\.]+)") {
            $CurrentVersion = $matches[1]
        }
    } catch {
        $CurrentVersion = "unknown"
    }
}

if ($CurrentVersion -eq $LatestVersion -and -not $Force) {
    Write-Host "You are already on the latest version ($LatestVersion)." -ForegroundColor Green
    exit 0
}

if ($CurrentVersion -ne "none" -and $CurrentVersion -ne "unknown") {
    Write-Host "Update available: $CurrentVersion -> $LatestVersion" -ForegroundColor Yellow
    $Choice = Read-Host "Update now? (Y/n)"
    if ($Choice -eq 'n') { exit 0 }
}

$ZipName = "nanovision_${LatestVersion}_windows_amd64.zip"
$Asset = $LatestRelease.assets | Where-Object { $_.name -eq $ZipName }

if (-not $Asset) {
    Write-Error "Release asset '$ZipName' not found."
    exit 1
}

$ZipPath = "$env:TEMP\$ZipName"
Write-Host "Downloading $LatestVersion..." -ForegroundColor Cyan
try {
    Invoke-WebRequest -Uri $Asset.browser_download_url -OutFile $ZipPath -ErrorAction Stop
} catch {
    Write-Error "Download failed: $($_.Exception.Message)"
    exit 1
}

if (-not (Test-Path $InstallDir)) { 
    New-Item -ItemType Directory -Path $InstallDir -Force | Out-Null 
}

Expand-Archive -Path $ZipPath -DestinationPath $InstallDir -Force

if (Test-Path "$InstallDir\nanovision.exe") {
    Unblock-File -Path "$InstallDir\nanovision.exe"
} else {
    Write-Error "Extraction failed: nanovision.exe not found."
    exit 1
}

Remove-Item $ZipPath

$UserPath = [Environment]::GetEnvironmentVariable("Path", "User")
if ($UserPath -notlike "*$InstallDir*") {
    [Environment]::SetEnvironmentVariable("Path", "$UserPath;$InstallDir", "User")
    $env:Path += ";$InstallDir"
    Write-Host "Added to PATH. You may need to restart your terminal." -ForegroundColor Yellow
}

Write-Host "nanovision $LatestVersion installed successfully." -ForegroundColor Green