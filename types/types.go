package types

import (
	"strconv"
)

type Task struct {
	TaskTitle string
	IsDone    bool
}

type TaskGroup struct {
	GroupTitle       string
	GroupDescription string
	Tasks            []Task
	TaskGroups       []TaskGroup
}

type DbDataT struct {
	TaskGroups []TaskGroup
}

func (g TaskGroup) FilterValue() string { return g.GroupTitle }
func (g TaskGroup) Title() string       { return g.GroupTitle }
func (g TaskGroup) Description() string { return g.GroupDescription }

func (t Task) FilterValue() string { return t.TaskTitle }
func (g Task) Title() string       { return g.TaskTitle }
func (g Task) Description() string { return strconv.FormatBool(g.IsDone) }
