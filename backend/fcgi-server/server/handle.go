package server

import (
	"fmt"
	"time"
	"strings"
)

const (
	sourceChan = 0
	relayChan  = 1

	source = 0
	relay = 1
)

type Handle struct {
	Path     string

	channel []chan []byte
	content bool
	cache   []byte
	err     error
	handler Handler
}

type HandleNode struct {
	left *HandleNode
	right *HandleNode

	height int
	factor int

	set *Handle
}

func (h *Handle) contentServer() {
	for {
		select {
		case c := <-h.channel[source]:
			if len(c) > 0 {
				h.cache = c
				h.content = true
			}
		default:
		}

		h.channel[relay] <- h.cache
	}
}

func NewHandle(path string, u Handler) *Handle {
	h := new(Handle)

	h.Path = path
	h.channel = make([]chan []byte, 2)
	h.content = false

	/* Unbuffered */
	h.channel[source] = make(chan []byte, 0)
	h.channel[relay] = make(chan []byte, 0)

	h.handler = u

	go h.contentServer()

	return h
}

func (h *Handle) Content() ([]byte, error) {
	select {
	case c := <-h.channel[relay]:
		if h.content == false {
			return nil, fmt.Errorf("No content available for %s", h.Path)
		} else {
			return c, nil
		}
	case <-time.After(1 * time.Second):
		return nil, fmt.Errorf("Content server timeout")
	}
}

var errlock = make(chan bool, 1)

func (h *Handle) GetErr() error {
	errlock <- true
	err := h.err
	<-errlock

	return err
}

func (h *Handle) SetErr(err error) {
	errlock <- true
	h.err = err
	<-errlock
}

func NewHandleNode(h *Handle) *HandleNode {
	n := new(HandleNode)

	if h != nil {
		n.set = h
		n.height = 1
	}

	return n
}

func (t *HandleNode) Find(path string) *Handle {
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

func (n *HandleNode) Insert(h *Handle) *HandleNode {
	c := 0

	if n == nil {
		return NewHandleNode(h)
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

func (n *HandleNode) UpdateHeight() {
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

func (n *HandleNode) Count() int {
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

func (n *HandleNode) Height() int {
	if n == nil {
		return 0
	} else {
		return n.height
	}
}

func (n *HandleNode) Balance() *HandleNode {
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

func (n *HandleNode) RotateLeft() *HandleNode {

	if n == nil {
		return nil
	}

	X := n
	Y := n.right

	/* inner child */
	T2 := Y.left

	Y.left = X
	X.right = T2

	X.UpdateHeight()
	Y.UpdateHeight()

	return Y
}

func (n *HandleNode) RotateRight() *HandleNode {

	if n == nil {
		return nil
	}

	X := n
	Y := n.left

	/* inner child */
	T2 := Y.right

	Y.right = X
	X.left = T2

	X.UpdateHeight()
	Y.UpdateHeight()

	return Y
}
