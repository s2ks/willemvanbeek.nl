package server

import (
	"testing"
	"fmt"
)

func newHandle(p string) *Handler {
	return NewHandler(p, nil)
}

func TestBst1(t *testing.T) {
	tree := NewHandlerNode(nil)

	tree = tree.Insert(newHandle("/"))
	tree = tree.Insert(newHandle("/"))
	tree = tree.Insert(newHandle("/"))
	tree = tree.Insert(newHandle("/"))
	tree = tree.Insert(newHandle("/"))
	tree = tree.Insert(newHandle("/"))

	PrintNode(tree)
	c := tree.Count()
	if c != 1 {
		t.Errorf("Assertion tree.Count() == 1 failed -- got %d", c)
	}


}

func TestBst2(t *testing.T) {
	tree := NewHandlerNode(nil)

	tree = tree.Insert(newHandle("/"))
	tree = tree.Insert(newHandle("/index"))
	tree = tree.Insert(newHandle("/aaaaaaaa"))
	tree = tree.Insert(newHandle("/hello/world"))
	tree = tree.Insert(newHandle("/path/to/some/file.file"))
	tree = tree.Insert(newHandle("/hello2"))

	tree = tree

	c := tree.Count()

	if c != 6 {
		t.Errorf("Assertion tree.Count() == 6 failed -- got %d", c)
	}

	PrintNode(tree)
}

func TestBst3(t *testing.T) {
	tree := NewHandlerNode(nil)

	tree = tree.Insert(newHandle("/"))
	tree = tree.Insert(newHandle("/"))
	tree = tree.Insert(newHandle("/hello2"))
	tree = tree.Insert(newHandle("/hello2"))
	tree = tree.Insert(newHandle("/index"))
	tree = tree.Insert(newHandle("/index"))
	tree = tree.Insert(newHandle("/index"))
	tree = tree.Insert(newHandle("/"))
	tree = tree.Insert(newHandle("/index"))
	tree = tree.Insert(newHandle("/aaaaaaaa"))
	tree = tree.Insert(newHandle("/hello/world"))
	tree = tree.Insert(newHandle("/path/to/some/file.file"))
	tree = tree.Insert(newHandle("/path/to/some/file.file"))
	tree = tree.Insert(newHandle("/hello2"))

	c := tree.Count()

	if c != 6 {
		t.Errorf("Assertion tree.Count() == 6 failed -- got %d", c)
	}

	PrintNode(tree)
}

func TestBst5(t *testing.T) {
	tree := NewHandlerNode(nil)

	tree = tree.Insert(newHandle("/"))
	tree = tree.Insert(newHandle("/1"))
	tree = tree.Insert(newHandle("/2"))
	tree = tree.Insert(newHandle("/3"))
	tree = tree.Insert(newHandle("/index"))
	tree = tree.Insert(newHandle("/index1"))
	tree = tree.Insert(newHandle("/index2"))
	tree = tree.Insert(newHandle("/index3"))
	tree = tree.Insert(newHandle("/index4"))
	tree = tree.Insert(newHandle("/aaaaaaaa"))
	tree = tree.Insert(newHandle("/hello/world"))
	tree = tree.Insert(newHandle("/path/to/some/file.file"))
	tree = tree.Insert(newHandle("/path/to/some/file.html"))
	tree = tree.Insert(newHandle("/path/to/some/file.php"))
	tree = tree.Insert(newHandle("/path/to/some/file.jpg"))
	tree = tree.Insert(newHandle("/hello2"))
	tree = tree.Insert(newHandle("/hello2/greeter1"))
	tree = tree.Insert(newHandle("/hello2/greeter1.cgi"))
	tree = tree.Insert(newHandle("/hello2/greeter_form"))
	tree = tree.Insert(newHandle("/hello2/goodbye"))

	c := tree.Count()

	if c != 20 {
		t.Errorf("Assertion tree.Count() == 20 failed -- got %d", c)
	}
	PrintNode(tree)
}

func TestBst7(t *testing.T) {
	tree := NewHandlerNode(nil)

	tree = tree.Insert(newHandle("/"))
	tree = tree.Insert(newHandle("/test123"))

	h1 := tree.Find("/")
	h2 := tree.Find("//")
	h3 := tree.Find("/test123")

	if h1 == nil || h2 != nil || h3 == nil {
		t.Errorf("Failed tree.Find() test")
	}

}

/* PrintNode helper */
func tree(prefix string, t *HandlerNode) {
	if t.left != nil && t.right != nil {
		fmt.Printf("%s%s%s (left) (%d)\n", prefix, "├── ", t.left.set.Path, t.left.factor)
		tree(fmt.Sprintf("%s%s", prefix, "│   "), t.left)

		fmt.Printf("%s%s%s (right) (%d)\n", prefix, "└── ", t.right.set.Path, t.right.factor)
		tree(fmt.Sprintf("%s%s", prefix, "│   "), t.right)
	} else {
		if t.left != nil {
			fmt.Printf("%s%s%s (left) (%d)\n", prefix, "└── ", t.left.set.Path, t.left.factor)
			tree(fmt.Sprintf("%s%s", prefix, "    "), t.left)
			return
		}
		if t.right != nil {
			fmt.Printf("%s%s%s (right) (%d)\n", prefix, "└── ", t.right.set.Path, t.right.factor)
			tree(fmt.Sprintf("%s%s", prefix, "    "), t.right)
			return
		}
	}
}

func PrintNode(t *HandlerNode) {
	if t.set != nil {
		fmt.Printf("%s (%d)\n", t.set.Path, t.factor)
		tree("", t)
	}
}

