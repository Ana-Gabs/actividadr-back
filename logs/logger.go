// ./logs/logger.go
package logs

import (
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
)

var Logger *logrus.Logger

func init() {
	Logger = logrus.New()

	logDir := "log" // Se crea dentro del root del proyecto

	// Crear el directorio si no existe
	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		err := os.MkdirAll(logDir, os.ModePerm)
		if err != nil {
			logrus.Fatalf("No se pudo crear el directorio de logs: %v", err)
		}
	}

	// Archivos de log
	errorLogPath := filepath.Join(logDir, "error.log")
	combinedLogPath := filepath.Join(logDir, "combined.log")
	allLogPath := filepath.Join(logDir, "all.log")

	// Crear archivos
	errorFile, _ := os.OpenFile(errorLogPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	combinedFile, _ := os.OpenFile(combinedLogPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	allFile, _ := os.OpenFile(allLogPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)

	// Configuración de logger
	Logger.SetFormatter(&logrus.JSONFormatter{})

	Logger.SetOutput(allFile) // Por defecto, manda todo a all.log

	// Agregar múltiples salidas
	Logger.AddHook(NewFileHook(errorFile, logrus.ErrorLevel))
	Logger.AddHook(NewFileHook(combinedFile, logrus.InfoLevel))
}

// FileHook permite enviar logs a archivos separados por nivel
type FileHook struct {
	Writer    *os.File
	LogLevels []logrus.Level
}

func NewFileHook(file *os.File, level logrus.Level) *FileHook {
	return &FileHook{
		Writer:    file,
		LogLevels: []logrus.Level{level},
	}
}

func (hook *FileHook) Fire(entry *logrus.Entry) error {
	line, err := entry.String()
	if err != nil {
		return err
	}
	_, err = hook.Writer.Write([]byte(line))
	return err
}

func (hook *FileHook) Levels() []logrus.Level {
	return hook.LogLevels
}
