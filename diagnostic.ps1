# Script de diagnostico rapido para SCP Go Identity Service

Write-Host "=== DIAGNOSTICO SCP GO IDENTITY SERVICE ===" -ForegroundColor Cyan
Write-Host ""

# 1. Verificar archivos necesarios
Write-Host "1. Verificando archivos..." -ForegroundColor Yellow
$files = @("scp-go-identity-service.exe", "config.json", "auth-test.html")
foreach ($file in $files) {
    if (Test-Path $file) {
        Write-Host "   ✓ $file" -ForegroundColor Green
    } else {
        Write-Host "   ✗ $file (FALTANTE)" -ForegroundColor Red
    }
}

# 2. Verificar configuracion
Write-Host ""
Write-Host "2. Verificando configuracion..." -ForegroundColor Yellow
if (Test-Path "config.json") {
    try {
        $config = Get-Content "config.json" | ConvertFrom-Json
        Write-Host "   ✓ config.json es JSON valido" -ForegroundColor Green
        Write-Host "   - Database: $($config.database.connectionString.Split(';')[0])" -ForegroundColor Gray
        Write-Host "   - Server: $($config.server.host):$($config.server.port)" -ForegroundColor Gray
    } catch {
        Write-Host "   ✗ config.json tiene formato invalido" -ForegroundColor Red
    }
} else {
    Write-Host "   ✗ config.json no encontrado" -ForegroundColor Red
}

# 3. Intentar iniciar servicio brevemente para probar
Write-Host ""
Write-Host "3. Probando inicio del servicio..." -ForegroundColor Yellow

$job = Start-Job -ScriptBlock {
    cd $using:PWD
    .\scp-go-identity-service.exe 2>&1
}

Start-Sleep -Seconds 5

# 4. Verificar si el servicio responde
Write-Host ""
Write-Host "4. Verificando conexion..." -ForegroundColor Yellow
try {
    $response = Invoke-WebRequest -Uri "http://localhost:8081/health" -TimeoutSec 5
    Write-Host "   ✓ Servicio responde correctamente" -ForegroundColor Green
    
    # Probar endpoint de debug si esta disponible
    try {
        $debugResponse = Invoke-WebRequest -Uri "http://localhost:8081/debug" -TimeoutSec 5
        $debugData = $debugResponse.Content | ConvertFrom-Json
        Write-Host "   ✓ Debug endpoint disponible" -ForegroundColor Green
        Write-Host "   - Database status: $($debugData.database)" -ForegroundColor Gray
    } catch {
        Write-Host "   - Debug endpoint no disponible (normal en version anterior)" -ForegroundColor Gray
    }
    
} catch {
    Write-Host "   ✗ No se puede conectar al servicio" -ForegroundColor Red
    Write-Host "   Error: $($_.Exception.Message)" -ForegroundColor Gray
    
    # Obtener logs del job
    $jobOutput = Receive-Job $job
    if ($jobOutput) {
        Write-Host ""
        Write-Host "   Salida del servicio:" -ForegroundColor Yellow
        $jobOutput | ForEach-Object { Write-Host "   > $_" -ForegroundColor Gray }
    }
}

# Limpiar
Stop-Job $job -ErrorAction SilentlyContinue
Remove-Job $job -ErrorAction SilentlyContinue

Write-Host ""
Write-Host "=== FIN DIAGNOSTICO ===" -ForegroundColor Cyan