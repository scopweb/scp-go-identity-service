package services

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/binary"
	"errors"
	"fmt"
	"hash"

	"golang.org/x/crypto/pbkdf2"
)

// PasswordService handles .NET Identity password verification for V2 and V3 formats.
type PasswordService struct{}

// NewPasswordService creates a new password service
func NewPasswordService() *PasswordService {
	return &PasswordService{}
}

// VerifyPassword checks a password against a .NET Identity hash, supporting both V2 and V3 formats.
func (ps *PasswordService) VerifyPassword(password, hashedPassword string) (bool, error) {
	if hashedPassword == "" {
		return false, errors.New("password hash is empty")
	}

	hashBytes, err := base64.StdEncoding.DecodeString(hashedPassword)
	if err != nil {
		return false, fmt.Errorf("failed to decode password hash: %v", err)
	}

	if len(hashBytes) == 0 {
		return false, errors.New("decoded hash is empty")
	}

	// The first byte is the format marker.
	formatMarker := hashBytes[0]
	switch formatMarker {
	case 0x00:
		// This is the ASP.NET Identity V2 format.
		return verifyV2(password, hashBytes)
	case 0x01:
		// This is the ASP.NET Identity V3 format.
		return verifyV3(password, hashBytes)
	default:
		return false, fmt.Errorf("unsupported password hash format marker: %d", formatMarker)
	}
}

// verifyV2 handles the V2 format: [0x00] || [salt (16 bytes)] || [hash (20 bytes)]
// The hash is PBKDF2 with HMAC-SHA1, 1000 iterations.
func verifyV2(password string, hashBytes []byte) (bool, error) {
	// V2 format: 1 (marker) + 16 (salt) + 20 (hash) = 37 bytes
	if len(hashBytes) != 37 {
		return false, fmt.Errorf("invalid V2 hash format: expected 37 bytes, got %d", len(hashBytes))
	}

	salt := hashBytes[1:17]
	storedHash := hashBytes[17:]

	// V2 uses 1000 iterations and SHA1
	iterations := 1000
	computedHash := pbkdf2.Key([]byte(password), salt, iterations, 20, sha1.New)

	return hmac.Equal(storedHash, computedHash), nil
}

// verifyV3 handles the V3 format. This version is specifically tailored to the discovery
// that the .NET 9 implementation is using HMAC-SHA512 but requesting a 32-byte subkey.
func verifyV3(password string, hashBytes []byte) (bool, error) {
	headerSize := 13
	if len(hashBytes) < headerSize {
		return false, fmt.Errorf("invalid V3 hash format: insufficient length for header, got %d", len(hashBytes))
	}

	// Read header fields
	prfIdentifier := binary.BigEndian.Uint32(hashBytes[1:5])
	iterations := int(binary.BigEndian.Uint32(hashBytes[5:9]))
	saltSize := int(binary.BigEndian.Uint32(hashBytes[9:13]))

	// Plausibility check for salt size
	if saltSize < 0 || headerSize+saltSize > len(hashBytes) {
		return false, fmt.Errorf("invalid V3 hash format: invalid salt size %d", saltSize)
	}

	// Determine PRF (hash function) from the identifier
	var prf func() hash.Hash
	switch prfIdentifier {
	case 0: // HMAC-SHA1
		prf = sha1.New
	case 1: // HMAC-SHA256
		prf = sha256.New
	case 2: // HMAC-SHA512
		prf = sha512.New
	default:
		return false, fmt.Errorf("unsupported V3 PRF identifier: %d", prfIdentifier)
	}

	// Extract salt and the stored subkey
	salt := hashBytes[headerSize : headerSize+saltSize]
	storedSubkey := hashBytes[headerSize+saltSize:]

	// --- The Crucial Fix ---
	// Based on the analysis, .NET is using HMAC-SHA512 but requesting a 32-byte (256-bit) subkey.
	// We will use the actual length of the stored subkey for the comparison.
	keyLength := len(storedSubkey)

	// As a final sanity check, if the PRF is SHA512, we expect a 32-byte key based on the user's .NET environment.
	// If it's different, the hash is likely from another source, but we proceed with the actual length.
	if prfIdentifier == 2 && keyLength != 32 {
		// This is unexpected given the latest findings, but we will still attempt verification.
		// A log here could be useful in a real-world scenario.
	}

	// If for some reason the stored subkey is empty, we can't proceed.
	if keyLength == 0 {
		return false, errors.New("invalid V3 hash: subkey is empty")
	}

	// Compute the hash of the provided password using the determined parameters.
	computedSubkey := pbkdf2.Key([]byte(password), salt, iterations, keyLength, prf)

	// Compare the computed hash with the stored one.
	return hmac.Equal(storedSubkey, computedSubkey), nil
}

// HashPassword creates a .NET Identity V3 compatible password hash using HMAC-SHA512.
// This is updated to match the user's provided hash parameters.
func (ps *PasswordService) HashPassword(password string) (string, error) {
	// Generate 16-byte salt (128 bits)
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return "", fmt.Errorf("failed to generate salt: %v", err)
	}

	iterations := 100000       // Match user's hash
	prfIdentifier := uint32(2) // HMAC-SHA512
	keyLength := 64

	hash := pbkdf2.Key([]byte(password), salt, iterations, keyLength, sha512.New)

	// Build the result byte array according to the V3 format.
	// Total length = 1 (marker) + 4 (prf) + 4 (iter) + 4 (salt len) + len(salt) + len(hash)
	result := make([]byte, 1+4+4+4+len(salt)+len(hash))

	result[0] = 0x01                                            // Format marker
	binary.BigEndian.PutUint32(result[1:5], prfIdentifier)      // PRF identifier
	binary.BigEndian.PutUint32(result[5:9], uint32(iterations)) // Iteration count
	binary.BigEndian.PutUint32(result[9:13], uint32(len(salt))) // Salt size

	copy(result[13:13+len(salt)], salt)
	copy(result[13+len(salt):], hash)

	// Encode the final byte array to base64
	return base64.StdEncoding.EncodeToString(result), nil
}
