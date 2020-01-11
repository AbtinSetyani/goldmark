package blocksUtil

import "github.com/anytypeio/go-anytype-library/pb/model"

type BlocksWriter struct {
	IsCurrentBlock bool
	blockBuffer    *model.Block
	TextBuffer     string
	MarksBuffer    []*model.BlockContentTextMark
	BlocksList     []*model.Block
}

func (b *BlocksWriter) Flush() (string, error) {
	return b.TextBuffer, nil
}