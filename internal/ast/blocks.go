package ast

type Document struct{ BaseNode }

func NewDocument() *Document {
	return &Document{
		BaseNode: New(NodeDocument),
	}
}

type BlockQuote struct{ BaseNode }

func NewBlockQuote() *BlockQuote {
	return &BlockQuote{
		BaseNode: New(NodeBlockQuote),
	}
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

type ListItem struct {
    BaseNode
    Indent int
}

func NewListItem() *ListItem {
	return &ListItem{
		BaseNode: New(NodeListItem),
	}
}

type CodeBlock struct {
	BaseNode
	Language string
	IsFenced bool
	Literal  string
}

func NewCodeBlock(isFenced bool) *CodeBlock {
	return &CodeBlock{
		BaseNode: New(NodeCodeBlock),
		IsFenced: isFenced,
	}
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

type ThematicBreak struct{ BaseNode }

func NewThematicBreak() *ThematicBreak {
	return &ThematicBreak{
		BaseNode: New(NodeThematicBreak),
	}
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

type Paragraph struct{ BaseNode }

func NewParagraph() *Paragraph {
	return &Paragraph{
		BaseNode: New(NodeParagraph),
	}
}
