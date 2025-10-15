# 🚀 Guía Rápida de Despliegue - SCP Go Identity Service

Esta guía describe cómo desplegar el servicio de forma segura en un entorno de producción.

## Para el Administrador del Servidor

### 1. Preparación (5 minutos)
Ejecuta el script de compilación desde el directorio del proyecto.
```batch
# Desde el directorio del proyecto
build.bat
```
Esto crea una carpeta `deploy` con todos los archivos necesarios para la instalación.

### 2. Configuración de Secretos (Método Recomendado)
En lugar de editar `config.json` directamente, configura los siguientes secretos como **variables de entorno del sistema** en el servidor donde se ejecutará el servicio.

1.  Abre "Propiedades del sistema" > "Variables de entorno...".
2.  En "Variables del sistema", crea dos nuevas variables:

    -   **`SCP_JWT_KEY`**: Pega aquí una clave secreta larga y compleja para firmar los tokens JWT.
        *Ejemplo de clave segura (¡genera la tuya!):* `B3z@_V-5!pL$G*m#q&R(s+T)u/W?X#Z`

    -   **`SCP_DB_CONNECTION_STRING`**: La cadena de conexión a la base de datos de producción.
        *Ejemplo:* `server=PROD_SERVER\\SQLEXPRESS;database=SCP_IDENTITY;user id=svc_user;password=P@ssw0rd_S3gur@;encrypt=true;TrustServerCertificate=false`

**Nota:** El archivo `config.json` dentro de la carpeta `deploy` puede dejarse como está, ya que las variables de entorno tendrán prioridad.

### 3. Instalación del Servicio (1 minuto)
Abre una terminal de PowerShell **como Administrador** y navega a la carpeta `deploy`. Luego, ejecuta:
```powershell
PowerShell -ExecutionPolicy Bypass -File install-service.ps1
```
El servicio se instalará y se iniciará automáticamente, cargando los secretos desde las variables de entorno.

### 4. Verificación
- 🧪 **Prueba de autenticación:** http://localhost:8081/test (¡PROBAR AQUÍ!)
- 🌐 **Página de estado:** http://localhost:8081/status 
- 📱 **API Health:** http://localhost:8081/health  
- Debe mostrar formularios funcionales y estado "online".

## Para el Desarrollador .NET

### URL del Servicio Interno
```
http://localhost:8081/authenticate
```

### Ejemplo de Uso en C#
```csharp
var client = new HttpClient();
var request = new {
    email = "usuario@empresa.com",
    password = "contraseña123", 
    generateToken = true
};

var json = JsonConvert.SerializeObject(request);
var content = new StringContent(json, Encoding.UTF8, "application/json");
var response = await client.PostAsync("http://localhost:8081/authenticate", content);
```

Ver ejemplo completo en: `dotnet-client-example.cs`

## URLs Importantes

| Endpoint | Propósito | Método | Tipo |
|----------|-----------|--------|------|
| `/test` | **Página de prueba de autenticación** 🧪 | GET | HTML |
| `/status` | Página "Estoy funcionando" 🌐 | GET | HTML |
| `/health` | Estado del servicio | GET | JSON |
| `/authenticate` | Autenticación | POST | JSON |
| `/` | Info del servicio | GET | JSON |

## Comandos de Servicio Windows

```powershell
# Ver estado
Get-Service SCPGoIdentityService

# Reiniciar
Restart-Service SCPGoIdentityService

# Detener
Stop-Service SCPGoIdentityService

# Iniciar
Start-Service SCPGoIdentityService

# Ver logs
Get-EventLog -LogName Application -Source SCPGoIdentityService -Newest 10
```

## Troubleshooting

### ❌ Error: "No se puede conectar al servicio"
- Verificar que el servicio esté corriendo: `Get-Service SCPGoIdentityService`
- Verificar el puerto: `netstat -an | findstr :8081`

### ❌ Error: "Database connection failed"
- Verificar cadena de conexión en config-server-internal.json
- Verificar que SQL Server esté corriendo
- Verificar permisos de la cuenta del servicio

### ❌ Error: "Invalid credentials" 
- Verificar que el usuario existe en AspNetUsers
- Verificar que la contraseña sea correcta
- Revisar logs del servicio

## 📞 Soporte
- Ver logs completos del servicio en Windows Event Viewer
- Aplicación > SCPGoIdentityService
- O usar PowerShell: `Get-EventLog -LogName Application -Source SCPGoIdentityService`