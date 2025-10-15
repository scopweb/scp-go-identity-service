package services

import (
	"database/sql"
	"fmt"
	"scp-go-identity-service/models"
	"log"
	"strings"

	_ "github.com/denisenkom/go-mssqldb"
)

// DatabaseService handles database operations
type DatabaseService struct {
	db *sql.DB
}

// NewDatabaseService creates a new database service
func NewDatabaseService(connectionString string) (*DatabaseService, error) {
	db, err := sql.Open("sqlserver", connectionString)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %v", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %v", err)
	}

	return &DatabaseService{db: db}, nil
}

// Close closes the database connection
func (ds *DatabaseService) Close() error {
	return ds.db.Close()
}

// TestConnection tests the database connection
func (ds *DatabaseService) TestConnection() error {
	return ds.db.Ping()
}

// GetUserByEmail retrieves a user by email
func (ds *DatabaseService) GetUserByEmail(email string) (*models.AspNetUser, error) {
	normalizedEmail := strings.ToUpper(email)
	query := `
		SELECT Id, UserName, NormalizedUserName, Email, NormalizedEmail, 
		       EmailConfirmed, PasswordHash, SecurityStamp, ConcurrencyStamp,
		       PhoneNumber, PhoneNumberConfirmed, TwoFactorEnabled, LockoutEnd,
		       LockoutEnabled, AccessFailedCount, FirstName, LastName, Culture
		FROM AspNetUsers 
		WHERE NormalizedEmail = @email
	`

	var user models.AspNetUser
	err := ds.db.QueryRow(query, sql.Named("email", normalizedEmail)).Scan(
		&user.Id,
		&user.UserName,
		&user.NormalizedUserName,
		&user.Email,
		&user.NormalizedEmail,
		&user.EmailConfirmed,
		&user.PasswordHash,
		&user.SecurityStamp,
		&user.ConcurrencyStamp,
		&user.PhoneNumber,
		&user.PhoneNumberConfirmed,
		&user.TwoFactorEnabled,
		&user.LockoutEnd,
		&user.LockoutEnabled,
		&user.AccessFailedCount,
		&user.FirstName,
		&user.LastName,
		&user.Culture,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // User not found
		}
		return nil, fmt.Errorf("failed to get user by email: %v", err)
	}

	return &user, nil
}

// GetUserRoles retrieves roles for a user
func (ds *DatabaseService) GetUserRoles(userID string) ([]string, error) {
	query := `
		SELECT r.Name
		FROM AspNetRoles r
		INNER JOIN AspNetUserRoles ur ON r.Id = ur.RoleId
		WHERE ur.UserId = @userId AND r.Name IS NOT NULL
	`

	rows, err := ds.db.Query(query, sql.Named("userId", userID))
	if err != nil {
		return nil, fmt.Errorf("failed to get user roles: %v", err)
	}
	defer rows.Close()

	var roles []string
	for rows.Next() {
		var roleName sql.NullString
		if err := rows.Scan(&roleName); err != nil {
			log.Printf("Error scanning role: %v", err)
			continue
		}
		if roleName.Valid {
			roles = append(roles, roleName.String)
		}
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating roles: %v", err)
	}

	return roles, nil
}

// GetUserClaims retrieves claims for a user
func (ds *DatabaseService) GetUserClaims(userID string) ([]models.AspNetUserClaim, error) {
	query := `
		SELECT Id, UserId, ClaimType, ClaimValue
		FROM AspNetUserClaims
		WHERE UserId = @userId
	`

	rows, err := ds.db.Query(query, sql.Named("userId", userID))
	if err != nil {
		return nil, fmt.Errorf("failed to get user claims: %v", err)
	}
	defer rows.Close()

	var claims []models.AspNetUserClaim
	for rows.Next() {
		var claim models.AspNetUserClaim
		if err := rows.Scan(&claim.Id, &claim.UserId, &claim.ClaimType, &claim.ClaimValue); err != nil {
			log.Printf("Error scanning claim: %v", err)
			continue
		}
		claims = append(claims, claim)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating claims: %v", err)
	}

	return claims, nil
}
