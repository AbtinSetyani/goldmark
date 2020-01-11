package goldmark_test

import (
	"testing"

	. "github.com/anytypeio/goldmark"
	"github.com/anytypeio/goldmark/parser"
	"github.com/anytypeio/goldmark/testutil"
)

func TestAttributeAndAutoHeadingID(t *testing.T) {
	markdown := New(
		WithParserOptions(
			parser.WithAttribute(),
			parser.WithAutoHeadingID(),
		),
	)
	testutil.DoTestCaseFile(markdown, "_test/options.txt", t)
}
