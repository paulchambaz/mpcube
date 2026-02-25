package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

type editCenterKeyMap struct {
	quit      key.Binding
	left      key.Binding
	up        key.Binding
	down      key.Binding
	top       key.Binding
	bottom    key.Binding
	edit      key.Binding
	editor    key.Binding
	sync      key.Binding
	revert    key.Binding
	revertAll key.Binding
	playPause key.Binding
	search    key.Binding
	openCover key.Binding
}

var editCenterKeys = editCenterKeyMap{
	quit:      key.NewBinding(key.WithKeys("q", "esc"), key.WithHelp("q", "discard")),
	left:      key.NewBinding(key.WithKeys("h", "left"), key.WithHelp("h", "panels")),
	up:        key.NewBinding(key.WithKeys("up", "k"), key.WithHelp("k", "up")),
	down:      key.NewBinding(key.WithKeys("down", "j"), key.WithHelp("j", "down")),
	top:       key.NewBinding(key.WithKeys("home", "g"), key.WithHelp("g", "top")),
	bottom:    key.NewBinding(key.WithKeys("end", "G"), key.WithHelp("G", "bottom")),
	edit:      key.NewBinding(key.WithKeys("i", "enter"), key.WithHelp("i", "edit field")),
	editor:    key.NewBinding(key.WithKeys("v"), key.WithHelp("v", "editor")),
	sync:      key.NewBinding(key.WithKeys("s"), key.WithHelp("s", "sync filenames")),
	revert:    key.NewBinding(key.WithKeys("r"), key.WithHelp("r", "revert")),
	revertAll: key.NewBinding(key.WithKeys("R"), key.WithHelp("R", "revert all")),
	playPause: key.NewBinding(key.WithKeys(" "), key.WithHelp("space", "play/pause")),
	search:    key.NewBinding(key.WithKeys("/"), key.WithHelp("/", "search")),
	openCover: key.NewBinding(key.WithKeys("o"), key.WithHelp("o", "open cover")),
}

func (k editCenterKeyMap) bindings() []helpEntry {
	return []helpEntry{
		bindingHelp(k.quit),
		bindingHelp(k.left),
		bindingHelp(k.up),
		bindingHelp(k.down),
		bindingHelp(k.top),
		bindingHelp(k.bottom),
		bindingHelp(k.edit),
		bindingHelp(k.editor),
		bindingHelp(k.sync),
		bindingHelp(k.revert),
		bindingHelp(k.revertAll),
		bindingHelp(k.openCover),
		bindingHelp(k.playPause),
		bindingHelp(globalKeys.seekForward),
		bindingHelp(globalKeys.seekBackward),
		bindingHelp(globalKeys.volumeUp),
		bindingHelp(globalKeys.volumeDown),
		{"?", "help"},
	}
}

type editSideKeyMap struct {
	right    key.Binding
	enter    key.Binding
	up       key.Binding
	down     key.Binding
	top      key.Binding
	bottom   key.Binding
	jumpNext key.Binding
	search   key.Binding
}

var editAlbumKeys = editSideKeyMap{
	right:    key.NewBinding(key.WithKeys("l", "right"), key.WithHelp("l", "editor")),
	enter:    key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "editor")),
	up:       key.NewBinding(key.WithKeys("up", "k"), key.WithHelp("k", "up")),
	down:     key.NewBinding(key.WithKeys("down", "j"), key.WithHelp("j", "down")),
	top:      key.NewBinding(key.WithKeys("home", "g"), key.WithHelp("g", "top")),
	bottom:   key.NewBinding(key.WithKeys("end", "G"), key.WithHelp("G", "bottom")),
	jumpNext: key.NewBinding(key.WithKeys("J"), key.WithHelp("J", "titles")),
	search:   key.NewBinding(key.WithKeys("/"), key.WithHelp("/", "search")),
}

