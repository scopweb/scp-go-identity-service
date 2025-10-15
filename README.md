# SCP Go Identity Service

Este es un servicio de autenticación en Go que se conecta a una base de datos de Microsoft .NET Identity existente para autenticar usuarios y generar tokens JWT. 

**🏢 Diseñado para entorno de servidor interno donde aplicaciones IIS/.NET pueden consumir el servicio de autenticación de forma independiente de internet.**

## Características

- ✅ Conexión a base de datos SQL Server con tablas de .NET Identity
- ✅ Verificación de contraseñas con hash de .NET Identity (PBKDF2-SHA256)
- ✅ Generación de tokens JWT compatible con .NET
- ✅ Recuperación de roles y claims de usuario
- ✅ Logging detallado de peticiones
- ✅ Endpoints RESTful

## Estructura del Proyecto

```
├── main.go                 # Punto de entrada del servicio
├── config/
│   └── config.go          # Configuración del servicio
├── models/
│   ├── user.go           # Modelos de tablas .NET Identity
│   └── auth.go           # Modelos de request/response
├── services/
│   ├── database.go       # Servicio de base de datos
│   ├── password.go       # Verificación de contraseñas .NET
│   └── jwt.go           # Generación de tokens JWT
├── handlers/
│   └── auth.go          # Manejadores HTTP
└── config.json          # Archivo de configuración
```

## Configuración

El servicio sigue una estrategia de configuración por capas para máxima seguridad y flexibilidad:

1.  **Archivo `config.json`**: Carga la configuración base. Ideal para valores no sensibles o de desarrollo.
2.  **Variables de Entorno**: Sobrescribe los valores del archivo si están presentes. **Este es el método preferido para producción.**

### Variables de Entorno (Recomendado para Producción)

-   `SCP_JWT_KEY`: Clave secreta para firmar los tokens JWT.
-   `SCP_DB_CONNECTION_STRING`: Cadena de conexión a la base de datos.

### Ejemplo de `config.json`

```json
{
  "database": {
    "connectionString": "server=localhost;database=DB_IDENTITY;user id=sa;password=YourPassword;encrypt=false"
  },
  "jwt": {
    "key": "una-clave-para-desarrollo-puede-ir-aqui-pero-sera-sobrescrita-en-produccion",
    "issuer": "https://www.domain.com",
    "audience": "https://www.domain.com",
    "expiryInHours": 12
  },
  "server": {
    "host": "0.0.0.0",
    "port": "8080"
  }
}
```

## 🚀 Instalación en Servidor Interno

### Prerrequisitos
- Windows Server con IIS
- SQL Server con base de datos .NET Identity existente
- Aplicaciones .NET corriendo en IIS que necesitan autenticación

### Instalación Rápida

1. **📥 Preparar el servicio:**
```batch
# Ejecutar el script de build
build.bat
```

2. **⚙️ Configurar para tu servidor:**
   - Editar `deploy/config-server-internal.json`
   - Actualizar la cadena de conexión a tu base de datos
   - Verificar puerto (por defecto: 8081)

3. **🔧 Instalar como Servicio de Windows:**
```powershell
# Ejecutar como Administrador
cd deploy
PowerShell -ExecutionPolicy Bypass -File install-service.ps1
```

4. **✅ Verificar instalación:**
   - Abrir http://localhost:8081/health
   - Revisar Windows Services para "SCP Go Identity Service"

### Configuración Manual (Desarrollo)

1. **Instalar dependencias:**
```bash
go mod tidy
```

2. **Configurar base de datos:**
   - Copiar `config-server-internal.json` a `config.json`
   - Actualizar cadena de conexión
   
3. **Ejecutar en desarrollo:**
```bash
go run main.go config-server-internal.json
```

## API Endpoints

### POST /authenticate

Autentica un usuario y opcionalmente genera un token JWT.

**Request:**
```json
{
  "email": "usuario@ejemplo.com",
  "password": "contraseña123",
  "generateToken": true
}
```

