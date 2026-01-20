package glogrus

import (
	"bytes"
	"context"
	"errors"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
	glogger "gorm.io/gorm/logger"
)

func TestRedactSQLParams(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "single string literal",
			input:    "SELECT * FROM users WHERE name = 'John Doe'",
			expected: "SELECT * FROM users WHERE name = '[REDACTED]'",
		},
		{
			name:     "multiple string literals",
			input:    "INSERT INTO users (name, email) VALUES ('John Doe', 'john@example.com')",
			expected: "INSERT INTO users (name, email) VALUES ('[REDACTED]', '[REDACTED]')",
		},
		{
			name:     "empty string literal",
			input:    "SELECT * FROM users WHERE name = ''",
			expected: "SELECT * FROM users WHERE name = '[REDACTED]'",
		},
		{
			name:     "string with escaped quotes",
			input:    "INSERT INTO users (name) VALUES ('O''Brien')",
			expected: "INSERT INTO users (name) VALUES ('[REDACTED]')",
		},
		{
			name:     "no string literals",
			input:    "SELECT * FROM users WHERE id = 123",
			expected: "SELECT * FROM users WHERE id = 123",
		},
		{
			name:     "mixed literals and numbers",
			input:    "SELECT * FROM users WHERE name = 'John' AND age > 25",
			expected: "SELECT * FROM users WHERE name = '[REDACTED]' AND age > 25",
		},
		{
			name:     "complex query with multiple literals",
			input:    "UPDATE users SET name = 'Jane', email = 'jane@example.com' WHERE id = 'abc-123'",
			expected: "UPDATE users SET name = '[REDACTED]', email = '[REDACTED]' WHERE id = '[REDACTED]'",
		},
		{
			name:     "string with special characters",
			input:    "INSERT INTO logs (message) VALUES ('Error: connection failed!')",
			expected: "INSERT INTO logs (message) VALUES ('[REDACTED]')",
		},
		{
			name:     "string with newlines",
			input:    "INSERT INTO text (content) VALUES ('Line 1\nLine 2')",
			expected: "INSERT INTO text (content) VALUES ('[REDACTED]')",
		},
		{
			name:     "json data",
			input:    "INSERT INTO data (json) VALUES ('{\"key\": \"value\"}')",
			expected: "INSERT INTO data (json) VALUES ('[REDACTED]')",
		},
		{
			name:     "empty input",
			input:    "",
			expected: "",
		},
		{
			name:     "consecutive string literals",
			input:    "SELECT 'foo' || 'bar'",
			expected: "SELECT '[REDACTED]' || '[REDACTED]'",
		},
		{
			name:     "string in LIKE clause",
			input:    "SELECT * FROM users WHERE name LIKE '%John%'",
			expected: "SELECT * FROM users WHERE name LIKE '[REDACTED]'",
		},
		{
			name:     "string in IN clause",
			input:    "SELECT * FROM users WHERE role IN ('admin', 'moderator', 'user')",
			expected: "SELECT * FROM users WHERE role IN ('[REDACTED]', '[REDACTED]', '[REDACTED]')",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := redactSQLParams(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestLoggerTrace_WithLogSQL(t *testing.T) {
	var buf bytes.Buffer
	logger := logrus.New()
	logger.SetOutput(&buf)
	logger.SetLevel(logrus.TraceLevel)
	logger.SetFormatter(&logrus.JSONFormatter{})

	l := New(Config{
		Logger:                    logger,
		SlowThreshold:             100 * time.Millisecond,
		IgnoreRecordNotFoundError: false,
		LogSQL:                    true,
	})

	ctx := context.Background()
	begin := time.Now()
	sql := "SELECT * FROM users WHERE name = 'secret123'"

	l.Trace(ctx, begin, func() (string, int64) {
		return sql, 1
	}, nil)

	output := buf.String()
	assert.Contains(t, output, `"sql":`)
	assert.Contains(t, output, "[REDACTED]")
	assert.NotContains(t, output, "secret123")
}

func TestLoggerTrace_WithoutLogSQL(t *testing.T) {
	var buf bytes.Buffer
	logger := logrus.New()
	logger.SetOutput(&buf)
	logger.SetLevel(logrus.TraceLevel)
	logger.SetFormatter(&logrus.JSONFormatter{})

	l := New(Config{
		Logger:                    logger,
		SlowThreshold:             100 * time.Millisecond,
		IgnoreRecordNotFoundError: false,
		LogSQL:                    false,
	})

	ctx := context.Background()
	begin := time.Now()
	sql := "SELECT * FROM users WHERE name = 'secret123'"

	l.Trace(ctx, begin, func() (string, int64) {
		return sql, 1
	}, nil)

	output := buf.String()
	// SQL field should not be included in JSON output when LogSQL is false
	assert.NotContains(t, output, `"sql":`)
	assert.NotContains(t, output, "secret123")
	assert.NotContains(t, output, "[REDACTED]")
}

func TestLoggerTrace_WithError(t *testing.T) {
	var buf bytes.Buffer
	logger := logrus.New()
	logger.SetOutput(&buf)
	logger.SetLevel(logrus.TraceLevel)
	logger.SetFormatter(&logrus.JSONFormatter{})

	l := New(Config{
		Logger:                    logger,
		SlowThreshold:             100 * time.Millisecond,
		IgnoreRecordNotFoundError: false,
		LogSQL:                    true,
	})

	ctx := context.Background()
	begin := time.Now()
	sql := "SELECT * FROM users WHERE id = 'bad-id'"
	testErr := errors.New("database error")

	l.Trace(ctx, begin, func() (string, int64) {
		return sql, 0
	}, testErr)

	output := buf.String()
	assert.Contains(t, output, "sql query error")
	assert.Contains(t, output, "[REDACTED]")
	assert.NotContains(t, output, "bad-id")
	assert.Contains(t, output, "database error")
}

func TestLoggerTrace_IgnoreRecordNotFound(t *testing.T) {
	var buf bytes.Buffer
	logger := logrus.New()
	logger.SetOutput(&buf)
	logger.SetLevel(logrus.TraceLevel)
	logger.SetFormatter(&logrus.JSONFormatter{})

	l := New(Config{
		Logger:                    logger,
		SlowThreshold:             100 * time.Millisecond,
		IgnoreRecordNotFoundError: true,
		LogSQL:                    true,
	})

	ctx := context.Background()
	begin := time.Now()
	sql := "SELECT * FROM users WHERE id = 'missing'"

	l.Trace(ctx, begin, func() (string, int64) {
		return sql, 0
	}, gorm.ErrRecordNotFound)

	output := buf.String()
	// Should not contain error message since we're ignoring record not found
	assert.NotContains(t, output, "sql query error")
	assert.NotContains(t, output, "record not found")
}

func TestLoggerTrace_SlowQuery(t *testing.T) {
	var buf bytes.Buffer
	logger := logrus.New()
	logger.SetOutput(&buf)
	logger.SetLevel(logrus.InfoLevel)
	logger.SetFormatter(&logrus.JSONFormatter{})

	l := New(Config{
		Logger:                    logger,
		SlowThreshold:             10 * time.Millisecond,
		IgnoreRecordNotFoundError: false,
		LogSQL:                    true,
	})

	ctx := context.Background()
	begin := time.Now().Add(-50 * time.Millisecond) // Simulate slow query
	sql := "SELECT * FROM users WHERE email = 'slow@example.com'"

	l.Trace(ctx, begin, func() (string, int64) {
		return sql, 100
	}, nil)

	output := buf.String()
	assert.Contains(t, output, "sql query slow")
	assert.Contains(t, output, "[REDACTED]")
	assert.NotContains(t, output, "slow@example.com")
}

func TestLoggerLogMode(t *testing.T) {
	l := &Logger{}
	result := l.LogMode(glogger.Info)
	assert.NotNil(t, result)
	assert.Equal(t, l, result)
}

func TestLoggerInfo(t *testing.T) {
	var buf bytes.Buffer
	logger := logrus.New()
	logger.SetOutput(&buf)
	logger.SetLevel(logrus.InfoLevel)

	l := New(Config{
		Logger: logger,
	})

	l.Info(context.Background(), "test info message: %s", "value")

	output := buf.String()
	assert.Contains(t, output, "test info message: value")
}

func TestLoggerWarn(t *testing.T) {
	var buf bytes.Buffer
	logger := logrus.New()
	logger.SetOutput(&buf)
	logger.SetLevel(logrus.WarnLevel)

	l := New(Config{
		Logger: logger,
	})

	l.Warn(context.Background(), "test warn message: %s", "value")

	output := buf.String()
	assert.Contains(t, output, "test warn message: value")
}

func TestLoggerError(t *testing.T) {
	var buf bytes.Buffer
	logger := logrus.New()
	logger.SetOutput(&buf)
	logger.SetLevel(logrus.ErrorLevel)

	l := New(Config{
		Logger: logger,
	})

	l.Error(context.Background(), "test error message: %s", "value")

	output := buf.String()
	assert.Contains(t, output, "test error message: value")
}

func TestLoggerZeroValue(t *testing.T) {
	var l Logger
	// Should not panic when using zero value
	require.NotPanics(t, func() {
		l.Info(context.Background(), "test")
	})
}

func TestNew_DefaultValues(t *testing.T) {
	l := New(Config{})

	// Should use defaults
	assert.NotNil(t, l)
	assert.NotNil(t, l.logger)
	assert.Equal(t, 500*time.Millisecond, l.slowThreshold)
	assert.False(t, l.ignoreRecordNotFoundError)
	assert.False(t, l.logSQL)
}

func TestNew_CustomValues(t *testing.T) {
	logger := logrus.New()
	l := New(Config{
		Logger:                    logger,
		SlowThreshold:             200 * time.Millisecond,
		IgnoreRecordNotFoundError: true,
		LogSQL:                    true,
	})

	assert.Equal(t, logger, l.logger)
	assert.Equal(t, 200*time.Millisecond, l.slowThreshold)
	assert.True(t, l.ignoreRecordNotFoundError)
	assert.True(t, l.logSQL)
}
