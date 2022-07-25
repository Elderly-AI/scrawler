package crawler

import (
	"context"
	"time"

	"github.com/Elderly-AI/scrawler/internal/pkg/model"
	desc "github.com/Elderly-AI/scrawler/pkg/proto/crawler"
)

func (i Implementation) GetLessonsCountByTags(ctx context.Context, req *desc.GetLessonsCountByTagsRequest) (*desc.GetLessonsCountByTagsResponse, error) {
	from := time.Now().AddDate(0, 0, -1)
	if req.GetFrom() != nil {
		from = req.GetFrom().AsTime()
	}
	to := time.Now()
	if req.GetTo() != nil {
		to = req.GetTo().AsTime()
	}
	tags := make([]model.Tag, len(req.GetTagIds()))
	for j, t := range req.TagIds {
		tags[j] = model.Tag{ID: t}
	}
	count, err := i.facade.GetLessonsCountByTags(ctx, from, to, tags)
	if err != nil {
		return nil, err
	}
	return &desc.GetLessonsCountByTagsResponse{
		Count: count,
	}, nil
}
