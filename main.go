package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/fmenozzi/hn/src/api"
	"github.com/fmenozzi/hn/src/cli"
	"github.com/fmenozzi/hn/src/formatting"
)

func FetchFrontPageItems(ranking api.FrontPageItemsRanking, limit int) []api.Item {
	client := api.MakeProdClient()
	frontPageItemIds, err := client.FetchFrontPageItemIds(ranking, limit)
	if err != nil {
		panic(fmt.Sprintf("error: %s\n", err))
	}
	frontPageItems, err := client.FetchItems(frontPageItemIds)
	if err != nil {
		panic(fmt.Sprintf("error: %s\n", err))
	}
	return frontPageItems
}

func FetchSearchItems(request api.SearchRequest) []api.Item {
	client := api.MakeProdClient()
	searchItemIds, err := client.Search(request)
	if err != nil {
		panic(fmt.Sprintf("error: %s\n", err))
	}
	searchItems, err := client.FetchItems(searchItemIds)
	if err != nil {
		panic(fmt.Sprintf("error: %s\n", err))
	}
	return searchItems
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
		request := api.SearchRequest{
			Query:   args.Query,
			Tags:    []string{"story"},
			Ranking: *args.RankingSearchResults,
			Limit:   args.Limit,
		}
		DisplayItems(FetchSearchItems(request), args.Styled)
	} else {
		DisplayItems(FetchFrontPageItems(*args.RankingFrontPage, args.Limit), args.Styled)
	}
}
