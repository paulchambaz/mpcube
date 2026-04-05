package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
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

type metadataTrack struct {
	position int
	title    string
}

type metadataSearchResultMsg struct {
	results []coverResult
	err     error
}

type metadataFetchResultMsg struct {
	tracks []metadataTrack
	err    error
}

type coverSearchResultMsg struct {
	results []coverResult
	err     error
}

type coverDownloadResultMsg struct {
	path            string
	ext             string
	err             error
	stageForInstall bool
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

func (ps *PlayerState) doCoverSearch(query string) {
	results, err := searchMusicBrainz(query)
	if err != nil {
		ps.editCoverError = err.Error()
		return
	}
	ps.editCoverResults = results
	ps.editCoverResultIdx = 0
	ps.editCoverResultOffset = 0
	if len(results) == 0 {
		ps.editCoverError = "no covers found"
		return
	}
	ps.editCoverError = ""
}

func (ps *PlayerState) coverFixResultOffset() {
	h := ps.editCoverResultsHeight()
	if h <= 0 {
		return
	}
	ps.editCoverResultOffset = clampOffset(ps.editCoverResultOffset, ps.editCoverResultIdx, h, min(1, h/4), len(ps.editCoverResults))
}

func (ps *PlayerState) editCoverResultsHeight() int {
	rightHeight := (ps.windowHeight - 6) / 3
	// Subtract 2 for search bar line + separator
	return rightHeight - 2
}

func (ps *PlayerState) editCoverDir() string {
	return filepath.Join(ps.config.MusicDir, ps.editAlbumOrig[3])
}

func fetchReleaseTracks(mbid string) ([]metadataTrack, error) {
	u := "https://musicbrainz.org/ws/2/release/" + mbid + "?inc=recordings&fmt=json"

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
		Media []struct {
			Tracks []struct {
				Position  int `json:"position"`
				Recording struct {
					Title string `json:"title"`
				} `json:"recording"`
			} `json:"tracks"`
		} `json:"media"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("musicbrainz: %w", err)
	}

	var tracks []metadataTrack
	// Flatten all media (handle multi-disc releases)
	for _, medium := range data.Media {
		for _, t := range medium.Tracks {
			tracks = append(tracks, metadataTrack{
				position: t.Position,
				title:    t.Recording.Title,
			})
		}
	}

	return tracks, nil
}

func (ps *PlayerState) doMetadataSearch(query string) {
	results, err := searchMusicBrainz(query)
	if err != nil {
		ps.editMetadataError = err.Error()
		return
	}
	ps.editMetadataResults = results
	ps.editMetadataResultIdx = 0
	ps.editMetadataResultOffset = 0
	if len(results) == 0 {
		ps.editMetadataError = "no results found"
		return
	}
	ps.editMetadataError = ""
}

func (ps *PlayerState) metadataFixResultOffset() {
	h := ps.editMetadataResultsHeight()
	if h <= 0 {
		return
	}
	ps.editMetadataResultOffset = clampOffset(ps.editMetadataResultOffset, ps.editMetadataResultIdx, h, min(1, h/4), len(ps.editMetadataResults))
}

func (ps *PlayerState) editMetadataResultsHeight() int {
	rightHeight := (ps.windowHeight - 6) / 3
	// Subtract 2 for search bar line + separator
	return rightHeight - 2
}
