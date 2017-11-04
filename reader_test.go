package bin_test

import (
	"bytes"
	"encoding/binary"
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
		rd.ReadUint16(&u)

		assert.Equal(t, uint16(0x2211), u)
	})
	t.Run("Respects set endianness", func(t *testing.T) {
		buf := &bytes.Buffer{}
		rd := &bin.Reader{R: buf, Endianness: binary.BigEndian}

		var u uint16
		buf.Write([]byte{0x11, 0x22})
		rd.ReadUint16(&u)
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

		run(t, "ReadInt8", func(t *testing.T, buf *bytes.Buffer, rd *bin.Reader) {
			var n int8
			buf.Write([]byte{0x69})
			rd.ReadInt8(&n)
			assert.EqualValues(t, 0x69, n)
		})
		run(t, "ReadInt16", func(t *testing.T, buf *bytes.Buffer, rd *bin.Reader) {
			var n int16
			buf.Write([]byte{0x69, 0x11})
			rd.ReadInt16(&n)
			assert.EqualValues(t, 0x6911, n)
		})
		run(t, "ReadInt32", func(t *testing.T, buf *bytes.Buffer, rd *bin.Reader) {
			var n int32
			buf.Write([]byte{0x69, 0x11, 0x22, 0x33})
			rd.ReadInt32(&n)
			assert.EqualValues(t, 0x69112233, n)
		})
		run(t, "ReadInt64", func(t *testing.T, buf *bytes.Buffer, rd *bin.Reader) {
			var n int64
			buf.Write([]byte{0x69, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77})
			rd.ReadInt64(&n)
			assert.EqualValues(t, 0x6911223344556677, n)
		})
	})
}
