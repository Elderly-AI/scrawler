package model

type Task struct {
	ID     uint64     `json:"task_log_id"`
	TaskID string     `json:"task_id"`
	Error  string     `json:"error"`
	Status TaskStatus `json:"status"`
}

type TaskStatus uint64

const (
	TaskStatusStarted TaskStatus = iota
	TaskStatusFailed
	TaskStatusFinished
)

type TaskDB struct {
	ID     uint64 `db:"task_log_id"`
	TaskID string `db:"task_id"`
	Status uint64 `db:"status"`
	Error  string `db:"error"`
}
