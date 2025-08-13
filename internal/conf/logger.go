package conf

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"time"
)

type Logger struct {
	//*zap.Logger
	*zap.SugaredLogger
}

// NewZapLogger 生产环境推荐配置
func NewZapLogger(conf *Conf) *Logger {
	// 1. 配置编码器 - 使用更可读的时间格式
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "timestamp",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder, // 改进：使用ISO8601格式
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// 2. 从配置获取日志路径和其他参数
	logPath := conf.Conf.GetString("log.path")
	if logPath == "" {
		logPath = "./logs/log.log"
	}

	maxSize := conf.Conf.GetInt("log.max_size")
	if maxSize == 0 {
		maxSize = 100
	}

	maxBackups := conf.Conf.GetInt("log.max_backups")
	if maxBackups == 0 {
		maxBackups = 30
	}

	maxAge := conf.Conf.GetInt("log.max_age")
	if maxAge == 0 {
		maxAge = 30
	}

	// 3. 配置日志轮转
	fileWriter := zapcore.AddSync(&lumberjack.Logger{
		Filename:   logPath,
		MaxSize:    maxSize, // MB
		MaxBackups: maxBackups,
		MaxAge:     maxAge, // days
		Compress:   conf.Conf.GetBool("log.compress"),
		LocalTime:  conf.Conf.GetBool("server.dev"), // 使用本地时间
	})

	// 4. 创建写入器 - 添加控制台输出
	var writers []zapcore.WriteSyncer
	writers = append(writers, fileWriter)

	// 开发环境添加控制台输出
	if conf.Conf.GetBool("server.dev") {
		writers = append(writers, zapcore.AddSync(os.Stdout))
	}

	// 5. 获取日志级别
	logLevel := zap.InfoLevel
	levelStr := conf.Conf.GetString("log.level")
	if level, err := zapcore.ParseLevel(levelStr); err == nil {
		logLevel = level
	}

	// 6. 创建核心
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		zapcore.NewMultiWriteSyncer(writers...),
		logLevel,
	)

	// 7. 采样配置
	core = zapcore.NewSamplerWithOptions(
		core,
		time.Second,
		5000, // 初始采样数100
		100,  // 之后采样间隔
	)

	// 8. 创建logger
	logger := zap.New(core,
		zap.AddCaller(),
		zap.AddStacktrace(zapcore.ErrorLevel),
		zap.AddCallerSkip(1),
		zap.Fields(zap.String("server", conf.Conf.GetString("server.name"))),
	)
	// 添加服务名称字段
	//serviceName := conf.Conf.GetString("service.name")
	//if serviceName != "" {
	//	logger = logger.With(zap.String("service", serviceName))
	//}
	return &Logger{
		SugaredLogger: logger.Sugar(),
	}
}
