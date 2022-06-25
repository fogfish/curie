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
	"net/url"
	"testing"

	"github.com/fogfish/curie"
	"github.com/fogfish/it"
)

func TestNew(t *testing.T) {
	for _, str := range []string{
		"a:",
		"a:b",
		"a:b/c",
		"a:b/c/d",
		"b",
		"b/c",
		"b/c/d",
	} {
		t.Run(fmt.Sprintf("(%s)", str), func(t *testing.T) {
			iri := curie.New(str)

			it.Ok(t).
				If(string(iri)).Equal(str).
				If(iri.Safe()).Equal("[" + str + "]")
			// If(curie.String(str).IRI().String()).Equal(str).
			// If(*curie.Safe(iri)).Equal(curie.String("[" + str + "]"))
		})
	}
}

func TestNewSafe(t *testing.T) {
	for _, str := range []string{
		"a:",
		"a:b",
		"a:b/c",
		"a:b/c/d",
		"b",
		"b/c",
		"b/c/d",
	} {
		t.Run(fmt.Sprintf("(%s)", str), func(t *testing.T) {
			iri := curie.New("[" + str + "]")

			it.Ok(t).
				If(string(iri)).Equal(str).
				If(iri.Safe()).Equal("[" + str + "]")
			// If(curie.String(str).IRI().String()).Equal(str).
			// If(*curie.Safe(iri)).Equal(curie.String("[" + str + "]"))
		})
	}
}

func TestNewFormat(t *testing.T) {
	iri := curie.New("a:b/c/d/%s", "e")

	it.Ok(t).
		If(string(iri)).Equal("a:b/c/d/e")
}

func TestNewBadFormat(t *testing.T) {
	for _, str := range []string{
		"a+b+c+d:",
		"b+c+d",
		"a:b+c+d",
		"a:b+c+d/e",
	} {
		t.Run(fmt.Sprintf("(%s)", str), func(t *testing.T) {
			iri := curie.New(str)

			it.Ok(t).
				If(string(iri)).Equal(str).
				If(iri.Safe()).Equal("[" + str + "]")
			// If(curie.String(str).IRI().String()).Equal(str)
		})
	}
}

func TestJSON(t *testing.T) {
	type Struct struct {
		ID curie.IRI `json:"id"`
	}

	for _, str := range []string{
		"",
		"a:",
		"a:b",
		"a:b/c",
		"a:b/c/d",
		"b",
		"b/c",
		"b/c/d",
	} {
		t.Run(fmt.Sprintf("(%s)", str), func(t *testing.T) {
			id := curie.New(str)

			send := Struct{ID: id}
			var recv Struct

			bytes, err1 := json.Marshal(send)
			err2 := json.Unmarshal(bytes, &recv)

			safe := id.Safe()

			it.Ok(t).
				If(err1).Should().Equal(nil).
				If(err2).Should().Equal(nil).
				If(recv).Should().Equal(send).
				If(string(bytes)).Should().Equal("{\"id\":\"" + safe + "\"}")
		})
	}
}

func TestJSONError(t *testing.T) {
	type Struct struct {
		ID curie.IRI `json:"id"`
	}
	var recv Struct

	err := json.Unmarshal([]byte("{\"id\":100}"), &recv)

	it.Ok(t).
		IfNotNil(err)
}

func TestIsEmpty(t *testing.T) {
	for str, val := range map[string]bool{
		"":    true,
		"a:":  false,
		"a:b": false,
		"b":   false,
		"b/c": false,
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
		"":      curie.EMPTY,
		"a:":    curie.PREFIX,
		"a:b":   curie.REFERENCE,
		"a:b/c": curie.REFERENCE,
		"b":     curie.REFERENCE,
		"b/c":   curie.REFERENCE,
		"b/c/d": curie.REFERENCE,
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
		"":      {},
		"a:":    {"a"},
		"a:b":   {"a", "b"},
		"a:b/c": {"a", "b/c"},
		"b":     {"", "b"},
		"b/c":   {"", "b/c"},
		"b/c/d": {"", "b/c/d"},
	} {
		t.Run(fmt.Sprintf("(%s)", str), func(t *testing.T) {
			iri := curie.New(str)

			it.Ok(t).
				If(curie.Seq(iri)).Equal(seq)
		})
	}
}

// func TestPath(t *testing.T) {
// 	for str, seq := range map[string]string{
// 		"":      "",
// 		"a:":    "a",
// 		"a:b":   "a/b",
// 		"a:b/c": "a/b/c",
// 		"b":     "b",
// 		"b/c":   "b/c",
// 		"b/c/d": "b/c/d",
// 	} {
// 		t.Run(fmt.Sprintf("(%s)", str), func(t *testing.T) {
// 			iri := curie.New(str)

// 			it.Ok(t).
// 				If(curie.Path(iri)).Equal(seq)
// 		})
// 	}
// }

