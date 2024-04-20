//
// Copyright (C) 2020 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/fogfish/curie
//

package reference_test

import (
	"fmt"
	"testing"

	"github.com/fogfish/curie/v2/internal/reference"
	"github.com/fogfish/it/v2"
)

func TestJoin(t *testing.T) {
	for _, ref := range []string{
		"",
		"b",
		"b/c",
		"b/c/d",
	} {
		t.Run(fmt.Sprintf("(%s)", ref), func(t *testing.T) {
			af1 := reference.Join(ref, '/', "x")
			af2 := reference.Join(ref, '/', "x", "y")
			af3 := reference.Join(ref, '/', "x", "y", "z")
			af4 := reference.Join(ref, '/')
			af5 := reference.Join(ref, '/', "")
			af6 := reference.Join(ref, '/', "", "")
			af7 := reference.Join(ref, '/', "", "", "")

			pfx := ref
			if len(ref) != 0 {
				pfx += "/"
			}

			it.Then(t).Should(
				it.Equal(af1, pfx+"x"),
				it.Equal(af2, pfx+"x/y"),
				it.Equal(af3, pfx+"x/y/z"),
				it.Equal(af4, ref),
				it.Equal(af5, ref),
				it.Equal(af6, ref),
				it.Equal(af7, ref),
			)
		})
	}
}

func TestSplit(t *testing.T) {
	it.Then(t).Should(
		it.Equal("",
			reference.Split("", '/', 1),
		),
		it.Equal("",
			reference.Split("", '/', 2),
		),
		it.Equal("",
			reference.Split("b", '/', 1),
		),
		it.Equal("",
			reference.Split("b", '/', 2),
		),
		it.Equal("b",
			reference.Split("b/c", '/', 1),
		),
		it.Equal("b",
			reference.Split("b/c", '/', 2),
		),
		it.Equal("b/c",
			reference.Split("b/c/d", '/', 1),
		),
		it.Equal("b",
			reference.Split("b/c/d", '/', 2),
		),
	)
}
