package goassert_test

import (
	"errors"
	. "github.com/threeq/goassert"
	"os"
	"regexp"
	"strings"
	"testing"
	"time"
)

func TestNot(t *testing.T) {
	res, msg := Empty("")
	if res != true || msg != "is Empty" {
		t.Error("Empty string is empty")
	}

	res, msg = Not(Empty)("")
	if res != false || msg != "Not is Empty" {
		t.Error("Not operational error")
	}

	res, msg = Not(Empty)("xxx")
	if res != true || msg != "Not is Empty" {
		t.Error("Not operational error")
	}
}

func TestAnd(t *testing.T) {
	res, msg := And(Empty, Nil)("")
	if res != false || msg != "is Nil" {
		t.Errorf("And operational error: %v,%s", res, msg)
	}

	res, msg = And(Empty, Nil)("xxxx")
	if res != false || msg != "is Empty" {
		t.Errorf("And operational error: %v,%s", res, msg)
	}

	res, msg = And(Nil, Empty)("xxxx")
	if res != false || msg != "is Nil" {
		t.Errorf("And operational error: %v,%s", res, msg)
	}

	res, msg = And(Not(Nil), Not(Empty))("xxxx")
	if res != true || msg != "" {
		t.Errorf("And operational error: %v,%s", res, msg)
	}

	res, msg = And(Not(Nil), Empty)("")
	if res != true || msg != "" {
		t.Errorf("And operational error: %v,%s", res, msg)
	}
}

func TestOr(t *testing.T) {
	res, msg := Or(Empty, Nil)("")
	if res != true || msg != "" {
		t.Errorf("And operational error: %v,%s", res, msg)
	}

	res, msg = Or(Empty, Nil)("xxxx")
	if res != false || msg != "(is Empty or is Nil)" {
		t.Errorf("And operational error: %v,%s", res, msg)
	}

	res, msg = Or(Nil, Empty)("xxxx")
	if res != false || msg != "(is Nil or is Empty)" {
		t.Errorf("And operational error: %v,%s", res, msg)
	}

	res, msg = Or(Not(Nil), Not(Empty))("xxxx")
	if res != true || msg != "" {
		t.Errorf("And operational error: %v,%s", res, msg)
	}

	res, msg = Or(Not(Nil), Empty)("")
	if res != true || msg != "" {
		t.Errorf("And operational error: %v,%s", res, msg)
	}

	res, msg = Or(Greater(7), Eq(7))(6)
	if res != false || msg != "(> 7 or == 7)" {
		t.Errorf("Or operational error: %v,%s", res, msg)
	}
	res, msg = Or(Greater(7), Eq(7))(7)
	if res != true || msg != "" {
		t.Errorf("Or operational error: %v,%s", res, msg)
	}
	res, msg = Or(Greater(7), Eq(7))(8)
	if res != true || msg != "" {
		t.Errorf("Or operational error: %v,%s", res, msg)
	}
}

func TestNil(t *testing.T) {

	res, msg := Nil(nil)
	if res != true {
		t.Error("nil test error")
	}
	if msg != "is Nil" {
		t.Error("nil msg error：" + msg)
	}

	res, msg = Nil("msg")
	if res {
		t.Error("not nil test error")
	}
	if msg != "is Nil" {
		t.Error("not nil msg error：" + msg)
	}
}

