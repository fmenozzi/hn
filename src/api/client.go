package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"sync"
)

const (
	prodHnUrl               string = "https://hacker-news.firebaseio.com/v0/"
	prodSearchPopularityUrl string = "http://hn.algolia.com/api/v1/search"
	prodSearchDateUrl       string = "http://hn.algolia.com/api/v1/search_by_date"
	maxStoriesLimit         int    = 500
)

type HnClient struct {
	client              http.Client
	hnUrl               string
	searchPopularityUrl string
	searchDateUrl       string
}

type HnClientBuilder interface {
	SetHnUrl(string) HnClientBuilder
	SetSearchPopularityUrl(string) HnClientBuilder
	SetSearchDateUrl(string) HnClientBuilder
	Build() HnClient
}

type concreteHnClientBuilder struct {
	hnclient HnClient
}

func (b *concreteHnClientBuilder) SetHnUrl(url string) HnClientBuilder {
	b.hnclient.hnUrl = url
	return b
}

func (b *concreteHnClientBuilder) SetSearchPopularityUrl(url string) HnClientBuilder {
	b.hnclient.searchPopularityUrl = url
	return b
}

func (b *concreteHnClientBuilder) SetSearchDateUrl(url string) HnClientBuilder {
	b.hnclient.searchDateUrl = url
	return b
}

func (b *concreteHnClientBuilder) Build() HnClient {
	return b.hnclient
}

func NewHnClientBuilder() HnClientBuilder {
	return &concreteHnClientBuilder{}
}

func MakeProdClient() HnClient {
	return NewHnClientBuilder().
		SetHnUrl(prodHnUrl).
		SetSearchPopularityUrl(prodSearchPopularityUrl).
		SetSearchDateUrl(prodSearchDateUrl).
		Build()
}

func (hn *HnClient) FetchFrontPageItemIds(ranking FrontPageItemsRanking, limit int) ([]ItemId, error) {
	if limit < 0 || limit > maxStoriesLimit {
		return nil, fmt.Errorf("invalid limit: %d\n", limit)
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
	response, err := hn.client.Get(fmt.Sprintf("%s/%s.json", hn.hnUrl, endpoint))
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	if response.StatusCode > 299 {
		return nil, fmt.Errorf("response failed with code %d\n", response.StatusCode)
	}

	var ids []ItemId
	if err := json.NewDecoder(response.Body).Decode(&ids); err != nil {
		return nil, err
	}

	if len(ids) <= limit {
		return ids, nil
	}
	return ids[:limit], nil
}

func (hn *HnClient) FetchItem(id ItemId) (*Item, error) {
	response, err := hn.client.Get(fmt.Sprintf("%s/item/%d.json", hn.hnUrl, id))
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	if response.StatusCode > 299 {
		return nil, fmt.Errorf("response failed with code %d\n", response.StatusCode)
	}

	var item Item
	if err := json.NewDecoder(response.Body).Decode(&item); err != nil {
		return nil, err
	}

	return &item, nil
}

func (hn *HnClient) FetchItems(ids []ItemId) ([]Item, error) {
	itemsMap := sync.Map{}
	wg := sync.WaitGroup{}
	errchan := make(chan error, len(ids)) // Buffered for non-blocking
	for _, id := range ids {
		wg.Add(1)
		go func(id ItemId, errchan chan error) {
			defer wg.Done()
			item, err := hn.FetchItem(id)
			if err != nil {
				errchan <- err
			} else {
				itemsMap.Store(id, *item)
			}
		}(id, errchan)
	}
	wg.Wait()
	close(errchan)
	err, errors := <-errchan
	if errors {
		// Return the first error we receive.
		return nil, err
	}

	items := make([]Item, len(ids))
	for i, id := range ids {
		mapitem, ok := itemsMap.Load(id)
		if !ok {
			panic(fmt.Sprintf("no item %d in items map\n", id))
		}
		item := mapitem.(Item)
		items[i] = item
	}

	return items, nil
}

func (hn *HnClient) Search(request SearchRequest) (*SearchResponse, error) {
	if request.Limit < 0 || request.Limit > maxStoriesLimit {
		return nil, fmt.Errorf("invalid limit: %d\n", request.Limit)
	}

	var endpoint string
	switch request.Ranking {
	case Popularity:
		endpoint = hn.searchPopularityUrl
	case Date:
		endpoint = hn.searchDateUrl
	}
	query := url.QueryEscape(request.Query)
	tags := url.QueryEscape(request.Tags)
	url := fmt.Sprintf("%s?query=%s&tags=%s&hitsPerPage=%d", endpoint, query, tags, request.Limit)
	response, err := hn.client.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	if response.StatusCode > 299 {
		return nil, fmt.Errorf("response failed with code %d\n", response.StatusCode)
	}

	var searchResponse SearchResponseJson
	if err := json.NewDecoder(response.Body).Decode(&searchResponse); err != nil {
		return nil, err
	}

	var hits []SearchResultJson
	if len(searchResponse.Hits) <= request.Limit {
		hits = searchResponse.Hits
	} else {
		hits = searchResponse.Hits[:request.Limit]
	}
	results := make([]SearchResult, len(hits))
	for i, hit := range hits {
		id, err := strconv.Atoi(hit.Id)
		if err != nil {
			return nil, err
		}
		highlightResultCommentText := hit.HighlightResult.CommentText.Value
		results[i] = SearchResult{
			Id:                         int32(id),
			HighlightResultCommentText: highlightResultCommentText,
		}
	}
	return &SearchResponse{
		Results: results,
	}, nil
}
