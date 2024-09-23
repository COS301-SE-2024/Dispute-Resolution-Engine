package utilities_test

import (
	"errors"
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

func TestLoggerWithRequestID(t *testing.T) {
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

    // Create a new logger and log with request ID
    logger := utilities.NewLogger().WithRequestID("12345")
    logger.Info("test message with request ID")

    // Close the writer and read the captured output
    w.Close()
    var buf [1024]byte
    n, _ := r.Read(buf[:])
    output := string(buf[:n])

    // Check if the output contains the expected log fields
    assert.Contains(t, output, `"timestamp":`)
    assert.Contains(t, output, `"level":"info"`)
    assert.Contains(t, output, `"message":"test message with request ID"`)
    assert.Contains(t, output, `"request_id":"12345"`)
}

func TestLoggerWithFields(t *testing.T) {
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

    // Create a new logger and log with custom fields
    fields := logrus.Fields{"field1": "value1", "field2": "value2"}
    logger := utilities.NewLogger().WithFields(fields)
    logger.Info("test message with fields")

    // Close the writer and read the captured output
    w.Close()
    var buf [1024]byte
    n, _ := r.Read(buf[:])
    output := string(buf[:n])

    // Check if the output contains the expected log fields
    assert.Contains(t, output, `"timestamp":`)
    assert.Contains(t, output, `"level":"info"`)
    assert.Contains(t, output, `"message":"test message with fields"`)
    assert.Contains(t, output, `"field1":"value1"`)
    assert.Contains(t, output, `"field2":"value2"`)
}

func TestLoggerWithError(t *testing.T) {
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

    // Create a new logger and log with an error
    err := errors.New("test error")
    logger := utilities.NewLogger().WithError(err)
    logger.Error("test message with error")

    // Close the writer and read the captured output
    w.Close()
    var buf [1024]byte
    n, _ := r.Read(buf[:])
    output := string(buf[:n])

    // Check if the output contains the expected log fields
    assert.Contains(t, output, `"timestamp":`)
    assert.Contains(t, output, `"level":"error"`)
    assert.Contains(t, output, `"message":"test message with error"`)
    assert.Contains(t, output, `"error":"test error"`)
}