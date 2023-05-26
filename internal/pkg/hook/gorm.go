package hook

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
	"time"
)

type SQLLog struct {
	level         logger.LogLevel
	SlowThreshold time.Duration
	logout        log.Logger
}

// NewSQLLog 需要返回 gorm.logger.Interface 的接口类
func NewSQLLog(logout log.Logger) *SQLLog {
	return &SQLLog{
		logout:        logout,
		SlowThreshold: 2 * time.Second,
	}
}
func (s *SQLLog) LogMode(level logger.LogLevel) logger.Interface {
	s.level = level
	return s
}

func (s *SQLLog) Info(ctx context.Context, msg string, data ...interface{}) {
	_ = log.WithContext(ctx, s.logout).Log(log.LevelInfo, log.DefaultMessageKey, "sql-info", "data", data, "tip", msg, "file", utils.FileWithLineNum())
}

func (s *SQLLog) Warn(ctx context.Context, msg string, data ...interface{}) {
	_ = log.WithContext(ctx, s.logout).Log(log.LevelWarn, log.DefaultMessageKey, "sql-warn", "data", data, "tip", msg, "file", utils.FileWithLineNum())
}

func (s *SQLLog) Error(ctx context.Context, msg string, data ...interface{}) {
	_ = log.WithContext(ctx, s.logout).Log(log.LevelError, log.DefaultMessageKey, "sql-error", "data", data, "tip", msg, "file", utils.FileWithLineNum())
}

// Trace 捕获trace信息
func (s *SQLLog) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	elapsed := time.Since(begin)
	switch {
	case err != nil && s.level >= logger.Error && (!errors.Is(err, gorm.ErrRecordNotFound)):
		sql, rows := fc()
		if rows == -1 {
			_ = log.WithContext(ctx, s.logout).Log(log.LevelError, log.DefaultMessageKey, "sql-error", "file", utils.FileWithLineNum(), "err", err, "elapsed", float64(elapsed.Nanoseconds())/1e6, "rows", rows, "sql", sql)
		}
	case elapsed > s.SlowThreshold && s.SlowThreshold != 0 && s.level >= logger.Warn:
		sql, rows := fc()
		slowLog := fmt.Sprintf("SLOW SQL >= %v", s.SlowThreshold)
		_ = log.WithContext(ctx, s.logout).Log(log.LevelWarn, log.DefaultMessageKey, "sql-warn", "file", utils.FileWithLineNum(), "slowLog", slowLog, "elapsed", float64(elapsed.Nanoseconds())/1e6, "rows", rows, "sql", sql)
	case s.level == logger.Info:
		sql, rows := fc()
		_ = log.WithContext(ctx, s.logout).Log(log.LevelWarn, log.DefaultMessageKey, "sql-info", "file", utils.FileWithLineNum(), "elapsed", float64(elapsed.Nanoseconds())/1e6, "rows", rows, "sql", sql)
	}
}

type TraceOrmPlugin struct {
	name string
}

func (t *TraceOrmPlugin) Name() string {
	return t.name
}

func (t *TraceOrmPlugin) WithName(name string) {
	t.name = name
}

// 注册钩子函数，上报 trace信息
func (t *TraceOrmPlugin) Initialize(db *gorm.DB) error {
	_ = db.Callback().Query().Before("gorm:query").Register("query_before", beforeSQL("query_span"))
	_ = db.Callback().Query().After("gorm:query").Register("query_after", afterSQL("query_span"))
	_ = db.Callback().Raw().Before("gorm:raw").Register("raw_before", beforeSQL("raw_span"))
	_ = db.Callback().Raw().After("gorm:raw").Register("raw_after", afterSQL("raw_span"))
	_ = db.Callback().Update().Before("gorm:update").Register("update_before", beforeSQL("update_span"))
	_ = db.Callback().Update().After("gorm:update").Register("update_after", afterSQL("update_span"))
	_ = db.Callback().Delete().Before("gorm:delete").Register("Delete_before", beforeSQL("Delete_span"))
	_ = db.Callback().Delete().After("gorm:delete").Register("Delete_after", afterSQL("Delete_span"))
	_ = db.Callback().Create().Before("gorm:create").Register("Create_before", beforeSQL("Create_span"))
	_ = db.Callback().Create().After("gorm:create").Register("Create_after", afterSQL("Create_span"))
	_ = db.Callback().Row().Before("gorm:row").Register("Row_before", beforeSQL("Row_span"))
	_ = db.Callback().Row().After("gorm:row").Register("Row_after", afterSQL("Row_span"))
	return nil
}

func beforeSQL(spanTag string) func(db *gorm.DB) {
	return func(db *gorm.DB) {
		if db.Statement.Context != nil {
			sqlCtx := db.Statement.Context
			// 再有traceID的方式下才生成新的
			if span := trace.SpanContextFromContext(sqlCtx); span.HasTraceID() {
				tracer := otel.Tracer("mysql")
				_, span := tracer.Start(sqlCtx, "mysql-"+spanTag, trace.WithAttributes(), trace.WithSpanKind(trace.SpanKindServer))
				db.InstanceSet(spanTag, span)
			}
		}
	}
}

func afterSQL(spanTag string) func(db *gorm.DB) {
	return func(db *gorm.DB) {
		// todo 记录SQL
		sql := db.Dialector.Explain(db.Statement.SQL.String(), db.Statement.Vars...)
		span, ok := db.InstanceGet(spanTag)
		if span != nil && ok {
			if _, exits := span.(trace.Span); exits {
				span := span.(trace.Span)
				span.SetAttributes(attribute.String("query", sql))
				span.End()
			}
		}
	}
}
