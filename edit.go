package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type EditFocus int

const (
	EditFocusCenter EditFocus = iota
	EditFocusAlbums
	EditFocusTitles
	EditFocusMetadata
	EditFocusCover
	EditFocusDownload
)

type editTrackState struct {
	Track string
	Title string
	File  string
}

const editAlbumFieldCount = 5

func (ps *PlayerState) enterEditMode() {
	if len(ps.musicData.Albums) == 0 {
		return
	}

	album := ps.musicData.Albums[ps.albumSelected]
	if len(album.Songs) == 0 {
		return
	}

	dir := filepath.Dir(album.Songs[0].URI)
	if dir == "." {
		dir = ""
	}

	ps.editAlbum = [5]string{
		album.Album,
		album.Artist,
		strconv.Itoa(album.Date),
		dir,
		"cover.jpg",
	}
	ps.editAlbumOrig = ps.editAlbum

	ps.editTracks = make([]editTrackState, len(album.Songs))
	ps.editTracksOrig = make([]editTrackState, len(album.Songs))
	for i, song := range album.Songs {
		t := editTrackState{
			Track: strconv.Itoa(song.Track),
			Title: song.Title,
			File:  filepath.Base(song.URI),
		}
		ps.editTracks[i] = t
		ps.editTracksOrig[i] = t
	}

	ps.editFocus = EditFocusCenter
	ps.editLastLeft = EditFocusAlbums
	ps.editLastRight = EditFocusMetadata
	ps.editFieldIdx = 0
	ps.editFieldOffset = 0
	ps.editTitleIdx = 0
	ps.editTitleOffset = 0
	ps.editInputBuf = ""
	ps.editInputPos = 0
	ps.mode = ModeEdit
	ps.editAlbumFixOffset()
	ps.editFixTitleOffset()
}

func (ps *PlayerState) exitEditMode() tea.Cmd {
	ps.editAlbum = [5]string{}
	ps.editAlbumOrig = [5]string{}
	ps.editTracks = nil
	ps.editTracksOrig = nil
	ps.editFieldIdx = 0
	ps.editFieldOffset = 0
	ps.editTitleIdx = 0
	ps.editTitleOffset = 0
	ps.editInputBuf = ""
	ps.editInputPos = 0
	ps.mode = ModeNormal
	client := ps.mpdClient
	return func() tea.Msg {
		client.Update("")
		musicData, _ := LoadMusicData(client)
		return libraryReloadMsg{musicData: musicData}
	}
}

func (ps *PlayerState) editFieldCount() int {
	return editAlbumFieldCount + len(ps.editTracks)*3
}

func (ps *PlayerState) editIsAlbumField(idx int) bool {
	return idx < editAlbumFieldCount
}

func (ps *PlayerState) editFieldToLine(idx int) int {
	if idx < editAlbumFieldCount {
		return idx
	}
	ti := (idx - editAlbumFieldCount) / 3
	fi := (idx - editAlbumFieldCount) % 3
	return editAlbumFieldCount + 1 + ti*4 + fi
}

func (ps *PlayerState) editTotalLines() int {
	if len(ps.editTracks) == 0 {
		return editAlbumFieldCount
	}
	return editAlbumFieldCount + 1 + len(ps.editTracks)*4 - 1
}

func (ps *PlayerState) editTrackIdx(idx int) int {
	return (idx - editAlbumFieldCount) / 3
}

func (ps *PlayerState) editTrackFieldIdx(idx int) int {
	return (idx - editAlbumFieldCount) % 3
}

func (ps *PlayerState) editCurrentLabel() string {
	if ps.editFieldIdx < editAlbumFieldCount {
		labels := [5]string{"Album", "Artist", "Date", "Dir", "Cover"}
		return labels[ps.editFieldIdx]
	}
	trackLabels := [3]string{"Track", "Title", "File"}
	return trackLabels[ps.editTrackFieldIdx(ps.editFieldIdx)]
}

func (ps *PlayerState) editCurrentValue() string {
	if ps.editFieldIdx < editAlbumFieldCount {
		return ps.editAlbum[ps.editFieldIdx]
	}
	ti := ps.editTrackIdx(ps.editFieldIdx)
	switch ps.editTrackFieldIdx(ps.editFieldIdx) {
	case 0:
		return ps.editTracks[ti].Track
	case 1:
		return ps.editTracks[ti].Title
	case 2:
		return ps.editTracks[ti].File
	}
	return ""
}

