package main

import (
	"fmt"
	"path/filepath"
	"strconv"
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

func (ps *PlayerState) exitEditMode() {
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

func (ps *PlayerState) editSyncFilenames() {
	ps.editAlbum[3] = ps.editAlbum[1] + " - " + ps.editAlbum[0]

	for i := range ps.editTracks {
		ext := filepath.Ext(ps.editTracksOrig[i].File)
		track, err := strconv.Atoi(ps.editTracks[i].Track)
		if err != nil {
			track = i + 1
		}
		ps.editTracks[i].File = fmt.Sprintf("%02d - %s%s", track, ps.editTracks[i].Title, ext)
	}
}
