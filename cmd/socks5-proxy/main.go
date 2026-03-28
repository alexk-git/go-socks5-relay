package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"go-socks5-relay/internal/config"
	"go-socks5-relay/internal/logger"
	"go-socks5-relay/internal/proxy"
)

func printUsage() {
	fmt.Printf(`SOCKS5 Proxy Server

Usage: %s [options]

Options:
  -debug, -d              Enable debug logging (same as -log-level debug)
  -config <file>          Path to configuration file (default: .env)
  -port <port>            Override port from config file
  -ip <ip>                Override IP address from config file
  -log-level <level>      Set log level: error, warn, info, debug (default: info)
  -help, -h               Show this help message

Environment variables:
  SOCKS5_DEBUG            Set to 1 or true to enable debug mode
  SOCKS5_CONFIG           Override config file path
  SOCKS5_LOG_LEVEL        Set log level

Examples:
  %s
  %s -debug
  %s -config ./custom.properties
  %s -log-level debug -port 8080
  %s -ip 127.0.0.1 -port 5432
`, os.Args[0], os.Args[0], os.Args[0], os.Args[0], os.Args[0], os.Args[0])
}

func main() {
	// Определяем флаги командной строки
	var (
		debugFlag      bool
		configPathFlag string
		portFlag       string
		ipFlag         string
		logLevelFlag   string
		showHelp       bool
	)

	flag.BoolVar(&debugFlag, "debug", false, "Enable debug logging")
	flag.BoolVar(&debugFlag, "d", false, "Enable debug logging (short)")
	flag.StringVar(&configPathFlag, "config", "", "Path to configuration file")
	flag.StringVar(&portFlag, "port", "", "Override port")
	flag.StringVar(&ipFlag, "ip", "", "Override IP address")
	flag.StringVar(&logLevelFlag, "log-level", "", "Log level: error, warn, info, debug")
	flag.BoolVar(&showHelp, "help", false, "Show help")
	flag.BoolVar(&showHelp, "h", false, "Show help (short)")

	flag.Parse()

	// Показываем справку если запрошено
	if showHelp {
		printUsage()
		return
	}

	// Определяем режим отладки (флаг имеет приоритет над переменной окружения)
	debugMode := debugFlag
	if !debugMode {
		debugMode = os.Getenv("SOCKS5_DEBUG") == "1" || os.Getenv("SOCKS5_DEBUG") == "true"
	}

	// Определяем уровень логирования
	logLevel := logLevelFlag
	if logLevel == "" {
		logLevel = os.Getenv("SOCKS5_LOG_LEVEL")
		if logLevel == "" {
			logLevel = "info"
		}
	}

	// Если включен debug, устанавливаем уровень логирования в debug
	if debugMode && logLevel == "info" {
		logLevel = "debug"
	}

	// Определяем путь к конфигурации
	configPath := config.GetConfigPath(configPathFlag)

	// Загружаем конфигурацию
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		log.Fatalf("Ошибка загрузки конфигурации: %v", err)
	}

	// Переопределяем IP если указан в флагах
	if ipFlag != "" {
		cfg.IP = ipFlag
		if err := cfg.Validate(); err != nil {
			log.Fatalf("Ошибка валидации конфигурации после переопределения IP: %v", err)
		}
	}

	// Переопределяем порт если указан в флагах
	if portFlag != "" {
		port, err := strconv.Atoi(portFlag)
		if err != nil {
			log.Fatalf("Ошибка: порт должен быть числом: %v", err)
		}
		cfg.Port = port
		if err := cfg.Validate(); err != nil {
			log.Fatalf("Ошибка валидации конфигурации после переопределения порта: %v", err)
		}
	}

	// Создаем логгер
	logger := logger.NewFilteredLogger(debugMode, logLevel)

	// Выводим информацию о запуске
	logger.PrintStartupInfo(cfg, configPath, logLevel, debugMode)

	// Создаем сервер с логирующим listener
	server := proxy.NewServer(cfg, logger)

	// Создаем контекст с сигналами для graceful shutdown
	ctx, cancel := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGHUP,
	)
	defer cancel()

	// Запускаем сервер
	if err := server.Start(ctx, cfg.Addr()); err != nil {
		logger.Errorf("Ошибка запуска сервера: %v", err)
		os.Exit(1)
	}

	logger.Info("Сервер завершил работу")
}
