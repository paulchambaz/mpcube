package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func (ps *PlayerState) renderEditView() string {
	sideWidth := ps.editSideWidth()
	centerWidth := ps.windowWidth - sideWidth*2 - 6

	titlesPanelHeight := (ps.windowHeight - 4) / 2
	albumsPanelHeight := ps.windowHeight - 4 - titlesPanelHeight

	centerHeight := ps.windowHeight - 1 - 2

	albumsPanel := ps.renderEditAlbumsPanel(sideWidth, albumsPanelHeight)
	titlesPanel := ps.renderEditTitlesPanel(sideWidth, titlesPanelHeight)
	leftColumn := lipgloss.JoinVertical(0.0, albumsPanel, titlesPanel)

	rightBottomHeight := (ps.windowHeight - 6) / 3
	rightMiddleHeight := rightBottomHeight
	rightTopHeight := ps.windowHeight - 6 - rightBottomHeight - rightMiddleHeight
	rightTopPanel := ps.renderEditEmptyPanel(" Metadata ", sideWidth, rightTopHeight)
	rightMiddlePanel := ps.renderEditEmptyPanel(" Cover ", sideWidth, rightMiddleHeight)
	rightBottomPanel := ps.renderEditEmptyPanel(" Download ", sideWidth, rightBottomHeight)
	rightColumn := lipgloss.JoinVertical(0.0, rightTopPanel, rightMiddlePanel, rightBottomPanel)

	centerPanel := ps.renderEditCenterPanel(centerWidth, centerHeight)
	shortcutBar := ps.renderEditShortcutBar(centerWidth + 2)
	centerColumn := lipgloss.JoinVertical(0.0, centerPanel, shortcutBar)

	return lipgloss.JoinHorizontal(0.0, leftColumn, centerColumn, rightColumn)
}

func (ps *PlayerState) editSideWidth() int {
	if ps.windowWidth > ps.config.WideThreshold {
		return ps.config.AlbumWidth - 2
	}
	return ps.windowWidth/5 - 2
}

func (ps *PlayerState) renderEditAlbumsPanel(width, height int) string {
	focused := ps.editFocus == EditFocusAlbums
	var borderColor lipgloss.Color
	if focused {
		borderColor = lipgloss.Color("9")
	} else {
		borderColor = lipgloss.Color("8")
	}

	title := " Albums "
	remainingWidth := width - len(title)
	leftPad := remainingWidth / 2
	rightPad := remainingWidth - leftPad
	topBorder := "┌" + strings.Repeat("─", max(0, leftPad)) + title + strings.Repeat("─", max(0, rightPad)) + "┐"
	topStyle := lipgloss.NewStyle().Foreground(borderColor)
	if focused {
		topStyle = topStyle.Bold(true)
	}

	contentStyle := lipgloss.NewStyle().
		Width(width).
		Height(height).
		Border(lipgloss.NormalBorder()).
		BorderTop(false).
		BorderForeground(borderColor)

	var content []string
	albums := ps.musicData.Albums
	visibleHeight := height

	for i := ps.albumOffset; i < len(albums) && i < ps.albumOffset+visibleHeight; i++ {
		style := lipgloss.NewStyle()
		isSelected := ps.albumSelected == i
		if isSelected && focused {
			style = style.Reverse(true).Foreground(lipgloss.Color("12"))
		} else if isSelected {
			style = style.Reverse(true).Foreground(lipgloss.Color("6"))
		} else {
			style = style.Foreground(lipgloss.Color("8"))
		}
		name := albums[i].Album
		if len(name) > width {
			name = name[:width]
		}
		content = append(content, style.Render(fmt.Sprintf("%-*s", width, name)))
	}

	return topStyle.Render(topBorder) + "\n" + contentStyle.Render(strings.Join(content, "\n"))
}

func (ps *PlayerState) renderEditTitlesPanel(width, height int) string {
	focused := ps.editFocus == EditFocusTitles
	var borderColor lipgloss.Color
	if focused {
		borderColor = lipgloss.Color("9")
	} else {
		borderColor = lipgloss.Color("8")
	}

	title := " Titles "
	remainingWidth := width - len(title)
	leftPad := remainingWidth / 2
	rightPad := remainingWidth - leftPad
	topBorder := "┌" + strings.Repeat("─", max(0, leftPad)) + title + strings.Repeat("─", max(0, rightPad)) + "┐"
	topStyle := lipgloss.NewStyle().Foreground(borderColor)
	if focused {
		topStyle = topStyle.Bold(true)
	}

	contentStyle := lipgloss.NewStyle().
		Width(width).
		Height(height).
		Border(lipgloss.NormalBorder()).
		BorderTop(false).
		BorderForeground(borderColor)

	var content []string
	for i := ps.editTitleOffset; i < len(ps.editTracks) && i < ps.editTitleOffset+height; i++ {
		style := lipgloss.NewStyle()
		isSelected := ps.editTitleIdx == i
		if isSelected && focused {
			style = style.Reverse(true).Foreground(lipgloss.Color("12"))
		} else if isSelected {
			style = style.Reverse(true).Foreground(lipgloss.Color("6"))
		} else {
			style = style.Foreground(lipgloss.Color("8"))
		}
		track := ps.editTracksOrig[i]
		line := fmt.Sprintf(" %2s - %s", track.Track, track.Title)
		if len(line) > width {
			line = line[:width]
		}
		content = append(content, style.Render(fmt.Sprintf("%-*s", width, line)))
	}

	return topStyle.Render(topBorder) + "\n" + contentStyle.Render(strings.Join(content, "\n"))
}

