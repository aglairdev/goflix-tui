package selfupdate

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/aglairdev/goflix/internal/debug"
	"github.com/aglairdev/goflix/internal/i18n"
	"github.com/aglairdev/goflix/internal/version"
)

// Mensagens

type CheckMsg struct{ Latest string }
type ResultMsg struct {
	Text string
	Err  bool
}

// Verificação de atualização

func Check() tea.Msg {
	time.Sleep(1000 * time.Millisecond)
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(version.RepoAPI)
	if err != nil {
		debug.LogErr("update check HTTP falhou: %v", err)
		return CheckMsg{}
	}
	defer resp.Body.Close()
	var payload struct {
		TagName string `json:"tag_name"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		debug.LogErr("update check JSON inválido: %v", err)
		return CheckMsg{}
	}
	if payload.TagName != "" && payload.TagName != version.Version {
		debug.Log("update disponível: %s (atual: %s)", payload.TagName, version.Version)
		return CheckMsg{Latest: payload.TagName}
	}
	debug.Log("update check: %s já é o mais recente", version.Version)
	return CheckMsg{}
}

func goBinPath() (string, error) {
	gobin, err := exec.Command("go", "env", "GOBIN").Output()
	if err == nil {
		if dir := strings.TrimSpace(string(gobin)); dir != "" {
			return filepath.Join(dir, "goflix"), nil
		}
	}
	gopath, err := exec.Command("go", "env", "GOPATH").Output()
	if err != nil {
		return "", err
	}
	dir := strings.TrimSpace(string(gopath))
	if dir == "" {
		return "", fmt.Errorf("GOPATH vazio")
	}
	return filepath.Join(dir, "bin", "goflix"), nil
}

func Update() tea.Cmd {
	return func() tea.Msg {
		cmd := exec.Command("go", "install", "github.com/aglairdev/goflix@latest")
		cmd.Stdout, cmd.Stderr = os.Stdout, os.Stderr
		if err := cmd.Run(); err != nil {
			debug.LogErr("go install falhou: %v", err)
			return ResultMsg{Text: i18n.T("update_error"), Err: true}
		}
		bin, err := goBinPath()
		if err != nil {
			debug.LogErr("não foi possível determinar destino do go install: %v", err)
			return ResultMsg{Text: i18n.T("update_error"), Err: true}
		}
		if _, err := os.Stat(bin); err != nil {
			debug.LogErr("binário atualizado não encontrado em %s: %v", bin, err)
			return ResultMsg{Text: i18n.T("update_error"), Err: true}
		}
		if err := exec.Command(bin, os.Args[1:]...).Start(); err != nil {
			debug.LogErr("falha ao reiniciar goflix atualizado: %v", err)
			return ResultMsg{Text: i18n.T("update_error"), Err: true}
		}
		return tea.QuitMsg{}
	}
}