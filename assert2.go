package ar

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
	"testing"
)

type Test struct {
	T       *testing.T
	F       *FTest
	section string
}
type FTest struct {
	T       *testing.T
	section string
}

func (test *Test) Section(s string) {
	test.section = s
	test.F.section = s
}
func (test *Test) logSection() {
	if test.section != "" {
		test.T.Logf("Nil check error in %s Section")
	}
}
func (test *FTest) logSection() {
	if test.section != "" {
		test.T.Logf("Nil check error in %s Section")
	}
}

// Nil tests
func (test *Test) IsNil(v interface{}, msgs ...interface{}) {
	if !testIsNil(v) {
		test.logSection()
		test.T.Error(msgs...)
	}
}
func (test *Test) IsNotNil(v interface{}, msgs ...interface{}) {
	if testIsNil(v) {
		test.logSection()
		test.T.Error(msgs...)
	}
}

// bool tests
func (test *Test) IsTrue(b bool, msgs ...interface{}) {
	if !b {
		test.logSection()
		test.T.Error(msgs...)
	}
}
func (test *Test) IsFalse(b bool, msgs ...interface{}) {
	if b {
		test.logSection()
		test.T.Error(msgs...)
	}
}

// Equality test
func (test *Test) AreEqual(x, y interface{}, msgs ...interface{}) {
	if !(reflect.DeepEqual(x, y) || strEqual(x, y)) {
		test.logSection()
		test.T.Log("Equality check failed")
		test.T.Error(msgs...)
	}
}
func (test *Test) AreNotEqual(x, y interface{}, msgs ...interface{}) {
	if reflect.DeepEqual(x, y) || strEqual(x, y) {
		test.logSection()
		test.T.Log("Inequality check failed")
		test.T.Error(msgs...)
	}
}

// String tests
func (test *Test) StartsWith(s, pre string, msgs ...interface{}) {
	if !strings.HasPrefix(s, pre) {
		test.logSection()
		test.T.Error(msgs...)
	}
}
func (test *Test) EndsWith(s, post string, msgs ...interface{}) {
	if !strings.HasSuffix(s, post) {
		test.logSection()
		test.T.Error(msgs...)
	}
}
func (test *Test) Matches(s, regex string, msgs ...interface{}) {
	matches, err := regexp.MatchString(regex, s)
	if err != nil {
		test.logSection()
		panic(err)
	} else if !matches {
		test.logSection()
		test.T.Error(msgs...)
	}
}
func (test *Test) NotMatches(s, regex string, msgs ...interface{}) {
	matches, err := regexp.MatchString(regex, s)
	if err != nil {
		test.logSection()
		panic(err)
	} else if matches {
		test.logSection()
		test.T.Error(msgs...)
	}
}

// Nil Format tests
func (test *FTest) IsNil(v interface{}, msgFormat string, msgs ...interface{}) {
	if !testIsNil(v) {
		test.logSection()
		test.T.Errorf(msgFormat, msgs...)
	}
}
func (test *FTest) IsNotNil(v interface{}, msgFormat string, msgs ...interface{}) {
	if testIsNil(v) {
		test.logSection()
		test.T.Errorf(msgFormat, msgs...)
	}
}

// bool tests
func (test *FTest) IsTrue(b bool, msgFormat string, msgs ...interface{}) {
	if !b {
		test.logSection()
		test.T.Errorf(msgFormat, msgs...)
	}
}
func (test *FTest) IsFalse(b bool, msgFormat string, msgs ...interface{}) {
	if b {
		test.logSection()
		test.T.Errorf(msgFormat, msgs...)
	}
}

// Equality test
func (test *FTest) AreEqual(x, y interface{}, msgFormat string, msgs ...interface{}) {
	if !(reflect.DeepEqual(x, y) || strEqual(x, y)) {
		test.logSection()
		test.T.Errorf(msgFormat, msgs...)
	}
}
func (test *FTest) AreNotEqual(x, y interface{}, msgFormat string, msgs ...interface{}) {
	if reflect.DeepEqual(x, y) || strEqual(x, y) {
		test.logSection()
		test.T.Errorf(msgFormat, msgs...)
	}
}

// String tests
func (test *FTest) StartsWith(s, pre, msgFormat string, msgs ...interface{}) {
	if !strings.HasPrefix(s, pre) {
		test.logSection()
		test.T.Errorf(msgFormat, msgs...)
	}
}
func (test *FTest) EndsWith(s, post, msgFormat string, msgs ...interface{}) {
	if !strings.HasSuffix(s, post) {
		test.T.Errorf(msgFormat, msgs...)
	}
}
func (test *FTest) Matches(s, regex, msgFormat string, msgs ...interface{}) {
	matches, err := regexp.MatchString(regex, s)
	if err != nil {
		test.logSection()
		panic(err)
	} else if !matches {
		test.logSection()
		test.T.Errorf(msgFormat, msgs...)
	}
}
func (test *FTest) NotMatches(s, regex, msgFormat string, msgs ...interface{}) {
	matches, err := regexp.MatchString(regex, s)
	if err != nil {
		test.logSection()
		panic(err)
	} else if matches {
		test.logSection()
		test.T.Errorf(msgFormat, msgs...)
	}
}

func testIsNil(v interface{}) bool {
	return v == nil || reflect.ValueOf(v).IsNil()
}

func strEqual(x, y interface{}) bool {
	return fmt.Sprint(x) == fmt.Sprint(y)
}

func Within(t *testing.T, f func(*Test)) {
	f(&Test{T: t, F: &FTest{T: t}})
}
