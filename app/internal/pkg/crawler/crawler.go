package crawler

import (
	"context"
	"time"

	"github.com/Elderly-AI/scrawler/internal/pkg/model"
)

type crawlerDB interface {
	AddTags(ctx context.Context, tags []model.Tag) (ids []uint64, err error)
	GetTagsWithPagination(ctx context.Context, search string, limit, offset uint64) (tags []model.TagDB, err error)
	AddLogsWithTags(ctx context.Context, logs []model.LessonLog) (err error)
	GetLessonsCountByTags(ctx context.Context, tagsExternalIDs []uint64, from, to time.Time) (count float64, err error)
}

type Facade struct {
	crawlerDB crawlerDB
}

func New(crawlerDB crawlerDB) Facade {
	return Facade{
		crawlerDB: crawlerDB,
	}
}
