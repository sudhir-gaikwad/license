package log

import (
	"fmt"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	logLevel zap.AtomicLevel
)

func InitLogger(levelString string) *zap.Logger {
	// Log to the console by default.
	logLevel = zap.NewAtomicLevel()
	encoderCfg := zap.NewProductionEncoderConfig()
	core := zapcore.NewCore(zapcore.NewJSONEncoder(encoderCfg),
		zapcore.Lock(os.Stdout),
		logLevel)

	setLogLevel(levelString)

	return zap.New(core, zap.AddCaller())
}

func setLogLevel(level string) {
	parsedLevel, err := zapcore.ParseLevel(level)
	if err != nil {
		// Fallback to logging at the info level.
		fmt.Printf("Falling back to the info log level. You specified: %s.\n",
			level)
		logLevel.SetLevel(zapcore.InfoLevel)
	} else {
		logLevel.SetLevel(parsedLevel)
	}
}
