package formatting

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
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
	Csv            = "csv"
)

func WritePlain(items []api.Item, clock Clock, w io.Writer) {
	writeItems(items, Plain, clock, w)
}

func WriteMarkdown(items []api.Item, clock Clock, w io.Writer) {
	writeItems(items, Markdown, clock, w)
}

func WriteJson(items []api.Item, w io.Writer) {
	// Note that we do not do any post-fetch formatting of output here: this is
	// basically a one-to-one mapping of the Item data as fetched. For example,
	// we do not fall back to post urls if the item does not have a url, and the
	// time is represented as the original Unix timestamp instead of the human-
	// readable relative time.
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "\t")
	if err := encoder.Encode(items); err != nil {
		panic(fmt.Sprintf("error formatting items as json: %s", err.Error()))
	}
}

func WriteCsv(items []api.Item, w io.Writer) {
	records := make([][]string, len(items))
	for i, item := range items {
		title := ""
		if item.Title != nil {
			title = *item.Title
		} else if item.Text != nil {
			title = *item.Text
		}
		url := fmt.Sprintf("%s%d", itemBaseUrl, item.Id)
		if item.Url != nil {
			url = *item.Url
		}
		score := int32(0)
		if item.Score != nil {
			score = *item.Score
		}
		comments := int32(0)
		if item.Descendants != nil {
			comments = *item.Descendants
		} else if item.Kids != nil {
			comments = int32(len(item.Kids))
		}
		records[i] = []string{
			strconv.FormatInt(int64(item.Id), 10),
			string(item.Type),
			*item.By,
			strconv.FormatInt(*item.Time, 10),
			title,
			url,
			strconv.FormatInt(int64(score), 10),
			strconv.FormatInt(int64(comments), 10),
		}
	}
	csvWriter := csv.NewWriter(w)
	if err := csvWriter.WriteAll(records); err != nil {
		panic(fmt.Sprintf("error formatting items as csv: %s", err.Error()))
	}
}

func writeItems(items []api.Item, style Style, clock Clock, w io.Writer) {
	for _, item := range items {
		switch item.Type {
		case api.Job:
			writeJobItem(&item, style, clock, w)
		case api.Story:
			writeStoryItem(&item, style, clock, w)
		case api.Poll:
			writePollItem(&item, style, clock, w)
		case api.PollOpt:
			writePollOptItem(&item, style, clock, w)
		case api.Comment:
			writeCommentItem(&item, style, clock, w)
		default:
			panic(fmt.Sprintf("invalid item type %s", item.Type))
		}
	}
}

func writeJobItem(job *api.Item, style Style, clock Clock, w io.Writer) {
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
		fmt.Fprintf(w, "HIRING: %s\n└─── %d %s %s\n", postUrl, score, ptsstr, time)
	case Markdown:
		fmt.Fprintf(w, "* **[HIRING: %s](%s)**\n* └─── %d %s %s\n", title, postUrl, score, ptsstr, time)
	default:
		panic(fmt.Sprintf("invalid style: %s\n", style))
	}
}

func writeStoryItem(story *api.Item, style Style, clock Clock, w io.Writer) {
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
		fmt.Fprintf(w, "%s\n└─── %d %s by %s %s | %d %s\n", url, score, ptsstr, by, time, comments, commentsstr)
	case Markdown:
		fmt.Fprintf(w, "* **[%s](%s)**\n* └─── %d %s by [%s](%s) %s | [%d %s](%s)\n", title, url, score, ptsstr, by, byUrl, time, comments, commentsstr, postUrl)
	default:
		panic(fmt.Sprintf("invalid style: %s\n", style))
	}
}

func writePollItem(poll *api.Item, style Style, clock Clock, w io.Writer) {
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
		fmt.Fprintf(w, "%s\n└─── %d %s by %s %s | %d %s\n", postUrl, score, ptsstr, by, time, comments, commentsstr)
	case Markdown:
		fmt.Fprintf(w, "* **[%s](%s)**\n* └─── %d %s by [%s](%s) %s | [%d %s](%s)\n", title, postUrl, score, ptsstr, by, byUrl, time, comments, commentsstr, postUrl)
	default:
		panic(fmt.Sprintf("invalid style: %s\n", style))
	}
}

func writePollOptItem(pollopt *api.Item, style Style, clock Clock, w io.Writer) {
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
		fmt.Fprintf(w, "%s\n└─── %d %s by %s %s\n", text, score, ptsstr, by, time)
	case Markdown:
		fmt.Fprintf(w, "* **[%s](%s)**\n* └─── %d %s by [%s](%s) %s\n", text, postUrl, score, ptsstr, by, byUrl, time)
	default:
		panic(fmt.Sprintf("invalid style: %s\n", style))
	}
}

func writeCommentItem(comment *api.Item, style Style, clock Clock, w io.Writer) {
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
		fmt.Fprintf(w, "%s\n└─── by %s %s | %d %s\n", text, by, time, comments, repliesstr)
	case Markdown:
		fmt.Fprintf(w, "* *[%s](%s)*\n* └─── by [%s](%s) %s | [%d %s](%s)\n", text, postUrl, by, byUrl, time, comments, repliesstr, postUrl)
	default:
		panic(fmt.Sprintf("invalid style: %s\n", style))
	}
}
