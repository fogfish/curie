//
// Copyright (C) 2020 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/fogfish/curie
//

package reference

import "strings"

// IRelative Ref as defined in IRI, RFC 3987

func Join(ref string, delim rune, segments ...string) string {
	if len(segments) == 0 {
		return ref
	}

	// custom strings.Join that ignores empty segments
	n := (len(segments) - 1)
	for i := 0; i < len(segments); i++ {
		n += len(segments[i])
	}

	if n == (len(segments) - 1) {
		return ref
	}

	var b strings.Builder
	b.Grow(n)
	b.WriteString(ref)

	for _, s := range segments {
		if s != "" {
			if b.Len() != 0 {
				b.WriteRune(delim)
			}
			b.WriteString(s)
		}
	}

	return b.String()
}

func Split(ref string, delim rune, n int) string {
	if n == 0 {
		return ref
	}

	x := strings.LastIndexFunc(ref,
		func(r rune) bool {
			if r == delim && n > 0 {
				n--

				return n == 0
			}

			return false
		},
	)

	if x == -1 {
		return ""
	}

	return ref[:x]
}
