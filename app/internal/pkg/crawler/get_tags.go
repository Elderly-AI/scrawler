package crawler

import (
	"context"

	"github.com/Elderly-AI/scrawler/internal/pkg/model"
)

func (f Facade) GetTags(ctx context.Context, search string, limit, offset uint64) ([]model.Tag, uint64, error) {
	dbTags, err := f.crawlerDB.GetTagsWithPagination(ctx, search, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	tags, total := convertDBTags(dbTags)
	return tags, total, nil
}

func convertDBTags(dbTags []model.TagDB) ([]model.Tag, uint64) {
	tags := make([]model.Tag, len(dbTags))
	var total uint64
	for i, t := range dbTags {
		tags[i] = model.Tag{
			ID:         t.ID,
			ExternalID: t.ExternalID,
			Title:      t.Title,
		}
		total = t.Total
	}
	return tags, total
}
