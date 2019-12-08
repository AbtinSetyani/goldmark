package extension

import (
	"testing"

	"github.com/enkogu/goldmark"
	"github.com/enkogu/goldmark/renderer/html"
	"github.com/enkogu/goldmark/testutil"
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
