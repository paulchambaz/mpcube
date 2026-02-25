package main

import (
	"fmt"
	"math"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func (ps *PlayerState) renderAlbumPanel(width, height int) string {
	var borderColor lipgloss.Color
	if ps.onAlbum {
		borderColor = lipgloss.Color("9")
	} else {
		borderColor = lipgloss.Color("8")
	}

	title := " Album "

	remainingWidth := width - len(title)
	leftPadding := remainingWidth / 2
	rightPadding := remainingWidth - leftPadding

	topBorder := "┌" + strings.Repeat("─", max(0, leftPadding)) + title + strings.Repeat("─", max(0, rightPadding)) + "┐"

	topBorderStyle := lipgloss.NewStyle().
		Foreground(borderColor)

	if ps.onAlbum {
		topBorderStyle = topBorderStyle.Bold(true)
	} else {
		topBorderStyle = topBorderStyle.Bold(false)
	}
	styledTopBorder := topBorderStyle.Render(topBorder)

	contentStyle := lipgloss.NewStyle().
		Width(width).
		Height(height).
		Border(lipgloss.NormalBorder()).
		BorderTop(false).
		BorderForeground(borderColor)

	var content []string
	visibleHeight := height

	albums := ps.musicData.Albums
	for i := ps.albumOffset; i < len(albums) && i < ps.albumOffset+visibleHeight; i++ {
		albumStyle := lipgloss.NewStyle()
		isPlaying := ps.albumPlaying != nil && *ps.albumPlaying == i
		isSelected := ps.albumSelected == i
		switch {
		case isPlaying && isSelected:
			albumStyle = albumStyle.Reverse(true).Foreground(lipgloss.Color("6"))
		case isPlaying:
			albumStyle = albumStyle.Reverse(true).Foreground(lipgloss.Color("2"))
		case isSelected:
			albumStyle = albumStyle.Reverse(true).Foreground(lipgloss.Color("12"))
		default:
			albumStyle = albumStyle.Foreground(lipgloss.Color("8"))
		}
		name := albums[i].Album
		if len(name) > width {
			name = name[:width]
		}
		line := albumStyle.Render(fmt.Sprintf("%-*s", width, name))
		content = append(content, line)
	}

	contentArea := contentStyle.Render(strings.Join(content, "\n"))
	return styledTopBorder + "\n" + contentArea
}

func (ps *PlayerState) renderTrackPanel(width, height int) string {
	var borderColor lipgloss.Color
	if !ps.onAlbum {
		borderColor = lipgloss.Color("9")
	} else {
		borderColor = lipgloss.Color("8")
	}

	title := " Title "

	remainingWidth := width - len(title)
	leftPadding := 4
	rightPadding := remainingWidth - leftPadding

	topBorder := "┌" + strings.Repeat("─", max(0, leftPadding)) + title + strings.Repeat("─", max(0, rightPadding)) + "┐"

	topBorderStyle := lipgloss.NewStyle().
		Foreground(borderColor)

	if !ps.onAlbum {
		topBorderStyle = topBorderStyle.Bold(true)
	} else {
		topBorderStyle = topBorderStyle.Bold(false)
	}
	styledTopBorder := topBorderStyle.Render(topBorder)

	contentStyle := lipgloss.NewStyle().
		Width(width).
		Height(height).
		Border(lipgloss.NormalBorder()).
		BorderTop(false).
		BorderForeground(borderColor)

	var content []string
	visibleHeight := height

	album := ps.musicData.Albums[ps.albumSelected]
	tracks := album.Songs
	for i := ps.trackOffset; i < len(tracks) && i < ps.trackOffset+visibleHeight; i++ {
		trackStyle := lipgloss.NewStyle()

		isPlaying := ps.albumPlaying != nil && *ps.albumPlaying == ps.albumSelected && ps.trackPlaying != nil && *ps.trackPlaying == i
		isSelected := !ps.onAlbum && ps.trackSelected == i
		switch {
		case isPlaying && isSelected:
			trackStyle = trackStyle.Reverse(true).Foreground(lipgloss.Color("6"))
		case isPlaying:
			trackStyle = trackStyle.Reverse(true).Foreground(lipgloss.Color("2"))
		case isSelected:
			trackStyle = trackStyle.Reverse(true).Foreground(lipgloss.Color("12"))
		default:
			trackStyle = trackStyle.Foreground(lipgloss.Color("8"))
		}

		track := tracks[i]

		left := fmt.Sprintf(" %2d -", track.Track)
		seconds := int(track.Duration.Seconds())
		right := fmt.Sprintf("%02d:%02d %s ", seconds/60, seconds%60, album.Artist)
		title := track.Title
		maxWidth := max(0, width-len(left)-len(right)-2)
		if len(title) > maxWidth {
			title = title[:maxWidth]
		}
		center := fmt.Sprintf("%-*s", maxWidth, title)

		fullLine := fmt.Sprintf("%s %s %s", left, center, right)
		if len(fullLine) > width {
			fullLine = fullLine[:width]
		}
		line := trackStyle.Render(fmt.Sprintf("%-*s", width, fullLine))
		content = append(content, line)
	}

	contentArea := contentStyle.Render(strings.Join(content, "\n"))
	return styledTopBorder + "\n" + contentArea
}

func (ps *PlayerState) renderSearchBar(width int) string {
	style := lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
	cursorStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("2")).Bold(true)

	query := ps.searchQuery
	maxQuery := width - 4
	if len(query) > maxQuery {
		query = query[len(query)-maxQuery:]
	}

	text := style.Render(" / "+query) + cursorStyle.Render("█")
	return lipgloss.NewStyle().Width(width).Height(1).Render(text)
}

