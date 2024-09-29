package logcfg

import (
	"bytes"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"os"
	"runtime"
	"testing"
)

func TestRunLoggerConfig_TableDriven(t *testing.T) {
	tests := []struct {
		name             string
		envLogs          string
		filePath         string
		expectedLogLevel logrus.Level
	}{
		{
			name:             "Test with 'debug' log level",
			envLogs:          "debug",
			filePath:         "./logTest.log",
			expectedLogLevel: logrus.DebugLevel,
		},
		{
			name:             "Test with 'info' log level",
			envLogs:          "info",
			filePath:         "./logTest.log",
			expectedLogLevel: logrus.InfoLevel,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			var logOutput bytes.Buffer
			logrus.SetOutput(&logOutput)

			RunLoggerConfig(tt.envLogs, tt.filePath)

			assert.Equal(t, tt.expectedLogLevel, logrus.GetLevel())

			formatter, ok := logrus.StandardLogger().Formatter.(*logrus.TextFormatter)
			assert.True(t, ok, "formatter should be of type *logrus.TextFormatter")
			assert.NotNil(t, formatter.CallerPrettyfier, "CallerPrettyfier should be set")

			frame := runtime.Frame{File: "/path/to/file.go", Line: 123, Function: "TestFunction"}
			_, file := formatter.CallerPrettyfier(&frame)
			expectedFile := fmt.Sprintf("file.go.123.TestFunction")
			assert.Equal(t, expectedFile, file, "CallerPrettyfier should format file information correctly")

			logrus.SetOutput(os.Stdout)
			logrus.SetLevel(logrus.InfoLevel)
			logrus.SetReportCaller(false)
			logrus.SetFormatter(&logrus.TextFormatter{})
		})
	}
}
