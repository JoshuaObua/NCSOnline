# ==============================================================================
# NCS Intranet & Online Platform - Podman Environment Automation Script
# ==============================================================================

Write-Host "==========================================================" -ForegroundColor Cyan
Write-Host "   National Council of Sports (NCS) - Podman Setup Tool   " -ForegroundColor Yellow
Write-Host "==========================================================" -ForegroundColor Cyan

# 1. Check Podman Installation
$podmanCmd = Get-Command podman -ErrorAction SilentlyContinue
if (-not $podmanCmd) {
    Write-Host "[!] Podman CLI not found on your system PATH." -ForegroundColor Yellow
    Write-Host "    You can install Podman for Windows using Winget:" -ForegroundColor Cyan
    Write-Host "    > winget install RedHat.Podman" -ForegroundColor White
    Write-Host "    Or download the installer from: https://podman.io" -ForegroundColor White
    Write-Host "    After installing, initialize and start the machine:" -ForegroundColor Cyan
    Write-Host "    > podman machine init" -ForegroundColor White
    Write-Host "    > podman machine start" -ForegroundColor White
    Write-Host ""
} else {
    Write-Host "[+] Podman CLI detected: $(podman --version)" -ForegroundColor Green
    
    # Check if podman machine is running (Windows WSL backend)
    $machineState = podman machine list --format "{{.Running}}" 2>$null
    if ($machineState -and ($machineState -notmatch "true")) {
        Write-Host "[*] Starting Podman virtual machine..." -ForegroundColor Cyan
        podman machine start
    }
}

# 2. Check .env configuration
$envFile = Join-Path $PSScriptRoot ".env"
$envExample = Join-Path $PSScriptRoot ".env.example"

if (-not (Test-Path $envFile)) {
    if (Test-Path $envExample) {
        Write-Host "[*] Creating .env from .env.example..." -ForegroundColor Cyan
        Copy-Item -Path $envExample -Destination $envFile
        Write-Host "[+] .env configuration file generated." -ForegroundColor Green
    }
} else {
    Write-Host "[+] .env configuration file found." -ForegroundColor Green
}

# 3. Choose compose provider (podman compose OR podman-compose)
$composeCmd = ""
if (Get-Command podman-compose -ErrorAction SilentlyContinue) {
    $composeCmd = "podman-compose -f podman-compose.yml"
} elseif (Get-Command podman -ErrorAction SilentlyContinue) {
    $composeCmd = "podman compose -f podman-compose.yml"
}

if ($composeCmd -ne "") {
    Write-Host "[*] Executing container build and startup using: $composeCmd" -ForegroundColor Cyan
    Invoke-Expression "$composeCmd up -d --build"
    
    Write-Host ""
    Write-Host "==========================================================" -ForegroundColor Green
    Write-Host "   NCS Intranet Stack successfully launched via Podman!   " -ForegroundColor Green
    Write-Host "==========================================================" -ForegroundColor Green
    Write-Host "Portal URL:       http://localhost:9081" -ForegroundColor Yellow
    Write-Host "Health Check:     http://localhost:9081/healthz" -ForegroundColor Yellow
    Write-Host "Credentials CSV:  Docs/Credentials.csv" -ForegroundColor Yellow
    Write-Host "==========================================================" -ForegroundColor Green
} else {
    Write-Host "[!] Podman Compose command not yet executable in this session." -ForegroundColor Yellow
    Write-Host "    Please ensure Podman Desktop / CLI is running, then run:" -ForegroundColor Cyan
    Write-Host "    podman compose -f podman-compose.yml up -d --build" -ForegroundColor White
}
