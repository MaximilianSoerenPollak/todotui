package types

type Task struct {
	Title  string
	IsDone bool
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

func (t Task) FilterValue() string { return t.Title }
