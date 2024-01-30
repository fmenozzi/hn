package api

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func WithJsonResponse(json string) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, json)
	})
}

func WithMultipleJsonResponses(responses map[string]string) *http.ServeMux {
	mux := http.NewServeMux()
	for path, json := range responses {
		json := json
		path := path
		mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, json)
		})
	}
	return mux
}

func WithFailedResponse(status int) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(status)
	})
}

func TestFetchFrontPageItemIdsSucceedsIfServerReturns200(t *testing.T) {
	server := httptest.NewServer(WithJsonResponse("[123, 456, 789]"))
	defer server.Close()
	client := NewHnClientBuilder().SetHnUrl(server.URL).Build()

	ids, err := client.FetchFrontPageItemIds(Top, 10)

	assert.Nil(t, err)
	assert.Equal(t, ids, []ItemId{123, 456, 789})
}

func TestFetchFrontPageItemIdsSucceedsWithLimitOfZero(t *testing.T) {
	server := httptest.NewServer(WithJsonResponse("[123, 456, 789]"))
	defer server.Close()
	client := NewHnClientBuilder().SetHnUrl(server.URL).Build()

	ids, err := client.FetchFrontPageItemIds(Top, 0)

	assert.Nil(t, err)
	assert.Empty(t, ids)
}

func TestFetchFrontPageItemIdsSucceedsWithLimitLessThanResponseSize(t *testing.T) {
	server := httptest.NewServer(WithJsonResponse("[123, 456, 789]"))
	defer server.Close()
	client := NewHnClientBuilder().SetHnUrl(server.URL).Build()

	ids, err := client.FetchFrontPageItemIds(Top, 1)

	assert.Nil(t, err)
	assert.Equal(t, ids, []ItemId{123})
}

func TestFetchFrontPageItemIdsFailsWithInvalidLimits(t *testing.T) {
	server := httptest.NewServer(WithJsonResponse("[]")) // unimportant
	defer server.Close()
	client := NewHnClientBuilder().SetHnUrl(server.URL).Build()

	_, err := client.FetchFrontPageItemIds(Top, -1)
	assert.NotNil(t, err)
	assert.ErrorContains(t, err, "invalid limit")

	_, err = client.FetchFrontPageItemIds(Top, maxStoriesLimit+1)
	assert.NotNil(t, err)
	assert.ErrorContains(t, err, "invalid limit")
}

func TestFetchFrontPageItemIdsFailsIfServerReturns500(t *testing.T) {
	server := httptest.NewServer(WithFailedResponse(http.StatusInternalServerError))
	defer server.Close()
	client := NewHnClientBuilder().SetHnUrl(server.URL).Build()

	_, err := client.FetchFrontPageItemIds(Top, 10)

	assert.NotNil(t, err)
	assert.ErrorContains(t, err, "500") // internal server error
}

func TestFetchFrontPageItemIdsFailsIfJsonCannotBeParsed(t *testing.T) {
	server := httptest.NewServer(WithJsonResponse("["))
	defer server.Close()
	client := NewHnClientBuilder().SetHnUrl(server.URL).Build()

	_, err := client.FetchFrontPageItemIds(Top, 10)

	assert.NotNil(t, err)
	assert.ErrorContains(t, err, "unexpected EOF")
}

func TestFetchItemSucceedsIfServerReturns200(t *testing.T) {
	server := httptest.NewServer(WithJsonResponse(`
	{
		"id": 123,
		"type": "story",
		"by": "username",
		"score": 456
	}
	`))
	defer server.Close()
	client := NewHnClientBuilder().SetHnUrl(server.URL).Build()

	item, err := client.FetchItem(123)

	username := "username"
	score := int32(456)
	assert.Nil(t, err)
	assert.Equal(t, item, &Item{
		Id:    123,
		Type:  Story,
		By:    &username,
		Score: &score,
	})
}

func TestFetchItemFailsIfServerReturns500(t *testing.T) {
	server := httptest.NewServer(WithFailedResponse(http.StatusInternalServerError))
	defer server.Close()
	client := NewHnClientBuilder().SetHnUrl(server.URL).Build()

	_, err := client.FetchItem(123)

	assert.NotNil(t, err)
	assert.ErrorContains(t, err, "500") // internal server error
}

func TestFetchItemFailsIfJsonCannotBeParsed(t *testing.T) {
	server := httptest.NewServer(WithJsonResponse("{"))
	defer server.Close()
	client := NewHnClientBuilder().SetHnUrl(server.URL).Build()

	_, err := client.FetchItem(123)

	assert.NotNil(t, err)
	assert.ErrorContains(t, err, "unexpected EOF")
}

func TestFetchItemsSucceedsIfServerReturns200(t *testing.T) {
	server := httptest.NewServer(WithMultipleJsonResponses(map[string]string{
		"/item/123.json": `{ "id": 123, "type": "story" }`,
		"/item/456.json": `{ "id": 456, "type": "job" }`,
		"/item/789.json": `{ "id": 789, "type": "poll" }`,
	}))
	defer server.Close()
	client := NewHnClientBuilder().SetHnUrl(server.URL).Build()

	items, err := client.FetchItems([]ItemId{123, 456, 789})

	assert.Nil(t, err)
	assert.Equal(t, items, []Item{
		{
			Id:   123,
			Type: Story,
		},
		{
			Id:   456,
			Type: Job,
		},
		{
			Id:   789,
			Type: Poll,
		},
	})
}

func TestFetchItemsFailsIfServerReturns500(t *testing.T) {
	server := httptest.NewServer(WithFailedResponse(http.StatusInternalServerError))
	defer server.Close()
	client := NewHnClientBuilder().SetHnUrl(server.URL).Build()

	_, err := client.FetchItems([]ItemId{123})

	assert.NotNil(t, err)
	assert.ErrorContains(t, err, "500") // internal server error
}

