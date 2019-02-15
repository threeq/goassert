package example

import (
	"github.com/threeq/goassert"
	"testing"
)

func TestExampleCondition(t *testing.T) {
	a := "你好"
	goassert.That(t, a).As("xxx").
		Equal("你好").
		StartsWith("").
		EndsWith("").
		Len(6).
		Contains("").
		NotContain("h")
}
