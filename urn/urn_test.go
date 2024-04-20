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

	"github.com/fogfish/curie/v2/urn"
	"github.com/fogfish/it/v2"
)

func TestNew(t *testing.T) {
	for expected, bag := range map[urn.URN][]string{
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

func TestNewFormat(t *testing.T) {
	for expected, bag := range map[urn.URN][]string{
		"urn:isbn":       {"isbn", ""},
		"urn:isbn:123":   {"isbn", "123"},
		"urn:isbn:1:2:3": {"isbn", "1:2:3"},
		"urn:isbn:1/2/3": {"isbn", "1/2/3"},
	} {
		t.Run(fmt.Sprintf("(%s)", expected), func(t *testing.T) {
			val := urn.New(bag[0], "%s", bag[1])
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
}

func TestIsEmpty(t *testing.T) {
	for urn, val := range map[urn.URN]bool{
		"":               true,
		"urn:":           true,
		"urn:isbn":       false,
		"urn:isbn:123":   false,
		"urn:isbn:1:2:3": false,
		"urn:isbn:1/2/3": false,
	} {
		t.Run(fmt.Sprintf("(%s)", urn), func(t *testing.T) {
			it.Then(t).
				Should(it.Equal(urn.IsEmpty(), val))
		})
	}
}
