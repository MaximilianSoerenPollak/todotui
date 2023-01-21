package models

import (

	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/maximiliansoerenpollak/todotui/db"
	"github.com/maximiliansoerenpollak/todotui/styles"
	"github.com/maximiliansoerenpollak/todotui/types"
)

type taskGroupsModel struct {

	list        list.Model
	inputs      []textinput.Model
	state       int // 0 -> List Mode | 1 -> Input Mode
	isFiltering bool
	focusIndex  int
	cursorMode  textinput.CursorMode
}

func (m taskGroupsModel) Init() tea.Cmd {
	return nil
}

func (m taskGroupsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		h, v := styles.DocStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	case tea.KeyMsg:
		switch keyPress := msg.String(); keyPress {
		case "ctrl+c":
			return m, tea.Quit
		}
		switch m.state {
		case 0:
			switch keypress := msg.String(); keypress {
			case "a":
				if !m.isFiltering {
					m.state = 1
					cmd := m.updateInputs(msg)
					return m, cmd
				}
			case "/":
				if len(m.list.Items()) > 0 {
					m.isFiltering = true
				}
			case "esc":
				if m.isFiltering {
					m.isFiltering = false
				}
			}
		case 1:
			switch keypress := msg.String(); keypress {
			case "ctrl+r":
				m.cursorMode++
				if m.cursorMode > textinput.CursorHide {
					m.cursorMode = textinput.CursorBlink
				}
				cmds := make([]tea.Cmd, len(m.inputs))
				for i := range m.inputs {
					cmds[i] = m.inputs[i].SetCursorMode(m.cursorMode)
				}
				return m, tea.Batch(cmds...)

			case "tab", "shift+tab", "enter", "up", "down":
				s := msg.String()

				if s == "enter" && m.focusIndex == len(m.inputs) {
					title := m.inputs[0].Value()
					description := m.inputs[1].Value()
					items := m.list.Items()
					m.state = 0
					m.list.SetItems(append(items, types.TaskGroup{GroupTitle: title, GroupDescription: description}))
					return m, nil
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
						m.inputs[i].PromptStyle = styles.TiFocusedStyle
						m.inputs[i].TextStyle = styles.TiFocusedStyle
						continue
					}
					m.inputs[i].Blur()
					m.inputs[i].PromptStyle = styles.TiNoStyle
					m.inputs[i].TextStyle = styles.TiNoStyle
				}

				return m, tea.Batch(cmds...)
			}

	}
    }
	if m.state == 1 {
		cmd := m.updateInputs(msg)
		return m, cmd
	}


	m.list, cmd = m.list.Update(msg)
    return m, cmd
}


func (m *taskGroupsModel) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))
	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

func (m taskGroupsModel) View() string {
	if m.state == 0 {

		return m.list.View()
	}
	var b strings.Builder

	for i := range m.inputs {
		b.WriteString(m.inputs[i].View())
		if i < len(m.inputs)-1 {
			b.WriteRune('\n')
		}
	}

	button := &styles.TiBlurredButton
	if m.focusIndex == len(m.inputs) {
		button = &styles.TiFocusedButton
	}
	fmt.Fprintf(&b, "\n\n%s\n\n", *button)

	b.WriteString(styles.TiHelpStyle.Render("cursor mode is "))
	b.WriteString(styles.TiCursorModeHelpStyle.Render(m.cursorMode.String()))
	b.WriteString(styles.TiHelpStyle.Render(" (ctrl+r to change style)"))

	return b.String()

}

func InitiateTaskGroupsList() taskGroupsModel {

	items := []list.Item{}
	for _, j := range db.MemData.TaskGroups {
		items = append(items, j)
	}

	m := taskGroupsModel{list: list.New(items, list.NewDefaultDelegate(), 0, 0), state: 0}

	inputs := make([]textinput.Model, 2)
	var t textinput.Model
	for i := range inputs {
		t = textinput.New()
		t.CursorStyle = styles.TiCursorStyle
		t.CharLimit = 32

		switch i {
		case 0:
			t.Placeholder = "Group Title"
			t.Focus()
			t.PromptStyle = styles.TiFocusedStyle
			t.TextStyle = styles.TiFocusedStyle
		case 1:
			t.Placeholder = "Description"
			t.CharLimit = 64
		}

		inputs[i] = t
	}
	m.inputs = inputs
	m.list.SetStatusBarItemName("TaskGroup", "TaskGroups")
	return m
}

// TODOs
// Delete TaskGroups
// Edit TaskGroups
// Cancel Adding New TaskGroup
// Switch To TasksList for the group
// Refactor Update Method
// Update Delegate and add Hints for all keybinds
