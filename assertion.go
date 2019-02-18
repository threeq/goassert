package goassert

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strings"
)

// TestingT is an interface wrapper around *testing.T
type TestingT interface {
	Errorf(format string, args ...interface{})
}

type FluentAssertion struct {
	t      TestingT
	actual interface{}
	name   string
}

func (assert *FluentAssertion) That(actual interface{}) *FluentAssertion {
	return &FluentAssertion{
		assert.t,
		actual,
		"",
	}
}

func (assert *FluentAssertion) Equal(expected interface{}, msgAndArgs ...interface{}) *FluentAssertion {

	actual := assert.actual
	res, msg := Eq(expected)(actual)
	if !res {
		if strings.HasPrefix(msg, "Invalid operation:") {
			Fail(assert, msg, msgAndArgs...)
		} else {
			diff := diff(expected, actual)
			expected, actual = formatUnequalValues(expected, actual)
			Fail(assert, fmt.Sprintf("Not equal: \n"+
				"expected: %s\n"+
				"actual  : %s%s", expected, actual, diff), msgAndArgs...)
		}
	}
	return assert
}

func (assert *FluentAssertion) EqualIgnoringCase(expected interface{}, msgAndArgs ...interface{}) *FluentAssertion {
	actual,ok1 := assert.actual.(string)
	exp,ok2 := expected.(string)

	if !ok1 {
		Fail(assert, "actual value is not string type")
		return assert
	}

	if !ok2 {
		Fail(assert, "expected value is not string type")
		return assert
	}

	if !strings.EqualFold(actual, exp) {
		Fail(assert, fmt.Sprintf("Not equal: \n"+
			"expected: %s\n"+
			"actual  : %s", expected, actual), msgAndArgs...)
	}
	return assert
}

func (assert *FluentAssertion) NotEqual(expected interface{}, msgAndArgs ...interface{}) *FluentAssertion {

	res, msg := NotEq(expected)(assert.actual)
	if !res {
		if strings.HasPrefix(msg, "Invalid operation:") {
			Fail(assert, msg, msgAndArgs...)
		} else {
			Fail(assert, fmt.Sprintf("Should not be: %#v\n", assert.actual), msgAndArgs...)
		}
	}
	return assert
}

// As() is used to describe the test and will be shown before the error message
func (assert *FluentAssertion) As(desc string) *FluentAssertion {
	assert.name = desc
	return assert
}

func (assert *FluentAssertion) StartsWith(expected interface{}, msgAndArgs ...interface{}) *FluentAssertion {
	et, ek := typeAndKind(expected)
	at, _ := typeAndKind(assert.actual)

	if et != at {
		Fail(assert, "Two different types are not supported", msgAndArgs...)
		return assert
	}

	if ek != reflect.Slice && ek != reflect.Array && ek != reflect.String {
		Fail(assert, "Unsupported type", msgAndArgs...)
		return assert
	}

	if ek == reflect.String {
		ab, eb := []byte(assert.actual.(string)), []byte(expected.(string))

		if !bytes.HasPrefix(ab, eb) {
			Fail(assert, fmt.Sprintf("Not startsWith: \n"+
				"expected: %s\n"+
				"actual  : %s", expected, assert.actual), msgAndArgs...)
			return assert
		}
	} else {
		as := reflect.ValueOf(assert.actual)
		es := reflect.ValueOf(expected)

		if as.Len()<es.Len() {
			Fail(assert, fmt.Sprintf("Not startsWith: \n"+
				"expected: %#v\n"+
				"actual  : %#v", expected, assert.actual), msgAndArgs...)
			return assert
		}

		for i := 0; i < es.Len(); i++ {
			if !ObjectsAreEqual(es.Index(i).Interface(), as.Index(i).Interface()) {
				Fail(assert, fmt.Sprintf("Not startsWith: \n"+
					"expected: %#v\n"+
					"actual  : %#v", expected, assert.actual), msgAndArgs...)
				return assert
			}
		}
	}
	return assert
}

