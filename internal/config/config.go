package config

import (
	"fmt"
	"net"
	"strconv"
)

// Config хранит параметры запуска SOCKS5-сервера
type Config struct {
	IP       string
	Port     int
	Username string
	Password string
}

// Validate проверяет корректность конфигурации
func (c *Config) Validate() error {
	// Проверяем IP
	if c.IP != "" && net.ParseIP(c.IP) == nil && c.IP != "localhost" {
		return fmt.Errorf("некорректный IP адрес: %s", c.IP)
	}

	// Проверяем порт
	if c.Port < 1 || c.Port > 65535 {
		return fmt.Errorf("порт должен быть в диапазоне 1-65535, получен: %d", c.Port)
	}

	// Проверяем учетные данные
	if c.Username == "" {
		return fmt.Errorf("имя пользователя не может быть пустым")
	}
	if c.Password == "" {
		return fmt.Errorf("пароль не может быть пустым")
	}

	return nil
}

// Addr возвращает адрес в формате host:port
func (c *Config) Addr() string {
	return net.JoinHostPort(c.IP, strconv.Itoa(c.Port))
}
