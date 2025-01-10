package logger

import (
	"fmt"
	"github.com/SongZihuan/huan-proxy/src/config"
	"github.com/SongZihuan/huan-proxy/src/utils"
	"io"
	"os"
)

type LoggerLevel string

const (
	LevelDebug LoggerLevel = "debug"
	LevelInfo  LoggerLevel = "info"
	LevelWarn  LoggerLevel = "warn"
	LevelError LoggerLevel = "error"
	LevelPanic LoggerLevel = "panic"
	LevelNone  LoggerLevel = "none"
)

type loggerLevel int64

const (
	levelDebug loggerLevel = 1
	levelInfo  loggerLevel = 2
	levelWarn  loggerLevel = 3
	levelError loggerLevel = 4
	levelPanic loggerLevel = 5
	levelNone  loggerLevel = 6
)

var levelMap = map[LoggerLevel]loggerLevel{
	LevelDebug: levelDebug,
	LevelInfo:  levelInfo,
	LevelWarn:  levelWarn,
	LevelError: levelError,
	LevelPanic: levelPanic,
	LevelNone:  levelNone,
}

type Logger struct {
	level      LoggerLevel
	logLevel   loggerLevel
	logTag     bool
	warnWriter io.Writer
	errWriter  io.Writer
	args0      string
	args0Name  string
}

var globalLogger *Logger = nil

func InitLogger() error {
	if !config.IsReady() {
		panic("config is not ready")
	}

	level := LoggerLevel(config.Config().Yaml.GlobalConfig.LogLevel)
	logLevel, ok := levelMap[level]
	if !ok {
		return fmt.Errorf("invalid log level: %s", level)
	}

	logger := &Logger{
		level:      level,
		logLevel:   logLevel,
		logTag:     config.Config().Yaml.LogTag.ToBool(true),
		warnWriter: os.Stdout,
		errWriter:  os.Stderr,
		args0:      utils.GetArgs0(),
		args0Name:  utils.GetArgs0Name(),
	}

	globalLogger = logger
	return nil
}

func IsReady() bool {
	return globalLogger != nil
}

func (l *Logger) Executablef(format string, args ...interface{}) string {
	str := fmt.Sprintf(format, args...)
	if str == "" {
		_, _ = fmt.Fprintf(l.warnWriter, "executable: %s\n", l.args0)
	} else {
		_, _ = fmt.Fprintf(l.warnWriter, "executable[%s]: %s\n", str, l.args0)
	}
	return l.args0
}

func (l *Logger) Executable() string {
	return l.Executablef("")
}

func (l *Logger) Tagf(format string, args ...interface{}) {
	l.TagSkipf(1, format, args...)
}

func (l *Logger) TagSkipf(skip int, format string, args ...interface{}) {
	if !l.logTag {
		return
	}

	funcName, file, _, line := utils.GetCallingFunctionInfo(skip + 1)

	str := fmt.Sprintf(format, args...)
	_, _ = fmt.Fprintf(l.warnWriter, "%s: TAG %s %s %s:%d\n", l.args0Name, str, funcName, file, line)
}

func (l *Logger) Debugf(format string, args ...interface{}) {
	if l.logLevel > levelDebug {
		return
	}

	str := fmt.Sprintf(format, args...)
	_, _ = fmt.Fprintf(l.warnWriter, "%s: %s\n", l.args0Name, str)
}

func (l *Logger) Infof(format string, args ...interface{}) {
	if l.logLevel > levelInfo {
		return
	}

	str := fmt.Sprintf(format, args...)
	_, _ = fmt.Fprintf(l.warnWriter, "%s: %s\n", l.args0Name, str)
}

func (l *Logger) Warnf(format string, args ...interface{}) {
	if l.logLevel > levelWarn {
		return
	}

	str := fmt.Sprintf(format, args...)
	_, _ = fmt.Fprintf(l.warnWriter, "%s: %s\n", l.args0Name, str)
}

func (l *Logger) Errorf(format string, args ...interface{}) {
	if l.logLevel > levelError {
		return
	}

	str := fmt.Sprintf(format, args...)
	_, _ = fmt.Fprintf(l.errWriter, "%s: %s\n", l.args0Name, str)
}

func (l *Logger) Panicf(format string, args ...interface{}) {
	if l.logLevel > levelPanic {
		return
	}

	str := fmt.Sprintf(format, args...)
	_, _ = fmt.Fprintf(l.errWriter, "%s: %s\n", l.args0Name, str)
}

func (l *Logger) Tag(args ...interface{}) {
	l.TagSkip(1, args...)
}

func (l *Logger) TagSkip(skip int, args ...interface{}) {
	if !l.logTag {
		return
	}

	funcName, file, _, line := utils.GetCallingFunctionInfo(skip + 1)

	str := fmt.Sprint(args...)
	_, _ = fmt.Fprintf(l.warnWriter, "%s: TAG %s %s %s:%d\n", l.args0Name, str, funcName, file, line)
}

func (l *Logger) Debug(args ...interface{}) {
	if l.logLevel > levelDebug {
		return
	}

	str := fmt.Sprint(args...)
	_, _ = fmt.Fprintf(l.warnWriter, "%s: %s\n", l.args0Name, str)
}

func (l *Logger) Info(args ...interface{}) {
	if l.logLevel > levelInfo {
		return
	}

	str := fmt.Sprint(args...)
	_, _ = fmt.Fprintf(l.warnWriter, "%s: %s\n", l.args0Name, str)
}

func (l *Logger) Warn(args ...interface{}) {
	if l.logLevel > levelWarn {
		return
	}

	str := fmt.Sprint(args...)
	_, _ = fmt.Fprintf(l.warnWriter, "%s: %\ns", l.args0Name, str)
}

func (l *Logger) Error(args ...interface{}) {
	if l.logLevel > levelError {
		return
	}

	str := fmt.Sprint(args...)
	_, _ = fmt.Fprintf(l.errWriter, "%s: %s\n", l.args0Name, str)
}

func (l *Logger) Panic(args ...interface{}) {
	if l.logLevel > levelPanic {
		return
	}

	str := fmt.Sprint(args...)
	_, _ = fmt.Fprintf(l.errWriter, "%s: %s\n", l.args0Name, str)
}
