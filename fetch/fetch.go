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
	Titles []string
	Lyrics []string
}

func Get(url string) (string, string) {
	res, err := http.Get(url)
	if err != nil {
		return "", ""
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", ""
	}

	data := string(body)
	title := parser.Title(data)
	lyrics := parser.Lyrics(data)

	return title, lyrics
}

func GetSong(artist, song string) tea.Cmd {
	return func() tea.Msg {
		artist, song = parser.FormatArgs(artist, song)
		url := fmt.Sprintf("https://genius.com/%s-%s-lyrics", artist, song)
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
			return ResMsg{"", ""}
		}

		defer res.Body.Close()
		body, err := io.ReadAll(res.Body)
		if err != nil {
			return ResMsg{"", ""}
		}

		data := string(body)
		tracklist := parser.AlbumList(data)
		var titles, lyrics []string

		for _, i := range tracklist {
			t, l := Get(i)

			if t == "" && l == "" {
				continue
			}

			titles = append(titles, t)
			lyrics = append(lyrics, l)
		}

		return AlbumMsg{titles, lyrics}
	}
}
