package util

import (
	"log"
	"math"
	"math/rand"
	"os"
	"time"
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

func CalcSampleSizeMethod1(total int) int {

	// Standard Deviation
	var o float64 = 0.5

	// Level of trustworthiness 99% high value 2.58 and 95% min value 1.96
	var z float64 = 2.58

	// Limit error acceptable from 1% to 9% . 5% is value standard
	var e float64 = 0.1

	// Total size
	var N float64 = float64(total)

	// Formula to calc representative sample
	n := (math.Pow(z, 2) * math.Pow(o, 2) * N) / (math.Pow(e, 2)*(N-1) + math.Pow(z, 2)*math.Pow(o, 2))

	return int(n)
}

func CalcSampleSizeMethod2(total int) int {

	// Security of 99%
	var Z float64 = 2.58

	// Proportion 50%
	var p float64 = 0.5

	var q float64 = 1 - p

	// Presition 5%
	var d float64 = 0.05

	// Total poblation
	var N float64 = float64(total)

	n := (N * math.Pow(Z, 2) * q * q) / (math.Pow(d, 2)*(N-1) + math.Pow(Z, 2)*p*q)

	return int(n)
}

func GetRandomNumberFrom(limit int) int {
	s := rand.NewSource(time.Now().UnixNano())
	r := rand.New(s)

	return r.Intn(100)
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
