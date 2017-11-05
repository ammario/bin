package bin_test

import (
	"bytes"
	"encoding/binary"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ammario/bin"
	"github.com/stretchr/testify/assert"
)

func TestWriter(t *testing.T) {
	t.Run("Default little endian", func(t *testing.T) {
		buf := &bytes.Buffer{}
		wr := &bin.Writer{W: buf}

		wr.Int16(0x1122)
		assert.Equal(t, []byte{0x22, 0x11}, buf.Bytes())
	})
	t.Run("Respects set endianness", func(t *testing.T) {
		buf := &bytes.Buffer{}
		wr := &bin.Writer{W: buf, Endianness: binary.BigEndian}

		wr.Int16(0x1122)
		assert.Equal(t, []byte{0x11, 0x22}, buf.Bytes())
	})
	t.Run("Keeps first non-nil error and proper count", func(t *testing.T) {
		fi, err := ioutil.TempFile("", "bintest")
		require.NoError(t, err)
		defer os.Remove(fi.Name())

		wr := &bin.Writer{W: fi}
		wr.Write([]byte("123456"))
		assert.Equal(t, 6, wr.N())
		assert.NoError(t, wr.Err())
		fi.Close()

		wr.Write([]byte("this doesn't get written"))
		assert.Equal(t, 6, wr.N())
		assert.Error(t, wr.Err())

	})
	t.Run("Ints", func(t *testing.T) {
		run := func(t *testing.T, name string, tf func(t *testing.T, buf *bytes.Buffer, wr *bin.Writer)) {
			t.Run(name, func(t *testing.T) {
				buf := &bytes.Buffer{}
				wr := &bin.Writer{W: buf, Endianness: binary.BigEndian}
				tf(t, buf, wr)
			})
		}

		run(t, "Int8", func(t *testing.T, buf *bytes.Buffer, wr *bin.Writer) {
			wr.Int8(0x69)
			assert.Equal(t, 1, buf.Len())
			assert.Equal(t, []byte{0x69}, buf.Bytes())
		})
		run(t, "Int16", func(t *testing.T, buf *bytes.Buffer, wr *bin.Writer) {
			wr.Int16(0x6060)
			assert.Equal(t, 2, buf.Len())
			assert.Equal(t, []byte{0x60, 0x60}, buf.Bytes())
		})
		run(t, "Int32", func(t *testing.T, buf *bytes.Buffer, wr *bin.Writer) {
			wr.Int32(0x11AABBCC)
			assert.Equal(t, 4, buf.Len())
			assert.Equal(t, []byte{0x11, 0xAA, 0xBB, 0xCC}, buf.Bytes())
		})
		run(t, "Int64", func(t *testing.T, buf *bytes.Buffer, wr *bin.Writer) {
			wr.Int64(0x0011AABBCCDDEEFF)
			assert.Equal(t, 8, buf.Len())
			assert.Equal(t, []byte{0x00, 0x11, 0xAA, 0xBB, 0xCC, 0xDD, 0xEE, 0xFF}, buf.Bytes())
		})
		run(t, "Uvarint", func(t *testing.T, buf *bytes.Buffer, wr *bin.Writer) {
			wr.Uvarint(0x11)
			assert.Equal(t, 1, buf.Len())
			assert.Equal(t, []byte{0x11}, buf.Bytes())
		})
		run(t, "Varint", func(t *testing.T, buf *bytes.Buffer, wr *bin.Writer) {
			wr.Uvarint(0x0FAB)
			assert.Equal(t, 2, buf.Len())
			assert.Equal(t, []byte{0xab, 0x1f}, buf.Bytes())
		})
	})
}

func BenchmarkWriter(b *testing.B) {
	buf := bytes.NewBuffer(make([]byte, 1024*1024))
	wr := &bin.Writer{W: buf}

	b.Run("Uint8", func(b *testing.B) {
		b.SetBytes(1)
		for i := 0; i < b.N; i++ {
			wr.Uint8(uint8(i))
		}
		b.StopTimer()
		buf.Reset()
	})
	b.Run("Uint64", func(b *testing.B) {
		b.SetBytes(8)
		for i := 0; i < b.N; i++ {
			wr.Uint64(uint64(i))
		}
		b.StopTimer()
		buf.Reset()
	})
	b.Run("Uvarint", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			wr.Uvarint(uint64(i))
		}
		b.StopTimer()
		buf.Reset()
	})
}

func BenchmarkEncodingBinaryWriter(b *testing.B) {
	buf := bytes.NewBuffer(make([]byte, 1024*1024))

	b.Run("Uint8", func(b *testing.B) {
		b.SetBytes(1)
		for i := 0; i < b.N; i++ {
			binary.Write(buf, binary.LittleEndian, uint8(i))
		}
		b.StopTimer()
		buf.Reset()
	})
	b.Run("Uint64", func(b *testing.B) {
		b.SetBytes(8)
		for i := 0; i < b.N; i++ {
			binary.Write(buf, binary.LittleEndian, uint64(i))
		}
		b.StopTimer()
		buf.Reset()
	})
	b.Run("Uvarint", func(b *testing.B) {
		var scratch [8]byte
		for i := 0; i < b.N; i++ {
			n := binary.PutUvarint(scratch[:], uint64(i))
			buf.Write(scratch[:n])
		}
	})
}
