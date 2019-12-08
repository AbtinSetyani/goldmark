package goldmark_test

import (
	"testing"

	. "github.com/enkogu/goldmark"
	"github.com/enkogu/goldmark/parser"
	"github.com/enkogu/goldmark/testutil"
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
