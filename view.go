package main

import (
	"fmt"
	"path/filepath"
	"strings"
)

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
		used := strings.Count(body, "\n") + 2
		if rem := m.height - used - 3; rem > 0 {
			body += strings.Repeat("\n", rem)
		}
	case screenRename:
		body = "\n  " + styleNormal.Render(t("rename_label")+": "+filepath.Base(m.renameTarget)) + "\n\n" +
			"  " + m.input.View() + "\n\n"
		used := strings.Count(body, "\n") + 2
		if rem := m.height - used - 3; rem > 0 {
			body += strings.Repeat("\n", rem)
		}

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

	debugSection := ""
	if m.debugFlash != "" {
		debugSection = "  " + m.debugFlash + "\n"
	}

	s := header + "\n" +
		styleDivider.Render(strings.Repeat("─", m.width)) + "\n" +
		body + footer +
		styleDivider.Render(strings.Repeat("─", m.width)) + "\n" +
		flash + debugSection
	lines := strings.Count(s, "\n")
	if rem := m.height - lines; rem > 0 {
		s += strings.Repeat("\n", rem)
	}
	return s
}
