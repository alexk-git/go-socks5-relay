package main

import (
	"context"
	"net"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestReadPropertiesFile(t *testing.T) {
	// Создаем временный файл конфигурации
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
		"key1": "value1",
		"key2": "value2",
		"key3": "value3 with spaces",
		"key4": "value4",
		"empty": "",
	}
	
	for k, v := range expected {
		if props[k] != v {
			t.Errorf("Key %q: expected %q, got %q", k, v, props[k])
		}
	}
	
	// Проверяем, что невалидные строки игнорируются
	if _, ok := props["invalid"]; ok {
		t.Error("Invalid line should be ignored")
	}
}

func TestConfigValidation(t *testing.T) {
	tests := []struct {
		name    string
		config  Config
		wantErr bool
	}{
		{
			name: "valid config",
			config: Config{
				IP:       "0.0.0.0",
				Port:     1080,
				Username: "user",
				Password: "pass",
			},
			wantErr: false,
		},
		{
			name: "invalid port",
			config: Config{
				IP:       "0.0.0.0",
				Port:     0,
				Username: "user",
				Password: "pass",
			},
			wantErr: true,
		},
		{
			name: "empty username",
			config: Config{
				IP:       "0.0.0.0",
				Port:     1080,
				Username: "",
				Password: "pass",
			},
			wantErr: true,
		},
		{
			name: "invalid IP",
			config: Config{
				IP:       "invalid",
				Port:     1080,
				Username: "user",
				Password: "pass",
			},
			wantErr: true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestIsNoisyError(t *testing.T) {
	tests := []struct {
		msg      string
		expected bool
	}{
		{"EOF", true},
		{"write: broken pipe", true},
		{"connection reset by peer", true},
		{"use of closed network connection", true},
		{"i/o timeout", true},
		{"context canceled", true},
		{"real error: something went wrong", false},
	}
	
	for _, tt := range tests {
		if result := isNoisyError(tt.msg); result != tt.expected {
			t.Errorf("isNoisyError(%q) = %v, expected %v", tt.msg, result, tt.expected)
		}
	}
}

func TestGetConfigPath(t *testing.T) {
	// Сохраняем оригинальные значения переменных окружения
	originalEnv := os.Getenv("SOCKS5_CONFIG")
	defer os.Setenv("SOCKS5_CONFIG", originalEnv)
	
	// Тест 1: через переменную окружения
	os.Setenv("SOCKS5_CONFIG", "/custom/path/config.properties")
	if path := getConfigPath(); path != "/custom/path/config.properties" {
		t.Errorf("Expected /custom/path/config.properties, got %s", path)
	}
	
	// Тест 2: без переменной окружения, проверяем существование файла в текущей директории
	os.Setenv("SOCKS5_CONFIG", "")
	// Создаем временный файл в текущей директории (осторожно!)
	// В реальном тесте лучше использовать временную директорию
}
