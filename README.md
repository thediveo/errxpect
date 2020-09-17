# errxpect

[![PkgGoDev](https://pkg.go.dev/badge/github.com/thediveo/errxpect)](https://pkg.go.dev/github.com/thediveo/errxpect)
[![GitHub](https://img.shields.io/github/license/thediveo/errxpect)](https://img.shields.io/github/license/thediveo/errxpect)
![build and test](https://github.com/thediveo/errxpect/workflows/build%20and%20test/badge.svg?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/thediveo/errxpect)](https://goreportcard.com/report/github.com/thediveo/errxpect)

## Tired?

Worn down by [Gomega's](https://github.com/onsi/gomega) noisy error testing
boilerplate for function returning multiple return values? Each time, given a
function returning multiple values and an error...

```go
func Foo(int) (string, bool, error) {
    return "", false, errors.New("DOH!")
}
```

...Gomega forces you to do break function call and test into separate steps,
requiring intermediate result variables (with most of them `_`s anyway):

```go
_, _, err := Foo(42)
Expect(err).To(HaveOccured())
```

## Errxpect!

Just `import . "github.com/thediveo/errxpect"` and then use `Errxpect(...)` in
place of `Expect(...)`. And enjoy more fluent error test assertions.

```go
Errxpect(Foo(42)).To(HaveOccured())
```

As Golang doesn't unpack multiple return values automatically when there is
another parameter present in a function call, error expectations with stack
offsets need to the phrased as follows using `.WithOffset(offset)`, keeping them
elegant:

```go
Errxpect(Foo(42)).WithOffset(1).To(HaveOccured())
```

## Copyright and License

`errxpect` is Copyright 2020 Harald Albrecht, and licensed under the Apache
License, Version 2.0.
