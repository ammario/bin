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
// It returns the amount of bytes written by this call.
func (w *Writer) Write(b []byte) int {
	if w.err != nil {
		return 0
	}
	n, err := w.W.Write(b)
	w.check(n, err)
	return n
}

func (w *Writer) Uint8(u uint8) {
	w.scratch[0] = u
	w.Write(w.scratch[:1])
}

func (w *Writer) Uint16(u uint16) {
	w.setDefaults()
	w.Endianness.PutUint16(w.scratch[:], u)
	w.Write(w.scratch[:2])
}

func (w *Writer) Uint32(u uint32) {
	w.setDefaults()
	w.Endianness.PutUint32(w.scratch[:], u)
	w.Write(w.scratch[:4])
}

func (w *Writer) Uint64(u uint64) {
	w.setDefaults()
	w.Endianness.PutUint64(w.scratch[:], u)
	w.Write(w.scratch[:8])
}

func (w *Writer) Uvarint(u uint64) {
	w.Write(w.scratch[:binary.PutUvarint(w.scratch[:], u)])
}

func (w *Writer) Varint(u int64) {
	w.Write(w.scratch[:binary.PutVarint(w.scratch[:], u)])
}

func (w *Writer) Int8(u int8) {
	w.Uint8(uint8(u))
}

func (w *Writer) Int16(u int16) {
	w.Uint16(uint16(u))
}

func (w *Writer) Int32(u int32) {
	w.Uint32(uint32(u))
}

func (w *Writer) Int64(u int64) {
	w.Uint64(uint64(u))
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