var editTitleKeys = editSideKeyMap{
	right:    key.NewBinding(key.WithKeys("l", "right"), key.WithHelp("l", "editor")),
	enter:    key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "editor")),
	up:       key.NewBinding(key.WithKeys("up", "k"), key.WithHelp("k", "up")),
	down:     key.NewBinding(key.WithKeys("down", "j"), key.WithHelp("j", "down")),
	top:      key.NewBinding(key.WithKeys("home", "g"), key.WithHelp("g", "top")),
	bottom:   key.NewBinding(key.WithKeys("end", "G"), key.WithHelp("G", "bottom")),
	jumpNext: key.NewBinding(key.WithKeys("K"), key.WithHelp("K", "albums")),
}

func (ps *PlayerState) handleEdit(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch ps.editFocus {
	case EditFocusCenter:
		return ps.handleEditCenter(msg)
	case EditFocusAlbums:
		return ps.handleEditAlbums(msg)
	case EditFocusTitles:
		return ps.handleEditTitles(msg)
	}
	return ps, nil
}

func (ps *PlayerState) handleEditCenter(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case key.Matches(msg, editCenterKeys.quit):
		ps.exitEditMode()
	case key.Matches(msg, editCenterKeys.left):
		if ps.editIsAlbumField(ps.editFieldIdx) {
			ps.editFocus = EditFocusAlbums
		} else {
			ps.editFocus = EditFocusTitles
		}
	case key.Matches(msg, editCenterKeys.up):
		ps.editCenterMoveUp()
	case key.Matches(msg, editCenterKeys.down):
		ps.editCenterMoveDown()
	case key.Matches(msg, editCenterKeys.top):
		ps.editFieldIdx = 0
		ps.editFieldOffset = 0
	case key.Matches(msg, editCenterKeys.bottom):
		ps.editFieldIdx = ps.editFieldCount() - 1
		ps.editFixCenterOffset()
	case key.Matches(msg, editCenterKeys.edit):
		ps.editInputBuf = ps.editCurrentValue()
		ps.editInputPos = len(ps.editInputBuf)
		ps.mode = ModeEditInput
	case key.Matches(msg, editCenterKeys.editor):
		return ps, ps.editOpenEditor()
	case key.Matches(msg, editCenterKeys.sync):
		ps.editSyncFilenames()
	case key.Matches(msg, editCenterKeys.revert):
		ps.editRevertField()
	case key.Matches(msg, editCenterKeys.revertAll):
		ps.editRevertAll()
	case key.Matches(msg, editCenterKeys.playPause):
		_ = ps.togglePlayPause()
	case key.Matches(msg, editCenterKeys.search):
		ps.editEnterSearch()
	case key.Matches(msg, editCenterKeys.openCover):
		ps.editOpenCover()
	}
	return ps, nil
}

func (ps *PlayerState) handleEditAlbums(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case key.Matches(msg, editCenterKeys.quit):
		ps.exitEditMode()
	case key.Matches(msg, editAlbumKeys.right), key.Matches(msg, editAlbumKeys.enter):
		ps.editFocus = EditFocusCenter
	case key.Matches(msg, editAlbumKeys.up):
		if ps.albumSelected > 0 {
			ps.albumSelected--
			ps.editAlbumFixOffset()
			ps.editLoadAlbum()
		}
	case key.Matches(msg, editAlbumKeys.down):
		if ps.albumSelected < len(ps.musicData.Albums)-1 {
			ps.albumSelected++
			ps.editAlbumFixOffset()
			ps.editLoadAlbum()
		} else {
			ps.editFocus = EditFocusTitles
		}
	case key.Matches(msg, editAlbumKeys.jumpNext):
		ps.editFocus = EditFocusTitles
	case key.Matches(msg, editAlbumKeys.top):
		ps.albumSelected = 0
		ps.albumOffset = 0
		ps.editLoadAlbum()
	case key.Matches(msg, editAlbumKeys.bottom):
		ps.albumSelected = len(ps.musicData.Albums) - 1
		ps.editAlbumFixOffset()
		ps.editLoadAlbum()
	case key.Matches(msg, editAlbumKeys.search):
		ps.editEnterSearch()
	}
	return ps, nil
}

