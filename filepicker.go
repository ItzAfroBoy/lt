package main

import (
	"errors"
	"os"
	"path"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type clearErrorMsg struct{}

func clearErrorAfter(t time.Duration) tea.Cmd {
	return tea.Tick(t, func(_ time.Time) tea.Msg {
		return clearErrorMsg{}
	})
}

func (m *model) updateFPModel(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	m.filepicker, cmd = m.filepicker.Update(msg)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+a":
			files, _ := os.ReadDir(m.filepicker.CurrentDirectory)
			for _, v := range files {
				if strings.HasSuffix(v.Name(), ".lt") {
					title, content := parseFile(path.Join(m.filepicker.CurrentDirectory, v.Name()))
					m.albumTitles = append(m.albumTitles, title)
					m.albumLyrics = append(m.albumLyrics, content)
				}
			}
			*albumMode = true
			m.state = "ui"
			return m, tea.Batch(tea.EnterAltScreen, m.UIInit(), tea.WindowSize())
		case "q":
			return m, tea.Quit
		}
	}

	if didSelect, path := m.filepicker.DidSelectFile(msg); didSelect {
		*albumMode = false
		m.title, m.content = parseFile(path)
		m.state = "ui"

		if *raw {
			m.state = "raw"
			return m, nil
		}

		return m, tea.Batch(tea.EnterAltScreen, m.UIInit(), tea.WindowSize())
	}

	if didSelect, path := m.filepicker.DidSelectDisabledFile(msg); didSelect {
		m.err = errors.New(path + " is not valid")
		return m, tea.Batch(cmd, clearErrorAfter(2*time.Second))
	}

	return m, cmd
}

func (m model) FPView() string {
	var s strings.Builder
	s.WriteString("\n  ")
	if m.err != nil {
		s.WriteString(m.filepicker.Styles.DisabledFile.Render(m.err.Error()))
	} else {
		s.WriteString("Pick a file or presss ctrl+a to show all lyrics:")
	}
	s.WriteString("\n\n" + m.filepicker.View() + "\n")
	return s.String()
}
