//
// goflix ~ gerenciador de vídeos no terminal
// © 2026 ~ AGL ~ github.com/aglairdev
// licença: MIT
//

package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

const (
	AppName = "goflix"
	Version = "v1.1.5"
	RepoAPI = "https://api.github.com/repos/aglairdev/goflix/releases/latest"
)

// Dependências

func checkDeps() {
	if _, err := exec.LookPath("mpv"); err != nil {
		fmt.Fprintln(os.Stderr, "✗ mpv não encontrado - instale via gerenciador de pacotes da sua distro.")
		os.Exit(1)
	}
	if _, err := exec.LookPath("ffprobe"); err != nil {
		fmt.Fprintln(os.Stderr, "· ffprobe não encontrado - duração dos vídeos não será exibida.")
		fmt.Fprintln(os.Stderr, "  Instale o ffmpeg via gerenciador de pacotes da sua distro.")
		time.Sleep(1500 * time.Millisecond)
	}
}

// Main

func main() {
	for _, arg := range os.Args[1:] {
		switch arg {
		case "-v":
			fmt.Printf("%s %s\n", AppName, Version)
			os.Exit(0)
		case "-d":
			debugMode = true
		case "-h":
			fmt.Fprintf(os.Stderr, "Usage: %s [flags]\n", os.Args[0])
			fmt.Fprintf(os.Stderr, "\nFlags:\n")
			fmt.Fprintf(os.Stderr, "  -v\tshow version\n")
			fmt.Fprintf(os.Stderr, "  -d\tmodo debug (verbose stderr)\n")
			fmt.Fprintf(os.Stderr, "  -h\tshow this help\n")
			os.Exit(0)
		default:
			fmt.Fprintf(os.Stderr, "unknown flag: %s\n\n", arg)
			fmt.Fprintf(os.Stderr, "Usage: %s [flags]\n", os.Args[0])
			fmt.Fprintf(os.Stderr, "  -v\tshow version\n")
			fmt.Fprintf(os.Stderr, "  -d\tmodo debug\n")
			fmt.Fprintf(os.Stderr, "  -h\tshow this help\n")
			os.Exit(1)
		}
	}

	checkDeps()

	if debugMode {
		var err error
		logPath := filepath.Join(cfgDir, "debug.log")
		logFile, err = os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			logFile = nil
		}
		now := time.Now().Format("2006-01-02 15:04:05")
		fmt.Fprintf(logFile, "--\n%s (início)\n", now)
		debug("modo debug iniciado ~ log: %s", logPath)
	}

	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		if logFile != nil {
			logFile.Close()
		}
		fmt.Fprintln(os.Stderr, "erro:", err)
		os.Exit(1)
	}
	if logFile != nil {
		now := time.Now().Format("2006-01-02 15:04:05")
		fmt.Fprintf(logFile, "%s (fim)\n--\n", now)
		logFile.Close()
	}
}