func TestEmpty(t *testing.T) {
	chWithValue := make(chan struct{}, 1)
	chWithValue <- struct{}{}
	var tiP *time.Time
	var tiNP time.Time
	var s *string
	var f *os.File
	sP := &s
	x := 1
	xP := &x

	type TString string
	type TStruct struct {
		x int
		s []int
	}

	// True
	if res, msg := Empty(""); res != true || msg != "is Empty" {
		t.Error("Empty string is empty")
	}
	if res, msg := Empty(nil); res != true || msg != "is Empty" {
		t.Error("Nil is empty")
	}
	if res, msg := Empty([]string{}); res != true || msg != "is Empty" {
		t.Error("Empty string array is empty")
	}
	if res, msg := Empty(0); res != true || msg != "is Empty" {
		t.Error("Zero int value is empty")
	}
	if res, msg := Empty(false); res != true || msg != "is Empty" {
		t.Error("False value is empty")
	}
	if res, msg := Empty(make(chan struct{})); res != true || msg != "is Empty" {
		t.Error("Channel without values is empty")
	}
	if res, msg := Empty(s); res != true || msg != "is Empty" {
		t.Error("Nil string pointer is empty")
	}
	if res, msg := Empty(f); res != true || msg != "is Empty" {
		t.Error("Nil os.File pointer is empty")
	}
	if res, msg := Empty(tiP); res != true || msg != "is Empty" {
		t.Error("Nil time.Time pointer is empty")
	}
	if res, msg := Empty(tiNP); res != true || msg != "is Empty" {
		t.Error("time.Time is empty")
	}
	if res, msg := Empty(TStruct{}); res != true || msg != "is Empty" {
		t.Error("struct with zero values is empty")
	}
	if res, msg := Empty(TString("")); res != true || msg != "is Empty" {
		t.Error("empty aliased string is empty")
	}
	if res, msg := Empty(sP); res != true || msg != "is Empty" {
		t.Error("ptr to nil value is empty")
	}

	// False
	if f, msg := Empty("something"); f != false || msg != "is Empty" {
		t.Error("Non Empty string is not empty")
	}
	if f, msg := Empty(errors.New("something")); f != false || msg != "is Empty" {
		t.Error("Non nil object is not empty")
	}
	if f, msg := Empty([]string{"something"}); f != false || msg != "is Empty" {
		t.Error("Non empty string array is not empty")
	}
	if f, msg := Empty(1); f != false || msg != "is Empty" {
		t.Error("Non-zero int value is not empty")
	}
	if f, msg := Empty(true); f != false || msg != "is Empty" {
		t.Error("True value is not empty")
	}
	if f, msg := Empty(chWithValue); f != false || msg != "is Empty" {
		t.Error("Channel with values is not empty")
	}
	if f, msg := Empty(TStruct{x: 1}); f != false || msg != "is Empty" {
		t.Error("struct with initialized values is empty")
	}
	if f, msg := Empty(TString("abc")); f != false || msg != "is Empty" {
		t.Error("non-empty aliased string is empty")
	}
	if f, msg := Empty(xP); f != false || msg != "is Empty" {
		t.Error("ptr to non-nil value is not empty")
	}
}

func TestTrue(t *testing.T) {
	res, msg := True(true)
	if res != true || msg != "is True" {
		t.Errorf("And operational error: %v,%s", res, msg)
	}
	res, msg = True(false)
	if res != false || msg != "is True" {
		t.Errorf("And operational error: %v,%s", res, msg)
	}

}

func TestFalse(t *testing.T) {
	res, msg := False(false)
	if res != true || msg != "is False" {
		t.Errorf("And operational error: %v,%s", res, msg)
	}
	res, msg = False(true)
	if res != false || msg != "is False" {
		t.Errorf("And operational error: %v,%s", res, msg)
	}
}

func TestZero(t *testing.T) {
	// correct test
	var res, msg = Zero(0)
	if res != true || msg != "" {
		t.Errorf("Zero error: %v,%s", res, msg)
	}
	res, msg = Zero("")
	if res != true || msg != "" {
		t.Errorf("Zero error: %v,%s", res, msg)
	}
	res, msg = Zero(nil)
	if res != true || msg != "" {
		t.Errorf("Zero error: %v,%s", res, msg)
	}

	// error test
	res, msg = Zero(1)
	if res != false || msg != "Should be zero, but was 1" {
		t.Errorf("Zero error: %v,%s", res, msg)
	}
	res, msg = Zero([]int{})
	if res != false || msg != "Should be zero, but was []" {
		t.Errorf("Zero error: %v,%s", res, msg)
	}
}

func TestNotZero(t *testing.T) {

	// error test
	var res, msg = NotZero(0)
	if res != false || msg != "Should not be zero, but was 0" {
		t.Errorf("Zero error: %v,%s", res, msg)
	}
	res, msg = NotZero("")
	if res != false || msg != "Should not be zero, but was " {
		t.Errorf("Zero error: %v,%s", res, msg)
	}
	res, msg = NotZero(nil)
	if res != false || msg != "Should not be zero, but was <nil>" {
		t.Errorf("Zero error: %v,%s", res, msg)
	}

	// correct test
	res, msg = NotZero(1)
	if res != true || msg != "" {
		t.Errorf("Zero error: %v,%s", res, msg)
	}
	res, msg = NotZero([]int{})
	if res != true || msg != "" {
		t.Errorf("Zero error: %v,%s", res, msg)
	}
}

