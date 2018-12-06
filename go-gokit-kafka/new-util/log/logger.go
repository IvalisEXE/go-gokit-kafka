package log

import (
	"log"
	"os"

	logkit "github.com/go-kit/kit/log"
	"gopkg.in/natefinch/lumberjack.v2"
)

const (
	//LogTime is log key for timestamp
	LogTime = "ts"
	//LogCaller is log key for source file name
	LogCaller = "caller"
	//LogMethod is log key for method name
	LogMethod = "method"
	//LogUser is log key for user
	LogUser = "user"
	//LogEmail is log key for email
	LogEmail = "email"
	//LogMobile is log key for mobile no
	LogMobile = "mobile"
	//LogRole is log key for role
	LogRole = "role"
	//LogTook is log key for call duration
	LogTook = "took"
	//LogInfo is log key for info
	LogInfo = "info"
	//LogError is log key for error
	LogError = "error"
	//LogService is log key for service name
	LogService = "service"
	//LogToken is log key for token
	LogToken = "token"
	//LogExit is log key for exit
	LogExit = "exit"
	//default file logger
	logFile = "service.log"
)

//File set default log to file
func File(file string) {
	logFile := &lumberjack.Logger{
		Filename:  file,
		MaxSize:   1, // megabytes
		LocalTime: true,
		Compress:  true, // disabled by default
	}

	log.SetOutput(logFile)
	log.SetFlags(log.LstdFlags)
}

//Logger returns default logger
func Logger() logkit.Logger {
	File(logFile)
	logger := logkit.NewLogfmtLogger(NewDefaultLogWriter())
	logger = logkit.With(logger, LogCaller, logkit.DefaultCaller)

	return logger
}

//StdLogger returns logger to stderr
func StdLogger() logkit.Logger {
	logger := logkit.NewLogfmtLogger(os.Stderr)
	logger = logkit.With(logger, LogTime, logkit.DefaultTimestampUTC, LogCaller, logkit.DefaultCaller)

	return logger
}

//FileLogger returns file logger
func FileLogger(file string) logkit.Logger {
	File(file)
	logger := logkit.NewLogfmtLogger(NewDefaultLogWriter())
	logger = logkit.With(logger, LogCaller, logkit.DefaultCaller)

	return logger
}
