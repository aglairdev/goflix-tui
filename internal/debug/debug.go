package debug

import (
	"fmt"
	"os"
	"sync"
	"time"
)

// Debug -d

var Enabled bool

var (
	mu      sync.Mutex
	pending []string
	logFile *os.File
)

func write(format string, args ...interface{}) {
	if !Enabled {
		return
	}
	msg := fmt.Sprintf("[goflix-debug] "+format, args...)
	fmt.Fprintln(os.Stderr, msg)
	mu.Lock()
	pending = append(pending, msg)
	mu.Unlock()
	if logFile != nil {
		fmt.Fprintln(logFile, msg)
	}
}

func Log(format string, args ...interface{}) {
	write(format, args...)
}

func LogErr(format string, args ...interface{}) {
	write(format, args...)
}

func Drain() []string {
	mu.Lock()
	defer mu.Unlock()
	if len(pending) == 0 {
		return nil
	}
	out := pending
	pending = nil
	return out
}

func Init(path string) {
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		logFile = nil
		return
	}
	logFile = f
	now := time.Now().Format("2006-01-02 15:04:05")
	fmt.Fprintf(logFile, "--\n%s (início)\n", now)
}

func Close() {
	if logFile == nil {
		return
	}
	now := time.Now().Format("2006-01-02 15:04:05")
	fmt.Fprintf(logFile, "%s (fim)\n--\n", now)
	logFile.Close()
	logFile = nil
}
