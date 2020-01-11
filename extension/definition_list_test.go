package extension

import (
	"testing"

	"github.com/anytypeio/goldmark"
	"github.com/anytypeio/goldmark/renderer/html"
	"github.com/anytypeio/goldmark/testutil"
)

func TestDefinitionList(t *testing.T) {
	markdown := goldmark.New(
		goldmark.WithRendererOptions(
			html.WithUnsafe(),
		),
		goldmark.WithExtensions(
			DefinitionList,
		),
	)
	testutil.DoTestCaseFile(markdown, "_test/definition_list.txt", t)
}