func (ps *PlayerState) editSetValue(val string) {
	if ps.editFieldIdx < editAlbumFieldCount {
		ps.editAlbum[ps.editFieldIdx] = val
		return
	}
	ti := ps.editTrackIdx(ps.editFieldIdx)
	switch ps.editTrackFieldIdx(ps.editFieldIdx) {
	case 0:
		ps.editTracks[ti].Track = val
	case 1:
		ps.editTracks[ti].Title = val
	case 2:
		ps.editTracks[ti].File = val
	}
}

func (ps *PlayerState) editIsModified(idx int) bool {
	if idx < editAlbumFieldCount {
		return ps.editAlbum[idx] != ps.editAlbumOrig[idx]
	}
	ti := (idx - editAlbumFieldCount) / 3
	fi := (idx - editAlbumFieldCount) % 3
	if ti >= len(ps.editTracks) {
		return false
	}
	switch fi {
	case 0:
		return ps.editTracks[ti].Track != ps.editTracksOrig[ti].Track
	case 1:
		return ps.editTracks[ti].Title != ps.editTracksOrig[ti].Title
	case 2:
		return ps.editTracks[ti].File != ps.editTracksOrig[ti].File
	}
	return false
}

func (ps *PlayerState) editLoadAlbum() {
	album := ps.musicData.Albums[ps.albumSelected]
	if len(album.Songs) == 0 {
		return
	}

	dir := filepath.Dir(album.Songs[0].URI)
	if dir == "." {
		dir = ""
	}

	ps.editAlbum = [5]string{
		album.Album,
		album.Artist,
		strconv.Itoa(album.Date),
		dir,
		"cover.jpg",
	}
	ps.editAlbumOrig = ps.editAlbum

	ps.editTracks = make([]editTrackState, len(album.Songs))
	ps.editTracksOrig = make([]editTrackState, len(album.Songs))
	for i, song := range album.Songs {
		t := editTrackState{
			Track: strconv.Itoa(song.Track),
			Title: song.Title,
			File:  filepath.Base(song.URI),
		}
		ps.editTracks[i] = t
		ps.editTracksOrig[i] = t
	}

	ps.editFieldIdx = 0
	ps.editFieldOffset = 0
	ps.editTitleIdx = 0
	ps.editTitleOffset = 0
}

func (ps *PlayerState) editRevertField() {
	idx := ps.editFieldIdx
	if idx < editAlbumFieldCount {
		ps.editAlbum[idx] = ps.editAlbumOrig[idx]
		return
	}
	ti := ps.editTrackIdx(idx)
	switch ps.editTrackFieldIdx(idx) {
	case 0:
		ps.editTracks[ti].Track = ps.editTracksOrig[ti].Track
	case 1:
		ps.editTracks[ti].Title = ps.editTracksOrig[ti].Title
	case 2:
		ps.editTracks[ti].File = ps.editTracksOrig[ti].File
	}
}

func (ps *PlayerState) editRevertAll() {
	ps.editAlbum = ps.editAlbumOrig
	for i := range ps.editTracks {
		ps.editTracks[i] = ps.editTracksOrig[i]
	}
}

func (ps *PlayerState) editTileNav(msg string) bool {
	switch msg {
	case "H":
		switch ps.editFocus {
		case EditFocusCenter:
			ps.editFocus = ps.editLastLeft
		case EditFocusMetadata, EditFocusCover, EditFocusDownload:
			ps.editLastRight = ps.editFocus
			ps.editFocus = EditFocusCenter
		}
		return true
	case "L":
		switch ps.editFocus {
		case EditFocusAlbums, EditFocusTitles:
			ps.editLastLeft = ps.editFocus
			ps.editFocus = EditFocusCenter
		case EditFocusCenter:
			ps.editFocus = ps.editLastRight
		}
		return true
	case "J":
		switch ps.editFocus {
		case EditFocusAlbums:
			ps.editFocus = EditFocusTitles
		case EditFocusMetadata:
			ps.editFocus = EditFocusCover
		case EditFocusCover:
			ps.editFocus = EditFocusDownload
		}
		return true
	case "K":
		switch ps.editFocus {
		case EditFocusTitles:
			ps.editFocus = EditFocusAlbums
		case EditFocusDownload:
			ps.editFocus = EditFocusCover
		case EditFocusCover:
			ps.editFocus = EditFocusMetadata
		}
		return true
	}
	return false
}