func (ps *PlayerState) renderSearchingBar(width int) string {
	style := lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
	countStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("2"))

	count := fmt.Sprintf("[%d/%d]", ps.searchMatchIdx+1, len(ps.searchMatches))
	query := ps.searchQuery
	prefix := " " + count + " / "
	maxQuery := width - len(prefix)
	if len(query) > maxQuery {
		query = query[len(query)-maxQuery:]
	}

	text := " " + countStyle.Render(count) + style.Render(" / "+query)
	return lipgloss.NewStyle().Width(width).Height(1).Render(text)
}

func (ps *PlayerState) renderVolumeBar(width int) string {
	contentStyle := lipgloss.NewStyle().
		Width(width).
		Height(1).
		Foreground(lipgloss.Color("8"))

	start := 5
	end := width - 6
	barWidth := end - start + 1
	ratio := float32(ps.volume) / 100.0
	cursorPos := int(ratio * float32(barWidth-1))

	bar := strings.Repeat("─", max(0, cursorPos)) + "█" + strings.Repeat("─", max(0, barWidth-cursorPos-1))

	return contentStyle.Render(fmt.Sprintf(" Vol %s %d%%", bar, ps.volume))
}

func (ps *PlayerState) renderProgressBar(width int) string {
	contentStyle := lipgloss.NewStyle().
		Width(width).
		Height(1).
		Foreground(lipgloss.Color("8"))

	if ps.albumPlaying == nil || ps.trackPlaying == nil {
		return contentStyle.Render(fmt.Sprintf(" 00:00 %s 00:00 ", strings.Repeat("─", max(0, width-14))))
	}

	track := ps.musicData.Albums[*ps.albumPlaying].Songs[*ps.trackPlaying]
	current := int(ps.position.Seconds())
	total := int(track.Duration.Seconds())

	start := 7
	end := width - 8
	barWidth := end - start + 1
	ratio := float64(current) / float64(total)
	cursorPos := int(math.Round(ratio * float64(barWidth-1)))

	bar := strings.Repeat("─", max(0, cursorPos)) + "█" + strings.Repeat("─", max(barWidth-cursorPos-1, 0))

	return contentStyle.Render(fmt.Sprintf(" %02d:%02d %s %02d:%02d ", current/60, current%60, bar, total/60, total%60))
}

func (ps *PlayerState) renderStatusBar() string {
	contentStyle := lipgloss.NewStyle().
		Width(8).
		Height(2)

	normalStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
	boldStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("2")).Bold(true)

	var shuffleText string
	if ps.shuffle {
		shuffleText = boldStyle.Render("shuffle")
	} else {
		shuffleText = normalStyle.Render("shuffle")
	}

	var repeatText string
	if ps.repeat {
		repeatText = boldStyle.Render("repeat")
	} else {
		repeatText = normalStyle.Render("repeat")
	}

	return contentStyle.Render(fmt.Sprintf("%s\n %s", shuffleText, repeatText))
}
