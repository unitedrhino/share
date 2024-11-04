package stores

import (
	"context"
	"errors"
	"fmt"
	"gitee.com/unitedrhino/share/utils"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm/logger"
	"time"
)

var (
	SlowThreshold             time.Duration = 500 * time.Millisecond
	ParameterizedQueries      bool          = false
	IgnoreRecordNotFoundError bool          = true
)

type Log struct {
	LogLevel logger.LogLevel
}

func NewLog(LogLevel logger.LogLevel) *Log {
	return &Log{LogLevel: LogLevel}
}

func (l *Log) LogMode(level logger.LogLevel) logger.Interface {
	newlogger := *l
	newlogger.LogLevel = level
	return &newlogger
}

func (l *Log) Info(ctx context.Context, s string, i ...interface{}) {
	if l.LogLevel >= logger.Info {
		logx.WithContext(ctx).WithCallerSkip(1).Infow(utils.FileWithLineNum(), logx.Field("body", fmt.Sprintf(s, i...)))
	}
}

func (l *Log) Warn(ctx context.Context, s string, i ...interface{}) {
	if l.LogLevel >= logger.Warn {
		logx.WithContext(ctx).WithCallerSkip(1).Errorw(utils.FileWithLineNum(), logx.Field("body", fmt.Sprintf(s, i...)))
	}
}

func (l *Log) Error(ctx context.Context, s string, i ...interface{}) {
	if l.LogLevel >= logger.Error {
		logx.WithContext(ctx).WithCallerSkip(1).Errorw(utils.FileWithLineNum(), logx.Field("body", fmt.Sprintf(s, i...)))
	}
}

func (l *Log) ParamsFilter(ctx context.Context, sql string, params ...interface{}) (string, []interface{}) {
	if ParameterizedQueries {
		return sql, nil
	}
	return sql, params
}

func (l *Log) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	if l.LogLevel <= logger.Silent {
		return
	}
	elapsed := time.Since(begin)
	useTime := fmt.Sprintf("%v ms", float64(elapsed.Nanoseconds())/1e6)
	switch {
	case err != nil && l.LogLevel >= logger.Error && (!errors.Is(err, logger.ErrRecordNotFound) || !IgnoreRecordNotFoundError):
		sql, rows := fc()
		logx.WithContext(ctx).Errorw("errorSql", logx.Field("call", utils.FileWithLineNum()), logx.Field("sql", sql),
			logx.Field("useTime", useTime), logx.Field("rows", rows), logx.Field("err", err))
	case elapsed > SlowThreshold && SlowThreshold != 0 && l.LogLevel >= logger.Warn:
		sql, rows := fc()
		slowLog := fmt.Sprintf("slowSql >= %v", SlowThreshold)
		logx.WithContext(ctx).Sloww(slowLog, logx.Field("call", utils.FileWithLineNum()), logx.Field("sql", sql),
			logx.Field("useTime", useTime), logx.Field("rows", rows))
	case l.LogLevel == logger.Info:
		sql, rows := fc()
		logx.WithContext(ctx).Infow("traceSql", logx.Field("call", utils.FileWithLineNum()), logx.Field("sql", sql),
			logx.Field("useTime", useTime), logx.Field("rows", rows))
	}
}
