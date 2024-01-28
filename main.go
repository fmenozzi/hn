package main

import (
	"fmt"
	"os"

	"github.com/fmenozzi/hn/src/api"
	"github.com/fmenozzi/hn/src/cli"
	"github.com/fmenozzi/hn/src/formatting"
)

func FetchFrontPageItems(ranking api.FrontPageItemsRanking, limit int) ([]api.Item, error) {
	client := api.MakeProdClient()
	frontPageItemIds, err := client.FetchFrontPageItemIds(ranking, limit)
	if err != nil {
		return nil, err
	}
	frontPageItems, err := client.FetchItems(frontPageItemIds)
	if err != nil {
		return nil, err
	}
	return frontPageItems, nil
}

func FetchSearchItems(request api.SearchRequest) ([]api.Item, error) {
	client := api.MakeProdClient()
	searchResponse, err := client.Search(request)
	if err != nil {
		return nil, err
	}
	searchItemIds := make([]api.ItemId, len(searchResponse.Results))
	for i, result := range searchResponse.Results {
		searchItemIds[i] = result.Id
	}
	searchItems, err := client.FetchItems(searchItemIds)
	if err != nil {
		return nil, err
	}
	return searchItems, nil
}

func DisplayItems(items []api.Item, style formatting.Style) {
	var clock formatting.RealClock
	switch style {
	case formatting.Plain:
		fmt.Print(formatting.DisplayPlain(items, &clock))
	case formatting.Markdown:
		fmt.Print(formatting.DisplayMarkdown(items, &clock))
	case formatting.Json:
		fmt.Print(formatting.DisplayJson(items))
	default:
		panic(fmt.Sprintf("invalid style: %s\n", style))
	}
}

func main() {
	args, err := cli.ArgsFromCli()
	if err != nil {
		fmt.Fprint(os.Stderr, err.Error())
		os.Exit(1)
	}

	if args.Version {
		fmt.Println(cli.Version)
		os.Exit(0)
	}

	if len(args.Query) != 0 {
		searchItems, err := FetchSearchItems(api.SearchRequest{
			Query:   args.Query,
			Tags:    args.Tags,
			Ranking: *args.RankingSearchResults,
			Limit:   args.Limit,
		})
		if err != nil {
			fmt.Fprint(os.Stderr, err.Error())
			os.Exit(1)
		}
		DisplayItems(searchItems, args.Style)
	} else {
		frontPageItems, err := FetchFrontPageItems(*args.RankingFrontPage, args.Limit)
		if err != nil {
			fmt.Fprint(os.Stderr, err.Error())
			os.Exit(1)
		}
		DisplayItems(frontPageItems, args.Style)
	}
}
