package selfupdate

import (
	"encoding/json"
	"net/http"
	"os"
	"os/exec"
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

func Update() tea.Cmd {
	return func() tea.Msg {
		cmd := exec.Command("go", "install", "github.com/aglairdev/goflix@latest")
		cmd.Stdout, cmd.Stderr = os.Stdout, os.Stderr
		if err := cmd.Run(); err != nil {
			debug.LogErr("go install falhou: %v", err)
			return ResultMsg{Text: i18n.T("update_error"), Err: true}
		}
		bin, err := os.Executable()
		if err != nil {
			debug.LogErr("não foi possível localizar o binário atual: %v", err)
			return ResultMsg{Text: i18n.T("update_error"), Err: true}
		}
		newBin, err := exec.LookPath("goflix")
		if err != nil {
			debug.LogErr("binário novo não encontrado no PATH: %v", err)
			return ResultMsg{Text: i18n.T("update_error"), Err: true}
		}
		if newBin != bin {
			if err := replaceBinary(newBin, bin); err != nil {
				debug.LogErr("falha ao mover binário atualizado para %s: %v", bin, err)
				return ResultMsg{Text: i18n.T("update_error"), Err: true}
			}
		}
		if err := exec.Command(bin, os.Args[1:]...).Start(); err != nil {
			debug.LogErr("falha ao reiniciar goflix atualizado: %v", err)
			return ResultMsg{Text: i18n.T("update_error"), Err: true}
		}
		return tea.QuitMsg{}
	}
}

func replaceBinary(newBin, bin string) error {
	if err := os.Rename(newBin, bin); err == nil {
		return nil
	}
	data, err := os.ReadFile(newBin)
	if err != nil {
		return err
	}
	tmp := bin + ".new"
	if err := os.WriteFile(tmp, data, 0755); err != nil {
		return err
	}
	if err := os.Rename(tmp, bin); err != nil {
		os.Remove(tmp)
		return err
	}
	os.Remove(newBin)
	return nil
}