package models

import (
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/maximiliansoerenpollak/todotui/db"
	"github.com/maximiliansoerenpollak/todotui/types"
)

type taskGroupsModel struct {
	list      list.Model
	textInput textinput.Model
	state     int
}

func (m taskGroupsModel) Init() tea.Cmd {
	return nil
}

func (m taskGroupsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	case tea.KeyMsg:
		switch m.state {
		case 0:
			switch keypress := msg.String(); keypress {
			case "a":
				m.state = 1
				m.textInput, cmd = m.textInput.Update(msg)
				return m, cmd
			}
		case 1:
			switch keypress := msg.String(); keypress {
			case "ctrl+c":
				return m, tea.Quit
			case "enter":
				if m.state == 1 {
					lis := m.list.Items()
					m.list.SetItems(append(lis, types.TaskGroup{GroupTitle: m.textInput.Value()}))
					m.state = 0
					m.textInput.SetValue("")
				}
			}

		}
	}

	if m.state == 1 {
		m.textInput, cmd = m.textInput.Update(msg)
		return m, cmd
	}

	m.list, cmd = m.list.Update(msg)

	return m, cmd
}

func (m taskGroupsModel) View() string {

	if m.state == 0 {
		return m.list.View()
	}
	return m.textInput.View()

}

func InitiateTaskGroupsList() taskGroupsModel {

	items := []list.Item{}

	for _, j := range db.MemData.TaskGroups {
		items = append(items, j)
	}
	fmt.Println(items)
	ti := textinput.New()
	ti.Placeholder = "Pikachu"
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 20

	return taskGroupsModel{list: list.New(items, list.NewDefaultDelegate(), 0, 0), state: 0, textInput: ti}
}
