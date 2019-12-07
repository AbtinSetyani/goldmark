package extension

import (
	"bytes"
	"regexp"

	"github.com/enkogu/goldmark"
	"github.com/enkogu/goldmark/ast"
	"github.com/enkogu/goldmark/parser"
	"github.com/enkogu/goldmark/text"
	"github.com/enkogu/goldmark/util"
)

var wwwURLRegxp = regexp.MustCompile(`^www\.[-a-zA-Z0-9@:%._\+~#=]{2,256}\.[a-z]{2,6}\b(?:[-a-zA-Z0-9@:%_\+.~#?&//=\(\);,'"]*)`)

var urlRegexp = regexp.MustCompile(`^(?:http|https|ftp):\/\/(?:www\.)?[-a-zA-Z0-9@:%._\+~#=]{2,256}\.[a-z]{2,6}\b([-a-zA-Z0-9@:%_\+.~#?&//=\(\);,'"]*)`)

type linkifyParser struct {
}

var defaultLinkifyParser = &linkifyParser{}

// NewLinkifyParser return a new InlineParser can parse
// text that seems like a URL.
func NewLinkifyParser() parser.InlineParser {
	return defaultLinkifyParser
}

func (s *linkifyParser) Trigger() []byte {
	// ' ' indicates any white spaces and a line head
	return []byte{' ', '*', '_', '~', '('}
}

var protoHTTP = []byte("http:")
var protoHTTPS = []byte("https:")
var protoFTP = []byte("ftp:")
var domainWWW = []byte("www.")

func (s *linkifyParser) Parse(parent ast.Node, block text.Reader, pc parser.Context) ast.Node {
	line, segment := block.PeekLine()
	consumes := 0
	start := segment.Start
	c := line[0]
	// advance if current position is not a line head.
	if c == ' ' || c == '*' || c == '_' || c == '~' || c == '(' {
		consumes++
		start++
		line = line[1:]
	}

	var m []int
	var protocol []byte
	var typ ast.AutoLinkType = ast.AutoLinkURL
	if bytes.HasPrefix(line, protoHTTP) || bytes.HasPrefix(line, protoHTTPS) || bytes.HasPrefix(line, protoFTP) {
		m = urlRegexp.FindSubmatchIndex(line)
	}
	if m == nil && bytes.HasPrefix(line, domainWWW) {
		m = wwwURLRegxp.FindSubmatchIndex(line)
		protocol = []byte("http")
	}
	if m != nil {
		lastChar := line[m[1]-1]
		if lastChar == '.' {
			m[1]--
		} else if lastChar == ')' {
			closing := 0
			for i := m[1] - 1; i >= m[0]; i-- {
				if line[i] == ')' {
					closing++
				} else if line[i] == '(' {
					closing--
				}
			}
			if closing > 0 {
				m[1] -= closing
			}
		} else if lastChar == ';' {
			i := m[1] - 2
			for ; i >= m[0]; i-- {
				if util.IsAlphaNumeric(line[i]) {
					continue
				}
				break
			}
			if i != m[1]-2 {
				if line[i] == '&' {
					m[1] -= m[1] - i
				}
			}
		}
	}
	if m == nil {
		typ = ast.AutoLinkEmail
		stop := util.FindEmailIndex(line)
		if stop < 0 {
			return nil
		}
		at := bytes.IndexByte(line, '@')
		m = []int{0, stop, at, stop - 1}
		if m == nil || bytes.IndexByte(line[m[2]:m[3]], '.') < 0 {
			return nil
		}
		lastChar := line[m[1]-1]
		if lastChar == '.' {
			m[1]--
		}
		if m[1] < len(line) {
			nextChar := line[m[1]]
			if nextChar == '-' || nextChar == '_' {
				return nil
			}
		}
	}
	if m == nil {
		return nil
	}
	if consumes != 0 {
		s := segment.WithStop(segment.Start + 1)
		ast.MergeOrAppendTextSegment(parent, s)
	}
	consumes += m[1]
	block.Advance(consumes)
	n := ast.NewTextSegment(text.NewSegment(start, start+m[1]))
	link := ast.NewAutoLink(typ, n)
	link.Protocol = protocol
	return link
}

func (s *linkifyParser) CloseBlock(parent ast.Node, pc parser.Context) {
	// nothing to do
}

type linkify struct {
}

type linkifyASTTransformer struct {
}

var defaultLinkifyASTTransformer = &linkifyASTTransformer{}

// NewLinkifyASTTransformer returns a new parser.ASTTransformer that
// is related to AutoLink.
func NewLinkifyASTTransformer() parser.ASTTransformer {
	return defaultLinkifyASTTransformer
}

func (a *linkifyASTTransformer) Transform(node *ast.Document, reader text.Reader, pc parser.Context) {
	var autoLinks []*ast.AutoLink
	_ = ast.Walk(node, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if entering {
			if autoLink, ok := n.(*ast.AutoLink); ok {
				autoLinks = append(autoLinks, autoLink)
			}
		}
		return ast.WalkContinue, nil

	})
	for _, autoLink := range autoLinks {
		nested := false
		for p := autoLink.Parent(); p != nil; p = p.Parent() {
			if _, ok := p.(*ast.Link); ok {
				nested = true
				break
			}
		}
		if nested {
			parent := autoLink.Parent()
			parent.ReplaceChild(parent, autoLink, ast.NewString(autoLink.Label(reader.Source())))
		}
	}
}

// Linkify is an extension that allow you to parse text that seems like a URL.
var Linkify = &linkify{}

func (e *linkify) Extend(m goldmark.Markdown) {
	m.Parser().AddOptions(
		parser.WithInlineParsers(
			util.Prioritized(NewLinkifyParser(), 999),
		),
		parser.WithASTTransformers(
			util.Prioritized(NewLinkifyASTTransformer(), 999),
		),
	)
}
