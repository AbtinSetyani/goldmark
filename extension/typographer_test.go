package extension

import (
	"testing"

	"github.com/enkogu/goldmark"
	"github.com/enkogu/goldmark/renderer/html"
	"github.com/enkogu/goldmark/testutil"
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
