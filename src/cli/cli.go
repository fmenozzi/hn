package cli

import (
	"flag"
	"fmt"

	"github.com/fmenozzi/hn/src/api"
	"github.com/fmenozzi/hn/src/formatting"
)

const (
	Version = "0.1.0"
	usage   = `A simple commandline hacker news client.

Options:
    -h, --help      show this help message and exit
    -v, --version   show program version information and exit
    -l, --limit     max number of results to fetch (default: 30)
    -s, --style     output style, one of plain|markdown|md|csv (default: plain)
    -r, --ranking   ranking method
                    top|new|best for front page items (default: top)
                    date|popularity for search result items (default: popularity)
    -q, --query     search query
    -t, --tags      filter search results on specific tags (default: story)

Notes:
    The output for --style=csv is: id,type,by,timestamp,title,url,score,comments

    Search tags are ANDed by default but can be ORed if between parentheses. For
    example, "author_pg,(story,poll)" filters on "author_pg AND (type=story OR type=poll)".
    See https://hn.algolia.com/api for more.
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

	// Output formatting style.
	Style formatting.Style

	// Search query for searching items via the Algolia API.
	Query string

	// Comma-separated list of tags for filtering search results.
	Tags string
}

func ArgsFromCli() (Args, error) {
	var version bool
	var limit int
	var stylestr string
	var ranking string
	var query string
	var tags string

	flag.Usage = func() { fmt.Print(usage) }
	flag.BoolVar(&version, "v", false, "")
	flag.BoolVar(&version, "version", false, "")
	flag.StringVar(&ranking, "r", "", "")
	flag.StringVar(&ranking, "ranking", "", "")
	flag.IntVar(&limit, "l", 30, "")
	flag.IntVar(&limit, "limit", 30, "")
	flag.StringVar(&stylestr, "s", "plain", "")
	flag.StringVar(&stylestr, "style", "plain", "")
	flag.StringVar(&query, "q", "", "")
	flag.StringVar(&query, "query", "", "")
	flag.StringVar(&tags, "t", "", "")
	flag.StringVar(&tags, "tags", "", "")

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
			return Args{}, fmt.Errorf("invalid search ranking: %s\n", ranking)
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
			return Args{}, fmt.Errorf("invalid front page ranking: %s\n", ranking)
		}
	}

	if len(tags) > 0 {
		// --tags only allowed if --query was also specified.
		if len(query) == 0 {
			return Args{}, fmt.Errorf("tags invalid without query\n")
		}
	} else if len(query) > 0 {
		// Default to stories.
		tags = "story"
	}

	var style formatting.Style
	switch stylestr {
	case "":
		fallthrough
	case "plain":
		style = formatting.Plain
	case "md":
		fallthrough
	case "markdown":
		style = formatting.Markdown
	case "csv":
		style = formatting.Csv
	default:
		return Args{}, fmt.Errorf("invalid style: %s\n", style)
	}

	return Args{
		Version:              version,
		RankingFrontPage:     frontPageRanking,
		RankingSearchResults: searchResultsRanking,
		Limit:                limit,
		Style:                style,
		Query:                query,
		Tags:                 tags,
	}, nil
}
