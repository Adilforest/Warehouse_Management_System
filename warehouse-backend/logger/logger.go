package logger

import (
	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"os"
)

var Log *logrus.Logger

// InitLogger инициализирует логгер
func InitLogger() {
	Log = logrus.New()
	Log.SetLevel(logrus.InfoLevel)

	// Логи в файл с ротацией
	logFile := &lumberjack.Logger{
		Filename:   "server.log",
		MaxSize:    10, // Мегабайты
		MaxBackups: 3,  // Максимальное количество резервных файлов
		MaxAge:     28, // Дней хранения
		Compress:   true,
	}

	// Логи в файл и консоль одновременно
	Log.SetOutput(io.MultiWriter(logFile, os.Stdout))

	Log.SetFormatter(&logrus.JSONFormatter{})
}

// LogInfo логирует информационные события
func LogInfo(event, message string, fields map[string]interface{}) {
	Log.WithFields(logrus.Fields(fields)).
		WithField("event", event).
		Info(message)
}

// LogError логирует ошибки с кастомными полями
func LogError(event, message string, fields map[string]interface{}) {
	Log.WithFields(logrus.Fields(fields)).
		WithField("event", event).
		Error(message)
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

// LogWarning логирует предупреждения с кастомными полями
func LogWarning(event, message string, fields map[string]interface{}) {
	Log.WithFields(logrus.Fields(fields)).
		WithField("event", event).
		Warn(message)
}
