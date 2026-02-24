package main

import (
	"math/rand"
	"time"

	"github.com/charmbracelet/bubbles/key"
)

type keyMap struct {
	quit          key.Binding
	left          key.Binding
	right         key.Binding
	up            key.Binding
	down          key.Binding
	top           key.Binding
	bottom        key.Binding
	enter         key.Binding
	space         key.Binding
	next          key.Binding
	prev          key.Binding
	volumeUp      key.Binding
	volumeDown    key.Binding
	clear         key.Binding
	seekForward   key.Binding
	seekBackward  key.Binding
	toggleShuffle key.Binding
	toggleRepeat  key.Binding
	update        key.Binding
	random        key.Binding
}

var keys = keyMap{
	quit:          key.NewBinding(key.WithKeys("q", "Q", "ctrl+c"), key.WithHelp("q", "quit")),
	left:          key.NewBinding(key.WithKeys("left", "h"), key.WithHelp("h", "switch to albums")),
	right:         key.NewBinding(key.WithKeys("right", "l"), key.WithHelp("l", "switch to titles")),
	up:            key.NewBinding(key.WithKeys("up", "k"), key.WithHelp("k", "up")),
	down:          key.NewBinding(key.WithKeys("down", "j"), key.WithHelp("j", "down")),
	top:           key.NewBinding(key.WithKeys("home", "g"), key.WithHelp("g", "top")),
	bottom:        key.NewBinding(key.WithKeys("end", "G"), key.WithHelp("G", "bottom")),
	enter:         key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "play")),
	space:         key.NewBinding(key.WithKeys(" "), key.WithHelp("space", "pause/play")),
	next:          key.NewBinding(key.WithKeys("n"), key.WithHelp("n", "next")),
	prev:          key.NewBinding(key.WithKeys("p"), key.WithHelp("p", "previous")),
	volumeUp:      key.NewBinding(key.WithKeys("=", "+"), key.WithHelp("+", "volume up")),
	volumeDown:    key.NewBinding(key.WithKeys("-"), key.WithHelp("-", "volume down")),
	clear:         key.NewBinding(key.WithKeys("x"), key.WithHelp("-", "clear")),
	seekForward:   key.NewBinding(key.WithKeys("."), key.WithHelp(".", "seek forward")),
	seekBackward:  key.NewBinding(key.WithKeys(","), key.WithHelp(",", "seek backward")),
	toggleShuffle: key.NewBinding(key.WithKeys("s"), key.WithHelp("s", "shuffle")),
	toggleRepeat:  key.NewBinding(key.WithKeys("r"), key.WithHelp("r", "repeat")),
	update:        key.NewBinding(key.WithKeys("U"), key.WithHelp("U", "update")),
	random:        key.NewBinding(key.WithKeys("R"), key.WithHelp("R", "random")),
}

func (ps *PlayerState) moveUp() {
	albums := ps.musicData.Albums

	if len(albums) == 0 {
		return
	}

	padding := min(ps.config.ScrollPadding, (ps.windowHeight-4)/4)

	if ps.onAlbum {
		if ps.albumSelected > 0 {
			ps.albumSelected--
		}

		if ps.albumSelected < ps.albumOffset+padding && ps.albumOffset > 0 {
			ps.albumOffset--
		}
	} else {
		if ps.trackSelected > 0 {
			ps.trackSelected--
		}

		if ps.trackSelected < ps.trackOffset+padding && ps.trackOffset > 0 {
			ps.trackOffset--
		}
	}
}

func (ps *PlayerState) moveDown() {
	albums := ps.musicData.Albums

	if len(albums) == 0 {
		return
	}

	padding := min(ps.config.ScrollPadding, (ps.windowHeight-4)/4)

	if ps.onAlbum {
		if ps.albumSelected < len(albums)-1 {
			ps.albumSelected++

			if ps.albumSelected > ps.albumOffset+ps.windowHeight-5-padding && ps.albumOffset < len(albums)-ps.windowHeight+4 {
				ps.albumOffset++
			}
		}
	} else {
		tracks := albums[ps.albumSelected].Songs

		if ps.trackSelected < len(tracks)-1 {
			ps.trackSelected++

			if ps.trackSelected > ps.trackOffset+ps.windowHeight-5-padding && ps.trackOffset < len(tracks)-ps.windowHeight+4 {
				ps.trackOffset++
			}
		}
	}
}

func (ps *PlayerState) moveTop() {
	albums := ps.musicData.Albums

	if len(albums) == 0 {
		return
	}

	if ps.onAlbum {
		ps.albumSelected = 0
		ps.albumOffset = 0
	} else {
		ps.trackSelected = 0
		ps.trackOffset = 0
	}
}

func (ps *PlayerState) moveBottom() {
	albums := ps.musicData.Albums

	if len(albums) == 0 {
		return
	}

	if ps.onAlbum {
		ps.albumSelected = len(ps.musicData.Albums) - 1
		ps.albumOffset = max(0, len(ps.musicData.Albums)-ps.windowHeight+4)
	} else {
		tracks := albums[ps.albumSelected].Songs
		ps.trackSelected = len(tracks) - 1
		ps.trackOffset = max(0, len(tracks)-ps.windowHeight+4)
	}
}

