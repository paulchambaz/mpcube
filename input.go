package main

import (
	"strings"

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

func (ps *PlayerState) handleKeyMsg(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if key.Matches(msg, globalKeys.forceQuit) {
		_ = ps.clear()
		return ps, tea.Quit
	}

	if ps.mode == ModeEditApply || ps.editCoverLoading {
		return ps, nil
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

	if ps.mode != ModeSearch && ps.mode != ModeEditInput && ps.mode != ModeEditSearch && ps.mode != ModeEditCoverInput {
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
	case ModeEdit:
		return ps.handleEdit(msg)
	case ModeEditInput:
		return ps.handleEditInput(msg)
	case ModeEditSearch:
		return ps.handleEditSearch(msg)
	case ModeEditSearching:
		return ps.handleEditSearching(msg)
	case ModeEditCoverInput:
		return ps.handleEditCoverInput(msg)
	case ModeEditCoverResults:
		return ps.handleEditCoverResults(msg)
	}
	return ps, nil
}
