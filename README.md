# bin

[![GoDoc](https://godoc.org/github.com/ammario/bin?status.svg)](https://godoc.org/github.com/ammario/bin)


A simple, ergonomic wrapper around `encoding/binary`.

It provides a `Writer` and `Reader` type.
Once an error is encountered, further Write or Read calls will be no-op and the error will be located in `Err()`.

In essence, this package makes code that looked like this:

```go
var nn int

n, err := thing.Write(a)
nn += n
if err != nil {
    return nn, err
}

n, err = thing.Write(b)
nn += n
if err != nil {
    return nn, err
}
n, err = thing.Write(c)
nn += n

return nn, err
```

look like this:

```go
wr := bin.Writer{W: thing}
wr.Write(a)
wr.Write(b)
wr.Write(c)
return wr.N(), wr.Err()
```

## Performance

Due to `bin`'s type-specific int read and write methods, it performs better than `encoding/binary`
on non-varint methods:

```
goos: linux
goarch: amd64
pkg: github.com/ammario/bin
BenchmarkReader/Uint8-8           	500000000	         3.74 ns/op	 267.25 MB/s
BenchmarkReader/Uint64-8          	500000000	         3.83 ns/op	2088.52 MB/s
BenchmarkReader/Uvarint-8         	200000000	         7.18 ns/op
BenchmarkEncodingBinaryReader/Uint8-8         	50000000	        28.4 ns/op	  35.24 MB/s
BenchmarkEncodingBinaryReader/Uint64-8        	50000000	        28.3 ns/op	 282.68 MB/s
BenchmarkEncodingBinaryReader/Uvarint-8       	300000000	         5.50 ns/op
BenchmarkWriter/Uint8-8                       	100000000	        12.3 ns/op	  81.21 MB/s
BenchmarkWriter/Uint64-8                      	100000000	        20.0 ns/op	 400.12 MB/s
BenchmarkWriter/Uvarint-8                     	100000000	        15.1 ns/op
BenchmarkEncodingBinaryWriter/Uint8-8         	50000000	        25.1 ns/op	  39.85 MB/s
BenchmarkEncodingBinaryWriter/Uint64-8        	30000000	        43.3 ns/op	 184.81 MB/s
BenchmarkEncodingBinaryWriter/Uvarint-8       	100000000	        10.3 ns/op
```
