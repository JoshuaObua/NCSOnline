package handlers

import "testing"

func TestOptionalDate(t *testing.T) {
	if v, err := parseOptionalDate(""); err != nil || v != nil {
		t.Fatal("blank date should be optional")
	}
	if _, err := parseOptionalDate("18/06/2026"); err == nil {
		t.Fatal("non-ISO date should fail")
	}
	if v, err := parseOptionalDate("2026-06-18"); err != nil || v == nil {
		t.Fatal("ISO date should parse")
	}
}
