package formatting

import (
	"fmt"
	"time"

	"github.com/fmenozzi/hn/src/api"
)

const itemBaseUrl = "https://news.ycombinator.com/item?id="

type Style string

const (
	Plain    Style = "plain"
	Markdown       = "markdown"
	Csv            = "csv"
)

func JobOutput(job *api.Item, style Style, clock Clock) string {
	score := *job.Score
	postUrl := fmt.Sprintf("%s%d", itemBaseUrl, job.Id)
	title := *job.Title
	time := GetRelativeTime(clock, time.Unix(*job.Time, 0))

	switch style {
	case Plain:
		return fmt.Sprintf("[%4d pts] [%17s] [       HIRING] %s\n", score, time, title)
	case Markdown:
		return fmt.Sprintf("* [%4d pts] [%17s] [       [HIRING](%s)] [%s](%s)\n", score, time, postUrl, title, postUrl)
	case Csv:
		return fmt.Sprintf("%d,job,%s,%d,\"%s\",%s,%d,0\n", job.Id, *job.By, *job.Time, title, postUrl, score)
	default:
		panic(fmt.Sprintf("invalid style: %s\n", style))
	}
}

func StoryOutput(story *api.Item, style Style, clock Clock) string {
	score := *story.Score
	comments := *story.Descendants
	postUrl := fmt.Sprintf("%s%d", itemBaseUrl, story.Id)
	title := *story.Title
	url := postUrl
	time := GetRelativeTime(clock, time.Unix(*story.Time, 0))

	if story.Url != nil && len(*story.Url) > 0 {
		url = *story.Url
	}

	switch style {
	case Plain:
		return fmt.Sprintf("[%4d pts] [%17s] [%4d comments] %s\n", score, time, comments, url)
	case Markdown:
		return fmt.Sprintf("* [%4d pts] [%17s] [[%4d comments](%s)] [%s](%s)\n", score, time, comments, postUrl, title, url)
	case Csv:
		return fmt.Sprintf("%d,story,%s,%d,\"%s\",%s,%d,%d\n", story.Id, *story.By, *story.Time, title, url, score, comments)
	default:
		panic(fmt.Sprintf("invalid style: %s\n", style))
	}
}

func PollOutput(poll *api.Item, style Style, clock Clock) string {
	score := *poll.Score
	comments := *poll.Descendants
	postUrl := fmt.Sprintf("%s%d", itemBaseUrl, poll.Id)
	title := *poll.Title
	time := GetRelativeTime(clock, time.Unix(*poll.Time, 0))

	switch style {
	case Plain:
		return fmt.Sprintf("[%4d pts] [%17s] [%4d comments] %s\n", score, time, comments, postUrl)
	case Markdown:
		return fmt.Sprintf("* [%4d pts] [%17s] [[%4d comments](%s)] [%s](%s)\n", score, time, comments, postUrl, title, postUrl)
	case Csv:
		return fmt.Sprintf("%d,poll,%s,%d,\"%s\",%s,%d,%d\n", poll.Id, *poll.By, *poll.Time, title, postUrl, score, comments)
	default:
		panic(fmt.Sprintf("invalid style: %s\n", style))
	}
}
