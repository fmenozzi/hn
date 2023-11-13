package main

import (
	"flag"
	"fmt"
	"strings"
	"sync"

	"github.com/fmenozzi/hn/src/api"
	"github.com/fmenozzi/hn/src/formatting"
)

type ItemMap struct {
	Items map[api.ItemId]api.Item
	mutex sync.Mutex
}

func (im *ItemMap) Add(id api.ItemId, item api.Item) {
	im.mutex.Lock()
	im.Items[id] = item
	im.mutex.Unlock()
}

func Output(stories *api.Stories, items *ItemMap, styled bool) string {
	var builder strings.Builder
	for _, id := range stories.Ids {
		item := items.Items[id]
		switch item.Type {
		case api.Job:
			builder.WriteString(formatting.JobOutput(&item, styled))
		case api.Story:
			builder.WriteString(formatting.StoryOutput(&item, styled))
		case api.Poll:
			builder.WriteString(formatting.PollOutput(&item, styled))
		}
	}
	return builder.String()
}

func main() {
	var limit int
	var styled bool
	var rankingstr string
	flag.IntVar(&limit, "l", 30, "Number of stories to fetch")
	flag.BoolVar(&styled, "s", false, "Whether to style output for mdcat")
	flag.StringVar(&rankingstr, "r", "top", "Ranking method (one of `top`, `best`, `new`)")
	flag.Parse()

	var ranking api.StoriesRanking
	switch rankingstr {
	case "top":
		ranking = api.Top
	case "best":
		ranking = api.Best
	case "new":
		ranking = api.New
	default:
		panic(fmt.Sprintf("invalid ranking option: %s", rankingstr))
	}

	client := api.MakeClient()
	stories, err := client.FetchStories(ranking, limit)
	if err != nil {
		panic(fmt.Sprintf("Error: %s\n", err))
	}

	items := ItemMap{Items: make(map[api.ItemId]api.Item)}
	var wg sync.WaitGroup
	for _, id := range stories.Ids {
		wg.Add(1)
		go func(id api.ItemId) {
			defer wg.Done()
			item, err := client.FetchItem(id)
			if err != nil {
				panic(fmt.Sprintf("Error: %s\n", err))
			}
			items.Add(id, *item)
		}(id)
	}
	wg.Wait()

	fmt.Print(Output(stories, &items, styled))
}
