package parser

import (
    "markee/internal/ast"
)

type DelimiterType int

const (
    DelimiterAsterisk DelimiterType = iota
    DelimiterUnderscore
    DelimiterOpenBracket
    DelimiterOpenImage
)

type Delimiter struct {
    Type          DelimiterType
    Count         int
    OriginalCount int
    IsActive      bool
    CanOpen       bool
    CanClose      bool
    ContentNode   *ast.Content
    Prev          *Delimiter
    Next          *Delimiter
}

type DelimiterStack struct {
    Top    *Delimiter
    Bottom *Delimiter
}

func NewDelimiterStack() *DelimiterStack {
    return &DelimiterStack{}
}

func (s *DelimiterStack) Push(d *Delimiter) {
    if s.Top == nil {
        s.Top = d
        s.Bottom = d
        d.Prev = nil
        d.Next = nil
    } else {
        d.Prev = s.Top
        d.Next = nil
        s.Top.Next = d
        s.Top = d
    }
}

func (s *DelimiterStack) Remove(d *Delimiter) {
    if d.Prev != nil {
        d.Prev.Next = d.Next
    } else {
        s.Bottom = d.Next
    }
    
    if d.Next != nil {
        d.Next.Prev = d.Prev
    } else {
        s.Top = d.Prev
    }
}

func (s *DelimiterStack) RemoveAbove(stackBottom *Delimiter) {
    var current *Delimiter
    if stackBottom != nil {
        current = stackBottom.Next
    } else {
        current = s.Bottom
    }
    
    for current != nil {
        next := current.Next
        s.Remove(current)
        current = next
    }
}
