# start-dev.ps1 — Arranque de desarrollo: backend Go (puerto 8080) + frontend Vite (puerto 3000)
# Uso: .\start-dev.ps1
# Cierra ambos procesos con Ctrl+C.

$ErrorActionPreference = "Stop"
$root = $PSScriptRoot

# Verificar config.yaml
if (-not (Test-Path "$root\server\config.yaml")) {
    Write-Host "ERROR: server\config.yaml no existe." -ForegroundColor Red
    Write-Host "Copia la plantilla y rellena jwt_secret:" -ForegroundColor Yellow
    Write-Host "  cp server\config.yaml.template server\config.yaml" -ForegroundColor Cyan
    exit 1
}

Write-Host "Arrancando backend (puerto 8080)..." -ForegroundColor Green
$backend = Start-Process -PassThru -NoNewWindow powershell -ArgumentList "-Command", "Set-Location '$root\server'; `$env:ENV='development'; go run main.go"

Write-Host "Arrancando frontend (puerto 3000)..." -ForegroundColor Green
$frontend = Start-Process -PassThru -NoNewWindow powershell -ArgumentList "-Command", "Set-Location '$root\client'; pnpm dev"

Write-Host ""
Write-Host "App corriendo en http://localhost:3000  (API en http://localhost:8080)" -ForegroundColor Cyan
Write-Host "Pulsa Ctrl+C para detener ambos procesos." -ForegroundColor Yellow

try {
    Wait-Process -Id $backend.Id
} finally {
    if (-not $backend.HasExited)  { Stop-Process -Id $backend.Id  -Force -ErrorAction SilentlyContinue }
    if (-not $frontend.HasExited) { Stop-Process -Id $frontend.Id -Force -ErrorAction SilentlyContinue }
    Write-Host "Procesos detenidos." -ForegroundColor Gray
}
