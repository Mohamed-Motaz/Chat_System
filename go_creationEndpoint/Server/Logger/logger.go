package Logger

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

const (
	SERVER = iota
	DATABASE
	MESSAGE_Q
	ELASTIC
)

const (
	LOG_INFO = iota
	LOG_ERROR
	LOG_DEBUG
)

const (
	ESSENTIAL     = 1
	NON_ESSENTIAL = 2
)

var (
	logger *log.Logger
)

var LOG_FILE_NAME string = fmt.Sprintf("%v-log.txt", time.Now().Unix())

func InitLogger(logFolder string) {
	path := filepath.Join(logFolder, LOG_FILE_NAME)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	logger = log.New(f, "", 0)
}

func print(essential int, format string, a ...interface{}) {
	fmt.Printf(format, a...)
	logger.Printf(format, a...)
}

func LogInfo(role int, essential int, format string, a ...interface{}) {
	format = beautifyLogs(role, essential, format, LOG_INFO)
	print(essential, format, a...)
}

func LogError(role int, essential int, format string, a ...interface{}) {
	format = beautifyLogs(role, essential, format, LOG_ERROR)
	print(essential, format, a...)
}

func LogDebug(role int, essential int, format string, a ...interface{}) {
	format = beautifyLogs(role, essential, format, LOG_DEBUG)
	print(essential, format, a...)
}

func FailOnError(role int, essential int, format string, a ...interface{}) {
	format = beautifyLogs(role, essential, format, LOG_ERROR)
	print(essential, format, a...)
	os.Exit(1)
}

func beautifyLogs(role int, essential int, format string, logType int) string {
	additionalInfo := determineRole(role)

	switch logType {
	case LOG_INFO:
		additionalInfo = Green + additionalInfo + "INFO: "
	case LOG_ERROR:
		additionalInfo = Red + additionalInfo + "ERROR: "
	case LOG_DEBUG:
		additionalInfo = Purple + additionalInfo + "DEBUG: "
	default:
		additionalInfo = Blue + additionalInfo + "DEFAULT: "
	}

	additionalInfo += makeTimestamp() + " -> "

	if format[len(format)-1] != '\n' {
		format += "\n"
	}
	format += Reset + "\n" //reset the terminal color

	return additionalInfo + format
}

func makeTimestamp() string {
	return time.Now().Format("01-02-2006 15:04:05")
}

func determineRole(role int) string {

	switch role {
	case MESSAGE_Q:
		return "MESSAGE_Q-> "
	case DATABASE:
		return "DATABASE-> "
	case SERVER:
		return "SERVER-> "
	case ELASTIC:
		return "ELASTIC-> "
	default:
		return "UNKNOWN-> "
	}
}
