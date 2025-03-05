package utils

import (
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"super-supply-chain/configs"
)

var Logger *zap.Logger

func InitZapLogger() *zap.Logger {
	// 1. 配置编码器（控制台和文件格式）
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder      // 时间格式
	consoleEncoder := zapcore.NewConsoleEncoder(encoderConfig) // 控制台格式
	jsonEncoder := zapcore.NewJSONEncoder(encoderConfig)       // 文件格式

	logFile := ""
	if configs.IsDev() {
		logFile = "logs/app.log"
	} else {
		logFile = "/mnt/logs/app.log"
	}

	// 2. 定义输出目标（控制台 + 文件）
	// 文件输出（带日志切割）
	fileWriter := zapcore.AddSync(&lumberjack.Logger{
		Filename:   logFile,
		MaxSize:    100,  // 单个文件最大100MB
		MaxBackups: 5,    // 保留5个备份
		MaxAge:     30,   // 保留30天
		Compress:   true, // 压缩旧日志
	})

	// 多端输出核心配置
	core := zapcore.NewTee(
		// 控制台输出（开发环境格式）
		zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), zap.DebugLevel),
		// 文件输出（生产环境JSON格式）
		zapcore.NewCore(jsonEncoder, fileWriter, zap.InfoLevel),
	)

	// 3. 创建Logger实例
	logger := zap.New(core, zap.AddCaller())
	return logger
}

func InitLogger() *zap.Logger {
	Logger = InitZapLogger()
	defer Logger.Sync()
	return Logger
}
