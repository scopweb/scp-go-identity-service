package main

import (
	"fmt"
	"scp-go-identity-service/config"
	"scp-go-identity-service/handlers"
	"scp-go-identity-service/services"
	"log"
	"net/http"
	"os"
)

func main() {
	// Load configuration
	configPath := "config.json"
	if len(os.Args) > 1 {
		configPath = os.Args[1]
	}

	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	log.Printf("🚀 Starting SCP Go Identity Service for Internal IIS Applications...")
	log.Printf("📍 Server Configuration:")
	log.Printf("   - Host: %s", cfg.Server.Host)
	log.Printf("   - Port: %s", cfg.Server.Port)
	log.Printf("   - Internal Service URL: http://%s:%s", cfg.Server.Host, cfg.Server.Port)
	log.Printf("🔐 JWT Configuration:")
	log.Printf("   - Issuer: %s", cfg.JWT.Issuer)
	log.Printf("   - Audience: %s", cfg.JWT.Audience)
	log.Printf("   - Token Expiry: %d hours", cfg.JWT.ExpiryInHours)

	// Initialize database service
	dbService, err := services.NewDatabaseService(cfg.Database.ConnectionString)
	if err != nil {
		log.Fatalf("Failed to initialize database service: %v", err)
	}
	defer dbService.Close()
	log.Printf("✅ Database connection established")

	// Initialize password service
	passwordService := services.NewPasswordService()
	log.Printf("✅ Password service initialized")

	// Initialize JWT service
	jwtService := services.NewJWTService(cfg.JWT)
	log.Printf("✅ JWT service initialized")

	// Initialize auth handler
	authHandler := handlers.NewAuthHandler(dbService, passwordService, jwtService)
	log.Printf("✅ Authentication handler initialized")

	// Setup HTTP routes
	mux := http.NewServeMux()

	// Authentication endpoint
	mux.HandleFunc("/authenticate", authHandler.Authenticate)

	// Health check endpoint (JSON)
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"status":"healthy","service":"scp-go-identity-service"}`)
	})

	// Status page - "Estoy funcionando" (HTML)
	mux.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		statusPage := `<!DOCTYPE html>
