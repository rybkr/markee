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

	LeaveDocument(node Node)
	LeaveBlockQuote(node Node)
	LeaveList(node Node)
	LeaveListItem(node Node)
	LeaveCodeBlock(node Node)
	LeaveHTMLBlock(node Node)
	LeaveThematicBreak(node Node)
	LeaveHeading(node Node)
	LeaveParagraph(node Node)
	LeaveCodeSpan(node Node)
	LeaveHTMLSpan(node Node)
	LeaveEmphasis(node Node)
	LeaveStrong(node Node)
	LeaveLink(node Node)
	LeaveImage(node Node)
	LeaveSoftBreak(node Node)
	LeaveLineBreak(node Node)
	LeaveContent(node Node)
}

type BaseVisitor struct{}

func (v *BaseVisitor) VisitDocument(node Node) VisitStatus      { return VisitStop }
func (v *BaseVisitor) VisitBlockQuote(node Node) VisitStatus    { return VisitStop }
func (v *BaseVisitor) VisitList(node Node) VisitStatus          { return VisitStop }
func (v *BaseVisitor) VisitListItem(node Node) VisitStatus      { return VisitStop }
func (v *BaseVisitor) VisitCodeBlock(node Node) VisitStatus     { return VisitStop }
func (v *BaseVisitor) VisitHTMLBlock(node Node) VisitStatus     { return VisitStop }
func (v *BaseVisitor) VisitThematicBreak(node Node) VisitStatus { return VisitStop }
func (v *BaseVisitor) VisitHeading(node Node) VisitStatus       { return VisitStop }
func (v *BaseVisitor) VisitParagraph(node Node) VisitStatus     { return VisitStop }
func (v *BaseVisitor) VisitCodeSpan(node Node) VisitStatus      { return VisitStop }
func (v *BaseVisitor) VisitHTMLSpan(node Node) VisitStatus      { return VisitStop }
func (v *BaseVisitor) VisitEmphasis(node Node) VisitStatus      { return VisitStop }
func (v *BaseVisitor) VisitStrong(node Node) VisitStatus        { return VisitStop }
func (v *BaseVisitor) VisitLink(node Node) VisitStatus          { return VisitStop }
func (v *BaseVisitor) VisitImage(node Node) VisitStatus         { return VisitStop }
func (v *BaseVisitor) VisitSoftBreak(node Node) VisitStatus     { return VisitStop }
func (v *BaseVisitor) VisitLineBreak(node Node) VisitStatus     { return VisitStop }
func (v *BaseVisitor) VisitContent(node Node) VisitStatus       { return VisitStop }

func (v *BaseVisitor) LeaveDocument(node Node)      { return }
func (v *BaseVisitor) LeaveBlockQuote(node Node)    { return }
func (v *BaseVisitor) LeaveList(node Node)          { return }
func (v *BaseVisitor) LeaveListItem(node Node)      { return }
func (v *BaseVisitor) LeaveCodeBlock(node Node)     { return }
func (v *BaseVisitor) LeaveHTMLBlock(node Node)     { return }
func (v *BaseVisitor) LeaveThematicBreak(node Node) { return }
func (v *BaseVisitor) LeaveHeading(node Node)       { return }
func (v *BaseVisitor) LeaveParagraph(node Node)     { return }
func (v *BaseVisitor) LeaveCodeSpan(node Node)      { return }
func (v *BaseVisitor) LeaveHTMLSpan(node Node)      { return }
func (v *BaseVisitor) LeaveEmphasis(node Node)      { return }
func (v *BaseVisitor) LeaveStrong(node Node)        { return }
func (v *BaseVisitor) LeaveLink(node Node)          { return }
func (v *BaseVisitor) LeaveImage(node Node)         { return }
func (v *BaseVisitor) LeaveSoftBreak(node Node)     { return }
func (v *BaseVisitor) LeaveLineBreak(node Node)     { return }
func (v *BaseVisitor) LeaveContent(node Node)       { return }

type VisitStatus int

const (
	VisitStop VisitStatus = iota
	VisitLastChild
	VisitChildrenDFS
)
