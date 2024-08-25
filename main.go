package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	lg "github.com/charmbracelet/lipgloss"
)

var artist *string
var title *string
var albumMode *bool
// var raw *bool
// var save *bool
var p *tea.Program

type model struct {
	focusIndex int
	inputs     []textinput.Model
	spinner spinner.Model
	index    int
	ready    bool
	viewport viewport.Model
	content  string
	albumTitles []string
	albumLyrics []string
	title       string
	state       string
	exit        bool
}

var (
	focusedStyle       = lg.NewStyle().Foreground(lg.Color("205"))
	blurredStyle       = lg.NewStyle().Foreground(lg.Color("240"))
	cursorStyle        = focusedStyle
	noStyle            = lg.NewStyle()
	helpStyle          = blurredStyle
	albumModeHelpStyle = lg.NewStyle().Foreground(lg.Color("244"))

	focusedButton = focusedStyle.Render("[ Submit ]")
	blurredButton = fmt.Sprintf("[ %s ]", blurredStyle.Render("Submit"))
)

func init() {
	artist = flag.String("artist", "none", "Artist of the song to fetch")
	title = flag.String("title", "none", "Title of the song to fetch")
	albumMode = flag.Bool("album", false, "Fetch lyrics for all songs in an album")
	// raw = flag.Bool("raw", false, "Show the raw text to the terminal")
	// save = flag.Bool("export", false, "Save your lyrics to a LT file")

	flag.Parse()
}

func main() {
	m := initalModel()
	p = tea.NewProgram(&m, tea.WithMouseAllMotion())
	if _, err := p.Run(); err != nil {
		fmt.Println("Couldn't run program:", err)
		os.Exit(1)
	}
}

func initalModel() model {
	m := model{}
	m.inputs = make([]textinput.Model, 2)
	
	for i := range m.inputs {
		t := textinput.New()
		t.Cursor.Style = cursorStyle

		switch i {
		case 0:
			t.Placeholder = "Artist"
			t.Focus()
			t.PromptStyle = focusedStyle
			t.TextStyle = focusedStyle
			if *artist != "none" {
				t.SetValue(*artist)
			}
		case 1:
			t.Placeholder = "Title"
			if *title != "none" {
				t.SetValue(*title)
			}
		}

		m.inputs[i] = t
	}

	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lg.NewStyle().Foreground(lg.Color("205"))
	m.spinner = s

	return m
}

func (m *model) Init() tea.Cmd {
	m.state = "input"
	return textinput.Blink
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var _m tea.Model
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "ctrl+c":
			return m, tea.Quit
		}
	}

	switch m.state {
	case "input":
		_m, cmd = m.updateInputsModel(msg)
		m = _m.(*model)
	case "spinner":
		_m, cmd = m.updateSpinnerModel(msg)
		m = _m.(*model)
	case "ui":
		_m, cmd = m.updateUIModel(msg)
		m = _m.(*model)
	}

	if m.exit {
		return m, tea.Quit
	}

	return m, cmd
}

func (m model) View() string {
	switch m.state {
	case "input":
		return m.inputsView()
	case "spinner":
		return m.spinnerView()
	case "ui":
		return m.UIView()
	default:
		return ""
	}
}
