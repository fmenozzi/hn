package main

import "fmt"

const itemBaseUrl = "https://news.ycombinator.com/item?id="

func JobOutput(job *Item, styled bool) string {
	score := *job.Score
	postUrl := fmt.Sprintf("%s%d", itemBaseUrl, job.Id)
	title := *job.Title

	if styled {
		return fmt.Sprintf("* [%4d pts] [       [HIRING](%s)] [%s](%s)\n", score, postUrl, title, postUrl)
	} else {
		return fmt.Sprintf("[%4d pts] [       HIRING] %s\n", score, title)
	}
}

func StoryOutput(story *Item, styled bool) string {
	score := *story.Score
	comments := *story.Descendants
	postUrl := fmt.Sprintf("%s%d", itemBaseUrl, story.Id)
	title := *story.Title
	url := postUrl

	if story.Url != nil {
		url = *story.Url
	}

	if styled {
		return fmt.Sprintf("* [%4d pts] [[%4d comments](%s)] [%s](%s)\n", score, comments, postUrl, title, url)
	} else {
		return fmt.Sprintf("[%4d pts] [%4d comments] %s\n", score, comments, url)
	}
}

func PollOutput(poll *Item, styled bool) string {
	score := *poll.Score
	comments := *poll.Descendants
	postUrl := fmt.Sprintf("%s%d", itemBaseUrl, poll.Id)
	title := *poll.Title

	if styled {
		return fmt.Sprintf("* [%4d pts] [[%4d comments](%s) [%s](%s)\n", score, comments, postUrl, title, postUrl)
	} else {
		return fmt.Sprintf("[%4d pts] [%4d comments] %s\n", score, comments, postUrl)
	}
}
