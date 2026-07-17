package video

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

	"github.com/aglairdev/goflix/internal/debug"
	"github.com/aglairdev/goflix/internal/i18n"
)

// Extensões de vídeo suportadas

var exts = map[string]bool{
	".mp4": true, ".mkv": true, ".avi": true, ".mov": true,
	".wmv": true, ".flv": true, ".webm": true, ".m4v": true,
	".ts": true, ".mpeg": true, ".mpg": true, ".3gp": true,
}

func MpvWatchDir() string {
	if xdg := os.Getenv("XDG_STATE_HOME"); xdg != "" {
		return filepath.Join(xdg, "mpv", "watch_later")
	}
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".local", "state", "mpv", "watch_later")
}

func MpvHash(path string) string {
	return strings.ToUpper(fmt.Sprintf("%x", md5.Sum([]byte(path))))
}

func GetResumePosition(path string) float64 {
	data, err := os.ReadFile(filepath.Join(MpvWatchDir(), MpvHash(path)))
	if err != nil {
		return 0
	}
	if m := regexp.MustCompile(`start=([0-9.]+)`).FindSubmatch(data); m != nil {
		v, _ := strconv.ParseFloat(string(m[1]), 64)
		return v
	}
	return 0
}

func ResetResumePosition(path string) {
	os.Remove(filepath.Join(MpvWatchDir(), MpvHash(path)))
}

// ffprobe

func GetDuration(path string) float64 {
	out, err := exec.Command("ffprobe", "-v", "error",
		"-show_entries", "format=duration",
		"-of", "default=noprint_wrappers=1:nokey=1",
		path).Output()
	if err != nil {
		debug.LogErr("ffprobe falhou para %s: %v", path, err)
		return 0
	}
	v, _ := strconv.ParseFloat(strings.TrimSpace(string(out)), 64)
	return v
}

func FormatTime(secs float64) string {
	if secs <= 0 {
		return i18n.T("dur_unknown")
	}
	s := int(math.Round(secs))
	if h := s / 3600; h > 0 {
		return fmt.Sprintf("%dh%02dm", h, (s%3600)/60)
	}
	return fmt.Sprintf("%dm", (s%3600)/60)
}

// Vídeos

type File struct {
	Path      string
	Name      string
	Duration  float64
	Resume    float64
	Watched   bool
	WatchedAt int64
}

func Load(dir string, watched map[string]int64) []File {
	entries, err := os.ReadDir(dir)
	if err != nil {
		debug.LogErr("falha ao ler diretório %s: %v", dir, err)
		return nil
	}
	var files []File
	for _, e := range entries {
		if e.IsDir() || !exts[strings.ToLower(filepath.Ext(e.Name()))] {
			continue
		}
		path := filepath.Join(dir, e.Name())
		ts, isWatched := watched[path]
		files = append(files, File{
			Path: path, Name: e.Name(),
			Duration: GetDuration(path), Resume: GetResumePosition(path),
			Watched: isWatched, WatchedAt: ts,
		})
	}
	sort.Slice(files, func(i, j int) bool { return files[i].Name < files[j].Name })
	return files
}
