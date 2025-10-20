# setup.ps1 - Prepares the development environment for the C++ examples
# MUST BE RUN AS ADMINISTRATOR from a 64-bit PowerShell terminal

# Check for Administrator privileges
if (-Not ([Security.Principal.WindowsPrincipal][Security.Principal.WindowsIdentity]::GetCurrent()).IsInRole([Security.Principal.WindowsBuiltInRole]::Administrator)) {
    Write-Warning "This script needs to be run as Administrator to install system-wide tools."
    Read-Host "Press Enter to exit"
    exit 1
}

# Check for 64-bit PowerShell
if ($env:PROCESSOR_ARCHITECTURE -eq 'x86') {
    Write-Warning "This setup script should be run from a 64-bit PowerShell terminal."
    Read-Host "Press Enter to exit"
    exit 1
}

Write-Host "--- Starting Environment Preparation ---" -ForegroundColor Green

# Install Chocolatey if not present
if (-not (Get-Command choco -ErrorAction SilentlyContinue)) {
    Write-Host "Installing Chocolatey..." -ForegroundColor Yellow
    Set-ExecutionPolicy Bypass -Scope Process -Force; [System.Net.ServicePointManager]::SecurityProtocol = [System.Net.ServicePointManager]::SecurityProtocol -bor 3072; iex ((New-Object System.Net.WebClient).DownloadString('https://community.chocolatey.org/install.ps1'))
} else {
    Write-Host "Chocolatey is already installed."
}

# Install required packages
$packages = @("mingw", "cmake", "python")

Write-Host "`nChecking for required packages..." -ForegroundColor Green
foreach ($pkg in $packages) {
    if (-not (choco list --local-only --exact $pkg | Select-String $pkg)) {
        Write-Host "Installing $pkg..." -ForegroundColor Yellow
        choco install $pkg -y
    } else {
        Write-Host "$pkg is already installed."
    }
}

Write-Host "`nInstalling gcovr (Python)..."
pip install --upgrade gcovr

Write-Host "`nEnvironment is ready!" -ForegroundColor Green
Write-Host "You may need to open a new PowerShell window for PATH changes to take effect."
Write-Host "Next, run one of the report generation scripts, like 'gen-gcov.ps1' or 'gen-lcov.ps1'."