// Package renderer renders the given AST to certain formats.
package renderer

import (
	"github.com/anytypeio/goldmark"
	"sync"

	"github.com/anytypeio/goldmark/ast"
	"github.com/anytypeio/goldmark/util"
)

// A BlocksConfig struct is a data structure that holds configuration of the BlocksRenderer.
type BlocksConfig struct {
	BlocksOptions       map[BlocksOptionName]interface{}
	NodeBlocksRenderers util.PrioritizedSlice
}

// NewBlocksConfig returns a new BlocksConfig
func NewBlocksConfig() *BlocksConfig {
	return &BlocksConfig{
		BlocksOptions:       map[BlocksOptionName]interface{}{},
		NodeBlocksRenderers: util.PrioritizedSlice{},
	}
}

type BlocksOption interface {
	SetBlocksConfig(*BlocksConfig)
}

// An BlocksOptionName is a name of the option.
type BlocksOptionName string

// An BlocksOption interface is a functional option type for the BlocksRenderer.
type BlocksBlocksOption interface {
	SetBlocksConfig(*BlocksConfig)
}

type withNodeBlocksRenderers struct {
	value []util.PrioritizedValue
}

func (o *withNodeBlocksRenderers) SetBlocksConfig(c *BlocksConfig) {
	c.NodeBlocksRenderers = append(c.NodeBlocksRenderers, o.value...)
}

// WithNodeBlocksRenderers is a functional option that allow you to add
// NodeBlocksRenderers to the blocksRenderer.
func WithNodeBlocksRenderers(ps ...util.PrioritizedValue) BlocksOption {
	return &withNodeBlocksRenderers{ps}
}

type withBlocksOption struct {
	name  BlocksOptionName
	value interface{}
}

func (o *withBlocksOption) SetBlocksConfig(c *BlocksConfig) {
	c.BlocksOptions[o.name] = o.value
}

// WithBlocksOption is a functional option that allow you to set
// an arbitrary option to the parser.
func WithBlocksOption(name BlocksOptionName, value interface{}) BlocksOption {
	return &withBlocksOption{name, value}
}

// A SetBlocksOptioner interface sets given option to the object.
type SetBlocksOptioner interface {
	// SetBlocksOption sets given option to the object.
	// Unacceptable options may be passed.
	// Thus implementations must ignore unacceptable options.
	SetBlocksOption(name BlocksOptionName, value interface{})
}

// NodeBlocksRendererFunc is a function that renders a given node.
type NodeBlocksRendererFunc func(bwr goldmark.BlocksWriter, source []byte, n ast.Node, entering bool) (ast.WalkStatus, error)

// A NodeBlocksRenderer interface offers NodeBlocksRendererFuncs.
type NodeBlocksRenderer interface {
	// BlocksRendererFuncs registers NodeBlocksRendererFuncs to given NodeBlocksRendererFuncRegisterer.
	RegisterFuncs(NodeBlocksRendererFuncRegisterer)
}

// A NodeBlocksRendererFuncRegisterer registers
type NodeBlocksRendererFuncRegisterer interface {
	// Register registers given NodeBlocksRendererFunc to this object.
	Register(ast.NodeKind, NodeBlocksRendererFunc)
}

// A BlocksRenderer interface renders given AST node to given
// writer with given BlocksRenderer.
type BlocksRenderer interface {
	BlocksRender(w goldmark.BlocksWriter, source []byte, n ast.Node) error

	// AddBlocksOptions adds given option to this blocksRenderer.
	AddBlocksOptions(...BlocksOption)
}

type blocksRenderer struct {
	config               *BlocksConfig
	options              map[BlocksOptionName]interface{}
	nodeBlocksRendererFuncsTmp map[ast.NodeKind]NodeBlocksRendererFunc
	maxKind              int
	nodeBlocksRendererFuncs    []NodeBlocksRendererFunc
	initSync             sync.Once
}

// NewBlocksRenderer returns a new BlocksRenderer with given options.
func NewBlocksRenderer(options ...BlocksOption) BlocksRenderer {
	config := NewBlocksConfig()
	for _, opt := range options {
		opt.SetBlocksConfig(config)
	}

	r := &blocksRenderer{
		options:              map[BlocksOptionName]interface{}{},
		config:               config,
		nodeBlocksRendererFuncsTmp: map[ast.NodeKind]NodeBlocksRendererFunc{},
	}

	return r
}

func (r *blocksRenderer) AddBlocksOptions(opts ...BlocksOption) {
	for _, opt := range opts {
		opt.SetBlocksConfig(r.config)
	}
}

func (r *blocksRenderer) Register(kind ast.NodeKind, v NodeBlocksRendererFunc) {
	r.nodeBlocksRendererFuncsTmp[kind] = v
	if int(kind) > r.maxKind {
		r.maxKind = int(kind)
	}
}

// Render renders the given AST node to the given writer with the given BlocksRenderer.
func (r *blocksRenderer) Render(w goldmark.BlocksWriter, source []byte, n ast.Node) error {
	r.initSync.Do(func() {
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
/*	if !ok {
		writer = bufio.NewWriter(w)
	}*/
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
		return err
	}
	return w.Flush()
}
