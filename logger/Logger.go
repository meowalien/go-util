package logger

import (
	"fmt"
	"github.com/kr/pretty"
	"github.com/meowalien/go-util"
	"io"
	"log"
	"os"
)

const (
	DEBUG uint8 = 1 << iota
	WARNING
	ERROR
	INFO
)
const (
	MUTE = uint8(0)
	ALL  = ^MUTE
)

var LogLevelMask = ALL

var InfoColor = go_util.FgGreen
var DebugColor = go_util.FgBlue
var WaringColor = go_util.FgYellow
var ErrorColor = go_util.FgRed
var TempColor = go_util.FgCyan
var TempLogOpen = true



type LoggerWrapper struct {
	// Debug logger should be used to log a message that will be useful for develop.
	DEBUG *Logger
	// Waring logger should be used to log an error that will not break business logic.
	WARNING *Logger
	// error logger should be used to log an error that will break business logic.
	ERROR      *Logger
	// error logger should be used to log an error that will break business logic.
	INFO      *Logger
	tempLogger *Logger
}

type Logger struct {
	IsMute bool
	log.Logger
	Color go_util.ColorCode
}

func (l *Logger) PrettyPrintln(v ...interface{}) {
	l.Println(pretty.Sprint(v...))
}

func (l *Logger) SetColor(color go_util.ColorCode) {
	l.Color = color
}

func (l *Logger) Output(calldepth int, s string) error {
	if l.IsMute {
		return nil
	}
	return l.Output(calldepth+1, go_util.ColorSting(s, l.Color))
}


// The TempLog should only use for debug, it will be close if the TempLogOpen parameter is false
// Se the settings in config/log.config.json
func (l *LoggerWrapper) TempLog() *Logger {
	if l.tempLogger == nil {
		if !TempLogOpen {
			return CreateMuteLogger()
		}
		l.tempLogger = &Logger{false,*log.New(os.Stdout, go_util.ColorSting("TEMP_LOG: ", TempColor), log.Ltime|log.Ldate|log.Lshortfile|log.Lmsgprefix), TempColor}
	}
	return l.tempLogger
}



type Setting struct {
	LogLevelMask uint8
	DebugColor   go_util.ColorCode
	WaringColor  go_util.ColorCode
	ErrorColor   go_util.ColorCode
	TempColor    go_util.ColorCode
	TempLogOpen  bool
}


// NewLoggerWrapper Create a new LoggerWrapper with given prefix,
// The prefix will be print before all log rows
func NewLoggerWrapper(prefix ,logFilePath string) *LoggerWrapper {
	var errorLogger , waringLogger , debugLogger , infoLogger *Logger

	if LogLevelMask != MUTE {
		fmt.Printf("Cteate logger: %s\n", prefix)
		if LogLevelMask & ERROR != 0{
			errorLogger = CreateErrorLogger(prefix,logFilePath+"error.log")
		}else{
			errorLogger =CreateMuteLogger()
		}

		if LogLevelMask & WARNING != 0{
			waringLogger = CreateWaringLogger(prefix,logFilePath+"waring.log")
		}else{
			waringLogger =CreateMuteLogger()
		}

		if LogLevelMask & DEBUG != 0{
			debugLogger =  CreateDebugLogger(prefix,logFilePath+"debug.log")
		}else{
			debugLogger =CreateMuteLogger()
		}

		if LogLevelMask & INFO  != 0{
			infoLogger = CreateInfoLogger(prefix,logFilePath+"info.log")
		}else{
			infoLogger =CreateMuteLogger()
		}

		return &LoggerWrapper{
			ERROR:   errorLogger,
			WARNING: waringLogger,
			DEBUG:   debugLogger,
			INFO:   infoLogger,
		}
	} else {
		fmt.Printf("Logger muted: %s\n",prefix)
		return NewMuteLoggerWrapper()
	}
}

// NewMuteLoggerWrapper create a mute logger that will do nothing when use
func NewMuteLoggerWrapper() *LoggerWrapper {
	return &LoggerWrapper{
		ERROR:   CreateMuteLogger(),
		WARNING: CreateMuteLogger(),
		DEBUG:   CreateMuteLogger(),
		INFO:   CreateMuteLogger(),
	}

}

// CreateMuteLogger create a Mute Logger, the mute logger will do nothing when used.
func CreateMuteLogger() *Logger {
	return &Logger{true ,*log.Default(),go_util.FgBlack}
}

// CreateErrorLogger create an Error Logger.
// error logger should be used to log an error that will break business logic.
func CreateErrorLogger(prefix string,logfile string) *Logger {
	return CreateLogger("ERROR: "+prefix,logfile,ErrorColor)
}

// CreateWaringLogger create a Waring Logger.
// Waring logger should be used to log an error that will not break business logic.
func CreateWaringLogger(prefix string,logfile string) *Logger {
	return CreateLogger("WARING: "+prefix,logfile,WaringColor)
}

// CreateDebugLogger create a Waring Logger.
// Debug logger should be used to log a message needed for debug.
func CreateDebugLogger(prefix string,logfile string) *Logger {
	return CreateLogger("DEBUG: "+prefix,logfile,DebugColor)
}

// CreateInfoLogger create an Info Logger.
// Info logger should be used to log a message that will be useful for develop.
func CreateInfoLogger(prefix string,logfile string) *Logger {
	return CreateLogger("INFO: "+prefix,logfile,InfoColor)
}



func CreateLogger(prefix string,logfile string , color go_util.ColorCode) *Logger {
	var writer io.Writer
	if logfile != ""{
		outputFile, err := os.OpenFile(logfile , os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			panic(err.Error())
		}
		writer = io.MultiWriter(outputFile, os.Stdout)
	}else{
		writer = io.Discard
	}
	return &Logger{false,*log.New(writer, fmt.Sprintf(go_util.ColorSting("%s: ", color), prefix), log.Ltime|log.Ldate|log.Lshortfile|log.Lmsgprefix), color}
}
