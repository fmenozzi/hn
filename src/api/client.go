package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const (
	baseUrl         string = "https://hacker-news.firebaseio.com/v0/"
	maxStoriesLimit int    = 500
)

type HnClient struct {
	client http.Client
}

func MakeClient() HnClient {
	return HnClient{
		http.Client{},
	}
}

func (hn *HnClient) FetchStories(ranking StoriesRanking, limit int) (*Stories, error) {
	if limit < 0 || limit > maxStoriesLimit {
		return nil, fmt.Errorf("Invalid limit: %d\n", limit)
	}
	var endpoint string
	switch ranking {
	case Top:
		endpoint = "topstories"
	case Best:
		endpoint = "beststories"
	case New:
		endpoint = "newstories"
	}
	response, err := hn.client.Get(fmt.Sprintf("%s/%s.json", baseUrl, endpoint))
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	if response.StatusCode > 299 {
		return nil, fmt.Errorf("Response failed with code %d\n", response.StatusCode)
	}
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	var ids []ItemId
	if err := json.Unmarshal(body, &ids); err != nil {
		return nil, fmt.Errorf("Failed to parse json body: %s\n", err)
	}
	return &Stories{
		Ids:     ids[:limit],
		Ranking: ranking,
	}, nil
}

func (hn *HnClient) FetchItem(id ItemId) (*Item, error) {
	response, err := hn.client.Get(fmt.Sprintf("%s/item/%d.json", baseUrl, id))
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	var item Item
	if err := json.Unmarshal(body, &item); err != nil {
		return nil, err
	}
	return &item, nil
}
