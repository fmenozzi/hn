package main

import (
	"flag"
	"fmt"
	"strings"
	"sync"

	"github.com/fmenozzi/hn/src/api"
	"github.com/fmenozzi/hn/src/formatting"
)

func Output(stories *api.Stories, items *sync.Map, styled bool) string {
	var builder strings.Builder
	for _, id := range stories.Ids {
		mapitem, _ := items.Load(id)
		item := mapitem.(api.Item)
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

	var items sync.Map
	var wg sync.WaitGroup
	for _, id := range stories.Ids {
		wg.Add(1)
		go func(id api.ItemId) {
			defer wg.Done()
			item, err := client.FetchItem(id)
			if err != nil {
				panic(fmt.Sprintf("Error: %s\n", err))
			}
			items.Store(id, *item)
		}(id)
	}
	wg.Wait()

	fmt.Print(Output(stories, &items, styled))
}
