package util

import (
	"strings"
	"fmt"
)

func ByteSubstituteMap(raw []byte, keyval map[string]string, delim string) ([]byte, error) {
	var err error

	buf := make([]byte, len(raw))
	copy(buf, raw)

	for key, val := range keyval {
		pattern := make([]byte, len(key) + (2 * len(delim)))
		w := copy(pattern[0:], []byte(delim))
		w += copy(pattern[w:], []byte(key))
		w += copy(pattern[w:], []byte(delim))

		keycount := strings.Count(string(buf), string(pattern))

		if keycount < 1 {
			err = fmt.Errorf("Pattern: \"%v\" not found in string: \"%v\"", string(pattern), string(buf))
			continue
		}

		dest := make([]byte, len(buf) + (keycount * (len(val) - w)))

		for i := 0; i < keycount; i++ {
			keystart := strings.Index(string(buf), string(pattern))
			keyend := keystart + len(pattern)


			start := buf[0:keystart]
			end := buf[keyend:]

			w = copy(dest[0:], start)
			w += copy(dest[w:], val)
			w += copy(dest[w:], end)

			buf = make([]byte, w)
			copy(buf[0:], dest[0:w])
		}
	}

	return buf, err
}

func SubstituteMap(str string, keyval map[string]string, delim string) (string, error) {
	var err error
	var buf []byte
	raw := []byte(str)

	buf, err = ByteSubstituteMap(raw, keyval, delim)

	return string(buf), err
}

func Substitute(str, key, val string, delim byte) (string, error) {
	raw := []byte(str)
	rawlen := len(raw)

	dest := make([]byte, rawlen + (len(val) - len(key)))
	written := 0

	keystart := -1
	keyend := -1

	for i := 0; i < rawlen; i++ {
		if raw[i] == delim && keystart < 0 {
			keystart = i
		} else if raw[i] == delim {
			keyend = i
			break
		}
	}

	if keystart < 0 {
		return str, fmt.Errorf("No variable found")
	} else if keyend < 0 {
		return str, fmt.Errorf("Syntax error at: \"%v\"", str)
	}

	start := raw[0:keystart]
	end := raw[keyend + 1:]
	keyname := raw[keystart + 1:keyend]
	keyval := []byte(val)

	if key == string(keyname) {
		written += copy(dest[written:], start)
		written += copy(dest[written:], keyval)
		written += copy(dest[written:], end)
		return string(dest[0:written]), nil
	} else {
		return str, fmt.Errorf("Unknown variable \"%%%v%%\"", keyname)
	}
}
