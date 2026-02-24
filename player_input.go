package main

import (
	"math/rand"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

func bindingHelp(b key.Binding) helpEntry {
	keys := make([]string, len(b.Keys()))
	for i, k := range b.Keys() {
		if k == " " {
			keys[i] = "█"
		} else {
			keys[i] = k
		}
	}
	return helpEntry{strings.Join(keys, ", "), b.Help().Desc}
}

type globalKeyMap struct {
	forceQuit    key.Binding
	help         key.Binding
	seekForward  key.Binding
	seekBackward key.Binding
	volumeUp     key.Binding
	volumeDown   key.Binding
}

var globalKeys = globalKeyMap{
	forceQuit:    key.NewBinding(key.WithKeys("ctrl+c"), key.WithHelp("ctrl+c", "quit")),
	help:         key.NewBinding(key.WithKeys("?"), key.WithHelp("?", "help")),
	seekForward:  key.NewBinding(key.WithKeys("."), key.WithHelp(".", "seek forward")),
	seekBackward: key.NewBinding(key.WithKeys(","), key.WithHelp(",", "seek backward")),
	volumeUp:     key.NewBinding(key.WithKeys("=", "+"), key.WithHelp("+", "volume up")),
	volumeDown:   key.NewBinding(key.WithKeys("-", "_"), key.WithHelp("-", "volume down")),
}

type normalKeyMap struct {
	quit          key.Binding
	left          key.Binding
	right         key.Binding
	up            key.Binding
	down          key.Binding
	top           key.Binding
	bottom        key.Binding
	enter         key.Binding
	playPause     key.Binding
	next          key.Binding
	prev          key.Binding
	clear         key.Binding
	toggleShuffle key.Binding
	toggleRepeat  key.Binding
	update        key.Binding
	random        key.Binding
	search        key.Binding
}

var normalKeys = normalKeyMap{
	quit:          key.NewBinding(key.WithKeys("q", "Q"), key.WithHelp("q", "quit")),
	left:          key.NewBinding(key.WithKeys("left", "h"), key.WithHelp("h", "albums")),
	right:         key.NewBinding(key.WithKeys("right", "l"), key.WithHelp("l", "tracks")),
	up:            key.NewBinding(key.WithKeys("up", "k"), key.WithHelp("k", "up")),
	down:          key.NewBinding(key.WithKeys("down", "j"), key.WithHelp("j", "down")),
	top:           key.NewBinding(key.WithKeys("home", "g"), key.WithHelp("g", "top")),
	bottom:        key.NewBinding(key.WithKeys("end", "G"), key.WithHelp("G", "bottom")),
	enter:         key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "start")),
	playPause:     key.NewBinding(key.WithKeys(" "), key.WithHelp("space", "play/pause")),
	next:          key.NewBinding(key.WithKeys("n"), key.WithHelp("n", "next track")),
	prev:          key.NewBinding(key.WithKeys("p"), key.WithHelp("p", "previous track")),
	clear:         key.NewBinding(key.WithKeys("x"), key.WithHelp("x", "clear")),
	toggleShuffle: key.NewBinding(key.WithKeys("s"), key.WithHelp("s", "shuffle")),
	toggleRepeat:  key.NewBinding(key.WithKeys("r"), key.WithHelp("r", "repeat")),
	update:        key.NewBinding(key.WithKeys("U"), key.WithHelp("U", "update")),
	random:        key.NewBinding(key.WithKeys("R"), key.WithHelp("R", "random")),
	search:        key.NewBinding(key.WithKeys("/"), key.WithHelp("/", "search")),
}

func (k normalKeyMap) bindings() []helpEntry {
	return []helpEntry{
		{"q, Q, ctrl+c", "quit"},
		bindingHelp(k.left),
		bindingHelp(k.right),
		bindingHelp(k.up),
		bindingHelp(k.down),
		bindingHelp(k.top),
		bindingHelp(k.bottom),
		bindingHelp(k.enter),
		bindingHelp(k.playPause),
		bindingHelp(k.next),
		bindingHelp(k.prev),
		bindingHelp(k.clear),
		bindingHelp(k.toggleShuffle),
		bindingHelp(k.toggleRepeat),
		bindingHelp(k.update),
		bindingHelp(k.random),
		bindingHelp(k.search),
		bindingHelp(globalKeys.seekForward),
		bindingHelp(globalKeys.seekBackward),
		bindingHelp(globalKeys.volumeUp),
		bindingHelp(globalKeys.volumeDown),
		{"?", "help"},
	}
}

type searchingKeyMap struct {
	nextMatch key.Binding
	prevMatch key.Binding
	confirm   key.Binding
	cancel    key.Binding
	reSearch  key.Binding
	playPause key.Binding
	quit      key.Binding
}

var searchingKeys = searchingKeyMap{
	nextMatch: key.NewBinding(key.WithKeys("n"), key.WithHelp("n", "next match")),
	prevMatch: key.NewBinding(key.WithKeys("p"), key.WithHelp("p", "previous match")),
	confirm:   key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "play album")),
	cancel:    key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "cancel")),
	reSearch:  key.NewBinding(key.WithKeys("/", "i"), key.WithHelp("/", "edit query")),
	playPause: key.NewBinding(key.WithKeys(" "), key.WithHelp("space", "play/pause")),
	quit:      key.NewBinding(key.WithKeys("q", "Q"), key.WithHelp("q", "quit")),
}

