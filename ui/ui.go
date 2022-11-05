package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	lg "github.com/charmbracelet/lipgloss"
)

func titleStyle() lg.Style {
	b := lg.RoundedBorder()
	b.Right = "|"

	return lg.NewStyle().BorderStyle(b).Padding(0, 1)
}

func infoStyle() lg.Style {
	b := lg.RoundedBorder()
	b.Left = "|"

	return titleStyle().Copy().BorderStyle(b)
}

func max(a, b int) int {
	if a > b {
		return a
	}

	return b
}

type Model struct {
	Title    string
	Content  string
	ready    bool
	viewport viewport.Model
}

func (m Model) headerView() string {
	title := titleStyle().Render(m.Title)
	line := strings.Repeat("-", max(0, m.viewport.Width-lg.Width(title)))

	return lg.JoinHorizontal(lg.Center, title, line)
}

func (m Model) footerView() string {
	info := infoStyle().Render(fmt.Sprintf("%3.f%%", m.viewport.ScrollPercent()*100))
	line := strings.Repeat("-", max(0, m.viewport.Width-lg.Width(info)))

	return lg.JoinHorizontal(lg.Center, line, info)
}

func InitialModel(title, lyrics string) Model {
	m := Model{Title: title, Content: lyrics}
	return m
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {

	case tea.KeyMsg:
		if k := msg.String(); k == "ctrl+c" || k == "q" || k == "esc" {
			return m, tea.Quit
		}

	case tea.WindowSizeMsg:
		headerHeight := lg.Height(m.headerView())
		footerHeight := lg.Height(m.footerView())
		verticalMarginHeight := headerHeight + footerHeight

		if !m.ready {
			m.viewport = viewport.New(msg.Width, msg.Height-verticalMarginHeight)
			m.viewport.YPosition = headerHeight
			m.viewport.HighPerformanceRendering = false
			m.viewport.SetContent(m.Content)
			m.ready = true
		} else {
			m.viewport.Width = msg.Width
			m.viewport.Height = msg.Height - verticalMarginHeight
		}
	}

	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	if !m.ready {
		return "\n Initializing..."
	}

	return fmt.Sprintf("%s\n%s\n%s", m.headerView(), m.viewport.View(), m.footerView())
}
