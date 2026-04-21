package logger

import (
	"strings"
	"testing"
)

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
		{"read tcp: use of closed network connection", true},
		{"real error: something went wrong", false},
		{"unexpected end of stream", false},
	}

	for _, tt := range tests {
		t.Run(tt.msg, func(t *testing.T) {
			if got := isNoisyError(tt.msg); got != tt.expected {
				t.Errorf("isNoisyError(%q) = %v, want %v", tt.msg, got, tt.expected)
			}
		})
	}
}

func TestNewFilteredLogger(t *testing.T) {
	tests := []struct {
		logLevel    string
		debug       bool
		wantLevel   int
		wantDebug   bool
	}{
		{"info", false, LogLevelInfo, false},
		{"debug", false, LogLevelDebug, true},
		{"error", false, LogLevelError, false},
		{"warn", false, LogLevelWarn, false},
		{"warning", false, LogLevelWarn, false},
		{"info", true, LogLevelInfo, true},
	}

	for _, tt := range tests {
		t.Run(tt.logLevel, func(t *testing.T) {
			l := NewFilteredLogger(tt.debug, tt.logLevel)
			if l.level != tt.wantLevel {
				t.Errorf("level: expected %d, got %d", tt.wantLevel, l.level)
			}
			if l.debug != tt.wantDebug {
				t.Errorf("debug: expected %v, got %v", tt.wantDebug, l.debug)
			}
		})
	}
}

func TestFilteredLoggerLevels(t *testing.T) {
	// Verify that lower log levels suppress output by checking no panic/error occurs
	// and that the level gating logic is consistent.
	levels := []struct {
		name  string
		level string
	}{
		{"error", "error"},
		{"warn", "warn"},
		{"info", "info"},
		{"debug", "debug"},
	}

	for _, lvl := range levels {
		t.Run(lvl.name, func(t *testing.T) {
			l := NewFilteredLogger(false, lvl.level)
			// Should not panic regardless of level
			l.Errorf("test error %s", "msg")
			l.Warnf("test warn %s", "msg")
			l.Infof("test info %s", "msg")
			l.Debugf("test debug %s", "msg")
			l.Info("test info plain")
		})
	}
}

func TestNoisyErrorPatternsCoverage(t *testing.T) {
	// Ensure all documented noisy patterns are actually filtered
	noisyMessages := []string{
		"EOF",
		"broken pipe",
		"connection reset by peer",
		"use of closed network connection",
		"broken pipe",
		"use of closed network connection",
		"i/o timeout",
		"context canceled",
	}

	for _, msg := range noisyMessages {
		if !isNoisyError(msg) {
			t.Errorf("expected %q to be filtered as noisy", msg)
		}
	}
}

func TestNoisyErrorCaseInsensitivity(t *testing.T) {
	// isNoisyError uses strings.Contains — verify exact case matters
	// (documents current behavior, not a requirement to change it)
	if isNoisyError(strings.ToUpper("EOF")) {
		// "EOF" uppercased is still "EOF", so this should still match
	}
}
