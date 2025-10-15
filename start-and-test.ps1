# Script simple para iniciar el servicio y abrir la pagina de estado
# Para desarrollo y pruebas rapidas

Write-Host "==================================" -ForegroundColor Cyan
Write-Host "  SCP Go Identity Service - Prueba  " -ForegroundColor Cyan
Write-Host "==================================" -ForegroundColor Cyan
Write-Host ""

# Verificar que el ejecutable existe
if (-not (Test-Path "scp-go-identity-service.exe")) {
    Write-Host "ERROR: No se encuentra scp-go-identity-service.exe" -ForegroundColor Red
    Write-Host "Ejecuta primero: go build -o scp-go-identity-service.exe main.go" -ForegroundColor Yellow
    exit 1
}

# Verificar que existe el archivo de configuracion
if (-not (Test-Path "config.json")) {
    Write-Host "ERROR: No se encuentra config.json" -ForegroundColor Red
    Write-Host "Copia config-example.json a config.json y configuralo" -ForegroundColor Yellow
    exit 1
}

Write-Host "1. Iniciando servicio..." -ForegroundColor Green
Write-Host "   Ejecutable: scp-go-identity-service.exe" -ForegroundColor Gray
Write-Host "   Config: config.json" -ForegroundColor Gray

# Iniciar el servicio en background
$job = Start-Job -ScriptBlock {
    Set-Location $using:PWD
    .\scp-go-identity-service.exe
}

# Esperar un momento para que inicie
Write-Host "2. Esperando a que el servicio inicie..." -ForegroundColor Yellow
Start-Sleep -Seconds 3

# Verificar si esta funcionando
try {
    $response = Invoke-WebRequest -Uri "http://localhost:8081/health" -Method GET -TimeoutSec 5
    if ($response.StatusCode -eq 200) {
        Write-Host "3. Servicio iniciado correctamente!" -ForegroundColor Green
        Write-Host "4. Abriendo pagina de estado..." -ForegroundColor Cyan
        
        # Abrir pagina de estado
        Start-Process "http://localhost:8081/status"
        
        Write-Host ""
        Write-Host "================================" -ForegroundColor Green
        Write-Host "  SERVICIO FUNCIONANDO" -ForegroundColor Green
        Write-Host "================================" -ForegroundColor Green
        Write-Host "Pagina de estado: http://localhost:8081/status" -ForegroundColor White
        Write-Host "API Health: http://localhost:8081/health" -ForegroundColor White
        Write-Host "Autenticacion: http://localhost:8081/authenticate" -ForegroundColor White
        Write-Host ""
        Write-Host "Presiona Ctrl+C para detener el servicio" -ForegroundColor Yellow
        
        # Mantener el servicio corriendo
        try {
            Wait-Job $job
        } finally {
            Write-Host ""
            Write-Host "Deteniendo servicio..." -ForegroundColor Yellow
            Stop-Job $job -PassThru | Remove-Job
        }
    } else {
        throw "Servicio respondio con codigo $($response.StatusCode)"
    }
} catch {
    Write-Host "ERROR: No se pudo conectar al servicio" -ForegroundColor Red
    Write-Host "Detalle: $($_.Exception.Message)" -ForegroundColor Yellow
    
    # Limpiar job
    Stop-Job $job -PassThru | Remove-Job
    exit 1
}