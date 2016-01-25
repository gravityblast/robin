package robin

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"sync"
	"time"
)

type fakeClient struct {
	httpClient *httpClient
}

func (fc *fakeClient) Do(req *http.Request) (*http.Response, error) {
	resp := &http.Response{}
	root := "./fixtures/html"
	f, err := os.Open(filepath.Join(root, filepath.FromSlash(path.Clean(req.URL.Path))))
	if err != nil {
		return nil, err
	}

	resp.Body = f

	return resp, nil
}

func httpFetcher(path string) Fetcher {
	url, _ := url.Parse(fmt.Sprintf("http://localhost:0000/%s", path))
	f := NewHTTPFetcher(url)
	f.client = &fakeClient{}

	return f
}

type fakeExporter struct {
	sync.Mutex
	items Items
}

func (fe *fakeExporter) Export(item Item) {
	rand.Seed(int64(time.Now().Nanosecond()))
	time.Sleep(3 * time.Second)
	err := json.NewEncoder(os.Stdout).Encode(item)
	if err != nil {
		log.Fatal(err)
	}

	fe.Lock()
	fe.items = append(fe.items, item)
	fe.Unlock()
}

func newTestItemScraper() *Scraper {
	f := httpFetcher("/item-1.html")
	s := newScraper("", f)

	title := NewTextExtractor("h1")
	desc := NewTextExtractor(".description")
	page := NewItemExtractor("page", DataExtractors{"title": title, "description": desc})

	name := NewTextExtractor(".name")
	surname := NewTextExtractor(".lastname")
	author := NewItemExtractor("author", DataExtractors{"name": name, "surname": surname})
	author.Selector = ".authors .author"

	s.Extractors = Extractors{"author": author, "page": page}

	return s
}

// func TestIntegration(t *testing.T) {
// 	// assert := assert.New(t)
// 	require := require.New(t)
// 	// itemScraper := newTestItemScraper()

// 	// r := NewRunner(map[string]*Scraper{"main": itemScraper})

// 	exp := &fakeExporter{}
// 	// err := r.Run("main", ex, 3)
// 	// if err != nil {
// 	// 	log.Fatal(err)
// 	// }

// 	c, err := NewConfigFromFile("fixtures/config.json")
// 	require.Nil(err)

// 	ss, err := c.Scrapers()
// 	require.Nil(err)

// 	r := NewRunner(ss, DefaultRunnerOptions)
// 	err = r.Run("index", exp)
// 	require.Nil(err)

// 	// expectedItems := Items{
// 	// 	{"_item": "page", "title": "Item 1", "description": "description 1"},
// 	// 	{"_item": "author", "name": "Foo", "surname": "McBar"},
// 	// 	{"_item": "author", "name": "Foo", "surname": "McBar"},
// 	// }

// 	// assert.Len(ex.items, len(expectedItems))
// 	// for _, expectedItem := range expectedItems {
// 	// 	assert.Contains(ex.items, expectedItem)
// 	// }
// }
