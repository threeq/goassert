# goassert

[![Build Status](https://travis-ci.org/threeq/goassert.svg?branch=master)](https://travis-ci.org/threeq/goassert)
[![codecov](https://codecov.io/gh/threeq/goassert/branch/master/graph/badge.svg)](https://codecov.io/gh/threeq/goassert)

一个流畅的断言库。实现参考 [testify](https://github.com/stretchr/testify) 和 [assertj](http://joel-costigliola.github.io/assertj)。

# Install

```shell
go get github.com/threeq/goassert
```

# Use

1. import `goassert` package
 
```go
import "github.com/threeq/goassert"
```

2. use `goassert`

use `That` function

```go
func TestExample(t *testing.T) {
	a := "你好"
	//
	goassert.That(t, a).As("xxx").
		Equal("你好").
		StartsWith("").
		EndsWith("").
		Len(6).
		Contains("").
		NotContain("h")
}
```

or use `Assertion` Object

```go
func TestExample(t *testing.T) {
	so := goassert.Assertion(t)
	a := "hello world"
	// do something
	so.That(a).As("xxx").
		Equal("hello world").
		StartsWith("").
		EndsWith("").
		Len(6).
		Contains("").
		NotContain("h")
}
```

## Use Condition

Assertion contain common assertions. 
But for the complex assertion (such as logical judgment, etc.) is not satisfied, can be implemented using `Condition`. 
`goassert` already contains a common `Condition` implementation.

```go
func TestExample(t *testing.T) {
	so := goassert.Assertion(t)
	a := "hello world"
    // do something
    so.That(a).
        Is(Not(Empty))
}
```

Use custom condition:

```go
cond1 = func(actual interface{}) (b bool, s string) {
    return actual.(string) == "123", "custom Condition"
}

func TestExample(t *testing.T) {
	so := goassert.Assertion(t)
	a := "hello world"
    // do something
    so.That(a).Is(cond1)
}
```

# Example

