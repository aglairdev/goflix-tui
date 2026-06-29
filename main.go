package main

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	AppName = "goflix"
	Version = "v1.1.0"
	RepoAPI = "https://api.github.com/repos/aglairdev/goflix/releases/latest"

	// Altere para qualquer cor hex válida hexadecimal
	ColorAccent = "#FF5FA7"
)

// Extensões de vídeo suportadas
var videoExts = map[string]bool{
	".mp4": true, ".mkv": true, ".avi": true, ".mov": true,
	".wmv": true, ".flv": true, ".webm": true, ".m4v": true,
	".ts": true, ".mpeg": true, ".mpg": true, ".3gp": true,
}

var debugMode bool

// Estilos
var (
	styleTitle      = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color(ColorAccent))
	styleVersion    = lipgloss.NewStyle().Foreground(lipgloss.Color("#AAAAAA"))
	styleDivider    = lipgloss.NewStyle().Foreground(lipgloss.Color(ColorAccent))
	styleFooterKey  = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color(ColorAccent))
	styleFooterDesc = lipgloss.NewStyle().Foreground(lipgloss.Color("#AAAAAA"))
	styleFooterSep  = lipgloss.NewStyle().Foreground(lipgloss.Color("#555555"))
	styleDir        = lipgloss.NewStyle().Foreground(lipgloss.Color(ColorAccent))
	styleProgress   = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFCC00"))
	styleWatched    = lipgloss.NewStyle().Faint(true).Foreground(lipgloss.Color("#5FAF5F"))
	styleNormal     = lipgloss.NewStyle().Foreground(lipgloss.Color("#DDDDDD"))
	styleMeta       = lipgloss.NewStyle().Foreground(lipgloss.Color("#AAAAAA"))
	styleSuccess    = lipgloss.NewStyle().Foreground(lipgloss.Color("#5FAF5F"))
	styleError      = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF5F5F"))
	styleLoading    = lipgloss.NewStyle().Faint(true).Foreground(lipgloss.Color("#FFCC00"))
	styleUpdate     = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#5FAF5F"))
)

func renderFooter(raw string) string {
	parts := strings.Split(raw, "  |  ")
	rendered := make([]string, len(parts))
	for i, part := range parts {
		kv := strings.SplitN(part, ": ", 2)
		if len(kv) == 2 {
			rendered[i] = styleFooterKey.Render(strings.TrimSpace(kv[0])) +
				styleFooterDesc.Render(": "+kv[1])
		} else {
			rendered[i] = styleFooterDesc.Render(part)
		}
	}
	return "  " + strings.Join(rendered, styleFooterSep.Render("  |  "))
}

// Config

// Caminhos calculados uma vez na inicialização
var (
	cfgDir      = initCfgDir()
	cfgFile     = filepath.Join(cfgDir, "config")
	watchedPath = filepath.Join(cfgDir, "watched")
)

func initCfgDir() string {
	if xdg := os.Getenv("XDG_CONFIG_HOME"); xdg != "" {
		return filepath.Join(xdg, "goflix")
	}
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".config", "goflix")
}

func mpvWatchDir() string {
	if xdg := os.Getenv("XDG_STATE_HOME"); xdg != "" {
		return filepath.Join(xdg, "mpv", "watch_later")
	}
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".local", "state", "mpv", "watch_later")
}

func ensureConfig() {
	os.MkdirAll(cfgDir, 0755)
	for _, f := range []string{cfgFile, watchedPath} {
		if _, err := os.Stat(f); os.IsNotExist(err) {
			os.WriteFile(f, nil, 0644)
		}
	}
}

func loadLines(path string) []string {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil
	}
	var lines []string
	for _, l := range strings.Split(string(data), "\n") {
		if l = strings.TrimSpace(l); l != "" && !strings.HasPrefix(l, "#") {
			lines = append(lines, l)
		}
	}
	return lines
}

func saveLines(path string, lines []string) error {
	return os.WriteFile(path, []byte(strings.Join(lines, "\n")+"\n"), 0644)
}

