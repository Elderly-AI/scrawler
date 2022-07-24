package crawler

import (
	"context"

	desc "github.com/Elderly-AI/scrawler/pkg/proto/crawler"
)

func (i Implementation) CheatExternalData(ctx context.Context, req *desc.CheatExternalDataRequest) (*desc.CheatExternalDataResponse, error) {
	tagsIDs, err := i.facade.CheatExternalTags(ctx)
	if err != nil {
		return nil, err
	}
	for _, tagID := range tagsIDs {
		err = i.facade.CheatExternalLessons(ctx, tagID)
		if err != nil {
			return nil, err
		}
	}
	return &desc.CheatExternalDataResponse{}, nil
}
