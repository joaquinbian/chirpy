package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestHashPassword(t *testing.T) {
	tests := []struct {
		name     string
		password string
		wantErr  bool
	}{
		{
			name:     "hash valid password",
			password: "mypassword123",
			wantErr:  false,
		},
		{
			name:     "hash password with special characters",
			password: "P@ssw0rd!#$%",
			wantErr:  false,
		},
		{
			name:     "hash empty password",
			password: "",
			wantErr:  false,
		},
		{
			name:     "hash long password",
			password: "thisIsAVeryLongPasswordWithLotsOfCharactersAndNumbers12345678901234567890",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hash, err := HashPassword(tt.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("HashPassword() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil && hash == "" {
				t.Errorf("HashPassword() returned empty hash")
			}
		})
	}
}

func TestComparePasswordHash(t *testing.T) {
	// Create a hash for testing
	password := "testpassword123"
	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("Failed to create hash for testing: %v", err)
	}

	tests := []struct {
		name     string
		password string
		hash     string
		want     bool
		wantErr  bool
	}{
		{
			name:     "correct password matches hash",
			password: password,
			hash:     hash,
			want:     true,
			wantErr:  false,
		},
		{
			name:     "incorrect password does not match",
			password: "wrongpassword",
			hash:     hash,
			want:     false,
			wantErr:  false,
		},
		{
			name:     "invalid hash returns false",
			password: password,
			hash:     "invalid_hash_string",
			want:     false,
			wantErr:  false,
		},
		{
			name:     "empty password with valid hash",
			password: "",
			hash:     hash,
			want:     false,
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ComparePasswordHash(tt.password, tt.hash)
			if (err != nil) != tt.wantErr {
				t.Errorf("ComparePasswordHash() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ComparePasswordHash() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHashPasswordConsistency(t *testing.T) {
	password := "consistent_test_password"

	hash1, err1 := HashPassword(password)
	if err1 != nil {
		t.Fatalf("First HashPassword failed: %v", err1)
	}

	hash2, err2 := HashPassword(password)
	if err2 != nil {
		t.Fatalf("Second HashPassword failed: %v", err2)
	}

	// Hashes should be different due to salt randomization
	if hash1 == hash2 {
		t.Errorf("HashPassword produced identical hashes for same password (should differ due to salt)")
	}

	// But both should match the same password
	ok1, _ := ComparePasswordHash(password, hash1)
	ok2, _ := ComparePasswordHash(password, hash2)

	if !ok1 || !ok2 {
		t.Errorf("Different hashes for same password should both validate correctly")
	}
}

func TestMakeJWT(t *testing.T) {
	userID := uuid.New()
	tokenSecret := "my-secret-key-for-testing"
	expiresIn := 1 * time.Hour

	tests := []struct {
		name      string
		userID    uuid.UUID
		secret    string
		expiresIn time.Duration
		wantErr   bool
	}{
		{
			name:      "create valid JWT",
			userID:    userID,
			secret:    tokenSecret,
			expiresIn: expiresIn,
			wantErr:   false,
		},
		{
			name:      "create JWT with different user ID",
			userID:    uuid.New(),
			secret:    tokenSecret,
			expiresIn: expiresIn,
			wantErr:   false,
		},
		{
			name:      "create JWT with short expiration",
			userID:    userID,
			secret:    tokenSecret,
			expiresIn: 1 * time.Minute,
			wantErr:   false,
		},
		{
			name:      "create JWT with long expiration",
			userID:    userID,
			secret:    tokenSecret,
			expiresIn: 24 * time.Hour,
			wantErr:   false,
		},
		{
			name:      "create JWT with empty secret",
			userID:    userID,
			secret:    "",
			expiresIn: expiresIn,
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := MakeJWT(tt.userID, tt.secret, tt.expiresIn)
			if (err != nil) != tt.wantErr {
				t.Errorf("MakeJWT() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil && token == "" {
				t.Errorf("MakeJWT() returned empty token")
			}
		})
	}
}

func TestValidateToken(t *testing.T) {
	userID := uuid.New()
	tokenSecret := "my-secret-key-for-testing"
	expiresIn := 1 * time.Hour

	// Create a valid token
	validToken, err := MakeJWT(userID, tokenSecret, expiresIn)
	if err != nil {
		t.Fatalf("Failed to create valid token: %v", err)
	}

	// Create an expired token
	expiredToken, err := MakeJWT(userID, tokenSecret, -1*time.Hour)
	if err != nil {
		t.Fatalf("Failed to create expired token: %v", err)
	}

	tests := []struct {
		name      string
		token     string
		secret    string
		wantID    uuid.UUID
		wantError bool
	}{
		{
			name:      "validate valid token",
			token:     validToken,
			secret:    tokenSecret,
			wantID:    userID,
			wantError: false,
		},
		{
			name:      "validate token with wrong secret",
			token:     validToken,
			secret:    "wrong-secret",
			wantID:    uuid.Nil,
			wantError: true,
		},
		{
			name:      "validate expired token",
			token:     expiredToken,
			secret:    tokenSecret,
			wantID:    uuid.Nil,
			wantError: true,
		},
		{
			name:      "validate malformed token",
			token:     "invalid.token.string",
			secret:    tokenSecret,
			wantID:    uuid.Nil,
			wantError: true,
		},
		{
			name:      "validate empty token",
			token:     "",
			secret:    tokenSecret,
			wantID:    uuid.Nil,
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id, err := ValidateToken(tt.token, tt.secret)
			if (err != nil) != tt.wantError {
				t.Errorf("ValidateToken() error = %v, wantError %v", err, tt.wantError)
				return
			}
			if !tt.wantError && id != tt.wantID {
				t.Errorf("ValidateToken() returned ID = %v, want %v", id, tt.wantID)
			}
			if tt.wantError && id != uuid.Nil {
				t.Errorf("ValidateToken() returned ID = %v, want uuid.Nil on error", id)
			}
		})
	}
}

func TestMakeAndValidateJWT(t *testing.T) {
	userID := uuid.New()
	tokenSecret := "integration-test-secret"
	expiresIn := 2 * time.Hour

	// Create a token
	token, err := MakeJWT(userID, tokenSecret, expiresIn)
	if err != nil {
		t.Fatalf("MakeJWT failed: %v", err)
	}

	// Validate the token
	validatedID, err := ValidateToken(token, tokenSecret)
	if err != nil {
		t.Fatalf("ValidateToken failed: %v", err)
	}

	// Check that the ID matches
	if validatedID != userID {
		t.Errorf("ValidateToken returned ID %v, want %v", validatedID, userID)
	}
}

func TestValidateTokenDifferentSecrets(t *testing.T) {
	userID := uuid.New()
	secret1 := "secret-one"
	secret2 := "secret-two"
	expiresIn := 1 * time.Hour

	// Create token with secret1
	token, err := MakeJWT(userID, secret1, expiresIn)
	if err != nil {
		t.Fatalf("MakeJWT failed: %v", err)
	}

	// Try to validate with secret2 - should fail
	_, err = ValidateToken(token, secret2)
	if err == nil {
		t.Errorf("ValidateToken should fail when using different secret")
	}

	// Validate with correct secret1 - should succeed
	validatedID, err := ValidateToken(token, secret1)
	if err != nil {
		t.Errorf("ValidateToken failed with correct secret: %v", err)
	}

	if validatedID != userID {
		t.Errorf("ValidateToken returned ID %v, want %v", validatedID, userID)
	}
}
