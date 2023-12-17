package main

import (
	"flag"
	"fmt"
	"os"
	"path"

	"github.com/ItzAfroBoy/lt/input"
	"github.com/ItzAfroBoy/lt/input/fileloader"
	"github.com/ItzAfroBoy/lt/input/parser"
	"github.com/ItzAfroBoy/lt/ui"
	"github.com/ItzAfroBoy/lt/ui/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

func getSavedLyrics() (string, string) {
	fm := fileloader.InitialModel()
	if _, err := tea.NewProgram(&fm).Run(); err != nil {
		fmt.Println("Couldn't run program:", err)
		os.Exit(1)
	}
	if fm.Exit {
		os.Exit(130)
	}

	body, err := os.ReadFile(fm.SelectedFile)
	if err != nil {
		fmt.Println("Couldn't open file:", err)
		os.Exit(1)
	}

	title, lyrics := parser.RawLyrics(string(body))
	return title, lyrics
}

func saveLyrics(title, lyrics string) {
	outpath := path.Join(parser.UserHomeDir(), "Saved Lyrics", fmt.Sprintf("%s.lt", title))
	output := fmt.Sprintf("%s\n\n%s\n", title, lyrics)
	if err := os.MkdirAll(path.Join(parser.UserHomeDir(), "Saved Lyrics"), 0o755); err != nil {
		fmt.Println("Couldn't create directory:", err)
		os.Exit(1)
	}

	if err := os.WriteFile(outpath, []byte(output), 0o755); err != nil {
		fmt.Println("Couldn't save lyrics:", err)
		os.Exit(1)
	}

	fmt.Printf("%s lyrics saved\n", title)
}

func displayLyrics(title, lyrics string) {
	um := ui.InitialModel(title, lyrics)
	if _, err := tea.NewProgram(um, tea.WithAltScreen(), tea.WithMouseCellMotion()).Run(); err != nil {
		fmt.Println("Couldn't run program:", err)
		os.Exit(1)
	}
}

func displayLyricsRaw(title, lyrics string) {
	fmt.Printf("%s\n\n%s\n", title, lyrics)
}

func main() {
	artist := flag.String("artist", "", "Artist of the song")
	title := flag.String("title", "", "Title of the song or album")
	skip := flag.Bool("skip", false, "Skip prompt")
	album := flag.Bool("album", false, "Album mode")
	raw := flag.Bool("raw", false, "Enable raw mode")
	save := flag.Bool("save", false, "Save lyrics")
	load := flag.Bool("load", false, "Load lyrics")
	var lyrics string

	flag.Parse()

	if *load {
		*title, lyrics = getSavedLyrics()
	}

	if !*skip && !*load {
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

	sm := spinner.InitialModel(*artist, *title, *album)
	if !*load {
		if _, err := tea.NewProgram(&sm).Run(); err != nil {
			fmt.Println("Couldn't run program:", err)
			os.Exit(1)
		}
	}

	if !*album {
		if sm.Title == "" || sm.Lyrics == "" {
			os.Exit(126)
		}
		*title, lyrics = sm.Title, sm.Lyrics

		if *raw {
			displayLyricsRaw(*title, lyrics)
		} else {
			displayLyrics(*title, lyrics)
		}

		if *save {
			saveLyrics(*title, lyrics)
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
