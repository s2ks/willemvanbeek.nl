package server

import (
	"strings"
)

/* AVL tree for path handlers */

type HandlerNode struct {
	left *HandlerNode
	right *HandlerNode

	height int
	factor int

	set *Handler
}

func NewHandlerNode(h *Handler) *HandlerNode {
	n := new(HandlerNode)

	if h != nil {
		n.set = h
		n.height = 1
	}

	return n
}

func (t *HandlerNode) Find(path string) *Handler {
	if t == nil {
		return nil
	}

	n := strings.Compare(path, t.set.Path)

	switch(n) {
	case -1:
		return t.left.Find(path)
	case 0:
		return t.set
	case 1:
		return t.right.Find(path)
	}

	return nil
}

func (n *HandlerNode) Insert(h *Handler) *HandlerNode {
	c := 0

	if n == nil {
		return NewHandlerNode(h)
	}

	if n.set != nil {
		c = strings.Compare(h.Path, n.set.Path)
	}

	switch(c) {
	case -1:
		n.left = n.left.Insert(h)
		break
	case 1:
		n.right = n.right.Insert(h)
		break
	case 0:
		n.set = h
		break
	}

	n.UpdateHeight()
	return n.Balance()
}

func (n *HandlerNode) UpdateHeight() {
	l := n.left.Height()
	r := n.right.Height()

	if l > r {
		n.height =  1 + l
	} else {
		n.height = 1 + r
	}

	/* update balance factor */
	n.factor = r - l
}

func (n *HandlerNode) Count() int {
	count := 1
	if n == nil {
		return 0
	}

	if n.left != nil {
		count += n.left.Count()
	}

	if n.right != nil {
		count += n.right.Count()
	}

	return count
}

func (n *HandlerNode) Height() int {
	if n == nil {
		return 0
	} else {
		return n.height
	}
}

func (n *HandlerNode) Balance() *HandlerNode {
	node := n

	if n == nil {
		return nil
	}

	/* right biased*/
	if n.factor > 1 && n.right.factor < 0 {
		/* RightLeft rotation */
		n.right = n.right.RotateRight()
		node = n.RotateLeft()
	} else if n.factor > 1 {
		/* Left rotation */
		node = n.RotateLeft()
	}

	/* left biased */
	if n.factor < -1 && n.left.factor > 0 {
		/* LeftRight rotation */
		n.left = n.left.RotateLeft()
		node = n.RotateRight()
	} else if n.factor < -1 {
		/* Right rotation */
		node = n.RotateRight()
	}

	return node
}

func (n *HandlerNode) RotateLeft() *HandlerNode {

	if n == nil {
		return nil
	}

	X := n
	Y := n.right

	/* inner child */
	T2 := Y.left

	/* rotate left */
	Y.left = X
	X.right = T2

	X.UpdateHeight()
	Y.UpdateHeight()

	return Y
}

func (n *HandlerNode) RotateRight() *HandlerNode {

	if n == nil {
		return nil
	}

	X := n
	Y := n.left

	/* inner child */
	T2 := Y.right

	/* rotate right */
	Y.right = X
	X.left = T2

	X.UpdateHeight()
	Y.UpdateHeight()

	return Y
}
