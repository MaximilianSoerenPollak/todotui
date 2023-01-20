package models

import (
	"fmt"

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

		case "enter":
			i, ok := m.list.SelectedItem().(types.Task)
			if ok {
				m.choice = i.Title
			}
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m tasksListModel) View() string {
	if m.choice != "" {
		return styles.QuitTextStyle.Render(fmt.Sprintf("Selected: %s", m.choice))
	}
	return "\n" + m.list.View()
}

func TasksListModel(parentTaskGroup types.TaskGroup) tasksListModel {
	return tasksListModel{parentTaskGroup: parentTaskGroup}
}
