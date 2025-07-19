package main

import (
	"fmt"
	"math"
	"time"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/fhs/gompd/v2/mpd"
)

type PlayerState struct {
	musicData *MusicData
	mpdClient *mpd.Client

	windowWidth  int
	windowHeight int
	onAlbum      bool

	albumSelected int
	albumOffset   int

	trackSelected int
	trackOffset   int

	playing      bool
	albumPlaying *int
	trackPlaying *int
	position     *time.Duration
	volume       int
	shuffle      bool
	repeat       bool

	retryCount int
}

func NewPlayerState(musicData *MusicData, mpdClient *mpd.Client) (*PlayerState, error) {
	ps := &PlayerState{
		musicData:     musicData,
		mpdClient:     mpdClient,
		onAlbum:       true,
		albumSelected: 0,
		trackSelected: 0,
	}

	if err := ps.updateMPDState(); err != nil {
		return nil, err
	}

	return ps, nil
}

func (ps *PlayerState) updateMPDState() error {
	status, err := ps.mpdClient.Status()
	if err != nil {
		return err
	}

	ps.playing = status["state"] == "play"
	ps.shuffle = status["random"] == "1"
	ps.repeat = status["repeat"] == "1"

	if vol := status["volume"]; vol != "" {
		_, err := fmt.Sscanf(vol, "%d", &ps.volume)
		if err != nil {
			return err
		}
	}

	if elapsed := status["elapsed"]; elapsed != "" {
		if secs, err := time.ParseDuration(elapsed + "s"); err == nil {
			ps.position = &secs
		}
	}

	currentSong, err := ps.mpdClient.CurrentSong()
	if err != nil {
		ps.albumPlaying = nil
		ps.trackPlaying = nil
		return err
	}

	songURI := currentSong["file"]
	albumTitle := currentSong["Album"]

	for albumIdx, album := range ps.musicData.Albums {
		if album.Album == albumTitle {
			ps.albumPlaying = &albumIdx
			for trackIdx, song := range album.Songs {
				if song.URI == songURI {
					ps.trackPlaying = &trackIdx
					return nil
				}
			}
		}
	}

	return nil
}

func (ps *PlayerState) Init() tea.Cmd {
	return ps.tickCmd()
}

func (ps *PlayerState) tickCmd() tea.Cmd {
	return tea.Tick(100*time.Millisecond, func(time.Time) tea.Msg {
		return tickMsg{}
	})
}

type tickMsg struct{}

func (ps *PlayerState) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		ps.windowWidth = msg.Width
		ps.windowHeight = msg.Height
		ps.resize()
		return ps, nil

	case tickMsg:
		if err := ps.updateMPDState(); err != nil {
			ps.retryCount++
			delay := time.Duration(math.Min(float64(ps.retryCount*ps.retryCount), 30)) * time.Second
			return ps, tea.Tick(delay, func(time.Time) tea.Msg { return tickMsg{} })
		}
		ps.retryCount = 0
		return ps, ps.tickCmd()

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.quit):
			_ = ps.clear()
			return ps, tea.Quit
		case key.Matches(msg, keys.left):
			ps.onAlbum = false
		case key.Matches(msg, keys.right):
			ps.onAlbum = true
		case key.Matches(msg, keys.up):
			ps.moveUp()
		case key.Matches(msg, keys.down):
			ps.moveDown()
		case key.Matches(msg, keys.top):
			ps.moveTop()
		case key.Matches(msg, keys.bottom):
			ps.moveBottom()
		case key.Matches(msg, keys.enter):
			_ = ps.playSelected()
		case key.Matches(msg, keys.space):
			_ = ps.togglePlayPause()
		case key.Matches(msg, keys.next):
			_ = ps.nextTrack()
		case key.Matches(msg, keys.prev):
			_ = ps.prevTrack()
		case key.Matches(msg, keys.volumeUp):
			_ = ps.volumeUp()
		case key.Matches(msg, keys.volumeDown):
			_ = ps.volumeDown()
		case key.Matches(msg, keys.clear):
			_ = ps.clear()
		case key.Matches(msg, keys.seekForward):
			_ = ps.seekForward()
		case key.Matches(msg, keys.seekBackward):
			_ = ps.seekBackward()
		case key.Matches(msg, keys.toggleShuffle):
			_ = ps.toggleShuffle()
		case key.Matches(msg, keys.toggleRepeat):
			_ = ps.toggleRepeat()
		case key.Matches(msg, keys.update):
			_ = ps.update()
		case key.Matches(msg, keys.random):
			_ = ps.random()
		}
	}
	return ps, nil
}
