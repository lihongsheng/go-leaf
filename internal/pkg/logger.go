package pkg

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/go-kratos/kratos/v2/transport/http"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"message-center/internal/conf"
	"message-center/internal/types"
	"os"
	"time"
)

var DefaultMessageKey = "msg"
var DefaultErrorKey = "error"
var fileLine = "file_line"

type ZapLog struct {
	log *zap.Logger
}

type Helper struct {
	log log.Logger
}

func NewHelper(log log.Logger) *Helper {
	return &Helper{log: log}
}

// NewZapFileLog 使用zap 记录日志
// 如果有需要，可以创建多个 不同级别的zap示例，来写入不同的文件中
// 后续把配置传进来，用配置设置定义
func NewZapFileLog(c conf.Log) log.Logger {
	// 自定义输出源,及切割文件的方式
	// 使用lumberjack
	var ws zapcore.WriteSyncer
	if c.MaxSize < 1 {
		c.MaxSize = 100
	}
	if c.MaxBackUp < 1 {
		c.MaxBackUp = 3
	}
	if c.MaxAge < 1 {
		c.MaxAge = 7
	}
	// 文件为空的话，默认标准输出接口
	if c.Path != "" {
		output := &lumberjack.Logger{
			Filename:   c.Path,
			MaxSize:    int(c.MaxSize),   // 日志文件最大 单位MB
			MaxBackups: int(c.MaxBackUp), // 最大备份文件数
			MaxAge:     int(c.MaxAge),    // 文件最大保存日期
			Compress:   false,
		}
		if c.MaxSize > 0 {
			output.MaxSize = int(c.MaxSize)
		}
		if c.MaxAge > 0 {
			output.MaxAge = int(c.MaxAge)
		}
		output.Compress = c.Compress
		ws = zapcore.AddSync(output)
	} else {
		ws = zapcore.AddSync(os.Stdout)
	}
	ec := zap.NewProductionEncoderConfig()
	ec.CallerKey = "file" // 展示打印日志的文件
	ec.TimeKey = "ts"     // 时间的key
	ec.LevelKey = "level" // 日志级别的key
	ec.MessageKey = "msg" // 日志提示的key
	// 时间格式
	ec.EncodeTime = func(time time.Time, encoder zapcore.PrimitiveArrayEncoder) {
		encoder.AppendString(time.Format("2006-01-02T15:04:05.000000Z"))
	}
	// json 格式输出
	en := zapcore.NewJSONEncoder(ec)
	core := zapcore.NewCore(en, ws, zapcore.InfoLevel)
	logger := zap.New(core, // 可以定义多个 输出zapcore.NewTee(core,errorFileCore)
		zap.AddCaller(),                      // 记录调用位置
		zap.AddCallerSkip(3),                 // 日志封装了一层往上跳3级记录文件执行位置
		zap.AddStacktrace(zapcore.WarnLevel), // 增加记录调用记录栈
		zap.Hooks(func(entry zapcore.Entry) error {
			// 或者利用zap 的 hooks进行拦截报警,错误信息较少
			// entry.Message
			return nil
		}),
	)
	l := &ZapLog{
		log: logger,
	}
	return l
}

// Log 必须是偶数
// 不要打印 Fatal 级别的日志会使程序down掉
func (l *ZapLog) Log(level log.Level, keyvals ...interface{}) error {
	keyLen := len(keyvals)
	if keyLen == 0 || keyLen%2 != 0 {
		l.log.Warn(fmt.Sprint("Keyvalues must appear in pairs: ", keyvals))
		return nil
	}
	var msg string

	//data := make([]zap.Field, 0, (keyLen/2)+1)
	data := make(map[string]interface{})
	common := make([]zap.Field, 0, (keyLen/2)+1)
	for i := 0; i < keyLen; i += 2 {
		if msgKey, ok := keyvals[i].(string); ok && msgKey == log.DefaultMessageKey {
			msg = fmt.Sprint(keyvals[i+1])
			continue
		}
		key := fmt.Sprint(keyvals[i])
		if _, exists := types.SkipLog[key]; exists {
			common = append(common, zap.Any(fmt.Sprint(keyvals[i]), keyvals[i+1]))
			continue
		}
		//data = append(data, zap.Any(fmt.Sprint(keyvals[i]), keyvals[i+1]))
		data[fmt.Sprint(keyvals[i])] = keyvals[i+1]
	}
	// var errorTip string
	// fileLine := tools.FileWithLineNumToStr()
	// msg 为错误标题，fileLine 为文件所在行数，error为错误信息集合
	// 访问的接口信息等需要再context里携带，可以在main的log初始化的时候加入func，在http和grpc的拦截器加入自定义的拦截器把要的数据写入到context
	var dataByte = make([]byte, 0)
	if len(data) > 0 {
		dataStr, err := json.Marshal(data)
		if err != nil {
			dataByte = []byte(fmt.Sprintf("%+v", data))
		}
		dataByte = dataStr
	}

	common = append(common, zap.String("data", string(dataByte)))
	switch level {
	case log.LevelDebug:
		l.log.Debug(msg, common...)
	case log.LevelInfo:
		common = append(common, zap.String(types.LogStackTrace, ""))
		l.log.Info(msg, common...)
	case log.LevelWarn:
		l.log.Warn(msg, common...)
	case log.LevelError:
		l.log.Error(msg, common...)
		// errorTip = fmt.Sprint(data)
	case log.LevelFatal:
		// 防止误使用Fatal，以免程序down掉
		l.log.Error(msg, common...)
		//l.log.Fatal(msg, data...)
	}
	// 报警信息接口
	return nil
}

func (l *ZapLog) Sync() error {
	return l.log.Sync()
}

func (l *ZapLog) Close() error {
	return l.Sync()
}

func GrpcPath() log.Valuer {
	return func(ctx context.Context) interface{} {
		if info, ok := transport.FromServerContext(ctx); ok {
			return info.Operation()
		}
		return ""
	}
}

func HTTPPath() log.Valuer {
	return func(ctx context.Context) interface{} {
		if info, ok := transport.FromServerContext(ctx); ok {
			kind := info.Kind().String()
			if kind == transport.KindHTTP.String() {
				if h, ok := info.(*http.Transport); ok {
					return h.Request().URL.Path
				}
			}
		}
		return ""
	}
}

func (h *Helper) Debug(msgTag string, data ...interface{}) {
	//logData := make([]interface{},len(data) + 2)
	//logData = data[:]
	data = append(data, log.DefaultMessageKey, msgTag)
	_ = h.log.Log(log.LevelDebug, data...)
}

func (h *Helper) Info(msgTag string, data ...interface{}) {
	data = append(data, log.DefaultMessageKey, msgTag)
	_ = h.log.Log(log.LevelInfo, data...)
}

func (h *Helper) Warn(msgTag string, data ...interface{}) {
	data = append(data, log.DefaultMessageKey, msgTag)
	_ = h.log.Log(log.LevelWarn, data...)
}

func (h *Helper) Error(msgTag string, data ...interface{}) {
	data = append(data, log.DefaultMessageKey, msgTag)
	_ = h.log.Log(log.LevelError, data...)
}

// GetLog 获取log,打印输出
func (h *Helper) GetLog() log.Logger {
	return h.log
}

// WithContext 绑定context携带trace信息
func (h *Helper) WithContext(ctx context.Context) *Helper {
	return &Helper{log: log.WithContext(ctx, h.log)}
}
