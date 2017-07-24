package octlog

import (
	"fmt"
	"log"
	"octlink/mirage/src/utils"
	"os"
)

const (
	PANIC_LEVEL int = iota
	FATAL_LEVEL
	ERROR_LEVEL
	WARN_LEVEL
	INFO_LEVEL
	DEBUG_LEVEL
)

var GDebugConfig DebugConfig

type DebugConfig struct {
	level int
}

type LogConfig struct {
	level   int
	logTime int64
	LogFile string
	fileFd  *os.File
	logger  *log.Logger
}

func getLogDir() string {
	if utils.IsFileExist("logs") {
		return "./logs/"
	} else {
		return "./"
	}
}

func InitLogConfig(logFile string, level int) *LogConfig {
	config := new(LogConfig)

	config.LogFile = getLogDir() + logFile
	config.level = level

	logfile, err := os.OpenFile(config.LogFile, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		fmt.Printf("%s\r\n", err.Error())
		return nil
	}

	logfile.Seek(0, 2)

	config.logger = log.New(logfile, "", log.Ldate|log.Ltime|log.Lshortfile)

	return config
}

func InitDebugConfig(level int) {
	GDebugConfig.level = level
	log.SetFlags(log.Lmicroseconds | log.Lshortfile)
}

// For File Logging
func (config *LogConfig) Debugf(format string, args ...interface{}) {
	if config.level >= DEBUG_LEVEL {
		config.logger.SetPrefix("DEBUG ")
		config.logger.Printf(format, args...)
	}
}

func (config *LogConfig) Infof(format string, args ...interface{}) {
	if config.level >= INFO_LEVEL {
		config.logger.SetPrefix("INFO ")
		config.logger.Printf(format, args...)
	}
}

func (config *LogConfig) Warnf(format string, args ...interface{}) {
	if config.level >= WARN_LEVEL {
		config.logger.SetPrefix("WARN ")
		config.logger.Printf(format, args...)
	}
}

func (config *LogConfig) Errorf(format string, args ...interface{}) {
	if config.level >= ERROR_LEVEL {
		config.logger.SetPrefix("ERROR ")
		config.logger.Printf(format, args...)
	}
}

func (config *LogConfig) Fatalf(format string, args ...interface{}) {
	if config.level >= FATAL_LEVEL {
		config.logger.SetPrefix("FATAL ")
		config.logger.Printf(format, args...)
	}
}

func (config *LogConfig) Panicf(format string, args ...interface{}) {
	if config.level >= PANIC_LEVEL {
		config.logger.SetPrefix("PANIC ")
		config.logger.Printf(format, args...)
	}
}

// For Debuging
func Debug(format string, args ...interface{}) {
	if GDebugConfig.level >= DEBUG_LEVEL {
		log.SetPrefix("DEBUG ")
		log.Output(2, fmt.Sprintf(format, args...))
	}
}

func Info(format string, args ...interface{}) {
	if GDebugConfig.level >= INFO_LEVEL {
		log.SetPrefix("INFO ")
		log.Output(2, fmt.Sprintf(format, args...))
	}
}

func Warn(format string, args ...interface{}) {
	if GDebugConfig.level >= WARN_LEVEL {
		log.SetPrefix("WARN ")
		log.Output(2, fmt.Sprintf(format, args...))
	}
}

func Error(format string, args ...interface{}) {
	if GDebugConfig.level >= ERROR_LEVEL {
		log.SetPrefix("ERROR ")
		log.Output(2, fmt.Sprintf(format, args...))
	}
}

func Fatal(format string, args ...interface{}) {
	if GDebugConfig.level >= FATAL_LEVEL {
		log.SetPrefix("FATAL ")
		log.Output(2, fmt.Sprintf(format, args...))
	}
}

func Panic(format string, args ...interface{}) {
	if GDebugConfig.level >= PANIC_LEVEL {
		log.SetPrefix("PANIC ")
		log.Output(2, fmt.Sprintf(format, args...))
	}
}
