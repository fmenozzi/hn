package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/fmenozzi/hn/src/api"
	"github.com/fmenozzi/hn/src/cli"
	"github.com/fmenozzi/hn/src/formatting"
)

var clock formatting.RealClock

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
	searchItemIds, err := client.Search(request)
	if err != nil {
		return nil, err
	}
	searchItems, err := client.FetchItems(searchItemIds)
	if err != nil {
		return nil, err
	}
	return searchItems, nil
}

func DisplayItems(items []api.Item, style formatting.Style) {
	var builder strings.Builder
	for _, item := range items {
		switch item.Type {
		case api.Job:
			builder.WriteString(formatting.JobOutput(&item, style, &clock))
		case api.Story:
			builder.WriteString(formatting.StoryOutput(&item, style, &clock))
		case api.Poll:
			builder.WriteString(formatting.PollOutput(&item, style, &clock))
		case api.PollOpt:
			builder.WriteString(formatting.PollOptOutput(&item, style, &clock))
		case api.Comment:
			builder.WriteString(formatting.CommentOutput(&item, style, &clock))
		default:
			panic(fmt.Sprintf("invalid item type %s", item.Type))
		}
	}
	fmt.Print(builder.String())
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
