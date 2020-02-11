package zerologger

import (
	"fmt"
	"runtime"
	"time"

	"github.com/rs/zerolog"
	lumber "gopkg.in/natefinch/lumberjack.v2"
)

//ZeroLogger -
type ZeroLogger struct {
	logger zerolog.Logger
}

//CreateLogger - create a logger
func CreateLogger(logfile string, maxSize, maxBackups, maxAge int, compress bool) ZeroLogger {

	fileLogger := &lumber.Logger{
		Filename:   logfile,
		MaxSize:    maxSize, // megabytes
		MaxBackups: maxBackups,
		MaxAge:     maxAge,   //days
		Compress:   compress, // disabled by default
	}

	zerolog.TimeFieldFormat = time.RFC3339
	loggerval := zerolog.New(fileLogger).With().Timestamp().Logger()
	return ZeroLogger{
		logger: loggerval,
	}

}

func log(logEvent *zerolog.Event, stackDepth int, uuid, format string, v ...interface{}) {
	_, file, line, ok := runtime.Caller(stackDepth + 1)
	if ok {
		logEvent.Str("file", file).Int("line", line).Str("UUID", uuid).Msg(fmt.Sprintf(format, v...))
	} else {
		logEvent.Str("file", "unknown").Str("UUID", uuid).Msg(fmt.Sprintf(format, v...))
	}
}

func logMap(logEvent *zerolog.Event, stackDepth int, uuid, message string, fields map[string]interface{}) {
	_, file, line, ok := runtime.Caller(stackDepth + 1)
	if ok {
		logEvent.Str("file", file).Int("line", line).Str("UUID", uuid).Fields(fields).Msg(message)
	} else {
		logEvent.Str("file", "unknown").Str("UUID", uuid).Fields(fields).Msg(message)
	}
}

//LogDebug - wrapper for debug log level
func (l *ZeroLogger) LogDebug(uuid, format string, v ...interface{}) {
	//runtime.Caller(1) returns the file name & line of the caller of this function
	log(l.logger.Debug(), 1, uuid, format, v)
}

//LogWarn - wrapper for warn log level
func (l *ZeroLogger) LogWarn(uuid, format string, v ...interface{}) {
	//runtime.Caller(1) returns the file name & line of the caller of this function
	log(l.logger.Warn(), 1, uuid, format, v)
}

//LogInfo - wrapper for info log level
func (l *ZeroLogger) LogInfo(uuid, format string, v ...interface{}) {
	//runtime.Caller(1) returns the file name & line of the caller of this function
	log(l.logger.Info(), 1, uuid, format, v)
}

//LogInfoMap - wrapper for info log level to log multiple fields (go map)
func (l *ZeroLogger) LogInfoMap(uuid, message string, fields map[string]interface{}) {
	//runtime.Caller(1) returns the file name & line of the caller of this function
	logMap(l.logger.Info(), 1, uuid, message, fields)
}

//LogWarnMap - wrapper for warn log level to log multiple fields (go map)
func (l *ZeroLogger) LogWarnMap(uuid, message string, fields map[string]interface{}) {
	//runtime.Caller(1) returns the file name & line of the caller of this function
	logMap(l.logger.Warn(), 1, uuid, message, fields)
}

//LogError - wrapper for error log level
func (l *ZeroLogger) LogError(uuid, format string, v ...interface{}) {
	//runtime.Caller(1) returns the file name & line of the caller of this function
	log(l.logger.Error(), 1, uuid, format, v)
}

//SD2LogError - same as LogError but gets caller information from stack depth level 2
func (l *ZeroLogger) SD2LogError(uuid, format string, v ...interface{}) {
	//runtime.Caller(1) returns the file name & line of the caller of this function
	log(l.logger.Error(), 2, uuid, format, v)
}

//LogFatal - wrapper for fatal log level.
//LogFatal will cause program to exit (os.Exit(1))
func (l *ZeroLogger) LogFatal(uuid, format string, v ...interface{}) {
	//runtime.Caller(1) returns the file name & line of the caller of this function
	log(l.logger.Fatal(), 1, uuid, format, v)
}
