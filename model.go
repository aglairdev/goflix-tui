package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"time"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

// Mensagens assíncronas

type loadDirMsg struct{ dir string }
type flashMsg struct {
	text string
	err  bool
}
type updateCheckMsg struct{ latest string }

// Telas

type screen int

const (
	screenMain screen = iota
	screenFiles
	screenInput
	screenLoading
	screenUpdate
	screenRename
)

var footerKey = map[screen]string{
	screenMain:   "footer_main",
	screenFiles:  "footer_files",
	screenInput:  "footer_input",
	screenRename: "footer_rename",
}

// footerKey para tela de arquivos com apenas diretórios (sem v/r)

var footerKeyDirs = "footer_files_dirs"

// Model

type model struct {
	screen        screen
	width, height int
	mainList      list.Model
	fileList      list.Model
	input         textinput.Model
	dirs          []string
	curDir        string
	watched       map[string]int64
	flash         string
	flashErr      bool
	debugFlash    string
	logDirEntry   bool
	quitting      bool
	pendingDir    string
	latestVer     string
	renameTarget  string
}

func initialModel() model {
	ensureConfig()
	loadTheme()
	inp := textinput.New()
	inp.CharLimit, inp.Width = 512, 60
	m := model{screen: screenLoading, input: inp, watched: loadWatched()}
	m.reloadDirs()
	if debugMode {
		m.debugFlash = "[goflix-debug] modo debug ativo"
	}
	return m
}

func newList(w, h int) list.Model {
	l := list.New([]list.Item{}, newDelegate(), w, h)
	l.SetShowTitle(false)
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(true)
	l.SetShowHelp(false)
	l.SetShowPagination(false)
	return l
}

func (m *model) reloadDirs() {
	m.dirs = loadLines(cfgFile)
	items := make([]list.Item, len(m.dirs))
	for i, d := range m.dirs {
		items[i] = dirItem{path: d}
	}
	if m.width == 0 || m.height == 0 {
		l := newList(0, 10)
		l.SetItems(items)
		m.mainList = l
		m.fileList = newList(0, 10)
		return
	}
	m.mainList.SetItems(items)
	m.mainList.SetDelegate(newDelegate())

}

func (m *model) loadDir(dir string) {
	m.curDir = dir
	m.watched = loadWatched()
	videos := loadVideos(dir, m.watched)

	h := m.height - 5
	if h < 1 {
		h = 10
	}

	if len(videos) == 0 {
		entries, _ := os.ReadDir(dir)
		var sub []list.Item
		for _, e := range entries {
			if e.IsDir() {
				sub = append(sub, dirItem{path: filepath.Join(dir, e.Name())})
			}
		}
		if len(sub) > 0 {
			l := newList(m.width, h)
			l.SetItems(sub)
			m.fileList = l
			return
		}
	}

	var inProgress, normal, watched []videoFile
	for _, v := range videos {
		switch {
		case v.watched:
			watched = append(watched, v)
		case v.resume > 0:
			inProgress = append(inProgress, v)
		default:
			normal = append(normal, v)
		}
	}
	sort.Slice(inProgress, func(i, j int) bool { return inProgress[i].resume > inProgress[j].resume })
	sort.Slice(watched, func(i, j int) bool { return watched[i].watchedAt > watched[j].watchedAt })

	var items []list.Item
	sep := videoItem{section: "sep"}
	for _, v := range inProgress {
		items = append(items, videoItem{file: v, section: "progress"})
	}
	if len(inProgress) > 0 && len(normal)+len(watched) > 0 {
		items = append(items, sep)
	}
	for _, v := range normal {
		items = append(items, videoItem{file: v, section: "normal"})
	}
	if len(normal) > 0 && len(watched) > 0 {
		items = append(items, sep)
	}
	for _, v := range watched {
		items = append(items, videoItem{file: v, section: "watched"})
	}

	l := newList(m.width, h)
	l.SetItems(items)
	m.fileList = l
	if m.logDirEntry {
		debug("%d vídeos em %s", len(videos), dir)
		m.logDirEntry = false
	}
}

func (m model) playFile(path string) tea.Cmd {
	dur, startedAt := getDuration(path), time.Now()
	return tea.ExecProcess(exec.Command("mpv", "--save-position-on-quit", path), func(err error) tea.Msg {
		if err != nil {
			debugErr("mpv retornou erro para %s: %v", path, err)
		}
		pos := getResumePosition(path)
		if dur > 0 && ((pos > 0 && pos/dur >= 0.90) || time.Since(startedAt).Seconds() >= dur*0.90) {
			setWatched(path, true)
		}
		return flashMsg{}
	})
}

func (m model) Init() tea.Cmd { return checkUpdate }