func (ps *PlayerState) handleEditTitles(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case key.Matches(msg, editCenterKeys.quit):
		ps.exitEditMode()
	case key.Matches(msg, editTitleKeys.right), key.Matches(msg, editTitleKeys.enter):
		ps.editFieldIdx = editAlbumFieldCount + ps.editTitleIdx*3
		if ps.editFieldIdx >= ps.editFieldCount() {
			ps.editFieldIdx = ps.editFieldCount() - 1
		}
		ps.editFixCenterOffset()
		ps.editFocus = EditFocusCenter
	case key.Matches(msg, editTitleKeys.up):
		if ps.editTitleIdx > 0 {
			ps.editTitleIdx--
			ps.editFixTitleOffset()
		} else {
			ps.editFocus = EditFocusAlbums
		}
	case key.Matches(msg, editTitleKeys.jumpNext):
		ps.editFocus = EditFocusAlbums
	case key.Matches(msg, editTitleKeys.down):
		if ps.editTitleIdx < len(ps.editTracks)-1 {
			ps.editTitleIdx++
			ps.editFixTitleOffset()
		}
	case key.Matches(msg, editTitleKeys.top):
		ps.editTitleIdx = 0
		ps.editTitleOffset = 0
	case key.Matches(msg, editTitleKeys.bottom):
		ps.editTitleIdx = len(ps.editTracks) - 1
		ps.editFixTitleOffset()
	}
	return ps, nil
}

type editSearchingKeyMap struct {
	nextMatch key.Binding
	prevMatch key.Binding
	confirm   key.Binding
	cancel    key.Binding
	reSearch  key.Binding
	quit      key.Binding
}

var editSearchingKeys = editSearchingKeyMap{
	nextMatch: key.NewBinding(key.WithKeys("n"), key.WithHelp("n", "next match")),
	prevMatch: key.NewBinding(key.WithKeys("p"), key.WithHelp("p", "previous match")),
	confirm:   key.NewBinding(key.WithKeys("enter", "l", "right"), key.WithHelp("enter", "confirm")),
	cancel:    key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "cancel")),
	reSearch:  key.NewBinding(key.WithKeys("/", "i"), key.WithHelp("/", "edit query")),
	quit:      key.NewBinding(key.WithKeys("q"), key.WithHelp("q", "discard")),
}

func (ps *PlayerState) handleEditSearch(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyEscape:
		ps.editCancelSearch()
	case tea.KeyEnter:
		ps.editConfirmSearch()
	case tea.KeyBackspace:
		ps.editSearchBackspace()
	case tea.KeySpace:
		ps.editSearchAddRune(' ')
	case tea.KeyRunes:
		ps.editSearchAddRune(msg.Runes[0])
	}
	return ps, nil
}

func (ps *PlayerState) handleEditSearching(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case key.Matches(msg, editSearchingKeys.nextMatch):
		ps.editNextMatch()
	case key.Matches(msg, editSearchingKeys.prevMatch):
		ps.editPrevMatch()
	case key.Matches(msg, editSearchingKeys.confirm):
		ps.editConfirmSearching()
	case key.Matches(msg, editSearchingKeys.cancel):
		ps.editCancelSearch()
	case key.Matches(msg, editSearchingKeys.reSearch):
		ps.mode = ModeEditSearch
	case key.Matches(msg, editSearchingKeys.quit):
		ps.exitEditMode()
	}
	return ps, nil
}

