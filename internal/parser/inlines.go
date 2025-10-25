package parser

import (
    "markee/internal/ast"
    "strings"
)

type InlineParser struct {
	input     string
	pos       int
	delims    *DelimiterStack
	container ast.Node
}

func NewInlineParser(container ast.Node, content string) *InlineParser {
	return &InlineParser{
		input:     content,
		pos:       0,
		delims:    NewDelimiterStack(),
		container: container,
	}
}

func ParseInlines(container ast.Node, content string) {
    ip := NewInlineParser(container, content)
    ip.parse()
}

func (p *InlineParser) parse() {
	for p.pos < len(p.input) {
		c := p.input[p.pos]
		switch c {
		case '\n':
			p.parseLineBreak()
		case '\\':
			p.parseBackslashEscape()
		case '`':
			p.parseBackticks()
		case '*', '_':
			p.parseDelimiterRun(c)
		case '[':
			p.parseOpenBracket()
		case '!':
			if p.peek(1) == '[' {
				p.parseOpenImage()
			} else {
				p.addText(string(c))
				p.pos++
			}
		case ']':
			p.parseCloseBracket()
		default:
			p.parseText()
		}
	}
	p.processEmphasis(nil)
}

func (p *InlineParser) peek(offset int) byte {
	pos := p.pos + offset
	if pos >= len(p.input) {
		return 0
	}
	return p.input[pos]
}

func (p *InlineParser) addText(s string) {
	textNode := ast.NewContent(s)
	p.container.AddChild(textNode)
}

func (p *InlineParser) parseText() {
	start := p.pos

	// Consume until we hit a special character
	for p.pos < len(p.input) {
		c := p.input[p.pos]
		if c == '\n' || c == '\\' || c == '`' || c == '*' || c == '_' ||
			c == '[' || c == ']' || c == '!' {
			break
		}
		p.pos++
	}

	if p.pos > start {
		p.addText(p.input[start:p.pos])
	}
}

func (p *InlineParser) parseLineBreak() {
	p.pos++ // consume \n

	// Check if previous node is text ending with spaces
	lastChild := p.container.LastChild()
	if textNode, ok := lastChild.(*ast.Content); ok {
		literal := textNode.Literal

		// Count trailing spaces
		spaces := 0
		for i := len(literal) - 1; i >= 0 && literal[i] == ' '; i-- {
			spaces++
		}

		if spaces >= 2 {
			// Hard break: remove trailing spaces, add hard break
			textNode.Literal = strings.TrimRight(literal, " ")
			p.container.AddChild(ast.NewLineBreak())
		} else {
			// Soft break
			p.container.AddChild(ast.NewSoftBreak())
		}
	} else {
		// No previous text, just soft break
		p.container.AddChild(ast.NewSoftBreak())
	}
}

func (p *InlineParser) parseBackslashEscape() {
	p.pos++ // consume \

	if p.pos < len(p.input) {
		next := p.input[p.pos]

		// Check if it's an escapable character
		if isEscapable(next) {
			p.addText(string(next))
			p.pos++
		} else if next == '\n' {
			// Backslash before newline = hard break
			p.container.AddChild(ast.NewLineBreak())
			p.pos++
		} else {
			// Not escapable, add literal backslash
			p.addText("\\")
		}
	} else {
		// Backslash at end of input
		p.addText("\\")
	}
}

func isEscapable(c byte) bool {
	return strings.ContainsRune("!\"#$%&'()*+,-./:;<=>?@[\\]^_`{|}~", rune(c))
}

func (p *InlineParser) parseBackticks() {
	start := p.pos

	// Count opening backticks
	openTicks := 0
	for p.pos < len(p.input) && p.input[p.pos] == '`' {
		openTicks++
		p.pos++
	}

	// Look for closing sequence of same length
	codeStart := p.pos

	for p.pos < len(p.input) {
		if p.input[p.pos] == '`' {
			closeTicks := 0
			closeStart := p.pos

			for p.pos < len(p.input) && p.input[p.pos] == '`' {
				closeTicks++
				p.pos++
			}

			if closeTicks == openTicks {
				// Found matching close
				code := p.input[codeStart:closeStart]

				// Strip one leading and trailing space if present
				code = strings.TrimSpace(code)
				if len(code) > 0 && code[0] == ' ' && code[len(code)-1] == ' ' {
					code = code[1 : len(code)-1]
				}

				code = collapseWhitespace(code)

				p.container.AddChild(ast.NewCodeSpan(code))
				return
			}
		} else {
			p.pos++
		}
	}

	// No matching close found, add as literal text
	p.pos = start
	p.addText(strings.Repeat("`", openTicks))
	p.pos += openTicks
}

