package test

import "testing"

// AssertEqual asserts that two objects are equal.
func AssertEqual(t *testing.T, a interface{}, b interface{}) {
	if a != b {
		t.Fatalf("%s != %s", a, b)
	}
}
