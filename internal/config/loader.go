package config

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// readPropertiesFile читает файл конфигурации формата key=value.
// Строки, начинающиеся с '#', и пустые строки игнорируются.
func readPropertiesFile(filename string) (map[string]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("не удалось открыть файл %s: %w", filename, err)
	}
	defer file.Close()

	// Оцениваем примерный размер для оптимизации выделения памяти
	fileInfo, _ := file.Stat()
	estimatedSize := 0
	if fileInfo != nil {
		estimatedSize = int(fileInfo.Size() / 50) // Примерно 50 байт на строку
	}
	properties := make(map[string]string, estimatedSize)

	scanner := bufio.NewScanner(file)
	lineNum := 0
	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())

		// Пропускаем пустые строки и комментарии
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Разделяем ключ и значение
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			log.Printf("Предупреждение: строка %d игнорируется (неверный формат): %s", lineNum, line)
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		// Проверяем, что ключ не пустой
		if key == "" {
			log.Printf("Предупреждение: строка %d игнорируется (пустой ключ): %s", lineNum, line)
			continue
		}

		properties[key] = value
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("ошибка чтения файла %s: %w", filename, err)
	}

	return properties, nil
}

// LoadConfig загружает и валидирует конфигурацию из файла
func LoadConfig(configPath string) (*Config, error) {
	props, err := readPropertiesFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("ошибка чтения конфигурации: %w", err)
	}

	// Функция-помощник для получения обязательных полей
	getRequired := func(key string) (string, error) {
		value, ok := props[key]
		if !ok {
			return "", fmt.Errorf("отсутствует обязательный ключ %q", key)
		}
		return value, nil
	}

	// Получаем обязательные поля
	ip, err := getRequired("ip")
	if err != nil {
		return nil, err
	}

	portStr, err := getRequired("port")
	if err != nil {
		return nil, err
	}

	port, err := strconv.Atoi(portStr)
	if err != nil {
		return nil, fmt.Errorf("порт должен быть числом: %w", err)
	}

	username, err := getRequired("username")
	if err != nil {
		return nil, err
	}

	password, err := getRequired("password")
	if err != nil {
		return nil, err
	}

	cfg := &Config{
		IP:       ip,
		Port:     port,
		Username: username,
		Password: password,
	}

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("неверная конфигурация: %w", err)
	}

	return cfg, nil
}

// GetConfigPath определяет путь к файлу конфигурации
func GetConfigPath(configPathFlag string) string {
	// Приоритет 1: флаг командной строки
	if configPathFlag != "" {
		return configPathFlag
	}

	// Приоритет 2: переменная окружения
	if configPath := os.Getenv("SOCKS5_CONFIG"); configPath != "" {
		return configPath
	}

	// Приоритет 3: текущая директория
	if _, err := os.Stat(".env"); err == nil {
		return ".env"
	}

	// Приоритет 4: директория исполняемого файла
	ex, err := os.Executable()
	if err == nil {
		return filepath.Join(filepath.Dir(ex), ".env")
	}

	// Последняя надежда - текущая директория
	return ".env"
}
