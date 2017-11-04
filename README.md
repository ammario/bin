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