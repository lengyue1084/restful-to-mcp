package logger

import (
	"context"

	"flow-bridge-mcp/internal/conf"

	logger2 "gorm.io/gorm/logger"

	"github.com/google/wire"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"time"
)

var ProviderSet = wire.NewSet(NewLogger, NewGormLogger)

// Interface Logger 定义了应用日志记录器的标准接口
type LoggerInterface interface {
	Debug(args ...interface{})
	Info(args ...interface{})
	Warn(args ...interface{})
	Error(args ...interface{})
	Fatal(args ...interface{})
	Panic(args ...interface{})

	DebugF(template string, args ...interface{})
	InfoF(template string, args ...interface{})
	WarnF(template string, args ...interface{})
	ErrorF(template string, args ...interface{})
	FatalF(template string, args ...interface{})
	PanicF(template string, args ...interface{})

	DebugWithContext(ctx context.Context, template string, args ...interface{})
	InfoWithContext(ctx context.Context, template string, args ...interface{})
	WarnWithContext(ctx context.Context, template string, args ...interface{})
	ErrorWithContext(ctx context.Context, template string, args ...interface{})
	FatalWithContext(ctx context.Context, template string, args ...interface{})
	PanicWithContext(ctx context.Context, template string, args ...interface{})

	With(args ...interface{}) *zap.SugaredLogger
	WithContext(ctx context.Context) *Logger
	WithError(err error) *Logger
	WithFields(fields map[string]interface{}) *Logger
	WithField(key string, value interface{}) *Logger
}

var _ LoggerInterface = (*Logger)(nil)

type Logger struct {
	*zap.SugaredLogger
}

func (l *Logger) WithContext(ctx context.Context) *Logger {
	var fields []interface{}
	if traceIdValue := ctx.Value("traceId"); traceIdValue != nil {
		if traceId, ok := traceIdValue.(string); ok && traceId != "" {
			fields = append(fields, "traceId", traceId)
		}
	}
	if spanIdValue := ctx.Value("spanId"); spanIdValue != nil {
		if spanId, ok := spanIdValue.(string); ok && spanId != "" {
			fields = append(fields, "spanId", spanId)
		}
	}
	l.SugaredLogger = l.SugaredLogger.With(fields...)
	return l
}

// NewLogger 生产环境推荐配置
func NewLogger(conf *conf.Conf) *Logger {
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

func (l *Logger) With(args ...interface{}) *zap.SugaredLogger {
	return l.SugaredLogger.With(args...)
}
func (l *Logger) Error(args ...interface{}) {
	l.SugaredLogger.Error(args...)
}
func (l *Logger) Info(args ...interface{}) {
	l.SugaredLogger.Info(args...)
}
func (l *Logger) Warn(args ...interface{}) {
	l.SugaredLogger.Warn(args...)
}
func (l *Logger) Debug(args ...interface{}) {
	l.SugaredLogger.Debug(args...)
}
func (l *Logger) Fatal(args ...interface{}) {
	l.SugaredLogger.Fatal(args...)
}
func (l *Logger) Panic(args ...interface{}) {
	l.SugaredLogger.Panic(args...)
}
func (l *Logger) WithError(err error) *Logger {
	l.SugaredLogger = l.SugaredLogger.With("error", err)
	return l
}
func (l *Logger) WithFields(fields map[string]interface{}) *Logger {
	l.SugaredLogger = l.SugaredLogger.With(fields)
	return l
}

func (l *Logger) WithField(key string, value interface{}) *Logger {
	l.SugaredLogger = l.SugaredLogger.With(key, value)
	return l
}
func (l *Logger) InfoF(template string, args ...interface{}) {
	l.SugaredLogger.Infof(template, args...)
}
func (l *Logger) InfoWithContext(ctx context.Context, template string, args ...interface{}) {
	l.WithContext(ctx)
	l.SugaredLogger.Infof(template, args...)
}
func (l *Logger) ErrorF(template string, args ...interface{}) {
	l.SugaredLogger.Errorf(template, args...)
}
func (l *Logger) ErrorWithContext(ctx context.Context, template string, args ...interface{}) {
	l.WithContext(ctx)
	l.SugaredLogger.Errorf(template, args...)
}

func (l *Logger) WarnF(template string, args ...interface{}) {
	l.SugaredLogger.Warnf(template, args...)
}
func (l *Logger) WarnWithContext(ctx context.Context, template string, args ...interface{}) {
	l.WithContext(ctx)
	l.SugaredLogger.Warnf(template, args...)
}
func (l *Logger) DebugF(template string, args ...interface{}) {
	l.SugaredLogger.Debugf(template, args...)
}
func (l *Logger) DebugWithContext(ctx context.Context, template string, args ...interface{}) {
	l.WithContext(ctx)
	l.SugaredLogger.Debugf(template, args...)
}
func (l *Logger) FatalF(template string, args ...interface{}) {
	l.SugaredLogger.Fatalf(template, args...)
}
func (l *Logger) FatalWithContext(ctx context.Context, template string, args ...interface{}) {
	l.WithContext(ctx)
	l.SugaredLogger.Fatalf(template, args...)
}
func (l *Logger) PanicF(template string, args ...interface{}) {
	l.SugaredLogger.Panicf(template, args...)
}
func (l *Logger) PanicWithContext(ctx context.Context, template string, args ...interface{}) {
	l.WithContext(ctx)
	l.SugaredLogger.Panicf(template, args...)
}

var _ logger2.Interface = (*GormLogger)(nil)

type GormLogger struct {
	log *Logger
}

func NewGormLogger(log *Logger) *GormLogger {
	return &GormLogger{log: log}
}

// LogMode log mode
func (g *GormLogger) LogMode(level logger2.LogLevel) logger2.Interface {
	return g
}

// Info print info
func (g *GormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	g.log.InfoWithContext(ctx, msg, data)
}

// Warn print warn messages
func (g *GormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	g.log.WarnWithContext(ctx, msg, data)
}

// Error print error messages
func (g *GormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	g.log.ErrorWithContext(ctx, msg, data)
}

// Trace print sql message
func (g *GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	elapsed := time.Since(begin)
	sql, rows := fc()

	g.log.With(
		zap.String("sql", sql),
		zap.Duration("elapsed", elapsed),
		zap.Int64("rows", rows),
	).Info("SQL executed")
}
