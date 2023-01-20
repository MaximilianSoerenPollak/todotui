package models

import (
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/maximiliansoerenpollak/todotui/db"
	"github.com/maximiliansoerenpollak/todotui/types"
)

type taskGroupsModel struct {
	list       list.Model
	textInputs []textinput.Model
	state      int
    cursorMode textinput.CursorMode
}

func (m taskGroupsModel) Init() tea.Cmd {
	return nil
}

func (m taskGroupsModel) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.textInputs))

	for i := range m.textInputs {
		m.textInputs[i], cmds[i] = m.textInputs[i].Update(msg)
	}
	return tea.Batch(cmds...)
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
				m.textInputs, cmd = m.textInputs.updateInputs(msg)
				return m, cmd
			}
		case 1:
			switch keypress := msg.String(); keypress {
            // Change cursorMode
            case "ctrl+r":
                m.cursorMode++
                if m.cursorMode > textinput.CursorHide {
                    m.cursorMode = textinput.CursorBlink
                }
                cmds := make ([]tea.Cmd, len(m.textInputs))
                for i := range m.textInputs {
                    cmds[i] = m.textInputs[i].SetCursorMode(m.cursorMode)
                }
                return m, tea.Batch(cmds...)
            // set cofucs to the next tab
            case "tab", "shit+tab", "enter", "up", "down":
                s :=msg.String()
                if s == "enter" && m.focusIndex == len(m.textInputs) {
                    return m, tea.Quit
                }
			case "ctrl+c":
				return m, tea.Quit
			// case "enter":
			// 	if m.state == 1 {
			// 		lis := m.list.Items()
			// 		m.list.SetItems(append(lis, types.TaskGroup{GroupTitle: m.textInputs.Value()}))
			// 		m.state = 0
			// 		m.textInputs.SetValue("")
				}
			}

		}
        return m, cmd
	}

	// if m.state == 1 {
	// 	m.textInputs, cmd = m.textInputs.Update(msg)
	// 	return m, cmd
	// }

	// m.list, cmd = m.list.Update(msg)



func (m taskGroupsModel) View() string {
	if m.state == 0 {
		return m.list.View()
	}

	var b strings.Builder

	for i := range m.textInputs {
		b.WriteString(m.textInputs[i].View())
		if i < len(m.textInputs)-1 {
			b.WriteRune('\n')
		}
	}
	return b.String()

	// return m.textInputs.View()

}

func InitiateTaskGroupsList() taskGroupsModel {

	items := []list.Item{}
	for _, j := range db.MemData.TaskGroups {
		items = append(items, j)
	}
	ti := textinput.New()
	ti.Placeholder = "Pikachu"
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 20
	m := taskGroupsModel{
		list:       list.New(items, list.NewDefaultDelegate(), 0, 0),
		textInputs: make([]textinput.Model, 2),
		state:      0,
	}
	var t textinput.Model
	for i := range m.textInputs {
		t = textinput.New()

		switch i {
		case 0:
			t.Placeholder = "TaskGroupName"
			t.Focus()
		case 1:
			t.Placeholder = "Description"
			t.Focus()
		}
		m.textInputs[i] = t
	}

	return m
}
