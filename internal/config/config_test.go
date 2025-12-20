package config

import (
	"os"
	"testing"
)

func TestValidate(t *testing.T) {
	tests := []struct {
		name    string
		cfg     *Config
		wantErr bool
	}{
		{
			name: "Valid config",
			cfg: &Config{
				WebhookSecret: "secret123",
			},
			wantErr: false,
		},
		{
			name: "Missing secret",
			cfg: &Config{
				WebhookSecret: "",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Validate(tt.cfg); (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetEnvOrDefault(t *testing.T) {
	key := "TRIGRA_TEST_KEY"
	def := "default_value"

	// Test default
	os.Unsetenv(key)
	if val := getEnvOrDefault(key, def); val != def {
		t.Errorf("getEnvOrDefault() = %v, want %v", val, def)
	}

	// Test env value
	expected := "env_value"
	os.Setenv(key, expected)
	defer os.Unsetenv(key)
	if val := getEnvOrDefault(key, def); val != expected {
		t.Errorf("getEnvOrDefault() = %v, want %v", val, expected)
	}
}
