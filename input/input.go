package input

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	lg "github.com/charmbracelet/lipgloss"
)

var (
	focusedStyle       = lg.NewStyle().Foreground(lg.Color("205"))
	blurredStyle       = lg.NewStyle().Foreground(lg.Color("240"))
	cursorStyle        = focusedStyle.Copy()
	noStyle            = lg.NewStyle()
	helpStyle          = blurredStyle.Copy()
	albumModeHelpStyle = lg.NewStyle().Foreground(lg.Color("244"))

	focusedButton = focusedStyle.Copy().Render("[ Submit ]")
	blurredButton = fmt.Sprintf("[ %s ]", blurredStyle.Render("Submit"))
)

type Model struct {
	focusIndex int
	inputs     []textinput.Model
	albumMode  bool
	Exit       bool
}

func InitialModel(artist, title string, albumMode bool) Model {
	m := Model{
		inputs:    make([]textinput.Model, 2),
		albumMode: albumMode,
	}

	for i := range m.inputs {
		t := textinput.New()
		t.Cursor.Style = cursorStyle

		switch i {
		case 0:
			t.Placeholder = "Artist"
			t.Focus()
			t.PromptStyle = focusedStyle
			t.TextStyle = focusedStyle
			if artist != "" {
				t.SetValue(artist)
			}

		case 1:
			t.Placeholder = "Title"
			if title != "" {
				t.SetValue(title)

			}
		}

		m.inputs[i] = t
	}

	return m
}

func (m *Model) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))

	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "ctrl+q", "esc":
			m.Exit = true
			return m, tea.Quit

		case "ctrl+r":
			if m.albumMode {
				m.albumMode = false
			} else {
				m.albumMode = true
			}

		case "tab", "shift+tab", "enter", "up", "down":
			s := msg.String()

			if s == "enter" && m.focusIndex == len(m.inputs) {
				return m, tea.Quit
			}

			if s == "up" || s == "shift+tab" {
				m.focusIndex--
			} else {
				m.focusIndex++
			}

			if m.focusIndex > len(m.inputs) {
				m.focusIndex = 0
			} else if m.focusIndex < 0 {
				m.focusIndex = len(m.inputs)
			}

			cmds := make([]tea.Cmd, len(m.inputs))
			for i := 0; i <= len(m.inputs)-1; i++ {
				if i == m.focusIndex {
					cmds[i] = m.inputs[i].Focus()
					m.inputs[i].PromptStyle = focusedStyle
					m.inputs[i].TextStyle = focusedStyle
					continue
				}

				m.inputs[i].Blur()
				m.inputs[i].PromptStyle = noStyle
				m.inputs[i].TextStyle = noStyle
			}

			return m, tea.Batch(cmds...)
		}
	}

	cmd := m.updateInputs(msg)
	return m, cmd
}

func (m Model) View() string {
	var b strings.Builder

	for i := range m.inputs {
		b.WriteString(m.inputs[i].View())
		if i < len(m.inputs)-1 {
			b.WriteRune('\n')
		}
	}

	button := &blurredButton
	if m.focusIndex == len(m.inputs) {
		button = &focusedButton
	}
	fmt.Fprintf(&b, "\n\n%s\n\n", *button)

	b.WriteString(helpStyle.Render("album mode is "))
	b.WriteString(albumModeHelpStyle.Render(m.albumModeString()))
	b.WriteString(helpStyle.Render(" (ctrl+r to change mode)"))

	return b.String()
}

func (m Model) albumModeString() string {
	if m.albumMode {
		return "on"
	}
	return "off"
}

func (m Model) Save() (string, string, bool) {
	return m.inputs[0].Value(), m.inputs[1].Value(), m.albumMode
}
