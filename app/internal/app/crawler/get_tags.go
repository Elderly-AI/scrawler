package crawler

import (
	"context"

	"github.com/Elderly-AI/scrawler/internal/pkg/model"
	desc "github.com/Elderly-AI/scrawler/pkg/proto/crawler"
)

const defaultLimit = 20

func (i Implementation) GetTags(ctx context.Context, req *desc.GetTagsRequest) (*desc.GetTagsResponse, error) {
	limit := req.GetPageSize()
	if limit == 0 {
		limit = defaultLimit
	}
	tags, total, err := i.facade.GetTags(ctx, req.GetSearch(), limit, req.GetPage()*limit)
	if err != nil {
		return nil, err
	}
	return &desc.GetTagsResponse{
		Total: total,
		Tags:  convertResponseTags(tags),
	}, nil
}

func convertResponseTags(tags []model.Tag) []*desc.GetTagsResponse_Tag {
	converted := make([]*desc.GetTagsResponse_Tag, len(tags))
	for i, t := range tags {
		converted[i] = &desc.GetTagsResponse_Tag{
			Id:         t.ID,
			Title:      t.Title,
			ExternalId: t.ExternalID,
		}
	}
	return converted
}
