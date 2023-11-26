package main

import (
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/fmenozzi/hn/src/api"
	"github.com/fmenozzi/hn/src/cli"
	"github.com/fmenozzi/hn/src/formatting"
)

func FetchFrontPageItems(ranking api.FrontPageItemsRanking, limit int) []api.Item {
	client := api.MakeClient()
	rankedStoriesIds, err := client.FetchRankedStoriesIds(ranking, limit)
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
		os.Exit(1)
	}

	if args.Version {
		fmt.Println(cli.Version)
		os.Exit(0)
	}

	if len(args.Query) != 0 {
		if args.RankingSearchResults == nil {
			fmt.Printf("error: invalid search result ranking: %v\n", args.RankingRawString)
			os.Exit(1)
		}
		fmt.Println("error: search is unimplemented")
		os.Exit(1)
	} else {
		if args.RankingFrontPage == nil {
			fmt.Printf("error: invalid front page ranking: %v\n", args.RankingRawString)
			os.Exit(1)
		}
		DisplayItems(FetchFrontPageItems(*args.RankingFrontPage, args.Limit), args.Styled)
	}
}
