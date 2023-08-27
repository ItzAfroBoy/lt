package main

import (
	"flag"
	"fmt"
	"os"
	"path"

	"github.com/ItzAfroBoy/lt/input"
	"github.com/ItzAfroBoy/lt/input/fileloader"
	"github.com/ItzAfroBoy/lt/ui"
	"github.com/ItzAfroBoy/lt/ui/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

func userHomeDir() string {
	dir, _ := os.UserHomeDir()
	return dir
}

func main() {
	artist := flag.String("artist", "", "Artist of the song")
	title := flag.String("title", "", "Title of the song or album")
	skip := flag.Bool("skip", false, "Skip prompt")
	album := flag.Bool("album", false, "Album mode")
	raw := flag.Bool("raw", false, "Enable raw mode")
	save := flag.Bool("save", false, "Save lyrics")
	load := flag.Bool("load", false, "Load lyrics")

	flag.Parse()

	if *load {
		fm := fileloader.InitialModel()
		if _, err := tea.NewProgram(&fm).Run(); err != nil {
			fmt.Println("Couldn't run program:", err)
			os.Exit(1)
		}
		if fm.Exit {
			os.Exit(130)
		}

		
	}

	if !*skip {
		im := input.InitialModel(*artist, *title, *album)
		if _, err := tea.NewProgram(&im).Run(); err != nil {
			fmt.Println("Couldn't run program:", err)
			os.Exit(1)
		}
		if im.Exit {
			os.Exit(130)
		}

		*artist, *title, *album = im.Save()
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

		if !*raw {
			um := ui.InitialModel(title, lyrics)
			if _, err := tea.NewProgram(um, tea.WithAltScreen(), tea.WithMouseCellMotion()).Run(); err != nil {
				fmt.Println("Couldn't run program:", err)
				os.Exit(1)
			}
		} else {
			fmt.Printf("%s\n\n%s\n", title, lyrics)
		}

		if *save {
			outpath := path.Join(userHomeDir(), "Saved Lyrics", fmt.Sprintf("%s.txt", title))
			output := fmt.Sprintf("%s\n\n%s\n", title, lyrics)
			if err := os.WriteFile(outpath, []byte(output), 0755); err != nil {
				fmt.Println("Couldn't save lyrics:", err)
				os.Exit(1)
			} else {
				fmt.Printf("%s lyrics saved\n", title)
			}
		}
	} else {
		if len(sm.AlbumTitles) == 0 || len(sm.AlbumLyrics) == 0 {
			os.Exit(126)
		}

		titles, lyrics := sm.AlbumTitles, sm.AlbumLyrics
		um := ui.AlbumInitialModel(titles, lyrics)
		if _, err := tea.NewProgram(um, tea.WithAltScreen(), tea.WithMouseCellMotion()).Run(); err != nil {
			fmt.Println("Couldn't run program:", err)
			os.Exit(1)
		}
	}
}
