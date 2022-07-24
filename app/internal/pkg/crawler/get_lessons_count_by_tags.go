package crawler

import (
	"context"
	"time"

	"github.com/Elderly-AI/scrawler/internal/pkg/model"
)

func (f Facade) GetLessonsCountByTags(ctx context.Context, from, to time.Time, tags []model.Tag) (float64, error) {
	tagsExternalIDs := make([]uint64, len(tags))
	for i, t := range tags {
		tagsExternalIDs[i] = t.ExternalID
	}
	return f.crawlerDB.GetLessonsCountByTags(ctx, tagsExternalIDs, from, to)
}
