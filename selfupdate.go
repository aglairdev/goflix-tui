package main

import (
	"encoding/json"
	"net/http"
	"os"
	"os/exec"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

// Verificação de atualização

func checkUpdate() tea.Msg {
	time.Sleep(1000 * time.Millisecond)
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(RepoAPI)
	if err != nil {
		debugErr("update check HTTP falhou: %v", err)
		return updateCheckMsg{}
	}
	defer resp.Body.Close()
	var payload struct {
		TagName string `json:"tag_name"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		debugErr("update check JSON inválido: %v", err)
		return updateCheckMsg{}
	}
	if payload.TagName != "" && payload.TagName != Version {
		debug("update disponível: %s (atual: %s)", payload.TagName, Version)
		return updateCheckMsg{latest: payload.TagName}
	}
	debug("update check: %s já é o mais recente", Version)
	return updateCheckMsg{}
}

// doUpdate executa go install e reinicia o processo com o binário atualizado

func doUpdate() tea.Cmd {
	return func() tea.Msg {
		cmd := exec.Command("go", "install", "github.com/aglairdev/goflix@latest")
		cmd.Stdout, cmd.Stderr = os.Stdout, os.Stderr
		if err := cmd.Run(); err != nil {
			debugErr("go install falhou: %v", err)
			return flashMsg{text: t("update_error"), err: true}
		}
		bin, _ := os.Executable()
		if newBin, err := exec.LookPath("goflix"); err == nil && newBin != bin {
			os.Remove(bin)
			os.Rename(newBin, bin)
		}
		exec.Command(bin, os.Args[1:]...).Start()
		return tea.QuitMsg{}
	}
}
