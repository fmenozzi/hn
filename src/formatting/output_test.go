package formatting

import (
	"bytes"
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
		Type:  api.Job,
		Score: intptr(1),
		By:    ptr("jobuser"),
		Time:  ptr(now.Add(-6 * time.Hour).Unix()), // 6 hours ago
		Title: ptr("Job title"),
	}

	story = api.Item{
		Id:          2,
		Type:        api.Story,
		Score:       intptr(10),
		By:          ptr("storyuser"),
		Time:        ptr(now.Add(-12 * 24 * time.Hour).Unix()), // 12 days ago
		Descendants: intptr(20),
		Title:       ptr("Story title"),
		Url:         ptr("www.story.url"),
	}

	poll = api.Item{
		Id:          3,
		Type:        api.Poll,
		Score:       intptr(100),
		By:          ptr("polluser"),
		Time:        ptr(now.Add(-40 * time.Minute).Unix()), // 40 minutes ago
		Descendants: intptr(200),
		Title:       ptr("Poll title"),
	}

	pollopt = api.Item{
		Id:    4,
		Type:  api.PollOpt,
		Score: intptr(1000),
		By:    ptr("polloptuser"),
		Time:  ptr(now.Add(-3 * 30 * 24 * time.Hour).Unix()), // 3 months ago
		Text:  ptr("Poll option text"),
	}

	comment = api.Item{
		Id:   5,
		Type: api.Comment,
		By:   ptr("commentuser"),
		Time: ptr(now.Add(-1 * 24 * time.Hour).Unix()), // a day ago
		Text: ptr("Comment text"),
		Kids: []api.ItemId{6, 7, 8, 9},
	}
)

func TestPlainOutput(t *testing.T) {
	var jobOutput, storyOutput, pollOutput, pollOptOutput, commentOutput bytes.Buffer

	writeJobItem(&job, Plain, &fakeClock, &jobOutput)
	writeStoryItem(&story, Plain, &fakeClock, &storyOutput)
	writePollItem(&poll, Plain, &fakeClock, &pollOutput)
	writePollOptItem(&pollopt, Plain, &fakeClock, &pollOptOutput)
	writeCommentItem(&comment, Plain, &fakeClock, &commentOutput)

	expectedJobOutput := "HIRING: https://news.ycombinator.com/item?id=1\n└─── 1 pt 6 hours ago\n"
	expectedStoryOutput := "www.story.url\n└─── 10 pts by storyuser 12 days ago | 20 comments\n"
	expectedPollOutput := "https://news.ycombinator.com/item?id=3\n└─── 100 pts by polluser 40 min ago | 200 comments\n"
	expectedpollOptOutput := "Poll option text\n└─── 1000 pts by polloptuser 3 months ago\n"
	expectedcommentOutput := "Comment text\n└─── by commentuser a day ago | 4 replies\n"

	assert.Equal(t, expectedJobOutput, jobOutput.String())
	assert.Equal(t, expectedStoryOutput, storyOutput.String())
	assert.Equal(t, expectedPollOutput, pollOutput.String())
	assert.Equal(t, expectedpollOptOutput, pollOptOutput.String())
	assert.Equal(t, expectedcommentOutput, commentOutput.String())
}

func TestMarkdownOutput(t *testing.T) {
	var jobOutput, storyOutput, pollOutput, pollOptOutput, commentOutput bytes.Buffer

	writeJobItem(&job, Markdown, &fakeClock, &jobOutput)
	writeStoryItem(&story, Markdown, &fakeClock, &storyOutput)
	writePollItem(&poll, Markdown, &fakeClock, &pollOutput)
	writePollOptItem(&pollopt, Markdown, &fakeClock, &pollOptOutput)
	writeCommentItem(&comment, Markdown, &fakeClock, &commentOutput)

	expectedJobOutput := "* **[HIRING: Job title](https://news.ycombinator.com/item?id=1)**\n* └─── 1 pt 6 hours ago\n"
	expectedStoryOutput := "* **[Story title](www.story.url)**\n* └─── 10 pts by [storyuser](https://news.ycombinator.com/user?id=storyuser) 12 days ago | [20 comments](https://news.ycombinator.com/item?id=2)\n"
	expectedPollOutput := "* **[Poll title](https://news.ycombinator.com/item?id=3)**\n* └─── 100 pts by [polluser](https://news.ycombinator.com/user?id=polluser) 40 min ago | [200 comments](https://news.ycombinator.com/item?id=3)\n"
	expectedpollOptOutput := "* **[Poll option text](https://news.ycombinator.com/item?id=4)**\n* └─── 1000 pts by [polloptuser](https://news.ycombinator.com/user?id=polloptuser) 3 months ago\n"
	expectedcommentOutput := "* *[Comment text](https://news.ycombinator.com/item?id=5)*\n* └─── by [commentuser](https://news.ycombinator.com/user?id=commentuser) a day ago | [4 replies](https://news.ycombinator.com/item?id=5)\n"

	assert.Equal(t, expectedJobOutput, jobOutput.String())
	assert.Equal(t, expectedStoryOutput, storyOutput.String())
	assert.Equal(t, expectedPollOutput, pollOutput.String())
	assert.Equal(t, expectedpollOptOutput, pollOptOutput.String())
	assert.Equal(t, expectedcommentOutput, commentOutput.String())
}

