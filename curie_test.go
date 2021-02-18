//
// Copyright (C) 2020 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/fogfish/curie
//

package curie_test

import (
	"encoding/json"
	"net/url"
	"testing"

	"github.com/fogfish/curie"
	"github.com/fogfish/it"
)

var (
	rZ curie.IRI = curie.New("")
	r0 curie.IRI = curie.New("a:")
	r1 curie.IRI = curie.New("b")
	r2 curie.IRI = curie.New("a:b")
	r3 curie.IRI = curie.New("a:b/c")
	r4 curie.IRI = curie.New("a:b/c/d")
	r5 curie.IRI = curie.New("a:b/c/d/e")

	w0 curie.IRI = curie.New("a+b+c+d:")
	w1 curie.IRI = curie.New("b+c+d")
	w2 curie.IRI = curie.New("a:b+c+d")
	w3 curie.IRI = curie.New("a:b+c+d/e")
)

func TestIRI(t *testing.T) {
	test := map[*curie.IRI][]string{
		&rZ: {},
		&r0: {"a"},
		&r1: {"", "b"},
		&r2: {"a", "b"},
		&r3: {"a", "b", "c"},
		&r4: {"a", "b", "c", "d"},
		&r5: {"a", "b", "c", "d", "e"},

		&w0: {"a+b+c+d"},
		&w1: {"", "b+c+d"},
		&w2: {"a", "b+c+d"},
		&w3: {"a", "b+c+d", "e"},
	}

	for k, v := range test {
		it.Ok(t).
			If(k.Seq()).Should().Equal(v)
	}
}

func TestSafeIRI(t *testing.T) {
	for k, v := range map[string]*curie.IRI{
		"[]":          &rZ,
		"[a:]":        &r0,
		"[b]":         &r1,
		"[a:b]":       &r2,
		"[a:b/c]":     &r3,
		"[a:b/c/d]":   &r4,
		"[a:b/c/d/e]": &r5,
	} {
		it.Ok(t).
			If(curie.New(k)).Should().Equal(*v).
			If(v.Safe()).Should().Equal(k)
	}
}

func TestIdentity(t *testing.T) {
	test := map[*curie.IRI]string{
		&rZ: "",
		&r0: "a:",
		&r1: "b",
		&r2: "a:b",
		&r3: "a:b/c",
		&r4: "a:b/c/d",
		&r5: "a:b/c/d/e",
		&w0: "a+b+c+d:",
		&w1: "b+c+d",
		&w2: "a:b+c+d",
		&w3: "a:b+c+d/e",
	}

	for k, v := range test {
		it.Ok(t).
			If(k.String()).Should().Equal(v)
	}
}

func TestThis(t *testing.T) {
	it.Ok(t).
		IfNotNil(r5.This()).
		IfTrue(r5.This().Eq(r5))
}

func TestOrigin(t *testing.T) {
	test := map[*curie.IRI][]curie.IRI{
		&rZ: {rZ, rZ, rZ, rZ, rZ, rZ},
		&r0: {r0, r0, r0, r0, r0, r0},
		// &r1: {rZ, rZ, r1, r1, r1, r1},
		&r2: {r0, r0, r2, r2, r2, r2},
		&r3: {r0, r0, r2, r3, r3, r3},
		&r4: {r0, r0, r2, r3, r4, r4},
		&r5: {r0, r0, r2, r3, r4, r5},
	}

	for k, v := range test {
		it.Ok(t).
			If(k.Origin()).Should().Equal(v[0]).
			If(k.Origin(1)).Should().Equal(v[1]).
			If(k.Origin(2)).Should().Equal(v[2]).
			If(k.Origin(3)).Should().Equal(v[3]).
			If(k.Origin(4)).Should().Equal(v[4]).
			If(k.Origin(5)).Should().Equal(v[5])
	}
}

