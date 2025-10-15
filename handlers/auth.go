package handlers

import (
	"encoding/json"
	"scp-go-identity-service/models"
	"scp-go-identity-service/services"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
)

// AuthHandler handles authentication requests
type AuthHandler struct {
	dbService       *services.DatabaseService
	passwordService *services.PasswordService
	jwtService      *services.JWTService
}

// NewAuthHandler creates a new authentication handler
func NewAuthHandler(dbService *services.DatabaseService, passwordService *services.PasswordService, jwtService *services.JWTService) *AuthHandler {
	return &AuthHandler{
		dbService:       dbService,
		passwordService: passwordService,
		jwtService:      jwtService,
	}
}

// Authenticate handles authentication requests
func (ah *AuthHandler) Authenticate(w http.ResponseWriter, r *http.Request) {
	requestID := uuid.New().String()[:8]

	// Log request details
	clientIP := r.RemoteAddr
	userAgent := r.Header.Get("User-Agent")
	contentType := r.Header.Get("Content-Type")

	log.Printf("[%s] Authentication request from IP: %s, UserAgent: %s, ContentType: %s",
		requestID, clientIP, userAgent, contentType)

	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		ah.sendErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed", requestID)
		return
	}

	var authReq models.AuthRequest
	if err := json.NewDecoder(r.Body).Decode(&authReq); err != nil {
		log.Printf("[%s] Failed to decode request body: %v", requestID, err)
		ah.sendErrorResponse(w, http.StatusBadRequest, "Invalid request format", requestID)
		return
	}

	// Validate required fields
	if authReq.Email == "" || authReq.Password == "" {
		log.Printf("[%s] Missing required fields: email or password", requestID)
		ah.sendErrorResponse(w, http.StatusBadRequest, "Email and password are required", requestID)
		return
	}

	// Get user from database
	user, err := ah.dbService.GetUserByEmail(authReq.Email)
	if err != nil {
		log.Printf("[%s] ❌ Database error while finding user %s: %v", requestID, authReq.Email, err)
		log.Printf("[%s] Error type: %T", requestID, err)
		ah.sendErrorResponse(w, http.StatusInternalServerError, "Database connection error", requestID)
		return
	}

	if user == nil {
		log.Printf("[%s] ❌ User not found: %s", requestID, authReq.Email)
		// To prevent user enumeration, we perform a dummy password check.
		// The cost of this check should be tuned to be similar to a real check.
		_, _ = ah.passwordService.VerifyPassword(authReq.Password, "dummyhashfortimingattackprevention")
		ah.sendErrorResponse(w, http.StatusUnauthorized, "Invalid credentials", requestID)
		return
	}

	// Check if user is locked out
	if user.LockoutEnabled && user.LockoutEnd.Valid && user.LockoutEnd.Time.After(time.Now()) {
		log.Printf("[%s] ❌ User account is locked out: %s", requestID, authReq.Email)
		ah.sendErrorResponse(w, http.StatusUnauthorized, "Account is locked", requestID)
		return
	}

	// Verify password
	if !user.PasswordHash.Valid || user.PasswordHash.String == "" {
		log.Printf("[%s] ❌ User has no password hash: %s", requestID, authReq.Email)
		ah.sendErrorResponse(w, http.StatusUnauthorized, "Invalid credentials", requestID)
		return
	}

	isValidPassword, err := ah.passwordService.VerifyPassword(authReq.Password, user.PasswordHash.String)
	if err != nil {
		log.Printf("[%s] Error verifying password for user %s: %v", requestID, authReq.Email, err)
		ah.sendErrorResponse(w, http.StatusInternalServerError, "Internal server error", requestID)
		return
	}

	if !isValidPassword {
		log.Printf("[%s] ❌ Invalid password for user: %s", requestID, authReq.Email)
		ah.sendErrorResponse(w, http.StatusUnauthorized, "Invalid credentials", requestID)
		return
	}

	log.Printf("[%s] ✅ Password verification successful for user: %s", requestID, authReq.Email)

	// Get user roles
	roles, err := ah.dbService.GetUserRoles(user.Id)
	if err != nil {
		log.Printf("[%s] Error getting roles for user %s: %v", requestID, user.Id, err)
		// Continue without roles rather than failing
		roles = []string{}
	}

	// Get user claims
	userClaims, err := ah.dbService.GetUserClaims(user.Id)
	if err != nil {
		log.Printf("[%s] Error getting claims for user %s: %v", requestID, user.Id, err)
		// Continue without claims rather than failing
		userClaims = []models.AspNetUserClaim{}
	}

	// Convert claims to response format
	claimsInfo := make([]models.ClaimInfo, len(userClaims))
	for i, claim := range userClaims {
		claimsInfo[i] = models.ClaimInfo{
			Type:  claim.ClaimType.String,
			Value: claim.ClaimValue.String,
		}
	}

	// Build user info
	userInfo := &models.UserInfo{
		Id:               user.Id,
		Email:            user.Email.String,
		UserName:         user.UserName.String,
		FirstName:        user.FirstName.String,
		LastName:         user.LastName.String,
		Culture:          user.Culture.String,
		IsEmailConfirmed: user.EmailConfirmed,
		PhoneNumber:      user.PhoneNumber.String,
	}

	// Generate JWT token if requested
	var token string
	var expiresAt *string
	if authReq.GenerateToken {
		log.Printf("[%s] Generating JWT token for user: %s", requestID, user.Id)

		token, err = ah.jwtService.GenerateToken(
			user.Id,
			user.Email.String,
			user.UserName.String,
			roles,
			user.FirstName.String,
			user.LastName.String,
			user.Culture.String,
		)
		if err != nil {
			log.Printf("[%s] Error generating token for user %s: %v", requestID, user.Id, err)
			ah.sendErrorResponse(w, http.StatusInternalServerError, "Failed to generate token", requestID)
			return
		}

		if token != "" {
			expirationTime := time.Now().Add(12 * time.Hour).UTC().Format(time.RFC3339)
			expiresAt = &expirationTime
			tokenPreview := token
			if len(token) > 40 {
				tokenPreview = token[:20] + "..." + token[len(token)-20:]
			}
			log.Printf("[%s] Token generated successfully. Length: %d, Preview: %s", requestID, len(token), tokenPreview)
		} else {
			log.Printf("[%s] ⚠ Generated token is empty", requestID)
		}
	} else {
		log.Printf("[%s] Token generation not requested", requestID)
	}

	// Build response
	response := models.AuthResponse{
		Success:   true,
		Message:   "Authentication successful",
		User:      userInfo,
		Roles:     roles,
		Claims:    claimsInfo,
		Token:     token,
		ExpiresAt: expiresAt,
	}

	log.Printf("[%s] ✅ Authentication successful for user: %s, Roles: %v", requestID, user.Email.String, roles)

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("[%s] Error encoding response: %v", requestID, err)
	}
}

// sendErrorResponse sends an error response
func (ah *AuthHandler) sendErrorResponse(w http.ResponseWriter, statusCode int, message, requestID string) {
	log.Printf("[%s] Sending error response: %d - %s", requestID, statusCode, message)

	response := models.ErrorResponse{
		Success: false,
		Message: message,
	}

	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("[%s] Error encoding error response: %v", requestID, err)
	}
}
