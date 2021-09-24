package main

import "testing"

func TestNumExtend(t *testing.T) {
	testData := []struct {
		str    string
		num    int
		result string
	}{
		{"", 0, ""},
		{"aaa", 0, "aaa"},
		{"aaasdsad", 0, "aaasdsad"},
		{"a.bb.c", 0, "a.bb.c"},
		{"aaa", 10, "aaa1"},
		{"", 10, "1"},
		{"aaa", 120, "aaa12"},
		{"aa.a", 120, "aa12.a"},
		{"a.bb.c", 20, "a.bb2.c"},
	}
	for _, d := range testData {
		got := numExtend(d.str, d.num)
		if got != d.result {
			t.Errorf("numExtend(%q, %d) = %q; want %q", d.str, d.num, got, d.result)
		}
	}
}
