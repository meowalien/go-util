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
)
const (
	MUTE = uint8(0)
	ALL  = ^MUTE
)

var LogLevelMask uint8

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
	log.Logger
	Color go_util.ColorCode
}

func (l *Logger) PrettyPrintln(v ...interface{}) {
	l.Println(pretty.Sprint(v...))
}

func (l *Logger) SetColor(color go_util.ColorCode) {
	l.Color = color
}

// Printf will color the text and act like log.Printf()
func (l *Logger) Printf(format string, v ...interface{}) {
	_ = l.Output(2, fmt.Sprintf(go_util.ColorSting(format, l.Color), v...))
}

// Print will color the text and act like log.Print()
func (l *Logger) Print(v ...interface{}) {
	_ = l.Output(2, go_util.ColorSting(fmt.Sprint(v...), l.Color))
}

// Println will color the text and act like log.Println()
func (l *Logger) Println(v ...interface{}) {
	_ = l.Output(2, go_util.ColorSting(fmt.Sprint(v...), l.Color))
}

func (l *Logger) Output(calldepth int, s string) error {
	return l.Output(calldepth+1, go_util.ColorSting(s, l.Color))
}
// The TempLog should only use for debug, it will be close if the TempLogOpen parameter is false
// Se the settings in config/log.config.json
func (l *LoggerWrapper) TempLog() *Logger {
	if l.tempLogger == nil {
		if !TempLogOpen {
			return CreateMuteLogger()
		}
		l.tempLogger = &Logger{*log.New(os.Stdout, go_util.ColorSting("TEMP_LOG: ", TempColor), log.Ltime|log.Ldate|log.Lshortfile|log.Lmsgprefix), TempColor}
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
	if LogLevelMask != MUTE {
		fmt.Printf("Cteate logger: %s\n", prefix)
		return &LoggerWrapper{
			ERROR:   CreateErrorLogger(prefix,logFilePath+"error.log"),
			WARNING: CreateWaringLogger(prefix,logFilePath+"waring.log"),
			DEBUG:   CreateDebugLogger(prefix,logFilePath+"debug.log"),
			INFO:   CreateInfoLogger(prefix,logFilePath+"info.log"),
		}
	} else {
		return NewMuteLoggerWrapper()
	}
}

// NewMuteLoggerWrapper create a mute logger that will do nothing when use
func NewMuteLoggerWrapper() *LoggerWrapper {
	return &LoggerWrapper{
		ERROR:   CreateMuteLogger(),
		WARNING: CreateMuteLogger(),
		DEBUG:   CreateMuteLogger(),
	}

}

// CreateMuteLogger create a Mute Logger, the mute logger will do nothing when used.
func CreateMuteLogger() *Logger {
	return &Logger{*log.New(io.Discard, "", log.LstdFlags), go_util.FgBlack}
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
	return CreateLogger("DEBUG: "+prefix,logfile,InfoColor)
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
	return &Logger{*log.New(writer, fmt.Sprintf(go_util.ColorSting("%s: ", color), prefix), log.Ltime|log.Ldate|log.Lshortfile|log.Lmsgprefix), color}
}
