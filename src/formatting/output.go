package formatting

import (
	"fmt"

	"github.com/fmenozzi/hn/src/api"
)

const itemBaseUrl = "https://news.ycombinator.com/item?id="

func JobOutput(job *api.Item, styled bool) string {
	score := *job.Score
	postUrl := fmt.Sprintf("%s%d", itemBaseUrl, job.Id)
	title := *job.Title

	if styled {
		return fmt.Sprintf("* [%4d pts] [       [HIRING](%s)] [%s](%s)\n", score, postUrl, title, postUrl)
	} else {
		return fmt.Sprintf("[%4d pts] [       HIRING] %s\n", score, title)
	}
}

func StoryOutput(story *api.Item, styled bool) string {
	score := *story.Score
	comments := *story.Descendants
	postUrl := fmt.Sprintf("%s%d", itemBaseUrl, story.Id)
	title := *story.Title
	url := postUrl

	if story.Url != nil && len(*story.Url) > 0 {
		url = *story.Url
	}

	if styled {
		return fmt.Sprintf("* [%4d pts] [[%4d comments](%s)] [%s](%s)\n", score, comments, postUrl, title, url)
	} else {
		return fmt.Sprintf("[%4d pts] [%4d comments] %s\n", score, comments, url)
	}
}

func PollOutput(poll *api.Item, styled bool) string {
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
