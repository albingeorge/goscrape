package logger

import (
	"os"
	"time"

	logrus_stack "github.com/Gurpartap/logrus-stack"
	"github.com/sirupsen/logrus"
)

//Logger ... logger wrapper struct
type Logger struct {
	Entry *logrus.Entry
}

// Init ...initialise logger
func Init() {
	logrus.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: time.RFC3339Nano,
	})

	logrus.SetOutput(os.Stdout)

	logrus.AddHook(logrus_stack.StandardHook())
}

// Get ...get instance of Logger
func Get() *Logger {
	return &Logger{
		Entry: logrus.NewEntry(logrus.StandardLogger()),
	}
}

//Info ... log info level
func (lgr *Logger) Info(traceCode string, fields map[string]interface{}) {
	lgr.Entry.WithFields(fields).Info(traceCode)
}