func (ps *PlayerState) playSelected() error {
	if ps.onAlbum {
		return ps.playAlbum(ps.albumSelected)
	} else {
		return ps.playTrack(ps.albumSelected, ps.trackSelected)
	}
}

func (ps *PlayerState) playAlbum(albumIdx int) error {
	if albumIdx >= len(ps.musicData.Albums) {
		return nil
	}

	if err := ps.mpdClient.Clear(); err != nil {
		return err
	}
	album := ps.musicData.Albums[albumIdx]

	for _, song := range album.Songs {
		if err := ps.mpdClient.Add(song.URI); err != nil {
			return err
		}
	}

	if err := ps.mpdClient.Play(-1); err != nil {
		return err
	}

	return ps.updateMPDState()
}

func (ps *PlayerState) playTrack(albumIdx, trackIdx int) error {
	if albumIdx >= len(ps.musicData.Albums) {
		return nil
	}

	album := ps.musicData.Albums[albumIdx]
	if trackIdx >= len(album.Songs) {
		return nil
	}

	if err := ps.mpdClient.Clear(); err != nil {
		return err
	}

	for _, song := range album.Songs {
		if err := ps.mpdClient.Add(song.URI); err != nil {
			return err
		}
	}

	if err := ps.mpdClient.Play(trackIdx); err != nil {
		return err
	}

	return ps.updateMPDState()
}

func (ps *PlayerState) togglePlayPause() error {
	if ps.playing {
		if err := ps.mpdClient.Pause(true); err != nil {
			return err
		}
	} else {
		if err := ps.mpdClient.Play(-1); err != nil {
			return err
		}
	}

	return ps.updateMPDState()
}

func (ps *PlayerState) nextTrack() error {
	if ps.albumPlaying == nil || ps.trackPlaying == nil {
		return nil
	}

	if *ps.trackPlaying == len(ps.musicData.Albums[*ps.albumPlaying].Songs)-1 {
		ps.albumPlaying = nil
		ps.trackPlaying = nil
	}

	if err := ps.mpdClient.Next(); err != nil {
		return err
	}

	return ps.updateMPDState()
}

func (ps *PlayerState) prevTrack() error {
	if ps.albumPlaying == nil || ps.trackPlaying == nil {
		return nil
	}

	if err := ps.mpdClient.Previous(); err != nil {
		return err
	}

	return ps.updateMPDState()
}

func (ps *PlayerState) volumeUp() error {
	newVol := min(100, ps.volume+ps.config.VolumeStep)
	if err := ps.mpdClient.SetVolume(newVol); err != nil {
		return err
	}
	return ps.updateMPDState()
}

func (ps *PlayerState) volumeDown() error {
	newVol := max(0, ps.volume-ps.config.VolumeStep)
	if err := ps.mpdClient.SetVolume(newVol); err != nil {
		return err
	}
	return ps.updateMPDState()
}

func (ps *PlayerState) toggleShuffle() error {
	if err := ps.mpdClient.Random(!ps.shuffle); err != nil {
		return err
	}
	return ps.updateMPDState()
}

func (ps *PlayerState) toggleRepeat() error {
	if err := ps.mpdClient.Repeat(!ps.repeat); err != nil {
		return err
	}
	return ps.updateMPDState()
}

func (ps *PlayerState) clear() error {
	if err := ps.mpdClient.Clear(); err != nil {
		return err
	}
	ps.albumPlaying = nil
	ps.trackPlaying = nil
	return ps.updateMPDState()
}

// MPD adds ~550ms of forward drift per seek during playback due to audio
// output buffer refill. Compensate so .,/,. returns to the same position.
const seekDrift = 550 * time.Millisecond

func (ps *PlayerState) seekForward() error {
	if ps.albumPlaying == nil || ps.trackPlaying == nil {
		return nil
	}

	seekDur := time.Duration(ps.config.SeekDuration) * time.Millisecond
	if *ps.position > ps.musicData.Albums[*ps.albumPlaying].Songs[*ps.trackPlaying].Duration-seekDur {
		return ps.nextTrack()
	}

	actualSeek := seekDur
	if ps.playing {
		actualSeek -= seekDrift
	}

	if err := ps.mpdClient.SeekCur(actualSeek, true); err != nil {
		return err
	}
	return ps.updateMPDState()
}

func (ps *PlayerState) seekBackward() error {
	if ps.albumPlaying == nil || ps.trackPlaying == nil {
		return nil
	}

	seekDur := time.Duration(ps.config.SeekDuration) * time.Millisecond
	if *ps.position < seekDur {
		if err := ps.mpdClient.SeekCur(0, false); err != nil {
			return err
		}
	} else {
		actualSeek := seekDur
		if ps.playing {
			actualSeek += seekDrift
		}
		if err := ps.mpdClient.SeekCur(-actualSeek, true); err != nil {
			return err
		}
	}
	return ps.updateMPDState()
}

func (ps *PlayerState) update() error {
	if _, err := ps.mpdClient.Update(""); err != nil {
		return err
	}

	musicData, err := LoadMusicData(ps.mpdClient)
	if err != nil {
		return err
	}

	ps.musicData = musicData

	return ps.updateMPDState()
}

func (ps *PlayerState) random() error {
	if len(ps.musicData.Albums) == 0 {
		return nil
	}
	albumIdx := rand.Intn(len(ps.musicData.Albums))

	return ps.playAlbum(albumIdx)
}
