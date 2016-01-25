package robin

import "github.com/PuerkitoBio/goquery"

type DataExtractor interface {
	Extract(*goquery.Selection) string
}

type DataExtractors map[string]DataExtractor

type textExtractor struct {
	Selector string
}

func NewTextExtractor(selector string) *textExtractor {
	return &textExtractor{
		Selector: selector,
	}
}

func (e *textExtractor) Extract(s *goquery.Selection) string {
	return s.Find(e.Selector).Text()
}

type attributeExtractor struct {
	Selector      string
	AttributeName string
}

func NewAttributeExtractor(selector string, attributeName string) *attributeExtractor {
	return &attributeExtractor{
		Selector:      selector,
		AttributeName: attributeName,
	}
}

func (e *attributeExtractor) Extract(s *goquery.Selection) string {
	v, _ := s.Find(e.Selector).Attr(e.AttributeName)
	return v
}
