package logger

import (
	"fmt"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
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
func getLogWritter(currentPath string) (zapcore.WriteSyncer, error) {
	file, err := os.OpenFile(currentPath+"/logs/filename.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		return nil, err
	}
	return zapcore.AddSync(file), nil
}

func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	// TODO: format time
	// encoderConfig.EncodeTime = zapcore.TimeEncoder(func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	// 	enc.AppendString(t.UTC().Format(""))
	// })
	encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	return zapcore.NewConsoleEncoder(encoderConfig)
}

func Init() error {
	currentPath, err := os.Getwd()
	if err != nil {
		return err
	}

	if err := createLogDirectory(currentPath); err != nil {
		return err
	}

	// TODO: add option to log to stdout via config
	writerSync, wErr := getLogWritter(currentPath)
	if wErr != nil {
		return wErr
	}

	encoder := getEncoder()
	core := zapcore.NewCore(encoder, writerSync, zapcore.DebugLevel)
	logg := zap.New(core, zap.AddCaller())

	zap.ReplaceGlobals(logg)

	return nil
}
