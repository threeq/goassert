package goassert

import (
	"errors"
	"testing"
)

// helper util
func failed(fn func(so *assertProxy)) bool {
	mockT := new(testing.T)
	so := Assertion(mockT)
	fn(so)
	return mockT.Failed()
}

type testStructDemo struct {
	f1 int
}

type testInterface interface {
	hh()
}

type testInterfaceImpl struct {
}

func (tif *testInterfaceImpl) hh() {

}

////////////////////////////

func TestThat(t *testing.T) {

	mockT := new(testing.T)

	a := "你好"
	That(mockT, a).As("xxx").
		Equal("你好").
		StartsWith("").
		EndsWith("").
		Len(6).
		Contains("").
		NotContain("h")
	if mockT.Failed() {
		t.Error("That function error")
	}
}

func TestAssertion(t *testing.T) {

	if failed(func(so *assertProxy) {
		so.That("hello world").As("xxx").
			Equal("hello world").
			StartsWith("h").
			EndsWith("d").
			Len(11).
			Contains("").
			NotContain("ss")
	}) {
		t.Error("That function error")
	}
}

func TestFluentAssertion_That(t *testing.T) {

	if !failed(func(so *assertProxy) {
		b1 := so.That("a").That("B")
		b1.Is(Empty)
	}) {
		t.Error("FluentAssertion.That string error")
	}
}

func TestFluentAssertion_Equal(t *testing.T) {
	mockT := new(testing.T)
	so := Assertion(mockT)

	so.That("1").Equal("1")
	if mockT.Failed() {
		t.Errorf("FluentAssertion.Equal is error")
	}

	mockT = new(testing.T)
	so = Assertion(mockT)
	so.That([]int{1, 2, 3}).Equal([]int{1, 2, 3})
	if mockT.Failed() {
		t.Errorf("FluentAssertion.Equal is error")
	}

	mockT = new(testing.T)
	so = Assertion(mockT)
	so.That([]int{1, 2, 3}).Equal([3]int{1, 2, 3})
	if !mockT.Failed() {
		t.Errorf("FluentAssertion.Equal is error")
	}

	mockT = new(testing.T)
	so = Assertion(mockT)
	so.That("1").Equal(1)
	if !mockT.Failed() {
		t.Errorf("FluentAssertion.Equal is error")
	}

	mockT = new(testing.T)
	so = Assertion(mockT)
	so.That(func() {}).Equal(func() {})
	if !mockT.Failed() {
		t.Errorf("FluentAssertion.Equal is error")
	}
}

func TestFluentAssertion_EqualIgnoringCase(t *testing.T) {
	var mockT *testing.T
	var so *assertProxy
	mockT = new(testing.T)
	so = Assertion(mockT)
	so.That("abc").EqualIgnoringCase("abc")
	if mockT.Failed() {
		t.Errorf("FluentAssertion.EqualIgnoringCase is error")
	}

	mockT = new(testing.T)
	so = Assertion(mockT)
	so.That("abc").EqualIgnoringCase("AbC")
	if mockT.Failed() {
		t.Errorf("FluentAssertion.EqualIgnoringCase is error")
	}
}

func TestFluentAssertion_NotEqual(t *testing.T) {
	var mockT *testing.T
	var so *assertProxy
	mockT = new(testing.T)
	so = Assertion(mockT)
	so.That("1").NotEqual("1")
	if !mockT.Failed() {
		t.Errorf("FluentAssertion.Equal is error")
	}

	mockT = new(testing.T)
	so = Assertion(mockT)
	so.That([]int{1, 2, 3}).NotEqual([]int{1, 2, 3})
	if !mockT.Failed() {
		t.Errorf("FluentAssertion.Equal is error")
	}

	mockT = new(testing.T)
	so = Assertion(mockT)
	so.That([]int{1, 2, 3}).NotEqual([3]int{1, 2, 3})
	if mockT.Failed() {
		t.Errorf("FluentAssertion.Equal is error")
	}

	mockT = new(testing.T)
	so = Assertion(mockT)
	so.That("1").NotEqual(1)
	if mockT.Failed() {
		t.Errorf("FluentAssertion.Equal is error")
	}

	mockT = new(testing.T)
	so = Assertion(mockT)
	so.That(func() {}).NotEqual(func() {})
	if !mockT.Failed() {
		t.Errorf("FluentAssertion.Equal is error")
	}
}

