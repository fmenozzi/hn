package main

import (
	"fmt"
	"strings"
	"sync"

	"github.com/fmenozzi/hn/src/api"
	"github.com/fmenozzi/hn/src/cli"
	"github.com/fmenozzi/hn/src/formatting"
)

func FetchFrontPageItems(args cli.Args) []api.Item {
	client := api.MakeClient()
	rankedStoriesIds, err := client.FetchRankedStoriesIds(args.Ranking, args.Limit)
	if err != nil {
		panic(fmt.Sprintf("Error: %s\n", err))
	}

	var itemsMap sync.Map
	var wg sync.WaitGroup
	for _, id := range rankedStoriesIds {
		wg.Add(1)
		go func(id api.ItemId) {
			defer wg.Done()
			item, err := client.FetchItem(id)
			if err != nil {
				panic(fmt.Sprintf("Error: %s\n", err))
			}
			itemsMap.Store(id, *item)
		}(id)
	}
	wg.Wait()

	rankedItems := make([]api.Item, len(rankedStoriesIds))
	for _, id := range rankedStoriesIds {
		mapitem, ok := itemsMap.Load(id)
		if !ok {
			panic(fmt.Sprintf("no item %d in items map", id))
		}
		item := mapitem.(api.Item)
		rankedItems = append(rankedItems, item)
	}

	return rankedItems
}

func DisplayItems(items []api.Item, styled bool) {
	var builder strings.Builder
	for _, item := range items {
		switch item.Type {
		case api.Job:
			builder.WriteString(formatting.JobOutput(&item, styled))
		case api.Story:
			builder.WriteString(formatting.StoryOutput(&item, styled))
		case api.Poll:
			builder.WriteString(formatting.PollOutput(&item, styled))
		}
	}
	fmt.Print(builder.String())
}

func main() {
	args, err := cli.ArgsFromCli()
	if err != nil {
		fmt.Printf("error: %s\n", err.Error())
	}

	DisplayItems(FetchFrontPageItems(args), args.Styled)
}
