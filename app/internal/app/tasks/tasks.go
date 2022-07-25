package tasks

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/Elderly-AI/scrawler/internal/pkg/model"
	desc "github.com/Elderly-AI/scrawler/pkg/proto/crawler"
)

type Facade interface {
	RunTask(ctx context.Context, taskID model.TaskID, runner func(ctx context.Context) error) error
}

type CrawlerFacade interface {
	CheatExternalData(ctx context.Context, req *desc.CheatExternalDataRequest) (*desc.CheatExternalDataResponse, error)
}

type Implementation struct {
	facade        Facade
	crawlerFacade CrawlerFacade
}

func New(facade Facade, crawlerFacade CrawlerFacade) Implementation {
	return Implementation{
		facade:        facade,
		crawlerFacade: crawlerFacade,
	}
}

func (i *Implementation) RunTasks(ctx context.Context, tasks []model.TaskID, interval time.Duration) {
	for {
		for _, task := range tasks {
			time.Sleep(time.Duration(rand.Intn(int(interval.Nanoseconds()))))
			runner := i.getRunner(task)
			err := i.facade.RunTask(ctx, task, runner)
			if err != nil {
				fmt.Println("Error in task runner: ", err.Error())
			}
		}
	}
}

func (i *Implementation) getRunner(task model.TaskID) func(ctx context.Context) error {
	switch task {
	case model.TaskIDCheat:
		return func(ctx context.Context) error {
			_, err := i.crawlerFacade.CheatExternalData(ctx, &desc.CheatExternalDataRequest{})
			return err
		}
	default:
		return func(ctx context.Context) error {
			return nil
		}
	}
}
