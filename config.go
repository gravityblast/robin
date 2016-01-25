package robin

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

type ConfigDecodeError struct {
	err error
}

func (e *ConfigDecodeError) Error() string {
	return fmt.Sprintf("error parsing config file: %s", e.err.Error())
}

func newConfigDecodeError(err error) *ConfigDecodeError {
	return &ConfigDecodeError{
		err: err,
	}
}

type ConfigFieldError struct {
	field   string
	message string
}

func (e *ConfigFieldError) Error() string {
	return fmt.Sprintf("%s %s", e.field, e.message)
}

func newConfigFieldError(f string, m string) *ConfigFieldError {
	return &ConfigFieldError{
		field:   f,
		message: m,
	}
}

type dataExtractorConfig struct {
	Attribute *string `json:"attribute"`
	Selector  string  `json:"selector"`
}

func (d *dataExtractorConfig) dataExtractor() DataExtractor {
	if d.Attribute == nil {
		return NewTextExtractor(d.Selector)
	} else {
		return NewAttributeExtractor(d.Selector, *d.Attribute)
	}
}

func (d *dataExtractorConfig) validate() error {
	if d.Selector == "" {
		return newConfigFieldError("selector", "can't be blank")
	}

	if d.Attribute != nil && *d.Attribute == "" {
		return newConfigFieldError("attribute", "can't be blank")
	}

	return nil
}

type extractorFieldsConfig map[string]*dataExtractorConfig

type extractorConfig struct {
	Selector     string
	FieldsConfig extractorFieldsConfig `json:"fields"`
}

func (e *extractorConfig) extractor(name string) Extractor {
	de := make(DataExtractors)

	for name, dec := range e.FieldsConfig {
		de[name] = dec.dataExtractor()
	}

	ie := NewItemExtractor(name, de)
	ie.Selector = e.Selector

	return ie
}

func (e *extractorConfig) validate() error {
	if len(e.FieldsConfig) == 0 {
		return newConfigFieldError("fields", "can't be empty")
	}

	return nil
}

type followerConfig struct {
	Selector    string                `json:"selector"`
	LinksConfig extractorFieldsConfig `json:"links"`
}

func (e *followerConfig) follower(name string) Follower {
	de := make(DataExtractors)

	for name, dec := range e.LinksConfig {
		de[name] = dec.dataExtractor()
	}

	ie := NewFollower(name, e.Selector, de)
	ie.Selector = e.Selector

	return ie
}

type scraperConfig struct {
	URL              string                      `json:"url"`
	ExtractorsConfig map[string]*extractorConfig `json:"extractors"`
	FollowersConfig  map[string]*followerConfig  `json:"followers"`
}

func (sc *scraperConfig) scraper(name string) (*Scraper, error) {
	fe, err := NewFetcher(sc.URL)
	if err != nil {
		return nil, err
	}

	s := newScraper(name, fe)
	for name, c := range sc.ExtractorsConfig {
		e := c.extractor(name)
		s.Extractors[name] = e
	}

	for name, c := range sc.FollowersConfig {
		fo := c.follower(name)
		s.Followers[name] = fo
	}

	return s, nil
}

func (sc *scraperConfig) validate() error {
	return nil
}

type ScrapersConfig map[string]scraperConfig

func (ssc ScrapersConfig) Scrapers() (map[string]*Scraper, error) {
	scrapers := make(map[string]*Scraper)
	for name, sc := range ssc {
		s, err := sc.scraper(name)
		if err != nil {
			return nil, err
		}

		scrapers[name] = s
	}

	return scrapers, nil
}

func NewConfigFromReader(r io.Reader) (ScrapersConfig, error) {
	var ssc ScrapersConfig
	err := json.NewDecoder(r).Decode(&ssc)
	if err != nil {
		fmt.Printf("--------> %+v\n", err)
		return nil, newConfigDecodeError(err)
	}

	return ssc, nil
}

func NewConfigFromFile(p string) (ScrapersConfig, error) {
	f, err := os.Open(p)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return NewConfigFromReader(f)
}
