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

func (w *Writer) setDefaults() {
	if w.Endianness == nil {
		w.Endianness = binary.LittleEndian
	}
}

// Write writes a byte slice to w.
func (w *Writer) Write(b []byte) {
	if w.err != nil {
		return
	}
	w.check(w.W.Write(b))
}

func (w *Writer) WriteUint8(u uint8) {
	w.scratch[0] = u
	w.Write(w.scratch[:1])
}

func (w *Writer) WriteUint16(u uint16) {
	w.setDefaults()
	w.Endianness.PutUint16(w.scratch[:], u)
	w.Write(w.scratch[:2])
}

func (w *Writer) WriteUint32(u uint32) {
	w.setDefaults()
	w.Endianness.PutUint32(w.scratch[:], u)
	w.Write(w.scratch[:4])
}

func (w *Writer) WriteUint64(u uint64) {
	w.setDefaults()
	w.Endianness.PutUint64(w.scratch[:], u)
	w.Write(w.scratch[:8])
}

func (w *Writer) WriteInt8(u int8) {
	w.WriteUint8(uint8(u))
}

func (w *Writer) WriteInt16(u int16) {
	w.WriteUint16(uint16(u))
}

func (w *Writer) WriteInt32(u int32) {
	w.WriteUint32(uint32(u))
}

func (w *Writer) WriteInt64(u int64) {
	w.WriteUint64(uint64(u))
}

// N returns the number of bytes successfully written.
func (w *Writer) N() int {
	return w.n
}

// Err returns the last non-nil error.
func (w *Writer) Err() error {
	return w.err
}

// Reset resets the internal byte count and error.
func (w *Writer) Reset() {
	w.n = 0
	w.err = nil
}