func TestFluentAssertion_StartsWith(t *testing.T) {
	mockT := new(testing.T)
	so := Assertion(mockT)
	so.That("你好").As("字符串").
		StartsWith("你")

	if mockT.Failed() {
		t.Error("FluentAssertion.StartsWith string error")
	}

	mockT = new(testing.T)
	so = Assertion(mockT)
	so.That("你好").As("字符串").
		StartsWith("好")

	if !mockT.Failed() {
		t.Error("FluentAssertion.StartsWith string error")
	}

	mockT = new(testing.T)
	so = Assertion(mockT)
	so.That([]string{"nihao", "你好", "你", "好"}).
		As("切片").
		StartsWith([]string{"nihao"}).
		StartsWith([]string{"nihao", "你好"})

	if mockT.Failed() {
		t.Error("FluentAssertion.StartsWith array/slice error")
	}

	mockT = new(testing.T)
	so = Assertion(mockT)
	so.That([]string{"nihao", "你好", "你", "好"}).
		As("切片").
		StartsWith([]string{"dfdfdf"}).
		StartsWith([]string{"nihao", "你好"})

	if !mockT.Failed() {
		t.Error("FluentAssertion.StartsWith array/slice error")
	}

	///////////////
	mockT = new(testing.T)
	so = Assertion(mockT)
	so.That([]string{"nihao", "你好", "你", "好"}).
		As("切片").
		StartsWith("sdfsdf")

	if !mockT.Failed() {
		t.Error("FluentAssertion.EndsWith array/slice error")
	}

	mockT = new(testing.T)
	so = Assertion(mockT)
	so.That([]string{"nihao", "你好", "你", "好"}).
		As("切片").
		StartsWith([]string{"nihao", "你好", "你", "好", "123213", "12321"})

	if !mockT.Failed() {
		t.Error("FluentAssertion.EndsWith array/slice error")
	}

	mockT = new(testing.T)
	so = Assertion(mockT)
	so.That(12323).
		StartsWith(123)

	if !mockT.Failed() {
		t.Error("FluentAssertion.EndsWith array/slice error")
	}
}

func TestFluentAssertion_EndsWith(t *testing.T) {
	mockT := new(testing.T)
	so := Assertion(mockT)
	so.That("你好").As("字符串").
		EndsWith("好")

	if mockT.Failed() {
		t.Error("FluentAssertion.EndsWith string error")
	}

	mockT = new(testing.T)
	so = Assertion(mockT)
	so.That("你好").As("字符串").
		EndsWith("你")

	if !mockT.Failed() {
		t.Error("FluentAssertion.EndsWith string error")
	}

	mockT = new(testing.T)
	so = Assertion(mockT)
	so.That([]string{"nihao", "你好", "你", "好"}).
		As("切片").
		EndsWith([]string{"好"}).
		EndsWith([]string{"你", "好"})

	if mockT.Failed() {
		t.Error("FluentAssertion.EndsWith array/slice error")
	}

	mockT = new(testing.T)
	so = Assertion(mockT)
	so.That([]string{"nihao", "你好", "你", "好"}).
		As("切片").
		EndsWith([]string{"dfdfdf"}).
		EndsWith([]string{"nihao", "你好"})

	if !mockT.Failed() {
		t.Error("FluentAssertion.EndsWith array/slice error")
	}

	///////////////
	mockT = new(testing.T)
	so = Assertion(mockT)
	so.That([]string{"nihao", "你好", "你", "好"}).
		As("切片").
		EndsWith("sdfsdf")

	if !mockT.Failed() {
		t.Error("FluentAssertion.EndsWith array/slice error")
	}

	mockT = new(testing.T)
	so = Assertion(mockT)
	so.That([]string{"nihao", "你好", "你", "好"}).
		As("切片").
		EndsWith([]string{"nihao", "你好", "你", "好", "123213", "12321"})

	if !mockT.Failed() {
		t.Error("FluentAssertion.EndsWith array/slice error")
	}

	mockT = new(testing.T)
	so = Assertion(mockT)
	so.That(12323).
		EndsWith(123)

	if !mockT.Failed() {
		t.Error("FluentAssertion.EndsWith array/slice error")
	}

}

