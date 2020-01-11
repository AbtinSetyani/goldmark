package extension

import (
	"testing"

	"github.com/anytypeio/goldmark"
	"github.com/anytypeio/goldmark/renderer/html"
	"github.com/anytypeio/goldmark/testutil"
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