func sanitizeFilename(s string) string {
	return strings.ReplaceAll(s, "/", "-")
}

func (ps *PlayerState) editSyncFilenames() {
	ps.editAlbum[3] = sanitizeFilename(ps.editAlbum[1] + " - " + ps.editAlbum[0])

	for i := range ps.editTracks {
		ext := filepath.Ext(ps.editTracksOrig[i].File)
		track, err := strconv.Atoi(ps.editTracks[i].Track)
		if err != nil {
			track = i + 1
		}
		ps.editTracks[i].File = fmt.Sprintf("%02d - %s%s", track, sanitizeFilename(ps.editTracks[i].Title), ext)
	}
}

// Apply pipeline

type applyOp int

const (
	applyOpTagWrite applyOp = iota
	applyOpRename
)

type applyCmd struct {
	fieldIdx int
	op       applyOp
	srcPath  string
	dstPath  string
	tags     map[string]string
}

var albumTagNames = [3]string{"Album", "Artist", "Date"}
var trackTagNames = [2]string{"Track Number", "Title"}

func (ps *PlayerState) editBuildApplyField(fieldIdx int) []applyCmd {
	baseDir := ps.config.MusicDir
	currentDir := ps.editAlbumOrig[3]
	var cmds []applyCmd

	if fieldIdx < 3 {
		for _, track := range ps.editTracksOrig {
			cmds = append(cmds, applyCmd{
				fieldIdx: fieldIdx,
				op:       applyOpTagWrite,
				srcPath:  filepath.Join(baseDir, currentDir, track.File),
				tags:     map[string]string{albumTagNames[fieldIdx]: ps.editAlbum[fieldIdx]},
			})
		}
	} else if fieldIdx == 3 {
		cmds = append(cmds, applyCmd{
			fieldIdx: fieldIdx,
			op:       applyOpRename,
			srcPath:  filepath.Join(baseDir, currentDir),
			dstPath:  filepath.Join(baseDir, ps.editAlbum[3]),
		})
	} else if fieldIdx >= editAlbumFieldCount {
		ti := ps.editTrackIdx(fieldIdx)
		fi := ps.editTrackFieldIdx(fieldIdx)
		if fi < 2 {
			cmds = append(cmds, applyCmd{
				fieldIdx: fieldIdx,
				op:       applyOpTagWrite,
				srcPath:  filepath.Join(baseDir, currentDir, ps.editTracksOrig[ti].File),
				tags:     map[string]string{trackTagNames[fi]: ps.editTrackValue(ti, fi)},
			})
		} else {
			cmds = append(cmds, applyCmd{
				fieldIdx: fieldIdx,
				op:       applyOpRename,
				srcPath:  filepath.Join(baseDir, currentDir, ps.editTracksOrig[ti].File),
				dstPath:  filepath.Join(baseDir, currentDir, ps.editTracks[ti].File),
			})
		}
	}

	return cmds
}

func (ps *PlayerState) editBuildApplyAll() []applyCmd {
	baseDir := ps.config.MusicDir
	currentDir := ps.editAlbumOrig[3]
	var cmds []applyCmd

	for idx := 0; idx < 3; idx++ {
		if !ps.editIsModified(idx) {
			continue
		}
		for _, track := range ps.editTracksOrig {
			cmds = append(cmds, applyCmd{
				fieldIdx: idx,
				op:       applyOpTagWrite,
				srcPath:  filepath.Join(baseDir, currentDir, track.File),
				tags:     map[string]string{albumTagNames[idx]: ps.editAlbum[idx]},
			})
		}
	}

	if ps.editIsModified(3) {
		cmds = append(cmds, applyCmd{
			fieldIdx: 3,
			op:       applyOpRename,
			srcPath:  filepath.Join(baseDir, currentDir),
			dstPath:  filepath.Join(baseDir, ps.editAlbum[3]),
		})
		currentDir = ps.editAlbum[3]
	}

	for ti := 0; ti < len(ps.editTracks); ti++ {
		baseIdx := editAlbumFieldCount + ti*3
		for fi := 0; fi < 2; fi++ {
			if !ps.editIsModified(baseIdx + fi) {
				continue
			}
			cmds = append(cmds, applyCmd{
				fieldIdx: baseIdx + fi,
				op:       applyOpTagWrite,
				srcPath:  filepath.Join(baseDir, currentDir, ps.editTracksOrig[ti].File),
				tags:     map[string]string{trackTagNames[fi]: ps.editTrackValue(ti, fi)},
			})
		}
		if ps.editIsModified(baseIdx + 2) {
			cmds = append(cmds, applyCmd{
				fieldIdx: baseIdx + 2,
				op:       applyOpRename,
				srcPath:  filepath.Join(baseDir, currentDir, ps.editTracksOrig[ti].File),
				dstPath:  filepath.Join(baseDir, currentDir, ps.editTracks[ti].File),
			})
		}
	}

	return cmds
}

