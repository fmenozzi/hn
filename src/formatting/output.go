package formatting

import (
	"encoding/json"
	"fmt"
	"strings"
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
	Json           = "json"
)

func DisplayPlain(items []api.Item, clock Clock) string {
	return display(items, Plain, clock)
}

func DisplayMarkdown(items []api.Item, clock Clock) string {
	return display(items, Markdown, clock)
}

func DisplayJson(items []api.Item) string {
	// Note that we do not do any post-fetch formatting of output here: this is
	// basically a one-to-one mapping of the Item data as fetched. For example,
	// we do not fall back to post urls if the item does not have a url, and the
	// time is represented as the original Unix timestamp instead of the human-
	// readable relative time.
	bytes, err := json.MarshalIndent(items, "", "\t")
	if err != nil {
		panic(fmt.Sprintf("error formatting items as json: %s", err.Error()))
	}
	return string(bytes)
}

func display(items []api.Item, style Style, clock Clock) string {
	var builder strings.Builder
	for _, item := range items {
		switch item.Type {
		case api.Job:
			builder.WriteString(jobOutput(&item, style, clock))
		case api.Story:
			builder.WriteString(storyOutput(&item, style, clock))
		case api.Poll:
			builder.WriteString(pollOutput(&item, style, clock))
		case api.PollOpt:
			builder.WriteString(pollOptOutput(&item, style, clock))
		case api.Comment:
			builder.WriteString(commentOutput(&item, style, clock))
		default:
			panic(fmt.Sprintf("invalid item type %s", item.Type))
		}
	}
	return builder.String()
}

func jobOutput(job *api.Item, style Style, clock Clock) string {
	score := *job.Score
	postUrl := fmt.Sprintf("%s%d", itemBaseUrl, job.Id)
	title := *job.Title
	time := GetRelativeTime(clock, time.Unix(*job.Time, 0))
	ptsstr := "pts"

	if score == 1 {
		ptsstr = "pt"
	}

	switch style {
	case Plain:
		return fmt.Sprintf("HIRING: %s\n└─── %d %s %s\n", postUrl, score, ptsstr, time)
	case Markdown:
		return fmt.Sprintf("* **[HIRING: %s](%s)**\n* └─── %d %s %s\n", title, postUrl, score, ptsstr, time)
	default:
		panic(fmt.Sprintf("invalid style: %s\n", style))
	}
}

func storyOutput(story *api.Item, style Style, clock Clock) string {
	by := *story.By
	score := *story.Score
	comments := *story.Descendants
	postUrl := fmt.Sprintf("%s%d", itemBaseUrl, story.Id)
	byUrl := fmt.Sprintf("%s%s", userBaseUrl, by)
	title := *story.Title
	url := postUrl
	time := GetRelativeTime(clock, time.Unix(*story.Time, 0))
	ptsstr := "pts"
	commentsstr := "comments"

	if story.Url != nil && len(*story.Url) > 0 {
		url = *story.Url
	}

	if score == 1 {
		ptsstr = "pt"
	}
	if comments == 1 {
		commentsstr = "comment"
	}

	switch style {
	case Plain:
		return fmt.Sprintf("%s\n└─── %d %s by %s %s | %d %s\n", url, score, ptsstr, by, time, comments, commentsstr)
	case Markdown:
		return fmt.Sprintf("* **[%s](%s)**\n* └─── %d %s by [%s](%s) %s | [%d %s](%s)\n", title, url, score, ptsstr, by, byUrl, time, comments, commentsstr, postUrl)
	default:
		panic(fmt.Sprintf("invalid style: %s\n", style))
	}
}

func pollOutput(poll *api.Item, style Style, clock Clock) string {
	by := *poll.By
	score := *poll.Score
	comments := *poll.Descendants
	postUrl := fmt.Sprintf("%s%d", itemBaseUrl, poll.Id)
	byUrl := fmt.Sprintf("%s%s", userBaseUrl, by)
	title := *poll.Title
	time := GetRelativeTime(clock, time.Unix(*poll.Time, 0))
	ptsstr := "pts"
	commentsstr := "comments"

	if score == 1 {
		ptsstr = "pt"
	}
	if comments == 1 {
		commentsstr = "comment"
	}

	switch style {
	case Plain:
		return fmt.Sprintf("%s\n└─── %d %s by %s %s | %d %s\n", postUrl, score, ptsstr, by, time, comments, commentsstr)
	case Markdown:
		return fmt.Sprintf("* **[%s](%s)**\n* └─── %d %s by [%s](%s) %s | [%d %s](%s)\n", title, postUrl, score, ptsstr, by, byUrl, time, comments, commentsstr, postUrl)
	default:
		panic(fmt.Sprintf("invalid style: %s\n", style))
	}
}

func pollOptOutput(pollopt *api.Item, style Style, clock Clock) string {
	by := *pollopt.By
	postUrl := fmt.Sprintf("%s%d", itemBaseUrl, pollopt.Id)
	byUrl := fmt.Sprintf("%s%s", userBaseUrl, by)
	text := *pollopt.Text
	time := GetRelativeTime(clock, time.Unix(*pollopt.Time, 0))
	score := *pollopt.Score
	ptsstr := "pts"

	// TODO: This is quite hacky but will do for now while we overhaul the UI.
	if len(text) > 70 {
		text = fmt.Sprintf("%s...", text[:70])
	}

	if score == 1 {
		ptsstr = "pt"
	}

	switch style {
	case Plain:
		return fmt.Sprintf("%s\n└─── %d %s by %s %s\n", text, score, ptsstr, by, time)
	case Markdown:
		return fmt.Sprintf("* **[%s](%s)**\n* └─── %d %s by [%s](%s) %s\n", text, postUrl, score, ptsstr, by, byUrl, time)
	default:
		panic(fmt.Sprintf("invalid style: %s\n", style))
	}
}

func commentOutput(comment *api.Item, style Style, clock Clock) string {
	by := *comment.By
	postUrl := fmt.Sprintf("%s%d", itemBaseUrl, comment.Id)
	byUrl := fmt.Sprintf("%s%s", userBaseUrl, by)
	text := *comment.Text
	time := GetRelativeTime(clock, time.Unix(*comment.Time, 0))
	comments := len(comment.Kids)
	repliesstr := "replies"

	// TODO: This is quite hacky but will do for now while we overhaul the UI.
	if len(text) > 70 {
		text = fmt.Sprintf("%s...", text[:70])
	}

	if comments == 1 {
		repliesstr = "reply"
	}

	switch style {
	case Plain:
		return fmt.Sprintf("%s\n└─── by %s %s | %d %s\n", text, by, time, comments, repliesstr)
	case Markdown:
		return fmt.Sprintf("* *[%s](%s)*\n* └─── by [%s](%s) %s | [%d %s](%s)\n", text, postUrl, by, byUrl, time, comments, repliesstr, postUrl)
	default:
		panic(fmt.Sprintf("invalid style: %s\n", style))
	}
}