func TestJsonOutput(t *testing.T) {
	var jsonOutput bytes.Buffer
	WriteJson([]api.Item{job, story, poll, pollopt, comment}, &jsonOutput)

	expectedJsonOutput := `[
	{
		"id": 1,
		"deleted": null,
		"type": "job",
		"by": "jobuser",
		"time": 9978400,
		"text": null,
		"dead": null,
		"parent": null,
		"poll": null,
		"kids": null,
		"url": null,
		"score": 1,
		"title": "Job title",
		"parts": null,
		"descendants": null
	},
	{
		"id": 2,
		"deleted": null,
		"type": "story",
		"by": "storyuser",
		"time": 8963200,
		"text": null,
		"dead": null,
		"parent": null,
		"poll": null,
		"kids": null,
		"url": "www.story.url",
		"score": 10,
		"title": "Story title",
		"parts": null,
		"descendants": 20
	},
	{
		"id": 3,
		"deleted": null,
		"type": "poll",
		"by": "polluser",
		"time": 9997600,
		"text": null,
		"dead": null,
		"parent": null,
		"poll": null,
		"kids": null,
		"url": null,
		"score": 100,
		"title": "Poll title",
		"parts": null,
		"descendants": 200
	},
	{
		"id": 4,
		"deleted": null,
		"type": "pollopt",
		"by": "polloptuser",
		"time": 2224000,
		"text": "Poll option text",
		"dead": null,
		"parent": null,
		"poll": null,
		"kids": null,
		"url": null,
		"score": 1000,
		"title": null,
		"parts": null,
		"descendants": null
	},
	{
		"id": 5,
		"deleted": null,
		"type": "comment",
		"by": "commentuser",
		"time": 9913600,
		"text": "Comment text",
		"dead": null,
		"parent": null,
		"poll": null,
		"kids": [
			6,
			7,
			8,
			9
		],
		"url": null,
		"score": null,
		"title": null,
		"parts": null,
		"descendants": null
	}
]
`
	assert.Equal(t, expectedJsonOutput, jsonOutput.String())
}

func TestCsvOutput(t *testing.T) {
	var csvOutput bytes.Buffer
	WriteCsv([]api.Item{job, story, poll, pollopt, comment}, &csvOutput)

	expectedCsvOutput :=
		`1,,job,jobuser,9978400,,,0,0,,,1,Job title,,0
2,,story,storyuser,8963200,,,0,0,,www.story.url,10,Story title,,20
3,,poll,polluser,9997600,,,0,0,,,100,Poll title,,200
4,,pollopt,polloptuser,2224000,Poll option text,,0,0,,,1000,,,0
5,,comment,commentuser,9913600,Comment text,,0,0,"6,7,8,9",,0,,,0
`
	assert.Equal(t, expectedCsvOutput, csvOutput.String())
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

	var storyOutput bytes.Buffer
	writeStoryItem(&story, Plain, &fakeClock, &storyOutput)

	expectedStoryOutput := "https://news.ycombinator.com/item?id=2\n└─── 10 pts by storyuser 12 days ago | 20 comments\n"

	assert.Equal(t, expectedStoryOutput, storyOutput.String())
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

	var jobOutput, storyOutput, pollOutput, pollOptOutput, commentOutput bytes.Buffer

	writeJobItem(&job, Plain, &fakeClock, &jobOutput)
	writeStoryItem(&story, Plain, &fakeClock, &storyOutput)
	writePollItem(&poll, Plain, &fakeClock, &pollOutput)
	writePollOptItem(&pollopt, Plain, &fakeClock, &pollOptOutput)
	writeCommentItem(&comment, Plain, &fakeClock, &commentOutput)

	expectedJobOutput := "HIRING: https://news.ycombinator.com/item?id=1\n└─── 1 pt 6 hours ago\n"
	expectedStoryOutput := "www.story.url\n└─── 1 pt by storyuser 12 days ago | 1 comment\n"
	expectedPollOutput := "https://news.ycombinator.com/item?id=3\n└─── 1 pt by polluser 40 min ago | 1 comment\n"
	expectedpollOptOutput := "Poll option text\n└─── 1 pt by polloptuser 3 months ago\n"
	expectedcommentOutput := "Comment text\n└─── by commentuser a day ago | 1 reply\n"

	assert.Equal(t, expectedJobOutput, jobOutput.String())
	assert.Equal(t, expectedStoryOutput, storyOutput.String())
	assert.Equal(t, expectedPollOutput, pollOutput.String())
	assert.Equal(t, expectedpollOptOutput, pollOptOutput.String())
	assert.Equal(t, expectedcommentOutput, commentOutput.String())
}
