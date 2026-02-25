package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

type helpEntry struct {
	key  string
	desc string
}

var searchHelpEntries = []helpEntry{
	{"ctrl+c", "quit"},
	{"esc", "cancel"},
	{"enter", "confirm"},
	{"backspace", "delete"},
	{"?", "help"},
}

func (ps *PlayerState) renderHelp() string {
	keyStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("2"))
	descStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
	titleStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("2")).Bold(true)

	var entries []helpEntry
	var modeLabel string
	switch ps.helpForMode {
	case ModeNormal:
		modeLabel = "Normal"
		entries = normalKeys.bindings()
	case ModeSearch:
		modeLabel = "Search"
		entries = searchHelpEntries
	case ModeSearching:
		modeLabel = "Searching"
		entries = searchingKeys.bindings()
	case ModeEdit, ModeEditInput, ModeEditSearch, ModeEditSearching:
		modeLabel = "Edit"
		entries = editCenterKeys.bindings()
	}

	maxKeyWidth := 0
	maxDescWidth := 0
	for _, e := range entries {
		if len(e.key) > maxKeyWidth {
			maxKeyWidth = len(e.key)
		}
		if len(e.desc) > maxDescWidth {
			maxDescWidth = len(e.desc)
		}
	}

	colGap := 6
	colWidth := maxKeyWidth + 2 + maxDescWidth
	availableHeight := ps.windowHeight - 4
	entriesPerCol := max(1, availableHeight/2)
	minCols := 1
	if len(entries) > 10 {
		minCols = 2
	}
	numCols := max(minCols, (len(entries)+entriesPerCol-1)/entriesPerCol)
	entriesPerCol = (len(entries) + numCols - 1) / numCols

	cellStyle := lipgloss.NewStyle().Width(colWidth)

	columns := make([][]string, numCols)
	for i, e := range entries {
		col := i / entriesPerCol
		padKey := fmt.Sprintf("%*s", maxKeyWidth, e.key)
		line := keyStyle.Render(padKey) + "  " + descStyle.Render(e.desc)
		columns[col] = append(columns[col], cellStyle.Render(line))
		columns[col] = append(columns[col], cellStyle.Render(""))
	}

	maxColHeight := 0
	for _, col := range columns {
		if len(col) > maxColHeight {
			maxColHeight = len(col)
		}
	}
	for i := range columns {
		for len(columns[i]) < maxColHeight {
			columns[i] = append(columns[i], cellStyle.Render(""))
		}
	}

	gap := strings.Repeat(" ", colGap)
	var contentLines []string
	for row := 0; row < maxColHeight; row++ {
		var parts []string
		for _, col := range columns {
			parts = append(parts, col[row])
		}
		contentLines = append(contentLines, strings.Join(parts, gap))
	}

	titleText := modeLabel + " Mode"
	title := titleStyle.Render(titleText)
	titlePad := max(0, (ps.windowWidth-len(titleText))/2)

	allLines := []string{strings.Repeat(" ", titlePad) + title, ""}
	allLines = append(allLines, contentLines...)

	totalHeight := len(allLines)
	topPad := max(0, (ps.windowHeight-totalHeight)/2)

	totalWidth := numCols*colWidth + (numCols-1)*colGap
	leftPad := max(0, (ps.windowWidth-totalWidth)/2)
	padStr := strings.Repeat(" ", leftPad)

	var result []string
	for i := 0; i < topPad; i++ {
		result = append(result, "")
	}
	result = append(result, allLines[0])
	result = append(result, allLines[1])
	for _, line := range allLines[2:] {
		result = append(result, padStr+line)
	}

	return lipgloss.NewStyle().
		Width(ps.windowWidth).
		Height(ps.windowHeight).
		Render(strings.Join(result, "\n"))
}