func collapseWhitespace(s string) string {
	var result strings.Builder
	prevSpace := false

	for _, c := range s {
		if c == ' ' || c == '\t' || c == '\n' || c == '\r' {
			if !prevSpace {
				result.WriteByte(' ')
				prevSpace = true
			}
		} else {
			result.WriteRune(c)
			prevSpace = false
		}
	}

	return result.String()
}

func (p *InlineParser) parseDelimiterRun(char byte) {
	start := p.pos

	// Count the run
	count := 0
	for p.pos < len(p.input) && p.input[p.pos] == char {
		count++
		p.pos++
	}

	// Determine if can open/close based on flanking rules
	canOpen, canClose := p.checkFlanking(start, p.pos, char)

	// Add text node with the delimiter characters
	textNode := ast.NewContent(strings.Repeat(string(char), count))
	p.container.AddChild(textNode)

	// Add to delimiter stack if it can open or close
	if canOpen || canClose {
		delimType := DelimiterAsterisk
		if char == '_' {
			delimType = DelimiterUnderscore
		}

		delim := &Delimiter{
			Type:          delimType,
			Count:         count,
			OriginalCount: count,
			IsActive:      true,
			CanOpen:       canOpen,
			CanClose:      canClose,
			ContentNode:      textNode,
		}

		p.delims.Push(delim)
	}
}

func (p *InlineParser) checkFlanking(start, end int, char byte) (canOpen, canClose bool) {
	// Get character before and after the run
	var before, after rune = ' ', ' '

	if start > 0 {
		before = rune(p.input[start-1])
	}

	if end < len(p.input) {
		after = rune(p.input[end])
	}

	// Check if before/after are whitespace or punctuation
	beforeIsWhitespace := isWhitespace(before)
	afterIsWhitespace := isWhitespace(after)
	beforeIsPunctuation := isPunctuation(before)
	afterIsPunctuation := isPunctuation(after)

	// Left-flanking: not followed by whitespace, and either
	// (a) not followed by punctuation, or
	// (b) followed by punctuation and preceded by whitespace or punctuation
	leftFlanking := !afterIsWhitespace &&
		(!afterIsPunctuation || beforeIsWhitespace || beforeIsPunctuation)

	// Right-flanking: not preceded by whitespace, and either
	// (a) not preceded by punctuation, or
	// (b) preceded by punctuation and followed by whitespace or punctuation
	rightFlanking := !beforeIsWhitespace &&
		(!beforeIsPunctuation || afterIsWhitespace || afterIsPunctuation)

	if char == '_' {
		// Underscore: more restrictive
		// Can open only if left-flanking and either not right-flanking or preceded by punctuation
		canOpen = leftFlanking && (!rightFlanking || beforeIsPunctuation)
		// Can close only if right-flanking and either not left-flanking or followed by punctuation
		canClose = rightFlanking && (!leftFlanking || afterIsPunctuation)
	} else {
		// Asterisk: simpler rules
		canOpen = leftFlanking
		canClose = rightFlanking
	}

	return
}

func isWhitespace(r rune) bool {
	return r == ' ' || r == '\t' || r == '\n' || r == '\r'
}

func isPunctuation(r rune) bool {
	// CommonMark defines punctuation characters
	return strings.ContainsRune("!\"#$%&'()*+,-./:;<=>?@[\\]^_`{|}~", r)
}

// In internal/parser/inline.go

func (p *InlineParser) parseOpenBracket() {
	// Add text node with '['
	textNode := ast.NewContent("[")
	p.container.AddChild(textNode)

	// Add to delimiter stack
	delim := &Delimiter{
		Type:          DelimiterOpenBracket,
		Count:         1,
		OriginalCount: 1,
		IsActive:      true,
		CanOpen:       true,
		CanClose:      false,
		ContentNode:      textNode,
	}

	p.delims.Push(delim)
	p.pos++
}

func (p *InlineParser) parseOpenImage() {
	// Add text node with '!['
	textNode := ast.NewContent("![")
	p.container.AddChild(textNode)

	// Add to delimiter stack
	delim := &Delimiter{
		Type:          DelimiterOpenImage,
		Count:         1,
		OriginalCount: 1,
		IsActive:      true,
		CanOpen:       true,
		CanClose:      false,
		ContentNode:      textNode,
	}

	p.delims.Push(delim)
	p.pos += 2 // consume '!['
}

