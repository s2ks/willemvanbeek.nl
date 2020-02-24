package util

import (
	"os"
)

func WriteToFile(file string, buf []byte) (int, error) {
	f, err := os.Create(file)

	defer f.Close()

	if err != nil {
		return 0, err
	}

	return f.Write(buf)
}