func TestURL(t *testing.T) {
	compact := curie.New("wikipedia:CURIE")
	url, err := curie.URL(curie.Namespaces{
		"wikipedia": "http://en.wikipedia.org/wiki/",
	}, compact)

	it.Ok(t).
		IfNil(err).
		If(url.String()).Equal("http://en.wikipedia.org/wiki/CURIE")
}

func TestURLCompatibility(t *testing.T) {
	uri := "https://example.com/a/b/c?de=fg&foo=bar"
	curi := curie.New(uri)

	expect, _ := url.Parse(uri)
	native, err := curie.URL(curie.Namespaces{}, curi)

	it.Ok(t).
		If(string(curi)).Equal(uri).
		If(curi.Safe()).Equal("[" + uri + "]").
		If(curie.Seq(curi)).Equal([]string{"https", "//example.com/a/b/c?de=fg&foo=bar"}).
		//
		IfNil(err).
		If(native).Equal(expect)
}

func TestURLConvert(t *testing.T) {
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
		uri, err := curie.URL(curie.Namespaces{"a": v[0]}, curi)

		it.Ok(t).
			If(err).Should().Equal(nil).
			If(uri).Should().Equal(expect)
	}
}

func TestEq(t *testing.T) {
	for _, str := range []string{
		"",
		"a:",
		"a:b",
		"a:b/c",
		"a:b/c/d",
		"b",
		"b/c",
		"b/c/d",
	} {
		t.Run(fmt.Sprintf("(%s)", str), func(t *testing.T) {
			a := curie.New(str)
			b := curie.New(str)
			c := curie.New("z:x/y")

			it.Ok(t).
				If(curie.Eq(a, b)).Should().Equal(true).
				If(curie.Eq(a, c)).Should().Equal(false)
		})
	}
}

func TestNotEq(t *testing.T) {
	r6 := curie.New("1:2:3:4:5:6")

	for _, str := range []string{
		"",
		"a:",
		"a:b",
		"a:b/c",
		"a:b/c/d",
		"b",
		"b/c",
		"b/c/d",
	} {
		t.Run(fmt.Sprintf("(%s)", str), func(t *testing.T) {
			a := curie.New(str)

			it.Ok(t).
				If(curie.Eq(a, r6)).Should().Equal(false)
		})
	}
}

func TestLt(t *testing.T) {
	for a, b := range map[string]string{
		"a:":    "b:",
		"a:b":   "a:c",
		"a:b/c": "a:b/d",
		"b":     "c",
		"b/c":   "b/d",
		"b/c/d": "b/c/e",
		"z:":    "z:x",
	} {
		t.Run(fmt.Sprintf("(%s)", a), func(t *testing.T) {
			it.Ok(t).
				If(curie.Lt(curie.New(a), curie.New(b))).Should().Equal(true).
				If(curie.Lt(curie.New(b), curie.New(a))).Should().Equal(false).
				If(curie.Lt(curie.New(a), curie.New(a))).Should().Equal(false)
		})
	}
}

func TestPrefixAndReference(t *testing.T) {
	for str, val := range map[string][]string{
		"":      {"", ""},
		"a:":    {"a", ""},
		"a:b":   {"a", "b"},
		"a:b/c": {"a", "b/c"},
		"b":     {"", "b"},
		"b/c":   {"", "b/c"},
	} {
		t.Run(fmt.Sprintf("(%s)", str), func(t *testing.T) {
			iri := curie.New(str)

			it.Ok(t).
				If(curie.Prefix(iri)).Should().Equal(val[0]).
				If(curie.Reference(iri)).Should().Equal(val[1])
		})
	}
}

func TestJoin(t *testing.T) {
	for _, str := range []string{
		"a:b",
		"a:b/c",
		"b",
		"b/c",
		"b/c/d",
	} {
		t.Run(fmt.Sprintf("(%s)", str), func(t *testing.T) {
			iri := curie.New(str)
			af1 := curie.Join(iri, "x")
			af2 := curie.Join(iri, "x", "y")
			af3 := curie.Join(iri, "x", "y", "z")
			af4 := curie.Join(iri)
			af5 := curie.Join(iri, "")
			af6 := curie.Join(iri, "", "")
			af7 := curie.Join(iri, "", "", "")

			it.Ok(t).
				If(string(af1)).Equal(str + "/x").
				If(string(af2)).Equal(str + "/x/y").
				If(string(af3)).Equal(str + "/x/y/z").
				If(string(af4)).Equal(str).
				If(string(af5)).Equal(str).
				If(string(af6)).Equal(str).
				If(string(af7)).Equal(str)
		})
	}
}

