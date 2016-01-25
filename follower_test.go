package robin

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type fakeQueue struct {
}

func TestFollower(t *testing.T) {
	assert := assert.New(t)

	sel := newTestSelection(t, `
	<body>
		<ul>
			<li>
				<h2>Item 1</h2>
				<a href="/item-1.html">Details</a>
			</li>
			<li>
				<h2>Item 2</h2>
				<a href="/item-2.html">Details</a>
			</li>
		</ul>
	</body>`)

	de := DataExtractors{
		"link": NewAttributeExtractor("a", "href"),
	}

	f := NewFollower("", "ul li", de)
	links := f.Links(sel)
	assert.Equal([]string{"/item-1.html", "/item-2.html"}, links)
}
