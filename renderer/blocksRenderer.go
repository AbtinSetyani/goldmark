// Package renderer renders the given AST to certain formats.
package renderer

import (
	"github.com/anytypeio/goldmark/ast"
	"github.com/anytypeio/goldmark/blocksUtil"
)


// A BlocksRenderer interface renders given AST node to given  writer with given BlocksRenderer.
type BlocksRenderer interface {
	BlocksRender(w blocksUtil.BlocksWriter, source []byte, n ast.Node) error
	// AddBlocksOptions adds given option to this blocksRenderer.
	//AddBlocksOptions(...BlocksOption)
}

type blocksRenderer struct {
	//config               *BlocksConfig
	//options              map[BlocksOptionName]interface{}
	//nodeBlocksRendererFuncsTmp map[ast.NodeKind]NodeBlocksRendererFunc
	maxKind              int
	//nodeBlocksRendererFuncs    []NodeBlocksRendererFunc
	//initSync             sync.Once
}


// Render renders the given AST node to the given writer with the given BlocksRenderer.
func (r *blocksRenderer) BlocksRender(w blocksUtil.BlocksWriter, source []byte, n ast.Node) error {
/*	r.initSync.Do(func() {
		r.options = r.config.BlocksOptions
		r.config.NodeBlocksRenderers.Sort()
		l := len(r.config.NodeBlocksRenderers)
		for i := l - 1; i >= 0; i-- {
			v := r.config.NodeBlocksRenderers[i]
			nr, _ := v.Value.(NodeBlocksRenderer)
			if se, ok := v.Value.(SetBlocksOptioner); ok {
				for oname, ovalue := range r.options {
					se.SetBlocksOption(oname, ovalue)
				}
			}
			nr.RegisterFuncs(r)
		}
		r.nodeBlocksRendererFuncs = make([]NodeBlocksRendererFunc, r.maxKind+1)
		for kind, nr := range r.nodeBlocksRendererFuncsTmp {
			r.nodeBlocksRendererFuncs[kind] = nr
		}
		r.config = nil
		r.nodeBlocksRendererFuncsTmp = nil
	})
	//writer, _ := w.NewBlocksWriter()

	err := ast.Walk(n, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		s := ast.WalkStatus(ast.WalkContinue)
		var err error
		f := r.nodeBlocksRendererFuncs[n.Kind()]
		if f != nil {
			s, err = f(w, source, n, entering)
		}
		return s, err
	})
	if err != nil {
		return nil, err
	}
*/
	//blocks := []*model.Block{}
	//block := model.Block{}
	//blocks = append(blocks, &block)

	return nil
}
