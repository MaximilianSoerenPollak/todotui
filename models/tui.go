package models

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/google/uuid"
	"github.com/maximiliansoerenpollak/todotui/db"
	"github.com/maximiliansoerenpollak/todotui/styles"
	"github.com/maximiliansoerenpollak/todotui/types"
)

type taskGroupsModel struct {
	list           list.Model
	inputs         []textinput.Model
	state          int // 0 -> List Mode | 1 -> Input Mode | 2 -> Editing Mode
	editInputs     []textinput.Model
	isFiltering    bool
	focusIndex     int
	editFocusIndex int
	cursorMode     textinput.CursorMode
	selected       types.TaskGroup
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
					m.inputs[0].SetValue("")
					m.inputs[1].SetValue("")
					m.inputs[0].Focus()
					m.focusIndex = 0
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
			case "d":
				index := m.list.Index()
				m.list.RemoveItem(index)
				return m, nil

			case "e":
				if !m.isFiltering {
					selected, ok := m.list.SelectedItem().(types.TaskGroup)
					if ok {
						m.selected = selected
						m.state = 2
					}
					cmd := m.updateInputs(msg)
					m.editInputs[0].SetValue("")
					m.editInputs[1].SetValue("")
					m.editInputs[0].Focus()
					m.editFocusIndex = 0
					return m, cmd
				}
			case "enter":
				parentTaskGroup, ok := m.list.SelectedItem().(types.TaskGroup)
				if ok {
					return initiateTasksListModel(parentTaskGroup, m), nil
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
					m.list.SetItems(append(items, types.TaskGroup{GroupTitle: title, GroupDescription: description, GroupId: uuid.NewString()}))
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
		case 2:
			switch keypress := msg.String(); keypress {
			case "ctrl+r":
				m.cursorMode++
				if m.cursorMode > textinput.CursorHide {
					m.cursorMode = textinput.CursorBlink
				}
				cmds := make([]tea.Cmd, len(m.editInputs))
				for i := range m.editInputs {
					cmds[i] = m.editInputs[i].SetCursorMode(m.cursorMode)
				}
				return m, tea.Batch(cmds...)

			case "tab", "shift+tab", "enter", "up", "down":
				s := msg.String()

				if s == "enter" && m.editFocusIndex == len(m.editInputs) {

					title := m.editInputs[0].Value()
					description := m.editInputs[1].Value()
					m.state = 0
					items := m.list.Items()
					for i, j := range items {
						l := j.(types.TaskGroup)
						if l.GroupId == m.selected.GroupId {
							l.GroupTitle = title
							l.GroupDescription = description
							items[i] = l
							break
						}

					}
					m.list.SetItems(items)
					return m, nil
				}
				if s == "up" || s == "shift+tab" {
					m.editFocusIndex--
				} else {
					m.editFocusIndex++
				}

				if m.editFocusIndex > len(m.editInputs) {
					m.editFocusIndex = 0
				} else if m.editFocusIndex < 0 {
					m.editFocusIndex = len(m.editInputs)
				}

				cmds := make([]tea.Cmd, len(m.editInputs))
				for i := 0; i <= len(m.editInputs)-1; i++ {
					if i == m.editFocusIndex {
						cmds[i] = m.editInputs[i].Focus()
						m.editInputs[i].PromptStyle = styles.TiFocusedStyle
						m.editInputs[i].TextStyle = styles.TiFocusedStyle
						continue
					}
					m.editInputs[i].Blur()
					m.editInputs[i].PromptStyle = styles.TiNoStyle
					m.editInputs[i].TextStyle = styles.TiNoStyle
				}

				return m, tea.Batch(cmds...)
			}

		}
	}
	if m.state == 1 {
		cmd := m.updateInputs(msg)
		return m, cmd
	}

	if m.state == 2 {
		cmd := m.updateEditInputs(msg)
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

func (m *taskGroupsModel) updateEditInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.editInputs))
	for i := range m.editInputs {
		m.editInputs[i], cmds[i] = m.editInputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

func (m taskGroupsModel) View() string {
	if m.state == 0 {
		return m.list.View()
	}

	if m.state == 2 {
		var b strings.Builder

		for i := range m.editInputs {
			b.WriteString(m.editInputs[i].View())
			if i < len(m.editInputs)-1 {
				b.WriteRune('\n')
			}
		}

		button := &styles.TiBlurredButton
		if m.focusIndex == len(m.editInputs) {
			button = &styles.TiFocusedButton
		}
		fmt.Fprintf(&b, "\n\n%s\n\n", *button)

		b.WriteString(styles.TiHelpStyle.Render("cursor mode is "))
		b.WriteString(styles.TiCursorModeHelpStyle.Render(m.cursorMode.String()))
		b.WriteString(styles.TiHelpStyle.Render(" (ctrl+r to change style)"))

		return b.String()
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
	editInputs := make([]textinput.Model, 2)
	for i := range editInputs {
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
		editInputs[i] = t
	}

	m.editInputs = editInputs
	m.inputs = inputs

	m.list.SetStatusBarItemName("TaskGroup", "TaskGroups")
	return m
}

// TODOs
// [x] Delete TaskGroups
// [x] Edit TaskGroups
// [] Cancel Adding New TaskGroup
// [x] Switch To TasksList for the group
// [] Refactor Update Method
// [] Update Delegate and add Hints for all keybinds
// [] Add appropriate titles for inputs windows
