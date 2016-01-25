package robin

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHTTPFetcher_New(t *testing.T) {
	u1, err := url.Parse("http://test.local/foo/bar/index.html")
	require.Nil(t, err)
	f1 := NewHTTPFetcher(u1)

	u2, err := url.Parse("baz.html")
	require.Nil(t, err)
	f2 := f1.New(u2)
	assert.Equal(t, "http://test.local/foo/bar/baz.html", f2.URL().String())

	u3, err := url.Parse("/index.html")
	require.Nil(t, err)
	f3 := f1.New(u3)
	assert.Equal(t, "http://test.local/index.html", f3.URL().String())
}
