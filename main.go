package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	config, err := LoadConfig()
	if err != nil {
		fmt.Printf("Failed to load config: %v\n", err)
		os.Exit(1)
	}

	client, err := NewMPDClient(config.MPDHost, config.MPDPort)
	if err != nil {
		fmt.Printf("Failed to connect to mpd: %v\n", err)
		os.Exit(1)
	}
	defer client.Close()

	musicData, err := LoadMusicData(client)
	if err != nil {
		fmt.Printf("Failed to load music data: %v\n", err)
		os.Exit(1)
	}

	player, err := NewPlayerState(config, musicData, client)
	if err != nil {
		fmt.Printf("Failed to connect to music daemon: %v\n", err)
		os.Exit(1)
	}

	program := tea.NewProgram(player, tea.WithAltScreen())
	if _, err := program.Run(); err != nil {
		fmt.Printf("Error running program: %v\n", err)
		os.Exit(1)
	}
}
