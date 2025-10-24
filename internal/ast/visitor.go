package ast

type Visitor interface {
	VisitDocument(node Node) VisitStatus
	VisitBlockQuote(node Node) VisitStatus
	VisitList(node Node) VisitStatus
	VisitListItem(node Node) VisitStatus
	VisitCodeBlock(node Node) VisitStatus
	VisitHTMLBlock(node Node) VisitStatus
	VisitThematicBreak(node Node) VisitStatus
	VisitHeading(node Node) VisitStatus
	VisitParagraph(node Node) VisitStatus
	VisitCodeSpan(node Node) VisitStatus
	VisitHTMLSpan(node Node) VisitStatus
	VisitEmphasis(node Node) VisitStatus
	VisitStrong(node Node) VisitStatus
	VisitLink(node Node) VisitStatus
	VisitImage(node Node) VisitStatus
	VisitSoftBreak(node Node) VisitStatus
	VisitLineBreak(node Node) VisitStatus
	VisitContent(node Node) VisitStatus
}

type BaseVisitor struct{}

func (v *BaseVisitor) VisitDocument(node Node) VisitStatus      { return VisitContinue }
func (v *BaseVisitor) VisitBlockQuote(node Node) VisitStatus    { return VisitContinue }
func (v *BaseVisitor) VisitList(node Node) VisitStatus          { return VisitContinue }
func (v *BaseVisitor) VisitListItem(node Node) VisitStatus      { return VisitContinue }
func (v *BaseVisitor) VisitCodeBlock(node Node) VisitStatus     { return VisitContinue }
func (v *BaseVisitor) VisitHTMLBlock(node Node) VisitStatus     { return VisitContinue }
func (v *BaseVisitor) VisitThematicBreak(node Node) VisitStatus { return VisitContinue }
func (v *BaseVisitor) VisitHeading(node Node) VisitStatus       { return VisitContinue }
func (v *BaseVisitor) VisitParagraph(node Node) VisitStatus     { return VisitContinue }
func (v *BaseVisitor) VisitCodeSpan(node Node) VisitStatus      { return VisitContinue }
func (v *BaseVisitor) VisitHTMLSpan(node Node) VisitStatus      { return VisitContinue }
func (v *BaseVisitor) VisitEmphasis(node Node) VisitStatus      { return VisitContinue }
func (v *BaseVisitor) VisitStrong(node Node) VisitStatus        { return VisitContinue }
func (v *BaseVisitor) VisitLink(node Node) VisitStatus          { return VisitContinue }
func (v *BaseVisitor) VisitImage(node Node) VisitStatus         { return VisitContinue }
func (v *BaseVisitor) VisitSoftBreak(node Node) VisitStatus     { return VisitContinue }
func (v *BaseVisitor) VisitLineBreak(node Node) VisitStatus     { return VisitContinue }
func (v *BaseVisitor) VisitContent(node Node) VisitStatus       { return VisitContinue }

type VisitStatus int

const (
	VisitContinue VisitStatus = iota
    VisitLastChild
	VisitStop
)