func TestLess(t *testing.T) {
	res, msg := Less(7)(5)
	if res != true || msg != "< 7" {
		t.Errorf("And operational error: %v,%s", res, msg)
	}
	res, msg = Less(7)(7)
	if res != false || msg != "< 7" {
		t.Errorf("And operational error: %v,%s", res, msg)
	}
	res, msg = Less(7)(8)
	if res != false || msg != "< 7" {
		t.Errorf("And operational error: %v,%s", res, msg)
	}
}

func TestGreater(t *testing.T) {
	res, msg := Greater(7)(5)
	if res != false || msg != "> 7" {
		t.Errorf("And operational error: %v,%s", res, msg)
	}
	res, msg = Greater(7)(7)
	if res != false || msg != "> 7" {
		t.Errorf("And operational error: %v,%s", res, msg)
	}
	res, msg = Greater(7)(8)
	if res != true || msg != "> 7" {
		t.Errorf("And operational error: %v,%s", res, msg)
	}
}

func TestLessEq(t *testing.T) {
	res, msg := LessEq(7)(5)
	if res != true || msg != "<= 7" {
		t.Errorf("And operational error: %v,%s", res, msg)
	}
	res, msg = LessEq(7)(7)
	if res != true || msg != "<= 7" {
		t.Errorf("And operational error: %v,%s", res, msg)
	}
	res, msg = LessEq(7)(8)
	if res != false || msg != "<= 7" {
		t.Errorf("And operational error: %v,%s", res, msg)
	}
}

func TestGreaterEq(t *testing.T) {
	res, msg := GreaterEq(7)(5)
	if res != false || msg != ">= 7" {
		t.Errorf("operational error: %v,%s", res, msg)
	}
	res, msg = GreaterEq(7)(7)
	if res != true || msg != ">= 7" {
		t.Errorf("operational error: %v,%s", res, msg)
	}
	res, msg = GreaterEq(7)(8)
	if res != true || msg != ">= 7" {
		t.Errorf("operational error: %v,%s", res, msg)
	}

}

func TestEq(t *testing.T) {
	var res, msg = Eq(2)(2)
	if res != true || msg != "" {
		t.Errorf("Eq error: %v,%s", res, msg)
	}
	res, msg = Eq("2")("2")
	if res != true || msg != "" {
		t.Errorf("Eq error: %v,%s", res, msg)
	}
	res, msg = Eq([2]int{1, 2})([2]int{1, 2})
	if res != true || msg != "" {
		t.Errorf("Eq error: %v,%s", res, msg)
	}
	res, msg = Eq([]int{1, 2})([]int{1, 2})
	if res != true || msg != "" {
		t.Errorf("Eq error: %v,%s", res, msg)
	}
	res, msg = Eq(2.2)(2.2)
	if res != true || msg != "" {
		t.Errorf("Eq error: %v,%s", res, msg)
	}
	type testObj struct {
		f1 int
		f2 string
	}
	res, msg = Eq(testObj{1, "2"})(testObj{1, "2"})
	if res != true || msg != "" {
		t.Errorf("Eq error: %v,%s", res, msg)
	}
	res, msg = Eq(&testObj{1, "2"})(&testObj{1, "2"})
	if res != true || msg != "" {
		t.Errorf("Eq error: %v,%s", res, msg)
	}

	// error
	res, msg = Eq(&testObj{1, "2"})(&testObj{1, "3"})
	if res != false || msg != `== &goassert_test.testObj{f1:1, f2:"2"}` {
		t.Errorf("Eq error: %v,%s", res, msg)
	}

	res, msg = Eq(2)("2")
	if res != false || msg != `== int(2)` {
		t.Errorf("Eq error: %v,%s", res, msg)
	}
	res, msg = Eq(2)(func() {})
	if res != false || !strings.Contains(msg, "Invalid operation:") {
		t.Errorf("Eq error: %v,%s", res, msg)
	}
	res, msg = Eq(func() {})(func() {})
	if res != false || !strings.Contains(msg, "Invalid operation:") {
		t.Errorf("Eq error: %v,%s", res, msg)
	}
}