func (ps *PlayerState) handleEditInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyEscape:
		ps.editInputBuf = ""
		ps.editInputPos = 0
		ps.mode = ModeEdit
	case tea.KeyEnter:
		ps.editSetValue(ps.editInputBuf)
		ps.editInputBuf = ""
		ps.editInputPos = 0
		ps.mode = ModeEdit
	case tea.KeyLeft:
		if ps.editInputPos > 0 {
			ps.editInputPos--
		}
	case tea.KeyRight:
		if ps.editInputPos < len(ps.editInputBuf) {
			ps.editInputPos++
		}
	case tea.KeyHome:
		ps.editInputPos = 0
	case tea.KeyEnd:
		ps.editInputPos = len(ps.editInputBuf)
	case tea.KeyBackspace:
		if ps.editInputPos > 0 {
			ps.editInputBuf = ps.editInputBuf[:ps.editInputPos-1] + ps.editInputBuf[ps.editInputPos:]
			ps.editInputPos--
		}
	case tea.KeyDelete:
		if ps.editInputPos < len(ps.editInputBuf) {
			ps.editInputBuf = ps.editInputBuf[:ps.editInputPos] + ps.editInputBuf[ps.editInputPos+1:]
		}
	case tea.KeySpace:
		ps.editInputBuf = ps.editInputBuf[:ps.editInputPos] + " " + ps.editInputBuf[ps.editInputPos:]
		ps.editInputPos++
	case tea.KeyRunes:
		ps.editInputBuf = ps.editInputBuf[:ps.editInputPos] + string(msg.Runes) + ps.editInputBuf[ps.editInputPos:]
		ps.editInputPos += len(string(msg.Runes))
	}
	return ps, nil
}

func (ps *PlayerState) editCenterMoveUp() {
	if ps.editFieldIdx > 0 {
		ps.editFieldIdx--
		ps.editFixCenterOffset()
	}
}

func (ps *PlayerState) editCenterMoveDown() {
	if ps.editFieldIdx < ps.editFieldCount()-1 {
		ps.editFieldIdx++
		ps.editFixCenterOffset()
	}
}

func (ps *PlayerState) editFixCenterOffset() {
	line := ps.editFieldToLine(ps.editFieldIdx)
	totalLines := ps.editTotalLines()
	panelHeight := ps.editCenterHeight()
	padding := min(ps.config.ScrollPadding, panelHeight/4)

	if line < ps.editFieldOffset+padding {
		ps.editFieldOffset = line - padding
	}
	if line >= ps.editFieldOffset+panelHeight-padding {
		ps.editFieldOffset = line - panelHeight + 1 + padding
	}
	ps.editFieldOffset = max(ps.editFieldOffset, 0)
	ps.editFieldOffset = min(ps.editFieldOffset, max(0, totalLines-panelHeight))
}

func (ps *PlayerState) editAlbumFixOffset() {
	panelHeight := ps.editAlbumsPanelHeight()
	padding := min(ps.config.ScrollPadding, panelHeight/4)

	if ps.albumSelected < ps.albumOffset+padding {
		ps.albumOffset = ps.albumSelected - padding
	}
	if ps.albumSelected >= ps.albumOffset+panelHeight-padding {
		ps.albumOffset = ps.albumSelected - panelHeight + 1 + padding
	}
	ps.albumOffset = max(ps.albumOffset, 0)
	ps.albumOffset = min(ps.albumOffset, max(0, len(ps.musicData.Albums)-panelHeight))
}

func (ps *PlayerState) editFixTitleOffset() {
	panelHeight := ps.editTitlesPanelHeight()
	padding := min(ps.config.ScrollPadding, panelHeight/4)

	if ps.editTitleIdx < ps.editTitleOffset+padding {
		ps.editTitleOffset = ps.editTitleIdx - padding
	}
	if ps.editTitleIdx >= ps.editTitleOffset+panelHeight-padding {
		ps.editTitleOffset = ps.editTitleIdx - panelHeight + 1 + padding
	}
	ps.editTitleOffset = max(ps.editTitleOffset, 0)
	ps.editTitleOffset = min(ps.editTitleOffset, max(0, len(ps.editTracks)-panelHeight))
}

