package crawler

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"

	"github.com/Elderly-AI/scrawler/internal/pkg/model"
)

const solveryUrl = "https://solvery.io/ru/mentors"
const regexpStartToken = "<script id=\"my-app-state\" type=\"application/json\">"
const regexpEndToken = "</script>"

func (f Facade) CheatExternalTags(ctx context.Context) ([]uint64, error) {
	page, err := getHtmlPage(solveryUrl)
	if err != nil {
		return nil, err
	}
	script, err := getPageScript(page)
	if err != nil {
		return nil, err
	}
	externalTags, err := unmarshalTags(script)
	if err != nil {
		return nil, err
	}
	tags := convertExternalTags(externalTags)
	return f.crawlerDB.AddTags(ctx, tags)
}

func getHtmlPage(url string) (page string, err error) {
	resp, err := http.Get(url)
	if err != nil {
		return
	}
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	return string(bytes), nil
}

func getPageScript(page string) (string, error) {
	re := regexp.MustCompile(regexpStartToken + ".*" + regexpEndToken)
	match := re.FindStringSubmatch(page)

	if len(match) != 1 {
		return "", errors.New("tags not found")
	}

	script := match[0]
	jsonString := script[len(regexpStartToken) : len(script)-len(regexpEndToken)]
	jsonString = strings.Replace(jsonString, "&q;", "\"", -1)
	return jsonString, nil
}

func unmarshalTags(raw string) ([]model.ExternalTag, error) {
	var tags model.ExternalTags
	err := json.Unmarshal([]byte(raw), &tags)
	if err != nil {
		return nil, err
	}
	return tags.Tags, nil
}

func convertExternalTags(externalTags []model.ExternalTag) []model.Tag {
	tags := make([]model.Tag, len(externalTags))
	for i, t := range externalTags {
		tags[i] = model.Tag{
			ExternalID: t.ID,
			Title:      t.Title,
		}
	}
	return tags
}
