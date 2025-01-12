package logger

import (
	"github.com/sirupsen/logrus"
	"os"
)

var Log *logrus.Logger

// InitLogger инициализирует логгер
func InitLogger() {
	Log = logrus.New()
	Log.SetFormatter(&logrus.JSONFormatter{}) // JSON формат
	file, err := os.OpenFile("server.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err == nil {
		Log.SetOutput(file)
	} else {
		Log.Warn("Failed to log to file. Using default stderr.")
	}
	Log.SetLevel(logrus.InfoLevel) // Уровень логирования
}

// LogRequest логирует HTTP-запрос с методом, путём и статусом
func LogRequest(method, path, status string) {
	Log.WithFields(logrus.Fields{
		"method": method,
		"path":   path,
		"status": status,
	}).Info("HTTP Request")
}

// LogDBError логирует ошибки, связанные с базой данных
func LogDBError(operation, message string, err error) {
	Log.WithFields(logrus.Fields{
		"operation": operation,
		"error":     err.Error(),
	}).Error(message)
}

// LogInfo логирует информационные события
func LogInfo(event, message string, fields map[string]interface{}) {
	entry := Log.WithFields(logrus.Fields{
		"event": event,
	})
	for key, value := range fields {
		entry = entry.WithField(key, value)
	}
	entry.Info(message)
}

// LogError логирует ошибки с кастомными полями
func LogError(event, message string, fields map[string]interface{}) {
	entry := Log.WithFields(logrus.Fields{
		"event": event,
	})
	for key, value := range fields {
		entry = entry.WithField(key, value)
	}
	entry.Error(message)
}