func (ps *PlayerState) editTrackValue(ti, fi int) string {
	switch fi {
	case 0:
		return ps.editTracks[ti].Track
	case 1:
		return ps.editTracks[ti].Title
	}
	return ""
}

func (ps *PlayerState) editStartApply(cmds []applyCmd) tea.Cmd {
	if len(cmds) == 0 {
		return nil
	}
	ps.applyQueue = cmds
	ps.applyProgress = 0
	ps.applyError = ""
	ps.applyReturnFocus = ps.editFocus
	ps.editFocus = EditFocusCenter
	ps.editFieldIdx = cmds[0].fieldIdx
	ps.editFixCenterOffset()
	ps.mode = ModeEditApply
	return ps.applyNextStep()
}

func (ps *PlayerState) applyNextStep() tea.Cmd {
	cmd := ps.applyQueue[ps.applyProgress]
	return func() tea.Msg {
		var err error
		switch cmd.op {
		case applyOpTagWrite:
			err = kid3WriteTags(cmd.srcPath, cmd.tags)
		case applyOpRename:
			err = os.Rename(cmd.srcPath, cmd.dstPath)
		}
		return applyStepMsg{err: err}
	}
}

func (ps *PlayerState) handleApplyStep(msg applyStepMsg) (tea.Model, tea.Cmd) {
	if msg.err != nil {
		ps.applyError = msg.err.Error()
		ps.applyQueue = nil
		ps.mode = ModeEdit
		return ps, nil
	}

	completedFieldIdx := ps.applyQueue[ps.applyProgress].fieldIdx
	ps.applyProgress++

	nextIsNewField := ps.applyProgress >= len(ps.applyQueue) ||
		ps.applyQueue[ps.applyProgress].fieldIdx != completedFieldIdx

	if nextIsNewField {
		ps.editApplyUpdateOrig(completedFieldIdx)
	}

	if ps.applyProgress < len(ps.applyQueue) {
		ps.editFieldIdx = ps.applyQueue[ps.applyProgress].fieldIdx
		ps.editFixCenterOffset()
		return ps, ps.applyNextStep()
	}

	return ps, ps.applyFinishCmd()
}

func (ps *PlayerState) editApplyUpdateOrig(fieldIdx int) {
	if fieldIdx < editAlbumFieldCount {
		ps.editAlbumOrig[fieldIdx] = ps.editAlbum[fieldIdx]
		return
	}
	ti := ps.editTrackIdx(fieldIdx)
	fi := ps.editTrackFieldIdx(fieldIdx)
	switch fi {
	case 0:
		ps.editTracksOrig[ti].Track = ps.editTracks[ti].Track
	case 1:
		ps.editTracksOrig[ti].Title = ps.editTracks[ti].Title
	case 2:
		ps.editTracksOrig[ti].File = ps.editTracks[ti].File
	}
}

func (ps *PlayerState) applyFinishCmd() tea.Cmd {
	client := ps.mpdClient
	return func() tea.Msg {
		client.Update("")
		return applyFinishMsg{}
	}
}

func (ps *PlayerState) handleApplyFinish(_ applyFinishMsg) (tea.Model, tea.Cmd) {
	ps.applyQueue = nil
	ps.editFocus = ps.applyReturnFocus
	ps.mode = ModeEdit
	return ps, nil
}
