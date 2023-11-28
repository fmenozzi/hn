package cli

import (
	"flag"
	"fmt"

	"github.com/fmenozzi/hn/src/api"
)

const (
	Version = "0.1.0"
	usage   = `A simple commandline hacker news client.

Options:
    -h, --help      show this help message and exit
    -v, --version   show program version information and exit
    -l, --limit     max number of results to fetch (default: 30)
    -s, --styled    if true, format output for piping into mdcat
    -r, --ranking   ranking method
                    top|new|best for front page items (default: top)
                    date|popularity for search result items (default: popularity)
    -q, --query     search query
`
)

type Args struct {
	// If true, version information was requested.
	Version bool

	// Ranking method for front page items.
	RankingFrontPage *api.FrontPageItemsRanking

	// Ranking method for search result items.
	RankingSearchResults *api.SearchItemsRanking

	// Max number of stories to fetch.
	Limit int

	// If true, format output for piping into `mdcat`.
	Styled bool

	// Search query for searching items via the Algolia API.
	Query string
}

func ArgsFromCli() (Args, error) {
	var version bool
	var limit int
	var styled bool
	var ranking string
	var query string

	flag.Usage = func() { fmt.Print(usage) }
	flag.BoolVar(&version, "v", false, "")
	flag.BoolVar(&version, "version", false, "")
	flag.StringVar(&ranking, "r", "", "")
	flag.StringVar(&ranking, "ranking", "", "")
	flag.IntVar(&limit, "l", 30, "")
	flag.IntVar(&limit, "limit", 30, "")
	flag.BoolVar(&styled, "s", false, "")
	flag.BoolVar(&styled, "styled", false, "")
	flag.StringVar(&query, "q", "", "")
	flag.StringVar(&query, "query", "", "")
	flag.Parse()

	var frontPageRanking *api.FrontPageItemsRanking
	var searchResultsRanking *api.SearchItemsRanking
	if len(query) > 0 {
		// --query was passed, interpret --ranking for search.
		switch ranking {
		case "":
			fallthrough
		case "popularity":
			searchResultsRanking = api.Popularity.ToPointer()
		case "date":
			searchResultsRanking = api.Date.ToPointer()
		default:
			return Args{}, fmt.Errorf("invalid search ranking: %s", ranking)
		}
	} else {
		// --query was not passed, interpret --ranking for front page.
		switch ranking {
		case "":
			fallthrough
		case "top":
			frontPageRanking = api.Top.ToPointer()
		case "new":
			frontPageRanking = api.New.ToPointer()
		case "best":
			frontPageRanking = api.Best.ToPointer()
		default:
			return Args{}, fmt.Errorf("invalid front page ranking: %s", ranking)
		}
	}

	return Args{
		Version:              version,
		RankingFrontPage:     frontPageRanking,
		RankingSearchResults: searchResultsRanking,
		Limit:                limit,
		Styled:               styled,
		Query:                query,
	}, nil
}