func TestFluentAssertion_Len(t *testing.T) {
	var mockT *testing.T
	var so *assertProxy

	mockT = new(testing.T)
	so = Assertion(mockT)
	so.That("1").Len(1)
	if mockT.Failed() {
		t.Errorf("FluentAssertion.Len is error")
	}

	mockT = new(testing.T)
	so = Assertion(mockT)
	so.That("").Len(0)
	if mockT.Failed() {
		t.Errorf("FluentAssertion.Len is error")
	}

	mockT = new(testing.T)
	so = Assertion(mockT)
	so.That([]int{}).Len(0)
	if mockT.Failed() {
		t.Errorf("FluentAssertion.Len is error")
	}

	mockT = new(testing.T)
	so = Assertion(mockT)
	so.That([]int{1, 2, 3}).Len(3)
	if mockT.Failed() {
		t.Errorf("FluentAssertion.Len is error")
	}

	mockT = new(testing.T)
	so = Assertion(mockT)
	so.That([]int{1, 2, 3}).Len(2)
	if !mockT.Failed() {
		t.Errorf("FluentAssertion.Len is error")
	}

	mockT = new(testing.T)
	so = Assertion(mockT)
	so.That(1).Len(2)
	if !mockT.Failed() {
		t.Errorf("FluentAssertion.Len is error")
	}
}

func TestFluentAssertion_Contains(t *testing.T) {

	if failed(func(so *assertProxy) {
		so.That("").Contains("")
	}) {
		t.Errorf("FluentAssertion.Contains is error")
	}

	if failed(func(so *assertProxy) {
		so.That("hello world").Contains("world")
	}) {
		t.Errorf("FluentAssertion.Contains is error")
	}

	if failed(func(so *assertProxy) {
		so.That([]int{1, 2, 3}).Contains(1)
	}) {
		t.Errorf("FluentAssertion.Contains is error")
	}

	if failed(func(so *assertProxy) {
		so.That([]int{1, 2, 3}).Contains(3)
	}) {
		t.Errorf("FluentAssertion.Contains is error")
	}

	if !failed(func(so *assertProxy) {
		so.That([]int{1, 2, 3}).Contains([]int{1})
	}) {
		t.Errorf("FluentAssertion.Contains is error")
	}

	if !failed(func(so *assertProxy) {
		so.That([]int{1, 2, 3}).Contains([]int{2})
	}) {
		t.Errorf("FluentAssertion.Contains is error")
	}

	if !failed(func(so *assertProxy) {
		so.That([]int{1, 2, 3}).Contains([]int{1, 2})
	}) {
		t.Errorf("FluentAssertion.Contains is error")
	}

	if !failed(func(so *assertProxy) {
		so.That([]int{1, 2, 3}).Contains([]int{1, 3})
	}) {
		t.Errorf("FluentAssertion.Contains is error")
	}

	if !failed(func(so *assertProxy) {
		so.That([]int{1, 2, 3}).Contains([]int{})
	}) {
		t.Errorf("FluentAssertion.Contains is error")
	}

	if !failed(func(so *assertProxy) {
		so.That([]int{1, 2, 3}).Contains([]int{4})
	}) {
		t.Errorf("FluentAssertion.Contains is error")
	}

	if !failed(func(so *assertProxy) {
		so.That([]int{1, 2, 3}).Contains([]int{1, 4})
	}) {
		t.Errorf("FluentAssertion.Contains is error")
	}
}

