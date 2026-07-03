package main

import (
	"crypto/md5"
	"fmt"
	"math"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

// Extensões de vídeo suportadas

var videoExts = map[string]bool{
	".mp4": true, ".mkv": true, ".avi": true, ".mov": true,
	".wmv": true, ".flv": true, ".webm": true, ".m4v": true,
	".ts": true, ".mpeg": true, ".mpg": true, ".3gp": true,
}

// MPV

func mpvWatchDir() string {
	if xdg := os.Getenv("XDG_STATE_HOME"); xdg != "" {
		return filepath.Join(xdg, "mpv", "watch_later")
	}
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".local", "state", "mpv", "watch_later")
}

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
		debugErr("ffprobe falhou para %s: %v", path, err)
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
		debugErr("falha ao ler diretório %s: %v", dir, err)
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
