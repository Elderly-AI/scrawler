package tasks

import (
	"context"
	"fmt"
	"time"

	"github.com/Elderly-AI/scrawler/internal/pkg/model"
)

type tasksDB interface {
	GetTask(ctx context.Context, taskID model.TaskID) (taskDB model.TaskDB, err error)
	CreateTaskLog(ctx context.Context, task model.TaskLog) (tasLogDB model.TaskLogDB, err error)
	UpdateTaskLog(ctx context.Context, task model.TaskLog) (tasLogDB model.TaskLogDB, err error)
	GetLatestFinishedTask(ctx context.Context, task model.Task) (taskLogsDB []model.TaskLogDB, err error)
}

type Facade struct {
	tasksDB tasksDB
}

func New(crawlerDB tasksDB) Facade {
	return Facade{
		tasksDB: crawlerDB,
	}
}

func (f Facade) RunTask(ctx context.Context, taskID model.TaskID, runner func(ctx context.Context) error) error {
	fmt.Println("run task", taskID, time.Now())
	taskDB, err := f.tasksDB.GetTask(ctx, taskID)
	task := convertTask(taskDB)
	if err != nil {
		return err
	}
	taskLogs, err := f.tasksDB.GetLatestFinishedTask(ctx, task)
	if len(taskLogs) > 0 && taskLogs[0].CreatedAt.Add(task.RunPeriod).After(time.Now()) {
		return nil
	}
	logDB, err := f.tasksDB.CreateTaskLog(ctx, model.TaskLog{
		TaskID: taskID,
	})
	if err != nil {
		return err
	}
	log := convertTaskLog(logDB)
	taskErr := runner(ctx)

	defer func() {
		if r := recover(); r != nil {
			taskErr = fmt.Errorf("panic in runner %s", r)
		}
	}()

	log.Status = model.TaskStatusFinished
	if taskErr != nil {
		log.Status = model.TaskStatusFailed
		log.Error = taskErr.Error()
	}
	_, err = f.tasksDB.UpdateTaskLog(ctx, log)
	return err
}

func convertTask(taskDB model.TaskDB) model.Task {
	return model.Task{
		TaskID:    model.TaskID(taskDB.TaskID),
		RunPeriod: time.Duration(taskDB.RunPeriod),
	}
}

func convertTaskLog(taskLogDB model.TaskLogDB) model.TaskLog {
	return model.TaskLog{
		ID:     taskLogDB.ID,
		TaskID: model.TaskID(taskLogDB.TaskID),
		Error:  taskLogDB.Error,
		Status: model.TaskLogStatus(taskLogDB.Status),
	}
}