func (k searchingKeyMap) bindings() []helpEntry {
	return []helpEntry{
		{"q, Q, ctrl+c", "quit"},
		bindingHelp(k.nextMatch),
		bindingHelp(k.prevMatch),
		bindingHelp(k.confirm),
		bindingHelp(k.cancel),
		bindingHelp(k.reSearch),
		{"?", "help"},
	}
}

func (ps *PlayerState) handleKeyMsg(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if key.Matches(msg, globalKeys.forceQuit) {
		_ = ps.clear()
		return ps, tea.Quit
	}

	if key.Matches(msg, globalKeys.help) {
		if ps.mode == ModeHelp {
			ps.mode = ps.helpForMode
			return ps, nil
		}
		ps.helpForMode = ps.mode
		ps.mode = ModeHelp
		return ps, nil
	}

	switch {
	case key.Matches(msg, globalKeys.seekForward):
		_ = ps.seekForward()
		return ps, nil
	case key.Matches(msg, globalKeys.seekBackward):
		_ = ps.seekBackward()
		return ps, nil
	case key.Matches(msg, globalKeys.volumeUp):
		_ = ps.volumeUp()
		return ps, nil
	case key.Matches(msg, globalKeys.volumeDown):
		_ = ps.volumeDown()
		return ps, nil
	}

	switch ps.mode {
	case ModeNormal:
		return ps.handleNormal(msg)
	case ModeSearch:
		return ps.handleSearch(msg)
	case ModeSearching:
		return ps.handleSearching(msg)
	case ModeHelp:
		return ps.handleHelp(msg)
	}
	return ps, nil
}

func (ps *PlayerState) handleNormal(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case key.Matches(msg, normalKeys.quit):
		_ = ps.clear()
		return ps, tea.Quit
	case key.Matches(msg, normalKeys.right):
		ps.onAlbum = false
		if len(ps.musicData.Albums) > 0 {
			lastTrack := len(ps.musicData.Albums[ps.albumSelected].Songs) - 1
			if ps.trackSelected > lastTrack {
				ps.trackSelected = lastTrack
			}
			if ps.trackOffset > lastTrack {
				ps.trackOffset = max(0, lastTrack-ps.windowHeight+5)
			}
		}
	case key.Matches(msg, normalKeys.left):
		ps.onAlbum = true
	case key.Matches(msg, normalKeys.up):
		ps.moveUp()
	case key.Matches(msg, normalKeys.down):
		ps.moveDown()
	case key.Matches(msg, normalKeys.top):
		ps.moveTop()
	case key.Matches(msg, normalKeys.bottom):
		ps.moveBottom()
	case key.Matches(msg, normalKeys.enter):
		_ = ps.playSelected()
	case key.Matches(msg, normalKeys.playPause):
		_ = ps.togglePlayPause()
	case key.Matches(msg, normalKeys.next):
		_ = ps.nextTrack()
	case key.Matches(msg, normalKeys.prev):
		_ = ps.prevTrack()
	case key.Matches(msg, normalKeys.clear):
		_ = ps.clear()
	case key.Matches(msg, normalKeys.toggleShuffle):
		_ = ps.toggleShuffle()
	case key.Matches(msg, normalKeys.toggleRepeat):
		_ = ps.toggleRepeat()
	case key.Matches(msg, normalKeys.update):
		_ = ps.update()
	case key.Matches(msg, normalKeys.random):
		_ = ps.random()
	case key.Matches(msg, normalKeys.search):
		ps.enterSearch()
	}
	return ps, nil
}

func (ps *PlayerState) handleSearch(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyEscape:
		ps.cancelSearch()
	case tea.KeyEnter:
		ps.confirmSearch()
	case tea.KeyBackspace:
		ps.searchBackspace()
	case tea.KeyRunes:
		ps.searchAddRune(msg.Runes[0])
	}
	return ps, nil
}

func (ps *PlayerState) handleSearching(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case key.Matches(msg, searchingKeys.nextMatch):
		ps.nextMatch()
	case key.Matches(msg, searchingKeys.prevMatch):
		ps.prevMatch()
	case key.Matches(msg, searchingKeys.confirm):
		ps.confirmSearching()
	case key.Matches(msg, searchingKeys.cancel):
		ps.cancelSearch()
	case key.Matches(msg, searchingKeys.reSearch):
		ps.mode = ModeSearch
	case key.Matches(msg, searchingKeys.playPause):
		_ = ps.togglePlayPause()
	case key.Matches(msg, searchingKeys.quit):
		_ = ps.clear()
		return ps, tea.Quit
	}
	return ps, nil
}

func (ps *PlayerState) handleHelp(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if msg.Type == tea.KeyEscape || msg.String() == "q" || msg.String() == "Q" {
		ps.mode = ps.helpForMode
	}
	return ps, nil
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
