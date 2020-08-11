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

/*
	Reads contents from multiple files and stores the contents in a single buffer.
	Useful for getting the contents of multiple template files for execution.
*/
func ReadFromFiles(files ...string) ([]byte, error) {
	var buf []byte
	var size int64
	var w int

	for _, file := range files {
		fi, err := os.Stat(file)

		if err != nil {
			return nil, err
		}

		size += fi.Size()
	}

	buf = make([]byte, size)

	for _, file := range files {
		f, err := os.Open(file)

		defer f.Close()

		if err != nil {
			return nil, err
		}

		n, err := f.Read(buf[w:])

		if err != nil {
			return nil, err
		}

		w += n
	}

	return buf[0:w], nil
}

func ReadFromFile(file string) ([]byte, error) {
	fi, err := os.Stat(file)

	if err != nil {
		return nil, err
	}

	size := fi.Size()

	buf := make([]byte, size)

	f, err := os.Open(file)

	defer f.Close()

	if err != nil {
		return nil, err
	}

	n, err := f.Read(buf)

	if err != nil {
		return nil, err
	}

	return buf[0:n], nil
}

/*
func ReadFromFile(file string) ([]byte, error) {
	return ReadFromFiles(file)
}
*/
