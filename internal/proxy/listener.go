package proxy

import (
	"net"
	"time"

	"go-socks5-relay/internal/logger"
)

// LoggingListener оборачивает net.Listener и логирует все подключения
type LoggingListener struct {
	net.Listener
	logger *logger.FilteredLogger
}

// NewLoggingListener создает новый логирующий listener
func NewLoggingListener(addr string, logger *logger.FilteredLogger) (*LoggingListener, error) {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}

	return &LoggingListener{
		Listener: listener,
		logger:   logger,
	}, nil
}

// Accept принимает соединение и логирует его
func (l *LoggingListener) Accept() (net.Conn, error) {
	conn, err := l.Listener.Accept()
	if err != nil {
		return nil, err
	}

	l.logger.Debugf("Новое соединение от %s", conn.RemoteAddr())
	return conn, nil
}

// LoggingConn оборачивает net.Conn и логирует закрытие соединения
type LoggingConn struct {
	net.Conn
	logger   *logger.FilteredLogger
	remoteAddr string
}

// NewLoggingConn создает новый логирующий коннект
func NewLoggingConn(conn net.Conn, logger *logger.FilteredLogger) *LoggingConn {
	return &LoggingConn{
		Conn:       conn,
		logger:     logger,
		remoteAddr: conn.RemoteAddr().String(),
	}
}

// Close закрывает соединение и логирует это
func (c *LoggingConn) Close() error {
	c.logger.Debugf("Соединение с %s закрыто", c.remoteAddr)
	return c.Conn.Close()
}

// SetDeadline устанавливает дедлайн
func (c *LoggingConn) SetDeadline(t time.Time) error {
	return c.Conn.SetDeadline(t)
}

// SetReadDeadline устанавливает дедлайн на чтение
func (c *LoggingConn) SetReadDeadline(t time.Time) error {
	return c.Conn.SetReadDeadline(t)
}

// SetWriteDeadline устанавливает дедлайн на запись
func (c *LoggingConn) SetWriteDeadline(t time.Time) error {
	return c.Conn.SetWriteDeadline(t)
}
