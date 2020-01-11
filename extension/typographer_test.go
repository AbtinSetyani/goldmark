package extension

import (
	"testing"

	"github.com/anytypeio/goldmark"
	"github.com/anytypeio/goldmark/renderer/html"
	"github.com/anytypeio/goldmark/testutil"
)

func TestTypographer(t *testing.T) {
	markdown := goldmark.New(
		goldmark.WithRendererOptions(
			html.WithUnsafe(),
		),
		goldmark.WithExtensions(
			Typographer,
		),
	)
	testutil.DoTestCaseFile(markdown, "_test/typographer.txt", t)
}
