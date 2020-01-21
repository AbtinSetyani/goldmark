package goldmark_test

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/anytypeio/goldmark/blocksUtil"
	"testing"

	. "github.com/anytypeio/goldmark"
	"github.com/anytypeio/goldmark/renderer/html"
)


func TestConvertBlocks(t *testing.T) {
	markdown := New(WithRendererOptions(
		html.WithXHTML(),
		html.WithUnsafe(),
	))
	source := []byte("## Hello world!\n Olol*ol*olo \n\n 123123")
	var b bytes.Buffer

	writer := bufio.NewWriter(&b)
	BR := blocksUtil.NewRWriter(writer)

	err := markdown.ConvertBlocks(source, BR)
	if err != nil {
		t.Error(err.Error())
	}

	fmt.Println("rState:", BR.GetBlocks())
	fmt.Println("b:", b.String())
}
