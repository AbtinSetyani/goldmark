// Package renderer renders the given AST to certain formats.
package blocksUtil

import (
	"bufio"
	"github.com/anytypeio/go-anytype-library/pb/model"
	"github.com/anytypeio/go-anytype-middleware/core/block/simple/text"
	"io"
)

// A RWriter is a subset of the bufio.Writer .
type RWriter interface {
	// TODO: LEGACY, remove it later
	io.Writer
	Available() int
	Buffered() int
	Flush() error
	WriteByte(c byte) error
	WriteRune(r rune) (size int, err error)
	WriteString(s string) (int, error)

	// Main part
	GetText () string
	AddTextToBuffer(s string)
	OpenNewBlock(content model.IsBlockContent)
	CloseCurrentBlock()
}

type rWriter struct {
	*bufio.Writer

	isCurrentBlock bool
	blockBuffer        *model.Block
	textBuffer         string
	marksBuffer        []model.BlockContentTextMark
	blocksList         []model.Block
	blockContentText   model.BlockContentText
}

func NewRWriter (writer *bufio.Writer) RWriter {
	return  &rWriter{Writer: writer}
}

func (rw *rWriter) GetText () string {
	return rw.textBuffer
}

func (rw *rWriter) AddTextToBuffer (text string) {
	rw.textBuffer += text
}

func (rw *rWriter) CloseCurrentBlock() {
	rw.isCurrentBlock = false

	newBlock := model.Block{
		Content: &model.BlockContentOfText{
			Text: &rw.blockContentText,
		},
	}

	rw.blocksList = append(rw.blocksList, newBlock)
	rw.blockBuffer = &model.Block{}
	rw.textBuffer = ""
}

func (rw *rWriter) OpenNewBlock(content model.IsBlockContent) {
	if rw.isCurrentBlock {
		rw.CloseCurrentBlock()
	}
	rw.isCurrentBlock = true
	rw.blockBuffer = &model.Block{
		//Id: "3",
		Content: content,
	}
}


func (rw *rWriter) OpenNewTextBlock(content model.IsBlockContent) {
	if rw.isCurrentBlock {
		rw.CloseCurrentBlock()
	}
	rw.isCurrentBlock = true
	rw.blockBuffer = &model.Block{
		//Id: "3",
		Content: content,
	}
}