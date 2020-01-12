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
	source := []byte("## Hello world!")
	var b bytes.Buffer

/*	writer := bufio.NewWriter(&b)
	rState := blocksUtil.NewRenderState()

	BR := blocksUtil.ExtendWriter(writer, rState)
*/

	writer := bufio.NewWriter(&b)
	BR := blocksUtil.NewRWriter(writer)


	err := markdown.ConvertBlocks(source, BR)
	if err != nil {
		t.Error(err.Error())
	}

	fmt.Println("rState:", BR)
	fmt.Println("b:", b.String())
/*	if b.String() != "<pre><code>a\n</code></pre>\n" {
		t.Errorf("%s \n---------\n %s", source, b.String())
	}*/
}
