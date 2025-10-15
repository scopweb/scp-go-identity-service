# Changelog

Todos los cambios notables en este proyecto serán documentados en este archivo.

El formato está basado en [Keep a Changelog](https://keepachangelog.com/es-ES/1.0.0/),
y este proyecto adhiere a [Versionado Semántico](https://semver.org/lang/es/).

## [0.1.0] - 2025-01-15

### 🎉 Lanzamiento Inicial

Primer lanzamiento público de SCP Go Identity Service.

### ✨ Características Principales

#### Autenticación
- Autenticación de usuarios contra base de datos .NET Identity existente
- Verificación de contraseñas compatible con PBKDF2 de .NET Identity (V2 y V3)
- Soporte para hash formato específico de .NET 9 con HMAC-SHA512
- Protección contra ataques de enumeración de usuarios (timing attack mitigation)
- Soporte para lockout de cuentas

#### JWT (JSON Web Tokens)
- Generación de tokens JWT con algoritmo HS512
- Claims personalizables (roles, email, nombre, cultura)
- Expiración configurable de tokens
- Validación de audiencia e issuer

#### API RESTful
- `POST /authenticate` - Endpoint de autenticación
- `GET /health` - Health check (JSON)
- `GET /status` - Página de estado visual (HTML)
- `GET /test` - Página de prueba interactiva (HTML)
- `GET /` - Información del servicio (JSON)

#### Seguridad
- Rate limiting por IP para prevenir ataques de fuerza bruta
- CORS habilitado y configurable
- Gestión segura de secretos mediante variables de entorno
- Logging detallado de todas las peticiones

#### Configuración
- Configuración jerárquica (archivo JSON + variables de entorno)
- Variables de entorno: `SCP_JWT_KEY`, `SCP_DB_CONNECTION_STRING`
- Hot-reload de configuración no implementado (requiere reinicio)

#### Base de Datos
- Conexión a SQL Server
- Soporte para tablas estándar de .NET Identity:
  - AspNetUsers
  - AspNetRoles
  - AspNetUserRoles
  - AspNetUserClaims
- Connection pooling automático
- Timeout configurable

#### Despliegue
- Instalación como servicio de Windows (`install-service.ps1`)
- Script de build automatizado (`build.bat`)
- Soporte para Docker (Dockerfile multi-stage)
- Configuración para servidor interno (localhost binding)

#### Herramientas y Utilidades
- Scripts de diagnóstico (PowerShell)
- Página de prueba interactiva con interfaz web
- Ejemplo de cliente .NET (C#)
- Script SQL para crear usuario de prueba

#### Documentación
- README completo en español
- Guía de despliegue paso a paso
- Guía de mantenimiento para asistentes de IA
- Documentación de arquitectura y decisiones técnicas
- Ejemplos de código para integración

### 🔧 Tecnologías

- **Lenguaje:** Go 1.21+
- **Base de Datos:** SQL Server (driver `go-mssqldb`)
- **JWT:** `github.com/golang-jwt/jwt/v5`
- **Criptografía:** `golang.org/x/crypto` (PBKDF2)
- **UUID:** `github.com/google/uuid`

### 📦 Dependencias

```go
require (
    github.com/denisenkom/go-mssqldb v0.12.3
    github.com/golang-jwt/jwt/v5 v5.3.0
    github.com/google/uuid v1.6.0
    golang.org/x/crypto v0.43.0
)
```

### 🐛 Problemas Conocidos

- La versión del API en `/` retorna "1.0" en lugar de "0.1" (inconsistencia menor)
- No hay implementación de refresh tokens
- El rate limiting es en memoria (se pierde al reiniciar el servicio)
- No hay soporte para Two-Factor Authentication (2FA) aunque .NET Identity lo soporte

### 📝 Notas

- **Compatibilidad:** Diseñado específicamente para .NET Identity con esquema estándar
- **Entorno:** Optimizado para servidores internos (no expuesto a internet)
- **Plataforma:** Principalmente Windows Server con IIS, pero funciona en Linux y Docker
- **Base de Datos:** Requiere una base de datos .NET Identity existente

### 🙏 Agradecimientos

Gracias a la comunidad de Go y .NET por las excelentes bibliotecas y documentación.

---

## [Unreleased]

### Planeado para Futuras Versiones

- [ ] Implementación de refresh tokens
- [ ] Soporte para 2FA/MFA
- [ ] Rate limiting persistente (Redis/base de datos)
- [ ] Métricas y observabilidad (Prometheus)
- [ ] Cache de usuarios autenticados
- [ ] Soporte para múltiples bases de datos Identity
- [ ] API para gestión de usuarios (crear, actualizar, eliminar)
- [ ] Integración con OAuth2/OIDC
- [ ] Soporte para claims dinámicos
- [ ] Hot-reload de configuración

---

[0.1.0]: https://github.com/scopweb/scp-go-identity-service/releases/tag/v0.1.0
[Unreleased]: https://github.com/scopweb/scp-go-identity-service/compare/v0.1.0...HEAD
