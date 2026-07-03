package main

import (
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.screen != screenLoading {
		mu.Lock()
		if len(pendingDebug) > 0 {
			m.debugFlash = strings.Join(pendingDebug, "\n")
			pendingDebug = nil
		}
		mu.Unlock()
	}

	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		h := msg.Height - 5
		if h < 1 {
			h = 10
		}
		m.mainList.SetSize(msg.Width, h)
		m.fileList.SetSize(msg.Width, h)
		return m, nil

	case updateCheckMsg:
		m.latestVer = msg.latest
		if msg.latest != "" {
			m.screen = screenUpdate
		} else {
			m.screen = screenMain
		}
		return m, nil

	case loadDirMsg:
		m.logDirEntry = true
		m.loadDir(msg.dir)
		m.pendingDir = ""
		m.screen = screenFiles
		return m, nil

	case flashMsg:
		m.flash, m.flashErr = msg.text, msg.err
		if m.screen == screenFiles && m.curDir != "" {
			m.loadDir(m.curDir)
		}
		return m, nil

	case tea.KeyMsg:
		m.flash = ""
		m.debugFlash = ""
		if m.screen == screenLoading {
			return m, nil
		}
		if m.screen == screenUpdate {
			if msg.String() == "u" {
				return m, doUpdate()
			}
			m.screen = screenMain
			return m, nil
		}
		switch m.screen {
		case screenInput:
			return m.updateInput(msg)
		case screenRename:
			return m.updateRename(msg)
		case screenFiles:
			return m.updateFiles(msg)
		default:
			return m.updateMain(msg)
		}
	}

	var cmd tea.Cmd
	switch m.screen {
	case screenMain:
		m.mainList, cmd = m.mainList.Update(msg)
	case screenFiles:
		m.fileList, cmd = m.fileList.Update(msg)
	}
	return m, cmd
}

func (m model) updateMain(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "ctrl+c":
		m.quitting = true
		return m, tea.Quit
	case "n":
		m.screen = screenInput
		m.input.SetValue("")
		m.input.Focus()
		return m, tea.Batch(tea.ClearScreen, textinput.Blink)
	case "d":
		if item, ok := m.mainList.SelectedItem().(dirItem); ok {
			removeDir(item.path)
			m.reloadDirs()
			n := len(m.mainList.Items())
			if n > 0 && m.mainList.Index() >= n {
				m.mainList.Select(n - 1)
			}
			m.flash = t("dir_removed") + ": " + filepath.Base(item.path)
		}
		return m, nil
	case "l":
		toggleLang()
		m.flash = t("lang_changed") + " " + langLabel[currentLang]
		m.reloadDirs()
		return m, nil
	case "t":
		currentTheme = (currentTheme + 1) % len(themes)
		applyTheme(themes[currentTheme].accent)
		saveTheme()
		m.flash = "Tema: " + themes[currentTheme].name
		m.reloadDirs()
		if m.screen == screenFiles && m.curDir != "" {
			m.loadDir(m.curDir)
		}
		return m, nil
	case "enter":
		if item, ok := m.mainList.SelectedItem().(dirItem); ok {
			m.pendingDir = item.path
			m.curDir = ""
			m.screen = screenLoading
			return m, func() tea.Msg { return loadDirMsg{dir: item.path} }
		}
	}
	var cmd tea.Cmd
	m.mainList, cmd = m.mainList.Update(msg)
	return m, cmd
}

