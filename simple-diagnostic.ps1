# Diagnostico simple para SCP Go Identity Service

Write-Host "=== DIAGNOSTICO ===" -ForegroundColor Cyan

# Verificar archivos
Write-Host "Archivos:" -ForegroundColor Yellow
if (Test-Path "scp-go-identity-service.exe") { Write-Host "✓ Ejecutable encontrado" -ForegroundColor Green } else { Write-Host "✗ Ejecutable faltante" -ForegroundColor Red }
if (Test-Path "config.json") { Write-Host "✓ Config encontrado" -ForegroundColor Green } else { Write-Host "✗ Config faltante" -ForegroundColor Red }

# Iniciar servicio para prueba
Write-Host ""
Write-Host "Iniciando servicio de prueba..." -ForegroundColor Yellow
$job = Start-Job { Set-Location $using:PWD; .\scp-go-identity-service.exe }
Start-Sleep 3

# Probar conexion
try {
    Invoke-WebRequest -Uri "http://localhost:8081/health" -TimeoutSec 3 | Out-Null
    Write-Host "✓ Servicio funciona correctamente" -ForegroundColor Green
} catch {
    Write-Host "✗ Error en el servicio: $($_.Exception.Message)" -ForegroundColor Red
    
    # Mostrar output del servicio
    $output = Receive-Job $job
    Write-Host "Output del servicio:" -ForegroundColor Yellow
    $output | ForEach-Object { Write-Host "$_" -ForegroundColor Gray }
}

# Limpiar
Stop-Job $job 2>$null
Remove-Job $job 2>$null

Write-Host "================" -ForegroundColor Cyan