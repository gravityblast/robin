package robin

import (
	"bytes"
	"encoding/json"
	"log"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDataExtractorConfig(t *testing.T) {
	tests := []struct {
		jsonText              string
		validationError       error
		expectedDataExtractor DataExtractor
		selector              string
	}{
		{
			`{}`,
			&ConfigFieldError{field: "selector", message: "can't be blank"},
			&textExtractor{},
			".container h1",
		},
		{
			`{"selector": "h1", "attribute":""}`,
			&ConfigFieldError{field: "attribute", message: "can't be blank"},
			&attributeExtractor{Selector: "h1"},
			".container h1",
		},
		{
			`{"selector":".container h1"}`,
			nil,
			&textExtractor{Selector: ".container h1"},
			".container h1",
		},
		{
			`{"selector":".container h1", "attribute":"class"}`,
			nil,
			&attributeExtractor{Selector: ".container h1", AttributeName: "class"},
			".container h1",
		},
	}

	for _, test := range tests {
		var c dataExtractorConfig
		r := bytes.NewBufferString(test.jsonText)
		err := json.NewDecoder(r).Decode(&c)
		if err != nil {
			log.Fatal(err)
		}

		err = c.validate()
		assert.Equal(t, test.validationError, err)

		de := c.dataExtractor()
		assert.Equal(t, test.expectedDataExtractor, de)
	}
}

func TestExtractorConfig(t *testing.T) {
	tests := []struct {
		jsonText              string
		validationError       error
		expectedItemExtractor Extractor
	}{
		{
			`{}`,
			&ConfigFieldError{field: "fields", message: "can't be empty"},
			&ItemExtractor{DataExtractors: DataExtractors{}},
		},
		{
			`{
				"selector": "ul li",
				"fields": {
					"title": {
						"selector": "h2"
					},
					"content": {
						"selector": "p"
					}
				}
			}`,
			nil,
			&ItemExtractor{
				Selector: "ul li",
				DataExtractors: DataExtractors{
					"title": &textExtractor{
						Selector: "h2",
					},
					"content": &textExtractor{
						Selector: "p",
					},
				},
			},
		},
	}

	for _, test := range tests {
		var c extractorConfig
		r := bytes.NewBufferString(test.jsonText)
		err := json.NewDecoder(r).Decode(&c)
		if err != nil {
			log.Fatal(err)
		}

		err = c.validate()
		assert.Equal(t, test.validationError, err)

		ie := c.extractor("")
		assert.Equal(t, test.expectedItemExtractor, ie)
	}
}

func TestScraperConfig(t *testing.T) {
	url, err := url.Parse("http://foo.local/bar")
	if err != nil {
		log.Fatal(err)
	}

	tests := []struct {
		jsonText        string
		validationError error
		expectedScraper *Scraper
	}{
		{
			`{
				"url": "http://foo.local/bar",
				"extractors": {
					"item_1": {
						"selector": ".container",
						"fields": {
							"title": {
								"selector": ".container h1"
							},
							"content": {
								"selector": ".container .content"
							}
						}
					}
				},
				"followers": {
					"pagination": {
						"selector": ".container .pagination",
						"page": {
							"selector": "li a",
							"attribute": "href"
						}
					}
				}
			}`,
			nil,
			&Scraper{
				Fetcher:  NewHTTPFetcher(url),
				Exporter: NewStdoutExporter(),
				Followers: Followers{
					"pagination": &follower{
						name:           "pagination",
						Selector:       ".container .pagination",
						DataExtractors: make(DataExtractors),
					},
				},
				Extractors: map[string]Extractor{
					"item_1": &ItemExtractor{
						name:     "item_1",
						Selector: ".container",
						DataExtractors: DataExtractors{
							"title": &textExtractor{
								Selector: ".container h1",
							},
							"content": &textExtractor{
								Selector: ".container .content",
							},
						},
					},
				},
			},
		},
	}

	for _, test := range tests {
		var c scraperConfig
		r := bytes.NewBufferString(test.jsonText)
		err := json.NewDecoder(r).Decode(&c)
		if err != nil {
			log.Fatal(err)
		}

		err = c.validate()
		assert.Equal(t, test.validationError, err)

		s, err := c.scraper("")
		assert.Nil(t, err)
		assert.Equal(t, test.expectedScraper, s)
	}
}

func TestScrapersConfig(t *testing.T) {
	content := `{
		"foo": {
			"url": "http://foo.local/bar",
			"extractors": {
				"item_1": {
					"selector": ".container",
					"fields": {
						"title": {
							"selector": ".container h1"
						},
						"content": {
							"selector": ".container .content"
						}
					}
				}
			}
		},
		"bar": {
			"url": "http://foo.local/bar",
			"extractors": {
				"item_1": {
					"selector": ".container",
					"fields": {
						"title": {
							"selector": ".container h1"
						},
						"content": {
							"selector": ".container .content"
						}
					}
				}
			}
		}
	}`

	var ssc ScrapersConfig
	err := json.NewDecoder(bytes.NewBufferString(content)).Decode(&ssc)
	if err != nil {
		log.Fatal(err)
	}

	assert.Len(t, ssc, 2)
}