func (m model) updateFiles(msg tea.KeyMsg) (tea.Model, tea.Cmd) {

	// v e r só funcionam quando há vídeos na lista (não em listagem de diretórios)

	isVideoList := !hasDirItems(m.fileList.Items())

	switch msg.String() {
	case "q":
		m.quitting = true
		return m, tea.Quit
	case "t":
		currentTheme = (currentTheme + 1) % len(themes)
		applyTheme(themes[currentTheme].accent)
		saveTheme()
		m.flash = "Tema: " + themes[currentTheme].name
		m.loadDir(m.curDir)
		return m, nil
	case "esc":
		for _, d := range m.dirs {
			if m.curDir == d {
				m.curDir, m.screen = "", screenMain
				m.reloadDirs()
				return m, nil
			}
		}
		m.loadDir(filepath.Dir(m.curDir))
		return m, nil
	case "v":
		if !isVideoList {
			return m, nil
		}
		if vi, ok := m.fileList.SelectedItem().(videoItem); ok && vi.section != "sep" {
			setWatched(vi.file.path, true)
			m.watched[vi.file.path] = time.Now().Unix()
			m.flash = t("marked_watched") + ": " + strings.TrimSuffix(vi.file.name, filepath.Ext(vi.file.name))
			m.loadDir(m.curDir)
		}
		return m, nil
	case "r":
		if !isVideoList {
			return m, nil
		}
		if vi, ok := m.fileList.SelectedItem().(videoItem); ok && vi.section != "sep" {
			setWatched(vi.file.path, false)
			resetResumePosition(vi.file.path)
			delete(m.watched, vi.file.path)
			m.flash = t("unmarked_watched") + ": " + strings.TrimSuffix(vi.file.name, filepath.Ext(vi.file.name))
			m.loadDir(m.curDir)
		}
		return m, nil
	case "a":
		var target string
		if di, ok := m.fileList.SelectedItem().(dirItem); ok {
			target = di.path
		} else if vi, ok := m.fileList.SelectedItem().(videoItem); ok && vi.section != "sep" {
			target = vi.file.path
		}
		if target != "" {
			m.renameTarget = target
			m.screen = screenRename
			m.input.SetValue(filepath.Base(target))
			m.input.Focus()
			return m, textinput.Blink
		}
		return m, nil
	case "enter":
		if di, ok := m.fileList.SelectedItem().(dirItem); ok {
			m.pendingDir = di.path
			m.curDir = ""
			m.screen = screenLoading
			return m, func() tea.Msg { return loadDirMsg{dir: di.path} }
		}
		if vi, ok := m.fileList.SelectedItem().(videoItem); ok && vi.section != "sep" {
			return m, m.playFile(vi.file.path)
		}
	}
	var cmd tea.Cmd
	m.fileList, cmd = m.fileList.Update(msg)
	return m, cmd
}

func (m model) updateInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.screen = screenMain
		m.input.SetValue("")
		return m, nil
	case "enter":
		path := strings.TrimSpace(m.input.Value())
		if strings.HasPrefix(path, "~/") {
			home, _ := os.UserHomeDir()
			path = filepath.Join(home, path[2:])
		}
		if err := addDir(path); err != nil {
			m.flash, m.flashErr = err.Error(), true
		} else {
			m.flash, m.flashErr = t("dir_added")+": "+filepath.Base(path), false
			m.screen = screenMain
			m.reloadDirs()
		}
		m.input.SetValue("")
		return m, nil
	}
	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

func (m model) updateRename(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.screen = screenFiles
		return m, nil
	case "enter":
		newName := strings.TrimSpace(m.input.Value())
		if newName == "" || newName == filepath.Base(m.renameTarget) {
			m.screen = screenFiles
			return m, nil
		}
		newPath := filepath.Join(filepath.Dir(m.renameTarget), newName)
		var fmsg flashMsg
		if err := os.Rename(m.renameTarget, newPath); err != nil {
			debugErr("rename falhou: %v", err)
			fmsg = flashMsg{text: err.Error(), err: true}
		} else {
			oldHL := filepath.Join(mpvWatchDir(), mpvHash(m.renameTarget))
			if _, err := os.Stat(oldHL); err == nil {
				debugErr(`histórico perdido: "%s" (watch_later existia)`, filepath.Base(m.renameTarget))
				os.Remove(oldHL)
			} else {
				debug("arquivo renomeado: %s → %s", filepath.Base(m.renameTarget), newName)
			}
			fmsg = flashMsg{text: t("renamed") + ": " + newName, err: false}
		}
		m.renameTarget = ""
		m.screen = screenFiles
		m.loadDir(m.curDir)
		return m, func() tea.Msg { return fmsg }
	}
	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	return m, cmd
}
