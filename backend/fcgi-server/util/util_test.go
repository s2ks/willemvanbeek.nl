package util

import (
	"testing"
	"fmt"
)

func TestSub1(t *testing.T) {
	str, err := Substitute("%WORD% World", "WORD", "Hello", '%')
	fmt.Println(str)
	if str != "Hello World" {
		t.Errorf("TestSub1 failed %v", err)
	}
}

func TestSub2(t *testing.T) {
	keyval := make(map[string]string)

	keyval["KEY"] = "Hello"
	keyval["KEY2"] = "And goodbye"
	str, err := SubstituteMap("%KEY%, %KEY%, %KEY%, world. %KEY2%.", keyval, "%")

	fmt.Println(str)

	if str != "Hello, Hello, Hello, world. And goodbye." {
		t.Errorf("TestSub2 failed %v", err)
	}
}

func TestSub3(t *testing.T) {
	keyval := make(map[string]string)

	keyval["KEY"] = "Hello"
	keyval["KEY2"] = "And goodbye"
	str, err := SubstituteMap("%KEY3%, %kEY%, %KEY% KEY% %KEY, world. %KEY2% KEY3.", keyval, "%")

	fmt.Println(str)

	if str != "%KEY3%, %kEY%, Hello KEY% %KEY, world. And goodbye KEY3." {
		t.Errorf("TestSub3 failed %v", err)
	}
}


func TestSub4(t *testing.T) {
	keyval := make(map[string]string)

	raw := "%KEY%, world. %KEY2%"

	keyval["KEY6"] = "Hello"
	keyval["KEY5"] = "And goodbye"
	str, err := SubstituteMap(raw, keyval, "%")

	fmt.Println(str)
	fmt.Println(err)

	if raw != str && err != nil {
		t.Errorf("TestSub4 failed")
	}
}