func TestJoinWithZero(t *testing.T) {
	it.Ok(t).
		If(string(curie.Join(curie.New(""), "x"))).Equal("x").
		If(string(curie.Join(curie.New(""), "x", "y"))).Equal("x/y").
		If(string(curie.Join(curie.New(""), "x", "y", "z"))).Equal("x/y/z").
		If(string(curie.Join(curie.New("a:"), "x"))).Equal("a:x").
		If(string(curie.Join(curie.New("a:"), "x", "y"))).Equal("a:x/y").
		If(string(curie.Join(curie.New("a:"), "x", "y", "z"))).Equal("a:x/y/z").
		If(string(curie.Join(curie.New("b"), "x"))).Equal("b/x").
		If(string(curie.Join(curie.New("b"), "x", "y"))).Equal("b/x/y").
		If(string(curie.Join(curie.New("b"), "x", "y", "z"))).Equal("b/x/y/z")
}

/*
func TestJoinImmutable(t *testing.T) {
	r3 := curie.New("a:b/c/d")
	rP := curie.Parent(r3)
	rN := curie.Join(rP, "t")

	it.Ok(t).
		If(r3.String()).Should().Equal("a:b/c/d").
		If(rP.String()).Should().Equal("a:b/c").
		If(rN.String()).Should().Equal("a:b/c/t")
}
*/

/*
func TestParent(t *testing.T) {
	it.Ok(t).
		If(curie.Parent(curie.New("a:b/c/d")).String()).Equal("a:b/c").
		If(curie.Parent(curie.New("a:b/c/d"), 1).String()).Equal("a:b/c").
		If(curie.Parent(curie.New("a:b/c/d"), 2).String()).Equal("a:b").
		If(curie.Parent(curie.New("a:b/c/d"), 3).String()).Equal("a:").
		If(curie.Parent(curie.New("a:b/c/d"), 4).String()).Equal("").
		If(curie.Parent(curie.New("a:b/c/d"), 5).String()).Equal("").
		If(curie.Parent(curie.New("a:b/c/d"), -1).String()).Equal("a:").
		If(curie.Parent(curie.New("a:b/c/d"), -2).String()).Equal("a:b").
		If(curie.Parent(curie.New("a:b/c/d"), -3).String()).Equal("a:b/c").
		If(curie.Parent(curie.New("a:b/c/d"), -4).String()).Equal("a:b/c/d").
		If(curie.Parent(curie.New("a:b/c/d"), -5).String()).Equal("a:b/c/d").
		If(curie.Parent(curie.New("b"), 1).String()).Equal("")
}
*/

/*
func TestChild(t *testing.T) {
	it.Ok(t).
		If(curie.Child(curie.New("a:b/c/d"))).Equal("d").
		If(curie.Child(curie.New("a:b/c/d"), 1)).Equal("d").
		If(curie.Child(curie.New("a:b/c/d"), 2)).Equal("c/d").
		If(curie.Child(curie.New("a:b/c/d"), 3)).Equal("b/c/d").
		If(curie.Child(curie.New("a:b/c/d"), 4)).Equal("a:b/c/d").
		If(curie.Child(curie.New("a:b/c/d"), -1)).Equal("b/c/d").
		If(curie.Child(curie.New("a:b/c/d"), -2)).Equal("c/d").
		If(curie.Child(curie.New("a:b/c/d"), -3)).Equal("d").
		If(curie.Child(curie.New("a:b/c/d"), -4)).Equal("")
}
*/

func TestTypeSafe(t *testing.T) {
	type A curie.IRI
	type B curie.IRI
	type C curie.IRI

	a := A(curie.New("a:"))
	b := B(curie.New("a:b"))
	c := C(curie.Join(curie.IRI(b), "c"))

	it.Ok(t).
		If(string(curie.IRI(a))).Should().Equal("a:").
		If(string(curie.IRI(b))).Should().Equal("a:b").
		If(string(curie.IRI(c))).Should().Equal("a:b/c")
}

func TestLinkedData(t *testing.T) {
	type Struct struct {
		ID curie.IRI  `json:"id"`
		LA *curie.IRI `json:"a,omitempty"`
		LB *curie.IRI `json:"b,omitempty"`
	}

	b := curie.New("a:b/c")
	test := map[*Struct]string{
		{ID: curie.New(""), LA: &b, LB: &b}:      "{\"id\":\"\",\"a\":\"[a:b/c]\",\"b\":\"[a:b/c]\"}",
		{ID: curie.New("a:"), LA: &b, LB: &b}:    "{\"id\":\"[a:]\",\"a\":\"[a:b/c]\",\"b\":\"[a:b/c]\"}",
		{ID: curie.New("b"), LA: &b, LB: &b}:     "{\"id\":\"[b]\",\"a\":\"[a:b/c]\",\"b\":\"[a:b/c]\"}",
		{ID: curie.New("a:b"), LA: &b, LB: &b}:   "{\"id\":\"[a:b]\",\"a\":\"[a:b/c]\",\"b\":\"[a:b/c]\"}",
		{ID: curie.New("a:b/c"), LA: &b, LB: &b}: "{\"id\":\"[a:b/c]\",\"a\":\"[a:b/c]\",\"b\":\"[a:b/c]\"}",
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
