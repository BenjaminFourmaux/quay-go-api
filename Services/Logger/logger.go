package Logger

import (
	"fmt"
	"strings"
	"sync"
	"time"
)

var introducer = " > "
var separator = " | "

type Level uint8

const (
	LevelDebug Level = iota
	LevelInfo
	LevelSuccess
	LevelWarning
	LevelError
	LevelSilent
)

var levelMu sync.RWMutex
var currentLevel = LevelDebug

func getCurrentDatetime() string {
	currentDatetime := time.Now()
	return currentDatetime.Format("15:04:05 02-01-2006")
}

func SetLevel(level Level) {
	levelMu.Lock()
	currentLevel = level
	levelMu.Unlock()
	fmt.Println("Log level set to: " + levelToString(level))
}

func GetLevel() Level {
	levelMu.RLock()
	level := currentLevel
	levelMu.RUnlock()

	return level
}

func shouldLog(level Level) bool {
	current := GetLevel()
	if current == LevelSilent {
		return false
	}

	return level >= current
}

func levelToString(level Level) string {
	switch level {
	case LevelDebug:
		return "DEBUG"
	case LevelInfo:
		return "INFO"
	case LevelSuccess:
		return "SUCCESS"
	case LevelWarning:
		return "WARNING"
	case LevelError:
		return "ERROR"
	default:
		return "LOG"
	}
}

func StringToLevel(level string) Level {
	switch strings.ToUpper(level) {
	case "SILENT":
		return LevelSilent
	case "DEBUG":
		return LevelDebug
	case "INFO":
		return LevelInfo
	case "SUCCESS":
		return LevelSuccess
	case "WARNING":
		return LevelWarning
	case "ERROR":
		return LevelError
	default:
		Error("Invalid log level: " + level)
		return LevelDebug
	}
}

func logAt(level Level, message string) {
	if !shouldLog(level) {
		return
	}

	fmt.Println(getCurrentDatetime() + introducer + levelToString(level) + separator + message)
}

func Debug(message string, args ...interface{}) {
	logAt(LevelDebug, fmt.Sprintf(message, args...))
}

func Info(message string) {
	logAt(LevelInfo, message)
}

func Success(message string, args ...interface{}) {
	logAt(LevelSuccess, fmt.Sprintf(message, args...))
}

func Warning(message string, args ...interface{}) {
	logAt(LevelWarning, fmt.Sprintf(message, args...))
}

func Error(message string, args ...interface{}) {
	logAt(LevelError, fmt.Sprintf(message, args...))
}

func Raise(err error) {
	Error(err.Error())
}

func Separator() {
	fmt.Println("----------------------------------")
}
