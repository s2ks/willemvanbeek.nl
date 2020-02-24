package util

type Buffer struct {
	buf []byte
	siz int64
}

func (b *Buffer) Write(p []byte) (n int, err error) {
	if b.buf == nil {
		b.buf = make([]byte, 0)
		b.siz = 0
	}

	n = 0
	err = nil

	size := len(p)
	for i := 0; i < size; i++ {
		b.buf = append(b.buf, p[i])
		b.siz++
		n++
	}

	return
}

func (b *Buffer) Bytes() []byte {
	return b.buf
}

func (b *Buffer) Size() int64 {
	return b.siz
}
