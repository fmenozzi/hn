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

func WithFailedResponse(status int) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(status)
	})
}

func TestFetchRankedStoriesIdsSucceedsIfServerReturns200(t *testing.T) {
	server := httptest.NewServer(WithJsonResponse("[123, 456, 789]"))
	defer server.Close()
	client := MakeClientForUrl(server.URL)

	rankedStoriesIds, err := client.FetchRankedStoriesIds(Top, 10)

	assert.Nil(t, err)
	assert.Equal(t, rankedStoriesIds, []ItemId{123, 456, 789})
}

func TestFetchRankedStoriesIdsSucceedsWithLimitOfZero(t *testing.T) {
	server := httptest.NewServer(WithJsonResponse("[123, 456, 789]"))
	defer server.Close()
	client := MakeClientForUrl(server.URL)

	rankedStoriesIds, err := client.FetchRankedStoriesIds(Top, 0)

	assert.Nil(t, err)
	assert.Empty(t, rankedStoriesIds)
}

func TestFetchRankedStoriesIdsFailsWithInvalidLimits(t *testing.T) {
	server := httptest.NewServer(WithJsonResponse("[]")) // unimportant
	defer server.Close()
	client := MakeClientForUrl(server.URL)

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
	client := MakeClientForUrl(server.URL)

	_, err := client.FetchRankedStoriesIds(Top, 10)

	assert.NotNil(t, err)
	assert.ErrorContains(t, err, "500") // internal server error
}

func TestFetchRankedStoriesIdsFailsIfJsonCannotBeParsed(t *testing.T) {
	server := httptest.NewServer(WithJsonResponse("["))
	defer server.Close()
	client := MakeClientForUrl(server.URL)

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
	client := MakeClientForUrl(server.URL)

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
	client := MakeClientForUrl(server.URL)

	_, err := client.FetchItem(123)

	assert.NotNil(t, err)
	assert.ErrorContains(t, err, "500") // internal server error
}

func TestFetchItemFailsIfJsonCannotBeParsed(t *testing.T) {
	server := httptest.NewServer(WithJsonResponse("{"))
	defer server.Close()
	client := MakeClientForUrl(server.URL)

	_, err := client.FetchItem(123)

	assert.NotNil(t, err)
	assert.ErrorContains(t, err, "unexpected end of JSON input")
}