func TestFluentAssertion_NotContain(t *testing.T) {
	if !failed(func(so *assertProxy) {
		so.That("").NotContain("")
	}) {
		t.Errorf("FluentAssertion.NotContain is error")
	}

	if !failed(func(so *assertProxy) {
		so.That("hello world").NotContain("world")
	}) {
		t.Errorf("FluentAssertion.NotContain is error")
	}

	if !failed(func(so *assertProxy) {
		so.That([]int{1, 2, 3}).NotContain(1)
	}) {
		t.Errorf("FluentAssertion.NotContain is error")
	}

	if !failed(func(so *assertProxy) {
		so.That([]int{1, 2, 3}).NotContain(3)
	}) {
		t.Errorf("FluentAssertion.NotContain is error")
	}

	if failed(func(so *assertProxy) {
		so.That([]int{1, 2, 3}).NotContain([]int{1})
	}) {
		t.Errorf("FluentAssertion.NotContain is error")
	}

	if failed(func(so *assertProxy) {
		so.That([]int{1, 2, 3}).NotContain([]int{2})
	}) {
		t.Errorf("FluentAssertion.NotContain is error")
	}

	if failed(func(so *assertProxy) {
		so.That([]int{1, 2, 3}).NotContain([]int{1, 2})
	}) {
		t.Errorf("FluentAssertion.NotContain is error")
	}

	if failed(func(so *assertProxy) {
		so.That([]int{1, 2, 3}).NotContain([]int{1, 3})
	}) {
		t.Errorf("FluentAssertion.NotContain is error")
	}

	if failed(func(so *assertProxy) {
		so.That([]int{1, 2, 3}).NotContain([]int{})
	}) {
		t.Errorf("FluentAssertion.NotContain is error")
	}

	if failed(func(so *assertProxy) {
		so.That([]int{1, 2, 3}).NotContain([]int{4})
	}) {
		t.Errorf("FluentAssertion.NotContain is error")
	}

	if failed(func(so *assertProxy) {
		so.That([]int{1, 2, 3}).NotContain([]int{1, 4})
	}) {
		t.Errorf("FluentAssertion.NotContain is error")
	}
}

func TestFluentAssertion_In(t *testing.T) {
	if failed(func(so *assertProxy) {
		so.That("").In("")
	}) {
		t.Errorf("FluentAssertion.In is error")
	}

	if failed(func(so *assertProxy) {
		so.That("").In("123")
	}) {
		t.Errorf("FluentAssertion.In is error")
	}

	if failed(func(so *assertProxy) {
		so.That("1").In("123")
	}) {
		t.Errorf("FluentAssertion.In is error")
	}

	if failed(func(so *assertProxy) {
		so.That(1).In([]int{1, 2, 3})
	}) {
		t.Errorf("FluentAssertion.In is error")
	}

	if !failed(func(so *assertProxy) {
		so.That("4").In("123")
	}) {
		t.Errorf("FluentAssertion.In is error")
	}

	if !failed(func(so *assertProxy) {
		so.That(4).In([]int{1, 2, 3})
	}) {
		t.Errorf("FluentAssertion.In is error")
	}
}

func TestFluentAssertion_NotIn(t *testing.T) {
	if !failed(func(so *assertProxy) {
		so.That("").NotIn("")
	}) {
		t.Errorf("FluentAssertion.NotIn is error")
	}

	if !failed(func(so *assertProxy) {
		so.That("").NotIn("123")
	}) {
		t.Errorf("FluentAssertion.NotIn is error")
	}

	if !failed(func(so *assertProxy) {
		so.That("1").NotIn("123")
	}) {
		t.Errorf("FluentAssertion.NotIn is error")
	}

	if !failed(func(so *assertProxy) {
		so.That(1).NotIn([]int{1, 2, 3})
	}) {
		t.Errorf("FluentAssertion.NotIn is error")
	}

	if failed(func(so *assertProxy) {
		so.That("4").NotIn("123")
	}) {
		t.Errorf("FluentAssertion.NotIn is error")
	}

	if failed(func(so *assertProxy) {
		so.That(4).NotIn([]int{1, 2, 3})
	}) {
		t.Errorf("FluentAssertion.NotIn is error")
	}
}

func TestFluentAssertion_Is(t *testing.T) {

	if failed(func(so *assertProxy) {
		so.That("你好").
			Is(Not(Empty))
	}) {
		t.Error("FluentAssertion.Is string error")
	}

	if !failed(func(so *assertProxy) {
		so.That("你好").
			Is(Empty)
	}) {
		t.Error("FluentAssertion.Is string error")
	}
}

func TestFluentAssertion_Not(t *testing.T) {
	if !failed(func(so *assertProxy) {
		so.That("你好").
			Not(Not(Empty))
	}) {
		t.Error("FluentAssertion.Not string error")
	}

	if failed(func(so *assertProxy) {
		so.That("你好").
			Not(Empty)
	}) {
		t.Error("FluentAssertion.Not string error")
	}
}

