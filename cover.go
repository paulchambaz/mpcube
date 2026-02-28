package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type coverResult struct {
	mbid         string
	releaseGroup string
	title        string
	artist       string
	date         string
	country      string
	status       string
	format       string
}

type coverSearchMsg struct {
	results []coverResult
	err     error
}

type coverDownloadMsg struct {
	path string
	ext  string
	err  error
}

func searchMusicBrainz(query string) ([]coverResult, error) {
	u := "https://musicbrainz.org/ws/2/release?query=" + strings.ReplaceAll(query, " ", "+") + "&fmt=json&limit=100"

	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "mpcube/1.0 (https://github.com/podcube/mpcube)")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("musicbrainz: %s", resp.Status)
	}

	var data struct {
		Releases []struct {
			ID           string `json:"id"`
			Title        string `json:"title"`
			Date         string `json:"date"`
			Country      string `json:"country"`
			Status       string `json:"status"`
			ArtistCredit []struct {
				Name string `json:"name"`
			} `json:"artist-credit"`
			ReleaseGroup struct {
				ID          string `json:"id"`
				PrimaryType string `json:"primary-type"`
			} `json:"release-group"`
			Media []struct {
				Format string `json:"format"`
			} `json:"media"`
		} `json:"releases"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("musicbrainz: %w", err)
	}

	var results []coverResult
	for _, r := range data.Releases {
		artist := ""
		if len(r.ArtistCredit) > 0 {
			artist = r.ArtistCredit[0].Name
		}
		format := ""
		if len(r.Media) > 0 {
			format = r.Media[0].Format
		}
		results = append(results, coverResult{
			mbid:         r.ID,
			releaseGroup: r.ReleaseGroup.ID,
			title:        r.Title,
			artist:       artist,
			date:         r.Date,
			country:      r.Country,
			status:       r.Status,
			format:       format,
		})
	}
	return results, nil
}

func downloadCoverArt(releaseGroupID, destPath string) (string, error) {
	u := "https://coverartarchive.org/release-group/" + releaseGroupID + "/front"

	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", "mpcube/1.0 (https://github.com/podcube/mpcube)")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("coverartarchive: %s", resp.Status)
	}

	ct := resp.Header.Get("Content-Type")
	ext := ".jpg"
	if strings.Contains(ct, "png") {
		ext = ".png"
	}

	finalPath := destPath + ext

	f, err := os.Create(finalPath)
	if err != nil {
		return "", err
	}
	defer f.Close()

	if _, err := io.Copy(f, resp.Body); err != nil {
		os.Remove(finalPath)
		return "", err
	}

	return ext, nil
}

func coverSearchCmd(query string) tea.Cmd {
	return func() tea.Msg {
		results, err := searchMusicBrainz(query)
		return coverSearchMsg{results: results, err: err}
	}
}

func coverDownloadCmd(releaseGroupID, destPath string) tea.Cmd {
	return func() tea.Msg {
		ext, err := downloadCoverArt(releaseGroupID, destPath)
		return coverDownloadMsg{path: destPath + ext, ext: ext, err: err}
	}
}

func (ps *PlayerState) handleCoverSearch(msg coverSearchMsg) (tea.Model, tea.Cmd) {
	ps.editCoverLoading = false
	ps.editCoverDownloading = false
	if msg.err != nil {
		ps.editCoverError = msg.err.Error()
		ps.mode = ModeEdit
		return ps, nil
	}
	ps.editCoverResults = msg.results
	ps.editCoverResultIdx = 0
	ps.editCoverResultOffset = 0
	if len(msg.results) == 0 {
		ps.editCoverError = "no covers found"
		ps.mode = ModeEdit
		return ps, nil
	}
	ps.editCoverError = ""
	ps.mode = ModeEditCoverResults
	return ps, nil
}

func (ps *PlayerState) handleCoverDownload(msg coverDownloadMsg) (tea.Model, tea.Cmd) {
	ps.editCoverLoading = false
	ps.editCoverDownloading = false
	if msg.err != nil {
		ps.editCoverError = msg.err.Error()
		ps.mode = ModeEditCoverResults
		return ps, nil
	}

	ps.editCoverPreviewPath = msg.path
	if ps.editCoverResultIdx < len(ps.editCoverResults) {
		ps.editCoverPreviewMBID = ps.editCoverResults[ps.editCoverResultIdx].releaseGroup
	}
	ps.editAlbum[4] = "cover" + msg.ext
	ps.editCoverPending = true
	ps.mode = ModeEditCoverResults

	if ps.editCoverOpenAfterDownload {
		ps.editCoverOpenAfterDownload = false
		c := exec.Command("xdg-open", msg.path)
		if err := c.Start(); err == nil {
			go c.Wait()
		}
	}

	return ps, nil
}

func (ps *PlayerState) coverFixResultOffset() {
	panelHeight := ps.editCoverResultsHeight()
	if panelHeight <= 0 {
		return
	}
	padding := min(1, panelHeight/4)

	if ps.editCoverResultIdx < ps.editCoverResultOffset+padding {
		ps.editCoverResultOffset = ps.editCoverResultIdx - padding
	}
	if ps.editCoverResultIdx >= ps.editCoverResultOffset+panelHeight-padding {
		ps.editCoverResultOffset = ps.editCoverResultIdx - panelHeight + 1 + padding
	}
	ps.editCoverResultOffset = max(ps.editCoverResultOffset, 0)
	ps.editCoverResultOffset = min(ps.editCoverResultOffset, max(0, len(ps.editCoverResults)-panelHeight))
}

func (ps *PlayerState) editCoverResultsHeight() int {
	rightHeight := (ps.windowHeight - 6) / 3
	// Subtract 2 for search bar line + separator
	return rightHeight - 2
}

func (ps *PlayerState) editCoverDir() string {
	return filepath.Join(ps.config.MusicDir, ps.editAlbumOrig[3])
}
