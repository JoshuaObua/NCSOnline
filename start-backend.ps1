#!/usr/bin/env pwsh
# NCS Backend Setup Script
# Runs after PostgreSQL 17 is installed via winget
# Sets up database, runs migrations, and starts the backend server

param(
    [string]$PgPassword = "NcsIntranetDb2026!Secure",
    [string]$DbName     = "ncsintranet",
    [string]$DbUser     = "ncsintranet_user",
    [string]$DbPass     = "NcsIntranetDb2026!Secure"
)

$ErrorActionPreference = "Stop"
Set-ExecutionPolicy -Scope Process -ExecutionPolicy Bypass

# Refresh PATH to pick up newly installed tools
$env:PATH = [System.Environment]::GetEnvironmentVariable('PATH','Machine') + ';' + [System.Environment]::GetEnvironmentVariable('PATH','User')

Write-Host "=== NCS Intranet Backend Setup ===" -ForegroundColor Cyan

# ── 1. Find psql ──────────────────────────────────────────────────────────────
$pgBin = "C:\Program Files\PostgreSQL\17\bin"
if (-not (Test-Path "$pgBin\psql.exe")) {
    # Try other versions
    $found = Get-ChildItem "C:\Program Files\PostgreSQL" -Directory -ErrorAction SilentlyContinue |
             Sort-Object Name -Descending |
             Select-Object -First 1
    if ($found) { $pgBin = "$($found.FullName)\bin" }
}

if (-not (Test-Path "$pgBin\psql.exe")) {
    Write-Error "PostgreSQL not found. Ensure it is installed first."
    exit 1
}

$env:PATH = "$pgBin;$env:PATH"
$env:PGPASSWORD = $PgPassword
Write-Host "[OK] Found psql at $pgBin" -ForegroundColor Green

# ── 2. Wait for PostgreSQL service ────────────────────────────────────────────
Write-Host "[..] Waiting for PostgreSQL service to start..." -ForegroundColor Yellow
$svc = Get-Service -Name "postgresql*" -ErrorAction SilentlyContinue | Select-Object -First 1
if ($svc) {
    if ($svc.Status -ne 'Running') {
        Start-Service $svc.Name
        Start-Sleep 3
    }
    Write-Host "[OK] PostgreSQL service: $($svc.Name) is $($svc.Status)" -ForegroundColor Green
} else {
    Write-Host "[WARN] No PostgreSQL service found - trying to connect anyway..." -ForegroundColor Yellow
}

# ── 3. Create database and user ───────────────────────────────────────────────
Write-Host "[..] Creating database user and database..." -ForegroundColor Yellow

$setupSql = @"
DO `$`$
BEGIN
   IF NOT EXISTS (SELECT FROM pg_catalog.pg_roles WHERE rolname = '$DbUser') THEN
      CREATE ROLE $DbUser LOGIN PASSWORD '$DbPass';
   END IF;
END
`$`$;

CREATE DATABASE $DbName OWNER $DbUser ENCODING 'UTF8' LC_COLLATE 'en_US.UTF-8' LC_CTYPE 'en_US.UTF-8' TEMPLATE template0;
GRANT ALL PRIVILEGES ON DATABASE $DbName TO $DbUser;
"@

# Try to run as postgres superuser
try {
    $setupSql | & "$pgBin\psql.exe" -U postgres -h localhost -p 5432 -c "SELECT 1" > $null 2>&1
    if ($LASTEXITCODE -eq 0) {
        $setupSql | & "$pgBin\psql.exe" -U postgres -h localhost -p 5432 2>&1
        Write-Host "[OK] Database and user created" -ForegroundColor Green
    }
} catch {
    Write-Host "[WARN] Could not create DB as postgres. Trying direct connection as $DbUser..." -ForegroundColor Yellow
}

# ── 4. Run migrations ─────────────────────────────────────────────────────────
Write-Host "[..] Running database migrations..." -ForegroundColor Yellow
$env:PGPASSWORD = $DbPass
$migrationsDir = "d:\NCSOnline\backend\migrations"
$migrationFiles = Get-ChildItem $migrationsDir -Filter "*.sql" | Sort-Object Name

$successCount = 0
$failCount = 0
foreach ($file in $migrationFiles) {
    $result = & "$pgBin\psql.exe" -U $DbUser -h localhost -p 5432 -d $DbName -f $file.FullName 2>&1
    if ($LASTEXITCODE -eq 0) {
        $successCount++
    } else {
        Write-Host "  [WARN] $($file.Name): $result" -ForegroundColor Yellow
        $failCount++
    }
}
Write-Host "[OK] Migrations: $successCount succeeded, $failCount warnings" -ForegroundColor Green

# ── 5. Start the backend ──────────────────────────────────────────────────────
Write-Host "[..] Starting NCS backend on port 9081..." -ForegroundColor Yellow
$backendExe = "d:\NCSOnline\backend\ncs-backend.exe"

if (-not (Test-Path $backendExe)) {
    Write-Host "[..] Building backend..." -ForegroundColor Yellow
    $goPath = (Get-Command go -ErrorAction SilentlyContinue)?.Source
    if (-not $goPath) { Write-Error "Go not found in PATH"; exit 1 }
    Push-Location d:\NCSOnline\backend
    & go build -o ncs-backend.exe ./cmd/server
    Pop-Location
}

# Load .env vars into process
Get-Content "d:\NCSOnline\.env" | ForEach-Object {
    if ($_ -match '^\s*([^#][^=]+)=(.*)$') {
        [System.Environment]::SetEnvironmentVariable($matches[1].Trim(), $matches[2].Trim(), 'Process')
    }
}

Write-Host "[OK] Starting backend..." -ForegroundColor Green
Write-Host "     URL: http://localhost:9081" -ForegroundColor Cyan
Write-Host "     Press Ctrl+C to stop" -ForegroundColor Gray
& $backendExe
