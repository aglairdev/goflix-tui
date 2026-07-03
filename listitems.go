package main

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
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

func (d dirItem) Title() string       { return styleDir.Render("› ") + filepath.Base(d.path) + "/" }
func (d dirItem) Description() string { return "" }
func (d dirItem) FilterValue() string { return filepath.Base(d.path) }

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
		BorderForeground(lipgloss.Color(currentAccent)).
		Foreground(lipgloss.Color("#FFFFFF")).
		Bold(true).
		Padding(0, 0, 0, 1)
	d.Styles.NormalTitle = d.Styles.NormalTitle.
		Border(lipgloss.HiddenBorder(), false, false, false, true).
		Padding(0, 0, 0, 1)
	return compactDelegate{d}
}
