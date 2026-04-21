package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestReadPropertiesFile(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "test.properties")

	content := `# Comment
key1=value1
key2 = value2
key3=value3 with spaces
  key4  =  value4

empty=
=invalid
`
	if err := os.WriteFile(configPath, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	props, err := readPropertiesFile(configPath)
	if err != nil {
		t.Fatalf("readPropertiesFile failed: %v", err)
	}

	expected := map[string]string{
		"key1":  "value1",
		"key2":  "value2",
		"key3":  "value3 with spaces",
		"key4":  "value4",
		"empty": "",
	}

	for k, v := range expected {
		if props[k] != v {
			t.Errorf("key %q: expected %q, got %q", k, v, props[k])
		}
	}

	if _, ok := props["invalid"]; ok {
		t.Error("line with empty key should be ignored")
	}
}

func TestReadPropertiesFileNotFound(t *testing.T) {
	_, err := readPropertiesFile("/nonexistent/path/config.env")
	if err == nil {
		t.Error("expected error for missing file, got nil")
	}
}

func TestConfigValidate(t *testing.T) {
	tests := []struct {
		name    string
		cfg     Config
		wantErr bool
	}{
		{
			name:    "valid config",
			cfg:     Config{IP: "0.0.0.0", Port: 1080, Username: "user", Password: "pass"},
			wantErr: false,
		},
		{
			name:    "valid localhost",
			cfg:     Config{IP: "localhost", Port: 1080, Username: "user", Password: "pass"},
			wantErr: false,
		},
		{
			name:    "port too low",
			cfg:     Config{IP: "0.0.0.0", Port: 0, Username: "user", Password: "pass"},
			wantErr: true,
		},
		{
			name:    "port too high",
			cfg:     Config{IP: "0.0.0.0", Port: 65536, Username: "user", Password: "pass"},
			wantErr: true,
		},
		{
			name:    "empty username",
			cfg:     Config{IP: "0.0.0.0", Port: 1080, Username: "", Password: "pass"},
			wantErr: true,
		},
		{
			name:    "empty password",
			cfg:     Config{IP: "0.0.0.0", Port: 1080, Username: "user", Password: ""},
			wantErr: true,
		},
		{
			name:    "invalid IP",
			cfg:     Config{IP: "invalid", Port: 1080, Username: "user", Password: "pass"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.cfg.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestConfigAddr(t *testing.T) {
	cfg := Config{IP: "127.0.0.1", Port: 1080}
	if got := cfg.Addr(); got != "127.0.0.1:1080" {
		t.Errorf("Addr() = %q, want %q", got, "127.0.0.1:1080")
	}
}

func TestGetConfigPath(t *testing.T) {
	origEnv := os.Getenv("SOCKS5_CONFIG")
	defer os.Setenv("SOCKS5_CONFIG", origEnv)

	// Priority 1: CLI flag overrides env var
	os.Setenv("SOCKS5_CONFIG", "/env/path")
	if got := GetConfigPath("/flag/path"); got != "/flag/path" {
		t.Errorf("flag priority: expected /flag/path, got %s", got)
	}

	// Priority 2: env var when no flag
	os.Setenv("SOCKS5_CONFIG", "/env/path")
	if got := GetConfigPath(""); got != "/env/path" {
		t.Errorf("env priority: expected /env/path, got %s", got)
	}

	// Priority 3/4: fallback returns non-empty path
	os.Setenv("SOCKS5_CONFIG", "")
	if got := GetConfigPath(""); got == "" {
		t.Error("fallback: expected non-empty path")
	}
}

func TestLoadConfig(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, ".env")

	content := "ip=127.0.0.1\nport=1080\nusername=user\npassword=pass\n"
	if err := os.WriteFile(configPath, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	cfg, err := LoadConfig(configPath)
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}

	if cfg.IP != "127.0.0.1" {
		t.Errorf("IP: expected 127.0.0.1, got %s", cfg.IP)
	}
	if cfg.Port != 1080 {
		t.Errorf("Port: expected 1080, got %d", cfg.Port)
	}
	if cfg.Username != "user" {
		t.Errorf("Username: expected user, got %s", cfg.Username)
	}
}

func TestLoadConfigMissingKey(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, ".env")

	// Missing password
	content := "ip=127.0.0.1\nport=1080\nusername=user\n"
	if err := os.WriteFile(configPath, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	if _, err := LoadConfig(configPath); err == nil {
		t.Error("expected error for missing password key, got nil")
	}
}