func (ps *PlayerState) renderEditEmptyPanel(title string, width, height int) string {
	borderColor := lipgloss.Color("8")

	remainingWidth := width - len(title)
	leftPad := remainingWidth / 2
	rightPad := remainingWidth - leftPad
	topBorder := "┌" + strings.Repeat("─", max(0, leftPad)) + title + strings.Repeat("─", max(0, rightPad)) + "┐"
	topStyle := lipgloss.NewStyle().Foreground(borderColor)

	contentStyle := lipgloss.NewStyle().
		Width(width).
		Height(height).
		Border(lipgloss.NormalBorder()).
		BorderTop(false).
		BorderForeground(borderColor)

	return topStyle.Render(topBorder) + "\n" + contentStyle.Render("")
}

func (ps *PlayerState) renderEditCenterPanel(width, height int) string {
	focused := ps.editFocus == EditFocusCenter
	var borderColor lipgloss.Color
	if focused {
		borderColor = lipgloss.Color("9")
	} else {
		borderColor = lipgloss.Color("8")
	}

	title := " Edit "
	remainingWidth := width - len(title)
	leftPad := 4
	rightPad := remainingWidth - leftPad
	topBorder := "┌" + strings.Repeat("─", max(0, leftPad)) + title + strings.Repeat("─", max(0, rightPad)) + "┐"
	topStyle := lipgloss.NewStyle().Foreground(borderColor)
	if focused {
		topStyle = topStyle.Bold(true)
	}

	contentStyle := lipgloss.NewStyle().
		Width(width).
		Height(height).
		Border(lipgloss.NormalBorder()).
		BorderTop(false).
		BorderForeground(borderColor)

	labelStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("2"))
	modStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("3"))
	albumLabels := [5]string{"Album", "Artist", "Date", "Dir", "Cover"}
	trackFieldLabels := [3]string{"Track", "Title", "File"}

	type renderLine struct {
		fieldIdx int
		isSep    bool
	}

	var lines []renderLine
	for i := 0; i < editAlbumFieldCount; i++ {
		lines = append(lines, renderLine{fieldIdx: i})
	}
	lines = append(lines, renderLine{isSep: true, fieldIdx: -1})
	for ti := 0; ti < len(ps.editTracks); ti++ {
		baseIdx := editAlbumFieldCount + ti*3
		for fi := 0; fi < 3; fi++ {
			lines = append(lines, renderLine{fieldIdx: baseIdx + fi})
		}
		if ti < len(ps.editTracks)-1 {
			lines = append(lines, renderLine{fieldIdx: -1})
		}
	}

	startLine := ps.editFieldOffset

	var content []string
	maxLabel := 8
	for li := startLine; li < len(lines) && len(content) < height; li++ {
		rl := lines[li]

		if rl.isSep {
			sep := " " + strings.Repeat("─", max(0, width-2))
			content = append(content, lipgloss.NewStyle().Foreground(lipgloss.Color("8")).Render(sep))
			continue
		}

		if rl.fieldIdx == -1 {
			content = append(content, "")
			continue
		}

		selected := focused && rl.fieldIdx == ps.editFieldIdx
		editInput := ps.mode == ModeEditInput && rl.fieldIdx == ps.editFieldIdx

		var label string
		var value string

		if rl.fieldIdx < editAlbumFieldCount {
			label = albumLabels[rl.fieldIdx]
			value = ps.editAlbum[rl.fieldIdx]
		} else {
			ti := ps.editTrackIdx(rl.fieldIdx)
			fi := ps.editTrackFieldIdx(rl.fieldIdx)
			label = trackFieldLabels[fi]
			switch fi {
			case 0:
				value = ps.editTracks[ti].Track
			case 1:
				value = ps.editTracks[ti].Title
			case 2:
				value = ps.editTracks[ti].File
			}
		}

		mod := ""
		if ps.editIsModified(rl.fieldIdx) {
			mod = " [mod]"
		}

		valWidth := max(0, width-maxLabel-3-len(mod))
		var line string
		if editInput {
			buf := ps.editInputBuf
			pos := ps.editInputPos
			visStart := 0
			visLen := valWidth - 1
			if pos > visStart+visLen {
				visStart = pos - visLen
			}
			if pos < visStart {
				visStart = pos
			}
			visEnd := min(len(buf), visStart+visLen)
			cursorInVis := pos - visStart

			before := buf[visStart : visStart+cursorInVis]
			var cursorChar string
			if pos < len(buf) {
				cursorChar = string(buf[pos])
			} else {
				cursorChar = " "
			}
			after := ""
			if pos+1 < len(buf) && pos+1 <= visEnd {
				after = buf[pos+1 : visEnd]
			}

			dimStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
			cursorStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("2")).Bold(true).Reverse(true)
			labelStr := labelStyle.Render(fmt.Sprintf(" %-*s", maxLabel, label))
			valStr := dimStyle.Render(before) + cursorStyle.Render(cursorChar) + dimStyle.Render(after)
			pad := max(0, valWidth-cursorInVis-1-len(after))
			line = labelStr + " " + valStr + strings.Repeat(" ", pad)
		} else if selected {
			selStyle := lipgloss.NewStyle().Reverse(true).Foreground(lipgloss.Color("12"))
			line = selStyle.Render(fmt.Sprintf(" %-*s %-*s", maxLabel, label, valWidth, value))
		} else {
			labelStr := labelStyle.Render(fmt.Sprintf(" %-*s", maxLabel, label))
			valStr := lipgloss.NewStyle().Foreground(lipgloss.Color("8")).Render(fmt.Sprintf("%-*s", valWidth, value))
			line = labelStr + " " + valStr
		}

		if mod != "" {
			line += modStyle.Render(mod)
		}

		content = append(content, line)
	}

	return topStyle.Render(topBorder) + "\n" + contentStyle.Render(strings.Join(content, "\n"))
}

