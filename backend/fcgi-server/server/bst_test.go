package server

import (
	"testing"
	"fmt"
)

func TestBst1(t *testing.T) {
	tree := NewHandleNode(nil)

	tree.Insert(NewHandle("/"))
	tree.Insert(NewHandle("/"))
	tree.Insert(NewHandle("/"))
	tree.Insert(NewHandle("/"))
	tree.Insert(NewHandle("/"))
	tree.Insert(NewHandle("/"))

	PrintNode(tree)
	c := tree.Count()
	if c != 1 {
		t.Errorf("Assertion tree.Count() == 1 failed -- got %d", c)
	}


}

func TestBst2(t *testing.T) {
	tree := NewHandleNode(nil)

	tree.Insert(NewHandle("/"))
	tree.Insert(NewHandle("/index"))
	tree.Insert(NewHandle("/aaaaaaaa"))
	tree.Insert(NewHandle("/hello/world"))
	tree.Insert(NewHandle("/path/to/some/file.file"))
	tree.Insert(NewHandle("/hello2"))

	c := tree.Count()

	if c != 6 {
		t.Errorf("Assertion tree.Count() == 6 failed -- got %d", c)
	}

	PrintNode(tree)
}

func TestBst3(t *testing.T) {
	tree := NewHandleNode(nil)

	tree.Insert(NewHandle("/"))
	tree.Insert(NewHandle("/"))
	tree.Insert(NewHandle("/"))
	tree.Insert(NewHandle("/index"))
	tree.Insert(NewHandle("/index"))
	tree.Insert(NewHandle("/index"))
	tree.Insert(NewHandle("/index"))
	tree.Insert(NewHandle("/aaaaaaaa"))
	tree.Insert(NewHandle("/hello/world"))
	tree.Insert(NewHandle("/path/to/some/file.file"))
	tree.Insert(NewHandle("/path/to/some/file.file"))
	tree.Insert(NewHandle("/hello2"))
	tree.Insert(NewHandle("/hello2"))
	tree.Insert(NewHandle("/hello2"))

	c := tree.Count()

	if c != 6 {
		t.Errorf("Assertion tree.Count() == 6 failed -- got %d", c)
	}

	PrintNode(tree)
}

func TestBst4(t *testing.T) {
	tree := NewHandleNode(nil)

	tree.Insert(NewHandle("/"))
	tree.Insert(NewHandle("/index"))
	tree.Insert(NewHandle("/aaaaaaaa"))
	tree.Insert(NewHandle("/hello/world"))
	tree.Insert(NewHandle("/path/to/some/file.file"))
	tree.Insert(NewHandle("/hello2"))

	c := tree.Count()

	if c != 6 {
		t.Errorf("Assertion tree.Count() == 6 failed -- got %d", c)
	}

	fmt.Printf("Before balance:\n")
	PrintNode(tree)

	tree = tree.Balance()
	fmt.Printf("After balance:\n")
	PrintNode(tree)

	c = tree.Count()

	if c != 6 {
		t.Errorf("Assertion tree.Count() == 6 failed -- got %d", c)
	}

}

func TestBst5(t *testing.T) {
	tree := NewHandleNode(nil)

	tree.Insert(NewHandle("/"))
	tree.Insert(NewHandle("/1"))
	tree.Insert(NewHandle("/2"))
	tree.Insert(NewHandle("/3"))
	tree.Insert(NewHandle("/index"))
	tree.Insert(NewHandle("/index1"))
	tree.Insert(NewHandle("/index2"))
	tree.Insert(NewHandle("/index3"))
	tree.Insert(NewHandle("/index4"))
	tree.Insert(NewHandle("/aaaaaaaa"))
	tree.Insert(NewHandle("/hello/world"))
	tree.Insert(NewHandle("/path/to/some/file.file"))
	tree.Insert(NewHandle("/path/to/some/file.html"))
	tree.Insert(NewHandle("/path/to/some/file.php"))
	tree.Insert(NewHandle("/path/to/some/file.jpg"))
	tree.Insert(NewHandle("/hello2"))
	tree.Insert(NewHandle("/hello2/greeter1"))
	tree.Insert(NewHandle("/hello2/greeter1.cgi"))
	tree.Insert(NewHandle("/hello2/greeter_form"))
	tree.Insert(NewHandle("/hello2/goodbye"))

	c := tree.Count()

	if c != 20 {
		t.Errorf("Assertion tree.Count() == 20 failed -- got %d", c)
	}

	fmt.Printf("Before balance:\n")
	PrintNode(tree)

	tree = tree.Balance()
	fmt.Printf("After balance:\n")
	PrintNode(tree)

	c = tree.Count()

	if c != 20 {
		t.Errorf("Assertion tree.Count() == 20 failed -- got %d", c)
	}

}

