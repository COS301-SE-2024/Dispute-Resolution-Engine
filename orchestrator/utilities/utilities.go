package utilities

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"
)

const (
	API_URL = "http://api:9000/event"
)

type APIReq struct {
    ID int64 `json:"id"`
	CurrentState string `json:"current_state"`
}

var Log *logrus.Logger

func init() {
    Log = logrus.New()
    Log.SetFormatter(&logrus.JSONFormatter{
        FieldMap: logrus.FieldMap{
            logrus.FieldKeyTime:  "timestamp",
            logrus.FieldKeyLevel: "level",
            logrus.FieldKeyMsg:   "message",
        },
    })
    Log.SetOutput(os.Stdout)
    Log.SetLevel(logrus.InfoLevel)
}

type Logger struct {
    *logrus.Entry
}

func NewLogger() *Logger {
    return &Logger{Log.WithFields(logrus.Fields{})}
}

func (l *Logger) WithRequestID(requestID string) *Logger {
    return &Logger{l.WithField("request_id", requestID)}
}

func (l *Logger) WithFields(fields logrus.Fields) *Logger {
    return &Logger{l.Entry.WithFields(fields)}
}

func (l *Logger) WithError(err error) *Logger {
    return &Logger{l.Entry.WithError(err)}
}

func (l *Logger) LogWithCaller() *Logger {
	if pc, file, line, ok := runtime.Caller(1); ok {
		fName := runtime.FuncForPC(pc).Name()
		splitFName := strings.Split(fName, ".")
		return &Logger{
			l.WithField("caller_file", file).
				WithField("caller_line", line).
				WithField("caller_function", splitFName[len(splitFName)-1]),
		}
	}
	return l
}

func StringToInt(s string) (int, error) {
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0, err
	}
	return i, nil
}

// This function makes an http POST request to the specified endpoint with the specified payload
func APIPostRequest(endpoint string, payload APIReq) (string, error) {
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}
	
	resp, err := http.Post(endpoint, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	
	return string(body), nil
}
