//
// goflix ~ gerenciador de vídeos no terminal
// © 2026 ~ AGL ~ github.com/aglairdev
// licença: MIT
//
// uso:   ./goflix
//        ./goflix -d     //debug
//        ./goflix -v     //versão
//        ./goflix -h     //ajuda
//

package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/aglairdev/goflix/internal/config"
	"github.com/aglairdev/goflix/internal/debug"
	"github.com/aglairdev/goflix/internal/ui"
	"github.com/aglairdev/goflix/internal/version"
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
			fmt.Printf("%s %s\n", version.AppName, version.Version)
			os.Exit(0)
		case "-d":
			debug.Enabled = true
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

	if debug.Enabled {
		logPath := filepath.Join(config.CfgDir, "debug.log")
		debug.Init(logPath)
		debug.Log("modo debug iniciado ~ log: %s", logPath)
	}

	p := tea.NewProgram(ui.New(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		debug.Close()
		fmt.Fprintln(os.Stderr, "erro:", err)
		os.Exit(1)
	}
	debug.Close()
}
