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
			If(curie.Seq(*k)).Should().Equal(v)
	}
}

func TestFormatIRI(t *testing.T) {
	it.Ok(t).
		If(curie.New("a:b/c/d/%s", "e").String()).Equal("a:b/c/d/e")
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
			If(v.Safe()).Should().Equal(k).
			If(curie.Safe(*v)).Should().Equal((*curie.String)(&k)).
			If(curie.String(k).IRI()).Should().Equal(*v)
	}
}

func TestURI(t *testing.T) {
	uri := "https://example.com/a/b/c?de=fg&foo=bar"
	curi := curie.New(uri)

	expect, _ := url.Parse(uri)
	native, err := curie.URI("https:", curi)

	it.Ok(t).
		If(curi.String()).Equal(uri).
		If(curi.Safe()).Equal("[" + uri + "]").
		If(curie.Seq(curi)).Equal([]string{"https", "", "", "example.com", "a", "b", "c?de=fg&foo=bar"}).
		//
		IfNil(err).
		If(native).Equal(expect)
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

func TestScheme(t *testing.T) {
	test := map[*curie.IRI]string{
		&rZ: "",
		&r0: "a",
		&r1: "",
		&r2: "a",
		&r3: "a",
		&r4: "a",
		&r5: "a",

		&w0: "a+b+c+d",
		&w1: "",
		&w2: "a",
		&w3: "a",
	}

	for k, v := range test {
		it.Ok(t).
			If(curie.Scheme(*k)).Should().Equal(v)
	}
}

func TestPrefix(t *testing.T) {
	test := map[*curie.IRI]string{
		&rZ: "",
		&r0: "a:",
		&r1: "b",
		&r2: "a:b",
		&r3: "a:b",
		&r4: "a:b/c",
		&r5: "a:b/c/d",
	}

	for k, v := range test {
		it.Ok(t).
			If(curie.Prefix(*k)).Should().Equal(v)
	}
}

func TestSplitAndPrefix(t *testing.T) {
	test := map[*curie.IRI][]string{
		&rZ: {"", "", "", "", "", ""},
		&r0: {"a:", "a:", "a:", "a:", "a:", "a:"},
		&r1: {"b", "b", "b", "b", "b", "b"},
		&r2: {"a:b", "a:b", "a:b", "a:b", "a:b", "a:b"},
		&r3: {"a:b/c", "a:b", "a:b", "a:b", "a:b", "a:b"},
		&r4: {"a:b/c/d", "a:b/c", "a:b", "a:b", "a:b", "a:b"},
		&r5: {"a:b/c/d/e", "a:b/c/d", "a:b/c", "a:b", "a:b", "a:b"},
	}

	for k, v := range test {
		it.Ok(t).
			If(curie.Prefix(curie.Split(*k, 0))).Should().Equal(v[0]).
			If(curie.Prefix(curie.Split(*k, 1))).Should().Equal(v[1]).
			If(curie.Prefix(curie.Split(*k, 2))).Should().Equal(v[2]).
			If(curie.Prefix(curie.Split(*k, 3))).Should().Equal(v[3]).
			If(curie.Prefix(curie.Split(*k, 4))).Should().Equal(v[4]).
			If(curie.Prefix(curie.Split(*k, 5))).Should().Equal(v[5])
	}
}

func TestSuffix(t *testing.T) {
	test := map[*curie.IRI]string{
		&rZ: "",
		&r0: "",
		&r1: "",
		&r2: "",
		&r3: "c",
		&r4: "d",
		&r5: "e",
	}

	for k, v := range test {
		it.Ok(t).
			If(curie.Suffix(*k)).Should().Equal(v)
	}
}

func TestSplitAndSuffix(t *testing.T) {
	test := map[*curie.IRI][]string{
		&rZ: {"", "", "", "", "", ""},
		&r0: {"", "", "", "", "", ""},
		&r1: {"", "", "", "", "", ""},
		&r2: {"", "", "", "", "", ""},
		&r3: {"", "c", "c", "c", "c", "c"},
		&r4: {"", "d", "c/d", "c/d", "c/d", "c/d"},
		&r5: {"", "e", "d/e", "c/d/e", "c/d/e", "c/d/e"},
	}

	for k, v := range test {
		it.Ok(t).
			If(curie.Suffix(curie.Split(*k, 0))).Should().Equal(v[0]).
			If(curie.Suffix(curie.Split(*k, 1))).Should().Equal(v[1]).
			If(curie.Suffix(curie.Split(*k, 2))).Should().Equal(v[2]).
			If(curie.Suffix(curie.Split(*k, 3))).Should().Equal(v[3]).
			If(curie.Suffix(curie.Split(*k, 4))).Should().Equal(v[4]).
			If(curie.Suffix(curie.Split(*k, 5))).Should().Equal(v[5])
	}
}

func TestParent(t *testing.T) {
	test := map[*curie.IRI]curie.IRI{
		&rZ: rZ,
		&r0: r0,
		&r1: r1,
		&r2: r2,
		&r3: r2,
		&r4: r3,
		&r5: r4,
	}

	for k, v := range test {
		it.Ok(t).
			If(curie.Parent(*k)).Should().Equal(v)
	}
}

func TestJoin(t *testing.T) {
	it.Ok(t).
		If(curie.Join(rZ, "a")).Should().Equal(r0).
		If(curie.Join(r0, "b")).Should().Equal(r2).
		If(curie.Join(r2, "c")).Should().Equal(r3).
		If(curie.Join(r3, "d")).Should().Equal(r4).
		If(curie.Join(r4, "e")).Should().Equal(r5)
}

func TestJoinRanked(t *testing.T) {
	it.Ok(t).
		If(curie.Join(rZ, "a/b/c/d/e")).Should().Equal(r5).
		If(curie.Join(r0, "b/c/d/e")).Should().Equal(r5).
		If(curie.Join(r2, "c/d/e")).Should().Equal(r5).
		If(curie.Join(r3, "d/e")).Should().Equal(r5).
		If(curie.Join(r4, "e")).Should().Equal(r5).
		If(curie.Join(rZ, "a:b/c/d/e")).Should().Equal(r5).
		If(curie.Join(r0, "b:c/d/e")).Should().Equal(r5).
		If(curie.Join(r2, "c:d/e")).Should().Equal(r5).
		If(curie.Join(r3, "d:e")).Should().Equal(r5).
		If(curie.Join(r4, "e:")).Should().Equal(r5)
}

func TestJoinImmutable(t *testing.T) {
	rN := curie.Join(curie.Parent(r3), "t")

	it.Ok(t).
		If(curie.Path(r3)).Should().Equal("a/b/c").
		If(curie.Path(rN)).Should().Equal("a/b/t")
}

func TestHeir(t *testing.T) {
	for k, v := range map[*curie.IRI][]curie.IRI{
		&rZ: {rZ, curie.New("")},
		&rZ: {r0, curie.New("a:")},
		&rZ: {r1, curie.New("b")},
		&rZ: {r5, curie.New("a:b/c/d/e")},
		&r0: {r5, curie.New("a:a/b/c/d/e")},
		&r1: {r5, curie.New("b/a/b/c/d/e")},
		&r2: {r5, curie.New("a:b/a/b/c/d/e")},
		&r3: {r5, curie.New("a:b/c/a/b/c/d/e")},
		&r4: {r5, curie.New("a:b/c/d/a/b/c/d/e")},
		&r5: {r5, curie.New("a:b/c/d/e/a/b/c/d/e")},
	} {
		it.Ok(t).
			If(curie.Heir(*k, v[0])).Should().Equal(v[1])
	}
}

func TestHeirImmutable(t *testing.T) {
	rN := curie.Heir(curie.Parent(r3), curie.New("t"))

	it.Ok(t).
		If(curie.Path(r3)).Should().Equal("a/b/c").
		If(curie.Path(rN)).Should().Equal("a/b/t")
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
		uri, err := curie.URI(v[0], curi)

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
			If(curie.Path(*k)).Should().Equal(v)
	}
}

