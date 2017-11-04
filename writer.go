package bin

import (
	"encoding/binary"
	"io"
)

// Writer wraps an io.Writer.
// The last non-nil error is located in Err().
// The total amount of bytes written is located in N().
type Writer struct {
	W io.Writer

	// Endianness will default to binary.LittleEndian.
	Endianness binary.ByteOrder

	scratch [8]byte

	err error
	n   int
}

func (w *Writer) check(n int, err error) {
	w.n += n
	if w.err == nil && err != nil {
		w.err = err
	}
}

// Write writes a byte slice to w.
func (w *Writer) Write(b []byte) {
	w.check(w.W.Write(b))
}

func (w *Writer) WriteUint8(u uint8) {
	w.W.Write
}

// N returns the number of bytes successfully written.
func (w *Writer) N() int {
	return w.n
}

// Reset resets the internal byte written count and error.
func (w *Writer) Reset() {
	w.n = 0
	w.err = nil
}
