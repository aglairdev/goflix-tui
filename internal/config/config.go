package config

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/aglairdev/goflix/internal/debug"
	"github.com/aglairdev/goflix/internal/i18n"
	"github.com/aglairdev/goflix/internal/theme"
)

// Config

var (
	CfgDir       = initCfgDir()
	cfgFile      = filepath.Join(CfgDir, "config")
	watchedPath  = filepath.Join(CfgDir, "watched")
	settingsPath = filepath.Join(CfgDir, "settings")
)

func initCfgDir() string {
	if xdg := os.Getenv("XDG_CONFIG_HOME"); xdg != "" {
		return filepath.Join(xdg, "goflix")
	}
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".config", "goflix")
}

func EnsureConfig() {
	os.MkdirAll(CfgDir, 0755)
	for _, f := range []string{cfgFile, watchedPath, settingsPath} {
		if _, err := os.Stat(f); os.IsNotExist(err) {
			os.WriteFile(f, nil, 0644)
		}
	}
	debug.Log("config dir: %s", CfgDir)
}

func loadLines(path string) []string {
	data, err := os.ReadFile(path)
	if err != nil {
		debug.LogErr("falha ao ler %s: %v", path, err)
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
	err := os.WriteFile(path, []byte(strings.Join(lines, "\n")+"\n"), 0644)
	if err != nil {
		debug.LogErr("falha ao salvar %s: %v", path, err)
	}
	return err
}

func ListDirs() []string {
	return loadLines(cfgFile)
}

func AddDir(path string) error {
	path = strings.TrimRight(path, "/")
	abs, err := filepath.Abs(path)
	if err != nil {
		return fmt.Errorf("%s", i18n.T("dir_invalid"))
	}
	info, err := os.Stat(abs)
	if err != nil || !info.IsDir() {
		return fmt.Errorf("%s", i18n.T("dir_invalid"))
	}
	dirs := loadLines(cfgFile)
	for _, d := range dirs {
		if d == abs {
			return fmt.Errorf("%s", i18n.T("dir_exists"))
		}
	}
	return saveLines(cfgFile, append(dirs, abs))
}

func RemoveDir(path string) {
	var nd []string
	for _, d := range loadLines(cfgFile) {
		if d != path {
			nd = append(nd, d)
		}
	}
	saveLines(cfgFile, nd)
	debug.Log("diretório removido: %s", path)
}

func LoadTheme() {
	data, err := os.ReadFile(settingsPath)
	if err != nil {
		debug.LogErr("settings não encontrado, usando padrão: %v", err)
		theme.Current = 0
		theme.Apply(theme.Themes[0].Accent)
		return
	}
	for _, l := range strings.Split(string(data), "\n") {
		if parts := strings.SplitN(l, "=", 2); len(parts) == 2 && parts[0] == "theme" {
			for i, th := range theme.Themes {
				if th.Name == parts[1] {
					theme.Current = i
					theme.Apply(th.Accent)
					return
				}
			}
		}
	}
	debug.Log("settings corrompido, usando tema padrão")
	theme.Current = 0
	theme.Apply(theme.Themes[0].Accent)
}

func SaveTheme() {
	os.WriteFile(settingsPath, []byte("theme="+theme.Themes[theme.Current].Name+"\n"), 0644)
	debug.Log("tema salvo: %s", theme.Themes[theme.Current].Name)
}

// Assistidos

func LoadWatched() map[string]int64 {
	w := map[string]int64{}
	for _, l := range loadLines(watchedPath) {
		if parts := strings.SplitN(l, "=", 2); len(parts) == 2 {
			ts, _ := strconv.ParseInt(parts[1], 10, 64)
			w[parts[0]] = ts
		} else {
			debug.LogErr("linha ignorada (formato inválido): %s", l)
		}
	}
	return w
}

func SaveWatched(w map[string]int64) {
	lines := make([]string, 0, len(w))
	for k, v := range w {
		lines = append(lines, fmt.Sprintf("%s=%d", k, v))
	}
	sort.Strings(lines)
	os.WriteFile(watchedPath, []byte(strings.Join(lines, "\n")+"\n"), 0644)
}

func SetWatched(path string, add bool) {
	w := LoadWatched()
	if add {
		w[path] = time.Now().Unix()
	} else {
		delete(w, path)
	}
	SaveWatched(w)
}
