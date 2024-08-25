package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	lg "github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
)

func titleStyle() lg.Style {
	b := lg.RoundedBorder()
	b.Right = "|"
	return lg.NewStyle().BorderStyle(b).Padding(0, 1)
}

func infoStyle() lg.Style {
	b := lg.RoundedBorder()
	b.Left = "|"
	return titleStyle().BorderStyle(b)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func (m model) headerView() string {
	title := titleStyle().Render(m.title)
	line := strings.Repeat("-", max(0, m.viewport.Width-lg.Width(title)))

	return lg.JoinHorizontal(lg.Center, title, line)
}

func (m model) footerView() string {
	info := infoStyle().Render(fmt.Sprintf("%3.f%%", m.viewport.ScrollPercent()*100))
	line := strings.Repeat("-", max(0, m.viewport.Width-lg.Width(info)))

	return lg.JoinHorizontal(lg.Center, line, info)
}

func (m *model) updateUIModel(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if k := msg.String(); k == "q" {
			return m, tea.Quit
		} else if (k == "right") && len(m.albumTitles) > 0 {
			m.index += 1

			if m.index > len(m.albumLyrics)-1 {
				m.index = 0
			}

			m.title = fmt.Sprintf("%s [%d/%d]", m.albumTitles[m.index], m.index+1, len(m.albumTitles))
			m.content = m.albumLyrics[m.index]
			m.viewport.SetContent(m.wordWrap())
			m.viewport.GotoTop()
		} else if (k == "left") && len(m.albumTitles) > 0 {
			m.index -= 1

			if m.index < 0 {
				m.index = len(m.albumLyrics) - 1
			}

			m.title = fmt.Sprintf("%s [%d/%d]", m.albumTitles[m.index], m.index+1, len(m.albumTitles))
			m.content = m.albumLyrics[m.index]
			m.viewport.SetContent(m.wordWrap())
			m.viewport.GotoTop()
		}

	case tea.WindowSizeMsg:
		headerHeight := lg.Height(m.headerView())
		footerHeight := lg.Height(m.footerView())
		verticalMarginHeight := headerHeight + footerHeight

		if !m.ready {
			m.viewport = viewport.New(msg.Width, msg.Height-verticalMarginHeight)
			m.viewport.YPosition = headerHeight
			m.viewport.HighPerformanceRendering = false
			m.viewport.SetContent(m.wordWrap())
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

func (m model) UIView() string {
	if !m.ready {
		return "\n Initializing..."
	}

	out := termenv.NewOutput(os.Stdout)
	out.SetWindowTitle(fmt.Sprintf("%s [%3.f%%]", m.title, m.viewport.ScrollPercent()*100))
	return fmt.Sprintf("%s\n%s\n%s", m.headerView(), m.viewport.View(), m.footerView())
}
