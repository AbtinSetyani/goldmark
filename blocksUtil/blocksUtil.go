// Package renderer renders the given AST to certain formats.
package blocksUtil

import (
	"bufio"
	"github.com/anytypeio/go-anytype-library/pb/model"
	"github.com/anytypeio/goldmark/util"
)

/*
type renderState struct {
	isCurrentBlock bool
	blockBuffer    *model.Block
	textBuffer     string
	marksBuffer    []model.BlockContentTextMark
	blocksList     []model.Block
}

type RenderState interface {
	SetText (text string)
	GetText () string
}

func  NewRenderState () RenderState {
	rs := &renderState{}
	return rs
}
*/
/*func (r *renderState) GetText () string {
	return r.textBuffer
}

func (r *renderState) SetText (text string) {
	r.textBuffer = text
}
*/


func ExtendWriter (writer *bufio.Writer) RWriter {
	return  &rWriter{Writer: writer}
}




type rWriter struct {
	*bufio.Writer

	isCurrentBlock bool
	blockBuffer    *model.Block
	textBuffer     string
	marksBuffer    []model.BlockContentTextMark
	blocksList     []model.Block
}

type RWriter interface {
	util.BufWriter
	//SetText (text string)
	//GetText () string
}

func NewRWriter (writer *bufio.Writer) RWriter {
	return  &rWriter{Writer: writer}
}


func (rw *rWriter) GetText () string {
	return rw.textBuffer
}

func (rw *rWriter) SetText (text string) {
	rw.textBuffer = text
}
