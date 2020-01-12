// Package renderer renders the given AST to certain formats.
package blocksUtil

import (
	"bufio"
	"github.com/anytypeio/go-anytype-library/pb/model"
)


type RenderState struct {
	isCurrentBlock bool
	blockBuffer    *model.Block
	textBuffer     string
	marksBuffer    []model.BlockContentTextMark
	blocksList     []model.Block
}

type Writer struct {
	*bufio.Writer
	rs     RenderState
}

func ExtendWriter (writer *bufio.Writer) (Writer) {
	return  Writer{writer, RenderState{}}
}
