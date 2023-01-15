package types

type Task struct {
	title  string
	isDone bool
}

type TaskGroup struct {
	title      string
	tasks      []Task
	taskGroups []TaskGroup
}
