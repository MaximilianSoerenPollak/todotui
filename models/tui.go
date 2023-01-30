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
			m, cmd, check := taskGroupViewState0(m, msg)
			if !check {
				return m, cmd
			}
		case 1:
			m, cmd, check := taskGroupViewState1(m, msg)
			if !check {
				return m, cmd
			}
		case 2:
			m, cmd, check := taskGroupViewState2(m, msg)
			if !check {
				return m, cmd
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
// [x] Refactor Update Method
// [] Update Delegate and add Hints for all keybinds
// [] Add appropriate titles for inputs windows
