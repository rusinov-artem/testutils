package testutils

import (
	"bufio"
	"bytes"
	"strings"
	"testing"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func SpyLogger() (*zap.Logger, *bytes.Buffer) {
	logger, _ := zap.NewDevelopment()
	logBuffer := bytes.NewBufferString("")
	logger = logger.WithOptions(zap.WrapCore(func(_ zapcore.Core) zapcore.Core {
		return zapcore.NewCore(
			zapcore.NewJSONEncoder(zap.NewDevelopmentEncoderConfig()),
			zapcore.AddSync(logBuffer),
			zapcore.DebugLevel,
		)
	}))
	return logger, logBuffer
}

type Logs struct {
	buf *bytes.Buffer
	t   *testing.T
}

func NewLogs(t *testing.T, b *bytes.Buffer) *Logs {
	return &Logs{buf: b, t: t}
}

func (l *Logs) Contains(str ...string) bool {
	l.t.Helper()
	sc := bufio.NewScanner(l.buf)
	sc.Split(bufio.ScanLines)

	for sc.Scan() {
		rec := sc.Text()

		found := true
		for i := range str {
			if !strings.Contains(rec, str[i]) {
				found = false
				break
			}
		}

		if found {
			return true
		}
	}

	l.t.Errorf("logs does not countain: %s", strings.Join(str, " | "))
	return false
}

func (l *Logs) String() string {
	return l.buf.String()
}

func (l *Logs) SetT(t *testing.T) {
	l.t = t
}