func (ps *PlayerState) View() string {
	if ps.windowWidth == 0 || ps.windowHeight == 0 {
		return "Loading..."
	}

	if ps.mode == ModeHelp {
		return ps.renderHelp()
	}

	if ps.mode == ModeEdit || ps.mode == ModeEditInput || ps.mode == ModeEditSearch || ps.mode == ModeEditSearching {
		return ps.renderEditView()
	}

	leftWidth := 0
	if ps.windowWidth > ps.config.WideThreshold {
		leftWidth = ps.config.AlbumWidth - 2
	} else {
		leftWidth = 2*ps.windowWidth/5 - 2
	}
	rightWidth := ps.windowWidth - leftWidth - 4
	panelHeight := ps.windowHeight - 4

	leftPanel := ps.renderAlbumPanel(leftWidth, panelHeight)
	rightPanel := ps.renderTrackPanel(rightWidth, panelHeight)

	mainView := lipgloss.JoinHorizontal(0.0, leftPanel, rightPanel)

	volumeWidth := 0
	if ps.windowWidth > ps.config.VolumeBarThreshold {
		volumeWidth = ps.config.VolumeBarWidth
	} else {
		volumeWidth = ps.windowWidth / 3
	}

	infoView := ps.renderInfoBar(ps.windowWidth - 8)
	volumeView := ps.renderVolumeBar(volumeWidth)
	barView := ps.renderProgressBar(ps.windowWidth - volumeWidth - 8)
	statusView := ps.renderStatusBar()

	barsView := lipgloss.JoinHorizontal(0.0, volumeView, barView)
	subView := lipgloss.JoinVertical(0.0, infoView, barsView)
	bottomView := lipgloss.JoinHorizontal(0.0, subView, statusView)

	return lipgloss.JoinVertical(0.0, mainView, bottomView)
}

func (ps *PlayerState) renderInfoBar(width int) string {
	switch ps.mode {
	case ModeSearch:
		return ps.renderSearchBar(width)
	case ModeSearching:
		return ps.renderSearchingBar(width)
	}

	contentStyle := lipgloss.NewStyle().
		Width(width).
		Height(1)
	normalStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
	boldStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("2"))

	if ps.albumPlaying == nil || ps.trackPlaying == nil {
		return normalStyle.Render(" Not playing")
	}

	album := ps.musicData.Albums[*ps.albumPlaying]
	track := album.Songs[*ps.trackPlaying]

	var playingText string
	if ps.playing {
		playingText = "Playing"
	} else {
		playingText = "Paused"
	}

	parts := []string{" ", playingText, " ", track.Title, " by ", album.Artist, " from ", album.Album}
	styles := []lipgloss.Style{normalStyle, normalStyle, normalStyle, boldStyle, normalStyle, boldStyle, normalStyle, boldStyle}

	currentLen := 0
	result := ""

	for i, part := range parts {
		if currentLen+len(part) <= width {
			result += styles[i].Render(part)
			currentLen += len(part)
		} else {
			remaining := width - 1 - currentLen
			if remaining > 0 {
				result += styles[i].Render(part[:remaining])
			}
			break
		}
	}

	return contentStyle.Render(fmt.Sprintf("%-*s", width, result))
}

func (ps *PlayerState) resize() {
	panelHeight := ps.windowHeight - 4
	padding := min(ps.config.ScrollPadding, panelHeight/4)

	if ps.albumSelected < ps.albumOffset+padding {
		ps.albumOffset = ps.albumSelected - padding
	}

	if ps.albumSelected >= ps.albumOffset+panelHeight-padding {
		ps.albumOffset = ps.albumSelected - panelHeight + 1 + padding
	}

	ps.albumOffset = max(ps.albumOffset, 0)
	ps.albumOffset = min(ps.albumOffset, max(0, len(ps.musicData.Albums)-panelHeight))

	if ps.trackSelected < ps.trackOffset+padding {
		ps.trackOffset = ps.trackSelected - padding
	}

	if ps.trackSelected >= ps.trackOffset+panelHeight-padding {
		ps.trackOffset = ps.trackSelected - panelHeight + 1 + padding
	}

	ps.trackOffset = max(ps.trackOffset, 0)
	ps.trackOffset = min(ps.trackOffset, max(0, len(ps.musicData.Albums[ps.albumSelected].Songs)-panelHeight))
}
