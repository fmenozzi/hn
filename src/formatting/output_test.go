package formatting

import (
	"fmt"
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
	jobOutput := jobOutput(&job, Plain, &fakeClock)
	storyOutput := storyOutput(&story, Plain, &fakeClock)
	pollOutput := pollOutput(&poll, Plain, &fakeClock)
	pollOptOutput := pollOptOutput(&pollopt, Plain, &fakeClock)
	commentOutput := commentOutput(&comment, Plain, &fakeClock)

	expectedJobOutput := "HIRING: https://news.ycombinator.com/item?id=1\n└─── 1 pt 6 hours ago\n"
	expectedStoryOutput := "www.story.url\n└─── 10 pts by storyuser 12 days ago | 20 comments\n"
	expectedPollOutput := "https://news.ycombinator.com/item?id=3\n└─── 100 pts by polluser 40 min ago | 200 comments\n"
	expectedpollOptOutput := "Poll option text\n└─── 1000 pts by polloptuser 3 months ago\n"
	expectedcommentOutput := "Comment text\n└─── by commentuser a day ago | 4 replies\n"

	assert.Equal(t, expectedJobOutput, jobOutput)
	assert.Equal(t, expectedStoryOutput, storyOutput)
	assert.Equal(t, expectedPollOutput, pollOutput)
	assert.Equal(t, expectedpollOptOutput, pollOptOutput)
	assert.Equal(t, expectedcommentOutput, commentOutput)
}

func TestMarkdownOutput(t *testing.T) {
	jobOutput := jobOutput(&job, Markdown, &fakeClock)
	storyOutput := storyOutput(&story, Markdown, &fakeClock)
	pollOutput := pollOutput(&poll, Markdown, &fakeClock)
	pollOptOutput := pollOptOutput(&pollopt, Markdown, &fakeClock)
	commentOutput := commentOutput(&comment, Markdown, &fakeClock)

	expectedJobOutput := "* **[HIRING: Job title](https://news.ycombinator.com/item?id=1)**\n* └─── 1 pt 6 hours ago\n"
	expectedStoryOutput := "* **[Story title](www.story.url)**\n* └─── 10 pts by [storyuser](https://news.ycombinator.com/user?id=storyuser) 12 days ago | [20 comments](https://news.ycombinator.com/item?id=2)\n"
	expectedPollOutput := "* **[Poll title](https://news.ycombinator.com/item?id=3)**\n* └─── 100 pts by [polluser](https://news.ycombinator.com/user?id=polluser) 40 min ago | [200 comments](https://news.ycombinator.com/item?id=3)\n"
	expectedpollOptOutput := "* **[Poll option text](https://news.ycombinator.com/item?id=4)**\n* └─── 1000 pts by [polloptuser](https://news.ycombinator.com/user?id=polloptuser) 3 months ago\n"
	expectedcommentOutput := "* *[Comment text](https://news.ycombinator.com/item?id=5)*\n* └─── by [commentuser](https://news.ycombinator.com/user?id=commentuser) a day ago | [4 replies](https://news.ycombinator.com/item?id=5)\n"

	assert.Equal(t, expectedJobOutput, jobOutput)
	assert.Equal(t, expectedStoryOutput, storyOutput)
	assert.Equal(t, expectedPollOutput, pollOutput)
	assert.Equal(t, expectedpollOptOutput, pollOptOutput)
	assert.Equal(t, expectedcommentOutput, commentOutput)
}

func TestStoryWithoutUrlFallbackToPostUrl(t *testing.T) {
	story := api.Item{
		Id:          2,
		Score:       intptr(10),
		By:          ptr("storyuser"),
		Time:        ptr(now.Add(-12 * 24 * time.Hour).Unix()), // 12 days ago
		Descendants: intptr(20),
		Title:       ptr("Story title"),
	}

	storyOutput := storyOutput(&story, Plain, &fakeClock)

	expectedStoryOutput := "https://news.ycombinator.com/item?id=2\n└─── 10 pts by storyuser 12 days ago | 20 comments\n"

	assert.Equal(t, expectedStoryOutput, storyOutput)
}

func TestSingularOutput(t *testing.T) {
	job, story, poll, pollopt, comment := job, story, poll, pollopt, comment

	job.Score = intptr(1)
	story.Score = intptr(1)
	story.Descendants = intptr(1)
	poll.Score = intptr(1)
	poll.Descendants = intptr(1)
	pollopt.Score = intptr(1)
	comment.Kids = []api.ItemId{3}

	fmt.Printf("story: %v\n", story)
	fmt.Printf("story.Url: %v\n", story.Url)

	jobOutput := jobOutput(&job, Plain, &fakeClock)
	storyOutput := storyOutput(&story, Plain, &fakeClock)
	pollOutput := pollOutput(&poll, Plain, &fakeClock)
	pollOptOutput := pollOptOutput(&pollopt, Plain, &fakeClock)
	commentOutput := commentOutput(&comment, Plain, &fakeClock)

	expectedJobOutput := "HIRING: https://news.ycombinator.com/item?id=1\n└─── 1 pt 6 hours ago\n"
	expectedStoryOutput := "www.story.url\n└─── 1 pt by storyuser 12 days ago | 1 comment\n"
	expectedPollOutput := "https://news.ycombinator.com/item?id=3\n└─── 1 pt by polluser 40 min ago | 1 comment\n"
	expectedpollOptOutput := "Poll option text\n└─── 1 pt by polloptuser 3 months ago\n"
	expectedcommentOutput := "Comment text\n└─── by commentuser a day ago | 1 reply\n"

	assert.Equal(t, expectedJobOutput, jobOutput)
	assert.Equal(t, expectedStoryOutput, storyOutput)
	assert.Equal(t, expectedPollOutput, pollOutput)
	assert.Equal(t, expectedpollOptOutput, pollOptOutput)
	assert.Equal(t, expectedcommentOutput, commentOutput)
}
