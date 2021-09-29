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
				If(iri.String()).Equal(str).
				If(iri.Safe()).Equal("[" + str + "]").
				If(curie.String(str).IRI().String()).Equal(str).
				If(*curie.Safe(iri)).Equal(curie.String("[" + str + "]"))
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
				If(iri.String()).Equal(str).
				If(iri.Safe()).Equal("[" + str + "]").
				If(curie.String(str).IRI().String()).Equal(str).
				If(*curie.Safe(iri)).Equal(curie.String("[" + str + "]"))
		})
	}
}

func TestNewFormat(t *testing.T) {
	it.Ok(t).
		If(curie.New("a:b/c/d/%s", "e").String()).Equal("a:b/c/d/e")
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
		"":      0,
		"a:":    1,
		"a:b":   2,
		"a:b/c": 3,
		"b":     2,
		"b/c":   3,
		"b/c/d": 4,
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
		"a:b/c": {"a", "b", "c"},
		"b":     {"", "b"},
		"b/c":   {"", "b", "c"},
		"b/c/d": {"", "b", "c", "d"},
	} {
		t.Run(fmt.Sprintf("(%s)", str), func(t *testing.T) {
			iri := curie.New(str)

			it.Ok(t).
				If(curie.Seq(iri)).Equal(seq)
		})
	}
}

func TestPath(t *testing.T) {
	for str, seq := range map[string]string{
		"":      "",
		"a:":    "a",
		"a:b":   "a/b",
		"a:b/c": "a/b/c",
		"b":     "b",
		"b/c":   "b/c",
		"b/c/d": "b/c/d",
	} {
		t.Run(fmt.Sprintf("(%s)", str), func(t *testing.T) {
			iri := curie.New(str)

			it.Ok(t).
				If(curie.Path(iri)).Equal(seq)
		})
	}
}

func TestURL(t *testing.T) {
	compact := curie.New("wikipedia:CURIE")
	url, err := curie.URL(map[string]string{
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
	native, err := curie.URL(map[string]string{}, curi)

	it.Ok(t).
		If(curi.String()).Equal(uri).
		If(curi.Safe()).Equal("[" + uri + "]").
		If(curie.Seq(curi)).Equal([]string{"https", "", "", "example.com", "a", "b", "c?de=fg&foo=bar"}).
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
		uri, err := curie.URL(map[string]string{"a": v[0]}, curi)

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

func TestPrefixAndSuffix(t *testing.T) {
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
				If(curie.Suffix(iri)).Should().Equal(val[1])
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

			it.Ok(t).
				If(af1.String()).Equal(str + "/x").
				If(af2.String()).Equal(str + "/x/y").
				If(af3.String()).Equal(str + "/x/y/z")
		})
	}
}

func TestJoinWithZero(t *testing.T) {
	it.Ok(t).
		If(curie.Join(curie.New(""), "x").String()).Equal("x").
		If(curie.Join(curie.New(""), "x", "y").String()).Equal("x/y").
		If(curie.Join(curie.New(""), "x", "y", "z").String()).Equal("x/y/z").
		If(curie.Join(curie.New("a:"), "x").String()).Equal("a:x").
		If(curie.Join(curie.New("a:"), "x", "y").String()).Equal("a:x/y").
		If(curie.Join(curie.New("a:"), "x", "y", "z").String()).Equal("a:x/y/z").
		If(curie.Join(curie.New("b"), "x").String()).Equal("b/x").
		If(curie.Join(curie.New("b"), "x", "y").String()).Equal("b/x/y").
		If(curie.Join(curie.New("b"), "x", "y", "z").String()).Equal("b/x/y/z")
}

func TestJoinImmutable(t *testing.T) {
	r3 := curie.New("a:b/c/d")
	rP := curie.Parent(r3)
	rN := curie.Join(rP, "t")

	it.Ok(t).
		If(r3.String()).Should().Equal("a:b/c/d").
		If(rP.String()).Should().Equal("a:b/c").
		If(rN.String()).Should().Equal("a:b/c/t")
}

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

func TestTypeSafe(t *testing.T) {
	type A curie.IRI
	type B curie.IRI
	type C curie.IRI

	a := A(curie.New("a:"))
	b := B(curie.New("a:b"))
	c := C(curie.Join(curie.IRI(b), "c"))

	it.Ok(t).
		If(curie.IRI(a).String()).Should().Equal("a:").
		If(curie.IRI(b).String()).Should().Equal("a:b").
		If(curie.IRI(c).String()).Should().Equal("a:b/c")
}

func TestLinkedData(t *testing.T) {
	type Struct struct {
		ID curie.IRI     `json:"id"`
		LA *curie.String `json:"a,omitempty"`
		LB *curie.IRI    `json:"b,omitempty"`
	}

	b := curie.New("a:b/c")
	test := map[*Struct]string{
		{ID: curie.New(""), LA: curie.Safe(b), LB: &b}:      "{\"id\":\"\",\"a\":\"[a:b/c]\",\"b\":\"[a:b/c]\"}",
		{ID: curie.New("a:"), LA: curie.Safe(b), LB: &b}:    "{\"id\":\"[a:]\",\"a\":\"[a:b/c]\",\"b\":\"[a:b/c]\"}",
		{ID: curie.New("b"), LA: curie.Safe(b), LB: &b}:     "{\"id\":\"[b]\",\"a\":\"[a:b/c]\",\"b\":\"[a:b/c]\"}",
		{ID: curie.New("a:b"), LA: curie.Safe(b), LB: &b}:   "{\"id\":\"[a:b]\",\"a\":\"[a:b/c]\",\"b\":\"[a:b/c]\"}",
		{ID: curie.New("a:b/c"), LA: curie.Safe(b), LB: &b}: "{\"id\":\"[a:b/c]\",\"a\":\"[a:b/c]\",\"b\":\"[a:b/c]\"}",
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

/*
func TestID(t *testing.T) {
	type Struct struct {
		ID curie.ID `json:"id"`
	}

	test := map[*Struct]string{
		{ID: curie.ID{PKey: curie.New("")}}:                              "{\"id\":\"\"}",
		{ID: curie.ID{PKey: curie.New("a:")}}:                            "{\"id\":\"[a:]\"}",
		{ID: curie.ID{PKey: curie.New("a:b")}}:                           "{\"id\":\"[a:b]\"}",
		{ID: curie.ID{PKey: curie.New("a:b/c")}}:                         "{\"id\":\"[a:b/c]\"}",
		{ID: curie.ID{PKey: curie.New("a:b"), SKey: curie.New("a:")}}:    "{\"id\":\"[a:b][a:]\"}",
		{ID: curie.ID{PKey: curie.New("a:b"), SKey: curie.New("a:b")}}:   "{\"id\":\"[a:b][a:b]\"}",
		{ID: curie.ID{PKey: curie.New("a:b"), SKey: curie.New("a:b/c")}}: "{\"id\":\"[a:b][a:b/c]\"}",
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
