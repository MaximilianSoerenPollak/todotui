package models

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/maximiliansoerenpollak/todotui/styles"
	"github.com/maximiliansoerenpollak/todotui/types"
	"log"
)

// What to do:
// 1. Have a model that CAN include a LIST of Tasks
// ? Just do []model.Task not sure if this will give us the wanted thing
// 2. Could then be able to just selecte the element we want to add? may make it much easier

type tasksListModel struct {
	list            list.Model
	inputs          []textinput.Model
	state           int // 0 -> List Mode | 1 -> Input Mode | 2 -> Editing Mode
	editInputs      []textinput.Model
	parentModel     tea.Model
	isFiltering     bool
	focusIndex      int
	editFocusIndex  int
	choice          string
	done            bool
	cursorMode      textinput.CursorMode
	selected        types.Task
	parentTaskGroup types.TaskGroup
}

func (m tasksListModel) Init() tea.Cmd {
	return nil
}

func (m tasksListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)
		return m, nil
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "ctrl+c":
			return m, tea.Quit
		}
		switch m.state {
		case 0:
			m, cmd, check := taskListViewState0(m, msg)
			if !check {
				return m, cmd
			}
			// case 1:
			// 	m, cmd, check := taskListViewState1(m, msg)
			// 	if !check {
			// 		return m, cmd
			// 	}
			// case 2:
			// 	m, cmd, check := taskListViewState2(m, msg)
			// 	if !check {
			// 		return m, cmd
			// 	}
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
	// // Go back, one layer up
	// case "backspace":
	// 	return m.parentModel, nil
	// case "a":
	// 	m.updateInputs(msg)
	// case "enter":
	// 	log.Println("Hello, I'm in the tasklistModel update in the enter function")
	// 	log.Printf("These are the input's here: msg %s", msg)
	// 	i, ok := m.list.SelectedItem().(types.Task)
	// 	log.Printf("What is the selecteditem: %T", m.list.SelectedItem())
	// 	log.Printf("Was it ok? ok: %v ", ok)
	// 	log.Printf("What is i, %T", i)
	// 	if ok {
	// 		m.choice = i.TaskTitle
	// 		m.updateInputs(msg)
	// 	}
	// 	log.Println("Selected item was not okay, doing nothing")
	// 	return m, nil
	// }
}

func (m tasksListModel) View() string {
	return styles.DocStyle.Render(m.list.View())
}
func initiateTasksListModel(parentTaskGroup types.TaskGroup, parentModel tea.Model) tasksListModel {
	items := []list.Item{}
	for _, j := range parentTaskGroup.Tasks {
		items = append(items, j)
	}
	m := tasksListModel{list: list.New(items, list.NewDefaultDelegate(), 0, 0), parentTaskGroup: parentTaskGroup, parentModel: parentModel}
	log.Println(m.list.Items())
	return m
}

// func initiateTasksListModel(parentTaskGroup types.TaskGroup, parentModel tea.Model) tasksListModel {
// 	m := tasksListModel{parentTaskGroup: parentTaskGroup, parentModel: parentModel}
// 	items := []list.Item{}
// 	for _, j := range parentTaskGroup.Tasks {
// 		items = append(items, j)
// 	}
// 	m.list = list.New(items, list.NewDefaultDelegate(), 0, 0)
//
// }

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
