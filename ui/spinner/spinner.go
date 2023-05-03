package spinner

import (
	"fmt"

	"github.com/ItzAfroBoy/lt/fetch"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	lg "github.com/charmbracelet/lipgloss"
)

type model struct {
	spinner     spinner.Model
	artist      string
	album       bool
	AlbumTitles []string
	AlbumLyrics []string
	Title       string
	Lyrics      string
	Exit        bool
}

func InitialModel(artist, title string, album bool) model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lg.NewStyle().Foreground(lg.Color("205"))

	m := model{spinner: s, artist: artist, Title: title, album: album}
	return m
}

func (m model) Init() tea.Cmd {
	if m.album {
		return tea.Batch(m.spinner.Tick, fetch.GetAlbum(m.artist, m.Title))
	} else {
		return tea.Batch(m.spinner.Tick, fetch.GetSong(m.artist, m.Title))
	}
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "ctrl+c", "ctrl+q":
			m.Exit = true
			return m, tea.Quit
		default:
			return m, nil
		}

	case fetch.AlbumMsg:
		m.AlbumTitles = msg.Titles
		m.AlbumLyrics = msg.Lyrics
		return m, tea.Quit

	case fetch.ResMsg:
		m.Title = msg.Title
		m.Lyrics = msg.Lyrics
		return m, tea.Quit

	default:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)

		return m, cmd
	}
}

func (m model) View() string {
	var str string

	if !m.Exit {
		str = fmt.Sprintf("%s Fetching lyrics...", m.spinner.View())
	} else {
		str = fmt.Sprintf("%s Quitting...\n", m.spinner.View())
	}

	return str
}
