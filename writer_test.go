package bin_test

import (
	"bytes"
	"encoding/binary"
	"testing"

	"github.com/ammario/bin"
	"github.com/stretchr/testify/assert"
)

func TestWriter(t *testing.T) {
	t.Run("Default little endian", func(t *testing.T) {
		buf := &bytes.Buffer{}
		wr := &bin.Writer{W: buf}

		wr.WriteInt16(0x1122)
		assert.Equal(t, []byte{0x22, 0x11}, buf.Bytes())
	})
	t.Run("Respects set endianness", func(t *testing.T) {
		buf := &bytes.Buffer{}
		wr := &bin.Writer{W: buf, Endianness: binary.BigEndian}

		wr.WriteInt16(0x1122)
		assert.Equal(t, []byte{0x11, 0x22}, buf.Bytes())
	})
	t.Run("Keeps proper count", func(t *testing.T) {
		buf := &bytes.Buffer{}
		wr := &bin.Writer{W: buf}

		wr.Write([]byte("123456"))
		assert.Equal(t, 6, wr.N())
	})
	t.Run("WriteInts", func(t *testing.T) {
		run := func(t *testing.T, name string, tf func(t *testing.T, buf *bytes.Buffer, wr *bin.Writer)) {
			t.Run(name, func(t *testing.T) {
				buf := &bytes.Buffer{}
				wr := &bin.Writer{W: buf, Endianness: binary.BigEndian}
				tf(t, buf, wr)
			})
		}

		run(t, "WriteInt8", func(t *testing.T, buf *bytes.Buffer, wr *bin.Writer) {
			wr.WriteInt8(0x69)
			assert.Equal(t, 1, buf.Len())
			assert.Equal(t, []byte{0x69}, buf.Bytes())
		})
		run(t, "WriteInt16", func(t *testing.T, buf *bytes.Buffer, wr *bin.Writer) {
			wr.WriteInt16(0x6060)
			assert.Equal(t, 2, buf.Len())
			assert.Equal(t, []byte{0x60, 0x60}, buf.Bytes())
		})
		run(t, "WriteInt32", func(t *testing.T, buf *bytes.Buffer, wr *bin.Writer) {
			wr.WriteInt32(0x11AABBCC)
			assert.Equal(t, 4, buf.Len())
			assert.Equal(t, []byte{0x11, 0xAA, 0xBB, 0xCC}, buf.Bytes())
		})
		run(t, "WriteInt64", func(t *testing.T, buf *bytes.Buffer, wr *bin.Writer) {
			wr.WriteInt64(0x0011AABBCCDDEEFF)
			assert.Equal(t, 8, buf.Len())
			assert.Equal(t, []byte{0x00, 0x11, 0xAA, 0xBB, 0xCC, 0xDD, 0xEE, 0xFF}, buf.Bytes())
		})
	})
}
