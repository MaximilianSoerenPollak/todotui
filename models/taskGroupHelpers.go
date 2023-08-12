package models

import (
	"strconv"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/google/uuid"
	"github.com/maximiliansoerenpollak/todotui/styles"
	"github.com/maximiliansoerenpollak/todotui/types"
)

func taskGroupViewState0(m taskGroupsModel, msg tea.KeyMsg) (tea.Model, tea.Cmd, bool) {

	switch keypress := msg.String(); keypress {
	case "a":
		if !m.isFiltering {
			m.state = 1
			cmd := m.updateInputs(msg)
			m.inputs[0].SetValue("")
			m.inputs[1].SetValue("")
			m.inputs[0].Focus()
			m.focusIndex = 0
			return m, cmd, false
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
		return m, nil, false

	case "e":
		if !m.isFiltering {
			selected, ok := m.list.SelectedItem().(types.TaskGroup)
			if ok {
				m.selected = selected
				m.state = 2
			}
			cmd := m.updateInputs(msg)
			m.editInputs[0].SetValue(selected.GroupTitle)
			m.editInputs[1].SetValue(selected.GroupDescription)
			m.editInputs[0].Focus()
			m.editFocusIndex = 0
			return m, cmd, false
		}
	case "enter":
		parentTaskGroup, ok := m.list.SelectedItem().(types.TaskGroup)
		if ok {
			return initiateTasksListModel(parentTaskGroup, m), nil, false
		}
	}
	return m, nil, true

}

func taskGroupViewState1(m taskGroupsModel, msg tea.KeyMsg) (tea.Model, tea.Cmd, bool) {
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
		return m, tea.Batch(cmds...), false

	case "tab", "shift+tab", "enter", "up", "down":
		s := msg.String()

		if s == "enter" && m.focusIndex == len(m.inputs) {
			title := m.inputs[0].Value()
			description := m.inputs[1].Value()
			items := m.list.Items()
			m.state = 0
			m.list.SetItems(append(items, types.TaskGroup{GroupTitle: title, GroupDescription: description, GroupId: uuid.NewString()}))
			return m, nil, false
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

		return m, tea.Batch(cmds...), false
	}
	return m, nil, true

}

func taskGroupViewState2(m taskGroupsModel, msg tea.KeyMsg) (tea.Model, tea.Cmd, bool) {
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
		return m, tea.Batch(cmds...), false

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
			return m, nil, false
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

		return m, tea.Batch(cmds...), false
	}
	return m, nil, true
}
// -------------- TASKS VIEW STATES -------------------------

func taskListViewState0(m tasksListModel, msg tea.KeyMsg) (tea.Model, tea.Cmd, bool) {
	switch keypress := msg.String(); keypress {
	case "a":
		if !m.isFiltering {
			m.state = 1
			cmd := m.updateInputs(msg)
			m.inputs[0].SetValue("") // out of Range error
			m.inputs[0].Focus()
			m.focusIndex = 0
			switch keypress := msg.String(); keypress {
			case "enter":
				m.done = !m.done
			}
			return m, cmd, false
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
		return m, nil, false

	case "e":
		if !m.isFiltering {
			selected, ok := m.list.SelectedItem().(types.Task)
			if ok {
				m.selected = selected
				m.state = 2
			}
			cmd := m.updateInputs(msg)
			m.editInputs[0].SetValue(selected.Title())
			newValue := m.convertDone(selected.IsDone)
			m.editInputs[1].SetValue(newValue)
			m.editInputs[0].Focus()
			m.editFocusIndex = 0
			return m, cmd, false
		}
	case "enter":
		parentTaskGroup, ok := m.list.SelectedItem().(types.TaskGroup)
		if ok {
			return initiateTasksListModel(parentTaskGroup, m), nil, false
		}
	}
	return m, nil, true

}

func (m tasksListModel) convertDone(b bool) string {
	return strconv.FormatBool(b)	
}
