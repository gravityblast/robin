package robin

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExtractor(t *testing.T) {
	sel := newTestSelection(t, `
	<body>
		<ul>
			<li>
				<h2>Item 1</h2>
				<img class="logo" src="http://localhost/item-1.jpg" />
				<span class="name">Foo 1</span>
				<span class="description">Bar 1</span>
			</li>
			<li>
				<h2>Item 2</h2>
				<img class="logo" src="http://localhost/item-2.jpg" />
				<span class="name">Foo 2</span>
				<span class="description">Bar 2</span>
			</li>
		</ul>
	</body>`)

	de := DataExtractors{
		"title":       NewTextExtractor("h2"),
		"name":        NewTextExtractor("span.name"),
		"description": NewTextExtractor("span.description"),
		"image":       NewAttributeExtractor("img.logo", "src"),
	}

	ie := NewItemExtractor("", de)
	ie.Selector = "ul li"

	expected := Items{
		{
			"title":       "Item 1",
			"image":       "http://localhost/item-1.jpg",
			"name":        "Foo 1",
			"description": "Bar 1",
		},
		{
			"title":       "Item 2",
			"image":       "http://localhost/item-2.jpg",
			"name":        "Foo 2",
			"description": "Bar 2",
		},
	}

	assert.Equal(t, expected, ie.Extract(sel))
}