func TestOriginNegative(t *testing.T) {
	test := map[*curie.IRI][]curie.IRI{
		&rZ: {rZ, rZ, rZ, rZ, rZ},
		&r0: {rZ, rZ, rZ, rZ, rZ},
		&r2: {r0, rZ, rZ, rZ, rZ},
		&r3: {r2, r0, rZ, rZ, rZ},
		&r4: {r3, r2, r0, rZ, rZ},
		&r5: {r4, r3, r2, r0, rZ},
	}

	for k, v := range test {
		it.Ok(t).
			If(k.Origin(-1)).Should().Equal(v[0]).
			If(k.Origin(-2)).Should().Equal(v[1]).
			If(k.Origin(-3)).Should().Equal(v[2]).
			If(k.Origin(-4)).Should().Equal(v[3]).
			If(k.Origin(-5)).Should().Equal(v[4])
	}
}

func TestPrefix(t *testing.T) {
	test := map[*curie.IRI][]string{
		&rZ: {"", "", "", "", "", ""},
		&r0: {"", "", "", "", "", ""},
		&r1: {"", "", "", "", "", ""},
		&r2: {"a:", "a:", "", "", "", ""},
		&r3: {"a:b", "a:b", "a:", "", "", ""},
		&r4: {"a:b/c", "a:b/c", "a:b", "a:", "", ""},
		&r5: {"a:b/c/d", "a:b/c/d", "a:b/c", "a:b", "a:", ""},
	}

	for k, v := range test {
		it.Ok(t).
			If(k.Prefix()).Should().Equal(v[0]).
			If(k.Prefix(1)).Should().Equal(v[1]).
			If(k.Prefix(2)).Should().Equal(v[2]).
			If(k.Prefix(3)).Should().Equal(v[3]).
			If(k.Prefix(4)).Should().Equal(v[4]).
			If(k.Prefix(5)).Should().Equal(v[5])
	}
}

func TestSuffix(t *testing.T) {
	test := map[*curie.IRI][]string{
		&rZ: {"", "", "", "", "", ""},
		&r0: {"a:", "a:", "a:", "a:", "a:", "a:"},
		&r1: {"b", "b", "b", "b", "b", "b"},
		&r2: {"b", "b", "a:b", "a:b", "a:b", "a:b"},
		&r3: {"c", "c", "b/c", "a:b/c", "a:b/c", "a:b/c"},
		&r4: {"d", "d", "c/d", "b/c/d", "a:b/c/d", "a:b/c/d"},
		&r5: {"e", "e", "d/e", "c/d/e", "b/c/d/e", "a:b/c/d/e"},
	}

	for k, v := range test {
		it.Ok(t).
			If(k.Suffix()).Should().Equal(v[0]).
			If(k.Suffix(1)).Should().Equal(v[1]).
			If(k.Suffix(2)).Should().Equal(v[2]).
			If(k.Suffix(3)).Should().Equal(v[3]).
			If(k.Suffix(4)).Should().Equal(v[4]).
			If(k.Suffix(5)).Should().Equal(v[5])
	}
}

func TestSuffixNegative(t *testing.T) {
	test := map[*curie.IRI][]string{
		&rZ: {"", "", "", "", ""},
		&r0: {"", "", "", "", ""},
		&r1: {"b", "", "", "", ""},
		&r2: {"b", "", "", "", ""},
		&r3: {"b/c", "c", "", "", ""},
		&r4: {"b/c/d", "c/d", "d", "", ""},
		&r5: {"b/c/d/e", "c/d/e", "d/e", "e", ""},
	}

	for k, v := range test {
		it.Ok(t).
			If(k.Suffix(-1)).Should().Equal(v[0]).
			If(k.Suffix(-2)).Should().Equal(v[1]).
			If(k.Suffix(-3)).Should().Equal(v[2]).
			If(k.Suffix(-4)).Should().Equal(v[3]).
			If(k.Suffix(-5)).Should().Equal(v[4]).
			If(k.Suffix(-6)).Should().Equal(v[4])
	}
}

