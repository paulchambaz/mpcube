package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/BurntSushi/toml"
)

const version = "1.0.0"

type Config struct {
	MPDHost string `toml:"mpd_host"`
	MPDPort int    `toml:"mpd_port"`

	VolumeStep   int `toml:"volume_step"`
	SeekDuration int `toml:"seek_duration"`

	TickInterval  int `toml:"tick_interval"`
	MaxRetryDelay int `toml:"max_retry_delay"`

	ScrollPadding      int `toml:"scroll_padding"`
	WideThreshold      int `toml:"wide_threshold"`
	AlbumWidth         int `toml:"album_width"`
	VolumeBarThreshold int `toml:"volume_bar_threshold"`
	VolumeBarWidth     int `toml:"volume_bar_width"`

	MusicDir string `toml:"music_dir"`
}

func DefaultConfig() Config {
	return Config{
		MPDHost:            "127.0.0.1",
		MPDPort:            6600,
		VolumeStep:         10,
		SeekDuration:       5000,
		TickInterval:       100,
		MaxRetryDelay:      30,
		ScrollPadding:      5,
		WideThreshold:      100,
		AlbumWidth:         40,
		VolumeBarThreshold: 90,
		VolumeBarWidth:     30,
	}
}

func LoadConfig() (Config, error) {
	cfg := DefaultConfig()

	home, err := os.UserHomeDir()
	if err == nil && cfg.MusicDir == "" {
		cfg.MusicDir = filepath.Join(home, "music")
	}
	if err == nil {
		configPath := filepath.Join(home, ".config", "mpcube", "config.toml")
		if _, err := os.Stat(configPath); err == nil {
			if _, err := toml.DecodeFile(configPath, &cfg); err != nil {
				return cfg, fmt.Errorf("could not parse %s: %w", configPath, err)
			}
		}
	}

	if v := os.Getenv("MPD_HOST"); v != "" {
		cfg.MPDHost = v
	}
	if v := os.Getenv("MPD_PORT"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			cfg.MPDPort = n
		}
	}
	if v := os.Getenv("MPCUBE_VOLUME_STEP"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			cfg.VolumeStep = n
		}
	}
	if v := os.Getenv("MPCUBE_SEEK_DURATION"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			cfg.SeekDuration = n
		}
	}
	if v := os.Getenv("MPCUBE_TICK_INTERVAL"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			cfg.TickInterval = n
		}
	}
	if v := os.Getenv("MPCUBE_MAX_RETRY_DELAY"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			cfg.MaxRetryDelay = n
		}
	}
	if v := os.Getenv("MPCUBE_SCROLL_PADDING"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			cfg.ScrollPadding = n
		}
	}
	if v := os.Getenv("MPCUBE_WIDE_THRESHOLD"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			cfg.WideThreshold = n
		}
	}
	if v := os.Getenv("MPCUBE_ALBUM_WIDTH"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			cfg.AlbumWidth = n
		}
	}
	if v := os.Getenv("MPCUBE_VOLUME_BAR_THRESHOLD"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			cfg.VolumeBarThreshold = n
		}
	}
	if v := os.Getenv("MPCUBE_VOLUME_BAR_WIDTH"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			cfg.VolumeBarWidth = n
		}
	}
	if v := os.Getenv("MPCUBE_MUSIC_DIR"); v != "" {
		cfg.MusicDir = v
	}

	var showVersion bool
	flag.StringVar(&cfg.MPDHost, "mpd-host", cfg.MPDHost, "MPD host address")
	flag.IntVar(&cfg.MPDPort, "mpd-port", cfg.MPDPort, "MPD port number")
	flag.IntVar(&cfg.VolumeStep, "volume-step", cfg.VolumeStep, "volume adjustment step")
	flag.IntVar(&cfg.SeekDuration, "seek-duration", cfg.SeekDuration, "seek duration in milliseconds")
	flag.IntVar(&cfg.TickInterval, "tick-interval", cfg.TickInterval, "UI refresh interval in milliseconds")
	flag.IntVar(&cfg.MaxRetryDelay, "max-retry-delay", cfg.MaxRetryDelay, "max MPD reconnect delay in seconds")
	flag.IntVar(&cfg.ScrollPadding, "scroll-padding", cfg.ScrollPadding, "scroll padding lines")
	flag.IntVar(&cfg.WideThreshold, "wide-threshold", cfg.WideThreshold, "terminal width for wide layout")
	flag.IntVar(&cfg.AlbumWidth, "album-width", cfg.AlbumWidth, "album panel width in wide layout")
	flag.IntVar(&cfg.VolumeBarThreshold, "volume-bar-threshold", cfg.VolumeBarThreshold, "terminal width for wide volume bar")
	flag.IntVar(&cfg.VolumeBarWidth, "volume-bar-width", cfg.VolumeBarWidth, "volume bar width in wide layout")
	flag.StringVar(&cfg.MusicDir, "music-dir", cfg.MusicDir, "MPD music directory path")
	flag.BoolVar(&showVersion, "version", false, "print version")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: mpcube [OPTIONS]\n\nOptions:\n")
		flag.VisitAll(func(f *flag.Flag) {
			prefix := "--"
			if len(f.Name) == 1 {
				prefix = "-"
			}
			if f.DefValue != "" && f.DefValue != "false" {
				fmt.Fprintf(os.Stderr, "  %s%s\n\t%s (default: %s)\n", prefix, f.Name, f.Usage, f.DefValue)
			} else {
				fmt.Fprintf(os.Stderr, "  %s%s\n\t%s\n", prefix, f.Name, f.Usage)
			}
		})
	}
	flag.Parse()

	if showVersion {
		fmt.Println(version)
		os.Exit(0)
	}

	return cfg, nil
}
