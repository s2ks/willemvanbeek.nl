package util

type Buffer struct {
	raw []byte
	size int
}

func (b *Buffer) Write(p []byte) (int, error) {
	if b.raw == nil {
		b.raw = make([]byte, 0)
		b.size = 0
	}

	buf := make([]byte, b.size + len(p))


	w := 0
	w += copy(buf[w:], b.raw)
	w += copy(buf[w:], p)

	b.raw = buf
	b.size = w

	return w, nil
}

func (b *Buffer) Bytes() []byte {
	return b.raw
}

func (b *Buffer) Size() int {
	return b.size
}
