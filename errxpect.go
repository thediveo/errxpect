package errxpect

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/onsi/gomega"
	"github.com/onsi/gomega/types"
)

// Errxpect wraps an actual multi-value allowing error assertions to be made on
// it:
//
//    func foo() (string, bool, error) { return "", false, errors.New("DOH!") }
//
//    Errxpect(foo()).To(HaveOccured())
//
// As Golang doesn't feature automatic multi-return value passing into a varargs
// function if there are additional parameters present, use WithOffset() on the
// return value of Errxpect, such as:
//    Errxpect(foo()).WithOffset(1).To(Succeed())
//
func Errxpect(actual ...interface{}) *ErrorAssertion {
	return &ErrorAssertion{
		actual: actual,
	}
}

// WithOffset replaces ExpectWithOffset() when using Errxpect in order to modify
// the call-stack offset when computing line numbers.
//
// This is most useful in helper functions that make assertions.  If you want
// Gomega's error message to refer to the calling line in the test (as opposed
// to the line in the helper function) set the argument of
// `Errxpect(...).WithOffset(offset)` appropriately.
func (errorassertion *ErrorAssertion) WithOffset(offset int) *ErrorAssertion {
	errorassertion.offset = offset
	return errorassertion
}

// ErrorAssertion implements the GomegaMatcher interface, while checking that
// actual multi-value returns are only consisting of the trailing error and all
// other return values must be zero.
type ErrorAssertion struct {
	offset int
	actual []interface{}
}

// match first checks that the actual multi-value non-error results are
// consistent with the error result (zero in case of an error, otherwise allowed
// to be non-zero). Only then is the user-specified matcher run, only with the
// trailing error result value. If necessary, the matcher is inverted when using
// with ShouldNot(), NotTo() and ToNot().
func (errorassertion *ErrorAssertion) match(matcher types.GomegaMatcher, invert bool, optionalDescription ...interface{}) bool {
	if !gomega.ExpectWithOffset(2+errorassertion.offset, errorassertion.actual).
		To(haveOnlyTrailingError(errorassertion.actual), optionalDescription...) {
		// Gosh! There is something rotten with the return values: a non-nil
		// error, yet the other return values aren't all zero.
		return false
	}
	if invert {
		return gomega.ExpectWithOffset(
			2+errorassertion.offset, errorassertion.actual[len(errorassertion.actual)-1]).
			To(gomega.Not(matcher), optionalDescription...)
	}
	return gomega.ExpectWithOffset(
		2+errorassertion.offset, errorassertion.actual[len(errorassertion.actual)-1]).
		To(matcher, optionalDescription...)
}

// Should makes an assertion that should be true.
func (errorassertion *ErrorAssertion) Should(matcher types.GomegaMatcher, optionalDescription ...interface{}) bool {
	return errorassertion.match(matcher, false, optionalDescription...)
}

// ShouldNot makes an assertion that should not be true.
func (errorassertion *ErrorAssertion) ShouldNot(matcher types.GomegaMatcher, optionalDescription ...interface{}) bool {
	return errorassertion.match(matcher, true, optionalDescription...)
}

// To makes an assertion that should be true.
func (errorassertion *ErrorAssertion) To(matcher types.GomegaMatcher, optionalDescription ...interface{}) bool {
	return errorassertion.match(matcher, false, optionalDescription...)
}

// NotTo makes an assertion that should not be true.
func (errorassertion *ErrorAssertion) NotTo(matcher types.GomegaMatcher, optionalDescription ...interface{}) bool {
	return errorassertion.match(matcher, true, optionalDescription...)
}

// ToNot makes an assertion that should not be true.
func (errorassertion *ErrorAssertion) ToNot(matcher types.GomegaMatcher, optionalDescription ...interface{}) bool {
	return errorassertion.match(matcher, true, optionalDescription...)
}

func haveOnlyTrailingError(actuals []interface{}) types.GomegaMatcher {
	return &trailingErrorMatcher{actuals: actuals}
}

type trailingErrorMatcher struct {
	actuals []interface{}
}

// Match vets the actual (function return) values: (1) all values, except for
// the trailing error value must be zero OR (2) if the trailing error is nil,
// then the other values are allowed to be non-zero.
//
// This matcher is kind of the opposite of Gomega's internal vetExtras(): while
// Gomega catches the case where a multi-return values function returns an
// error, we catch we case that when we expect errors, a multi-return function
// additionally returns actual values other than just the error value.
func (te *trailingErrorMatcher) Match(interface{}) (bool, error) {
	// No return values at all?! That's ... unexpected!
	if len(te.actuals) == 0 {
		return false, errors.New("No return values")
	}
	// If the final error value is zero, then we allow the other return values
	// to be whatever they like.
	err := te.actuals[len(te.actuals)-1]
	if err == nil {
		return true, nil
	}
	zeroValue := reflect.Zero(reflect.TypeOf(err)).Interface()
	if reflect.DeepEqual(zeroValue, err) {
		return true, nil
	}
	// There's an error, so all other return values must be zero.
	for idx, actual := range te.actuals[:len(te.actuals)-1] {
		if actual != nil {
			zeroValue := reflect.Zero(reflect.TypeOf(actual)).Interface()
			if !reflect.DeepEqual(zeroValue, actual) {
				return false, fmt.Errorf(
					"Unexpected non-nil/non-zero actual non-error argument at index %d:\n\t<%T>: %#v",
					idx+1, actual, actual)
			}
		}
	}
	return true, nil
}

func (te *trailingErrorMatcher) FailureMessage(actual interface{}) string        { return "" }
func (te *trailingErrorMatcher) NegatedFailureMessage(actual interface{}) string { return "" }
