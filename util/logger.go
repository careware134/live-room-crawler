package util

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
	"gopkg.in/robfig/cron.v2"
	"io"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
)

func InitLog() {
	// Create a lumberjack logger that writes to a file
	logPath := "./log/app.log"
	home, _ := user.Current()
	if home != nil {
		logPath = filepath.Join(home.HomeDir, "Documents", "STDH2D", "crawler", "app.log")
	}

	logger := &lumberjack.Logger{
		Filename:   logPath,
		MaxSize:    10, // Max size in MB before rolling over
		MaxBackups: 5,  // Max number of old log files to keep
		MaxAge:     30, // Max number of days to keep old log files
		Compress:   true,
	}
	multiWriter := io.MultiWriter(os.Stderr, logger)

	log.SetFormatter(&log.TextFormatter{
		ForceColors:      true,
		FullTimestamp:    true,
		TimestampFormat:  time.RFC3339,
		CallerPrettyfier: myCallerPrettyfier,
	})
	log.SetLevel(log.InfoLevel)
	log.SetOutput(multiWriter)

	// Create a cron job that runs every day at midnight
	c := cron.New()
	c.AddFunc("@midnight", func() {
		// Close the current log file and create a new one
		logger.Rotate()
	})
	c.Start()
}

func myCallerPrettyfier(frame *runtime.Frame) (function string, file string) {
	// Get the filename without the full path
	fileParts := strings.Split(frame.File, "/")
	filename := fileParts[len(fileParts)-1]

	// Return the function name and filename with line number
	return frame.Function, fmt.Sprintf("%s:%d", filename, frame.Line)
}

func Logger() *log.Entry {
	_, file, line, ok := runtime.Caller(1)
	if !ok {
		panic("Could not get context info for logger!")
	}

	filename := file[strings.LastIndex(file, "/")+1:] + ":" + strconv.Itoa(line)
	//functionName := runtime.FuncForPC(pc).Name()
	//fn := functionName[strings.LastIndex(functionName, ".")+1:]
	return log.WithField("file", filename)
}
