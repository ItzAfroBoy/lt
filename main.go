package main

import (
	"flag"
	"fmt"
	"os"
	"path"

	"github.com/charmbracelet/bubbles/filepicker"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	lg "github.com/charmbracelet/lipgloss"
)

var artist *string
var title *string
var albumMode *bool
var raw *bool
var save *bool
var load *bool
var p *tea.Program

type model struct {
	filepicker  filepicker.Model
	err         error
	focusIndex  int
	inputs      []textinput.Model
	spinner     spinner.Model
	index       int
	ready       bool
	viewport    viewport.Model
	content     string
	albumTitles []string
	albumLyrics []string
	title       string
	state       string
}

var (
	focusedStyle       = lg.NewStyle().Foreground(lg.Color("205"))
	blurredStyle       = lg.NewStyle().Foreground(lg.Color("240"))
	cursorStyle        = focusedStyle
	noStyle            = lg.NewStyle()
	helpStyle          = blurredStyle
	albumModeHelpStyle = lg.NewStyle().Foreground(lg.Color("244"))

	focusedButton = focusedStyle.Render("[ Submit ]")
	blurredButton = blurredStyle.Render("[ Submit ]")
)

func init() {
	artist = flag.String("artist", "none", "Artist of the song to fetch")
	title = flag.String("title", "none", "Title of the song to fetch")
	albumMode = flag.Bool("album", false, "Fetch lyrics for all songs in an album")
	raw = flag.Bool("raw", false, "Show the raw text to the terminal")
	save = flag.Bool("export", false, "Save your lyrics to a LT file")
	load = flag.Bool("import", false, "Load your lyrics from an LT file")

	flag.Parse()
	if *raw && *albumMode {
		os.Exit(1)
	}
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

	fp := filepicker.New()
	fp.AllowedTypes = []string{".lt"}
	fp.CurrentDirectory = path.Join(userHomeDir(), "Saved Lyrics")
	m.filepicker = fp

	return m
}

func (m *model) Init() tea.Cmd {
	if *load {
		m.state = "filepicker"
		return m.filepicker.Init()
	} else if *artist != "none" && *title != "none" {
		m.state = "spinner"
		formatArgs()
		return m.spinnerInit()
	}

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
	case clearErrorMsg:
		m.err = nil
	}

	switch m.state {
	case "filepicker":
		_m, cmd = m.updateFPModel(msg)
		m = _m.(*model)
	case "input":
		_m, cmd = m.updateInputsModel(msg)
		m = _m.(*model)
	case "spinner":
		_m, cmd = m.updateSpinnerModel(msg)
		m = _m.(*model)
	case "ui":
		_m, cmd = m.updateUIModel(msg)
		m = _m.(*model)
	case "raw":
		return m, tea.Sequence(tea.ClearScreen, tea.Println(m.title, m.content), tea.Quit)
	}

	return m, cmd
}

func (m model) View() string {
	switch m.state {
	case "filepicker":
		return m.FPView()
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

func (m model) saveLyrics() {
	outpath := path.Join(userHomeDir(), "Saved Lyrics", fmt.Sprintf("%s.lt", m.title))
	output := fmt.Sprintf("%s\n\n%s\n", m.title, m.content)
	if err := os.MkdirAll(path.Join(userHomeDir(), "Saved Lyrics"), 0o755); err != nil {
		fmt.Println("Couldn't create directory:", err)
		os.Exit(1)
	}

	if err := os.WriteFile(outpath, []byte(output), 0o755); err != nil {
		fmt.Println("Couldn't save lyrics:", err)
		os.Exit(1)
	}
}
