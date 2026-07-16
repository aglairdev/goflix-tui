package theme

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

const DefaultAccent = "#CBA6F7"

type Theme struct {
	Name   string
	Accent string
}

var Themes = []Theme{
	{Name: "catppuccin", Accent: "#CBA6F7"},
	{Name: "cyberpunk", Accent: "#00FF9C"},
	{Name: "gruvbox", Accent: "#FE8019"},
	{Name: "nord", Accent: "#88C0D0"},
	{Name: "netflix", Accent: "#E50914"},
}

var Current int
var CurrentAccent = DefaultAccent

// Estilos

var (
	Title       = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color(DefaultAccent))
	VersionText = lipgloss.NewStyle().Foreground(lipgloss.Color("#AAAAAA"))
	Divider     = lipgloss.NewStyle().Foreground(lipgloss.Color(DefaultAccent))
	FooterKey   = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color(DefaultAccent))
	FooterDesc  = lipgloss.NewStyle().Foreground(lipgloss.Color("#AAAAAA"))
	FooterSep   = lipgloss.NewStyle().Foreground(lipgloss.Color("#555555"))
	Dir         = lipgloss.NewStyle().Foreground(lipgloss.Color(DefaultAccent))
	Progress    = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFCC00"))
	Watched     = lipgloss.NewStyle().Faint(true).Foreground(lipgloss.Color("#5FAF5F"))
	Normal      = lipgloss.NewStyle().Foreground(lipgloss.Color("#DDDDDD"))
	Meta        = lipgloss.NewStyle().Foreground(lipgloss.Color("#AAAAAA"))
	Success     = lipgloss.NewStyle().Foreground(lipgloss.Color("#5FAF5F"))
	Error       = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF5F5F"))
	Loading     = lipgloss.NewStyle().Faint(true).Foreground(lipgloss.Color("#FFCC00"))
	Update      = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#5FAF5F"))
)

func Apply(accent string) {
	c := lipgloss.Color(accent)
	Title = Title.Foreground(c)
	Divider = Divider.Foreground(c)
	FooterKey = FooterKey.Foreground(c)
	Dir = Dir.Foreground(c)
	CurrentAccent = accent
}

func RenderFooter(raw string) string {
	parts := strings.Split(raw, "  |  ")
	rendered := make([]string, len(parts))
	for i, part := range parts {
		kv := strings.SplitN(part, ": ", 2)
		if len(kv) == 2 {
			rendered[i] = FooterKey.Render(strings.TrimSpace(kv[0])) +
				FooterDesc.Render(": "+kv[1])
		} else {
			rendered[i] = FooterDesc.Render(part)
		}
	}
	return "  " + strings.Join(rendered, FooterSep.Render("  |  "))
}
