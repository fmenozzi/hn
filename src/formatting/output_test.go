package formatting

import (
	"testing"
	"time"

	"github.com/fmenozzi/hn/src/api"
	"github.com/stretchr/testify/assert"
)

type Addressable interface {
	int32 | int64 | string
}

func ptr[T Addressable](t T) *T {
	return &t
}

func intptr(i int) *int32 {
	return ptr(int32(i))
}

var (
	now       = time.Unix(10000000, 0)
	fakeClock = FakeClock{now}

	job = api.Item{
		Id:    1,
		Score: intptr(1),
		By:    ptr("jobuser"),
		Time:  ptr(now.Add(-6 * time.Hour).Unix()), // 6 hours ago
		Title: ptr("Job title"),
	}

	story = api.Item{
		Id:          2,
		Score:       intptr(10),
		By:          ptr("storyuser"),
		Time:        ptr(now.Add(-12 * 24 * time.Hour).Unix()), // 12 days ago
		Descendants: intptr(20),
		Title:       ptr("Story title"),
		Url:         ptr("www.story.url"),
	}

	poll = api.Item{
		Id:          3,
		Score:       intptr(100),
		By:          ptr("polluser"),
		Time:        ptr(now.Add(-40 * time.Minute).Unix()), // 40 minutes ago
		Descendants: intptr(200),
		Title:       ptr("Poll title"),
	}

	pollopt = api.Item{
		Id:    4,
		Score: intptr(1000),
		By:    ptr("polloptuser"),
		Time:  ptr(now.Add(-3 * 30 * 24 * time.Hour).Unix()), // 3 months ago
		Text:  ptr("Poll option text"),
	}

	comment = api.Item{
		Id:   5,
		By:   ptr("commentuser"),
		Time: ptr(now.Add(-1 * 24 * time.Hour).Unix()), // a day ago
		Text: ptr("Comment text"),
		Kids: []api.ItemId{6, 7, 8, 9},
	}
)

func TestPlainOutput(t *testing.T) {
	jobOutput := JobOutput(&job, Plain, &fakeClock)
	storyOutput := StoryOutput(&story, Plain, &fakeClock)
	pollOutput := PollOutput(&poll, Plain, &fakeClock)
	pollOptOutput := PollOptOutput(&pollopt, Plain, &fakeClock)
	commentOutput := CommentOutput(&comment, Plain, &fakeClock)

	expectedJobOutput := "[   1 pts] [   6 hours ago] [       HIRING] Job title\n"
	expectedStoryOutput := "[  10 pts] [   12 days ago] [  20 comments] www.story.url\n"
	expectedPollOutput := "[ 100 pts] [    40 min ago] [ 200 comments] https://news.ycombinator.com/item?id=3\n"
	expectedPollOptOutput := "[1000 pts] [  3 months ago] [             ] Poll option text\n"
	expectedCommentOutput := "[        ] [     a day ago] [             ] Comment text\n"

	assert.Equal(t, expectedJobOutput, jobOutput)
	assert.Equal(t, expectedStoryOutput, storyOutput)
	assert.Equal(t, expectedPollOutput, pollOutput)
	assert.Equal(t, expectedPollOptOutput, pollOptOutput)
	assert.Equal(t, expectedCommentOutput, commentOutput)
}

func TestMarkdownOutput(t *testing.T) {
	jobOutput := JobOutput(&job, Markdown, &fakeClock)
	storyOutput := StoryOutput(&story, Markdown, &fakeClock)
	pollOutput := PollOutput(&poll, Markdown, &fakeClock)
	pollOptOutput := PollOptOutput(&pollopt, Markdown, &fakeClock)
	commentOutput := CommentOutput(&comment, Markdown, &fakeClock)

	expectedJobOutput := "* [   1 pts] [   6 hours ago] [       [HIRING](https://news.ycombinator.com/item?id=1)] [Job title](https://news.ycombinator.com/item?id=1)\n"
	expectedStoryOutput := "* [  10 pts] [   12 days ago] [[  20 comments](https://news.ycombinator.com/item?id=2)] [Story title](www.story.url)\n"
	expectedPollOutput := "* [ 100 pts] [    40 min ago] [[ 200 comments](https://news.ycombinator.com/item?id=3)] [Poll title](https://news.ycombinator.com/item?id=3)\n"
	expectedPollOptOutput := "* [1000 pts] [  3 months ago] [             ] [Poll option text](https://news.ycombinator.com/item?id=4)\n"
	expectedCommentOutput := "* [        ] [     a day ago] [             ] [Comment text](https://news.ycombinator.com/item?id=5)\n"

	assert.Equal(t, expectedJobOutput, jobOutput)
	assert.Equal(t, expectedStoryOutput, storyOutput)
	assert.Equal(t, expectedPollOutput, pollOutput)
	assert.Equal(t, expectedPollOptOutput, pollOptOutput)
	assert.Equal(t, expectedCommentOutput, commentOutput)
}

func TestCsvOutput(t *testing.T) {
	jobOutput := JobOutput(&job, Csv, &fakeClock)
	storyOutput := StoryOutput(&story, Csv, &fakeClock)
	pollOutput := PollOutput(&poll, Csv, &fakeClock)
	pollOptOutput := PollOptOutput(&pollopt, Csv, &fakeClock)
	commentOutput := CommentOutput(&comment, Csv, &fakeClock)

	expectedJobOutput := "1,job,jobuser,9978400,\"Job title\",https://news.ycombinator.com/item?id=1,1,0\n"
	expectedStoryOutput := "2,story,storyuser,8963200,\"Story title\",www.story.url,10,20\n"
	expectedPollOutput := "3,poll,polluser,9997600,\"Poll title\",https://news.ycombinator.com/item?id=3,100,200\n"
	expectedPollOptOutput := "4,pollopt,polloptuser,2224000,\"Poll option text\",https://news.ycombinator.com/item?id=4,1000,0\n"
	expectedCommentOutput := "5,comment,commentuser,9913600,\"Comment text\",https://news.ycombinator.com/item?id=5,0,4\n"

	assert.Equal(t, expectedJobOutput, jobOutput)
	assert.Equal(t, expectedStoryOutput, storyOutput)
	assert.Equal(t, expectedPollOutput, pollOutput)
	assert.Equal(t, expectedPollOptOutput, pollOptOutput)
	assert.Equal(t, expectedCommentOutput, commentOutput)
}

func TestStoryWithoutUrlFallbackToPostUrl(t *testing.T) {
	story = api.Item{
		Id:          2,
		Score:       intptr(10),
		By:          ptr("storyuser"),
		Time:        ptr(now.Add(-12 * 24 * time.Hour).Unix()), // 12 days ago
		Descendants: intptr(20),
		Title:       ptr("Story title"),
	}

	storyOutput := StoryOutput(&story, Csv, &fakeClock)

	expectedStoryOutput := "2,story,storyuser,8963200,\"Story title\",https://news.ycombinator.com/item?id=2,10,20\n"

	assert.Equal(t, expectedStoryOutput, storyOutput)
}
