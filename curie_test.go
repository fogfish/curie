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
	"fmt"
	"testing"

	"github.com/fogfish/curie"
	"github.com/fogfish/it"
)

const (
	// schema:prefix#suffix
	i000 = ""
	i100 = "a:"
	i010 = "b"
	i110 = "a:b"
	i120 = "a:b/c"
	i130 = "a:b/c/d"

	i011 = "b#c"
	i111 = "a:b#c"
	i112 = "a:b#c/d"
	i121 = "a:b/c#d"
	i122 = "a:b/c#d/e"
	i133 = "a:b/c/d#e/f/g"

	w100 = "a+b+c+d:"
	w010 = "b+c+d"
	w110 = "a:b+c+d"
	w120 = "a:b+c+d/e"
)

var IRIs = []string{i000, i100, i010, i110, i120, i130, i011, i111, i112, i121, i122, i133}

func TestNew(t *testing.T) {
	for _, str := range IRIs {
		t.Run(fmt.Sprintf("(%s)", str), func(t *testing.T) {
			iri := curie.New(str)

			it.Ok(t).
				If(iri.String()).Equal(str).
				If(iri.Safe()).Equal("[" + str + "]").
				If(curie.String(str).IRI().String()).Equal(str)
		})

	}
}

func TestNewSafe(t *testing.T) {
	for _, str := range IRIs {
		t.Run(fmt.Sprintf("(%s)", str), func(t *testing.T) {
			iri := curie.New("[" + str + "]")

			it.Ok(t).
				If(iri.String()).Equal(str).
				If(iri.Safe()).Equal("[" + str + "]")
		})

	}
}

func TestNewWrong(t *testing.T) {
	for _, str := range []string{w100, w010, w110, w120} {
		t.Run(fmt.Sprintf("(%s)", str), func(t *testing.T) {
			iri := curie.New(str)

			it.Ok(t).
				If(iri.String()).Equal(str).
				If(iri.Safe()).Equal("[" + str + "]").
				If(curie.String(str).IRI().String()).Equal(str)
		})

	}
}

func TestJSON(t *testing.T) {
	type Struct struct {
		ID curie.IRI `json:"id"`
	}

	for _, str := range IRIs {
		t.Run(fmt.Sprintf("(%s)", str), func(t *testing.T) {
			id := curie.New(str)

			send := Struct{ID: id}
			var recv Struct

			bytes, err1 := json.Marshal(send)
			err2 := json.Unmarshal(bytes, &recv)

			safe := id.Safe()
			if safe == "[]" {
				safe = ""
			}

			it.Ok(t).
				If(err1).Should().Equal(nil).
				If(err2).Should().Equal(nil).
				If(recv).Should().Equal(send).
				If(string(bytes)).Should().Equal("{\"id\":\"" + safe + "\"}")
		})
	}
}

func TestIsEmpty(t *testing.T) {
	for str, val := range map[string]bool{
		i000: true,
		i100: false,
		i010: false,
		i110: false,
		i120: false,
		i130: false,
		i011: false,
		i111: false,
		i112: false,
		i121: false,
		i122: false,
		i133: false,
	} {
		t.Run(fmt.Sprintf("(%s)", str), func(t *testing.T) {
			iri := curie.New(str)

			it.Ok(t).
				If(curie.IsEmpty(iri)).Equal(val)
		})
	}
}

func TestRank(t *testing.T) {
	for str, seq := range map[string]int{
		i000: 0,
		i100: 1,
		i010: 2,
		i110: 2,
		i120: 3,
		i130: 4,
		i011: 3,
		i111: 3,
		i112: 4,
		i121: 4,
		i122: 5,
		i133: 7,
	} {

		t.Run(fmt.Sprintf("(%s)", str), func(t *testing.T) {
			iri := curie.New(str)

			it.Ok(t).
				If(curie.Rank(iri)).Equal(seq)
		})
	}
}

func TestSeq(t *testing.T) {
	for str, seq := range map[string][]string{
		i000: {},
		i100: {"a"},
		i010: {"", "b"},
		i110: {"a", "b"},
		i120: {"a", "b", "c"},
		i130: {"a", "b", "c", "d"},
		i011: {"", "b", "c"},
		i111: {"a", "b", "c"},
		i112: {"a", "b", "c", "d"},
		i121: {"a", "b", "c", "d"},
		i122: {"a", "b", "c", "d", "e"},
		i133: {"a", "b", "c", "d", "e", "f", "g"},
	} {

		t.Run(fmt.Sprintf("(%s)", str), func(t *testing.T) {
			iri := curie.New(str)

			it.Ok(t).
				If(curie.Seq(iri)).Equal(seq)
		})
	}
}

/*
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

func TestReScheme(t *testing.T) {
	tZ := curie.ReScheme(rZ, "t")
	t0 := curie.ReScheme(r0, "t")
	t1 := curie.ReScheme(r1, "t")
	t5 := curie.ReScheme(r5, "t")

	it.Ok(t).
		If(tZ.Safe()).Equal("[t:]").
		If(t0.Safe()).Equal("[t:]").
		If(t1.Safe()).Equal("[t:b]").
		If(t5.Safe()).Equal("[t:b/c/d/e]")
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
		&r0: {rZ, curie.New("a:")},
		&r1: {rZ, curie.New("b")},
		&r5: {rZ, curie.New("a:b/c/d/e")},
	} {
		it.Ok(t).
			If(curie.Heir(v[0], *k)).Should().Equal(v[1])
	}

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
*/
