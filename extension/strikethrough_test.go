package extension

import (
	"testing"

	"github.com/enkogu/goldmark"
	"github.com/enkogu/goldmark/renderer/html"
	"github.com/enkogu/goldmark/testutil"
)

func TestStrikethrough(t *testing.T) {
	markdown := goldmark.New(
		goldmark.WithRendererOptions(
			html.WithUnsafe(),
		),
		goldmark.WithExtensions(
			Strikethrough,
		),
	)
	testutil.DoTestCaseFile(markdown, "_test/strikethrough.txt", t)
}
