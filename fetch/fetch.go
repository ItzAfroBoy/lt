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

func Get(artist, song string) tea.Cmd {
	return func() tea.Msg {
		artist, song = parser.Args(artist, song)
		res, err := http.Get(fmt.Sprintf("https://genius.com/%s-%s-lyrics", artist, song))
		if err != nil {
			return ResMsg{"", ""}
		}

		defer res.Body.Close()
		body, err := io.ReadAll(res.Body)
		if err != nil {
			return ResMsg{"", ""}
		}

		data := string(body)
		title := parser.Title(data)
		lyrics := parser.Lyrics(data)

		return ResMsg{title, lyrics}
	}
}
