package models

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/maximiliansoerenpollak/todotui/styles"
	"github.com/maximiliansoerenpollak/todotui/types"
	"github.com/charmbracelet/bubbles/textinput"
)

type tasksListModel struct {
	list           list.Model
	inputs         []textinput.Model
	state          int // 0 -> List Mode | 1 -> Input Mode | 2 -> Editing Mode
	editInputs     []textinput.Model
	parentModel   tea.Model
	isFiltering    bool
	focusIndex     int
	editFocusIndex int
	cursorMode     textinput.CursorMode
	selected       types.TaskGroup

}

func (m tasksListModel) Init() tea.Cmd {
	return nil
}

func (m tasksListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)
		return m, nil
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "ctrl+c":
			return m, tea.Quit
		case "home":
			return m.parentModel, nil
		case "a":
			m.updateInputs()
		case "enter":
			i, ok := m.list.SelectedItem().(types.Task)
			if ok {
				m.choice = i.TaskTitle
				m.updateInputs(msg)
			}
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m tasksListModel) View() string {
	return styles.DocStyle.Render(m.list.View())
}

func initiateTasksListModel(parentTaskGroup types.TaskGroup, parentModel tea.Model) tasksListModel {
	m := tasksListModel{parentTaskGroup: parentTaskGroup, parentModel: parentModel}
	items := []list.Item{}
	for _, j := range parentTaskGroup.Tasks {
		items = append(items, j)
	}
	m.list = list.New(items, list.NewDefaultDelegate(), 0, 0)
	return m
}

func (m *tasksListModel) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))
	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

func (m *tasksListModel) updateEditInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.editInputs))
	for i := range m.editInputs {
		m.editInputs[i], cmds[i] = m.editInputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}