func addDir(path string) error {
	path = strings.TrimRight(path, "/")
	abs, err := filepath.Abs(path)
	if err != nil {
		return fmt.Errorf("%s", t("dir_invalid"))
	}
	info, err := os.Stat(abs)
	if err != nil || !info.IsDir() {
		return fmt.Errorf("%s", t("dir_invalid"))
	}
	dirs := loadLines(cfgFile)
	for _, d := range dirs {
		if d == abs {
			return fmt.Errorf("%s", t("dir_exists"))
		}
	}
	return saveLines(cfgFile, append(dirs, abs))
}

func removeDir(path string) {
	var nd []string
	for _, d := range loadLines(cfgFile) {
		if d != path {
			nd = append(nd, d)
		}
	}
	saveLines(cfgFile, nd)
}

// Assistidos

func loadWatched() map[string]int64 {
	w := map[string]int64{}
	for _, l := range loadLines(watchedPath) {
		if parts := strings.SplitN(l, "=", 2); len(parts) == 2 {
			ts, _ := strconv.ParseInt(parts[1], 10, 64)
			w[parts[0]] = ts
		}
	}
	return w
}

func saveWatched(w map[string]int64) {
	lines := make([]string, 0, len(w))
	for k, v := range w {
		lines = append(lines, fmt.Sprintf("%s=%d", k, v))
	}
	sort.Strings(lines)
	os.WriteFile(watchedPath, []byte(strings.Join(lines, "\n")+"\n"), 0644)
}

// setWatched adiciona (add=true) ou remove (add=false) um vídeo do histórico
func setWatched(path string, add bool) {
	w := loadWatched()
	if add {
		w[path] = time.Now().Unix()
	} else {
		delete(w, path)
	}
	saveWatched(w)
}

// MPV

// mpvHash retorna o MD5 uppercase que o mpv usa para nomear arquivos de retomada
func mpvHash(path string) string {
	return strings.ToUpper(fmt.Sprintf("%x", md5.Sum([]byte(path))))
}

func getResumePosition(path string) float64 {
	data, err := os.ReadFile(filepath.Join(mpvWatchDir(), mpvHash(path)))
	if err != nil {
		return 0
	}
	if m := regexp.MustCompile(`start=([0-9.]+)`).FindSubmatch(data); m != nil {
		v, _ := strconv.ParseFloat(string(m[1]), 64)
		return v
	}
	return 0
}

func resetResumePosition(path string) {
	os.Remove(filepath.Join(mpvWatchDir(), mpvHash(path)))
}

// ffprobe

func getDuration(path string) float64 {
	out, err := exec.Command("ffprobe", "-v", "error",
		"-show_entries", "format=duration",
		"-of", "default=noprint_wrappers=1:nokey=1",
		path).Output()
	if err != nil {
		return 0
	}
	v, _ := strconv.ParseFloat(strings.TrimSpace(string(out)), 64)
	return v
}

func formatTime(secs float64) string {
	if secs <= 0 {
		return t("dur_unknown")
	}
	s := int(math.Round(secs))
	if h := s / 3600; h > 0 {
		return fmt.Sprintf("%dh%02dm", h, (s%3600)/60)
	}
	return fmt.Sprintf("%dm", (s%3600)/60)
}

// Vídeos

type videoFile struct {
	path      string
	name      string
	duration  float64
	resume    float64
	watched   bool
	watchedAt int64
}

func loadVideos(dir string, watched map[string]int64) []videoFile {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil
	}
	var files []videoFile
	for _, e := range entries {
		if e.IsDir() || !videoExts[strings.ToLower(filepath.Ext(e.Name()))] {
			continue
		}
		path := filepath.Join(dir, e.Name())
		ts, isWatched := watched[path]
		files = append(files, videoFile{
			path: path, name: e.Name(),
			duration: getDuration(path), resume: getResumePosition(path),
			watched: isWatched, watchedAt: ts,
		})
	}
	sort.Slice(files, func(i, j int) bool { return files[i].name < files[j].name })
	return files
}

// hasDirItems retorna true se a lista contém apenas dirItem (sem vídeos)
func hasDirItems(items []list.Item) bool {
	for _, item := range items {
		if _, ok := item.(dirItem); ok {
			return true
		}
		if _, ok := item.(videoItem); ok {
			return false
		}
	}
	return false
}

// Itens da lista

// dirItem representa diretórios raiz e subdiretórios
type dirItem struct{ path string }