func TestParent(t *testing.T) {
	test := map[*curie.IRI][]curie.IRI{
		&rZ: {rZ, rZ, rZ, rZ, rZ, rZ},
		&r0: {rZ, rZ, rZ, rZ, rZ, rZ},
		&r1: {rZ, rZ, rZ, rZ, rZ, rZ},
		&r2: {r0, r0, rZ, rZ, rZ, rZ},
		&r3: {r2, r2, r0, rZ, rZ, rZ},
		&r4: {r3, r3, r2, r0, rZ, rZ},
		&r5: {r4, r4, r3, r2, r0, rZ},
	}

	for k, v := range test {
		it.Ok(t).
			If(k.Parent()).Should().Equal(v[0]).
			If(k.Parent(1)).Should().Equal(v[1]).
			If(k.Parent(2)).Should().Equal(v[2]).
			If(k.Parent(3)).Should().Equal(v[3]).
			If(k.Parent(4)).Should().Equal(v[4]).
			If(k.Parent(5)).Should().Equal(v[5])
	}
}

func TestParentNegative(t *testing.T) {
	test := map[*curie.IRI][]curie.IRI{
		&rZ: {rZ, rZ, rZ, rZ, rZ},
		&r0: {r0, r0, r0, r0, r0},
		&r1: {rZ, r1, r1, r1, r1},
		&r2: {r0, r2, r2, r2, r2},
		&r3: {r0, r2, r3, r3, r3},
		&r4: {r0, r2, r3, r4, r4},
		&r5: {r0, r2, r3, r4, r5},
	}

	for k, v := range test {
		it.Ok(t).
			If(k.Parent(-1)).Should().Equal(v[0]).
			If(k.Parent(-2)).Should().Equal(v[1]).
			If(k.Parent(-3)).Should().Equal(v[2]).
			If(k.Parent(-4)).Should().Equal(v[3]).
			If(k.Parent(-5)).Should().Equal(v[4]).
			If(k.Parent(-6)).Should().Equal(v[4])

	}
}

func TestJoin(t *testing.T) {
	it.Ok(t).
		If(rZ.Join("a")).Should().Equal(r0).
		If(r0.Join("b")).Should().Equal(r2).
		If(r2.Join("c")).Should().Equal(r3).
		If(r3.Join("d")).Should().Equal(r4).
		If(r4.Join("e")).Should().Equal(r5)
}

func TestJoinRanked(t *testing.T) {
	it.Ok(t).
		If(rZ.Join("a/b/c/d/e")).Should().Equal(r5).
		If(r0.Join("b/c/d/e")).Should().Equal(r5).
		If(r2.Join("c/d/e")).Should().Equal(r5).
		If(r3.Join("d/e")).Should().Equal(r5).
		If(r4.Join("e")).Should().Equal(r5).
		If(rZ.Join("a:b/c/d/e")).Should().Equal(r5).
		If(r0.Join("b:c/d/e")).Should().Equal(r5).
		If(r2.Join("c:d/e")).Should().Equal(r5).
		If(r3.Join("d:e")).Should().Equal(r5).
		If(r4.Join("e:")).Should().Equal(r5)
}

func TestJoinImmutable(t *testing.T) {
	rN := r3.Parent().Join("t")

	it.Ok(t).
		If(r3.Path()).Should().Equal("a/b/c").
		If(rN.Path()).Should().Equal("a/b/t")
}

func TestHeir(t *testing.T) {
	for k, v := range map[*curie.IRI][]curie.IRI{
		&rZ: {r5, curie.New("a:b/c/d/e")},
		&r0: {r5, curie.New("a:a/b/c/d/e")},
		&r1: {r5, curie.New("b/a/b/c/d/e")},
		&r2: {r5, curie.New("a:b/a/b/c/d/e")},
		&r3: {r5, curie.New("a:b/c/a/b/c/d/e")},
		&r4: {r5, curie.New("a:b/c/d/a/b/c/d/e")},
		&r5: {r5, curie.New("a:b/c/d/e/a/b/c/d/e")},
	} {
		it.Ok(t).
			If(k.Heir(v[0])).Should().Equal(v[1])
	}
}

