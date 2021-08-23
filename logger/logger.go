package logger

import (
	"fmt"
	"os"
	"time"

	"github.com/celtics-auto/ebiten-server/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	DEVELOPMENT_STRING = "development"
)

/*
	Creates log dir if does not exist
*/
func createLogDirectory(currentPath string) error {
	dirPath := fmt.Sprintf("%s/logs", currentPath)

	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		return os.Mkdir("logs", os.ModePerm)
	}

	return nil
}

// TODO: customize filename
func getFileWritter(currentPath string) (zapcore.WriteSyncer, error) {
	file, err := os.OpenFile(currentPath+"/logs/filename.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		return nil, err
	}
	return zapcore.AddSync(file), nil
}

func getWriters(cfg *config.Logger) (zapcore.WriteSyncer, error) {
	writersArray := []zapcore.WriteSyncer{}

	if cfg.File {
		currentPath, err := os.Getwd()
		if err != nil {
			return nil, err
		}

		if err := createLogDirectory(currentPath); err != nil {
			return nil, err
		}

		fileSyncer, wErr := getFileWritter(currentPath)
		if wErr != nil {
			return nil, wErr
		}

		writersArray = append(writersArray, fileSyncer)
	}

	if cfg.Stdout {
		writersArray = append(writersArray, zapcore.AddSync(os.Stdout))
	}

	allWriters := zap.CombineWriteSyncers(writersArray...)

	return allWriters, nil
}

func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()

	encoderConfig.EncodeTime = zapcore.TimeEncoder(func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.UTC().Format("02-01-2006 15:04:05"))
	})

	encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	return zapcore.NewConsoleEncoder(encoderConfig)
}

func getLogLevel(env string) zapcore.Level {
	if env == DEVELOPMENT_STRING {
		return zapcore.DebugLevel
	}

	return zapcore.WarnLevel
}

func Init(cfg *config.Logger, env string) error {
	allWriters, err := getWriters(cfg)
	if err != nil {
		return err
	}

	logLevel := getLogLevel(env)
	encoder := getEncoder()
	core := zapcore.NewCore(encoder, allWriters, logLevel)
	logg := zap.New(core, zap.AddCaller())

	zap.ReplaceGlobals(logg)

	return nil
}
