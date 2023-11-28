package formatting

import (
	"fmt"

	"github.com/fmenozzi/hn/src/api"
)

const itemBaseUrl = "https://news.ycombinator.com/item?id="

type Style string

const (
	Plain    Style = "plain"
	Markdown       = "markdown"
	Csv            = "csv"
)

func JobOutput(job *api.Item, style Style) string {
	score := *job.Score
	postUrl := fmt.Sprintf("%s%d", itemBaseUrl, job.Id)
	title := *job.Title

	switch style {
	case Plain:
		return fmt.Sprintf("[%4d pts] [       HIRING] %s\n", score, title)
	case Markdown:
		return fmt.Sprintf("* [%4d pts] [       [HIRING](%s)] [%s](%s)\n", score, postUrl, title, postUrl)
	case Csv:
		return fmt.Sprintf("%d,job,%s,\"%s\",%s,%d,0\n", job.Id, *job.By, title, postUrl, score)
	default:
		panic(fmt.Sprintf("invalid style: %s\n", style))
	}
}

func StoryOutput(story *api.Item, style Style) string {
	score := *story.Score
	comments := *story.Descendants
	postUrl := fmt.Sprintf("%s%d", itemBaseUrl, story.Id)
	title := *story.Title
	url := postUrl

	if story.Url != nil && len(*story.Url) > 0 {
		url = *story.Url
	}

	switch style {
	case Plain:
		return fmt.Sprintf("[%4d pts] [%4d comments] %s\n", score, comments, url)
	case Markdown:
		return fmt.Sprintf("* [%4d pts] [[%4d comments](%s)] [%s](%s)\n", score, comments, postUrl, title, url)
	case Csv:
		return fmt.Sprintf("%d,story,%s,\"%s\",%s,%d,%d\n", story.Id, *story.By, title, url, score, comments)
	default:
		panic(fmt.Sprintf("invalid style: %s\n", style))
	}
}

func PollOutput(poll *api.Item, style Style) string {
	score := *poll.Score
	comments := *poll.Descendants
	postUrl := fmt.Sprintf("%s%d", itemBaseUrl, poll.Id)
	title := *poll.Title

	switch style {
	case Plain:
		return fmt.Sprintf("[%4d pts] [%4d comments] %s\n", score, comments, postUrl)
	case Markdown:
		return fmt.Sprintf("* [%4d pts] [[%4d comments](%s) [%s](%s)\n", score, comments, postUrl, title, postUrl)
	case Csv:
		return fmt.Sprintf("%d,poll,%s,\"%s\",%s,%d,%d\n", poll.Id, *poll.By, title, postUrl, score, comments)
	default:
		panic(fmt.Sprintf("invalid style: %s\n", style))
	}
}
