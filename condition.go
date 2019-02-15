package goassert

import (
	"fmt"
	"reflect"
	"strings"
)

// Condition type define
/*
Custom conditions can be used in assertions.

Example：
	var cond1 = func(actual interface{}) (bool, string) {
		// do something ...
	}

	goassert.That(t, var1).Is(cond1)
 */
type Condition func(actual interface{}) (bool, string)

/*
Logical operation：and

	And(cond1, cond2, ...)
	And(Empty, True)
	And(Empty, True, False)
*/
func And(conditions ...Condition) Condition {
	return func(actual interface{}) (bool, string) {
		for _, cond := range conditions {
			ok, msg := cond(actual)
			if !ok {
				return false, msg
			}
		}
		return true, ""
	}
}

// Logical operation：or
//
// Or(cond1, cond2, ...)
// Or(Empty, True)
// Or(Empty, True, False)
func Or(conditions ...Condition) Condition {
	return func(actual interface{}) (bool, string) {
		var msgs []string
		for _, cond := range conditions {
			ok, msg := cond(actual)
			if ok {
				return true, ""
			}
			msgs = append(msgs, msg)
		}
		return false, "(" + strings.Join(msgs, " or ") + ")"
	}
}

// Logical operation：not
//
// Not(condition)
// Not(True)
// Not(Empty)
// And(Not(Empty), Not(Nil))
func Not(condition Condition) Condition {
	return func(actual interface{}) (bool, string) {
		flag, msg := condition(actual)
		return !flag, "Not " + msg
	}
}

// Judge is nil
func Nil(actual interface{}) (bool, string) {
	return actual == nil, "is Nil"
}

// Judge is empty
func Empty(actual interface{}) (bool, string) {
	name := "is Empty"
	// get nil case out of the way
	if actual == nil {
		return true, name
	}

	objValue := reflect.ValueOf(actual)

	switch objValue.Kind() {
	// collection types are empty when they have no element
	case reflect.Array, reflect.Chan, reflect.Map, reflect.Slice:
		return objValue.Len() == 0, name
	// pointers are empty if nil or if the actual they point to is empty
	case reflect.Ptr:
		if objValue.IsNil() {
			return true, name
		}
		deref := objValue.Elem().Interface()
		return Empty(deref)
	// for all other types, compare against the zero actual
	default:
		zero := reflect.Zero(objValue.Type())
		return reflect.DeepEqual(actual, zero.Interface()), name
	}
}

// Judge is true
func True(actual interface{}) (bool, string) {
	return true == actual.(bool), "is True"
}

// Judge is false
func False(actual interface{}) (bool, string) {
	return false == actual.(bool), "is False"
}

// Judge is Zero value
func Zero(actual interface{}) (b bool, s string) {
	if actual != nil && !reflect.DeepEqual(actual, reflect.Zero(reflect.TypeOf(actual)).Interface()) {
		return false, fmt.Sprintf("Should be zero, but was %v", actual)
	}
	return true, ""
}

// Judge NOT Zero value
func NotZero(actual interface{}) (b bool, s string) {
	if actual == nil || reflect.DeepEqual(actual, reflect.Zero(reflect.TypeOf(actual)).Interface()) {
		return false, fmt.Sprintf("Should not be zero, but was %v", actual)
	}
	return true, ""
}

// Numerical operation: <
//
// Less(3)
func Less(expected interface{}) Condition {
	return func(actual interface{}) (b bool, s string) {
		b = false
		s = "< not support type: string, struct, pointer"

		if numberType(expected) && numberType(actual) {
			b, s = numberLess(expected, actual)
		}

		return
	}
}

// Numerical operation: >
//
// Greater(3)
func Greater(expected interface{}) Condition {
	return func(actual interface{}) (b bool, s string) {
		b = false
		s = "> not support type: string, struct, pointer"

		if numberType(expected) && numberType(actual) {
			b, s = numberLess(expected, actual)
			eq, _ := numberEq(expected, actual)
			b = !b && !eq
			s = strings.Replace(s, "<", ">", -1)
		}
		return
	}
}

// Numerical operation: <=
//
// LessEq(3)
func LessEq(expected interface{}) Condition {
	return func(actual interface{}) (b bool, s string) {
		b = false
		s = "<= not support type: string, struct, pointer"

		if numberType(expected) && numberType(actual) {
			less, _ := numberLess(expected, actual)
			eq, es := numberEq(expected, actual)
			b = less || eq
			s = strings.Replace(es, "==", "<=", -1)
		}
		return
	}
}

// Numerical operation: >=
//
// GreaterEq(3)
func GreaterEq(expected interface{}) Condition {
	return func(actual interface{}) (b bool, s string) {
		b = false
		s = ">= not support type: string, struct, pointer"

		if numberType(expected) && numberType(actual) {
			less, _ := numberLess(expected, actual)
			eq, es := numberEq(expected, actual)
			b = eq || !less
			s = strings.Replace(es, "==", ">=", -1)
		}
		return
	}
}

// Equal operation: ==
//
// Eq(3)
// Eq("3")
// Eq("hello")
func Eq(expected interface{}) Condition {
	return func(actual interface{}) (bool, string) {
		if err := validateEqualArgs(expected, actual); err != nil {
			return false, fmt.Sprintf("Invalid operation: %#v == %#v (%s)",
				expected, actual, err)
		}

		if !ObjectsAreEqual(expected, actual) {
			expected, actual = formatUnequalValues(expected, actual)
			return false, fmt.Sprintf("== %s", expected)
		}

		return true, ""
	}
}

// NOT Equal operation: !=
//
// NotEq(3)
// NotEq("3")
// NotEq("hello")
func NotEq(expected interface{}) Condition {
	return func(actual interface{}) (b bool, s string) {
		if err := validateEqualArgs(expected, actual); err != nil {
			return false, fmt.Sprintf("Invalid operation: %#v == %#v (%s)",
				expected, actual, err)
		}

		if ObjectsAreEqual(expected, actual) {
			expected, actual = formatUnequalValues(expected, actual)
			return false, fmt.Sprintf("!= %s", expected)
		}

		return true, ""
	}
}

// Regexp asserts that a specified regexp matches a string.
//
//  NotRegexp(regexp.MustCompile("starts"))
//  NotRegexp("^start")
func Regexp(rx interface{}) Condition {
	return func(actual interface{}) (b bool, s string) {
		match := matchRegexp(rx, actual)

		if !match {
			return false, fmt.Sprintf("Expect \"%v\" to match \"%v\"", actual, rx)
		}

		return true, ""
	}
}

// NotRegexp asserts that a specified regexp does not match a string.
//
//  NotRegexp(regexp.MustCompile("starts"))
//  NotRegexp("^start")
func NotRegexp(rx interface{}) Condition {
	return func(actual interface{}) (b bool, s string) {
		match := matchRegexp(rx, actual)

		if match {
			return false, fmt.Sprintf("Expect \"%v\" to NOT match \"%v\"", actual, rx)
		}

		return true, ""
	}
}

func Len(length int) Condition {
	return func(actual interface{}) (b bool, s string) {
		ok, l := getLen(actual)
		if !ok {
			return false, fmt.Sprintf("\"%s\" could not be applied builtin len()", actual)
		}

		if l != length {
			return false, fmt.Sprintf("\"%s\" should have %d item(s), but has %d", actual, length, l)
		}
		return true, ""
	}
}