func (p *InlineParser) parseCloseBracket() {
	p.pos++ // consume ']'

	// Look for link or image
	opener := p.findLinkOpener()

	if opener == nil {
		// No opener found, add literal ]
		p.addText("]")
		return
	}

	if !opener.IsActive {
		// Opener is inactive, remove it and add literal ]
		p.delims.Remove(opener)
		p.addText("]")
		return
	}

	// Try to parse link/image destination and title
	isImage := opener.Type == DelimiterOpenImage

	// Try inline link: ]( ... )
	if p.peek(0) == '(' {
		if dest, title, ok := p.parseLinkDestination(); ok {
			p.createLinkOrImage(opener, dest, title, isImage)
			return
		}
	}

	// Try reference link: ][ref] or ]
	if _, ok := p.parseLinkLabel(); ok {
		// TODO: Look up reference in reference map
		// For now, just remove opener and add literal ]
		p.delims.Remove(opener)
		p.addText("]")
		return
	}

	// No valid link, remove opener and add literal ]
	p.delims.Remove(opener)
	p.addText("]")
}

func (p *InlineParser) findLinkOpener() *Delimiter {
	// Search backwards in delimiter stack for [ or ![
	current := p.delims.Top

	for current != nil {
		if current.Type == DelimiterOpenBracket || current.Type == DelimiterOpenImage {
			return current
		}
		current = current.Prev
	}

	return nil
}

func (p *InlineParser) parseLinkDestination() (dest, title string, ok bool) {
	if p.peek(0) != '(' {
		return "", "", false
	}

	p.pos++ // consume '('
	p.skipWhitespace()

	// Parse destination
	if p.peek(0) == '<' {
		// Angle-bracket enclosed destination
		dest, ok = p.parseAngleBracketDest()
		if !ok {
			return "", "", false
		}
	} else {
		// Regular destination
		dest, ok = p.parseRegularDest()
		if !ok {
			return "", "", false
		}
	}

	// Optional whitespace
	initialPos := p.pos
	p.skipWhitespace()

	// Optional title
	if p.peek(0) == '"' || p.peek(0) == '\'' || p.peek(0) == '(' {
		if t, ok := p.parseLinkTitle(); ok {
			title = t
		} else {
			// Failed to parse title, backtrack
			p.pos = initialPos
		}
	}

	// Optional whitespace
	p.skipWhitespace()

	// Must end with )
	if p.peek(0) != ')' {
		return "", "", false
	}

	p.pos++ // consume ')'
	return dest, title, true
}

func (p *InlineParser) parseAngleBracketDest() (string, bool) {
	p.pos++ // consume '<'
	start := p.pos

	for p.pos < len(p.input) {
		c := p.input[p.pos]

		if c == '>' {
			dest := p.input[start:p.pos]
			p.pos++ // consume '>'
			return dest, true
		}

		if c == '<' || c == '\n' || c == '\\' {
			return "", false
		}

		p.pos++
	}

	return "", false
}

func (p *InlineParser) parseRegularDest() (string, bool) {
	start := p.pos
	parenDepth := 0

	for p.pos < len(p.input) {
		c := p.input[p.pos]

		if c == '\\' && p.pos+1 < len(p.input) {
			// Escaped character
			p.pos += 2
			continue
		}

		if isWhitespace(rune(c)) {
			break
		}

		if c == '(' {
			parenDepth++
			p.pos++
		} else if c == ')' {
			if parenDepth == 0 {
				break
			}
			parenDepth--
			p.pos++
		} else {
			p.pos++
		}
	}

	dest := p.input[start:p.pos]
	if len(dest) == 0 {
		return "", false
	}

	return dest, true
}

func (p *InlineParser) parseLinkTitle() (string, bool) {
	delimiter := p.input[p.pos]
	closingDelim := delimiter

	if delimiter == '(' {
		closingDelim = ')'
	}

	p.pos++ // consume opening delimiter
	start := p.pos

	for p.pos < len(p.input) {
		c := p.input[p.pos]

		if c == '\\' && p.pos+1 < len(p.input) {
			// Escaped character
			p.pos += 2
			continue
		}

		if c == closingDelim {
			title := p.input[start:p.pos]
			p.pos++ // consume closing delimiter
			return title, true
		}

		p.pos++
	}

	return "", false
}

