// Copyright 2019 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package diff

import (
	"strings"
	"testing"
)

var formatTests = []struct {
	text1          string
	text2          string
	diff           string
	suppressCommon string
}{
	{"a b c", "a b d e f", "a b -c +d +e +f", "3-c 3+d 4+e 5+f"},
	{"", "a b c", "+a +b +c", "1+a 2+b 3+c"},
	{"a b c", "", "-a -b -c", "1-a 2-b 3-c"},
	{"a b c", "d e f", "-a -b -c +d +e +f", "1-a 2-b 3-c 1+d 2+e 3+f"},
	{"a b c d e f", "a b d e f", "a b -c d e f", "3-c"},
	{"a b c e f", "a b c d e f", "a b c +d e f", "4+d"},
}

func TestFormat(t *testing.T) {
	for _, tt := range formatTests {
		// Turn spaces into \n.
		text1 := strings.ReplaceAll(tt.text1, " ", "\n")
		if text1 != "" {
			text1 += "\n"
		}
		text2 := strings.ReplaceAll(tt.text2, " ", "\n")
		if text2 != "" {
			text2 += "\n"
		}
		compare(t, format, text1, text2, tt.diff)
		compare(t, suppressCommon, text1, text2, tt.suppressCommon)
	}
}

func format(text1, text2 string) string {
	return Format(text1, text2)
}

func suppressCommon(text1, text2 string) string {
	return Format(text1, text2, OptSuppressCommon())
}

func compare(t *testing.T, testFn func(string, string) string, text1, text2 string, want string) {
	t.Helper()
	got := testFn(text1, text2)
	// Cut final \n, cut spaces, turn remaining \n into spaces.
	got = strings.ReplaceAll(strings.ReplaceAll(strings.TrimSuffix(got, "\n"), " ", ""), "\n", " ")
	if got != want {
		t.Errorf("diff(%q, %q) = %q, want %q", text1, text2, got, want)
	}
}
