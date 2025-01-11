//
// Copyright (C) 2020 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/fogfish/curie
//

package urn_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/fogfish/curie/v2"
	"github.com/fogfish/curie/v2/urn"
	"github.com/fogfish/it/v2"
)

func TestNew(t *testing.T) {
	for expected, bag := range map[urn.URN][]string{
		"urn:":           {"", ""},
		"urn:isbn":       {"isbn", ""},
		"urn:isbn:123":   {"isbn", "123"},
		"urn:isbn:1:2:3": {"isbn", "1:2:3"},
		"urn:isbn:1/2/3": {"isbn", "1/2/3"},
	} {
		t.Run(fmt.Sprintf("(%s)", expected), func(t *testing.T) {
			val := urn.New(bag[0], bag[1])
			schema, ref := val.Split()

			it.Then(t).Should(
				it.Equal(val, expected),
				it.Equal(schema, bag[0]),
				it.Equal(ref, bag[1]),
				it.Equal(val.Schema(), schema),
				it.Equal(val.Reference(), ref),
			)
		})
	}
}

func TestBasePath(t *testing.T) {
	for input, expected := range map[urn.URN][]string{
		"":               {"", ""},
		"urn:":           {"", "urn:"},
		"urn:isbn":       {"", "urn:isbn"},
		"urn:isbn:b":     {"b", "urn:isbn"},
		"urn:isbn:b:c":   {"c", "urn:isbn:b"},
		"urn:isbn:b:c:d": {"d", "urn:isbn:b:c"},
	} {
		t.Run(fmt.Sprintf("(%s)", expected), func(t *testing.T) {
			base := urn.Base(input)
			path := urn.Path(input)

			it.Then(t).Should(
				it.Equal(base, expected[0]),
				it.Equal(string(path), expected[1]),
			)
		})
	}
}

func TestHeadTail(t *testing.T) {
	for input, expected := range map[urn.URN][]string{
		"":               {"", ""},
		"urn:":           {"", "urn:"},
		"urn:isbn":       {"", "urn:isbn"},
		"urn:isbn:b":     {"b", "urn:isbn"},
		"urn:isbn:b:c":   {"b", "urn:isbn:c"},
		"urn:isbn:b:c:d": {"b", "urn:isbn:c:d"},
	} {
		t.Run(fmt.Sprintf("(%s)", expected), func(t *testing.T) {
			head := urn.Head(input)
			tail := urn.Tail(input)

			it.Then(t).Should(
				it.Equal(head, expected[0]),
				it.Equal(string(tail), expected[1]),
			)
		})
	}
}

func TestCodec(t *testing.T) {
	type Struct struct {
		ID urn.URN `json:"id"`
	}

	for _, id := range []urn.URN{
		"",
		"urn:isbn",
		"urn:isbn:123",
		"urn:isbn:1:2:3",
		"urn:isbn:1/2/3",
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
				it.Equiv(bytes, []byte("{\"id\":\""+string(id)+"\"}")),
			)
		})
	}
}

func TestCodecFail(t *testing.T) {
	type Struct struct {
		ID urn.URN `json:"id"`
	}

	for _, id := range []urn.URN{
		"isbn",
		"isbn:123",
		"/1/2/3",
	} {
		t.Run(fmt.Sprintf("Encode (%s)", id), func(t *testing.T) {
			send := Struct{ID: id}

			_, err1 := json.Marshal(send)

			it.Then(t).ShouldNot(
				it.Nil(err1),
			)
		})

		t.Run(fmt.Sprintf("Decode (%s)", id), func(t *testing.T) {
			var recv Struct

			err2 := json.Unmarshal([]byte("{\"id\":\""+string(id)+"\"}"), &recv)

			it.Then(t).ShouldNot(
				it.Nil(err2),
			)
		})
	}

	t.Run("Invalid type", func(t *testing.T) {
		var recv Struct

		err := json.Unmarshal([]byte("{\"id\":10}"), &recv)
		it.Then(t).ShouldNot(
			it.Nil(err),
		)
	})
}

func TestJoin(t *testing.T) {
	for _, id := range []urn.URN{
		"urn:",
		"urn:isbn",
		"urn:isbn:123",
		"urn:isbn:1:2:3",
		"urn:isbn:1/2/3",
	} {
		t.Run(fmt.Sprintf("(%s)", id), func(t *testing.T) {
			it.Then(t).Should(
				it.Equal(id.Join("x").Cut(1), id),
			)
		})
	}
}

func TestUrn2Iri(t *testing.T) {
	for URN, IRI := range map[urn.URN]curie.IRI{
		"urn:":           "",
		"urn:isbn":       "isbn:",
		"urn:isbn:123":   "isbn:123",
		"urn:isbn:1:2:3": "isbn:1/2/3",
	} {
		t.Run(fmt.Sprintf("(%s)", URN), func(t *testing.T) {
			it.Then(t).Should(
				it.Equal(urn.ToIRI(URN), IRI),
				it.Equal(urn.ToURN(IRI), URN),
			)
		})
	}
}