func TestFetchItemsFailsIfJsonCannotBeParsed(t *testing.T) {
	server := httptest.NewServer(WithJsonResponse("{"))
	defer server.Close()
	client := NewHnClientBuilder().SetHnUrl(server.URL).Build()

	_, err := client.FetchItems([]ItemId{123})

	assert.NotNil(t, err)
	assert.ErrorContains(t, err, "unexpected EOF")
}

func TestSearchSucceedsIfServerReturns200(t *testing.T) {
	server := httptest.NewServer(WithJsonResponse(`
	{
		"hits": [
			{ "objectID": "123" },
			{ "objectID": "456" },
			{ "objectID": "789" }
		]
	}
	`))
	defer server.Close()
	client := NewHnClientBuilder().SetSearchPopularityUrl(server.URL).Build()

	response, err := client.Search(SearchRequest{
		Query:   "query", // unimportant
		Tags:    "story",
		Ranking: Popularity,
		Limit:   30,
	})

	assert.Nil(t, err)
	assert.Equal(t, &SearchResponse{Results: []SearchResult{
		{Id: 123},
		{Id: 456},
		{Id: 789},
	}}, response)
}

func TestSearchSucceedsWhenSortingByDate(t *testing.T) {
	server := httptest.NewServer(WithJsonResponse(`
	{
		"hits": [
			{ "objectID": "123" },
			{ "objectID": "456" },
			{ "objectID": "789" }
		]
	}
	`))
	defer server.Close()
	client := NewHnClientBuilder().SetSearchDateUrl(server.URL).Build()

	response, err := client.Search(SearchRequest{
		Query:   "query", // unimportant
		Tags:    "story",
		Ranking: Date,
		Limit:   30,
	})

	assert.Nil(t, err)
	assert.Equal(t, &SearchResponse{Results: []SearchResult{
		{Id: 123},
		{Id: 456},
		{Id: 789},
	}}, response)
}

func TestSearchSucceedsWithMultiWordQuery(t *testing.T) {
	server := httptest.NewServer(WithJsonResponse(`
	{
		"hits": [
			{ "objectID": "123" },
			{ "objectID": "456" },
			{ "objectID": "789" }
		]
	}
	`))
	defer server.Close()
	client := NewHnClientBuilder().SetSearchPopularityUrl(server.URL).Build()

	response, err := client.Search(SearchRequest{
		Query:   "multi word query",
		Tags:    "story",
		Ranking: Popularity,
		Limit:   30,
	})

	assert.Nil(t, err)
	assert.Equal(t, &SearchResponse{Results: []SearchResult{
		{Id: 123},
		{Id: 456},
		{Id: 789},
	}}, response)
}

func TestSearchSucceedsWithALimitOfZero(t *testing.T) {
	server := httptest.NewServer(WithJsonResponse(`
	{
		"hits": [
			{ "objectID": "123" },
			{ "objectID": "456" },
			{ "objectID": "789" }
		]
	}
	`))
	defer server.Close()
	client := NewHnClientBuilder().SetSearchPopularityUrl(server.URL).Build()

	response, err := client.Search(SearchRequest{
		Query:   "query", // unimportant
		Tags:    "story",
		Ranking: Popularity,
		Limit:   0,
	})

	assert.Nil(t, err)
	assert.Empty(t, response.Results)
}

func TestSearchSucceedsWithALimitLessThanResponseSize(t *testing.T) {
	server := httptest.NewServer(WithJsonResponse(`
	{
		"hits": [
			{ "objectID": "123" },
			{ "objectID": "456" },
			{ "objectID": "789" }
		]
	}
	`))
	defer server.Close()
	client := NewHnClientBuilder().SetSearchPopularityUrl(server.URL).Build()

	response, err := client.Search(SearchRequest{
		Query:   "query", // unimportant
		Tags:    "story",
		Ranking: Popularity,
		Limit:   1,
	})

	assert.Nil(t, err)
	assert.Equal(t, &SearchResponse{Results: []SearchResult{
		{Id: 123},
	}}, response)
}

func TestSearchFailsWhenServerReturns500(t *testing.T) {
	server := httptest.NewServer(WithFailedResponse(http.StatusInternalServerError))
	defer server.Close()
	client := NewHnClientBuilder().SetSearchPopularityUrl(server.URL).Build()

	_, err := client.Search(SearchRequest{
		Query:   "query", // unimportant
		Tags:    "story",
		Ranking: Popularity,
		Limit:   30,
	})

	assert.NotNil(t, err)
	assert.ErrorContains(t, err, "500") // internal server error
}

func TestSearchFailsWithInvalidLimits(t *testing.T) {
	server := httptest.NewServer(WithJsonResponse("[]")) // unimportant
	defer server.Close()
	client := NewHnClientBuilder().SetSearchPopularityUrl(server.URL).Build()

	_, err := client.Search(SearchRequest{Limit: -1})
	assert.NotNil(t, err)
	assert.ErrorContains(t, err, "invalid limit")

	_, err = client.Search(SearchRequest{Limit: maxStoriesLimit + 1})
	assert.NotNil(t, err)
	assert.ErrorContains(t, err, "invalid limit")
}

func TestSearchFailsIfJsonCannotBeParsed(t *testing.T) {
	server := httptest.NewServer(WithJsonResponse("["))
	defer server.Close()
	client := NewHnClientBuilder().SetSearchPopularityUrl(server.URL).Build()

	_, err := client.Search(SearchRequest{
		Query:   "query", // unimportant
		Tags:    "story",
		Ranking: Popularity,
		Limit:   30,
	})

	assert.NotNil(t, err)
	assert.ErrorContains(t, err, "unexpected EOF")
}