func (p *InlineParser) parseLinkLabel() (string, bool) {
	if p.peek(0) != '[' {
		return "", false
	}

	p.pos++ // consume '['
	start := p.pos

	for p.pos < len(p.input) && p.pos-start < 1000 { // Max 999 chars
		c := p.input[p.pos]

		if c == '\\' && p.pos+1 < len(p.input) {
			p.pos += 2
			continue
		}

		if c == '[' {
			return "", false // No nested brackets
		}

		if c == ']' {
			label := p.input[start:p.pos]
			p.pos++ // consume ']'
			return label, true
		}

		p.pos++
	}

	return "", false
}

func (p *InlineParser) skipWhitespace() {
	for p.pos < len(p.input) {
		c := p.input[p.pos]
		if c != ' ' && c != '\t' && c != '\n' && c != '\r' {
			break
		}
		p.pos++
	}
}

func (p *InlineParser) createLinkOrImage(opener *Delimiter, dest, title string, isImage bool) {
	// Get all nodes between opener and current position
	var inlineNodes []ast.Node

	// Find the opener's text node in the container
	openerNode := opener.ContentNode
	collectingNodes := false

	var current ast.Node
	for current = p.container.FirstChild(); current != nil; current = current.NextSibling() {
		if current == openerNode {
			collectingNodes = true
			continue // Don't include the opener itself
		}

		if collectingNodes {
			inlineNodes = append(inlineNodes, current)
		}
	}

	// Create link or image node
	var linkNode ast.Node
	if isImage {
		// For images, extract alt text from inline nodes
		altText := extractTextContent(inlineNodes)
		linkNode = ast.NewImage(dest, title, altText)
	} else {
		linkNode = ast.NewLink(dest, title)

		// Add inline nodes as children
		for _, node := range inlineNodes {
			linkNode.AddChild(node)
		}
	}

	// Remove nodes between opener and closer from container
	current = openerNode.NextSibling()
	for current != nil {
		next := current.NextSibling()
		p.container.RemoveChild(current)
		current = next
	}

	// Replace opener text node with link/image node
	p.container.ReplaceChild(openerNode, linkNode)

	// Process emphasis on the link content
	// Save current container
	savedContainer := p.container
	p.container = linkNode
	p.processEmphasis(opener)
	p.container = savedContainer

	// Remove opener from delimiter stack
	p.delims.Remove(opener)

	// If link (not image), deactivate all [ delimiters before opener
	if !isImage {
		current := p.delims.Bottom
		for current != nil && current != opener {
			if current.Type == DelimiterOpenBracket {
				current.IsActive = false
			}
			current = current.Next
		}
	}
}

func extractTextContent(nodes []ast.Node) string {
	var result strings.Builder

	for _, node := range nodes {
		if textNode, ok := node.(*ast.Content); ok {
			result.WriteString(textNode.Literal)
		}
	}

	return result.String()
}
// In internal/parser/inline.go

