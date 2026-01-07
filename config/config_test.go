// config/config_test.go
package config

import (
	"os"
	"testing"
)

func TestGetEnv(t *testing.T) {
	tests := []struct {
		name     string
		key      string
		value    string
		fallback string
		want     string
	}{
		{
			name:     "env var exists",
			key:      "TEST_KEY",
			value:    "test_value",
			fallback: "default_value",
			want:     "test_value",
		},
		{
			name:     "env var does not exist",
			key:      "NON_EXISTENT_KEY",
			value:    "",
			fallback: "default_value",
			want:     "default_value",
		},
		{
			name:     "env var empty",
			key:      "EMPTY_KEY",
			value:    "",
			fallback: "default_value",
			want:     "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set environment variable jika ada value
			if tt.value != "" || tt.name == "env var empty" {
				os.Setenv(tt.key, tt.value)
				defer os.Unsetenv(tt.key)
			} else {
				// Pastikan env var tidak ada
				os.Unsetenv(tt.key)
			}

			got := getEnv(tt.key, tt.fallback)
			if got != tt.want {
				t.Errorf("getEnv() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLoadEnv(t *testing.T) {
	tests := []struct {
		name       string
		envVars    map[string]string
		wantPort   string
		wantDBHost string
		shouldLoad bool
	}{
		{
			name: "load with custom env vars",
			envVars: map[string]string{
				"PORT":                  "8080",
				"DB_HOST":               "test-host",
				"DB_USER":               "test-user",
				"DB_PASSWORD":           "test-pass",
				"JWT_SECRET":            "test-secret",
				"JWT_EXPIRY":            "120",
				"REFRESH_TOKEN_EXPIRED": "48h",
			},
			wantPort:   "8080",
			wantDBHost: "test-host",
			shouldLoad: true,
		},
		{
			name:       "load with default values",
			envVars:    map[string]string{},
			wantPort:   "3030",      // default value
			wantDBHost: "localhost", // default value
			shouldLoad: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup environment variables
			for key, value := range tt.envVars {
				os.Setenv(key, value)
				defer os.Unsetenv(key)
			}

			// Clear previous config
			AppConfig = nil

			// Load environment
			LoadEnv()

			// Verify config is not nil
			if AppConfig == nil {
				t.Fatal("AppConfig should not be nil after LoadEnv()")
			}

			// Test specific fields
			if AppConfig.AppPort != tt.wantPort {
				t.Errorf("AppConfig.AppPort = %v, want %v", AppConfig.AppPort, tt.wantPort)
			}

			if AppConfig.DBHost != tt.wantDBHost {
				t.Errorf("AppConfig.DBHost = %v, want %v", AppConfig.DBHost, tt.wantDBHost)
			}

			// Test that all fields are populated
			if AppConfig.DBUser == "" {
				t.Error("DBUser should not be empty")
			}

			if AppConfig.DBPassword == "" {
				t.Error("DBPassword should not be empty")
			}

			if AppConfig.JWTSecret == "" {
				t.Error("JWTSecret should not be empty")
			}

			if AppConfig.JWTExpire == "" {
				t.Error("JWTExpire should not be empty")
			}

			if AppConfig.JWTRefreshToken == "" {
				t.Error("JWTRefreshToken should not be empty")
			}
		})
	}
}

func TestLoadEnvWithEnvFile(t *testing.T) {
	// Create a temporary .env file
	envContent := `PORT=9999
DB_HOST=test-db-host
DB_USER=test-db-user
DB_PASSWORD=test-db-pass
JWT_SECRET=test-jwt-secret
JWT_EXPIRY=3600
REFRESH_TOKEN_EXPIRED=72h`

	tmpFile, err := os.CreateTemp("", "test*.env")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.WriteString(envContent); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tmpFile.Close()

	// Change to temp directory and back
	oldWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}
	defer os.Chdir(oldWd)

	tmpDir := tmpFile.Name()
	os.Chdir(tmpDir)

	// Clear previous config
	AppConfig = nil

	// Load environment
	LoadEnv()

	// Verify values from .env file
	if AppConfig.AppPort != "9999" {
		t.Errorf("AppPort = %v, want 9999", AppConfig.AppPort)
	}

	if AppConfig.DBHost != "test-db-host" {
		t.Errorf("DBHost = %v, want test-db-host", AppConfig.DBHost)
	}

	if AppConfig.DBUser != "test-db-user" {
		t.Errorf("DBUser = %v, want test-db-user", AppConfig.DBUser)
	}

	if AppConfig.DBPassword != "test-db-pass" {
		t.Errorf("DBPassword = %v, want test-db-pass", AppConfig.DBPassword)
	}

	if AppConfig.JWTSecret != "test-jwt-secret" {
		t.Errorf("JWTSecret = %v, want test-jwt-secret", AppConfig.JWTSecret)
	}

	if AppConfig.JWTExpire != "3600" {
		t.Errorf("JWTExpire = %v, want 3600", AppConfig.JWTExpire)
	}

	if AppConfig.JWTRefreshToken != "72h" {
		t.Errorf("JWTRefreshToken = %v, want 72h", AppConfig.JWTRefreshToken)
	}
}

func TestConfigSingleton(t *testing.T) {
	// Clear config
	AppConfig = nil

	// First load
	LoadEnv()
	config1 := AppConfig

	// Second load should return same instance
	LoadEnv()
	config2 := AppConfig

	if config1 != config2 {
		t.Error("LoadEnv() should maintain singleton instance")
	}
}