<html lang="es">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>SCP Go Identity Service - Estado</title>
    <style>
        body {
            font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
            margin: 0;
            padding: 0;
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            display: flex;
            justify-content: center;
            align-items: center;
            min-height: 100vh;
            color: white;
        }
        .container {
            text-align: center;
            background: rgba(255, 255, 255, 0.1);
            padding: 3rem;
            border-radius: 20px;
            box-shadow: 0 8px 32px rgba(0, 0, 0, 0.2);
            backdrop-filter: blur(10px);
            border: 1px solid rgba(255, 255, 255, 0.2);
            max-width: 600px;
            width: 90%;
        }
        .status-icon {
            font-size: 4rem;
            margin-bottom: 1rem;
            animation: pulse 2s infinite;
        }
        @keyframes pulse {
            0% { transform: scale(1); }
            50% { transform: scale(1.1); }
            100% { transform: scale(1); }
        }
        h1 {
            font-size: 2.5rem;
            margin-bottom: 1rem;
            text-shadow: 2px 2px 4px rgba(0, 0, 0, 0.3);
        }
        .service-info {
            background: rgba(255, 255, 255, 0.1);
            padding: 1.5rem;
            border-radius: 10px;
            margin: 2rem 0;
            text-align: left;
        }
        .endpoint {
            margin: 0.5rem 0;
            font-family: 'Courier New', monospace;
            background: rgba(0, 0, 0, 0.2);
            padding: 0.5rem;
            border-radius: 5px;
        }
        .timestamp {
            font-size: 0.9rem;
            opacity: 0.8;
            margin-top: 2rem;
        }
        .version {
            position: absolute;
            bottom: 20px;
            right: 20px;
            font-size: 0.8rem;
            opacity: 0.6;
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="status-icon">✅</div>
        <h1>¡Estoy Funcionando!</h1>
        <p><strong>SCP Go Identity Service</strong> está operativo y listo para recibir peticiones.</p>
        
        <div class="service-info">
            <h3>📋 Endpoints Disponibles:</h3>
            <div class="endpoint">POST /authenticate - Autenticación de usuarios</div>
            <div class="endpoint">GET /health - Verificación de estado (JSON)</div>
            <div class="endpoint">GET /status - Esta página</div>
            <div class="endpoint">GET /test - Página de prueba de autenticación</div>
            <div class="endpoint">GET / - Información del servicio (JSON)</div>
        </div>
        
        <div class="service-info" style="text-align: center;">
            <a href="/test" style="display: inline-block; background: rgba(255,255,255,0.2); color: white; padding: 10px 20px; text-decoration: none; border-radius: 5px; margin: 10px;">🧪 Probar Autenticación</a>
        </div>

        <div class="service-info">
            <h3>🔧 Información del Servidor:</h3>
            <p><strong>Host:</strong> ` + cfg.Server.Host + `</p>
            <p><strong>Puerto:</strong> ` + cfg.Server.Port + `</p>
            <p><strong>JWT Issuer:</strong> ` + cfg.JWT.Issuer + `</p>
            <p><strong>Expiry:</strong> ` + fmt.Sprintf("%d", cfg.JWT.ExpiryInHours) + ` horas</p>
        </div>

        <div class="timestamp">
            🕒 Última verificación: <span id="timestamp"></span>
        </div>
    </div>

    <div class="version">v0.1</div>

    <script>
        // Actualizar timestamp
        function updateTimestamp() {
            document.getElementById('timestamp').textContent = new Date().toLocaleString('es-ES');
        }
        updateTimestamp();
        setInterval(updateTimestamp, 1000);

        // Auto-refresh cada 30 segundos
        setTimeout(() => {
            window.location.reload();
        }, 30000);
    </script>
</body>
</html>`
		fmt.Fprint(w, statusPage)
	})

	// Authentication test page (HTML)
	mux.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "auth-test.html")
	})

	// Debug endpoint - database connection test
	mux.HandleFunc("/debug", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		// Test database connection
		if err := dbService.TestConnection(); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, `{"status":"error","database":"disconnected","error":"%s"}`, err.Error())
		} else {
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, `{"status":"ok","database":"connected","config":"loaded"}`)
		}
	}) // Root endpoint (JSON info)
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"message":"SCP Go Identity Service","version":"0.1","endpoints":["/authenticate","/health","/status","/test"],"status":"running"}`)
	})

	// Add CORS middleware
	corsHandler := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}

			next.ServeHTTP(w, r)
		})
	}

	// Add logging middleware
	loggingHandler := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Printf("%s %s %s", r.Method, r.URL.Path, r.RemoteAddr)
			next.ServeHTTP(w, r)
		})
	}

	// Wrap with middleware
	handler := loggingHandler(corsHandler(handlers.RateLimitMiddleware(mux)))

	// Start server
	address := fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)
	log.Printf("🚀 Server starting on %s", address)
	log.Printf("📋 Available endpoints:")
	log.Printf("   POST /authenticate - User authentication")
	log.Printf("   GET  /health       - Health check (JSON)")
	log.Printf("   GET  /status       - Status page (HTML)")
	log.Printf("   GET  /test         - Authentication test page (HTML) - ¡PROBAR AQUÍ!")
	log.Printf("   GET  /             - Service information (JSON)")
	log.Printf("")
	log.Printf("🧪 Para probar autenticación: http://%s:%s/test", cfg.Server.Host, cfg.Server.Port)
	log.Printf("🌐 Para verificar estado: http://%s:%s/status", cfg.Server.Host, cfg.Server.Port)

	if err := http.ListenAndServe(address, handler); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
