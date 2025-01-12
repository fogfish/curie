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

	"github.com/fogfish/curie/v2"
	"github.com/fogfish/it/v2"
)

func TestNew(t *testing.T) {
	for expected, bag := range map[curie.IRI][]string{
		"":        {"", ""},
		"a:":      {"a", ""},
		"a:b":     {"a", "b"},
		"a:b/c":   {"a", "b/c"},
		"a:b/c/d": {"a", "b/c/d"},
		"b":       {"", "b"},
		"b/c":     {"", "b/c"},
		"b/c/d":   {"", "b/c/d"},
	} {
		t.Run(fmt.Sprintf("(%s)", expected), func(t *testing.T) {
			val := curie.New(bag[0], bag[1])
			schema, ref := val.Split()

			safeExp := "[" + string(expected) + "]"
			if len(expected) == 0 {
				safeExp = ""
			}

			it.Then(t).Should(
				it.Equal(val, expected),
				it.Equal(val.Safe(), safeExp),
				it.Equal(schema, bag[0]),
				it.Equal(ref, bag[1]),
				it.Equal(val.Schema(), schema),
				it.Equal(val.Reference(), ref),
			)
		})
	}
}

func TestNamespaceIRI(t *testing.T) {
	a := curie.Namespace("a")
	it.Then(t).Should(
		it.Equal(a.IRI("b/c"), "a:b/c"),
	)
}

func TestBasePath(t *testing.T) {
	for input, expected := range map[curie.IRI][]string{
		"":        {"", ""},
		"a:":      {"", "a:"},
		"a:b":     {"b", "a:"},
		"a:b/c":   {"c", "a:b"},
		"a:b/c/d": {"d", "a:b/c"},
		"b":       {"b", ""},
		"b/c":     {"c", "b"},
		"b/c/d":   {"d", "b/c"},
	} {
		t.Run(fmt.Sprintf("(%s)", expected), func(t *testing.T) {
			base := curie.Base(input)
			path := curie.Path(input)

			it.Then(t).Should(
				it.Equal(base, expected[0]),
				it.Equal(string(path), expected[1]),
			)
		})
	}
}

func TestHeadTail(t *testing.T) {
	for input, expected := range map[curie.IRI][]string{
		"":        {"", ""},
		"a:":      {"", "a:"},
		"a:b":     {"b", "a:"},
		"a:b/c":   {"b", "a:c"},
		"a:b/c/d": {"b", "a:c/d"},
		"b":       {"b", ""},
		"b/c":     {"b", "c"},
		"b/c/d":   {"b", "c/d"},
	} {
		t.Run(fmt.Sprintf("(%s)", expected), func(t *testing.T) {
			head := curie.Head(input)
			tail := curie.Tail(input)

			it.Then(t).Should(
				it.Equal(head, expected[0]),
				it.Equal(string(tail), expected[1]),
			)
		})
	}
}

func TestCodec(t *testing.T) {
	type Struct struct {
		ID curie.IRI `json:"id"`
	}

	for _, id := range []curie.IRI{
		"",
		"a:",
		"a:b",
		"a:b/c",
		"a:b/c/d",
		"b",
		"b/c",
		"b/c/d",
	} {
		t.Run(fmt.Sprintf("(%s)", id), func(t *testing.T) {
			send := Struct{ID: id}
			var recv Struct

			bytes, err1 := json.Marshal(send)
			err2 := json.Unmarshal(bytes, &recv)

			it.Then(t).Should(
				it.Nil(err1),
				it.Nil(err2),
				it.Equiv(send, recv),
				it.Equiv(bytes, []byte("{\"id\":\""+id+"\"}")),
			)
		})
	}

	for _, id := range []curie.IRI{
		"isbn",
		"isbn:123",
		"/1/2/3",
	} {
		t.Run(fmt.Sprintf("Decode (%s)", id), func(t *testing.T) {
			var recv Struct

			err2 := json.Unmarshal([]byte("{\"id\":\"["+string(id)+"]\"}"), &recv)

			it.Then(t).Should(
				it.Nil(err2),
				it.Equal(recv.ID, id),
			)
		})
	}
}

func TestCodecFail(t *testing.T) {
	type Struct struct {
		ID curie.IRI `json:"id"`
	}

	t.Run("Invalid type", func(t *testing.T) {
		var recv Struct

		err := json.Unmarshal([]byte("{\"id\":10}"), &recv)
		it.Then(t).ShouldNot(
			it.Nil(err),
		)
	})
}

