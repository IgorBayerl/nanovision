# setup.ps1 - Prepares the environment for the C# project.
# It verifies the .NET SDK installation and installs it via Chocolatey if missing.
# MUST BE RUN AS ADMINISTRATOR from a PowerShell terminal.

# Check for Administrator privileges, which are required for installation.
if (-Not ([Security.Principal.WindowsPrincipal][Security.Principal.WindowsIdentity]::GetCurrent()).IsInRole([Security.Principal.WindowsBuiltInRole]::Administrator)) {
    Write-Error "This script must be run as Administrator to install the .NET SDK."
    Read-Host "Press Enter to exit"; exit 1
}

Write-Host "--- C# Environment Preparation ---" -ForegroundColor Green

# Check for .NET SDK
Write-Host "`nChecking for the .NET SDK..." -ForegroundColor Green

if (Get-Command dotnet -ErrorAction SilentlyContinue) {
    Write-Host ".NET SDK is already installed. Displaying details:" -ForegroundColor Cyan
    dotnet --info
    Write-Host "`n✅ The environment is ready." -ForegroundColor Green
    Write-Host "You can now run a report generation script like '.\gen-csharp-cobertura.ps1'."
} else {
    # --- Installation Logic ---
    Write-Warning ".NET SDK was NOT found. Attempting to install via Chocolatey..."

    # 1. Install Chocolatey if it's not present
    if (-not (Get-Command choco -ErrorAction SilentlyContinue)) {
        Write-Host "`nChocolatey package manager not found. Installing Chocolatey first..." -ForegroundColor Yellow
        Set-ExecutionPolicy Bypass -Scope Process -Force; [System.Net.ServicePointManager]::SecurityProtocol = [System.Net.ServicePointManager]::SecurityProtocol -bor 3072; iex ((New-Object System.Net.WebClient).DownloadString('https://community.chocolatey.org/install.ps1'))
    } else {
        Write-Host "`nChocolatey is already installed."
    }

    # 2. Install the .NET 8 SDK using Chocolatey
    Write-Host "`nInstalling .NET 9 SDK (dotnet-9.0-sdk)..." -ForegroundColor Yellow
    choco install dotnet-9.0-sdk -y

    if ($LASTEXITCODE -ne 0) {
        Write-Error "`n❌ .NET SDK installation failed. Please check the Chocolatey output above for errors."
    } else {
        Write-Host "`n✅ .NET SDK has been installed." -ForegroundColor Green
        Write-Warning "You MUST open a new PowerShell terminal for the system PATH changes to take effect."
        Write-Host "After opening a new terminal, you can re-run this script to verify the installation or run a report script."
    }
}

# Finalization
Read-Host "Press Enter to exit"