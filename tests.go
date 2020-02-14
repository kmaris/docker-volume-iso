package main

import testing

func NewIsoDriver(t *testing.T) {
	testsLocation := "tests"

	got := newIsoDriver(testsLocation)
	if got.volumesRoot != testsLocation {
		t.Error("newIsoDriver(%s) = %s; want %s", testsLocation, got.volumesRoot, testsLocation)
	}
}

