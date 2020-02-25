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
}

type HandleNode struct {
	left *HandleNode
	right *HandleNode
	parent *HandleNode

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

func NewHandle(path string) *Handle {
	h := new(Handle)

	h.Path = path
	h.channel = make([]chan []byte, 2)
	h.content = false

	/* Unbuffered */
	h.channel[source] = make(chan []byte, 0)
	h.channel[relay] = make(chan []byte, 0)

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

func NewHandleNode(parent *HandleNode) *HandleNode {
	t := new(HandleNode)
	t.parent = parent

	return t
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

func (t *HandleNode) Insert(h *Handle) {
	if t.set == nil {
		t.set = h
		return
	}

	n := strings.Compare(h.Path, t.set.Path)

	switch(n) {
	case -1:
		if t.left == nil {
			t.left = NewHandleNode(t)
		}

		t.left.Insert(h)
		break
	case 0:
		t.set = h
		break
	case 1:
		if t.right == nil {
			t.right = NewHandleNode(t)
		}
		t.right.Insert(h)
		break
	}

	return
}

func (t *HandleNode) Count() (n int) {
	n = 0
	if t == nil {
		return
	}

	/* root node */
	if t.parent == nil {
		n += 1
	}

	if t.left != nil {
		n += 1 + t.left.Count()
	}

	if t.right != nil {
		n += 1 + t.right.Count()
	}

	return
}

func (t *HandleNode) Root() *HandleNode {
	root := t

	for root.parent != nil {
		root = root.parent
	}

	return root
}

func (t *HandleNode) Balance() *HandleNode {
	var lc int
	var rc int

	var node *HandleNode

	node = t

	rc = node.right.Count()
	lc = node.left.Count()

	delta := lc - rc

	/* right biased */
	if delta < -1 {
		node.right.Balance()
		if node.RotateLeft() {
			node = node.parent
		}
	}

	/* left biased */
	if delta > 1 {
		node.left.Balance()

		if node.RotateRight() {
			node = node.parent
		}
	}

	return node
}

/*
		A*
	       / \
	      B   C
	     / \ / \
	    1  2 3  4

	  Rotate left
	        C
	       / \
	      A*  4
	     / \
	    B   3
	   / \
	  1   2


	  Rotate right
	  	B
	       / \
	      1   A*
	         / \
		2   C
		   / \
		  3   4
*/

func (t *HandleNode) RotateLeft() bool {
	A := t

	if A == nil {
		return false
	}

	C := t.right

	if C == nil {
		return false
	}

	P := A.parent

	/* A is not the root node */
	if P != nil {
		if P.left == A {
			P.left = C
		}
		if P.right == A {
			P.right = C
		}
	}

	A.right = C.left
	C.parent = A.parent
	A.parent = C
	C.left = A

	return true
}

func (t *HandleNode) RotateRight() bool {
	A := t

	if A == nil {
		return false
	}

	B := t.left

	if B == nil {
		return false
	}

	P := A.parent

	/* A is not the root node */
	if P != nil {
		if P.left == A {
			P.left = B
		}

		if P.right == A {
			P.right = B
		}
	}

	A.left = B.right
	B.parent = A.parent
	A.parent = B
	B.right = A

	return true
}
