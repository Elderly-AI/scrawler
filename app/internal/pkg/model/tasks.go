package model

import (
	"time"
)

type Task struct {
	TaskID    TaskID
	RunPeriod time.Duration
}

type TaskDB struct {
	TaskID    string `db:"task_id"`
	RunPeriod uint64 `db:"run_period"`
}

type TaskLog struct {
	ID     uint64
	TaskID TaskID
	Error  string
	Status TaskLogStatus
}

type TaskLogStatus uint64

const (
	TaskStatusStarted TaskLogStatus = iota
	TaskStatusFailed
	TaskStatusFinished
)

type TaskID string

const (
	TaskIDCheat TaskID = "cheat_task"
)

type TaskLogDB struct {
	ID        uint64    `db:"task_log_id"`
	TaskID    string    `db:"task_id"`
	Status    uint64    `db:"status"`
	Error     string    `db:"error"`
	CreatedAt time.Time `db:"created_at"`
}