func TestNotEq(t *testing.T) {
	var res, msg = NotEq(2)(2)
	if res != false || msg != "!= 2" {
		t.Errorf("Eq error: %v,%s", res, msg)
	}
	res, msg = NotEq("2")("2")
	if res != false || msg != `!= "2"` {
		t.Errorf("Eq error: %v,%s", res, msg)
	}
	res, msg = NotEq([2]int{1, 2})([2]int{1, 2})
	if res != false || msg != "!= [2]int{1, 2}" {
		t.Errorf("Eq error: %v,%s", res, msg)
	}
	res, msg = NotEq([]int{1, 2})([]int{1, 2})
	if res != false || msg != "!= []int{1, 2}" {
		t.Errorf("Eq error: %v,%s", res, msg)
	}
	res, msg = NotEq(2.2)(2.2)
	if res != false || msg != "!= 2.2" {
		t.Errorf("Eq error: %v,%s", res, msg)
	}
	type testObj struct {
		f1 int
		f2 string
	}
	res, msg = NotEq(testObj{1, "2"})(testObj{1, "2"})
	if res != false || msg != `!= goassert_test.testObj{f1:1, f2:"2"}` {
		t.Errorf("Eq error: %v,%s", res, msg)
	}
	res, msg = NotEq(&testObj{1, "2"})(&testObj{1, "2"})
	if res != false || msg != `!= &goassert_test.testObj{f1:1, f2:"2"}` {
		t.Errorf("Eq error: %v,%s", res, msg)
	}

	res, msg = NotEq(&testObj{1, "2"})(&testObj{1, "3"})
	if res != true || msg != `` {
		t.Errorf("Eq error: %v,%s", res, msg)
	}

	res, msg = NotEq(2)("2")
	if res != true || msg != `` {
		t.Errorf("Eq error: %v,%s", res, msg)
	}
	res, msg = NotEq(2)(func() {})
	if res != false || !strings.Contains(msg, "Invalid operation:") {
		t.Errorf("Eq error: %v,%s", res, msg)
	}
	res, msg = NotEq(func() {})(func() {})
	if res != false || !strings.Contains(msg, "Invalid operation:") {
		t.Errorf("Eq error: %v,%s", res, msg)
	}
}

func TestRegexp(t *testing.T) {
	// regexp string
	var res, msg = Regexp("^hello")("hello world!")
	if res != true || msg != "" {
		t.Errorf("Regexp error: %v,%s", res, msg)
	}
	res, msg = Regexp("hello")("hello world!")
	if res != true || msg != "" {
		t.Errorf("Regexp error: %v,%s", res, msg)
	}
	res, msg = Regexp(`^hello`)("hello world!")
	if res != true || msg != "" {
		t.Errorf("Regexp error: %v,%s", res, msg)
	}
	res, msg = Regexp("^\\w")("hello world!")
	if res != true || msg != "" {
		t.Errorf("Regexp error: %v,%s", res, msg)
	}
	res, msg = Regexp("\\w!$")("hello world!")
	if res != true || msg != "" {
		t.Errorf("Regexp error: %v,%s", res, msg)
	}

	// regexp object
	var reg, _ = regexp.Compile("^hello")
	res, msg = Regexp(reg)("hello world!")
	if res != true || msg != "" {
		t.Errorf("Regexp error: %v,%s", res, msg)
	}
	reg, _ = regexp.Compile("hello")
	res, msg = Regexp(reg)("hello world!")
	if res != true || msg != "" {
		t.Errorf("Regexp error: %v,%s", res, msg)
	}
	reg, _ = regexp.Compile(`^hello`)
	res, msg = Regexp(reg)("hello world!")
	if res != true || msg != "" {
		t.Errorf("Regexp error: %v,%s", res, msg)
	}
	reg, _ = regexp.Compile("^\\w")
	res, msg = Regexp(reg)("hello world!")
	if res != true || msg != "" {
		t.Errorf("Regexp error: %v,%s", res, msg)
	}
	reg, _ = regexp.Compile("\\w!$")
	res, msg = Regexp(reg)("hello world!")
	if res != true || msg != "" {
		t.Errorf("Regexp error: %v,%s", res, msg)
	}
	reg, _ = regexp.Compile(`\w!$`)
	res, msg = Regexp(reg)("hello world!")
	if res != true || msg != "" {
		t.Errorf("Regexp error: %v,%s", res, msg)
	}

	///////////////////////////////////////////////////////////
	res, msg = Regexp("\\w!!$")("hello world!")
	if res != false || msg != `Expect "hello world!" to match "\w!!$"` {
		t.Errorf("Regexp error: %v,%s", res, msg)
	}
	reg, _ = regexp.Compile(`\w!!$`)
	res, msg = Regexp(reg)("hello world!")
	if res != false || msg != `Expect "hello world!" to match "\w!!$"` {
		t.Errorf("Regexp error: %v,%s", res, msg)
	}
}

