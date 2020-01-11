package extension

import (
	"testing"

	"github.com/anytypeio/goldmark"
	"github.com/anytypeio/goldmark/renderer/html"
	"github.com/anytypeio/goldmark/testutil"
)

func TestFootnote(t *testing.T) {
	markdown := goldmark.New(
		goldmark.WithRendererOptions(
			html.WithUnsafe(),
		),
		goldmark.WithExtensions(
			Footnote,
		),
	)
	testutil.DoTestCaseFile(markdown, "_test/footnote.txt", t)
}
