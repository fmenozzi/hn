package formatting

import (
	"testing"

	"github.com/fmenozzi/hn/src/api"
	"github.com/stretchr/testify/assert"
)

type Addressable interface {
	api.ItemId | int | string
}

func ptr[T Addressable](t T) *T {
	return &t
}

func intptr(i int) *int32 {
	return ptr(int32(i))
}

func strptr(s string) *string {
	return ptr(s)
}

var (
	job = api.Item{
		Id:    1,
		Score: intptr(1),
		By:    strptr("jobuser"),
		Title: strptr("Job title"),
	}

	story = api.Item{
		Id:          2,
		Score:       intptr(10),
		By:          strptr("storyuser"),
		Descendants: intptr(20),
		Title:       strptr("Story title"),
		Url:         strptr("www.story.url"),
	}

	poll = api.Item{
		Id:          3,
		Score:       intptr(100),
		By:          strptr("polluser"),
		Descendants: intptr(200),
		Title:       strptr("Poll title"),
	}
)

func TestPlainOutput(t *testing.T) {
	jobOutput := JobOutput(&job, Plain)
	storyOutput := StoryOutput(&story, Plain)
	pollOutput := PollOutput(&poll, Plain)

	expectedJobOutput := "[   1 pts] [       HIRING] Job title\n"
	expectedStoryOutput := "[  10 pts] [  20 comments] www.story.url\n"
	expectedPollOutput := "[ 100 pts] [ 200 comments] https://news.ycombinator.com/item?id=3\n"

	assert.Equal(t, expectedJobOutput, jobOutput)
	assert.Equal(t, expectedStoryOutput, storyOutput)
	assert.Equal(t, expectedPollOutput, pollOutput)
}

func TestMarkdownOutput(t *testing.T) {
	jobOutput := JobOutput(&job, Markdown)
	storyOutput := StoryOutput(&story, Markdown)
	pollOutput := PollOutput(&poll, Markdown)

	expectedJobOutput := "* [   1 pts] [       [HIRING](https://news.ycombinator.com/item?id=1)] [Job title](https://news.ycombinator.com/item?id=1)\n"
	expectedStoryOutput := "* [  10 pts] [[  20 comments](https://news.ycombinator.com/item?id=2)] [Story title](www.story.url)\n"
	expectedPollOutput := "* [ 100 pts] [[ 200 comments](https://news.ycombinator.com/item?id=3)] [Poll title](https://news.ycombinator.com/item?id=3)\n"

	assert.Equal(t, expectedJobOutput, jobOutput)
	assert.Equal(t, expectedStoryOutput, storyOutput)
	assert.Equal(t, expectedPollOutput, pollOutput)
}

func TestCsvOutput(t *testing.T) {
	jobOutput := JobOutput(&job, Csv)
	storyOutput := StoryOutput(&story, Csv)
	pollOutput := PollOutput(&poll, Csv)

	expectedJobOutput := "1,job,jobuser,\"Job title\",https://news.ycombinator.com/item?id=1,1,0\n"
	expectedStoryOutput := "2,story,storyuser,\"Story title\",www.story.url,10,20\n"
	expectedPollOutput := "3,poll,polluser,\"Poll title\",https://news.ycombinator.com/item?id=3,100,200\n"

	assert.Equal(t, expectedJobOutput, jobOutput)
	assert.Equal(t, expectedStoryOutput, storyOutput)
	assert.Equal(t, expectedPollOutput, pollOutput)
}

func TestStoryWithoutUrlFallbackToPostUrl(t *testing.T) {
	story = api.Item{
		Id:          2,
		Score:       intptr(10),
		By:          strptr("storyuser"),
		Descendants: intptr(20),
		Title:       strptr("Story title"),
	}

	storyOutput := StoryOutput(&story, Csv)

	expectedStoryOutput := "2,story,storyuser,\"Story title\",https://news.ycombinator.com/item?id=2,10,20\n"

	assert.Equal(t, expectedStoryOutput, storyOutput)
}
