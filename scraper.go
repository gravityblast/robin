package robin

import (
	"fmt"
	"net/url"

	"github.com/PuerkitoBio/goquery"
)

type Scraper struct {
	Name       string
	Fetcher    Fetcher
	Exporter   Exporter
	Extractors Extractors
	Followers  Followers
}

func newScraper(name string, fetcher Fetcher) *Scraper {
	return &Scraper{
		Name:       name,
		Fetcher:    fetcher,
		Exporter:   NewStdoutExporter(),
		Extractors: make(Extractors),
		Followers:  make(Followers),
	}
}

func (s *Scraper) Scrape(log AppLogger, exp Exporter, q *exportQueue) error {
	res := s.Fetcher.Fetch(log)
	if err := res.Error(); err != nil {
		log.Error(fmt.Sprintf("%v", err))
		return err
	}
	r := res.ReadCloser()
	defer r.Close()

	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		log.Error(fmt.Sprintf("%v", err))
		return err
	}

	for name, ext := range s.Extractors {
		log.Debug(fmt.Sprintf("extract `%s` from `%s`", name, s.Fetcher.URL()))
		items := ext.Extract(doc.Selection)
		for _, item := range items {
			item["_item"] = name
			j := &exportJob{
				exporter: exp,
				item:     item,
			}
			q.push(j)
		}
	}

	for name, follower := range s.Followers {
		links := follower.Links(doc.Selection)
		log.Debug(fmt.Sprintf("Using follower `%s`", name))
		for _, href := range links {
			url, err := url.Parse(href)
			if err != nil {
				log.Error(fmt.Sprintf("Invalid URL `%s`: %s", href, err))
			}
			log.Debug(fmt.Sprintf("Link extracted: %s", url))
			fetcher := s.Fetcher.New(url)
			log.Debug(fmt.Sprintf("Following URL: %s", fetcher.URL()))
		}
	}

	q.close()

	return nil
}
