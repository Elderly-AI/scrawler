package tasks

import (
	"context"
	"github.com/jmoiron/sqlx"

	"github.com/Elderly-AI/scrawler/internal/pkg/model"
)

type Database struct {
	conn *sqlx.DB
}

func New(conn *sqlx.DB) Database {
	return Database{
		conn,
	}
}

func (d Database) GetTask(ctx context.Context, taskID model.TaskID) (taskDB model.TaskDB, err error) {
	var tasksDB []model.TaskDB
	query := `SELECT task_id, run_period  FROM tasks
			  WHERE deleted_at IS NULL AND task_id = $1 LIMIT 1`
	err = d.conn.Select(&tasksDB, query, string(taskID))
	if err != nil {
		return model.TaskDB{}, err
	}
	if len(tasksDB) == 1 {
		return tasksDB[0], nil
	}
	return model.TaskDB{}, nil
}

func (d Database) CreateTaskLog(ctx context.Context, task model.TaskLog) (taskLogDB model.TaskLogDB, err error) {
	var taskLogsDB []model.TaskLogDB
	query := `INSERT INTO tasks_logs (task_id) VALUES ($1) RETURNING task_log_id, task_id, status, error, created_at`
	err = d.conn.Select(&taskLogsDB, query, task.TaskID)
	if err != nil {
		return model.TaskLogDB{}, err
	}
	if len(taskLogsDB) == 1 {
		return taskLogsDB[0], nil
	}
	return model.TaskLogDB{}, nil
}

func (d Database) UpdateTaskLog(ctx context.Context, task model.TaskLog) (taskLogDB model.TaskLogDB, err error) {
	var taskLogsDB []model.TaskLogDB
	query := `UPDATE tasks_logs SET status = $2, error = $3 WHERE task_id = $1 RETURNING task_log_id, task_id, status, error, created_at`
	err = d.conn.Select(&taskLogsDB, query, task.TaskID, uint64(task.Status), task.Error)
	if err != nil {
		return model.TaskLogDB{}, err
	}
	if len(taskLogsDB) == 1 {
		return taskLogsDB[0], nil
	}
	return model.TaskLogDB{}, nil
}

func (d Database) GetLatestFinishedTask(ctx context.Context, task model.Task) (taskLogsDB []model.TaskLogDB, err error) {
	query := `SELECT task_log_id, task_id, status, error, created_at FROM tasks_logs 
			  WHERE deleted_at IS NULL AND status = task_status_finished() AND task_id = $1 
			  ORDER BY created_at ASC LIMIT 1`
	err = d.conn.Select(&taskLogsDB, query, string(task.TaskID))
	return
}