func (d dirItem) Title() string       { return styleDir.Render("› ") + filepath.Base(d.path) + "/" }
func (d dirItem) Description() string { return "" }
func (d dirItem) FilterValue() string { return filepath.Base(d.path) }

// videoItem representa um vídeo na lista, com seção visual
type videoItem struct {
	file    videoFile
	section string // "progress" | "normal" | "watched" | "sep"
}

func (v videoItem) FilterValue() string { return v.file.name }
func (v videoItem) Description() string { return "" }

func (v videoItem) Title() string {
	if v.section == "sep" {
		return styleDivider.Render(strings.Repeat("─", 64))
	}
	name := strings.TrimSuffix(v.file.name, filepath.Ext(v.file.name))
	dur, pos := formatTime(v.file.duration), formatTime(v.file.resume)
	switch v.section {
	case "progress":
		return styleNormal.Render(name) +
			styleMeta.Render(fmt.Sprintf("  %s %s %s  ", pos, t("progress"), dur)) +
			styleProgress.Render("▶ "+t("continue"))
	case "watched":
		return styleWatched.Render(name) +
			styleMeta.Render(fmt.Sprintf("  %s  ", dur)) +
			styleWatched.Render("✓ "+t("watched_label"))
	default:
		return styleNormal.Render(name) + styleMeta.Render(fmt.Sprintf("  %s", dur))
	}
}

// Delegate

type compactDelegate struct{ list.DefaultDelegate }

func newDelegate() compactDelegate {
	d := list.NewDefaultDelegate()
	d.ShowDescription = false
	d.SetSpacing(0)
	d.Styles.SelectedTitle = d.Styles.SelectedTitle.
		Border(lipgloss.NormalBorder(), false, false, false, true).
		BorderForeground(lipgloss.Color(ColorAccent)).
		Foreground(lipgloss.Color("#FFFFFF")).
		Bold(true).
		Padding(0, 0, 0, 1)
	d.Styles.NormalTitle = d.Styles.NormalTitle.
		Border(lipgloss.HiddenBorder(), false, false, false, true).
		Padding(0, 0, 0, 1)
	return compactDelegate{d}
}

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
	quitting      bool
	pendingDir    string
	latestVer     string
	renameTarget  string
}

func initialModel() model {
	ensureConfig()
	inp := textinput.New()
	inp.CharLimit, inp.Width = 512, 60
	m := model{screen: screenLoading, input: inp, watched: loadWatched()}
	m.reloadDirs()
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
}

// Verificação de atualização

