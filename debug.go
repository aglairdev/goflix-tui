package main

import (
	"fmt"
	"os"
	"sync"
)

// Debug -d

var debugMode bool
var logFile *os.File

var (
	mu           sync.Mutex
	pendingDebug []string
)

func debug(format string, args ...interface{}) {
	if !debugMode {
		return
	}
	msg := fmt.Sprintf("[goflix-debug] "+format, args...)
	fmt.Fprintln(os.Stderr, msg)
	mu.Lock()
	pendingDebug = append(pendingDebug, msg)
	mu.Unlock()
	if logFile != nil {
		fmt.Fprintln(logFile, msg)
	}
}

func debugErr(format string, args ...interface{}) {
	if !debugMode {
		return
	}
	msg := fmt.Sprintf("[goflix-debug] "+format, args...)
	fmt.Fprintln(os.Stderr, msg)
	mu.Lock()
	pendingDebug = append(pendingDebug, msg)
	mu.Unlock()
	if logFile != nil {
		fmt.Fprintln(logFile, msg)
	}
}
