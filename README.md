# testa
Readable non-ambiguous testing library for Golang

# assert package
Provides assert functionality.

There are two kinds of asserts. Both mark the function under test as having failed, but one stops the execution of the test (fatal asserts) and the other one allows it to continue.

Non-fatal assert function construction

```go
func TestExampleFunc(t *testing.T) {
    // Constructs the non-fatal assert function
    assert := testa.New(t)
}
```

Fatal assert function construction

```go
func TestExampleFunc(t *testing.T) {
    // Constructs the fatal assert function
    require := testa.NewFatal(t)
}
```

Usage

```go
func TestExampleFunc(t *testing.T) {
    want := 5
    got := ExampleFunc()

    // assert got value equals want
    assert(got).Equals(want)
}
```

```go
func TestExampleFunc(t *testing.T) {
    got, err := ExampleFunc()

    // assert err value is nil
    require(err).IsNil()
    // assert got value is not nil
    require(got).IsNotNil()
}
```

# Licence
This project is licensed under the terms of the MIT license.

# Sources

* [Testify](https://github.com/stretchr/testify)