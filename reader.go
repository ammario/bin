package bin

import (
	"encoding/binary"
	"io"
)

// Reader wraps an io.Reader.
// The last non-nil error is located in Err().
// The total amount of bytes read is located in N().
type Reader struct {
	R io.Reader

	// Endianness will default to binary.LittleEndian.
	Endianness binary.ByteOrder

	scratch [8]byte

	err error
	n   int
}

func (rd *Reader) check(n int, err error) {
	rd.n += n
	if rd.err == nil && err != nil {
		rd.err = err
	}
}

func (rd *Reader) setDefaults() {
	if rd.Endianness == nil {
		rd.Endianness = binary.LittleEndian
	}
}

func (rd *Reader) Read(p []byte) int {
	if rd.err != nil {
		return 0
	}
	n, err := rd.R.Read(p)
	rd.check(n, err)
	return n
}

func (rd *Reader) ReadUint8(u *uint8) {
	rd.Read(rd.scratch[:1])
	if rd.err == nil {
		*u = rd.scratch[0]
	}
}

func (rd *Reader) ReadUint16(u *uint16) {
	rd.Read(rd.scratch[:2])
	if rd.err == nil {
		rd.setDefaults()
		*u = rd.Endianness.Uint16(rd.scratch[:])
	}
}

func (rd *Reader) ReadUint32(u *uint32) {
	rd.Read(rd.scratch[:4])
	if rd.err == nil {
		rd.setDefaults()
		*u = rd.Endianness.Uint32(rd.scratch[:])
	}
}

func (rd *Reader) ReadUint64(u *uint64) {
	rd.Read(rd.scratch[:8])
	if rd.err == nil {
		rd.setDefaults()
		*u = rd.Endianness.Uint64(rd.scratch[:])
	}
}

func (rd *Reader) ReadInt8(u *int8) {
	rd.Read(rd.scratch[:1])
	if rd.err == nil {
		*u = int8(rd.scratch[0])
	}
}

func (rd *Reader) ReadInt16(u *int16) {
	rd.Read(rd.scratch[:2])
	if rd.err == nil {
		*u = int16(rd.Endianness.Uint16(rd.scratch[:]))
	}
}

func (rd *Reader) ReadInt32(u *int32) {
	rd.Read(rd.scratch[:4])
	if rd.err == nil {
		*u = int32(rd.Endianness.Uint32(rd.scratch[:]))
	}
}

func (rd *Reader) ReadInt64(u *int64) {
	rd.Read(rd.scratch[:8])
	if rd.err == nil {
		*u = int64(rd.Endianness.Uint64(rd.scratch[:]))
	}
}

// N returns the number of bytes successfully read.
func (rd *Reader) N() int {
	return rd.n
}

// Err returns the last non-nil error.
func (rd *Reader) Err() error {
	return rd.err
}

// Reset resets the internal byte count and error.
func (rd *Reader) Reset() {
	rd.n = 0
	rd.err = nil
}
