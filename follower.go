package robin

import "github.com/PuerkitoBio/goquery"

type Follower interface {
	Name() string
	Links(*goquery.Selection) []string
}

type Followers map[string]Follower

type follower struct {
	name           string
	Fetcher        Fetcher
	Selector       string
	DataExtractors DataExtractors
}

func NewFollower(name string, selector string, de DataExtractors) *follower {
	return &follower{
		name:           name,
		Selector:       selector,
		DataExtractors: de,
	}
}

func (f *follower) Name() string {
	return f.name
}

func (f *follower) Links(sel *goquery.Selection) []string {
	if f.Selector != "" {
		sel = sel.Find(f.Selector)
	}

	var links []string

	sel.Each(func(i int, s *goquery.Selection) {
		for _, ext := range f.DataExtractors {
			links = append(links, ext.Extract(s))
		}
	})

	return links
}
