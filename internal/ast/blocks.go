package ast

type Document struct{ BaseNode }

func NewDocument() *Document {
	return &Document{
		BaseNode: New(NodeDocument),
	}
}

func (d *Document) Accept(v Visitor) {
	v.VisitDocument(d)
}

type BlockQuote struct{ BaseNode }

func NewBlockQuote() *BlockQuote {
	return &BlockQuote{
		BaseNode: New(NodeBlockQuote),
	}
}

func (b *BlockQuote) Accept(v Visitor) {
	v.VisitBlockQuote(b)
}

type List struct {
	BaseNode
	IsOrdered bool
	IsTight   bool
	StartNum  int
	Delimiter byte
}

func NewList(isOrdered bool) *List {
	return &List{
		BaseNode:  New(NodeList),
		IsOrdered: isOrdered,
		IsTight:   true,
		StartNum:  1,
	}
}

func (l *List) Accept(v Visitor) {
	v.VisitList(l)
}

type ListItem struct {
	BaseNode
	Indent int
}

func NewListItem(indent int) *ListItem {
	return &ListItem{
		BaseNode: New(NodeListItem),
		Indent:   indent,
	}
}

func (l *ListItem) Accept(v Visitor) {
	v.VisitListItem(l)
}

type CodeBlock struct {
	BaseNode
	Language  string
	IsFenced  bool
	FenceChar byte
}

func NewCodeBlock(isFenced bool) *CodeBlock {
	return &CodeBlock{
		BaseNode: New(NodeCodeBlock),
		IsFenced: isFenced,
	}
}

func (c *CodeBlock) Accept(v Visitor) {
	v.VisitCodeBlock(c)
}

type HTMLBlock struct {
	BaseNode
	Literal string
}

func NewHTMLBlock(literal string) *HTMLBlock {
	return &HTMLBlock{
		BaseNode: New(NodeHTMLBlock),
		Literal:  literal,
	}
}

func (h *HTMLBlock) Accept(v Visitor) {
	v.VisitHTMLBlock(h)
}

type ThematicBreak struct{ BaseNode }

func NewThematicBreak() *ThematicBreak {
	return &ThematicBreak{
		BaseNode: New(NodeThematicBreak),
	}
}

func (t *ThematicBreak) Accept(v Visitor) {
	v.VisitThematicBreak(t)
}

type Heading struct {
	BaseNode
	Level int
}

func NewHeading(level int) *Heading {
	return &Heading{
		BaseNode: New(NodeHeading),
		Level:    level,
	}
}

func (h *Heading) Accept(v Visitor) {
	v.VisitHeading(h)
}

type Paragraph struct{ BaseNode }

func NewParagraph() *Paragraph {
	return &Paragraph{
		BaseNode: New(NodeParagraph),
	}
}

func (p *Paragraph) Accept(v Visitor) {
	v.VisitParagraph(p)
}

type BlankLine struct{ BaseNode }

func NewBlankLine() *BlankLine {
	return &BlankLine{
		BaseNode: New(NodeBlankLine),
	}
}

func (b *BlankLine) Accept(v Visitor) {
	v.VisitBlankLine(b)
}
