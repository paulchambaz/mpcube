package main

import (
	"sort"
	"strings"
)

func fuzzyScore(query, target string) int {
	query = strings.ToLower(query)
	target = strings.ToLower(target)
	qi := 0
	score := 0
	prevIdx := -1

	for i := 0; i < len(target) && qi < len(query); i++ {
		if target[i] == query[qi] {
			score++
			if prevIdx >= 0 && i == prevIdx+1 {
				score += 3
			}
			if i == 0 || target[i-1] == ' ' {
				score += 5
			}
			prevIdx = i
			qi++
		}
	}

	if qi < len(query) {
		return -1
	}
	return score
}

func (ps *PlayerState) runSearch() {
	ps.searchMatches = nil
	ps.searchMatchIdx = 0
	if ps.searchQuery == "" {
		return
	}

	type result struct {
		idx   int
		score int
	}
	var results []result

	for i, album := range ps.musicData.Albums {
		s1 := fuzzyScore(ps.searchQuery, album.Artist+" "+album.Album)
		s2 := fuzzyScore(ps.searchQuery, album.Album+" "+album.Artist)
		score := max(s1, s2)
		if score >= 0 {
			results = append(results, result{i, score})
		}
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].score > results[j].score
	})

	for _, r := range results {
		ps.searchMatches = append(ps.searchMatches, r.idx)
	}
}

func (ps *PlayerState) jumpToMatch(idx int) {
	if len(ps.searchMatches) == 0 {
		return
	}
	ps.searchMatchIdx = idx
	ps.albumSelected = ps.searchMatches[idx]
	ps.trackSelected = 0
	ps.trackOffset = 0
	ps.resize()
}

func (ps *PlayerState) enterSearch() {
	ps.searchSavedAlbum = ps.albumSelected
	ps.searchSavedOffset = ps.albumOffset
	ps.searchQuery = ""
	ps.searchMatches = nil
	ps.searchMatchIdx = 0
	ps.mode = ModeSearch
}

func (ps *PlayerState) cancelSearch() {
	ps.albumSelected = ps.searchSavedAlbum
	ps.albumOffset = ps.searchSavedOffset
	ps.searchQuery = ""
	ps.searchMatches = nil
	ps.mode = ModeNormal
}

func (ps *PlayerState) confirmSearch() {
	if ps.searchQuery == "" {
		ps.cancelSearch()
		return
	}
	if len(ps.searchMatches) > 0 {
		ps.mode = ModeSearching
	}
}

func (ps *PlayerState) searchAddRune(r rune) {
	ps.searchQuery += string(r)
	ps.runSearch()
	if len(ps.searchMatches) > 0 {
		ps.jumpToMatch(0)
	}
}

func (ps *PlayerState) searchBackspace() {
	if len(ps.searchQuery) > 0 {
		ps.searchQuery = ps.searchQuery[:len(ps.searchQuery)-1]
	}
	if ps.searchQuery == "" {
		ps.albumSelected = ps.searchSavedAlbum
		ps.albumOffset = ps.searchSavedOffset
		return
	}
	ps.runSearch()
	if len(ps.searchMatches) > 0 {
		ps.jumpToMatch(0)
	}
}

func (ps *PlayerState) nextMatch() {
	if len(ps.searchMatches) == 0 {
		return
	}
	idx := (ps.searchMatchIdx + 1) % len(ps.searchMatches)
	ps.jumpToMatch(idx)
}

func (ps *PlayerState) prevMatch() {
	if len(ps.searchMatches) == 0 {
		return
	}
	idx := ps.searchMatchIdx - 1
	if idx < 0 {
		idx = len(ps.searchMatches) - 1
	}
	ps.jumpToMatch(idx)
}

func (ps *PlayerState) confirmSearching() {
	_ = ps.playAlbum(ps.albumSelected)
	ps.mode = ModeNormal
}
