// Package renderer renders the given AST to certain formats.
package blocksUtil

import (
	"bufio"
	"github.com/anytypeio/go-anytype-library/pb/model"
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
	CloseTextBlock()
	GetBlocks() []*model.Block

	GetMarkStart () int
	SetMarkStart ()

	AddMark (mark model.BlockContentTextMark)

	GetBlockTest() *model.Block
	OpenNewTextBlock (model.BlockContentTextStyle)

	SetIsNumberedList(isNumbered bool)
	GetIsNumberedList() (isNumbered bool)
}

type contentBuff struct {

}

type rWriter struct {
	*bufio.Writer

	isNumberedList    bool

	textBuffer        string
	marksBuffer       []*model.BlockContentTextMark
	marksStartQueue   []int
	textStylesQueue   []model.BlockContentTextStyle
	blocks            []*model.Block
}

func (rw *rWriter) SetMarkStart () {
	rw.marksStartQueue = append(rw.marksStartQueue, len(rw.textBuffer))
}

func (rw *rWriter) GetMarkStart () int {
	return rw.marksStartQueue[len(rw.marksStartQueue) - 1]
}


func (rw *rWriter) AddMark (mark model.BlockContentTextMark) {
	s := rw.marksStartQueue
	rw.marksStartQueue = s[:len(s)-1]
	rw.marksBuffer = append(rw.marksBuffer, &mark)
}

func (rw *rWriter) OpenNewTextBlock (style model.BlockContentTextStyle) {
	rw.textStylesQueue = append(rw.textStylesQueue, style)
}

func (rw *rWriter) GetBlockTest() *model.Block {
	newBlock := model.Block{
		Content: &model.BlockContentOfText{
			Text: &model.BlockContentText{
				//TODO Style: content.style,
				Text: "123123",
			},
		},
	}
	return &newBlock
}

func (rw *rWriter) GetBlocks() []*model.Block {
	return rw.blocks
}

func (rw *rWriter) SetIsNumberedList (isNumbered bool) {
	rw.isNumberedList = isNumbered
}

func (rw *rWriter) GetIsNumberedList() (isNumbered bool) {
	return rw.isNumberedList
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
}

func (rw *rWriter) CloseTextBlock() {
	s := rw.textStylesQueue
	var style model.BlockContentTextStyle;

	if len(rw.textStylesQueue) > 0 {
		style, rw.textStylesQueue = s[len(s)-1], s[:len(s)-1]
	}

	newBlock := model.Block{
		Content: &model.BlockContentOfText{
			Text: &model.BlockContentText{
				Text: rw.textBuffer,
				Style: style,
				Marks: &model.BlockContentTextMarks{
					Marks: rw.marksBuffer,
				},
			},
		},
	}
	rw.blocks = append(rw.blocks, &newBlock)
	rw.marksStartQueue = []int{}
	rw.marksBuffer = []*model.BlockContentTextMark{}
	rw.textBuffer = ""
}

func (rw *rWriter) OpenNewBlock(content model.IsBlockContent) {
}