func TestFluentAssertion_IsType(t *testing.T) {
	if failed(func(so *assertProxy) {
		so.That("你好").IsType(string(""))
	}) {
		t.Error("FluentAssertion.IsType string error")
	}

	if failed(func(so *assertProxy) {
		so.That(1).IsType(int(0))
	}) {
		t.Error("FluentAssertion.IsType string error")
	}

	if !failed(func(so *assertProxy) {
		so.That(1).IsType(int64(0))
	}) {
		t.Error("FluentAssertion.IsType string error")
	}

	if !failed(func(so *assertProxy) {
		so.That(1).IsType(int32(0))
	}) {
		t.Error("FluentAssertion.IsType string error")
	}

	if failed(func(so *assertProxy) {
		so.That(testStructDemo{1231}).IsType(testStructDemo{})
	}) {
		t.Error("FluentAssertion.IsType string error")
	}

	if failed(func(so *assertProxy) {
		so.That(&testStructDemo{1231}).IsType(&testStructDemo{})
	}) {
		t.Error("FluentAssertion.IsType string error")
	}

	if !failed(func(so *assertProxy) {
		so.That(&testStructDemo{1231}).IsType(testStructDemo{})
	}) {
		t.Error("FluentAssertion.IsType string error")
	}
}

func TestFluentAssertion_AllOf(t *testing.T) {
	var cond1 Condition = func(actual interface{}) (b bool, s string) {
		return actual.(string) == "123", "custom Condition"
	}

	if failed(func(so *assertProxy) {
		so.That("123").AllOf(Not(Empty), Len(3), Eq("123"), cond1)
	}) {
		t.Error("FluentAssertion.AllOf string error")
	}

	if !failed(func(so *assertProxy) {
		so.That("123").AllOf(Not(Empty), Len(3), Eq("13"), cond1)
	}) {
		t.Error("FluentAssertion.AllOf string error")
	}
}

func TestFluentAssertion_AnyOf(t *testing.T) {
	var cond1 Condition = func(actual interface{}) (b bool, s string) {
		b = actual.(string) == "123"
		return

	}
	if failed(func(so *assertProxy) {
		so.That("123").AnyOf(Not(Empty), Len(3), Eq("123"), cond1)
	}) {
		t.Error("FluentAssertion.AnyOf string error")
	}

	if failed(func(so *assertProxy) {
		so.That("123").AnyOf(Not(Empty), Len(3), Eq("13"), cond1)
	}) {
		t.Error("FluentAssertion.AnyOf string error")
	}

	if !failed(func(so *assertProxy) {
		so.That("123").AnyOf(Empty, Len(4), Eq("13"))
	}) {
		t.Error("FluentAssertion.AnyOf string error")
	}
}

func TestFluentAssertion_As(t *testing.T) {
	if failed(func(so *assertProxy) {
		so.That("123").As("ok")
	}) {
		t.Error("FluentAssertion.As string error")
	}
}

func TestFluentAssertion_DirExists(t *testing.T) {
	if failed(func(so *assertProxy) {
		so.That("/usr").DirExists()
	}) {
		t.Error("FluentAssertion.DirExists string error")
	}

	if !failed(func(so *assertProxy) {
		so.That("/xxx-xx-x-x-x-x").DirExists()
	}) {
		t.Error("FluentAssertion.DirExists string error")
	}
}

func TestFluentAssertion_FileExists(t *testing.T) {
	if failed(func(so *assertProxy) {
		so.That("/bin/ls").FileExists()
	}) {
		t.Error("FluentAssertion.FileExists string error")
	}

	if !failed(func(so *assertProxy) {
		so.That("/xxx-xx-x-x-x-x").FileExists()
	}) {
		t.Error("FluentAssertion.FileExists string error")
	}
}

func TestFluentAssertion_HasMessage(t *testing.T) {
	if failed(func(so *assertProxy) {
		so.That(errors.New("error")).HasMessage("error")
	}) {
		t.Error("FluentAssertion.HasMessage error")
	}

	if !failed(func(so *assertProxy) {
		so.That(errors.New("error error123")).HasMessage("errorerror")
	}) {
		t.Error("FluentAssertion.HasMessage error")
	}
}

func TestFluentAssertion_Implements(t *testing.T) {
	if failed(func(so *assertProxy) {
		so.That(&testInterfaceImpl{}).Implements((*testInterface)(nil))
	}) {
		t.Error("FluentAssertion.Implements error")
	}

	if !failed(func(so *assertProxy) {
		so.That(testInterfaceImpl{}).Implements((*testInterface)(nil))
	}) {
		t.Error("FluentAssertion.Implements error")
	}

	if !failed(func(so *assertProxy) {
		so.That(testStructDemo{}).Implements((*testInterface)(nil))
	}) {
		t.Error("FluentAssertion.Implements error")
	}

	if !failed(func(so *assertProxy) {
		so.That(&testStructDemo{}).Implements((*testInterface)(nil))
	}) {
		t.Error("FluentAssertion.Implements error")
	}
}

