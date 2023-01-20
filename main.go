package main

import (
	// "fmt"
	// "os"

	// tea "github.com/charmbracelet/bubbletea"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/maximiliansoerenpollak/todotui/db"
	"github.com/maximiliansoerenpollak/todotui/models"
	"github.com/maximiliansoerenpollak/todotui/types"
)

func main() {
	db.Write(db.MemData)
	db.MemData.TaskGroups = append(db.MemData.TaskGroups, types.TaskGroup{GroupTitle: "TaskGroup 1", Tasks: []types.Task{}, GroupDescription: "Description 1"})
	m := models.InitiateTaskGroupsList()
	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
	// models.InitTui2()
}