// checkUpdate consulta a API do GitHub em background
func checkUpdate() tea.Msg {
	time.Sleep(1000 * time.Millisecond)
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(RepoAPI)
	if err != nil {
		return updateCheckMsg{}
	}
	defer resp.Body.Close()
	var payload struct {
		TagName string `json:"tag_name"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return updateCheckMsg{}
	}
	if payload.TagName != "" && payload.TagName != Version {
		return updateCheckMsg{latest: payload.TagName}
	}
	return updateCheckMsg{}
}

// doUpdate executa go install e reinicia o processo com o binário atualizado
func doUpdate() tea.Cmd {
	return func() tea.Msg {
		cmd := exec.Command("go", "install", "github.com/aglairdev/goflix@latest")
		cmd.Stdout, cmd.Stderr = os.Stdout, os.Stderr
		if err := cmd.Run(); err != nil {
			return flashMsg{text: t("update_error"), err: true}
		}
		bin, _ := os.Executable()
		exec.Command(bin, os.Args[1:]...).Start()
		return tea.QuitMsg{}
	}
}

func (m model) Init() tea.Cmd { return checkUpdate }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	m.flash = ""

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
		return m, textinput.Blink
	case "d":
		if item, ok := m.mainList.SelectedItem().(dirItem); ok {
			removeDir(item.path)
			m.reloadDirs()
			m.flash = t("dir_removed") + ": " + filepath.Base(item.path)
		}
		return m, nil
	case "l":
		toggleLang()
		m.flash = t("lang_changed") + " " + langLabel[currentLang]
		m.reloadDirs()
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
			fmsg = flashMsg{text: err.Error(), err: true}
		} else {
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

func (m model) playFile(path string) tea.Cmd {
	dur, startedAt := getDuration(path), time.Now()
	return tea.ExecProcess(exec.Command("mpv", "--save-position-on-quit", path), func(err error) tea.Msg {
		pos := getResumePosition(path)
		if dur > 0 && ((pos > 0 && pos/dur >= 0.90) || time.Since(startedAt).Seconds() >= dur*0.90) {
			setWatched(path, true)
		}
		return flashMsg{}
	})
}

func (m model) View() string {
	if m.quitting {
		return "  " + styleVersion.Render(AppName+" ꕤ - "+t("bye")) + "\n"
	}

	header := "  " + styleTitle.Render(AppName+" ꕤ") + "  " + styleVersion.Render(Version)
	if m.screen == screenFiles && m.curDir != "" && m.pendingDir == "" {
		header += "  " + styleDir.Render(filepath.Base(m.curDir)+"/")
	}

	var body string
	switch m.screen {
	case screenMain:
		if len(m.dirs) == 0 {
			body = "\n  " + styleNormal.Render(t("no_dirs")) + "\n"
		} else {
			body = "\n" + m.mainList.View() + "\n"
		}
		used := strings.Count(body, "\n") + 2
		if rem := m.height - used - 3; rem > 0 {
			body += strings.Repeat("\n", rem)
		}
	case screenFiles:
		if len(m.fileList.Items()) == 0 {
			body = "\n  " + styleMeta.Render(t("no_video")) + "\n"
		} else {
			body = "\n" + m.fileList.View() + "\n"
		}
		used := strings.Count(body, "\n") + 2
		if rem := m.height - used - 3; rem > 0 {
			body += strings.Repeat("\n", rem)
		}
	case screenInput:
		body = "  " + styleNormal.Render(t("prompt_dir")) + "\n\n" +
			"  " + m.input.View() + "\n\n"
	case screenRename:
		body = "\n  " + styleNormal.Render(t("rename_label")+": "+filepath.Base(m.renameTarget)) + "\n\n" +
			"  " + m.input.View() + "\n\n"
	case screenLoading:
		body = "  " + styleLoading.Render("⟳  "+t("loading")) + "\n"
	case screenUpdate:
		body = "\n  " + fmt.Sprintf("ꕤ %s: %s  (%s: %s)\n\n",
			t("update_available"), styleUpdate.Render(m.latestVer),
			t("update_current"), styleError.Render(Version)) +
			"  " + styleVersion.Render(t("update_prompt")) + "\n"
	}

	footer := ""
	if m.screen == screenFiles && hasDirItems(m.fileList.Items()) {
		footer = renderFooter(t(footerKeyDirs)) + "\n"
	} else if key, ok := footerKey[m.screen]; ok {
		footer = renderFooter(t(key)) + "\n"
	}

	flash := "\n"
	if m.flash != "" {
		style, prefix := styleSuccess, "✓  "
		if m.flashErr {
			style, prefix = styleError, "✗  "
		}
		flash = "  " + style.Render(prefix+m.flash) + "\n"
	}

	s := header + "\n" +
		styleDivider.Render(strings.Repeat("─", m.width)) + "\n" +
		body + footer +
		styleDivider.Render(strings.Repeat("─", m.width)) + "\n" +
		flash + "\n"
	lines := strings.Count(s, "\n")
	if rem := m.height - lines; rem > 0 {
		s += strings.Repeat("\n", rem)
	}
	return s
}

// Debug -d

func debug(format string, args ...interface{}) {
	if debugMode {
		fmt.Fprintf(os.Stderr, "[goflix-debug] "+format+"\n", args...)
	}
}

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
			fmt.Fprintf(os.Stderr, "  -d\tdebug mode (verbose stderr)\n")
			fmt.Fprintf(os.Stderr, "  -h\tshow this help\n")
			os.Exit(0)
		default:
			fmt.Fprintf(os.Stderr, "unknown flag: %s\n\n", arg)
			fmt.Fprintf(os.Stderr, "Usage: %s [flags]\n", os.Args[0])
			fmt.Fprintf(os.Stderr, "  -v\tshow version\n")
			fmt.Fprintf(os.Stderr, "  -d\tdebug mode\n")
			fmt.Fprintf(os.Stderr, "  -h\tshow this help\n")
			os.Exit(1)
		}
	}

	checkDeps()
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintln(os.Stderr, "erro:", err)
		os.Exit(1)
	}
}