package ast

type Visitor interface {
	VisitDocument(node Node)
	VisitBlockQuote(node Node)
	VisitList(node Node)
	VisitListItem(node Node)
	VisitCodeBlock(node Node)
	VisitHTMLBlock(node Node)
	VisitThematicBreak(node Node)
	VisitHeading(node Node)
	VisitParagraph(node Node)
	VisitCodeSpan(node Node)
	VisitHTMLSpan(node Node)
	VisitEmphasis(node Node)
	VisitStrong(node Node)
	VisitLink(node Node)
	VisitImage(node Node)
	VisitSoftBreak(node Node)
	VisitLineBreak(node Node)
	VisitContent(node Node)
}

type BaseVisitor struct{}

func (v *BaseVisitor) VisitDocument(node Node)      {}
func (v *BaseVisitor) VisitBlockQuote(node Node)    {}
func (v *BaseVisitor) VisitList(node Node)          {}
func (v *BaseVisitor) VisitListItem(node Node)      {}
func (v *BaseVisitor) VisitCodeBlock(node Node)     {}
func (v *BaseVisitor) VisitHTMLBlock(node Node)     {}
func (v *BaseVisitor) VisitThematicBreak(node Node) {}
func (v *BaseVisitor) VisitHeading(node Node)       {}
func (v *BaseVisitor) VisitParagraph(node Node)     {}
func (v *BaseVisitor) VisitCodeSpan(node Node)      {}
func (v *BaseVisitor) VisitHTMLSpan(node Node)      {}
func (v *BaseVisitor) VisitEmphasis(node Node)      {}
func (v *BaseVisitor) VisitStrong(node Node)        {}
func (v *BaseVisitor) VisitLink(node Node)          {}
func (v *BaseVisitor) VisitImage(node Node)         {}
func (v *BaseVisitor) VisitSoftBreak(node Node)     {}
func (v *BaseVisitor) VisitLineBreak(node Node)     {}
func (v *BaseVisitor) VisitContent(node Node)       {}

func WalkChildren(v Visitor, n Node) {
    for child := n.FirstChild(); child != nil; child = child.NextSibling() {
        child.Accept(v)
    }
}

func WalkLastChild(v Visitor, n Node) {
    if child := n.LastChild(); child != nil {
        child.Accept(v)
    }
}
