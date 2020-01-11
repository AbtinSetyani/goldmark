package blocksUtil

import "github.com/anytypeio/go-anytype-library/pb/model"

// BlocksWriter represents
type BlocksWriter struct {
	isCurrentBlock bool
	blockBuffer    *model.Block
	textBuffer     string
	marksBuffer    []*model.BlockContentTextMark
	blocksList     []*model.Block
}

/*type BlocksWriter interface {
	NewBlocksWriter() *BlocksWriter
}
*/
func  NewBlocksWriter() BlocksWriter {
	return BlocksWriter{}
}

func (b *BlocksWriter) Flush() (string, error) {
	return b.textBuffer, nil
}