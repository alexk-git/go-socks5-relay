package proxy

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/things-go/go-socks5"

	"go-socks5-relay/internal/config"
	"go-socks5-relay/internal/logger"
)

// Server представляет собой SOCKS5 сервер с поддержкой graceful shutdown
type Server struct {
	socksServer *socks5.Server
	listener    net.Listener
	logger      *logger.FilteredLogger
	stopChan    chan struct{}
}

// NewServer создает новый экземпляр сервера
func NewServer(cfg *config.Config, logger *logger.FilteredLogger) *Server {
	// Настраиваем аутентификацию
	credentials := socks5.StaticCredentials{
		cfg.Username: cfg.Password,
	}

	// Создаем SOCKS5 сервер
	socksServer := socks5.NewServer(
		socks5.WithCredential(credentials),
		socks5.WithLogger(logger),
	)

	return &Server{
		socksServer: socksServer,
		logger:      logger,
		stopChan:    make(chan struct{}),
	}
}

// Start запускает сервер
func (s *Server) Start(ctx context.Context, addr string) error {
	// Создаем логирующий listener
	var err error
	s.listener, err = NewLoggingListener(addr, s.logger)
	if err != nil {
		return fmt.Errorf("не удалось создать listener: %w", err)
	}

	// Запускаем горутину для graceful shutdown
	go s.handleShutdown(ctx)

	s.logger.Infof("SOCKS5-прокси запущен на %s (TCP + UDP)", addr)

	// Запускаем сервер (блокирующий вызов)
	if err := s.socksServer.Serve(s.listener); err != nil {
		// Проверяем, не была ли ошибка вызвана остановкой сервера
		select {
		case <-s.stopChan:
			s.logger.Info("Сервер остановлен корректно")
			return nil
		default:
			return fmt.Errorf("ошибка работы SOCKS5-сервера: %w", err)
		}
	}

	return nil
}

// handleShutdown ожидает сигнал завершения и gracefully останавливает сервер
func (s *Server) handleShutdown(ctx context.Context) {
	<-ctx.Done()

	s.logger.Info("Получен сигнал завершения, останавливаем сервер...")

	// Создаем таймаут для graceful shutdown
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Закрываем listener, чтобы новые соединения не принимались
	if err := s.listener.Close(); err != nil {
		s.logger.Errorf("Ошибка при закрытии listener: %v", err)
	}

	// Даем время на завершение активных соединений
	<-shutdownCtx.Done()
	if shutdownCtx.Err() == context.DeadlineExceeded {
		s.logger.Info("Таймаут graceful shutdown, принудительное завершение")
	}

	close(s.stopChan)
}
