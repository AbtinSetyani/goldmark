package parser

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/enkogu/goldmark/ast"
	"github.com/enkogu/goldmark/text"
	"github.com/enkogu/goldmark/util"
)

var linkLabelStateKey = NewContextKey()

type linkLabelState struct {
	ast.BaseInline

	Segment text.Segment

	IsImage bool

	Prev *linkLabelState

	Next *linkLabelState

	First *linkLabelState

	Last *linkLabelState
}

func newLinkLabelState(segment text.Segment, isImage bool) *linkLabelState {
	return &linkLabelState{
		Segment: segment,
		IsImage: isImage,
	}
}

func (s *linkLabelState) Text(source []byte) []byte {
	return s.Segment.Value(source)
}

func (s *linkLabelState) Dump(source []byte, level int) {
	fmt.Printf("%slinkLabelState: \"%s\"\n", strings.Repeat("    ", level), s.Text(source))
}

var kindLinkLabelState = ast.NewNodeKind("LinkLabelState")

func (s *linkLabelState) Kind() ast.NodeKind {
	return kindLinkLabelState
}

func pushLinkLabelState(pc Context, v *linkLabelState) {
	tlist := pc.Get(linkLabelStateKey)
	var list *linkLabelState
	if tlist == nil {
		list = v
		v.First = v
		v.Last = v
		pc.Set(linkLabelStateKey, list)
	} else {
		list = tlist.(*linkLabelState)
		l := list.Last
		list.Last = v
		l.Next = v
		v.Prev = l
	}
}

func removeLinkLabelState(pc Context, d *linkLabelState) {
	tlist := pc.Get(linkLabelStateKey)
	var list *linkLabelState
	if tlist == nil {
		return
	}
	list = tlist.(*linkLabelState)

	if d.Prev == nil {
		list = d.Next
		if list != nil {
			list.First = d
			list.Last = d.Last
			list.Prev = nil
			pc.Set(linkLabelStateKey, list)
		} else {
			pc.Set(linkLabelStateKey, nil)
		}
	} else {
		d.Prev.Next = d.Next
		if d.Next != nil {
			d.Next.Prev = d.Prev
		}
	}
	if list != nil && d.Next == nil {
		list.Last = d.Prev
	}
	d.Next = nil
	d.Prev = nil
	d.First = nil
	d.Last = nil
}

type linkParser struct {
}

var defaultLinkParser = &linkParser{}

// NewLinkParser return a new InlineParser that parses links.
func NewLinkParser() InlineParser {
	return defaultLinkParser
}

func (s *linkParser) Trigger() []byte {
	return []byte{'!', '[', ']'}
}

var linkDestinationRegexp = regexp.MustCompile(`\s*([^\s].+)`)
var linkTitleRegexp = regexp.MustCompile(`\s+(\)|["'\(].+)`)
var linkBottom = NewContextKey()

func (s *linkParser) Parse(parent ast.Node, block text.Reader, pc Context) ast.Node {
	line, segment := block.PeekLine()
	if line[0] == '!' {
		if len(line) > 1 && line[1] == '[' {
			block.Advance(1)
			pc.Set(linkBottom, pc.LastDelimiter())
			return processLinkLabelOpen(block, segment.Start+1, true, pc)
		}
		return nil
	}
	if line[0] == '[' {
		pc.Set(linkBottom, pc.LastDelimiter())
		return processLinkLabelOpen(block, segment.Start, false, pc)
	}

	// line[0] == ']'
	tlist := pc.Get(linkLabelStateKey)
	if tlist == nil {
		return nil
	}
	last := tlist.(*linkLabelState).Last
	if last == nil {
		return nil
	}
	block.Advance(1)
	removeLinkLabelState(pc, last)
	if s.containsLink(last) { // a link in a link text is not allowed
		ast.MergeOrReplaceTextSegment(last.Parent(), last, last.Segment)
		return nil
	}
	labelValue := block.Value(text.NewSegment(last.Segment.Start+1, segment.Start))
	if util.IsBlank(labelValue) && !last.IsImage {
		ast.MergeOrReplaceTextSegment(last.Parent(), last, last.Segment)
		return nil
	}

	c := block.Peek()
	l, pos := block.Position()
	var link *ast.Link
	var hasValue bool
	if c == '(' { // normal link
		link = s.parseLink(parent, last, block, pc)
	} else if c == '[' { // reference link
		link, hasValue = s.parseReferenceLink(parent, last, block, pc)
		if link == nil && hasValue {
			ast.MergeOrReplaceTextSegment(last.Parent(), last, last.Segment)
			return nil
		}
	}

	if link == nil {
		// maybe shortcut reference link
		block.SetPosition(l, pos)
		ssegment := text.NewSegment(last.Segment.Stop, segment.Start)
		maybeReference := block.Value(ssegment)
		ref, ok := pc.Reference(util.ToLinkReference(maybeReference))
		if !ok {
			ast.MergeOrReplaceTextSegment(last.Parent(), last, last.Segment)
			return nil
		}
		link = ast.NewLink()
		s.processLinkLabel(parent, link, last, pc)
		link.Title = ref.Title()
		link.Destination = ref.Destination()
	}
	if last.IsImage {
		last.Parent().RemoveChild(last.Parent(), last)
		return ast.NewImage(link)
	}
	last.Parent().RemoveChild(last.Parent(), last)
	return link
}

