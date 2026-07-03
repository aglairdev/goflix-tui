package main

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

const defaultAccent = "#CBA6F7"

type theme struct {
	name   string
	accent string
}

var themes = []theme{
	{name: "catppuccin", accent: "#CBA6F7"},
	{name: "cyberpunk", accent: "#00FF9C"},
	{name: "gruvbox", accent: "#FE8019"},
	{name: "nord", accent: "#88C0D0"},
	{name: "netflix", accent: "#E50914"},
}

var currentTheme int
var currentAccent = defaultAccent

// Estilos

var (
	styleTitle      = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color(defaultAccent))
	styleVersion    = lipgloss.NewStyle().Foreground(lipgloss.Color("#AAAAAA"))
	styleDivider    = lipgloss.NewStyle().Foreground(lipgloss.Color(defaultAccent))
	styleFooterKey  = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color(defaultAccent))
	styleFooterDesc = lipgloss.NewStyle().Foreground(lipgloss.Color("#AAAAAA"))
	styleFooterSep  = lipgloss.NewStyle().Foreground(lipgloss.Color("#555555"))
	styleDir        = lipgloss.NewStyle().Foreground(lipgloss.Color(defaultAccent))
	styleProgress   = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFCC00"))
	styleWatched    = lipgloss.NewStyle().Faint(true).Foreground(lipgloss.Color("#5FAF5F"))
	styleNormal     = lipgloss.NewStyle().Foreground(lipgloss.Color("#DDDDDD"))
	styleMeta       = lipgloss.NewStyle().Foreground(lipgloss.Color("#AAAAAA"))
	styleSuccess    = lipgloss.NewStyle().Foreground(lipgloss.Color("#5FAF5F"))
	styleError      = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF5F5F"))
	styleLoading    = lipgloss.NewStyle().Faint(true).Foreground(lipgloss.Color("#FFCC00"))
	styleUpdate     = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#5FAF5F"))
)

func applyTheme(accent string) {
	c := lipgloss.Color(accent)
	styleTitle = styleTitle.Foreground(c)
	styleDivider = styleDivider.Foreground(c)
	styleFooterKey = styleFooterKey.Foreground(c)
	styleDir = styleDir.Foreground(c)
	currentAccent = accent
}

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
