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

func TestFetchRankedStoriesIdsSucceedsIfServerReturns200(t *testing.T) {
	server := httptest.NewServer(WithJsonResponse("[123, 456, 789]"))
	defer server.Close()
	client := NewHnClientBuilder().SetHnUrl(server.URL).Build()

	rankedStoriesIds, err := client.FetchRankedStoriesIds(Top, 10)

	assert.Nil(t, err)
	assert.Equal(t, rankedStoriesIds, []ItemId{123, 456, 789})
}

func TestFetchRankedStoriesIdsSucceedsWithLimitOfZero(t *testing.T) {
	server := httptest.NewServer(WithJsonResponse("[123, 456, 789]"))
	defer server.Close()
	client := NewHnClientBuilder().SetHnUrl(server.URL).Build()

	rankedStoriesIds, err := client.FetchRankedStoriesIds(Top, 0)

	assert.Nil(t, err)
	assert.Empty(t, rankedStoriesIds)
}

func TestFetchRankedStoriesIdsSucceedsWithLimitLessThanResponseSize(t *testing.T) {
	server := httptest.NewServer(WithJsonResponse("[123, 456, 789]"))
	defer server.Close()
	client := NewHnClientBuilder().SetHnUrl(server.URL).Build()

	rankedStoriesIds, err := client.FetchRankedStoriesIds(Top, 1)

	assert.Nil(t, err)
	assert.Equal(t, rankedStoriesIds, []ItemId{123})
}

func TestFetchRankedStoriesIdsFailsWithInvalidLimits(t *testing.T) {
	server := httptest.NewServer(WithJsonResponse("[]")) // unimportant
	defer server.Close()
	client := NewHnClientBuilder().SetHnUrl(server.URL).Build()

	_, err := client.FetchRankedStoriesIds(Top, -1)
	assert.NotNil(t, err)
	assert.ErrorContains(t, err, "Invalid limit")

	_, err = client.FetchRankedStoriesIds(Top, MaxStoriesLimit+1)
	assert.NotNil(t, err)
	assert.ErrorContains(t, err, "Invalid limit")
}

func TestFetchRankedStoriesIdsFailsIfServerReturns500(t *testing.T) {
	server := httptest.NewServer(WithFailedResponse(http.StatusInternalServerError))
	defer server.Close()
	client := NewHnClientBuilder().SetHnUrl(server.URL).Build()

	_, err := client.FetchRankedStoriesIds(Top, 10)

	assert.NotNil(t, err)
	assert.ErrorContains(t, err, "500") // internal server error
}

func TestFetchRankedStoriesIdsFailsIfJsonCannotBeParsed(t *testing.T) {
	server := httptest.NewServer(WithJsonResponse("["))
	defer server.Close()
	client := NewHnClientBuilder().SetHnUrl(server.URL).Build()

	_, err := client.FetchRankedStoriesIds(Top, 10)

	assert.NotNil(t, err)
	assert.ErrorContains(t, err, "unexpected end of JSON input")
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
	assert.ErrorContains(t, err, "unexpected end of JSON input")
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
	assert.ErrorContains(t, err, "unexpected end of JSON input")
}

func TestSearchSucceedsIfServerReturns200(t *testing.T) {
	server := httptest.NewServer(WithJsonResponse(`
	{
		"hits": [
			{ "story_id": 123 },
			{ "story_id": 456 },
			{ "story_id": 789 }
		]
	}
	`))
	defer server.Close()
	client := NewHnClientBuilder().SetSearchPopularityUrl(server.URL).Build()

	items, err := client.Search(SearchRequest{
		Query:   "query", // unimportant
		Tags:    []string{"story"},
		Ranking: Popularity,
		Limit:   30,
	})

	assert.Nil(t, err)
	assert.Equal(t, &SearchItems{
		Ids:     []ItemId{123, 456, 789},
		Ranking: Popularity,
	}, items)
}