func (s *linkParser) containsLink(last *linkLabelState) bool {
	if last.IsImage {
		return false
	}
	var c ast.Node
	for c = last; c != nil; c = c.NextSibling() {
		if _, ok := c.(*ast.Link); ok {
			return true
		}
	}
	return false
}

func processLinkLabelOpen(block text.Reader, pos int, isImage bool, pc Context) *linkLabelState {
	start := pos
	if isImage {
		start--
	}
	state := newLinkLabelState(text.NewSegment(start, pos+1), isImage)
	pushLinkLabelState(pc, state)
	block.Advance(1)
	return state
}

func (s *linkParser) processLinkLabel(parent ast.Node, link *ast.Link, last *linkLabelState, pc Context) {
	var bottom ast.Node
	if v := pc.Get(linkBottom); v != nil {
		bottom = v.(ast.Node)
	}
	pc.Set(linkBottom, nil)
	ProcessDelimiters(bottom, pc)
	for c := last.NextSibling(); c != nil; {
		next := c.NextSibling()
		parent.RemoveChild(parent, c)
		link.AppendChild(link, c)
		c = next
	}
}

func (s *linkParser) parseReferenceLink(parent ast.Node, last *linkLabelState, block text.Reader, pc Context) (*ast.Link, bool) {
	_, orgpos := block.Position()
	block.Advance(1) // skip '['
	line, segment := block.PeekLine()
	endIndex := util.FindClosure(line, '[', ']', false, true)
	if endIndex < 0 {
		return nil, false
	}

	block.Advance(endIndex + 1)
	ssegment := segment.WithStop(segment.Start + endIndex)
	maybeReference := block.Value(ssegment)
	if util.IsBlank(maybeReference) { // collapsed reference link
		ssegment = text.NewSegment(last.Segment.Stop, orgpos.Start-1)
		maybeReference = block.Value(ssegment)
	}

	ref, ok := pc.Reference(util.ToLinkReference(maybeReference))
	if !ok {
		return nil, true
	}

	link := ast.NewLink()
	s.processLinkLabel(parent, link, last, pc)
	link.Title = ref.Title()
	link.Destination = ref.Destination()
	return link, true
}

func (s *linkParser) parseLink(parent ast.Node, last *linkLabelState, block text.Reader, pc Context) *ast.Link {
	block.Advance(1) // skip '('
	block.SkipSpaces()
	var title []byte
	var destination []byte
	var ok bool
	if block.Peek() == ')' { // empty link like '[link]()'
		block.Advance(1)
	} else {
		destination, ok = parseLinkDestination(block)
		if !ok {
			return nil
		}
		block.SkipSpaces()
		if block.Peek() == ')' {
			block.Advance(1)
		} else {
			title, ok = parseLinkTitle(block)
			if !ok {
				return nil
			}
			block.SkipSpaces()
			if block.Peek() == ')' {
				block.Advance(1)
			} else {
				return nil
			}
		}
	}

	link := ast.NewLink()
	s.processLinkLabel(parent, link, last, pc)
	link.Destination = destination
	link.Title = title
	return link
}

func parseLinkDestination(block text.Reader) ([]byte, bool) {
	block.SkipSpaces()
	line, _ := block.PeekLine()
	buf := []byte{}
	if block.Peek() == '<' {
		i := 1
		for i < len(line) {
			c := line[i]
			if c == '\\' && i < len(line)-1 && util.IsPunct(line[i+1]) {
				buf = append(buf, '\\', line[i+1])
				i += 2
				continue
			} else if c == '>' {
				block.Advance(i + 1)
				return line[1:i], true
			}
			buf = append(buf, c)
			i++
		}
		return nil, false
	}
	opened := 0
	i := 0
	for i < len(line) {
		c := line[i]
		if c == '\\' && i < len(line)-1 && util.IsPunct(line[i+1]) {
			buf = append(buf, '\\', line[i+1])
			i += 2
			continue
		} else if c == '(' {
			opened++
		} else if c == ')' {
			opened--
			if opened < 0 {
				break
			}
		} else if util.IsSpace(c) {
			break
		}
		buf = append(buf, c)
		i++
	}
	block.Advance(i)
	return line[:i], len(line[:i]) != 0
}

func parseLinkTitle(block text.Reader) ([]byte, bool) {
	block.SkipSpaces()
	opener := block.Peek()
	if opener != '"' && opener != '\'' && opener != '(' {
		return nil, false
	}
	closer := opener
	if opener == '(' {
		closer = ')'
	}
	line, _ := block.PeekLine()
	pos := util.FindClosure(line[1:], opener, closer, false, true)
	if pos < 0 {
		return nil, false
	}
	pos += 2 // opener + closer
	block.Advance(pos)
	return line[1 : pos-1], true
}

func (s *linkParser) CloseBlock(parent ast.Node, block text.Reader, pc Context) {
	tlist := pc.Get(linkLabelStateKey)
	if tlist == nil {
		return
	}
	for s := tlist.(*linkLabelState); s != nil; {
		next := s.Next
		removeLinkLabelState(pc, s)
		s.Parent().ReplaceChild(s.Parent(), s, ast.NewTextSegment(s.Segment))
		s = next
	}
}
