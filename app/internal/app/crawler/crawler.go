package crawler

import (
	"context"
	"time"

	"github.com/Elderly-AI/scrawler/internal/pkg/model"
	desc "github.com/Elderly-AI/scrawler/pkg/proto/crawler"
)

type Facade interface {
	GetTags(ctx context.Context, search string, limit, offset uint64) ([]model.Tag, uint64, error)
	GetLessonsCountByTags(ctx context.Context, from, to time.Time, tags []model.Tag) (float64, error)

	CheatExternalTags(ctx context.Context) ([]uint64, error)
	CheatExternalLessons(ctx context.Context, tagID uint64) error
}

type Implementation struct {
	facade Facade
	desc.UnimplementedCrawlerServer
}

func New(facade Facade) Implementation {
	return Implementation{
		facade: facade,
	}
}
