@echo off
REM Script para compilar y preparar SCP Go Identity Service para despliegue en Windows.

echo.
echo ========================================
echo  SCP Go Identity Service - Build & Deploy
echo ========================================
echo.

REM --- 1. Limpieza del directorio de despliegue anterior ---
if exist "deploy" (
    echo 1. Limpiando directorio de despliegue anterior...
    RMDIR /S /Q deploy
)

REM --- 2. Compilacion del ejecutable ---
echo.
echo 2. Compilando para Windows x64...
set CGO_ENABLED=0
set GOOS=windows
set GOARCH=amd64
go build -ldflags="-s -w" -o scp-go-identity-service.exe main.go

if %errorlevel% neq 0 (
    echo.
    echo ❌ ERROR: Fallo en la compilacion.
    pause
    exit /b 1
)
echo    ✅ Ejecutable 'scp-go-identity-service.exe' creado.

REM --- 3. Preparando la carpeta de despliegue ---
echo.
echo 3. Creando estructura de despliegue en la carpeta 'deploy'...
mkdir deploy

copy scp-go-identity-service.exe deploy\ >nul
del scp-go-identity-service.exe

REM Copia la configuracion de servidor y la renombra a 'config.json' que es la que busca el ejecutable.
if exist "config-server-internal.json" (
    copy config-server-internal.json deploy\config.json >nul
    echo    ✅ 'config-server-internal.json' copiado como 'config.json'.
) else (
    copy config.json deploy\config.json >nul
    echo    ⚠️  ADVERTENCIA: No se encontro 'config-server-internal.json'. Se uso 'config.json' por defecto.
)

copy install-service.ps1 deploy\ >nul
echo    ✅ Script 'install-service.ps1' copiado.

echo.
echo ========================================
echo  ✅ COMPILACION COMPLETADA
echo ========================================
echo.
echo 📁 Archivos listos en la carpeta 'deploy'.
echo    Contenido:
dir deploy /b
echo.
echo 🚀 PROXIMOS PASOS:
echo.
echo    1. Ve a la carpeta 'deploy'.
echo    2. Para PRODUCCION, NO edites 'config.json' con datos sensibles.
echo       En su lugar, configura las variables de entorno en el sistema donde correra el servicio.
echo       - SCP_JWT_KEY: La clave secreta para firmar los tokens.
echo       - SCP_DB_CONNECTION_STRING: La cadena de conexion a la base de datos.
echo    3. Abre una terminal de PowerShell COMO ADMINISTRADOR en la carpeta 'deploy'.
echo    4. Ejecuta el siguiente comando para instalar el servicio:
echo.
echo       PowerShell -ExecutionPolicy Bypass -File .\install-service.ps1
echo.
echo 🌐 El servicio (configurado para usar las variables de entorno) se iniciara automaticamente con Windows.
echo.

pause