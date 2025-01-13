package logger

import (
	"io"
)

func Executablef(format string, args ...interface{}) string {
	if !IsReady() {
		return ""
	}
	return globalLogger.Executablef(format, args...)
}

func Executable() string {
	if !IsReady() {
		return ""
	}
	return globalLogger.Executable()
}

func Tagf(format string, args ...interface{}) {
	if !IsReady() {
		return
	}
	globalLogger.TagSkipf(1, format, args...)
}

func Debugf(format string, args ...interface{}) {
	if !IsReady() {
		return
	}
	globalLogger.Debugf(format, args...)
}

func Infof(format string, args ...interface{}) {
	if !IsReady() {
		return
	}
	globalLogger.Infof(format, args...)
}

func Warnf(format string, args ...interface{}) {
	if !IsReady() {
		return
	}
	globalLogger.Warnf(format, args...)
}

func Errorf(format string, args ...interface{}) {
	if !IsReady() {
		return
	}
	globalLogger.Errorf(format, args...)
}

func Panicf(format string, args ...interface{}) {
	if !IsReady() {
		return
	}
	globalLogger.Panicf(format, args...)
}

func Tag(args ...interface{}) {
	if !IsReady() {
		return
	}
	globalLogger.TagSkip(1, args...)
}

func Debug(args ...interface{}) {
	if !IsReady() {
		return
	}
	globalLogger.Debug(args...)
}

func Info(args ...interface{}) {
	if !IsReady() {
		return
	}
	globalLogger.Info(args...)
}

func Warn(args ...interface{}) {
	if !IsReady() {
		return
	}
	globalLogger.Warn(args...)
}

func Error(args ...interface{}) {
	if !IsReady() {
		return
	}
	globalLogger.Error(args...)
}

func Panic(args ...interface{}) {
	if !IsReady() {
		return
	}
	globalLogger.Panic(args...)
}

func DebugWriter() io.Writer {
	if !IsReady() {
		return DefaultWarnWriter
	}
	return globalLogger.DebugWriter()
}

func InfoWriter() io.Writer {
	if !IsReady() {
		return DefaultWarnWriter
	}
	return globalLogger.InfoWriter()
}

func WarningWriter() io.Writer {
	if !IsReady() {
		return DefaultWarnWriter
	}
	return globalLogger.WarningWriter()
}

func TagWriter() io.Writer {
	if !IsReady() {
		return DefaultWarnWriter
	}
	return globalLogger.TagWriter()
}

func ErrorWriter() io.Writer {
	if !IsReady() {
		return DefaultWarnWriter
	}
	return globalLogger.ErrorWriter()
}
func PanicWriter() io.Writer {
	if !IsReady() {
		return DefaultWarnWriter
	}
	return globalLogger.PanicWriter()
}
