package models

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/maximiliansoerenpollak/todotui/styles"
	"github.com/maximiliansoerenpollak/todotui/types"
)

type tasksListModel struct {
	list            list.Model
	parentModel     tea.Model
	parentTaskGroup types.TaskGroup
	choice          string
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
		case "enter":
			i, ok := m.list.SelectedItem().(types.Task)
			if ok {
				m.choice = i.TaskTitle
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
