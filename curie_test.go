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
	"strings"
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

var (
	IRIs        = []string{i000, i100, i010, i110, i120, i130, i011, i111, i112, i121, i122, i133}
	nonZeroIRIs = []string{i110, i120, i130, i011, i111, i112, i121, i122, i133}
)

func TestNew(t *testing.T) {
	for _, str := range IRIs {
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

func TestNewFormat(t *testing.T) {
	it.Ok(t).
		If(curie.New("a:b/c/d/%s", "e").String()).Equal("a:b/c/d/e")
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

func TestPath(t *testing.T) {
	for str, seq := range map[string]string{
		i000: "",
		i100: "a",
		i010: "b",
		i110: "a/b",
		i120: "a/b/c",
		i130: "a/b/c/d",
		i011: "b/c",
		i111: "a/b/c",
		i112: "a/b/c/d",
		i121: "a/b/c/d",
		i122: "a/b/c/d/e",
		i133: "a/b/c/d/e/f/g",
	} {
		t.Run(fmt.Sprintf("(%s)", str), func(t *testing.T) {
			iri := curie.New(str)

			it.Ok(t).
				If(curie.Path(iri)).Equal(seq)
		})
	}
}

func TestURI(t *testing.T) {
	compact := curie.New("wikipedia:CURIE")
	url, err := curie.URI("http://en.wikipedia.org/wiki/", compact)

	it.Ok(t).
		IfNil(err).
		If(url.String()).Equal("http://en.wikipedia.org/wiki/CURIE")
}

func TestURICompatibility(t *testing.T) {
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

func TestURIConvert(t *testing.T) {
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

func TestEq(t *testing.T) {
	for _, str := range IRIs {
		t.Run(fmt.Sprintf("(%s)", str), func(t *testing.T) {
			a := curie.New(str)
			b := curie.New(str)
			c := curie.New(w120)

			it.Ok(t).
				If(curie.Eq(a, b)).Should().Equal(true).
				If(curie.Eq(a, c)).Should().Equal(false)
		})
	}
}

func TestNotEq(t *testing.T) {
	r6 := curie.New("1:2:3:4:5:6")

	for _, str := range IRIs {
		t.Run(fmt.Sprintf("(%s)", str), func(t *testing.T) {
			a := curie.New(str)

			it.Ok(t).
				If(curie.Eq(a, r6)).Should().Equal(false)
		})
	}
}

func TestLt(t *testing.T) {
	for a, b := range map[string]string{
		i000:      "b:",
		i100:      "b:",
		i010:      "c",
		i110:      "a:c",
		i120:      "a:b/d",
		i130:      "a:b/c/e",
		i011:      "b#d",
		i111:      "a:b#d",
		i112:      "a:b#c/e",
		i121:      "a:b/c#e",
		i122:      "a:b/c#d/f",
		i133:      "a:b/c/d#e/f/h",
		"a:x/x/a": "a:x/x/x/a",
	} {
		t.Run(fmt.Sprintf("(%s)", a), func(t *testing.T) {
			it.Ok(t).
				If(curie.Lt(curie.New(a), curie.New(b))).Should().Equal(true).
				If(curie.Lt(curie.New(b), curie.New(a))).Should().Equal(false).
				If(curie.Lt(curie.New(a), curie.New(a))).Should().Equal(false)
		})
	}
}

func TestScheme(t *testing.T) {
	for str, val := range map[string]string{
		i000: "",
		i100: "a",
		i010: "",
		i110: "a",
		i120: "a",
		i130: "a",
		i011: "",
		i111: "a",
		i112: "a",
		i121: "a",
		i122: "a",
		i133: "a",
	} {
		t.Run(fmt.Sprintf("(%s)", str), func(t *testing.T) {
			iri := curie.New(str)

			it.Ok(t).
				If(curie.Scheme(iri)).Should().Equal(val)
		})
	}
}

func TestNewScheme(t *testing.T) {
	for _, str := range IRIs {
		t.Run(fmt.Sprintf("(%s)", str), func(t *testing.T) {
			iri := curie.NewScheme(curie.New(str), "t")

			it.Ok(t).
				If(curie.Scheme(iri)).Should().Equal("t")
		})
	}
}

func TestPrefixAndParent(t *testing.T) {
	for str, val := range map[string]string{
		i000: "",
		i100: "a:",
		i010: "b",
		i110: "a:b",
		i120: "a:b/c",
		i130: "a:b/c/d",
		i011: "b",
		i111: "a:b",
		i112: "a:b",
		i121: "a:b/c",
		i122: "a:b/c",
		i133: "a:b/c/d",
	} {
		t.Run(fmt.Sprintf("(%s)", str), func(t *testing.T) {
			iri := curie.New(str)

			it.Ok(t).
				If(curie.Prefix(iri)).Equal(val).
				If(curie.Parent(iri).String()).Equal(val)
		})
	}
}

func TestSuffixAndChild(t *testing.T) {
	for str, val := range map[string]string{
		i000: "",
		i100: "",
		i010: "",
		i110: "",
		i120: "",
		i130: "",
		i011: "c",
		i111: "c",
		i112: "c/d",
		i121: "d",
		i122: "d/e",
		i133: "e/f/g",
	} {
		t.Run(fmt.Sprintf("(%s)", str), func(t *testing.T) {
			iri := curie.New(str)

			it.Ok(t).
				If(curie.Suffix(iri)).Equal(val).
				If(curie.Path(curie.Child(iri))).Equal(val)
		})
	}
}

func TestChildScheme(t *testing.T) {
	it.Ok(t).
		If(curie.Child(curie.New("")).String()).Equal("").
		If(curie.Child(curie.New("a:b#c")).String()).Equal("c:").
		If(curie.Child(curie.New("a:b#c/d")).String()).Equal("c:d").
		If(curie.Child(curie.New("a:b#c/d/e")).String()).Equal("c:d/e")
}

func TestSplit(t *testing.T) {
	iri := curie.New(i133)

	for n, val := range map[int]string{
		0: "a:b/c/d/e/f/g",
		1: "a:b/c/d/e/f#g",
		2: "a:b/c/d/e#f/g",
		3: "a:b/c/d#e/f/g",
		4: "a:b/c#d/e/f/g",
		5: "a:b#c/d/e/f/g",
		6: "a:b#c/d/e/f/g",
		7: "a:b#c/d/e/f/g",
		8: "a:b#c/d/e/f/g",
		9: "a:b#c/d/e/f/g",
	} {
		t.Run(fmt.Sprintf("(%s, %d)", i133, n), func(t *testing.T) {
			newVal := curie.Split(iri, n)

			it.Ok(t).
				If(newVal.String()).Equal(val)
		})
	}
}

func TestSplitNoSuffix(t *testing.T) {
	iri := curie.New("a:b/c/d/e/f/g")

	for n, val := range map[int]string{
		0: "a:b/c/d/e/f/g",
		1: "a:b/c/d/e/f#g",
		2: "a:b/c/d/e#f/g",
		3: "a:b/c/d#e/f/g",
		4: "a:b/c#d/e/f/g",
		5: "a:b#c/d/e/f/g",
		6: "a:b#c/d/e/f/g",
		7: "a:b#c/d/e/f/g",
		8: "a:b#c/d/e/f/g",
		9: "a:b#c/d/e/f/g",
	} {
		t.Run(fmt.Sprintf("(%s, %d)", i133, n), func(t *testing.T) {
			newVal := curie.Split(iri, n)

			it.Ok(t).
				If(newVal.String()).Equal(val)
		})
	}
}

func TestSplitZero(t *testing.T) {
	it.Ok(t).
		If(curie.Split(curie.New(i000), 1).String()).Equal("").
		If(curie.Split(curie.New(i100), 1).String()).Equal("a:").
		If(curie.Split(curie.New(i010), 1).String()).Equal("b")
}

func TestJoin(t *testing.T) {
	for _, str := range nonZeroIRIs {
		t.Run(fmt.Sprintf("(%s)", str), func(t *testing.T) {
			iri := curie.New(str)
			af1 := curie.Join(iri, "x")
			af2 := curie.Join(iri, "x", "y")
			af3 := curie.Join(iri, "x", "y", "z")

			bas := strings.ReplaceAll(str, "#", "/")

			it.Ok(t).
				If(af1.String()).Equal(bas + "#x").
				If(af2.String()).Equal(bas + "#x/y").
				If(af3.String()).Equal(bas + "#x/y/z").
				If(curie.Prefix(af1)).Equal(bas).
				If(curie.Prefix(af2)).Equal(bas).
				If(curie.Prefix(af3)).Equal(bas).
				IfTrue(curie.Eq(curie.Parent(af1), curie.New(bas))).
				IfTrue(curie.Eq(curie.Parent(af2), curie.New(bas))).
				IfTrue(curie.Eq(curie.Parent(af3), curie.New(bas)))
		})
	}
}

func TestJoinWithZero(t *testing.T) {
	it.Ok(t).
		If(curie.Join(curie.New(i000), "x").String()).Equal("x").
		If(curie.Join(curie.New(i000), "x", "y").String()).Equal("x#y").
		If(curie.Join(curie.New(i000), "x", "y", "z").String()).Equal("x#y/z").
		If(curie.Join(curie.New(i100), "x").String()).Equal("a:x").
		If(curie.Join(curie.New(i100), "x", "y").String()).Equal("a:x#y").
		If(curie.Join(curie.New(i100), "x", "y", "z").String()).Equal("a:x#y/z").
		If(curie.Join(curie.New(i010), "x").String()).Equal("b#x").
		If(curie.Join(curie.New(i010), "x", "y").String()).Equal("b#x/y").
		If(curie.Join(curie.New(i010), "x", "y", "z").String()).Equal("b#x/y/z")
}

func TestJoinImmutable(t *testing.T) {
	r3 := curie.New(i133)
	rP := curie.Parent(r3)
	rN := curie.Join(rP, "t")

	it.Ok(t).
		If(r3.String()).Should().Equal("a:b/c/d#e/f/g").
		If(rP.String()).Should().Equal("a:b/c/d").
		If(rN.String()).Should().Equal("a:b/c/d#t")
}

func TestHeir(t *testing.T) {
	for _, str := range nonZeroIRIs {
		t.Run(fmt.Sprintf("(%s)", str), func(t *testing.T) {
			iri := curie.New(str)
			af1 := curie.Heir(iri, curie.New("x:"))
			af2 := curie.Heir(iri, curie.New("x:y"))
			af3 := curie.Heir(iri, curie.New("x:y#z"))

			bas := strings.ReplaceAll(str, "#", "/")

			it.Ok(t).
				If(af1.String()).Equal(bas + "#x").
				If(af2.String()).Equal(bas + "#x/y").
				If(af3.String()).Equal(bas + "#x/y/z").
				If(curie.Prefix(af1)).Equal(bas).
				If(curie.Prefix(af2)).Equal(bas).
				If(curie.Prefix(af3)).Equal(bas).
				IfTrue(curie.Eq(curie.Parent(af1), curie.New(bas))).
				IfTrue(curie.Eq(curie.Parent(af2), curie.New(bas))).
				IfTrue(curie.Eq(curie.Parent(af3), curie.New(bas)))
		})
	}
}

func TestHeirWithZero(t *testing.T) {
	it.Ok(t).
		If(curie.Heir(curie.New(i000), curie.New("x:")).String()).Equal("x").
		If(curie.Heir(curie.New(i000), curie.New("x:y")).String()).Equal("x#y").
		If(curie.Heir(curie.New(i000), curie.New("x:y#z")).String()).Equal("x#y/z").
		If(curie.Heir(curie.New(i100), curie.New("x:")).String()).Equal("a:x").
		If(curie.Heir(curie.New(i100), curie.New("x:y")).String()).Equal("a:x#y").
		If(curie.Heir(curie.New(i100), curie.New("x:y#z")).String()).Equal("a:x#y/z").
		If(curie.Heir(curie.New(i010), curie.New("x:")).String()).Equal("b#x").
		If(curie.Heir(curie.New(i010), curie.New("x:y")).String()).Equal("b#x/y").
		If(curie.Heir(curie.New(i010), curie.New("x:y#z")).String()).Equal("b#x/y/z")
}

func TestHeirWithZeroSuffix(t *testing.T) {
	it.Ok(t).
		If(curie.Heir(curie.New(i111), curie.New("")).String()).Equal("a:b#c")
}

func TestHeirImmutable(t *testing.T) {
	r3 := curie.New(i133)
	rP := curie.Parent(r3)
	rN := curie.Heir(rP, curie.New("t"))

	it.Ok(t).
		If(r3.String()).Should().Equal("a:b/c/d#e/f/g").
		If(rP.String()).Should().Equal("a:b/c/d").
		If(rN.String()).Should().Equal("a:b/c/d#t")
}

func TestTypeSafe(t *testing.T) {
	type A curie.IRI
	type B curie.IRI
	type C curie.IRI

	a := A(curie.New(i100))
	b := B(curie.New(i110))
	c := C(curie.Join(curie.IRI(b), "c"))

	it.Ok(t).
		If(curie.IRI(a).String()).Should().Equal(i100).
		If(curie.IRI(b).String()).Should().Equal(i110).
		If(curie.IRI(c).String()).Should().Equal(i111)
}

func TestLinkedData(t *testing.T) {
	type Struct struct {
		ID curie.IRI     `json:"id"`
		LA *curie.String `json:"a,omitempty"`
		LB *curie.IRI    `json:"b,omitempty"`
	}

	b := curie.New(i120)
	test := map[*Struct]string{
		{ID: curie.New(i000), LA: curie.Safe(b), LB: &b}: "{\"id\":\"\",\"a\":\"[a:b/c]\",\"b\":\"[a:b/c]\"}",
		{ID: curie.New(i100), LA: curie.Safe(b), LB: &b}: "{\"id\":\"[a:]\",\"a\":\"[a:b/c]\",\"b\":\"[a:b/c]\"}",
		{ID: curie.New(i010), LA: curie.Safe(b), LB: &b}: "{\"id\":\"[b]\",\"a\":\"[a:b/c]\",\"b\":\"[a:b/c]\"}",
		{ID: curie.New(i110), LA: curie.Safe(b), LB: &b}: "{\"id\":\"[a:b]\",\"a\":\"[a:b/c]\",\"b\":\"[a:b/c]\"}",
		{ID: curie.New(i120), LA: curie.Safe(b), LB: &b}: "{\"id\":\"[a:b/c]\",\"a\":\"[a:b/c]\",\"b\":\"[a:b/c]\"}",
		{ID: curie.New(i111), LA: curie.Safe(b), LB: &b}: "{\"id\":\"[a:b#c]\",\"a\":\"[a:b/c]\",\"b\":\"[a:b/c]\"}",
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
