# Script para instalar SCP Go Identity Service como un servicio de Windows.
# DEBE ser ejecutado como Administrador.

param(
    [string]$ServiceName = "SCPGoIdentityService",
    [string]$DisplayName = "SCP Go Identity Service",
    [string]$Description = "Servicio de autenticacion Go para .NET Identity."
)

# --- 1. Verificar Permisos de Administrador ---
$currentUser = New-Object Security.Principal.WindowsPrincipal([Security.Principal.WindowsIdentity]::GetCurrent())
if (-not $currentUser.IsInRole([Security.Principal.WindowsBuiltInRole]::Administrator)) {
    Write-Error "Este script debe ser ejecutado como Administrador."
    Write-Host "Por favor, abre una nueva terminal de PowerShell como Administrador y vuelve a ejecutarlo."
    pause
    exit 1
}

# --- 2. Definir Rutas Absolutas ---
# $PSScriptRoot es la carpeta donde se encuentra este script.
$scriptPath = $PSScriptRoot
$executablePath = Join-Path $scriptPath "scp-go-identity-service.exe"
$configPath = Join-Path $scriptPath "config.json"

# Verificar que los archivos necesarios existen en la misma carpeta que el script.
if (-not (Test-Path $executablePath)) {
    Write-Error "No se encuentra el ejecutable 'scp-go-identity-service.exe' en la carpeta '$scriptPath'."
    pause
    exit 1
}
if (-not (Test-Path $configPath)) {
    Write-Error "No se encuentra el archivo de configuracion 'config.json' en la carpeta '$scriptPath'."
    pause
    exit 1
}

# --- 3. Crear el Comando del Servicio con Rutas Absolutas ---
# Esto asegura que el servicio siempre encuentre sus archivos, sin importar el directorio de trabajo.
$serviceCommand = """$executablePath"" ""$configPath"""

Write-Host "Preparando para instalar el servicio..."
Write-Host "  - Servicio: $ServiceName"
Write-Host "  - Ejecutable: $executablePath"
Write-Host "  - Configuracion: $configPath"
Write-Host ""

try {
    # --- 4. Detener y Eliminar el Servicio si ya Existe (para una instalacion limpia) ---
    $existingService = Get-Service -Name $ServiceName -ErrorAction SilentlyContinue
    if ($existingService) {
        Write-Host "⚠️  Servicio existente encontrado. Se detendra y reinstalara para asegurar una configuracion limpia." -ForegroundColor Yellow
        Stop-Service -Name $ServiceName -Force -ErrorAction SilentlyContinue
        Remove-Service -Name $ServiceName -Force -ErrorAction SilentlyContinue
        Write-Host "   Servicio anterior eliminado. Esperando unos segundos..."
        Start-Sleep -Seconds 5
    }

    # --- 5. Crear el Nuevo Servicio ---
    Write-Host "Creando nuevo servicio '$DisplayName'..." -ForegroundColor Green
    New-Service -Name $ServiceName `
                -BinaryPathName $serviceCommand `
                -DisplayName $DisplayName `
                -Description $Description `
                -StartupType Automatic

    # --- 6. Configurar Recuperacion Automatica ante Fallos ---
    Write-Host "Configurando recuperacion automatica..."
    # Si el servicio falla, intentara reiniciarse 3 veces.
    sc.exe failure $ServiceName reset= 86400 actions= restart/60000/restart/60000/restart/60000

    # --- 7. Iniciar el Servicio ---
    Write-Host "Iniciando servicio..."
    Start-Service -Name $ServiceName

    # --- 8. Verificar Estado Final ---
    Start-Sleep -Seconds 2
    $service = Get-Service -Name $ServiceName
    Write-Host ""
    Write-Host "✅ ¡Servicio instalado y arrancado con exito!" -ForegroundColor Green
    Write-Host "   - Nombre: $($service.Name)"
    Write-Host "   - Estado: $($service.Status)"
    Write-Host "   - Tipo de Inicio: $($service.StartType)"
    Write-Host ""
    Write-Host "El servicio ahora se iniciara automaticamente con Windows."

} catch {
    Write-Error "❌ ERROR al instalar el servicio: $($_.Exception.Message)"
    Write-Error "Asegurate de estar ejecutando este script como Administrador."
    pause
    exit 1
}

Write-Host ""
Write-Host "Comandos utiles:" -ForegroundColor Cyan
Write-Host "  - Ver estado: Get-Service $ServiceName"
Write-Host "  - Detener:    Stop-Service $ServiceName"
Write-Host "  - Iniciar:    Start-Service $ServiceName"
Write-Host "  - Eliminar:   Remove-Service $ServiceName"
Write-Host ""
pause