package validation

import "testing"

func TestIsEmailValid(t *testing.T) {
	if !IsEmailValid("a@b.com") {
		t.Error(`IsEmailValid("a@b.com") = false`)
	}
	if !IsEmailValid("owalmsley@sheffield.ac.uk") {
		t.Error(`IsEmailValid("owalmsley@sheffield.ac.uk") = false`)
	}
}
func TestNotIsEmailValid(t *testing.T) {
	if IsEmailValid("ab.com") {
		t.Error(`IsEmailValid("ab.com") = true`)
	}
}