func (assert *FluentAssertion) EndsWith(expected interface{}, msgAndArgs ...interface{}) *FluentAssertion {
	et, ek := typeAndKind(expected)
	at, _ := typeAndKind(assert.actual)

	if et != at {
		Fail(assert, "Two different types are not supported", msgAndArgs...)
		return assert
	}

	if ek != reflect.Slice && ek != reflect.Array && ek != reflect.String {
		Fail(assert, "Unsupported type", msgAndArgs...)
		return assert
	}

	if ek == reflect.String {
		ab, eb := []byte(assert.actual.(string)), []byte(expected.(string))

		if !bytes.HasSuffix(ab, eb) {
			Fail(assert, fmt.Sprintf("Not endsWith: \n"+
				"expected: %s\n"+
				"actual  : %s", expected, assert.actual), msgAndArgs...)
			return assert
		}
	} else {
		as := reflect.ValueOf(assert.actual)
		es := reflect.ValueOf(expected)

		if as.Len()<es.Len() {
			Fail(assert, fmt.Sprintf("Not endsWith: \n"+
				"expected: %#v\n"+
				"actual  : %#v", expected, assert.actual), msgAndArgs...)
			return assert
		}

		for i := 0; i < es.Len(); i++ {
			ei := es.Len()-i-1
			ai := as.Len()-i-1
			if !ObjectsAreEqual(es.Index(ei).Interface(), as.Index(ai).Interface()) {
				Fail(assert, fmt.Sprintf("Not endsWith: \n"+
					"expected: %#v\n"+
					"actual  : %#v", expected, assert.actual), msgAndArgs...)
				return assert
			}
		}
	}
	return assert
}

// Len asserts that the specified object has specific length.
// Len also fails if the object has a type that len() not accept.
func (assert *FluentAssertion) Len(length int, msgAndArgs ...interface{}) *FluentAssertion {
	ok, l := getLen(assert.actual)
	if !ok {
		Fail(assert, fmt.Sprintf("\"%s\" could not be applied builtin len()", assert.actual), msgAndArgs...)
	}

	if l != length {
		Fail(assert, fmt.Sprintf("\"%s\" should have %d item(s), but has %d", assert.actual, length, l), msgAndArgs...)
	}
	return assert
}


// Contains asserts that the specified string, list(array, slice...) or map contains the
// specified substring or element.
//
//    assert.Contains("World")
func (assert *FluentAssertion) Contains(expected interface{}, msgAndArgs ...interface{}) *FluentAssertion {
	ok, found := includeElement(assert.actual, expected)
	if !ok {
		Fail(assert, fmt.Sprintf("\"%s\" could not be applied builtin len()", assert.actual), msgAndArgs...)
	}
	if !found {
		Fail(assert, fmt.Sprintf("\"%s\" does not contain \"%s\"", assert.actual, expected), msgAndArgs...)
	}

	return assert
}

func (assert *FluentAssertion) NotContain(expected interface{}, msgAndArgs ...interface{}) *FluentAssertion {
	ok, found := includeElement(assert.actual, expected)
	if !ok {
		Fail(assert, fmt.Sprintf("\"%s\" could not be applied builtin len()", assert.actual), msgAndArgs...)
	}
	if found {
		Fail(assert, fmt.Sprintf("\"%s\" does contain \"%s\"", expected, assert.actual), msgAndArgs...)
	}

	return assert
}

func (assert *FluentAssertion) In(expected interface{}, msgAndArgs ...interface{}) *FluentAssertion {
	ok, found := includeElement(expected, assert.actual)
	if !ok {
		Fail(assert, fmt.Sprintf("\"%s\" could not be applied builtin len()", expected), msgAndArgs...)
	}
	if !found {
		Fail(assert, fmt.Sprintf("\"%s\" does not in \"%s\"", assert.actual, expected), msgAndArgs...)
	}

	return assert
}

func (assert *FluentAssertion) NotIn(expected interface{}, msgAndArgs ...interface{}) *FluentAssertion {
	ok, found := includeElement(expected, assert.actual)
	if !ok {
		Fail(assert, fmt.Sprintf("\"%s\" could not be applied builtin len()", expected), msgAndArgs...)
	}
	if found {
		Fail(assert, fmt.Sprintf("\"%s\" does in \"%s\"", assert.actual, expected), msgAndArgs...)
	}

	return assert
}

func (assert *FluentAssertion) HasMessage(expected string, msgAndArgs ...interface{}) *FluentAssertion {
	if assert.actual == nil {
		Fail(assert, "An error is expected but got nil.", msgAndArgs...)
		return assert
	}
	theError, ok := assert.actual.(error)
	if !ok {
		Fail(assert, "Object is not error type.", msgAndArgs...)
		return assert
	}
	actual := theError.Error()
	// don't need to use deep equals here, we know they are both strings
	if expected != actual {
		Fail(assert, fmt.Sprintf("Error message not equal:\n"+
			"expected: %q\n"+
			"actual  : %q", expected, actual), msgAndArgs...)
	}
	return assert
}

