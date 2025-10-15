# Guía de Mantenimiento para el Asistente de IA (Claude)

Hola, Claude. Esta es tu guía de referencia para mantener y actualizar el `scp-go-identity-service`.

## 1. Resumen del Servicio

Este es un servicio de autenticación escrito en Go. Su propósito principal es actuar como un puente para sistemas de **.NET Identity**, permitiendo la autenticación de usuarios desde una base de datos existente de .NET Identity y la generación de tokens JWT.

**Tecnologías Clave:**
- **Lenguaje:** Go
- **Base de Datos:** SQL Server (compatible con esquema .NET Identity)
- **Hashing de Contraseña:** PBKDF2 (compatible con formatos V2 y V3 de .NET Identity)
- **Tokens:** JWT (`github.com/golang-jwt/jwt/v5`)
- **Servidor:** `net/http` estándar de Go

## 2. Arquitectura y Flujo de Autenticación

1.  **Petición HTTP:** El cliente envía un `POST` a `/authenticate` con `email` y `password`.
2.  **Middleware:** La petición pasa por los siguientes middlewares en `main.go`:
    -   `loggingHandler`: Registra la petición.
    -   `corsHandler`: Aplica cabeceras CORS.
    -   `RateLimitMiddleware`: Limita las peticiones por IP (`handlers/ratelimit.go`).
3.  **Handler:** `authHandler.Authenticate` (`handlers/auth.go`) gestiona la lógica.
4.  **Búsqueda de Usuario:** Se busca al usuario por email en la base de datos (`services/database.go`).
5.  **Protección Enumeración:** Si el usuario no existe, se ejecuta una comparación de hash "falsa" para mantener un tiempo de respuesta constante y evitar ataques de enumeración.
6.  **Verificación de Contraseña:** `passwordService.VerifyPassword` (`services/password.go`) compara la contraseña proporcionada con el hash almacenado. **Esta es una parte crítica (ver sección 4).**
7.  **Generación de Token:** Si la autenticación es exitosa y se solicita, `jwtService.GenerateToken` (`services/jwt.go`) crea un token JWT firmado con **HS512**.
8.  **Respuesta:** Se devuelve una respuesta JSON con el resultado.

## 3. Gestión de la Configuración y Secretos

La configuración es jerárquica para máxima seguridad. El proceso de carga en `config/config.go` es:

1.  Lee el archivo `config.json`.
2.  **Sobrescribe** valores si existen las siguientes variables de entorno:
    -   `SCP_JWT_KEY`: Para la clave secreta de firma de JWT.
    -   `SCP_DB_CONNECTION_STRING`: Para la cadena de conexión a la base de datos.

**Tu Tarea en Mantenimiento:**
- **Nunca** escribas secretos directamente en `config.json`.
- Al desplegar o dar instrucciones de ejecución, siempre utiliza el método de variables de entorno para los secretos.

## 4. Punto Crítico de Mantenimiento: Hashing de Contraseñas

El componente más frágil de este servicio es la compatibilidad con el hashing de contraseñas de .NET Identity, ya que puede cambiar entre versiones de .NET.

**El Problema:**
La función `verifyV3` en `services/password.go` está diseñada para una implementación específica de .NET que usa `HMAC-SHA512` pero genera una subclave de solo `32 bytes` (en lugar de los 64 esperados).

**Si la autenticación de contraseñas comienza a fallar, sigue estos pasos de diagnóstico:**

1.  **Obtén un Hash de Muestra:** Pide al usuario que genere un hash de contraseña para una clave conocida (ej: `123456`) en su entorno .NET actual.
2.  **Analiza el Hash:**
    -   Decodifica la cadena Base64 del hash.
    -   Examina los bytes del encabezado (los primeros 13 bytes):
        -   `bytes[1:5]`: Identificador PRF (0=SHA1, 1=SHA256, 2=SHA512).
        -   `bytes[5:9]`: Número de iteraciones.
        -   `bytes[9:13]`: Longitud del salt.
    -   Determina la **longitud de la subclave** (hash real), que son los bytes restantes después del salt.
3.  **Actualiza el Código:**
    -   Abre `services/password.go`.
    -   Modifica la función `verifyV3` para que coincida con los parámetros que has descubierto. Lo más probable es que solo necesites ajustar la lógica que determina el `keyLength` o el `prf` a utilizar en la llamada a `pbkdf2.Key`.

## 5. Proceso de Build y Despliegue

-   **Build:** Ejecuta `build.bat`. Este script compila el ejecutable para Windows y lo coloca en una carpeta `deploy` junto con los scripts de instalación.
-   **Despliegue:**
    1.  El administrador del sistema debe configurar las variables de entorno `SCP_JWT_KEY` y `SCP_DB_CONNECTION_STRING` en el servidor de producción.
    2.  Desde la carpeta `deploy`, se debe ejecutar `install-service.ps1` como Administrador para instalarlo como un servicio de Windows.

**Tu Tarea en Mantenimiento:**
- Si añades nuevos archivos necesarios para la ejecución (ej. plantillas HTML), asegúrate de **actualizar `build.bat`** para que los copie a la carpeta `deploy`.

## 6. Estructura del Código

-   `main.go`: Punto de entrada, configuración de rutas y middlewares.
-   `/handlers`: Manejadores de peticiones HTTP.
    -   `auth.go`: Lógica principal de autenticación.
    -   `ratelimit.go`: Middleware de limitación de tasa.
-   `/services`: Lógica de negocio desacoplada.
    -   `database.go`: Interacción con la base de datos.
    -   `password.go`: Verificación de contraseñas. **(Punto crítico)**
    -   `jwt.go`: Generación y validación de tokens JWT.
-   `/models`: Estructuras de datos (Go structs).
-   `/config`: Lógica de carga de configuración.
-   `build.bat`: Script de compilación y empaquetado.
-   `install-service.ps1`: Script de instalación como servicio de Windows.
-   `README.md` y `DEPLOYMENT-GUIDE.md`: Documentación para el usuario final. Mantenla actualizada.