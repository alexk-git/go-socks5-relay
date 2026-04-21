package proxy

import (
	"context"
	"net"
	"strconv"
	"testing"
	"time"

	"go-socks5-relay/internal/config"
	"go-socks5-relay/internal/logger"
)

func newTestLogger() *logger.FilteredLogger {
	return logger.NewFilteredLogger(false, "error")
}

func freeAddr(t *testing.T) (string, int) {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	addr := ln.Addr().String()
	ln.Close()
	_, portStr, _ := net.SplitHostPort(addr)
	port, _ := strconv.Atoi(portStr)
	return addr, port
}

func TestNewServer(t *testing.T) {
	cfg := &config.Config{IP: "127.0.0.1", Port: 1080, Username: "user", Password: "pass"}
	s := NewServer(cfg, newTestLogger())
	if s == nil {
		t.Fatal("NewServer returned nil")
	}
	if s.socksServer == nil {
		t.Error("socksServer should not be nil")
	}
	if s.stopChan == nil {
		t.Error("stopChan should not be nil")
	}
}

func TestServerListens(t *testing.T) {
	addr, port := freeAddr(t)
	cfg := &config.Config{IP: "127.0.0.1", Port: port, Username: "user", Password: "pass"}
	s := NewServer(cfg, newTestLogger())

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go s.Start(ctx, addr)

	var conn net.Conn
	var err error
	for i := 0; i < 20; i++ {
		time.Sleep(10 * time.Millisecond)
		conn, err = net.DialTimeout("tcp", addr, 100*time.Millisecond)
		if err == nil {
			break
		}
	}
	if err != nil {
		t.Fatalf("server did not start listening: %v", err)
	}
	conn.Close()
}

func TestNewLoggingListener(t *testing.T) {
	ln, err := NewLoggingListener("127.0.0.1:0", newTestLogger())
	if err != nil {
		t.Fatalf("NewLoggingListener failed: %v", err)
	}
	defer ln.Close()

	addr := ln.Addr().String()

	accepted := make(chan net.Conn, 1)
	go func() {
		conn, _ := ln.Accept()
		accepted <- conn
	}()

	client, err := net.DialTimeout("tcp", addr, time.Second)
	if err != nil {
		t.Fatalf("could not connect: %v", err)
	}
	defer client.Close()

	select {
	case conn := <-accepted:
		if conn == nil {
			t.Fatal("Accept returned nil conn")
		}
		conn.Close()
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for Accept")
	}
}

func TestLoggingListenerReturnsLoggingConn(t *testing.T) {
	ln, err := NewLoggingListener("127.0.0.1:0", newTestLogger())
	if err != nil {
		t.Fatal(err)
	}
	defer ln.Close()

	accepted := make(chan net.Conn, 1)
	go func() {
		conn, err := ln.Accept()
		if err != nil {
			accepted <- nil
			return
		}
		accepted <- conn
	}()

	client, err := net.DialTimeout("tcp", ln.Addr().String(), time.Second)
	if err != nil {
		t.Fatal(err)
	}
	defer client.Close()

	select {
	case conn := <-accepted:
		if conn == nil {
			t.Fatal("Accept returned nil conn")
		}
		if _, ok := conn.(*LoggingConn); !ok {
			t.Errorf("Accept should return *LoggingConn, got %T", conn)
		}
		conn.Close()
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for Accept")
	}
}

func TestLoggingConnClose(t *testing.T) {
	server, client := net.Pipe()
	defer client.Close()

	lc := NewLoggingConn(server, newTestLogger())
	if err := lc.Close(); err != nil {
		t.Errorf("Close failed: %v", err)
	}
}

func TestLoggingConnRemoteAddr(t *testing.T) {
	ln, err := NewLoggingListener("127.0.0.1:0", newTestLogger())
	if err != nil {
		t.Fatal(err)
	}
	defer ln.Close()

	accepted := make(chan net.Conn, 1)
	go func() {
		conn, _ := ln.Accept()
		accepted <- conn
	}()

	client, err := net.DialTimeout("tcp", ln.Addr().String(), time.Second)
	if err != nil {
		t.Fatal(err)
	}
	defer client.Close()

	select {
	case conn := <-accepted:
		lc := conn.(*LoggingConn)
		if lc.remoteAddr == "" {
			t.Error("remoteAddr should not be empty")
		}
		lc.Close()
	case <-time.After(time.Second):
		t.Fatal("timed out")
	}
}
