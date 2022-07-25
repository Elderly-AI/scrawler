package crawler

import (
	"context"
	"time"

	"github.com/Elderly-AI/scrawler/internal/pkg/model"
)

func (f Facade) GetLessonsCountByTags(ctx context.Context, from, to time.Time, tags []model.Tag) (float64, error) {
	tagsIDs := make([]uint64, len(tags))
	for i, t := range tags {
		tagsIDs[i] = t.ID
	}
	return f.crawlerDB.GetLessonsCountByTags(ctx, tagsIDs, from, to)
}