func TestHeirImmutable(t *testing.T) {
	rN := r3.Parent().Heir(curie.New("t"))

	it.Ok(t).
		If(r3.Path()).Should().Equal("a/b/c").
		If(rN.Path()).Should().Equal("a/b/t")
}

func TestCURIE2URI(t *testing.T) {
	for compact, v := range map[string][]string{
		"":          {"https://example.com/", ""},
		"a:":        {"https://example.com/", "https://example.com/"},
		"a:b":       {"https://example.com/", "https://example.com/b"},
		"b":         {"https://example.com/", "b"},
		"b/c/d/e":   {"https://example.com/", "b/c/d/e"},
		"a:b/c/d/e": {"https://example.com#", "https://example.com#b/c/d/e"},
	} {
		curi := curie.New(compact)
		expect, _ := url.Parse(v[1])
		uri, err := curi.URI(v[0])

		it.Ok(t).
			If(err).Should().Equal(nil).
			If(uri).Should().Equal(expect)
	}
}

func TestPath(t *testing.T) {
	test := map[*curie.IRI]string{
		&rZ: "",
		&r0: "a",
		&r1: "b",
		&r2: "a/b",
		&r3: "a/b/c",
		&r4: "a/b/c/d",
		&r5: "a/b/c/d/e",
	}

	for k, v := range test {
		it.Ok(t).
			If(k.Path()).Should().Equal(v)
	}
}

func TestEq(t *testing.T) {
	test := []curie.IRI{r0, r1, r2, r3, r4, r5}

	for _, v := range test {
		it.Ok(t).
			If(v.Eq(v)).Should().Equal(true).
			If(v.Eq(w0)).Should().Equal(false)
	}
}

func TestNotEq(t *testing.T) {
	r6 := curie.New("1:2:3:4:5:6")
	test := []curie.IRI{r0, r1, r2, r3, r4, r5}

	for _, v := range test {
		it.Ok(t).
			If(v.Eq(r6)).Should().Equal(false).
			If(v.Eq(v.Parent().Join("t"))).Should().Equal(false)
	}
}

func TestLt(t *testing.T) {
	for a, b := range map[string]string{
		"":        "b:",
		"a:":      "b:",
		"a":       "b",
		"a/a":     "a/b",
		"a:a":     "a:b",
		"a:x/a":   "a:x/b",
		"a:x/x/a": "a:x/x/x/a",
	} {
		it.Ok(t).
			If(curie.New(a).Lt(curie.New(b))).Should().Equal(true)
	}
}

func TestJSON(t *testing.T) {
	type Struct struct {
		ID    curie.IRI `json:"id"`
		Title string    `json:"title"`
	}

	test := map[*Struct]string{
		{ID: rZ, Title: "t"}: "{\"id\":\"\",\"title\":\"t\"}",
		{ID: r0, Title: "t"}: "{\"id\":\"[a:]\",\"title\":\"t\"}",
		{ID: r1, Title: "t"}: "{\"id\":\"[b]\",\"title\":\"t\"}",
		{ID: r2, Title: "t"}: "{\"id\":\"[a:b]\",\"title\":\"t\"}",
		{ID: r3, Title: "t"}: "{\"id\":\"[a:b/c]\",\"title\":\"t\"}",
	}

	for eg, expect := range test {
		in := Struct{}

		bytes, err1 := json.Marshal(eg)
		err2 := json.Unmarshal(bytes, &in)

		it.Ok(t).
			If(err1).Should().Equal(nil).
			If(err2).Should().Equal(nil).
			If(*eg).Should().Equal(in).
			If(string(bytes)).Should().Equal(expect)
	}
}

func TestTypeSafe(t *testing.T) {
	type A struct{ curie.IRI }
	type B struct{ curie.IRI }
	type C struct{ curie.IRI }

	a := A{curie.New("a:")}
	b := B{curie.New("a:b")}
	c := C{b.IRI.Join("c")}

	it.Ok(t).
		If(a.IRI).Should().Equal(r0).
		If(b.IRI).Should().Equal(r2).
		If(c.IRI).Should().Equal(r3)
}
