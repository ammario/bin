package bin

import "io"

type byteReader struct {
	R io.Reader
	b [1]byte
}

func (b *byteReader) ReadByte() (byte, error) {
	_, err := b.R.Read(b.b[:])
	if err != nil {
		return 0, err
	}
	return b.b[0], nil
}
