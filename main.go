package main

import (
	"flag"
	"fmt"
	"strings"
	"sync"
)

type ItemMap struct {
	Items map[ItemId]Item
	mutex sync.Mutex
}

func (im *ItemMap) Add(id ItemId, item Item) {
	im.mutex.Lock()
	im.Items[id] = item
	im.mutex.Unlock()
}

func Output(stories *Stories, items *ItemMap, styled bool) string {
	var builder strings.Builder
	for _, id := range stories.Ids {
		item := items.Items[id]
		switch item.Type {
		case Job:
			builder.WriteString(JobOutput(&item, styled))
		case Story:
			builder.WriteString(StoryOutput(&item, styled))
		case Poll:
			builder.WriteString(PollOutput(&item, styled))
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

	var ranking StoriesRanking
	switch rankingstr {
	case "top":
		ranking = Top
	case "best":
		ranking = Best
	case "new":
		ranking = New
	default:
		panic(fmt.Sprintf("invalid ranking option: %s", rankingstr))
	}

	client := MakeClient()
	stories, err := client.FetchStories(ranking, limit)
	if err != nil {
		panic(fmt.Sprintf("Error: %s\n", err))
	}

	items := ItemMap{Items: make(map[ItemId]Item)}
	var wg sync.WaitGroup
	for _, id := range stories.Ids {
		wg.Add(1)
		go func(id ItemId) {
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
