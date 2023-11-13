package cli

import (
	"flag"
	"fmt"

	"github.com/fmenozzi/hn/src/api"
)

type Args struct {
	// Ranking method for fetched stories.
	Ranking api.StoriesRanking

	// Max number of stories to fetch.
	Limit int

	// If true, format output for piping into `mdcat`.
	Styled bool
}

func ArgsFromCli() (Args, error) {
	var limit int
	var styled bool
	var rankingstr string

	flag.IntVar(&limit, "l", 30, "Max number of stories to fetch")
	flag.BoolVar(&styled, "s", false, "If true, format output for piping into `mdcat`")
	flag.StringVar(&rankingstr, "r", "top", "Ranking method (one of `top`, `new`, `best`)")
	flag.Parse()

	var ranking api.StoriesRanking
	switch rankingstr {
	case "top":
		ranking = api.Top
	case "new":
		ranking = api.New
	case "best":
		ranking = api.Best
	default:
		return Args{}, fmt.Errorf("invalid ranking: %s", rankingstr)
	}

	return Args{
		Ranking: ranking,
		Limit:   limit,
		Styled:  styled,
	}, nil
}
