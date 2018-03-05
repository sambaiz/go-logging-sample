package log

import (
	"time"

	"github.com/sambaiz/go-logging-sample/config"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	lumberjack "gopkg.in/natefinch/lumberjack.v2"
)

// Logger outputs logs to LOG_PATH
var Logger = newLogger()

// Shutdown should be called before shutdown app
func Shutdown() error {
	return Logger.Sync()
}

func newLogger() *zap.Logger {
	zapConfig := zap.NewProductionConfig()
	encoderConfig := zapConfig.EncoderConfig
	encoderConfig.EncodeTime = jstTimeEncoder
	enc := zapcore.NewJSONEncoder(encoderConfig)
	w := zapcore.AddSync(
		&lumberjack.Logger{
			Filename:   viper.GetString(config.LogPath),
			MaxSize:    viper.GetInt(config.LogRotateMaxSize), // MB
			MaxBackups: viper.GetInt(config.LogRotateMaxBackups),
			MaxAge:     viper.GetInt(config.LogRotateMaxDays),
		},
	)
	return zap.New(
		zapcore.NewCore(enc, w, zapConfig.Level),
	)
}

func jstTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	const layout = "2006-01-02 15:04:05"
	jst := time.FixedZone("Asia/Tokyo", 9*60*60)
	enc.AppendString(t.In(jst).Format(layout))
}
