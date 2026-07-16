package ui

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/aglairdev/goflix/internal/i18n"
	"github.com/aglairdev/goflix/internal/theme"
	"github.com/aglairdev/goflix/internal/version"
)

func (m Model) View() string {
	if m.quitting {
		return "  " + theme.VersionText.Render(version.AppName+" ꕤ - "+i18n.T("bye")) + "\n"
	}

	header := "  " + theme.Title.Render(version.AppName+" ꕤ") + "  " + theme.VersionText.Render(version.Version)
	if m.screen == screenFiles && m.curDir != "" && m.pendingDir == "" {
		header += "  " + theme.Dir.Render(filepath.Base(m.curDir)+"/")
	}

	var body string
	switch m.screen {
	case screenMain:
		if len(m.dirs) == 0 {
			body = "\n  " + theme.Normal.Render(i18n.T("no_dirs")) + "\n"
		} else {
			body = "\n" + m.mainList.View() + "\n"
		}
		used := strings.Count(body, "\n") + 2
		if rem := m.height - used - 3; rem > 0 {
			body += strings.Repeat("\n", rem)
		}
	case screenFiles:
		if len(m.fileList.Items()) == 0 {
			body = "\n  " + theme.Meta.Render(i18n.T("no_video")) + "\n"
		} else {
			body = "\n" + m.fileList.View() + "\n"
		}
		used := strings.Count(body, "\n") + 2
		if rem := m.height - used - 3; rem > 0 {
			body += strings.Repeat("\n", rem)
		}
	case screenInput:
		body = "  " + theme.Normal.Render(i18n.T("prompt_dir")) + "\n\n" +
			"  " + m.input.View() + "\n\n"
		used := strings.Count(body, "\n") + 2
		if rem := m.height - used - 3; rem > 0 {
			body += strings.Repeat("\n", rem)
		}
	case screenRename:
		body = "\n  " + theme.Normal.Render(i18n.T("rename_label")+": "+filepath.Base(m.renameTarget)) + "\n\n" +
			"  " + m.input.View() + "\n\n"
		used := strings.Count(body, "\n") + 2
		if rem := m.height - used - 3; rem > 0 {
			body += strings.Repeat("\n", rem)
		}

	case screenLoading:
		body = "  " + theme.Loading.Render("⟳  "+i18n.T("loading")) + "\n"
	case screenUpdate:
		body = "\n  " + fmt.Sprintf("ꕤ %s: %s  (%s: %s)\n\n",
			i18n.T("update_available"), theme.Update.Render(m.latestVer),
			i18n.T("update_current"), theme.Error.Render(version.Version)) +
			"  " + theme.VersionText.Render(i18n.T("update_prompt")) + "\n"
	}

	footer := ""
	if m.screen == screenFiles && hasDirItems(m.fileList.Items()) {
		footer = theme.RenderFooter(i18n.T(footerKeyDirs)) + "\n"
	} else if key, ok := footerKey[m.screen]; ok {
		footer = theme.RenderFooter(i18n.T(key)) + "\n"
	}

	flash := "\n"
	if m.flash != "" {
		style, prefix := theme.Success, "✓  "
		if m.flashErr {
			style, prefix = theme.Error, "✗  "
		}
		flash = "  " + style.Render(prefix+m.flash) + "\n"
	}

	debugSection := ""
	if m.debugFlash != "" {
		debugSection = "  " + m.debugFlash + "\n"
	}

	s := header + "\n" +
		theme.Divider.Render(strings.Repeat("─", m.width)) + "\n" +
		body + footer +
		theme.Divider.Render(strings.Repeat("─", m.width)) + "\n" +
		flash + debugSection
	lines := strings.Count(s, "\n")
	if rem := m.height - lines; rem > 0 {
		s += strings.Repeat("\n", rem)
	}
	return s
}