func (p *InlineParser) processEmphasis(stackBottom *Delimiter) {
    // Track openers_bottom for each delimiter type
    // Indexed by: [delimiter_type][length % 3][can_open]
    openersBottom := make(map[DelimiterType]map[int]map[bool]*Delimiter)
    
    // Initialize current position
    var currentPosition *Delimiter
    if stackBottom != nil {
        currentPosition = stackBottom.Next
    } else {
        currentPosition = p.delims.Bottom
    }
    
    // Repeat until we run out of potential closers
    for currentPosition != nil {
        // Move forward to find first potential closer with * or _
        for currentPosition != nil {
            if (currentPosition.Type == DelimiterAsterisk || 
                currentPosition.Type == DelimiterUnderscore) && 
                currentPosition.CanClose {
                break
            }
            currentPosition = currentPosition.Next
        }
        
        if currentPosition == nil {
            break
        }
        
        // Look back for matching opener
        opener := p.findMatchingOpener(currentPosition, stackBottom, openersBottom)
        
        if opener != nil {
            // Found a matching opener!
            
            // Determine if strong or regular emphasis
            // Strong if both >= 2, otherwise regular
            isStrong := opener.Count >= 2 && currentPosition.Count >= 2
            useCount := 1
            if isStrong {
                useCount = 2
            }
            
            // Create emphasis node
            var emphNode ast.Node
            if isStrong {
                emphNode = ast.NewStrong()
            } else {
                emphNode = ast.NewEmphasis()
            }
            
            // Move nodes between opener and closer into emphasis node
            p.moveNodesBetween(opener.ContentNode, currentPosition.ContentNode, emphNode)
            
            // Insert emphasis node after opener
            p.container.InsertAfter(opener.ContentNode, emphNode)
            
            // Remove delimiters between opener and closer
            p.removeDelimitersBetween(opener, currentPosition)
            
            // Remove useCount delimiters from opener and closer
            opener.Count -= useCount
            currentPosition.Count -= useCount
            
            // Update text nodes
            if opener.Count > 0 {
                opener.ContentNode.Literal = strings.Repeat(
                    string(opener.ContentNode.Literal[0]), 
                    opener.Count,
                )
            } else {
                // Remove empty opener
                p.container.RemoveChild(opener.ContentNode)
                p.delims.Remove(opener)
            }
            
            if currentPosition.Count > 0 {
                currentPosition.ContentNode.Literal = strings.Repeat(
                    string(currentPosition.ContentNode.Literal[0]), 
                    currentPosition.Count,
                )
            } else {
                // Remove empty closer and advance
                next := currentPosition.Next
                p.container.RemoveChild(currentPosition.ContentNode)
                p.delims.Remove(currentPosition)
                currentPosition = next
                continue
            }
            
        } else {
            // No matching opener found
            
            // Set openers_bottom
            if openersBottom[currentPosition.Type] == nil {
                openersBottom[currentPosition.Type] = make(map[int]map[bool]*Delimiter)
            }
            lengthMod3 := currentPosition.OriginalCount % 3
            if openersBottom[currentPosition.Type][lengthMod3] == nil {
                openersBottom[currentPosition.Type][lengthMod3] = make(map[bool]*Delimiter)
            }
            
            canOpen := currentPosition.CanOpen
            if currentPosition.Prev != nil {
                openersBottom[currentPosition.Type][lengthMod3][canOpen] = currentPosition.Prev
            } else {
                openersBottom[currentPosition.Type][lengthMod3][canOpen] = stackBottom
            }
            
            // If closer is not a potential opener, remove it
            if !currentPosition.CanOpen {
                next := currentPosition.Next
                p.delims.Remove(currentPosition)
                currentPosition = next
            } else {
                currentPosition = currentPosition.Next
            }
        }
    }
    
    // Remove all delimiters above stack_bottom
    p.delims.RemoveAbove(stackBottom)
}

func (p *InlineParser) findMatchingOpener(
    closer *Delimiter, 
    stackBottom *Delimiter,
    openersBottom map[DelimiterType]map[int]map[bool]*Delimiter,
) *Delimiter {
    
    // Get the openers_bottom for this delimiter type
    var bottomDelim *Delimiter
    if openersBottom[closer.Type] != nil {
        lengthMod3 := closer.OriginalCount % 3
        if openersBottom[closer.Type][lengthMod3] != nil {
            bottomDelim = openersBottom[closer.Type][lengthMod3][closer.CanOpen]
        }
    }
    
    // Start search from closer, going backwards
    current := closer.Prev
    
    for current != nil && current != stackBottom && current != bottomDelim {
        // Check if it matches
        if current.Type == closer.Type && current.CanOpen && current.IsActive {
            // Found potential match
            
            // Check the "sum of lengths is multiple of 3" rule
            // If both can open and close, and their sum is multiple of 3,
            // they don't match unless one of them is also a multiple of 3
            if (current.CanOpen && current.CanClose) || 
               (closer.CanOpen && closer.CanClose) {
                sumLength := current.OriginalCount + closer.OriginalCount
                if sumLength%3 == 0 && 
                   current.OriginalCount%3 != 0 && 
                   closer.OriginalCount%3 != 0 {
                    // Skip this opener
                    current = current.Prev
                    continue
                }
            }
            
            return current
        }
        
        current = current.Prev
    }
    
    return nil
}

func (p *InlineParser) moveNodesBetween(start, end ast.Node, dest ast.Node) {
    current := start.NextSibling()
    
    for current != nil && current != end {
        next := current.NextSibling()
        p.container.RemoveChild(current)
        dest.AddChild(current)
        current = next
    }
}

func (p *InlineParser) removeDelimitersBetween(start, end *Delimiter) {
    current := start.Next
    
    for current != nil && current != end {
        next := current.Next
        p.delims.Remove(current)
        current = next
    }
}