func TestSearchSucceedsWhenSortingByDate(t *testing.T) {
	server := httptest.NewServer(WithJsonResponse(`
	{
		"hits": [
			{ "story_id": 123 },
			{ "story_id": 456 },
			{ "story_id": 789 }
		]
	}
	`))
	defer server.Close()
	client := NewHnClientBuilder().SetSearchDateUrl(server.URL).Build()

	items, err := client.Search(SearchRequest{
		Query:   "query", // unimportant
		Tags:    []string{"story"},
		Ranking: Date,
		Limit:   30,
	})

	assert.Nil(t, err)
	assert.Equal(t, &SearchItems{
		Ids:     []ItemId{123, 456, 789},
		Ranking: Date,
	}, items)
}

func TestSearchSucceedsWithMultiWordQuery(t *testing.T) {
	server := httptest.NewServer(WithJsonResponse(`
	{
		"hits": [
			{ "story_id": 123 },
			{ "story_id": 456 },
			{ "story_id": 789 }
		]
	}
	`))
	defer server.Close()
	client := NewHnClientBuilder().SetSearchPopularityUrl(server.URL).Build()

	items, err := client.Search(SearchRequest{
		Query:   "multi word query",
		Tags:    []string{"story"},
		Ranking: Popularity,
		Limit:   30,
	})

	assert.Nil(t, err)
	assert.Equal(t, &SearchItems{
		Ids:     []ItemId{123, 456, 789},
		Ranking: Popularity,
	}, items)
}

func TestSearchSucceedsWithALimitOfZero(t *testing.T) {
	server := httptest.NewServer(WithJsonResponse(`
	{
		"hits": [
			{ "story_id": 123 },
			{ "story_id": 456 },
			{ "story_id": 789 }
		]
	}
	`))
	defer server.Close()
	client := NewHnClientBuilder().SetSearchPopularityUrl(server.URL).Build()

	items, err := client.Search(SearchRequest{
		Query:   "query", // unimportant
		Tags:    []string{"story"},
		Ranking: Popularity,
		Limit:   0,
	})

	assert.Nil(t, err)
	assert.Equal(t, &SearchItems{
		Ids:     []ItemId{},
		Ranking: Popularity,
	}, items)
}

func TestSearchSucceedsWithALimitLessThanResponseSize(t *testing.T) {
	server := httptest.NewServer(WithJsonResponse(`
	{
		"hits": [
			{ "story_id": 123 },
			{ "story_id": 456 },
			{ "story_id": 789 }
		]
	}
	`))
	defer server.Close()
	client := NewHnClientBuilder().SetSearchPopularityUrl(server.URL).Build()

	items, err := client.Search(SearchRequest{
		Query:   "query", // unimportant
		Tags:    []string{"story"},
		Ranking: Popularity,
		Limit:   1,
	})

	assert.Nil(t, err)
	assert.Equal(t, &SearchItems{
		Ids:     []ItemId{123},
		Ranking: Popularity,
	}, items)
}

func TestSearchFailsWhenServerReturns500(t *testing.T) {
	server := httptest.NewServer(WithFailedResponse(http.StatusInternalServerError))
	defer server.Close()
	client := NewHnClientBuilder().SetSearchPopularityUrl(server.URL).Build()

	_, err := client.Search(SearchRequest{
		Query:   "query", // unimportant
		Tags:    []string{"story"},
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
	assert.ErrorContains(t, err, "Invalid limit")

	_, err = client.Search(SearchRequest{Limit: MaxStoriesLimit + 1})
	assert.NotNil(t, err)
	assert.ErrorContains(t, err, "Invalid limit")
}

func TestSearchFailsIfJsonCannotBeParsed(t *testing.T) {
	server := httptest.NewServer(WithJsonResponse("["))
	defer server.Close()
	client := NewHnClientBuilder().SetSearchPopularityUrl(server.URL).Build()

	_, err := client.Search(SearchRequest{
		Query:   "query", // unimportant
		Tags:    []string{"story"},
		Ranking: Popularity,
		Limit:   30,
	})

	assert.NotNil(t, err)
	assert.ErrorContains(t, err, "unexpected end of JSON input")
}