func (assert *FluentAssertion) IsType(expectedType interface{}, msgAndArgs ...interface{}) *FluentAssertion {
	if !ObjectsAreEqual(reflect.TypeOf(assert.actual), reflect.TypeOf(expectedType)) {
		Fail(assert, fmt.Sprintf("Object expected to be of type %v, but was %v",
			reflect.TypeOf(expectedType), reflect.TypeOf(assert.actual)), msgAndArgs...)
	}
	return assert
}

func (assert *FluentAssertion) Implements(interfaceObject interface{}, msgAndArgs ...interface{}) *FluentAssertion {
	interfaceType := reflect.TypeOf(interfaceObject).Elem()

	if assert.actual == nil {
		Fail(assert, fmt.Sprintf("Cannot check if nil implements %v", interfaceType), msgAndArgs...)
		return assert
	}
	if !reflect.TypeOf(assert.actual).Implements(interfaceType) {
		Fail(assert, fmt.Sprintf("%T must implement %v", assert.actual, interfaceType), msgAndArgs...)
		return assert
	}
	return assert
}

func (assert *FluentAssertion) Is(condition Condition, msgAndArgs ...interface{}) *FluentAssertion {
	result, msg := condition(assert.actual)
	if !result {
		Fail(assert, fmt.Sprintf("Should Is: \n"+
			"expected  : True\n"+
			"actual    : False\n"+
			"condition : %s\n"+
			"value     : %s", msg, assert.actual), msgAndArgs...)
	}
	return assert
}

func (assert *FluentAssertion) Not(condition Condition, msgAndArgs ...interface{}) *FluentAssertion {
	result, msg := condition(assert.actual)
	if result {
		Fail(assert, fmt.Sprintf("Should NOT Is: \n"+
			"expected  : True\n"+
			"actual    : False\n"+
			"condition : Not %s\n"+
			"value     : %s", msg, assert.actual), msgAndArgs...)
	}
	return assert
}

func (assert *FluentAssertion) AllOf(conditions ...Condition) *FluentAssertion {
	result, msg := And(conditions...)(assert.actual)
	if !result {
		Fail(assert, fmt.Sprintf("Should Is: \n"+
			"expected  : True\n"+
			"actual    : False\n"+
			"condition : AllOf %s\n"+
			"value     : %s", msg, assert.actual))
	}
	return assert
}

func (assert *FluentAssertion) AnyOf(conditions ...Condition) *FluentAssertion {
	result, msg := Or(conditions...)(assert.actual)
	if !result {
		Fail(assert, fmt.Sprintf("Should Is: \n"+
			"expected  : True\n"+
			"actual    : False\n"+
			"condition : AnyOf %s\n"+
			"value     : %s", msg, assert.actual))
	}
	return assert
}

// Panics asserts that the code inside the specified PanicTestFunc panics.
//
//   assert.Panics(t, func(){ GoCrazy() })
func (assert *FluentAssertion) Panics(f PanicTestFunc, msgAndArgs ...interface{}) *FluentAssertion {

	if funcDidPanic, panicValue := didPanic(f); !funcDidPanic {
		Fail(assert, fmt.Sprintf("func %#v should panic\n\tPanic value:\t%#v", f, panicValue), msgAndArgs...)
	}

	return assert
}


// Regexp asserts that a specified regexp matches a string.
//
//  assert.Regexp(t, regexp.MustCompile("start"), "it's starting")
//  assert.Regexp(t, "start...$", "it's not starting")
func (assert *FluentAssertion) Regexp(rx interface{}, msgAndArgs ...interface{}) *FluentAssertion {

	match := matchRegexp(rx, assert.actual)

	if !match {
		Fail(assert, fmt.Sprintf("Expect \"%v\" to match \"%v\"", assert.actual, rx), msgAndArgs...)
	}

	return assert
}

// NotRegexp asserts that a specified regexp does not match a string.
//
//  assert.NotRegexp(regexp.MustCompile("starts"), "it's starting")
//  assert.NotRegexp("^start", "it's not starting")
func (assert *FluentAssertion) NotRegexp(rx interface{}, msgAndArgs ...interface{}) *FluentAssertion {
	match := matchRegexp(rx, assert.actual)

	if match {
		Fail(assert, fmt.Sprintf("Expect \"%v\" to NOT match \"%v\"", assert.actual, rx), msgAndArgs...)
	}

	return assert

}

