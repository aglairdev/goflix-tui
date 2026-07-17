package ui

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"

	"github.com/aglairdev/goflix/internal/i18n"
	"github.com/aglairdev/goflix/internal/theme"
	"github.com/aglairdev/goflix/internal/video"
)

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

type dirItem struct{ path string }

func (d dirItem) Title() string       { return theme.Dir.Render("› ") + filepath.Base(d.path) + "/" }
func (d dirItem) Description() string { return "" }
func (d dirItem) FilterValue() string { return filepath.Base(d.path) }

type videoItem struct {
	file    video.File
	section string // "progress" | "normal" | "watched" | "sep"
}

func (v videoItem) FilterValue() string { return v.file.Name }
func (v videoItem) Description() string { return "" }

func (v videoItem) Title() string {
	if v.section == "sep" {
		return theme.Divider.Render(strings.Repeat("─", 64))
	}
	name := strings.TrimSuffix(v.file.Name, filepath.Ext(v.file.Name))
	dur, pos := video.FormatTime(v.file.Duration), video.FormatTime(v.file.Resume)
	switch v.section {
	case "progress":
		return theme.Normal.Render(name) +
			theme.Meta.Render(fmt.Sprintf("  %s %s %s  ", pos, i18n.T("progress"), dur)) +
			theme.Progress.Render("▶ "+i18n.T("continue"))
	case "watched":
		return theme.Watched.Render(name) +
			theme.Meta.Render(fmt.Sprintf("  %s  ", dur)) +
			theme.Watched.Render("✓ "+i18n.T("watched_label"))
	default:
		return theme.Normal.Render(name) + theme.Meta.Render(fmt.Sprintf("  %s", dur))
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
		BorderForeground(lipgloss.Color(theme.CurrentAccent)).
		Foreground(lipgloss.Color("#FFFFFF")).
		Bold(true).
		Padding(0, 0, 0, 1)
	d.Styles.NormalTitle = d.Styles.NormalTitle.
		Border(lipgloss.HiddenBorder(), false, false, false, true).
		Padding(0, 0, 0, 1)
	return compactDelegate{d}
}
