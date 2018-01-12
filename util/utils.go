package util

import (
	"log"
	"math"
	"os"
)

type LogLevel int

const (
	LOG_INFO LogLevel = iota + 1
	LOG_WARNING
	LOG_DEBUG
)

type Logger struct {
	logLevel LogLevel
}

func NewLogger() *Logger {
	logLevel := getLogLevelDefault()
	logger := Logger{
		logLevel: logLevel,
	}

	return &logger
}

func CalcSampleSize(total int) float64 {

	var n float64

	// Standard Deviation
	var o float64 = 0.5

	// Level of trustworthiness 99% high value 2.58 and 95% min value 1.96
	var z float64 = 2.58

	// Limit error acceptable from 1% to 9% . 5% is value standard
	var e float64 = 0.5

	var N float64 = float64(total)

	// Formula to calc representative sample
	n = (math.Pow(z, 2) * math.Pow(o, 2) * N) / (math.Pow(e, 2)*(N-1) + math.Pow(z, 2)*math.Pow(o, 2))

	return n * 100
}

func (l *Logger) Info(v ...interface{}) {
	if l.logLevel == LOG_DEBUG || l.logLevel == LOG_WARNING || l.logLevel == LOG_INFO {
		log.Println("[INFO]", v)
	}
}

func (l *Logger) Debug(v ...interface{}) {
	if l.logLevel == LOG_DEBUG {
		log.Println("[DEBUG]", v)
	}
}

func (l *Logger) Warning(v ...interface{}) {
	if l.logLevel == LOG_DEBUG || l.logLevel == LOG_WARNING {
		log.Println("[WARNING]", v)
	}
}

func getLogLevelDefault() LogLevel {
	var logLevelDefault LogLevel

	logLevelEnv := os.Getenv("LOG_LEVEL")

	if len(logLevelEnv) == 0 {
		logLevelDefault = LOG_INFO
	} else {
		switch logLevelEnv {
		case "INFO":
			logLevelDefault = LOG_INFO
		case "DEBUG":
			logLevelDefault = LOG_DEBUG
		case "WARNING":
			logLevelDefault = LOG_WARNING
		}
	}

	return logLevelDefault
}