**Response (éxito):**
```json
{
  "success": true,
  "message": "Authentication successful",
  "user": {
    "id": "user-id-guid",
    "email": "usuario@ejemplo.com",
    "userName": "usuario",
    "firstName": "Nombre",
    "lastName": "Apellido",
    "culture": "es-ES",
    "isEmailConfirmed": true,
    "phoneNumber": "+1234567890"
  },
  "roles": ["Admin", "User"],
  "claims": [
    {
      "type": "custom-claim",
      "value": "custom-value"
    }
  ],
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "expiresAt": "2024-01-01T12:00:00Z"
}
```

**Response (error):**
```json
{
  "success": false,
  "message": "Invalid credentials"
}
```

### GET /test

**🧪 Página de prueba de autenticación - ¡Diseño minimalista oscuro!**

Página HTML para probar la autenticación con formulario de usuario y contraseña.

**URL:** `http://localhost:8081/test`

Características:
- ✅ Diseño minimalista y serio en tema oscuro
- ✅ Formulario completo de autenticación
- ✅ Visualización de respuestas del servidor
- ✅ Indicador de estado del servicio en tiempo real
- ✅ Validación de campos y manejo de errores

### GET /status

**🌐 Página web "Estoy funcionando" - ¡Abrir en navegador!**

Página HTML bonita para verificar rápidamente que el servicio está funcionando.

**URL:** `http://localhost:8081/status`

Características:
- ✅ Interfaz visual atractiva
- ✅ Información del servidor en tiempo real
- ✅ Auto-refresh cada 30 segundos
- ✅ Lista de todos los endpoints disponibles
- ✅ Enlace directo a página de pruebas

### GET /health

Verifica el estado del servicio (formato JSON para aplicaciones).

**Response:**
```json
{
  "status": "healthy",
  "service": "scp-go-identity-service"
}
```

### GET /

Información general del servicio (formato JSON).

**Response:**
```json
{
  "message": "SCP Go Identity Service",
  "version": "1.0",
  "endpoints": ["/authenticate", "/health", "/status"],
  "status": "running"
}
```

## Características de Seguridad

- ✅ **Algoritmo de Firma JWT Robusto**: Utiliza `HS512` para la firma de tokens, alineado con los estándares de hashing de contraseñas.
- ✅ **Protección contra Enumeración de Usuarios**: El servicio mitiga los ataques de enumeración de usuarios al garantizar tiempos de respuesta similares para usuarios existentes y no existentes.
- ✅ **Limitación de Tasa (Rate Limiting)**: Incorpora un middleware que limita las solicitudes por IP para proteger contra ataques de fuerza bruta y DoS.
- ✅ **Gestión Segura de Secretos**: Prioriza el uso de variables de entorno (`JJP_JWT_KEY`, `JJP_DB_CONNECTION_STRING`) para cargar configuraciones sensibles, evitando que queden expuestas en archivos.
- ✅ **Validación de Audiencia de Token**: Valida que los tokens JWT estén destinados específicamente a este servicio.
- ✅ Verificación de contraseñas usando el mismo algoritmo que .NET Identity.
- ✅ Soporte para lockout de cuentas.
- ✅ Tokens JWT con expiración configurable.
- ✅ CORS habilitado para aplicaciones web.
- ✅ Logging de todas las peticiones de autenticación.

## Formato de Token JWT

El token JWT incluye los siguientes claims:

- `nameid`: ID del usuario
- `email`: Email del usuario
- `unique_name`: Nombre de usuario
- `role`: Roles del usuario (array)
- `given_name`: Nombre (opcional)
- `family_name`: Apellido (opcional)
- `culture`: Cultura del usuario (opcional)

## Logging

El servicio registra información detallada incluyendo:

- ID de petición único para tracking
- IP del cliente y User-Agent
- Resultado de autenticación (éxito/fallo)
- Errores de base de datos y validación
- Generación de tokens

## Dependencias

- `github.com/golang-jwt/jwt/v5` - Generación y validación JWT
- `github.com/denisenkom/go-mssqldb` - Driver SQL Server
- `github.com/google/uuid` - Generación de UUIDs
- `golang.org/x/crypto` - Verificación de contraseñas PBKDF2

## 🔗 Integración desde Aplicaciones IIS/.NET

### Consumir desde aplicación ASP.NET

