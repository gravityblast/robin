package robin

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type Response interface {
	ReadCloser() io.ReadCloser
	Error() error
}

type response struct {
	r   io.ReadCloser
	err error
}

func (r *response) ReadCloser() io.ReadCloser {
	return r.r
}

func (r *response) Error() error {
	return r.err
}

type Fetcher interface {
	Fetch(AppLogger) Response
	New(*url.URL) Fetcher
	URL() *url.URL
}

func NewFetcher(u string) (Fetcher, error) {
	url, err := url.Parse(u)
	if err != nil {
		return nil, err
	}

	return NewHTTPFetcher(url), nil
}

type httpClient interface {
	Do(*http.Request) (*http.Response, error)
}

type HTTPFetcher struct {
	client httpClient
	u      *url.URL
}

func NewHTTPFetcher(u *url.URL) *HTTPFetcher {
	client := &http.Client{}
	return &HTTPFetcher{
		client: client,
		u:      u,
	}
}

func (f *HTTPFetcher) Fetch(log AppLogger) Response {
	r := &response{}
	log.Info(fmt.Sprintf("fetching `%s`", f.u.String()))
	req, err := http.NewRequest("GET", f.u.String(), nil)
	if err != nil {
		r.err = err
		log.Error(fmt.Sprintf("%v", err))
		return r
	}

	resp, err := f.client.Do(req)
	if err != nil {
		r.err = err
		log.Error(fmt.Sprintf("%v", err))
		return r
	}

	r.r = resp.Body

	return r
}

func (f *HTTPFetcher) URL() *url.URL {
	return f.u
}

func (f *HTTPFetcher) New(u *url.URL) Fetcher {
	return NewHTTPFetcher(f.u.ResolveReference(u))
}
