// Copyright 2019 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package diff provides basic text comparison (like Unix's diff(1)).
package diff

import (
	"fmt"
	"strings"
)

type config struct {
	suppressCommon bool
}

type option func(*config)

// OptSuppressCommon suppresses common lines.
func OptSuppressCommon() option {
	return func(c *config) {
		c.suppressCommon = true
	}
}

// Format returns a formatted diff of the two texts,
// showing the entire text and the minimum line-level
// additions and removals to turn text1 into text2.
// (That is, lines only in text1 appear with a leading -,
// and lines only in text2 appear with a leading +.)
func Format(text1, text2 string, options ...option) string {
	var c config
	for _, option := range options {
		option(&c)
	}
	if text1 != "" && !strings.HasSuffix(text1, "\n") {
		text1 += "(missing final newline)"
	}
	lines1 := strings.Split(text1, "\n")
	lines1 = lines1[:len(lines1)-1] // remove empty string after final line
	if text2 != "" && !strings.HasSuffix(text2, "\n") {
		text2 += "(missing final newline)"
	}
	lines2 := strings.Split(text2, "\n")
	lines2 = lines2[:len(lines2)-1] // remove empty string after final line

	// Naive dynamic programming algorithm for edit distance.
	// https://en.wikipedia.org/wiki/Wagner–Fischer_algorithm
	// dist[i][j] = edit distance between lines1[:len(lines1)-i] and lines2[:len(lines2)-j]
	// (The reversed indices make following the minimum cost path
	// visit lines in the same order as in the text.)
	dist := make([][]int, len(lines1)+1)
	for i := range dist {
		dist[i] = make([]int, len(lines2)+1)
		if i == 0 {
			for j := range dist[0] {
				dist[0][j] = j
			}
			continue
		}
		for j := range dist[i] {
			if j == 0 {
				dist[i][0] = i
				continue
			}
			cost := dist[i][j-1] + 1
			if cost > dist[i-1][j]+1 {
				cost = dist[i-1][j] + 1
			}
			if lines1[len(lines1)-i] == lines2[len(lines2)-j] {
				if cost > dist[i-1][j-1] {
					cost = dist[i-1][j-1]
				}
			}
			dist[i][j] = cost
		}
	}

	var buf strings.Builder
	i, j := len(lines1), len(lines2)
	for i > 0 || j > 0 {
		cost := dist[i][j]
		if i > 0 && j > 0 && cost == dist[i-1][j-1] && lines1[len(lines1)-i] == lines2[len(lines2)-j] {
			if !c.suppressCommon {
				k := len(lines1) - i
				fmt.Fprintf(&buf, " %s\n", lines1[k])
			}
			i--
			j--
		} else if i > 0 && cost == dist[i-1][j]+1 {
			k := len(lines1) - i
			if c.suppressCommon {
				fmt.Fprint(&buf, k+1)
			}
			fmt.Fprintf(&buf, "-%s\n", lines1[k])
			i--
		} else {
			k := len(lines2) - j
			if c.suppressCommon {
				fmt.Fprint(&buf, k+1)
			}
			fmt.Fprintf(&buf, "+%s\n", lines2[k])
			j--
		}
	}
	return buf.String()
}