```csharp
// En tu aplicación IIS, agregar el cliente del servicio
var identityService = new SCPIdentityService("http://localhost:8081");

// Autenticar usuario
var result = await identityService.AuthenticateAsync(email, password, true);

if (result.Success)
{
    // Usuario autenticado - crear sesión/cookies
    Session["UserId"] = result.User.Id;
    Session["UserRoles"] = result.Roles;
    Session["AuthToken"] = result.Token;
}
```

### Web.config - Configuración de conexión

```xml
<appSettings>
    <add key="SCP.IdentityService.BaseUrl" value="http://localhost:8081" />
    <add key="SCP.IdentityService.Timeout" value="30" />
</appSettings>
```

### Ejemplo completo disponible en:
📄 `dotnet-client-example.cs` - Implementación completa del cliente .NET

## 🛡️ Seguridad para Servidor Interno

- ✅ Servicio bind a localhost (no expuesto externamente)
- ✅ Sin dependencias de internet
- ✅ Comunicación HTTP interna entre IIS y servicio Go
- ✅ Misma base de datos .NET Identity existente
- ✅ Tokens JWT para sesiones seguras

## Compatibilidad

Compatible con aplicaciones .NET Identity existentes:
- ✅ Tablas de base de datos estándar de .NET Identity  
- ✅ Hash de contraseñas PBKDF2-SHA256
- ✅ Estructura de roles y claims
- ✅ Tokens JWT compatibles con .NET
- ✅ Funciona con SQL Server local o remoto
- ✅ Integración transparente con aplicaciones IIS

## Desarrollo

Para desarrollo local con recarga automática:

```bash
# Instalar air para hot reload
go install github.com/cosmtrek/air@latest

# Ejecutar con recarga automática
air
```

## Docker

```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o scp-go-identity-service main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/scp-go-identity-service .
COPY --from=builder /app/config.json .
EXPOSE 8080
CMD ["./scp-go-identity-service"]
```

## Nota Crítica sobre el Hashing de Contraseñas

Durante el desarrollo de este servicio, se encontró un comportamiento específico y no estándar en la implementación de `PasswordHasher` de **.NET 9**. Es crucial documentarlo para futuras referencias y mantenimiento.

### El Descubrimiento

El sistema .NET 9 utilizado para generar los hashes de contraseña produce un hash que, aunque sigue el formato V3 de ASP.NET Identity, tiene una peculiaridad:

- **Identificador PRF:** El hash se marca con el identificador `0x00000002`, que corresponde a **HMAC-SHA512**.
- **Longitud de la Subclave:** A pesar de usar HMAC-SHA512 (que normalmente produce una subclave de 64 bytes), la implementación de .NET solicita y genera una subclave de solo **32 bytes** (256 bits).

Esta combinación es inusual. El hash es internamente consistente pero no sigue la convención estándar donde la longitud de la subclave coincide con la salida natural del algoritmo hash.

### Implementación en Go (`services/password.go`)

Para garantizar la compatibilidad, la función `verifyV3` en este proyecto fue diseñada específicamente para replicar este comportamiento:

1.  **Lee el Identificador PRF** del hash para determinar el algoritmo a utilizar (ej. `sha512.New`).
2.  Lee los demás parámetros como las iteraciones y el salt.
3.  Extrae la subclave almacenada.
4.  **Utiliza la longitud real de la subclave extraída** como el parámetro `keyLength` para la función `pbkdf2.Key`.

### ¿Qué hacer si la autenticación falla en el futuro?

Si las contraseñas comienzan a fallar después de una actualización del sistema .NET, es muy probable que la lógica de hashing haya cambiado. Los pasos para diagnosticar serían:

1.  **Obtener un hash de prueba:** Genera un hash para una contraseña conocida (ej. "123456") en el nuevo entorno .NET.
2.  **Analizar el nuevo hash:** Decodifícalo de Base64 y examina sus componentes:
    *   Identificador PRF (¿Ha cambiado el algoritmo? ¿Es ahora SHA256 o SHA1?).
    *   Número de iteraciones.
    *   Longitud del salt.
    *   **Longitud de la subclave resultante.**
3.  **Actualizar `services/password.go`:** Modifica la función `verifyV3` para que coincida con los nuevos parámetros descubiertos. Es probable que solo necesites ajustar la lógica que determina el `keyLength` o el `prf`.