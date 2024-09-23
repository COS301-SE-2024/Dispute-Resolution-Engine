package utilities_test

import (
	"os"
	"testing"

	"orchestrator/utilities"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestInit(t *testing.T) {
	// Backup the original stdout
	originalStdout := os.Stdout
	defer func() { os.Stdout = originalStdout }()

	// Create a pipe to capture stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Reinitialize the logger to capture the output
	utilities.Log = logrus.New()
	utilities.Log.SetFormatter(&logrus.JSONFormatter{
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyTime:  "timestamp",
			logrus.FieldKeyLevel: "level",
			logrus.FieldKeyMsg:   "message",
		},
	})
	utilities.Log.SetOutput(os.Stdout)
	utilities.Log.SetLevel(logrus.InfoLevel)

	// Log a test message
	utilities.Log.Info("test message")

	// Close the writer and read the captured output
	w.Close()
	var buf [1024]byte
	n, _ := r.Read(buf[:])
	output := string(buf[:n])

	// Check if the output contains the expected log fields
	assert.Contains(t, output, `"timestamp":`)
	assert.Contains(t, output, `"level":"info"`)
	assert.Contains(t, output, `"message":"test message"`)
}

func TestLogWithCaller(t *testing.T) {
    // Backup the original stdout
    originalStdout := os.Stdout
    defer func() { os.Stdout = originalStdout }()

    // Create a pipe to capture stdout
    r, w, _ := os.Pipe()
    os.Stdout = w

    // Reinitialize the logger to capture the output
    utilities.Log = logrus.New()
    utilities.Log.SetFormatter(&logrus.JSONFormatter{
        FieldMap: logrus.FieldMap{
            logrus.FieldKeyTime:  "timestamp",
            logrus.FieldKeyLevel: "level",
            logrus.FieldKeyMsg:   "message",
        },
    })
    utilities.Log.SetOutput(os.Stdout)
    utilities.Log.SetLevel(logrus.InfoLevel)

    // Create a new logger and log with caller information
    logger := utilities.NewLogger().LogWithCaller()
    logger.Info("test message with caller")

    // Close the writer and read the captured output
    w.Close()
    var buf [1024]byte
    n, _ := r.Read(buf[:])
    output := string(buf[:n])

    // Check if the output contains the expected log fields
    assert.Contains(t, output, `"timestamp":`)
    assert.Contains(t, output, `"level":"info"`)
    assert.Contains(t, output, `"message":"test message with caller"`)
    assert.Contains(t, output, `"caller_file":`)
    assert.Contains(t, output, `"caller_line":`)
    assert.Contains(t, output, `"caller_function":"TestLogWithCaller"`)
}