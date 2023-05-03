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
	title := flag.String("title", "", "Title of the song or album")
	skip := flag.Bool("skip", false, "Skip prompt")
	album := flag.Bool("album", false, "Album mode")

	flag.Parse()

	if !*skip {
		im := input.InitialModel(*artist, *title)
		if _, err := tea.NewProgram(&im).Run(); err != nil {
			fmt.Println("Couldn't run program:", err)
			os.Exit(1)
		}
		if im.Exit {
			os.Exit(130)
		}

		*artist, *title = im.Save()
	}

	if *artist == "" || *title == "" {
		os.Exit(126)
	}

	sm := spinner.InitialModel(*artist, *title, *album)
	if _, err := tea.NewProgram(&sm).Run(); err != nil {
		fmt.Println("Couldn't run program:", err)
		os.Exit(1)
	}

	if !*album {
		if sm.Title == "" || sm.Lyrics == "" {
			os.Exit(126)
		}

		title, lyrics := sm.Title, sm.Lyrics

		um := ui.InitialModel(title, lyrics)
		if _, err := tea.NewProgram(um, tea.WithAltScreen(), tea.WithMouseCellMotion()).Run(); err != nil {
			fmt.Println("Couldn't run program:", err)
			os.Exit(1)
		}
	} else {
		if len(sm.AlbumTitles) == 0 || len(sm.AlbumLyrics) == 0 {
			os.Exit(126)
		}

		titles, lyrics := sm.AlbumTitles, sm.AlbumLyrics

		um := ui.AlbumInitialModel(titles, lyrics)
		if _, err := tea.NewProgram(um, tea.WithAltScreen(),tea.WithMouseCellMotion()).Run(); err != nil {
			fmt.Println("Couldn't run program:", err)
			os.Exit(1)
		}
	}
}
