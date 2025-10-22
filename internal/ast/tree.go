package ast

func (n *Node) AppendChild(child *Node) {
    child.Parent = n
    if n.LastChild != nil {
        n.LastChild.Next = child
        child.Prev = n.LastChild
        n.LastChild = child
    } else {
        n.FirstChild = child
        n.LastChild = child
    }
}

func (n *Node) InsertAfter(sibling *Node) {
    sibling.Parent = n.Parent
    sibling.Prev = n
    sibling.Next = n.Next
    if n.Next != nil {
        n.Next.Prev = sibling
    }
    n.Next = sibling
    if n.Parent != nil && n.Parent.LastChild == n {
        n.Parent.LastChild = sibling
    }
}

func (n *Node) Unlink() {
    if n.Prev != nil {
        n.Prev.Next = n.Next
    } else if n.Parent != nil {
        n.Parent.FirstChild = n.next
    }

    if n.Next != nil {
        n.Next.Prev = n.Prev
    } else if n.Parent != nil {
        n.Parent.LastChild = n.Prev
    }

    n.Parent = nil
    n.Next = nil
    n.Prev = nil
}
