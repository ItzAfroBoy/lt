package fetch

import (
	"fmt"
	"io"
	"net/http"

	"github.com/ItzAfroBoy/lt/input/parser"
	tea "github.com/charmbracelet/bubbletea"
)

type ResMsg struct {
	Title  string
	Lyrics string
}

type AlbumMsg struct {
	Titles     []string
	Lyrics     []string
	AlbumTitle string
}

func Get(url string) (string, string) {
	res, err := http.Get(url)
	if err != nil {
		return "Error", err.Error()
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "Error", err.Error()
	}

	data := string(body)
	title := parser.Title(data, false)
	lyrics := parser.Lyrics(data)

	return title, lyrics
}

func GetSong(artist, title string) tea.Cmd {
	return func() tea.Msg {
		artist, title = parser.FormatArgs(artist, title)
		url := fmt.Sprintf("https://genius.com/%s-%s-lyrics", artist, title)
		title, lyrics := Get(url)
		if title == "" && lyrics == "" {
			return ResMsg{"", ""}
		}

		return ResMsg{title, lyrics}
	}
}

func GetAlbum(artist, album string) tea.Cmd {
	return func() tea.Msg {
		artist, album = parser.FormatArgs(artist, album)
		res, err := http.Get(fmt.Sprintf("https://genius.com/albums/%s/%s", artist, album))
		if err != nil {
			return ResMsg{"Error", err.Error()}
		}

		defer res.Body.Close()
		body, err := io.ReadAll(res.Body)
		if err != nil {
			return ResMsg{"Error", err.Error()}
		}

		tracklist := parser.AlbumList(string(body))
		albumTitle := parser.Title(string(body), true)
		var titles, lyrics []string

		for _, song := range tracklist {
			title, lyric := Get(song)

			if title == "" && lyric == "" {
				continue
			}

			titles = append(titles, title)
			lyrics = append(lyrics, lyric)
		}

		return AlbumMsg{titles, lyrics, albumTitle}
	}
}
