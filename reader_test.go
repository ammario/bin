package bin_test

import (
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"io"
	"testing"

	"github.com/ammario/bin"
	"github.com/stretchr/testify/assert"
)

func TestReader(t *testing.T) {
	t.Run("Default little endian", func(t *testing.T) {
		buf := &bytes.Buffer{}
		rd := &bin.Reader{R: buf}
		var u uint16
		buf.Write([]byte{0x11, 0x22})
		rd.Uint16(&u)

		assert.Equal(t, uint16(0x2211), u)
	})
	t.Run("Respects set endianness", func(t *testing.T) {
		buf := &bytes.Buffer{}
		rd := &bin.Reader{R: buf, Endianness: binary.BigEndian}

		var u uint16
		buf.Write([]byte{0x11, 0x22})
		rd.Uint16(&u)
		assert.Equal(t, uint16(0x1122), u)
	})

	t.Run("Ints", func(t *testing.T) {
		run := func(t *testing.T, name string, tf func(t *testing.T, buf *bytes.Buffer, rd *bin.Reader)) {
			t.Run(name, func(t *testing.T) {
				buf := &bytes.Buffer{}
				rd := &bin.Reader{R: buf, Endianness: binary.BigEndian}
				tf(t, buf, rd)
			})
		}

		run(t, "Int8", func(t *testing.T, buf *bytes.Buffer, rd *bin.Reader) {
			var n int8
			buf.Write([]byte{0x69})
			rd.Int8(&n)
			assert.EqualValues(t, 0x69, n)
		})
		run(t, "Int16", func(t *testing.T, buf *bytes.Buffer, rd *bin.Reader) {
			var n int16
			buf.Write([]byte{0x69, 0x11})
			rd.Int16(&n)
			assert.EqualValues(t, 0x6911, n)
		})
		run(t, "Int32", func(t *testing.T, buf *bytes.Buffer, rd *bin.Reader) {
			var n int32
			buf.Write([]byte{0x69, 0x11, 0x22, 0x33})
			rd.Int32(&n)
			assert.EqualValues(t, 0x69112233, n)
		})
		run(t, "Int64", func(t *testing.T, buf *bytes.Buffer, rd *bin.Reader) {
			var n int64
			buf.Write([]byte{0x69, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77})
			rd.Int64(&n)
			assert.EqualValues(t, 0x6911223344556677, n)
		})
		// TODO.. better varint testing
		run(t, "Uvarint", func(t *testing.T, buf *bytes.Buffer, rd *bin.Reader) {
			buf.Write([]byte{0x11})
			var u uint64
			rd.Uvarint(&u)
			assert.EqualValues(t, 0x11, u)
		})
		run(t, "Varint", func(t *testing.T, buf *bytes.Buffer, rd *bin.Reader) {
			buf.Write([]byte{50})
			var u int64
			rd.Varint(&u)
			assert.EqualValues(t, 25, u)
		})
	})
}

func bigByteReader() *bytes.Reader {
	buf := make([]byte, 1024*1024)
	rand.Read(buf)
	return bytes.NewReader(buf)
}

func BenchmarkReader(b *testing.B) {
	bytRdr := bigByteReader()
	rdr := &bin.Reader{R: bytRdr}

	b.Run("Uint8", func(b *testing.B) {
		var u uint8
		b.SetBytes(1)
		for i := 0; i < b.N; i++ {
			rdr.Uint8(&u)
		}
		b.StopTimer()
		bytRdr.Seek(0, io.SeekStart)
	})
	b.Run("Uint64", func(b *testing.B) {
		var u uint64
		b.SetBytes(8)
		for i := 0; i < b.N; i++ {
			rdr.Uint64(&u)
		}
		b.StopTimer()
		bytRdr.Seek(0, io.SeekStart)
	})
	b.Run("Uvarint", func(b *testing.B) {
		var u uint64
		for i := 0; i < b.N; i++ {
			rdr.Uvarint(&u)
		}
		b.StopTimer()
		bytRdr.Seek(0, io.SeekStart)
	})
}

func BenchmarkEncodingBinaryReader(b *testing.B) {
	bytRdr := bigByteReader()

	b.Run("Uint8", func(b *testing.B) {
		var u uint8
		b.SetBytes(1)
		for i := 0; i < b.N; i++ {
			binary.Read(bytRdr, binary.LittleEndian, &u)
		}
		b.StopTimer()
		bytRdr.Seek(0, io.SeekStart)
	})
	b.Run("Uint64", func(b *testing.B) {
		var u uint64
		b.SetBytes(8)
		for i := 0; i < b.N; i++ {
			binary.Read(bytRdr, binary.LittleEndian, &u)
		}
		b.StopTimer()
		bytRdr.Seek(0, io.SeekStart)
	})
	b.Run("Uvarint", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			binary.ReadUvarint(bytRdr)
		}
		b.StopTimer()
		bytRdr.Seek(0, io.SeekStart)
	})
}
