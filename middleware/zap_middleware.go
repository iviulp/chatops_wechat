package middleware

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"time"
	"yuguosheng/int/mychatops/config"
)

var MyLogger *zap.Logger

func init() {
	var myEncodeTime zapcore.TimeEncoder
	myEncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format(config.GetSystemConf().LogConf.LogTimeFormat))
	}

	logfile, err := os.OpenFile(config.GetSystemConf().LogConf.LogOutPutPath, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		panic("Failed to Open Log File")
	}
	defer logfile.Close()

	zapconfig := zap.Config{
		Level:       zap.NewAtomicLevelAt(zap.DebugLevel),
		Development: false,
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		Encoding: "console",
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:          "打印时间",
			CallerKey:        "方法",
			LevelKey:         "日志级别",
			NameKey:          "logger",
			FunctionKey:      zapcore.OmitKey,
			MessageKey:       "信息",
			StacktraceKey:    "堆栈信息",
			LineEnding:       zapcore.DefaultLineEnding,
			EncodeLevel:      zapcore.CapitalColorLevelEncoder,
			EncodeTime:       myEncodeTime,
			EncodeDuration:   zapcore.SecondsDurationEncoder,
			EncodeCaller:     zapcore.ShortCallerEncoder,
			ConsoleSeparator: " ",
		},
		OutputPaths:      []string{"stdout", logfile.Name()},
		ErrorOutputPaths: []string{"stderr"},
	}
	MyLogger, _ = zapconfig.Build()
	defer func(MyLogger *zap.Logger) {
		err := MyLogger.Sync()
		if err != nil {

		}
	}(MyLogger)
}

// LoggerToFile 日志记录到文件
func LoggerToFile() gin.HandlerFunc {

	return func(c *gin.Context) {
		// 开始时间
		startTime := time.Now()
		// 处理请求
		c.Next()
		// 结束时间
		endTime := time.Now()
		// 执行时间
		latencyTime := endTime.Sub(startTime)
		// 请求方式
		reqMethod := c.Request.Method
		// 请求路由
		reqUri := c.Request.RequestURI
		// 请求值
		reqForm := c.Request.Form
		// 状态码
		statusCode := c.Writer.Status()
		// 请求IP
		clientIP := c.ClientIP()

		MyLogger.Debug("Request",
			zap.Any("status_code", statusCode),
			zap.Any("latency_time", latencyTime),
			zap.Any("client_ip", clientIP),
			zap.Any("req_method", reqMethod),
			zap.Any("req_uri", reqUri),
			zap.Any("req_form", reqForm),
		)

	}
}