func TestBst6(t *testing.T) {
	tree := NewHandleNode(nil)

	tree.Insert(NewHandle("/"))
	tree.Insert(NewHandle("/contact"))
	tree.Insert(NewHandle("/beelden"))
	tree.Insert(NewHandle("/beelden2"))
	tree.Insert(NewHandle("/beelden/steen"))
	tree.Insert(NewHandle("/beelden/hout"))
	tree.Insert(NewHandle("/beelden/metaal"))
	tree.Insert(NewHandle("/login"))
	tree.Insert(NewHandle("/add"))
	tree.Insert(NewHandle("/delete"))
	tree.Insert(NewHandle("/edit"))
	tree.Insert(NewHandle("/index"))

	c := tree.Count()

	if c != 12 {
		t.Errorf("Assertion tree.Count() == 12 failed -- got %d", c)
	}

	fmt.Printf("Before balance:\n")
	PrintNode(tree)

	tree = tree.Balance()
	fmt.Printf("After balance:\n")
	PrintNode(tree)

	c = tree.Count()

	if c != 12 {
		t.Errorf("Assertion tree.Count() == 12 failed -- got %d", c)
	}

}

func TestBst7(t *testing.T) {
	tree := NewHandleNode(nil)

	tree.Insert(NewHandle("/"))
	tree.Insert(NewHandle("/contact"))
	tree.Insert(NewHandle("/beelden"))
	tree.Insert(NewHandle("/beelden2"))
	tree.Insert(NewHandle("/beelden/steen"))
	tree.Insert(NewHandle("/beelden/hout"))
	tree.Insert(NewHandle("/beelden/metaal"))
	tree.Insert(NewHandle("/login"))
	tree.Insert(NewHandle("/add"))
	tree.Insert(NewHandle("/delete"))
	tree.Insert(NewHandle("/edit"))
	tree.Insert(NewHandle("/index"))

	tree = tree.Balance()


	h1 := tree.Find("/")
	h2 := tree.Find("/beelden/s")
	h3 := tree.Find("/i-do-not-exist")
	h4 := tree.Find("hello world")
	h5 := tree.Find("//")
	h6 := tree.Find("/contact/")
	h7 := tree.Find("////beelden//////steen")
	h8 := tree.Find("beelden/steen")
	h9 := tree.Find("/ /login")
	ha := tree.Find("/login")
	hb := tree.Find("/beelden/steen")
	hc := tree.Find("/beelden/steen/")

	if h1 == nil || ha == nil || hb == nil {
		t.Errorf("Failed tree.Find() test")
	}

	if h2 != nil || h3 != nil || h4 != nil || h5 != nil || h6 != nil || h7 != nil || h8 != nil ||
	h9 != nil || hc != nil {
		t.Errorf("Failed tree.Find() test")
	}


}

/* PrintNode helper */
func tree(prefix string, t *HandleNode) {
	if t.left != nil && t.right != nil {
		fmt.Printf("%s%s%s (left)\n", prefix, "├── ", t.left.set.Path)
		tree(fmt.Sprintf("%s%s", prefix, "│   "), t.left)

		fmt.Printf("%s%s%s (right)\n", prefix, "└── ", t.right.set.Path)
		tree(fmt.Sprintf("%s%s", prefix, "│   "), t.right)
	} else {
		if t.left != nil {
			fmt.Printf("%s%s%s (left)\n", prefix, "└── ", t.left.set.Path)
			tree(fmt.Sprintf("%s%s", prefix, "    "), t.left)
			return
		}
		if t.right != nil {
			fmt.Printf("%s%s%s (right)\n", prefix, "└── ", t.right.set.Path)
			tree(fmt.Sprintf("%s%s", prefix, "    "), t.right)
			return
		}
	}
}

func PrintNode(t *HandleNode) {
	if t.set != nil {
		fmt.Printf("%s\n", t.set.Path)
		tree("", t)
	}
}

