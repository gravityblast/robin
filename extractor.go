package robin

import "github.com/PuerkitoBio/goquery"

type Extractor interface {
	Name() string
	Extract(sel *goquery.Selection) Items
}

type Extractors map[string]Extractor

type Item map[string]string

type Items []Item

type ItemExtractor struct {
	name           string
	Selector       string
	DataExtractors DataExtractors
}

func NewItemExtractor(name string, de DataExtractors) *ItemExtractor {
	return &ItemExtractor{
		name:           name,
		DataExtractors: de,
	}
}

func (e *ItemExtractor) Extract(sel *goquery.Selection) Items {
	var items Items

	if e.Selector != "" {
		sel = sel.Find(e.Selector)
	}

	sel.Each(func(i int, s *goquery.Selection) {
		item := make(Item)
		for name, ext := range e.DataExtractors {
			item[name] = ext.Extract(s)
		}
		items = append(items, item)
	})

	return items

}

func (e *ItemExtractor) Name() string {
	return e.name
}