func (ps *PlayerState) editShortcutBarPad(width int) int {
	shortcuts := [][2]string{
		{"i", "edit"}, {"v", "editor"}, {"s", "sync"},
		{"r", "revert"}, {"q", "quit"}, {"U", "apply"},
	}
	currentLen := 1
	for _, sc := range shortcuts {
		currentLen += len(sc[0]) + 1 + len(sc[1]) + 2
	}
	return max(0, (width-currentLen)/2)
}

func (ps *PlayerState) renderEditSearchBar(width int) string {
	style := lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
	cursorStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("2")).Bold(true)

	pad := ps.editShortcutBarPad(width)
	query := ps.searchQuery
	maxQuery := width - pad - 3
	if len(query) > maxQuery {
		query = query[len(query)-maxQuery:]
	}

	text := strings.Repeat(" ", pad) + style.Render("/ "+query) + cursorStyle.Render("█")
	return lipgloss.NewStyle().Width(width).Height(1).Render(text)
}

func (ps *PlayerState) renderEditSearchingBar(width int) string {
	style := lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
	countStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("2"))

	pad := ps.editShortcutBarPad(width)
	count := fmt.Sprintf("[%d/%d]", ps.searchMatchIdx+1, len(ps.searchMatches))
	query := ps.searchQuery
	maxQuery := width - pad - len(count) - 4
	if len(query) > maxQuery {
		query = query[len(query)-maxQuery:]
	}

	text := strings.Repeat(" ", pad) + countStyle.Render(count) + style.Render(" / "+query)
	return lipgloss.NewStyle().Width(width).Height(1).Render(text)
}

func (ps *PlayerState) renderEditShortcutBar(width int) string {
	if ps.mode == ModeEditSearch {
		return ps.renderEditSearchBar(width)
	}
	if ps.mode == ModeEditSearching {
		return ps.renderEditSearchingBar(width)
	}

	keyStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("2"))
	descStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("8"))

	type shortcut struct {
		key  string
		desc string
	}

	hasModified := false
	for i := 0; i < ps.editFieldCount(); i++ {
		if ps.editIsModified(i) {
			hasModified = true
			break
		}
	}

	var shortcuts []shortcut
	if ps.mode == ModeEditInput {
		shortcuts = []shortcut{
			{"enter", "confirm"},
			{"esc", "cancel"},
		}
	} else {
		shortcuts = []shortcut{
			{"i", "edit"},
			{"v", "editor"},
			{"s", "sync"},
			{"r", "revert"},
			{"q", "quit"},
			{"U", "apply"},
		}
	}

	var parts []string
	currentLen := 1
	for _, sc := range shortcuts {
		entry := keyStyle.Render(sc.key) + " " + descStyle.Render(sc.desc)
		entryLen := len(sc.key) + 1 + len(sc.desc) + 2
		if currentLen+entryLen > width {
			break
		}
		parts = append(parts, entry)
		currentLen += entryLen
	}

	inner := strings.Join(parts, "  ")
	if hasModified {
		modStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("3"))
		modTag := modStyle.Render("[modified]")
		modLen := len("[modified]") + 2
		if currentLen+modLen <= width {
			inner += "  " + modTag
			currentLen += modLen
		}
	}

	pad := max(0, (width-currentLen)/2)
	text := strings.Repeat(" ", pad) + inner

	return lipgloss.NewStyle().Width(width).Height(1).Render(text)
}

