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
//    Errxpect(func() (string, bool, error) { return "", false, errors.New("DOH!") }()).To(Equal("foo"))
func Errxpect(actual ...interface{}) gomega.Assertion {
	return ErrxpectWithOffset(0, actual...)
}

// ErrxpectWithOffset wraps an actual multi-value allowing error assertions to
// be made on it:
//
//    ErrxpectWithOffset(func() (string, bool, error) { return "", false, errors.New("DOH!") }()).To(Equal("foo"))
//
// Unlike `Errxpect`, `ErrxpectWithOffset` takes an additional integer argument
// that is used to modify the call-stack offset when computing line numbers.
//
// This is most useful in helper functions that make assertions.  If you want
// Gomega's error message to refer to the calling line in the test (as opposed
// to the line in the helper function) set the first argument of
// `ErrxpectWithOffset` appropriately.
func ErrxpectWithOffset(offset int, actual ...interface{}) gomega.Assertion {
	return &errorAssertion{
		offset: offset,
		actual: actual,
	}
}

// errorAssertion implements the GomegaMatcher interface, while checking that
// actual multi-value returns are only consisting of the trailing error and all
// other return values must be zero.
type errorAssertion struct {
	offset int
	actual []interface{}
}

// match first checks that the actual multi-value non-error results are
// consistent with the error result (zero in case of an error, otherwise allowed
// to be non-zero). Only then is the user-specified matcher run, only with the
// trailing error result value. If necessary, the matcher is inverted when using
// with ShouldNot(), NotTo() and ToNot().
func (errorassertion *errorAssertion) match(matcher types.GomegaMatcher, invert bool, optionalDescription ...interface{}) bool {
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

func (errorassertion *errorAssertion) Should(matcher types.GomegaMatcher, optionalDescription ...interface{}) bool {
	return errorassertion.match(matcher, false, optionalDescription...)
}

func (errorassertion *errorAssertion) ShouldNot(matcher types.GomegaMatcher, optionalDescription ...interface{}) bool {
	return errorassertion.match(matcher, true, optionalDescription...)
}

func (errorassertion *errorAssertion) To(matcher types.GomegaMatcher, optionalDescription ...interface{}) bool {
	return errorassertion.match(matcher, false, optionalDescription...)
}

func (errorassertion *errorAssertion) NotTo(matcher types.GomegaMatcher, optionalDescription ...interface{}) bool {
	return errorassertion.match(matcher, true, optionalDescription...)
}

func (errorassertion *errorAssertion) ToNot(matcher types.GomegaMatcher, optionalDescription ...interface{}) bool {
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
