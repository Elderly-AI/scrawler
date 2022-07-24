package crawler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/Elderly-AI/scrawler/internal/pkg/model"
)

const solveryGetLessonsUrl = "https://solvery.io/api/mentoring/listing/get?locale=ru"
const solveryGetLessonData = `{"filters":{"visibility":"ALL_USERS","rate":{},"discounts":[],"isAvailableOnly":false,"tagIds":[%d],"language":["ru"]},"pagination":{"limit":10,"page":%d},"sort":{"rate":true}}`
const defaultPageSize = 10

func (f Facade) CheatExternalLessons(ctx context.Context, tagID uint64) error {
	entities, err := processBatchLessonsRequests(tagID)
	if err != nil {
		return err
	}
	logs := convertEntities(entities)
	return f.crawlerDB.AddLogsWithTags(ctx, logs)
}

func convertEntities(entities []model.ExternalEntity) []model.LessonLog {
	logs := make([]model.LessonLog, len(entities))
	for i, e := range entities {
		logs[i] = model.LessonLog{
			ExternalID:   e.ID,
			LessonsCount: e.MentorProfile.Statistics.Sessions.Count,
			Tags:         convertExternalTags(e.MentorProfile.Tags),
		}
	}
	return logs
}

func processBatchLessonsRequests(tagID uint64) ([]model.ExternalEntity, error) {
	var entities []model.ExternalEntity
	page := uint64(0)
	total := uint64(0)
	for page*defaultPageSize <= total {
		page += 1
		resp, err := processLessonsRequest(page, tagID)
		if err != nil {
			return nil, err
		}
		total = resp.Result.Count
		entities = append(entities, resp.Result.Entities...)
	}
	return entities, nil
}

func processLessonsRequest(page uint64, tagID uint64) (*model.ExternalLessonResponse, error) {
	var jsonStr = []byte(fmt.Sprintf(solveryGetLessonData, tagID, page))
	req, err := http.NewRequest("POST", solveryGetLessonsUrl, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	var response model.ExternalLessonResponse
	err = json.Unmarshal(body, &response)
	return &response, err
}
