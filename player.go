package main

import (
	"fmt"
	"math"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/fhs/gompd/v2/mpd"
)

type Mode int

const (
	ModeNormal Mode = iota
	ModeSearch
	ModeSearching
	ModeHelp
	ModeEdit
	ModeEditInput
	ModeEditSearch
	ModeEditSearching
	ModeEditApply
	ModeEditCoverInput
	ModeEditCoverResults
)

type PlayerState struct {
	config    Config
	musicData *MusicData
	mpdClient *mpd.Client

	windowWidth  int
	windowHeight int
	onAlbum     bool
	mode        Mode
	helpForMode Mode

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

	searchQuery       string
	searchMatches     []int
	searchMatchIdx    int
	searchSavedAlbum  int
	searchSavedOffset int

	retryCount int

	editFocus       EditFocus
	editLastLeft    EditFocus
	editLastRight   EditFocus
	editFieldIdx    int
	editFieldOffset int
	editTitleIdx    int
	editTitleOffset int
	editAlbum       [5]string
	editAlbumOrig   [5]string
	editTracks      []editTrackState
	editTracksOrig  []editTrackState
	editCorrupted      []bool
	editHasCoverFile   bool
	editCoverFile      string
	editHasEmbeddedArt   bool
	editStripEmbeddedArt bool
	editCoverSearch           string
	editCoverResults          []coverResult
	editCoverResultIdx        int
	editCoverResultOffset     int
	editCoverLoading          bool
	editCoverError            string
	editCoverPreviewPath      string
	editCoverPreviewMBID      string
	editCoverPending            bool
	editCoverDownloading        bool
	editCoverOpenAfterDownload  bool
	editInputBuf       string
	editInputPos       int

	applyQueue       []applyCmd
	applyProgress    int
	applyError       string
	applyReturnFocus EditFocus
}

func NewPlayerState(config Config, musicData *MusicData, mpdClient *mpd.Client) (*PlayerState, error) {
	ps := &PlayerState{
		config:        config,
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
	return tea.Tick(time.Duration(ps.config.TickInterval)*time.Millisecond, func(time.Time) tea.Msg {
		return tickMsg{}
	})
}

type tickMsg struct{}

type applyStepMsg struct {
	err error
}

type applyFinishMsg struct {
	musicData *MusicData
}

type libraryReloadMsg struct {
	musicData *MusicData
}

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
			delay := time.Duration(math.Min(float64(ps.retryCount*ps.retryCount), float64(ps.config.MaxRetryDelay))) * time.Second
			return ps, tea.Tick(delay, func(time.Time) tea.Msg { return tickMsg{} })
		}
		ps.retryCount = 0
		return ps, ps.tickCmd()

	case editorFinishedMsg:
		ps.handleEditorFinished(msg)
		return ps, nil

	case applyStepMsg:
		return ps.handleApplyStep(msg)

	case applyFinishMsg:
		return ps.handleApplyFinish(msg)

	case coverSearchMsg:
		return ps.handleCoverSearch(msg)

	case coverDownloadMsg:
		return ps.handleCoverDownload(msg)

	case libraryReloadMsg:
		if msg.musicData != nil {
			ps.musicData = msg.musicData
			ps.updateMPDState()
			ps.resize()
		}
		return ps, nil

	case tea.KeyMsg:
		return ps.handleKeyMsg(msg)
	}
	return ps, nil
}

