# Script para probar rapidamente el SCP Go Identity Service
# Abre la pagina de estado en el navegador

param(
    [string]$Url = "http://localhost:8081/test"
)

Write-Host "Verificando SCP Go Identity Service..." -ForegroundColor Cyan
Write-Host "URL: $Url" -ForegroundColor Yellow

try {
    # Verificar si el servicio responde
    $response = Invoke-WebRequest -Uri "http://localhost:8081/health" -Method GET -TimeoutSec 5
    
    if ($response.StatusCode -eq 200) {
        Write-Host "Servicio esta funcionando correctamente!" -ForegroundColor Green
        Write-Host "Abriendo pagina de prueba de autenticacion en el navegador..." -ForegroundColor Cyan
        
        # Abrir en el navegador predeterminado
        Start-Process $Url
        
        Write-Host ""
        Write-Host "Comandos utiles:" -ForegroundColor Magenta
        Write-Host "- Verificar estado del servicio Windows: Get-Service SCPGoIdentityService" -ForegroundColor White
        Write-Host "- Ver logs: Get-EventLog -LogName Application -Source SCPGoIdentityService -Newest 5" -ForegroundColor White
        Write-Host "- Reiniciar servicio: Restart-Service SCPGoIdentityService" -ForegroundColor White
    } else {
        Write-Host "El servicio respondio con codigo: $($response.StatusCode)" -ForegroundColor Red
    }
} catch {
    Write-Host "Error al conectar con el servicio:" -ForegroundColor Red
    Write-Host "   $($_.Exception.Message)" -ForegroundColor Yellow
    Write-Host ""
    Write-Host "Posibles soluciones:" -ForegroundColor Cyan
    Write-Host "1. Verificar que el servicio esta corriendo: Get-Service SCPGoIdentityService" -ForegroundColor White
    Write-Host "2. Si no esta instalado como servicio, ejecutar: .\scp-go-identity-service.exe" -ForegroundColor White
    Write-Host "3. Verificar la configuracion en config.json" -ForegroundColor White
}

Write-Host ""