package main

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/fhs/gompd/v2/mpd"
)

type rawEntry struct {
	uri      string
	title    string
	artist   string
	album    string
	date     int
	track    int
	duration time.Duration
}

func NewMPDClient(host string, port int) (*mpd.Client, error) {
	address := fmt.Sprintf("%s:%d", host, port)
	client, err := mpd.Dial("tcp", address)
	if err != nil {
		return nil, fmt.Errorf("could not connect to mpd at %s: %w", address, err)
	}
	return client, nil
}

func LoadMusicData(client *mpd.Client) (*MusicData, error) {
	attrs, err := client.ListAllInfo("")
	if err != nil {
		return nil, fmt.Errorf("could not list songs: %w", err)
	}

	var rawEntries []rawEntry
	for _, attr := range attrs {
		if attr["file"] == "" {
			continue
		}
		raw := parseRawEntry(attr)
		rawEntries = append(rawEntries, raw)
	}

	albumMap := make(map[string]*Album)
	for _, raw := range rawEntries {
		albumKey := fmt.Sprintf("%s|||%s", raw.artist, raw.album)

		album, exists := albumMap[albumKey]
		if !exists {
			album = &Album{
				Artist: raw.artist,
				Album:  raw.album,
				Date:   raw.date,
				Songs:  []Song{},
			}
			albumMap[albumKey] = album
		}

		if len(raw.artist) < len(album.Artist) {
			album.Artist = raw.artist
		}

		album.Songs = append(album.Songs, Song{
			URI:      raw.uri,
			Title:    raw.title,
			Track:    raw.track,
			Duration: raw.duration,
		})
	}

	var albums []Album
	for _, album := range albumMap {
		albums = append(albums, *album)
	}

	sortAlbums(albums)
	for i := range albums {
		sortSongs(albums[i].Songs)
	}

	return &MusicData{Albums: albums}, nil
}

func parseRawEntry(attr mpd.Attrs) rawEntry {
	uri := attr["file"]

	title := attr["Title"]
	if title == "" {
		title = uri
	}

	artist := attr["Artist"]
	if artist == "" {
		artist = "Unknown Artist"
	}

	album := attr["Album"]
	if album == "" {
		album = "Unknown Album"
	}

	date := 0
	if dateStr := attr["Date"]; dateStr != "" {
		if yearStr := strings.Split(dateStr, "-")[0]; yearStr != "" {
			if parsed, err := strconv.Atoi(yearStr); err == nil {
				date = parsed
			}
		}
	}

	track := 0
	if trackStr := attr["Track"]; trackStr != "" {
		if trackNum := strings.Split(trackStr, "/")[0]; trackNum != "" {
			if parsed, err := strconv.Atoi(trackNum); err == nil {
				track = parsed
			}
		}
	}

	var duration time.Duration
	if timeStr := attr["Time"]; timeStr != "" {
		if seconds, err := strconv.Atoi(timeStr); err == nil {
			duration = time.Duration(seconds) * time.Second
		}
	}

	return rawEntry{
		uri:      uri,
		title:    title,
		artist:   artist,
		album:    album,
		date:     date,
		track:    track,
		duration: duration,
	}
}

func sortAlbums(albums []Album) {
	sort.Slice(albums, func(i, j int) bool {
		a, b := albums[i], albums[j]

		if a.Artist == "" && b.Artist != "" {
			return true
		}

		if a.Artist != "" && b.Artist == "" {
			return false
		}

		aArtist := strings.TrimPrefix(a.Artist, "The ")
		bArtist := strings.TrimPrefix(b.Artist, "The ")

		if aArtist != bArtist {
			return aArtist < bArtist
		}

		if a.Date != b.Date {
			return a.Date < b.Date
		}

		return a.Album < b.Album
	})
}

func sortSongs(songs []Song) {
	sort.Slice(songs, func(i, j int) bool {
		a, b := songs[i], songs[j]

		if a.Track != b.Track {
			return a.Track < b.Track
		}

		return a.URI < b.URI
	})
}
