package main

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

func (m model) spinnerInit() tea.Cmd {
	if *albumMode {
		return tea.Batch(m.spinner.Tick, getAlbum)
	}
	return tea.Batch(m.spinner.Tick, getSong)
}

func (m *model) updateSpinnerModel(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case albumMsg:
		m.albumTitles = msg.titles
		m.albumLyrics = msg.lyrics
		m.title = m.albumTitles[0]
		m.content = m.albumLyrics[0]
		m.state = "ui"
		return m, tea.Batch(tea.EnterAltScreen, m.UIInit(), tea.WindowSize())
	case resMsg:
		m.title = msg.title
		m.content = msg.lyrics
		if *raw {
			m.state = "raw"
			return m, nil
		}

		m.state = "ui"
		return m, tea.Batch(tea.EnterAltScreen, m.UIInit(), tea.WindowSize())
	default:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)

		return m, cmd
	}
}

func (m model) spinnerView() string {
	var str string
	if !m.exit {
		str = fmt.Sprintf("%s Fetching lyrics...", m.spinner.View())
	} else {
		str = fmt.Sprintf("%s Quitting...\n", m.spinner.View())
	}

	return str
}
