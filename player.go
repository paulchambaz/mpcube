package main

import (
	"fmt"
	"math"
	"strconv"
	"strings"
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
	ModeEditCoverInput
	ModeEditMetadataInput
	ModeEditApply
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
	editCoverError            string
	editCoverPreviewPath      string
	editCoverPreviewMBID      string
	editCoverPreviewResultIdx int
	editCoverPending          bool
	editCoverSearching        bool
	editCoverDownloading      bool
	editMetadataSearch        string
	editMetadataResults       []coverResult
	editMetadataResultIdx     int
	editMetadataResultOffset  int
	editMetadataError         string
	editMetadataPending       bool
	editMetadataSearching     bool
	editInputBuf       string
	editInputPos       int

	editingUUID              string
	editingOriginalTrackURIs []string
	editingOriginalDir       string
	editingCurrentDir        string

	applyError       string
	applyQueue       []applyCmd
	applyProgress    int
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

func (ps *PlayerState) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		ps.windowWidth = msg.Width
		ps.windowHeight = msg.Height
		ps.resize()
		return ps, nil

	case tickMsg:
		// Process apply queue if in apply mode
		if ps.mode == ModeEditApply {
			return ps.handleApplyTick()
		}

		// Normal MPD state update
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

	case metadataSearchResultMsg:
		ps.editMetadataSearching = false
		if msg.err != nil {
			ps.editMetadataError = msg.err.Error()
		} else {
			sorted := sortCoverResultsByScore(msg.results, ps.editMetadataSearch)
			ps.editMetadataResults = sorted
			ps.editMetadataResultIdx = 0
			ps.editMetadataResultOffset = 0
			if len(msg.results) == 0 {
				ps.editMetadataError = "no results found"
			} else {
				ps.editMetadataError = ""
			}
		}
		return ps, nil

	case metadataFetchResultMsg:
		ps.editMetadataSearching = false
		if msg.err != nil {
			ps.editMetadataError = msg.err.Error()
			return ps, nil
		}
		// Ignore if results were cleared (e.g., user exited editor during fetch)
		if len(ps.editMetadataResults) == 0 || ps.editMetadataResultIdx >= len(ps.editMetadataResults) {
			return ps, nil
		}
		// Stage the fetched metadata
		selected := ps.editMetadataResults[ps.editMetadataResultIdx]
		ps.editAlbum[0] = selected.title
		ps.editAlbum[1] = selected.artist
		if len(selected.date) >= 4 {
			ps.editAlbum[2] = selected.date[:4]
		} else {
			ps.editAlbum[2] = selected.date
		}
		// Match tracks by track number
		for i := range ps.editTracks {
			trackNumStr := strings.TrimSpace(ps.editTracks[i].Track)
			trackNum, err := strconv.Atoi(trackNumStr)
			if err != nil {
				continue
			}
			for _, mbTrack := range msg.tracks {
				if mbTrack.position == trackNum {
					ps.editTracks[i].Track = strconv.Itoa(mbTrack.position)
					ps.editTracks[i].Title = mbTrack.title
					break
				}
			}
		}
		ps.editMetadataPending = true
		ps.editMetadataError = ""
		return ps, nil

	case coverSearchResultMsg:
		ps.editCoverSearching = false
		if msg.err != nil {
			ps.editCoverError = msg.err.Error()
		} else {
			sorted := sortCoverResultsByScore(msg.results, ps.editCoverSearch)
			ps.editCoverResults = sorted
			ps.editCoverResultIdx = 0
			ps.editCoverResultOffset = 0
			if len(msg.results) == 0 {
				ps.editCoverError = "no covers found"
			} else {
				ps.editCoverError = ""
			}
		}
		return ps, nil

	case coverDownloadResultMsg:
		ps.editCoverDownloading = false
		if msg.err != nil {
			ps.editCoverError = msg.err.Error()
			return ps, nil
		}
		// Ignore if results were cleared (e.g., user exited editor during download)
		if len(ps.editCoverResults) == 0 || ps.editCoverResultIdx >= len(ps.editCoverResults) {
			return ps, nil
		}
		// Update preview path
		ps.editCoverPreviewPath = msg.path + msg.ext
		ps.editCoverPreviewMBID = ps.editCoverResults[ps.editCoverResultIdx].releaseGroup
		ps.editCoverPreviewResultIdx = ps.editCoverResultIdx
		ps.editAlbum[4] = "cover" + msg.ext
		// Only stage if requested (enter key, not 'o' key)
		if msg.stageForInstall {
			ps.editCoverPending = true
		}
		ps.editCoverError = ""
		return ps, nil

	case tea.KeyMsg:
		return ps.handleKeyMsg(msg)
	}
	return ps, nil
}

