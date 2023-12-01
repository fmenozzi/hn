package formatting

import (
	"fmt"
	"time"

	"github.com/fmenozzi/hn/src/api"
)

const (
	itemBaseUrl = "https://news.ycombinator.com/item?id="
	userBaseUrl = "https://news.ycombinator.com/user?id="
)

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
		return fmt.Sprintf("HIRING: %s\n└─── %d pts %s\n", postUrl, score, time)
	case Markdown:
		return fmt.Sprintf("* **[HIRING: %s](%s)**\n* └─── %d pts %s\n", title, postUrl, score, time)
	case Csv:
		return fmt.Sprintf("%d,job,%s,%d,\"%s\",%s,%d,0\n", job.Id, *job.By, *job.Time, title, postUrl, score)
	default:
		panic(fmt.Sprintf("invalid style: %s\n", style))
	}
}

func StoryOutput(story *api.Item, style Style, clock Clock) string {
	by := *story.By
	score := *story.Score
	comments := *story.Descendants
	postUrl := fmt.Sprintf("%s%d", itemBaseUrl, story.Id)
	byUrl := fmt.Sprintf("%s%s", userBaseUrl, by)
	title := *story.Title
	url := postUrl
	time := GetRelativeTime(clock, time.Unix(*story.Time, 0))

	if story.Url != nil && len(*story.Url) > 0 {
		url = *story.Url
	}

	switch style {
	case Plain:
		return fmt.Sprintf("%s\n└─── %d pts by %s %s | %d comments\n", url, score, by, time, comments)
	case Markdown:
		return fmt.Sprintf("* **[%s](%s)**\n* └─── %d pts by [%s](%s) %s | [%d comments](%s)\n", title, url, score, by, byUrl, time, comments, postUrl)
	case Csv:
		return fmt.Sprintf("%d,story,%s,%d,\"%s\",%s,%d,%d\n", story.Id, *story.By, *story.Time, title, url, score, comments)
	default:
		panic(fmt.Sprintf("invalid style: %s\n", style))
	}
}

func PollOutput(poll *api.Item, style Style, clock Clock) string {
	by := *poll.By
	score := *poll.Score
	comments := *poll.Descendants
	postUrl := fmt.Sprintf("%s%d", itemBaseUrl, poll.Id)
	byUrl := fmt.Sprintf("%s%s", userBaseUrl, by)
	title := *poll.Title
	time := GetRelativeTime(clock, time.Unix(*poll.Time, 0))

	switch style {
	case Plain:
		return fmt.Sprintf("%s\n└─── %d pts by %s %s | %d comments\n", postUrl, score, by, time, comments)
	case Markdown:
		return fmt.Sprintf("* **[%s](%s)**\n* └─── %d pts by [%s](%s) %s | [%d comments](%s)\n", title, postUrl, score, by, byUrl, time, comments, postUrl)
	case Csv:
		return fmt.Sprintf("%d,poll,%s,%d,\"%s\",%s,%d,%d\n", poll.Id, *poll.By, *poll.Time, title, postUrl, score, comments)
	default:
		panic(fmt.Sprintf("invalid style: %s\n", style))
	}
}

func PollOptOutput(pollopt *api.Item, style Style, clock Clock) string {
	by := *pollopt.By
	postUrl := fmt.Sprintf("%s%d", itemBaseUrl, pollopt.Id)
	byUrl := fmt.Sprintf("%s%s", userBaseUrl, by)
	text := *pollopt.Text
	time := GetRelativeTime(clock, time.Unix(*pollopt.Time, 0))
	score := *pollopt.Score

	// TODO: This is quite hacky but will do for now while we overhaul the UI.
	if len(text) > 70 {
		text = fmt.Sprintf("%s...", text[:70])
	}

	switch style {
	case Plain:
		return fmt.Sprintf("%s\n└─── %d pts by %s %s\n", text, score, by, time)
	case Markdown:
		return fmt.Sprintf("* **[%s](%s)**\n* └─── %d pts by [%s](%s) %s\n", text, postUrl, score, by, byUrl, time)
	case Csv:
		// TODO: commas in text are not always respected, find a way around that.
		return fmt.Sprintf("%d,pollopt,%s,%d,\"%s\",%s,%d,0\n", pollopt.Id, *pollopt.By, *pollopt.Time, *pollopt.Text, postUrl, score)
	default:
		panic(fmt.Sprintf("invalid style: %s\n", style))
	}
}

func CommentOutput(comment *api.Item, style Style, clock Clock) string {
	by := *comment.By
	postUrl := fmt.Sprintf("%s%d", itemBaseUrl, comment.Id)
	byUrl := fmt.Sprintf("%s%s", userBaseUrl, by)
	text := *comment.Text
	time := GetRelativeTime(clock, time.Unix(*comment.Time, 0))
	score := 0
	comments := len(comment.Kids)

	// TODO: This is quite hacky but will do for now while we overhaul the UI.
	if len(text) > 70 {
		text = fmt.Sprintf("%s...", text[:70])
	}

	switch style {
	case Plain:
		return fmt.Sprintf("%s\n└─── by %s %s | %d replies\n", text, by, time, comments)
	case Markdown:
		return fmt.Sprintf("* *[%s](%s)*\n* └─── by [%s](%s) %s | [%d replies](%s)\n", text, postUrl, by, byUrl, time, comments, postUrl)
	case Csv:
		// TODO: commas in text are not always respected, find a way around that.
		return fmt.Sprintf("%d,comment,%s,%d,\"%s\",%s,%d,%d\n", comment.Id, *comment.By, *comment.Time, *comment.Text, postUrl, score, comments)
	default:
		panic(fmt.Sprintf("invalid style: %s\n", style))
	}
}