func (ps *PlayerState) editCenterHeight() int {
	return ps.windowHeight - 1 - 2
}

func (ps *PlayerState) editAlbumsPanelHeight() int {
	return ps.windowHeight - 4 - (ps.windowHeight-4)/2
}

func (ps *PlayerState) editTitlesPanelHeight() int {
	return (ps.windowHeight - 4) / 2
}

type editorFinishedMsg struct {
	err      error
	tempFile string
}

func (ps *PlayerState) editOpenCover() {
	coverPath := filepath.Join(ps.config.MusicDir, ps.editAlbum[3], ps.editAlbum[4])
	c := exec.Command("xdg-open", coverPath)
	if err := c.Start(); err != nil {
		return
	}
	go c.Wait()
}

func (ps *PlayerState) editOpenEditor() tea.Cmd {
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "vi"
	}

	tmpFile, err := os.CreateTemp("", "mpcube-edit-*.txt")
	if err != nil {
		return nil
	}

	var content strings.Builder
	album := ps.musicData.Albums[ps.albumSelected]
	fmt.Fprintf(&content, "# mpcube edit — %s by %s\n", album.Album, album.Artist)
	fmt.Fprintf(&content, "# Lines starting with # are ignored. Save and quit to apply.\n\n")
	fmt.Fprintf(&content, "Album: %s\n", ps.editAlbum[0])
	fmt.Fprintf(&content, "Artist: %s\n", ps.editAlbum[1])
	fmt.Fprintf(&content, "Date: %s\n", ps.editAlbum[2])
	fmt.Fprintf(&content, "Directory: %s\n", ps.editAlbum[3])
	fmt.Fprintf(&content, "Cover: %s\n\n", ps.editAlbum[4])
	fmt.Fprintf(&content, "# Track | Title | Filename\n")
	for _, t := range ps.editTracks {
		fmt.Fprintf(&content, "%s | %s | %s\n", t.Track, t.Title, t.File)
	}

	if _, err := tmpFile.WriteString(content.String()); err != nil {
		tmpFile.Close()
		os.Remove(tmpFile.Name())
		return nil
	}
	tmpFile.Close()

	c := exec.Command(editor, tmpFile.Name())
	return tea.ExecProcess(c, func(err error) tea.Msg {
		return editorFinishedMsg{err: err, tempFile: tmpFile.Name()}
	})
}

func (ps *PlayerState) handleEditorFinished(msg editorFinishedMsg) {
	defer os.Remove(msg.tempFile)
	if msg.err != nil {
		return
	}

	data, err := os.ReadFile(msg.tempFile)
	if err != nil {
		return
	}

	inTracks := false
	trackIdx := 0
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		if strings.Contains(line, "|") {
			inTracks = true
			parts := strings.SplitN(line, "|", 3)
			if len(parts) == 3 && trackIdx < len(ps.editTracks) {
				ps.editTracks[trackIdx].Track = strings.TrimSpace(parts[0])
				ps.editTracks[trackIdx].Title = strings.TrimSpace(parts[1])
				ps.editTracks[trackIdx].File = strings.TrimSpace(parts[2])
				trackIdx++
			}
			continue
		}

		if inTracks {
			continue
		}

		if val, ok := strings.CutPrefix(line, "Album:"); ok {
			ps.editAlbum[0] = strings.TrimSpace(val)
		} else if val, ok := strings.CutPrefix(line, "Artist:"); ok {
			ps.editAlbum[1] = strings.TrimSpace(val)
		} else if val, ok := strings.CutPrefix(line, "Date:"); ok {
			ps.editAlbum[2] = strings.TrimSpace(val)
		} else if val, ok := strings.CutPrefix(line, "Directory:"); ok {
			ps.editAlbum[3] = strings.TrimSpace(val)
		} else if val, ok := strings.CutPrefix(line, "Cover:"); ok {
			ps.editAlbum[4] = strings.TrimSpace(val)
		}
	}
}
