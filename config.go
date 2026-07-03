package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
)

// Config

var (
	cfgDir       = initCfgDir()
	cfgFile      = filepath.Join(cfgDir, "config")
	watchedPath  = filepath.Join(cfgDir, "watched")
	settingsPath = filepath.Join(cfgDir, "settings")
)

func initCfgDir() string {
	if xdg := os.Getenv("XDG_CONFIG_HOME"); xdg != "" {
		return filepath.Join(xdg, "goflix")
	}
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".config", "goflix")
}

func ensureConfig() {
	os.MkdirAll(cfgDir, 0755)
	for _, f := range []string{cfgFile, watchedPath, settingsPath} {
		if _, err := os.Stat(f); os.IsNotExist(err) {
			os.WriteFile(f, nil, 0644)
		}
	}
	debug("config dir: %s", cfgDir)
}

func loadLines(path string) []string {
	data, err := os.ReadFile(path)
	if err != nil {
		debugErr("falha ao ler %s: %v", path, err)
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
		debugErr("falha ao salvar %s: %v", path, err)
	}
	return err
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
	debug("diretório removido: %s", path)
}

func loadTheme() {
	data, err := os.ReadFile(settingsPath)
	if err != nil {
		debugErr("settings não encontrado, usando padrão: %v", err)
		currentTheme = 0
		applyTheme(themes[0].accent)
		return
	}
	for _, l := range strings.Split(string(data), "\n") {
		if parts := strings.SplitN(l, "=", 2); len(parts) == 2 && parts[0] == "theme" {
			for i, t := range themes {
				if t.name == parts[1] {
					currentTheme = i
					applyTheme(t.accent)
					return
				}
			}
		}
	}
	debug("settings corrompido, usando tema padrão")
	currentTheme = 0
	applyTheme(themes[0].accent)
}

func saveTheme() {
	os.WriteFile(settingsPath, []byte("theme="+themes[currentTheme].name+"\n"), 0644)
	debug("tema salvo: %s", themes[currentTheme].name)
}

// Assistidos

func loadWatched() map[string]int64 {
	w := map[string]int64{}
	for _, l := range loadLines(watchedPath) {
		if parts := strings.SplitN(l, "=", 2); len(parts) == 2 {
			ts, _ := strconv.ParseInt(parts[1], 10, 64)
			w[parts[0]] = ts
		} else {
			debugErr("linha ignorada (formato inválido): %s", l)
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

func setWatched(path string, add bool) {
	w := loadWatched()
	if add {
		w[path] = time.Now().Unix()
	} else {
		delete(w, path)
	}
	saveWatched(w)
}