func TestURL(t *testing.T) {
	prefixes := curie.Namespaces{
		"wikipedia": "http://en.wikipedia.org/wiki/",
	}
	t.Run("KnownPrefix", func(t *testing.T) {
		url, err := curie.URL(prefixes, curie.IRI("wikipedia:CURIE"))

		it.Then(t).Should(
			it.Nil(err),
			it.Equal(url.String(), "http://en.wikipedia.org/wiki/CURIE"),
		)
	})

	t.Run("UnknownPrefix", func(t *testing.T) {
		url, err := curie.URL(prefixes, curie.IRI("wiki:CURIE"))

		it.Then(t).Should(
			it.Nil(err),
			it.Equal(url.String(), "wiki:CURIE"),
		)
	})

	t.Run("PercentEncoded", func(t *testing.T) {
		url, err := curie.URL(prefixes, curie.IRI("wikipedia:Ῥόδος_%1F"))

		it.Then(t).Should(
			it.Nil(err),
			it.Equal(url.String(), "http://en.wikipedia.org/wiki/%E1%BF%AC%CF%8C%CE%B4%CE%BF%CF%82_%1F"),
		)
	})

	t.Run("Invalid URL", func(t *testing.T) {
		v := curie.URI(prefixes, curie.IRI("%2f:first_path_segment_in_URL_cannot_contain_colon"))

		it.Then(t).Should(
			it.Equal(v, "%2f:first_path_segment_in_URL_cannot_contain_colon"),
		)
	})
}

func TestFromURL(t *testing.T) {
	prefixes := curie.Namespaces{
		"wikipedia": "http://en.wikipedia.org/wiki/",
	}
	absolute := "http://en.wikipedia.org/wiki/CURIE"
	compact := curie.IRI("wikipedia:CURIE")

	t.Run("Identity", func(t *testing.T) {
		uri := curie.URI(prefixes, curie.FromURI(prefixes, absolute))

		it.Then(t).Should(
			it.Equal(uri, absolute),
		)
	})

	t.Run("KnownPrefix", func(t *testing.T) {
		iri := curie.FromURI(prefixes, absolute)

		it.Then(t).Should(
			it.Equal(iri, compact),
		)
	})

	t.Run("UnknownPrefix", func(t *testing.T) {
		iri := curie.FromURI(curie.Namespaces{}, absolute)

		it.Then(t).Should(
			it.Equal(iri, curie.IRI(absolute)),
		)
	})

	t.Run("PercentEncoded", func(t *testing.T) {
		iri := curie.FromURI(prefixes, `http://en.wikipedia.org/wiki/%E1%BF%AC%CF%8C%CE%B4%CE%BF%CF%82`)

		it.Then(t).Should(
			it.Equal(iri, curie.IRI("wikipedia:Ῥόδος")),
		)
	})

	t.Run("PercentEncodedCorrupted", func(t *testing.T) {
		iri := curie.FromURI(prefixes, `http://en.wikipedia.org/wiki/%%`)

		it.Then(t).Should(
			it.Equal(iri, curie.IRI(`wikipedia:%%`)),
		)
	})
}

func TestURLCompatibility(t *testing.T) {
	uri := "https://example.com/a/b/c?de=fg&foo=bar"
	curi := curie.IRI(uri)

	expect, _ := url.Parse(uri)
	native, err := curie.URL(curie.Namespaces{}, curi)
	schema, path := curie.Split(curi)

	it.Then(t).Should(
		it.Equal(string(curi), uri),
		it.Equal(curi.Safe(), "["+uri+"]"),
		it.Equal(schema, "https"),
		it.Equal(path, "//example.com/a/b/c?de=fg&foo=bar"),

		//
		it.Nil(err),
		it.Equiv(native, expect),
	)
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
		curi := curie.IRI(compact)
		expect, _ := url.Parse(v[1])
		uri, err := curie.URL(curie.Namespaces{"a": v[0]}, curi)

		it.Then(t).Should(
			it.Nil(err),
			it.Equiv(uri, expect),
		)
	}
}

func TestTypeSafe(t *testing.T) {
	type A curie.IRI
	type B curie.IRI
	type C curie.IRI

	a := A(curie.New("a", ""))
	b := B(curie.New("a", "b"))
	c := C(curie.Join(curie.IRI(b), "c"))

	it.Then(t).Should(
		it.Equal(string(curie.IRI(a)), "a:"),
		it.Equal(string(curie.IRI(b)), "a:b"),
		it.Equal(string(curie.IRI(c)), "a:b/c"),
	)
}

func TestJoin(t *testing.T) {
	for _, id := range []curie.IRI{
		"",
		"a:",
		"a:b",
		"a:b/c",
		"a:b/c/d",
		"b",
		"b/c",
		"b/c/d",
	} {
		t.Run(fmt.Sprintf("(%s)", id), func(t *testing.T) {
			it.Then(t).Should(
				it.Equal(id.Join("x").Cut(1), id),
			)
		})
	}
}
