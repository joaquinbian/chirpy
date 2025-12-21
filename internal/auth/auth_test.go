package auth

import "testing"

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
