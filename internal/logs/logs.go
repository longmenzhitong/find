// package logs implements methods for recording program running events.
package logs

import (
	"find/internal/config"
	"fmt"
	"io"
	"log"
	"os"
)

// enabled is the switch of log.
var enabled bool

// levelCode describes log level.
var levelCode int

const (
	levelCodeDebug = 0
	levelCodeInfo  = 1
	levelCodeWarn  = 2
	levelCodeError = 3
)

func init() {
	enabled = config.Conf.Log.Enabled
	if !enabled {
		return
	}

	err := setLog()
	if err != nil {
		fmt.Printf("set log error: %s\n", err.Error())
		return
	}

	levelCode, err = getLevelCode()
	if err != nil {
		fmt.Printf("get log level code error: %s\n", err.Error())
		return
	}
}

// setLog is used to set golang log options.
func setLog() error {
	file, err := os.OpenFile(config.Conf.Log.Path, os.O_CREATE|os.O_APPEND|os.O_RDWR, os.ModePerm)
	if err != nil {
		return fmt.Errorf("open log error: %v", err)
	}

	log.SetOutput(io.MultiWriter(file))
	log.SetFlags(log.LstdFlags)

	return nil
}

// getLevelCode is used to get level code(which is convenient to compare) by level.
func getLevelCode() (int, error) {
	switch config.Conf.Log.Level {
	case "debug":
		return levelCodeDebug, nil
	case "info":
		return levelCodeInfo, nil
	case "warn":
		return levelCodeWarn, nil
	case "error":
		return levelCodeError, nil
	}
	return -1, fmt.Errorf("invalid log level: %s", config.Conf.Log.Level)
}

// Debug is used to record log of debug level.
func Debug(format string, v ...interface{}) {
	if enabled && levelCode <= levelCodeDebug {
		log.Printf("[debug] "+format, v...)
	}
}

// Info is used to record log of info level.
func Info(format string, v ...interface{}) {
	if enabled && levelCode <= levelCodeInfo {
		log.Printf("[info] "+format, v...)
	}
}

// Warn is used to record log of warn level.
func Warn(format string, v ...interface{}) {
	if enabled && levelCode <= levelCodeWarn {
		log.Printf("[warn] "+format, v...)
	}
}

// Error is used to record log of error level.
// Only error log uses standard output which means it'll directly show to the user.
func Error(format string, v ...interface{}) {
	if enabled && levelCode <= levelCodeError {
		log.Printf("[error] "+format, v...)
	}
	fmt.Printf(format, v...)
}
