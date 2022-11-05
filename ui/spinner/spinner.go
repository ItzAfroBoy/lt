package spinner

import (
	"fmt"

	"github.com/ItzAfroBoy/lt/fetch"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	lg "github.com/charmbracelet/lipgloss"
)

type model struct {
	spinner spinner.Model
	artist string
	song string
	Title string
	Lyrics string
	Exit bool
}

func InitialModel(artist, song string) model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lg.NewStyle().Foreground(lg.Color("205"))
	
	m := model{spinner: s, artist: artist, song: song}
	return m
}

func (m model) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, fetch.Get(m.artist, m.song))
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
	str := fmt.Sprintf("%s Fetching lyrics...", m.spinner.View())
	if m.Exit {
		return fmt.Sprintf("%s Quitting...\n", m.spinner.View())
	}

	return str
}