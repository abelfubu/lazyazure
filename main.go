package main

import (
	"fmt"
	"lazyaz/internal/models"
	pullrequests "lazyaz/internal/pull-requests"
	workitems "lazyaz/internal/work-items"
	"log"
	"strings"

	"github.com/atotto/clipboard"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
	"github.com/pkg/browser"
)

var (
	normal    = lipgloss.Color("#EEEEEE")
	subtle    = lipgloss.Color("#99A9C9")
	highlight = lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"}
	special   = lipgloss.AdaptiveColor{Light: "#43BF6D", Dark: "#73F59F"}

	base = lipgloss.NewStyle().Foreground(normal)

	activeTabBorder = lipgloss.Border{
		Top:         "─",
		Bottom:      " ",
		Left:        "│",
		Right:       "│",
		TopLeft:     "╭",
		TopRight:    "╮",
		BottomLeft:  "┘",
		BottomRight: "└",
	}

	tabBorder = lipgloss.Border{
		Top:         "─",
		Bottom:      "─",
		Left:        "│",
		Right:       "│",
		TopLeft:     "╭",
		TopRight:    "╮",
		BottomLeft:  "┴",
		BottomRight: "┴",
	}

	tabGap = tab.
		BorderTop(false).
		BorderLeft(false).
		BorderRight(false)

	tabBase = lipgloss.NewStyle().
		BorderForeground(highlight).
		Padding(0, 1)

	tab = tabBase.
		Foreground(subtle).
		Border(tabBorder, true)

	activeTab = tabBase.Border(activeTabBorder, true).Foreground(normal).Bold(true)

	tabs = []string{"(W) Work Items", "(P) Pull Requests"}
)

type Model struct {
	list         list.Model
	preview      viewport.Model
	activePane   int
	width        int
	height       int
	selectedItem int
	renderer     *glamour.TermRenderer
	tabIndex     int
}

func initialModel() Model {
	l := list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)
	l.Title = "Work Items"
	l.SetShowStatusBar(false)
	l.SetShowTitle(false)
	l.StartSpinner()

	vp := viewport.New(0, 0)
	vp.Style = lipgloss.NewStyle().
		BorderLeft(true).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("62")).
		Padding(0)

	renderer, _ := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(80),
	)

	return Model{
		list:       l,
		preview:    vp,
		activePane: 0,
		renderer:   renderer,
		tabIndex:   0,
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(textinput.Blink, workitems.FetchWorkItems)
}

func handleResponseMsg[T list.Item](m *Model, msg []T) tea.Cmd {
	var items []list.Item

	m.list.SetItems([]list.Item{})
	m.list.ResetFilter()

	for _, item := range msg {
		items = append(items, item)
	}
	cmd := m.list.SetItems(items)

	if i, ok := m.list.SelectedItem().(models.UiItem); ok {
		m.preview.SetContent(i.GetPreview(m.renderer))
		m.selectedItem = i.GetID()
	}

	return cmd
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		horizontalMargin := 4
		listWidth := (m.width / 2) - horizontalMargin
		previewWidth := m.width - listWidth - horizontalMargin*2
		m.list.SetSize(listWidth, m.height-4)
		m.preview.Width = previewWidth
		m.preview.Height = m.height - 4

	case tea.KeyMsg:
		if m.list.FilterState() == list.Filtering {
			break
		}

		switch msg.String() {
		case "ctrl+y":
			if i, ok := m.list.SelectedItem().(models.UiItem); ok {
				if err := clipboard.WriteAll(fmt.Sprintf("%d", i.GetID())); err != nil {
					log.Fatalf("Failed to copy to clipboard: %v", err)
				}
			}
		case "ctrl+c", "q":
			return m, tea.Quit
		case "enter":
			if i, ok := m.list.SelectedItem().(models.UiItem); ok {
				browser.OpenURL(i.GetURL())
			}
		case "w":
			m.tabIndex = 0
			cmds = append(cmds, workitems.FetchWorkItems)
		case "p":
			m.tabIndex = 1
			cmds = append(cmds, pullrequests.FetchPullRequests)
		}

	case workitems.WorkItemsResponseMsg:
		return m, handleResponseMsg(&m, msg)
	case pullrequests.PullRequestResponseMsg:
		return m, handleResponseMsg(&m, msg)
	}

	newListModel, cmd := m.list.Update(msg)
	m.list = newListModel
	cmds = append(cmds, cmd)

	if i, ok := m.list.SelectedItem().(models.UiItem); ok {
		m.preview.SetContent(i.GetPreview(m.renderer))
		m.selectedItem = i.GetID()
	}

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	var tabViews []string

	for i, text := range tabs {
		if i == m.tabIndex {
			tabViews = append(tabViews, activeTab.Render(text))
		} else {
			tabViews = append(tabViews, tab.Render(text))
		}
	}

	tabView := lipgloss.JoinHorizontal(lipgloss.Top, tabViews...)
	gap := tabGap.Render(strings.Repeat(" ", max(0, m.width-lipgloss.Width(tabView)-8)))
	tabView = lipgloss.JoinHorizontal(lipgloss.Bottom, tabView, gap)

	listView := lipgloss.NewStyle().
		MarginLeft(2).
		MarginRight(2).
		Width((m.width / 2) - 4).
		Render(m.list.View())

	previewView := lipgloss.NewStyle().
		MarginLeft(2).
		MarginRight(2).
		Render(m.preview.View())

	body := lipgloss.JoinHorizontal(0, listView, previewView)
	header := lipgloss.NewStyle().Padding(0, 3).Render(tabView)

	return lipgloss.JoinVertical(0, header, body)
}

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v", err)
	}
}
