package goldmark_test

import (
	"bytes"
	"testing"

	goldmark "github.com/anytypeio/goldmark"
)


func TestEndsWithNonSpaceCharacter(t *testing.T) {
	markdown := goldmark.New()
	source := []byte("```\na\n```")
	var b bytes.Buffer
	err := markdown.Convert(source, &b)
	if err != nil {
		t.Error(err.Error())
	}
	if b.String() != "<pre><code>a\n</code></pre>\n" {
		t.Errorf("%s \n---------\n %s", source, b.String())
	}
}
