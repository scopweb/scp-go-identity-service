package models

// AuthRequest represents the authentication request
type AuthRequest struct {
	Email         string `json:"email" binding:"required,email"`
	Password      string `json:"password" binding:"required"`
	GenerateToken bool   `json:"generateToken"`
}

// AuthResponse represents the authentication response
type AuthResponse struct {
	Success   bool        `json:"success"`
	Message   string      `json:"message"`
	User      *UserInfo   `json:"user,omitempty"`
	Roles     []string    `json:"roles,omitempty"`
	Claims    []ClaimInfo `json:"claims,omitempty"`
	Token     string      `json:"token,omitempty"`
	ExpiresAt *string     `json:"expiresAt,omitempty"`
}

// UserInfo represents user information in the response
type UserInfo struct {
	Id               string `json:"id"`
	Email            string `json:"email"`
	UserName         string `json:"userName"`
	FirstName        string `json:"firstName"`
	LastName         string `json:"lastName"`
	Culture          string `json:"culture"`
	IsEmailConfirmed bool   `json:"isEmailConfirmed"`
	PhoneNumber      string `json:"phoneNumber"`
}

// ClaimInfo represents claim information in the response
type ClaimInfo struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

// ErrorResponse represents error response
type ErrorResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}
