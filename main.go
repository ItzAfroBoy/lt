package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/ItzAfroBoy/lt/input"
	"github.com/ItzAfroBoy/lt/ui"
	"github.com/ItzAfroBoy/lt/ui/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	artist := flag.String("artist", "", "Artist of the song")
	song := flag.String("song", "", "Name of the song")
	skip := flag.Bool("skip", false, "Skip prompt")

	flag.Parse()

	if !*skip {
		im := input.InitialModel(*artist, *song)
		if err := tea.NewProgram(&im).Start(); err != nil {
			fmt.Println("Couldn't run program:", err)
			os.Exit(1)
		}
		if im.Exit {
			os.Exit(130)
		}

		*artist, *song = im.Save()
	}
	
	if *artist == "" || *song == "" {
		os.Exit(126)
	}
	
	sm := spinner.InitialModel(*artist, *song)
	if err := tea.NewProgram(&sm).Start(); err != nil {
		fmt.Println("Couldn't run program:", err)
		os.Exit(1)
	}

	if sm.Title == "" || sm.Lyrics == "" {
		os.Exit(126)
	}

	title, lyrics := sm.Title, sm.Lyrics

	um := ui.InitialModel(title, lyrics)
	if err := tea.NewProgram(um, tea.WithAltScreen(), tea.WithMouseCellMotion()).Start(); err != nil {
		fmt.Println("Couldn't run program:", err)
		os.Exit(1)
	}
}
