package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const (
	baseUrl         string = "https://hacker-news.firebaseio.com/v0/"
	MaxStoriesLimit int    = 500
)

type HnClient struct {
	client http.Client
	url    string
}

func MakeClient() HnClient {
	return MakeClientForUrl(baseUrl)
}

func MakeClientForUrl(url string) HnClient {
	return HnClient{
		http.Client{},
		url,
	}
}

func (hn *HnClient) FetchRankedStoriesIds(ranking FrontPageItemsRanking, limit int) ([]ItemId, error) {
	if limit < 0 || limit > MaxStoriesLimit {
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
	response, err := hn.client.Get(fmt.Sprintf("%s/%s.json", hn.url, endpoint))
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
		return nil, err
	}

	if len(ids) <= limit {
		return ids, nil
	}
	return ids[:limit], nil
}

func (hn *HnClient) FetchItem(id ItemId) (*Item, error) {
	response, err := hn.client.Get(fmt.Sprintf("%s/item/%d.json", hn.url, id))
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

	var item Item
	if err := json.Unmarshal(body, &item); err != nil {
		return nil, err
	}

	return &item, nil
}
