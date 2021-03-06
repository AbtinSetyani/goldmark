package extension

import (
	"testing"

	"github.com/enkogu/goldmark"
	"github.com/enkogu/goldmark/renderer/html"
	"github.com/enkogu/goldmark/testutil"
)

func TestTable(t *testing.T) {
	markdown := goldmark.New(
		goldmark.WithRendererOptions(
			html.WithUnsafe(),
		),
		goldmark.WithExtensions(
			Table,
		),
	)
	testutil.DoTestCaseFile(markdown, "_test/table.txt", t)
}