func TestNotRegexp(t *testing.T) {
	// regexp string
	var res, msg = NotRegexp("^hello")("hello world!")
	if res != false || msg != `Expect "hello world!" to NOT match "^hello"` {
		t.Errorf("Regexp error: %v,%s", res, msg)
	}
	res, msg = NotRegexp("hello")("hello world!")
	if res != false || msg != `Expect "hello world!" to NOT match "hello"` {
		t.Errorf("Regexp error: %v,%s", res, msg)
	}
	res, msg = NotRegexp(`^hello`)("hello world!")
	if res != false || msg != `Expect "hello world!" to NOT match "^hello"` {
		t.Errorf("Regexp error: %v,%s", res, msg)
	}
	res, msg = NotRegexp("^\\w")("hello world!")
	if res != false || msg != `Expect "hello world!" to NOT match "^\w"` {
		t.Errorf("Regexp error: %v,%s", res, msg)
	}
	res, msg = NotRegexp("\\w!$")("hello world!")
	if res != false || msg != `Expect "hello world!" to NOT match "\w!$"` {
		t.Errorf("Regexp error: %v,%s", res, msg)
	}

	// regexp object
	var reg, _ = regexp.Compile("^hello")
	res, msg = NotRegexp(reg)("hello world!")
	if res != false || msg != `Expect "hello world!" to NOT match "^hello"` {
		t.Errorf("Regexp error: %v,%s", res, msg)
	}
	reg, _ = regexp.Compile("hello")
	res, msg = NotRegexp(reg)("hello world!")
	if res != false || msg != `Expect "hello world!" to NOT match "hello"` {
		t.Errorf("Regexp error: %v,%s", res, msg)
	}
	reg, _ = regexp.Compile(`^hello`)
	res, msg = NotRegexp(reg)("hello world!")
	if res != false || msg != `Expect "hello world!" to NOT match "^hello"` {
		t.Errorf("Regexp error: %v,%s", res, msg)
	}
	reg, _ = regexp.Compile("^\\w")
	res, msg = NotRegexp(reg)("hello world!")
	if res != false || msg != `Expect "hello world!" to NOT match "^\w"` {
		t.Errorf("Regexp error: %v,%s", res, msg)
	}
	reg, _ = regexp.Compile("\\w!$")
	res, msg = NotRegexp(reg)("hello world!")
	if res != false || msg != `Expect "hello world!" to NOT match "\w!$"` {
		t.Errorf("Regexp error: %v,%s", res, msg)
	}
	reg, _ = regexp.Compile(`\w!$`)
	res, msg = NotRegexp(reg)("hello world!")
	if res != false || msg != `Expect "hello world!" to NOT match "\w!$"` {
		t.Errorf("Regexp error: %v,%s", res, msg)
	}

	///////////////////////////////////////////////////////////
	res, msg = NotRegexp("\\w!!$")("hello world!")
	if res != true || msg != `` {
		t.Errorf("Regexp error: %v,%s", res, msg)
	}
	reg, _ = regexp.Compile(`\w!!$`)
	res, msg = NotRegexp(reg)("hello world!")
	if res != true || msg != `` {
		t.Errorf("Regexp error: %v,%s", res, msg)
	}
}

func TestLen(t *testing.T) {
	var res, msg = Len(2)("12")
	if res != true && msg != "" {
		t.Errorf("Len error")
	}

	res, msg = Len(2)(3)
	if res != false && msg != `"3" could not be applied builtin len()` {
		t.Errorf("Len error")
	}
}
