package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime/debug"
	"strings"
	"super-supply-chain/configs"
	"time"
)

func InitZapLogger() *zap.Logger {
	// 日志切割配置
	lumberJackLogger := &lumberjack.Logger{
		Filename:   "./logs/app.log", // 日志文件路径
		MaxSize:    100,              // 单文件最大 100MB
		MaxBackups: 5,                // 保留 5 个旧文件
		MaxAge:     30,               // 保留 30 天
		Compress:   true,             // 压缩旧文件
	}

	// Zap 编码器配置
	var encoderConfig zapcore.EncoderConfig
	var encoder zapcore.Encoder
	var writeSyncer zapcore.WriteSyncer
	if configs.IsDev() {
		encoderConfig = zap.NewDevelopmentEncoderConfig()
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
		writeSyncer = zapcore.AddSync(os.Stdout)
	} else {
		encoderConfig = zap.NewProductionEncoderConfig()
		encoder = zapcore.NewJSONEncoder(encoderConfig)
		writeSyncer = zapcore.AddSync(lumberJackLogger)
	}

	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder   // 时间格式
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder // 日志级别显示为大写（INFO）

	// 创建 Zap 核心组件
	core := zapcore.NewCore(
		encoder,       // JSON 格式
		writeSyncer,   // 输出到文件
		zap.InfoLevel, // 日志级别
	)

	// 构建 Logger 实例
	logger := zap.New(core, zap.AddCaller())
	return logger
}

func GinZapLogger(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		// 处理请求
		c.Next()

		// 记录日志
		cost := time.Since(start)
		logger.Info("HTTP Request",
			zap.Int("status", c.Writer.Status()),
			zap.String("method", c.Request.Method),
			zap.String("path", path),
			zap.String("query", query),
			zap.String("ip", c.ClientIP()),
			zap.String("user-agent", c.Request.UserAgent()),
			zap.Duration("cost", cost),
		)
	}
}

func GinZapRecovery(logger *zap.Logger, stack bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// 处理 broken pipe 等网络错误
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") {
							brokenPipe = true
						}
					}
				}

				dumpBytes, _ := httputil.DumpRequest(c.Request, false)

				// 记录错误详情
				if !brokenPipe {
					logger.Error("Recovered from panic",
						zap.Any("error", err),
						zap.String("request", string(dumpBytes)),
						zap.String("stack", string(debug.Stack())),
					)
				}

				// 返回 500 响应
				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()
		c.Next()
	}
}
