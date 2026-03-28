package logger

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
	"syscall"

	"go-socks5-relay/internal/config"
)

// LogLevel константы для уровней логирования
const (
	LogLevelError = iota
	LogLevelWarn
	LogLevelInfo
	LogLevelDebug
)

// FilteredLogger оборачивает стандартный logger и подавляет шумные ошибки
type FilteredLogger struct {
	inner *log.Logger
	debug bool
	level int
}

// NewFilteredLogger создает новый логгер с фильтрацией шумных ошибок
func NewFilteredLogger(debug bool, logLevel string) *FilteredLogger {
	level := LogLevelInfo // по умолчанию

	switch strings.ToLower(logLevel) {
	case "error":
		level = LogLevelError
	case "warn", "warning":
		level = LogLevelWarn
	case "info":
		level = LogLevelInfo
	case "debug":
		level = LogLevelDebug
		debug = true // debug режим автоматически включает debug уровень
	}

	return &FilteredLogger{
		inner: log.New(os.Stdout, "[socks5] ", log.LstdFlags),
		debug: debug,
		level: level,
	}
}

// Errorf логирует ошибку, если она не является шумной или включен debug режим
func (l *FilteredLogger) Errorf(format string, args ...interface{}) {
	if l.level < LogLevelError {
		return
	}

	msg := fmt.Sprintf(format, args...)

	// В debug режиме логируем все ошибки
	if l.debug {
		l.inner.Printf("[D]: %s", msg)
		return
	}

	// Фильтруем шумные ошибки
	if isNoisyError(msg) {
		return
	}

	l.inner.Printf("[E]: %s", msg)
}

// Warnf логирует предупреждение
func (l *FilteredLogger) Warnf(format string, args ...interface{}) {
	if l.level >= LogLevelWarn {
		l.inner.Printf("[W]: "+format, args...)
	}
}

// Infof логирует информационное сообщение
func (l *FilteredLogger) Infof(format string, args ...interface{}) {
	if l.level >= LogLevelInfo {
		l.inner.Printf("[I]: "+format, args...)
	}
}

// Info логирует информационное сообщение (без форматирования)
func (l *FilteredLogger) Info(msg string) {
	if l.level >= LogLevelInfo {
		l.inner.Printf("[I]: %s", msg)
	}
}

// Debugf логирует отладочное сообщение
func (l *FilteredLogger) Debugf(format string, args ...interface{}) {
	if l.level >= LogLevelDebug {
		l.inner.Printf("[D]: "+format, args...)
	}
}

// PrintStartupInfo выводит информацию о запуске
func (l *FilteredLogger) PrintStartupInfo(cfg *config.Config, configPath, logLevel string, debugMode bool) {
	l.Info("=== SOCKS5 Proxy Server ===")
	l.Infof("Адрес: %s", cfg.Addr())
	l.Infof("Аутентификация: Username/Password")
	l.Infof("Пользователь: %q", cfg.Username)
	l.Infof("Конфигурация: %s", configPath)
	l.Infof("Уровень логирования: %s", logLevel)
	l.Infof("Debug режим: %v", debugMode)
	l.Info("===========================")
}

// isNoisyError возвращает true для ошибок, которые являются штатным поведением
func isNoisyError(msg string) bool {
	noisyPatterns := []string{
		io.EOF.Error(),
		"broken pipe",
		"connection reset by peer",
		"use of closed network connection",
		syscall.EPIPE.Error(),
		net.ErrClosed.Error(),
		"i/o timeout",
		"context canceled",
	}

	for _, pattern := range noisyPatterns {
		if strings.Contains(msg, pattern) {
			return true
		}
	}
	return false
}
