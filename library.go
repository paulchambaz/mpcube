package main

import "time"

type MusicData struct {
	Albums []Album
}

type Album struct {
	Artist string
	Album  string
	Date   int
	Songs  []Song
}

type Song struct {
	URI      string
	Title    string
	Track    int
	Duration time.Duration
}