func TestFluentAssertion_JSONEq(t *testing.T) {
	if failed(func(so *assertProxy) {
		so.That(`{"1":"1", "2":2}`).JSONEq(`{"2":2, "1":"1"}`)
	}) {
		t.Errorf("FluentAssertion.JSONEq error")
	}

	if failed(func(so *assertProxy) {
		so.That(`["foo", {"hello": "world", "nested": "hash"}]`).
			JSONEq(`["foo", {"nested": "hash", "hello": "world"}]`)
	}) {
		t.Errorf("FluentAssertion.JSONEq error")
	}

	if !failed(func(so *assertProxy) {
		so.That(`{"1":"1", "2":2}`).JSONEq(`{"2":2, "1":"11"}`)
	}) {
		t.Errorf("FluentAssertion.JSONEq error")
	}
}

func TestFluentAssertion_NotRegexp(t *testing.T) {
	if !failed(func(so *assertProxy) {
		so.That("hello world").NotRegexp("^he.*d$")
	}) {
		t.Errorf("FluentAssertion.NotRegexp error")
	}

	if failed(func(so *assertProxy) {
		so.That("hello world").NotRegexp("^hedd.*d$")
	}) {
		t.Errorf("FluentAssertion.NotRegexp error")
	}
}

func TestFluentAssertion_Regexp(t *testing.T) {
	if failed(func(so *assertProxy) {
		so.That("hello world").Regexp("^he.*d$")
	}) {
		t.Errorf("FluentAssertion.Regexp error")
	}

	if !failed(func(so *assertProxy) {
		so.That("hello world").Regexp("^hedd.*d$")
	}) {
		t.Errorf("FluentAssertion.Regexp error")
	}
}
func TestFluentAssertion_Zero(t *testing.T) {
	if failed(func(so *assertProxy) {
		so.That(0).Zero()
	}) {
		t.Errorf("FluentAssertion.Zero error")
	}

	if failed(func(so *assertProxy) {
		so.That("").Zero()
	}) {
		t.Errorf("FluentAssertion.Zero error")
	}

	if failed(func(so *assertProxy) {
		so.That(nil).Zero()
	}) {
		t.Errorf("FluentAssertion.Zero error")
	}

	if !failed(func(so *assertProxy) {
		so.That(1).Zero()
	}) {
		t.Errorf("FluentAssertion.Zero error")
	}
}

func TestFluentAssertion_NotZero(t *testing.T) {
	if !failed(func(so *assertProxy) {
		so.That(0).NotZero()
	}) {
		t.Errorf("FluentAssertion.NotZero error")
	}

	if !failed(func(so *assertProxy) {
		so.That("").NotZero()
	}) {
		t.Errorf("FluentAssertion.NotZero error")
	}

	if !failed(func(so *assertProxy) {
		so.That(nil).NotZero()
	}) {
		t.Errorf("FluentAssertion.NotZero error")
	}

	if failed(func(so *assertProxy) {
		so.That(1).NotZero()
	}) {
		t.Errorf("FluentAssertion.NotZero error")
	}
}

func TestFluentAssertion_Panics(t *testing.T) {
	if !failed(func(so *assertProxy) {
		so.That(nil).Panics(func() {})
	}) {
		t.Errorf("FluentAssertion.Panics error")
	}

	if failed(func(so *assertProxy) {
		so.That(nil).Panics(func() {
			panic("has panic")
		})
	}) {
		t.Errorf("FluentAssertion.Panics error")
	}
}

func TestAssertProxy_Panics(t *testing.T) {
	mockT := new(testing.T)
	assert := Assertion(mockT)
	assert.Panics(func() {})
	if !mockT.Failed() {
		t.Errorf("assertProxy.Panics error")
	}

	mockT = new(testing.T)
	assert = Assertion(mockT)
	assert.Panics(func() {
		panic("has panic")
	})
	if mockT.Failed() {
		t.Errorf("assertProxy.Panics error")
	}
}

func TestAssertProxy_That(t *testing.T) {
	so := Assertion(t).That("")
	if so.t != t {
		t.Errorf("assertProxy.That error")
	}
}