func TestEq(t *testing.T) {
	test := []curie.IRI{r0, r1, r2, r3, r4, r5}

	for _, v := range test {
		it.Ok(t).
			If(curie.Eq(v, v)).Should().Equal(true).
			If(curie.Eq(v, w0)).Should().Equal(false)
	}
}

func TestNotEq(t *testing.T) {
	r6 := curie.New("1:2:3:4:5:6")
	test := []curie.IRI{r0, r1, r2, r3, r4, r5}

	for _, v := range test {
		it.Ok(t).
			If(curie.Eq(v, r6)).Should().Equal(false).
			If(curie.Eq(v, curie.Join(curie.Parent(v), "t"))).Should().Equal(false)
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
			If(curie.Lt(curie.New(a), curie.New(b))).Should().Equal(true).
			If(curie.Lt(curie.New(b), curie.New(a))).Should().Equal(false)
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
	type A curie.IRI
	type B curie.IRI
	type C curie.IRI

	a := A(curie.New("a:"))
	b := B(curie.New("a:b"))
	c := C(curie.Join(curie.IRI(b), "c"))

	it.Ok(t).
		If(curie.IRI(a)).Should().Equal(r0).
		If(curie.IRI(b)).Should().Equal(r2).
		If(curie.IRI(c)).Should().Equal(r3)
}

func TestLinkedData(t *testing.T) {
	type Struct struct {
		ID curie.IRI     `json:"id"`
		LA *curie.String `json:"a,omitempty"`
		LB *curie.IRI    `json:"b,omitempty"`
	}

	test := map[*Struct]string{
		{ID: rZ, LA: curie.Safe(r3), LB: &r3}: "{\"id\":\"\",\"a\":\"[a:b/c]\",\"b\":\"[a:b/c]\"}",
		{ID: r0, LA: curie.Safe(r3), LB: &r3}: "{\"id\":\"[a:]\",\"a\":\"[a:b/c]\",\"b\":\"[a:b/c]\"}",
		{ID: r1, LA: curie.Safe(r3), LB: &r3}: "{\"id\":\"[b]\",\"a\":\"[a:b/c]\",\"b\":\"[a:b/c]\"}",
		{ID: r2, LA: curie.Safe(r3), LB: &r3}: "{\"id\":\"[a:b]\",\"a\":\"[a:b/c]\",\"b\":\"[a:b/c]\"}",
		{ID: r3, LA: curie.Safe(r3), LB: &r3}: "{\"id\":\"[a:b/c]\",\"a\":\"[a:b/c]\",\"b\":\"[a:b/c]\"}",
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
