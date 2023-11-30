package api

type ItemId = int32

type ItemType string

const (
	Job     ItemType = "job"
	Story            = "story"
	Comment          = "comment"
	Poll             = "poll"
	PollOpt          = "pollopt"
)

type Item struct {
	// The item's unique id.
	Id ItemId `json:"id"`

	// `true` if the item is deleted.
	Deleted *bool `json:"deleted"`

	// The type of the item. One of "job", "story", "comment", "poll", or "pollopt".
	Type ItemType `json:"type"`

	// The username of the item's author.
	By *string `json:"by"`

	// Creation date of the item in Unix time.
	Time *int64 `json:"time"`

	// The comment, story, or poll text in HTML.
	Text *string `json:"text"`

	// `true` if the item is dead.
	Dead *bool `json:"dead"`

	// The comment's parent: either another comment or the relevant story.
	Parent *ItemId `json:"parent"`

	// The pollopt's associated poll.
	Poll *ItemId `json:"poll"`

	// The ids of the item's comments, in ranked display order.
	Kids []ItemId `json:"kids"`

	// The url of the story.
	Url *string `json:"url"`

	// The story's score, or the votes for a pollopt.
	Score *int32 `json:"score"`

	// The title of the story, poll, or job in HTML
	Title *string `json:"title"`

	// A list of related pollopts, in display order.
	Parts []ItemId `json:"parts"`

	// In the case of stories or polls, the total comment count.
	Descendants *int32 `json:"descendants"`
}

type FrontPageItemsRanking int

const (
	Top FrontPageItemsRanking = iota
	Best
	New
)

func (r FrontPageItemsRanking) ToPointer() *FrontPageItemsRanking {
	return &r
}

type SearchItemsRanking int

const (
	Date SearchItemsRanking = iota
	Popularity
)

func (r SearchItemsRanking) ToPointer() *SearchItemsRanking {
	return &r
}

type SearchRequest struct {
	Query   string
	Tags    string
	Ranking SearchItemsRanking
	Limit   int
}

type SearchResultJson struct {
	Id string `json:"objectID"`
}

type SearchResponseJson struct {
	Hits []SearchResultJson `json:"hits"`
}