// Zero asserts that i is the zero value for its type.
func (assert *FluentAssertion) Zero(msgAndArgs ...interface{}) *FluentAssertion {
	i := assert.actual
	if i != nil && !reflect.DeepEqual(i, reflect.Zero(reflect.TypeOf(i)).Interface()) {
		Fail(assert, fmt.Sprintf("Should be zero, but was %v", i), msgAndArgs...)
	}
	return assert
}

// NotZero asserts that i is not the zero value for its type.
func (assert *FluentAssertion) NotZero(msgAndArgs ...interface{}) *FluentAssertion {
	i := assert.actual
	if i == nil || reflect.DeepEqual(i, reflect.Zero(reflect.TypeOf(i)).Interface()) {
		Fail(assert, fmt.Sprintf("Should not be zero, but was %v", i), msgAndArgs...)
	}
	return assert
}


// FileExists checks whether a file exists in the given path. It also fails if the path points to a directory or there is an error when trying to check the file.
func (assert *FluentAssertion) FileExists(msgAndArgs ...interface{}) *FluentAssertion {
	path := assert.actual.(string)
	info, err := os.Lstat(path)
	if err != nil {
		if os.IsNotExist(err) {
			Fail(assert, fmt.Sprintf("unable to find file %q", path), msgAndArgs...)
			return assert
		}
		Fail(assert, fmt.Sprintf("error when running os.Lstat(%q): %s", path, err), msgAndArgs...)
		return assert
	}
	if info.IsDir() {
		Fail(assert, fmt.Sprintf("%q is a directory", path), msgAndArgs...)
		return assert
	}
	return assert
}

// DirExists checks whether a directory exists in the given path. It also fails if the path is a file rather a directory or there is an error checking whether it exists.
func (assert *FluentAssertion) DirExists(msgAndArgs ...interface{}) *FluentAssertion {
	path := assert.actual.(string)
	info, err := os.Lstat(path)
	if err != nil {
		if os.IsNotExist(err) {
			Fail(assert, fmt.Sprintf("unable to find file %q", path), msgAndArgs...)
			return assert
		}
		Fail(assert, fmt.Sprintf("error when running os.Lstat(%q): %s", path, err), msgAndArgs...)
		return assert
	}
	if !info.IsDir() {
		Fail(assert, fmt.Sprintf("%q is a file", path), msgAndArgs...)
		return assert
	}
	return assert
}

// JSONEq asserts that two JSON strings are equivalent.
//
//  assert.JSONEq(t, `{"hello": "world", "foo": "bar"}`, `{"foo": "bar", "hello": "world"}`)
func (assert *FluentAssertion) JSONEq(expected string, msgAndArgs ...interface{}) *FluentAssertion {

	var expectedJSONAsInterface, actualJSONAsInterface interface{}

	if err := json.Unmarshal([]byte(expected), &expectedJSONAsInterface); err != nil {
		 Fail(assert, fmt.Sprintf("Expected value ('%s') is not valid json.\nJSON parsing error: '%s'", expected, err.Error()), msgAndArgs...)
		 return assert
	}

	actual := assert.actual.(string)
	if err := json.Unmarshal([]byte(actual), &actualJSONAsInterface); err != nil {
		Fail(assert, fmt.Sprintf("Input ('%s') needs to be valid json.\nJSON parsing error: '%s'", actual, err.Error()), msgAndArgs...)
		return assert
	}

	assert.That(actualJSONAsInterface).Equal(expectedJSONAsInterface, msgAndArgs...)
	return assert
}

type assertProxy struct {
	t TestingT
}

func (tp *assertProxy) That(that interface{}) *FluentAssertion {
	return &FluentAssertion{
		tp.t,
		that,
		"",
	}
}

// Panics asserts that the code inside the specified PanicTestFunc panics.
//
//   assert.Panics(t, func(){ GoCrazy() })
func (tp *assertProxy) Panics(f PanicTestFunc, msgAndArgs ...interface{}) *assertProxy {

	assert := &FluentAssertion{
		tp.t,
		nil,
		"",
	}
	if funcDidPanic, panicValue := didPanic(f); !funcDidPanic {
		Fail(assert, fmt.Sprintf("func %#v should panic\n\tPanic value:\t%#v", f, panicValue), msgAndArgs...)
	}

	return tp
}

func Assertion(t TestingT) *assertProxy {
	return &assertProxy{
		t,
	}
}

func That(t TestingT, actual interface{}) *FluentAssertion {
	return &FluentAssertion{
		t,
		actual,
		"",
	}
}
