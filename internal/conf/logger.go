package conf

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"time"
)

type Logger struct {
	*zap.Logger
}

// NewZapLogger 生产环境推荐配置
func NewZapLogger(conf *Conf) *Logger {
	// 1. 配置编码器
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.EpochTimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// 2. 配置日志轮转
	fileWriter := zapcore.AddSync(&lumberjack.Logger{
		Filename:   "./logs/log.log",
		MaxSize:    100, // MB
		MaxBackups: 30,
		MaxAge:     30, // days
		Compress:   true,
	})

	// 3. 创建核心
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		zapcore.NewMultiWriteSyncer(fileWriter),
		zap.InfoLevel,
	)

	// 4. 采样配置
	core = zapcore.NewSamplerWithOptions(
		core,
		time.Second,
		100, // 初始采样数
		100, // 之后采样间隔
		//zapcore.SamplerHook(func(dropped zapcore.SamplingDecision) {
		//	if dropped {
		//		metrics.Incr("log.dropped")
		//	}
		//}),
	)

	return &Logger{
		Logger: zap.New(core,
			zap.AddCaller(),
			zap.AddStacktrace(zapcore.ErrorLevel),
		),
	}
	//return zap.New(core,
	//	zap.AddCaller(),
	//	zap.AddStacktrace(zapcore.ErrorLevel),
	//)
}
