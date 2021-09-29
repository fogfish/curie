//
// Copyright (C) 2020 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/fogfish/curie
//

/*

Package curie implements the type for compact URI. It defines a generic syntax
for expressing URIs by abbreviated literal as defined by the W3C.
https://www.w3.org/TR/2010/NOTE-curie-20101216/
*/
package curie

import (
	"encoding/json"
	"fmt"
	"net/url"
	"path"
	"strings"
)

//------------------------------------------------------------------------------
//
// IRI
//
//------------------------------------------------------------------------------

/*

IRI is compact URI, defined as superset of XML QNames.
  safe_curie  :=   '[' curie ']'
  curie       :=   [ [ prefix ] ':' ] suffix
  prefix      :=   NCName
  suffix      :=   NCName [ / suffix ]
*/
type IRI struct {
	// sequence of IRI segments
	seq []string
}

/*

String transform CURIE to string
*/
func (iri IRI) String() string {
	return join(iri.seq)
}

/*

Safe transforms CURIE to safe string
*/
func (iri IRI) Safe() string {
	val := iri.String()
	if val == "" {
		return ""
	}

	return "[" + val + "]"
}

/*

MarshalJSON `IRI ⟼ "[prefix:suffix]"`
*/
func (iri IRI) MarshalJSON() ([]byte, error) {
	if Rank(iri) == 0 {
		return json.Marshal("")
	}

	return json.Marshal(iri.Safe())
}

/*

UnmarshalJSON `"[prefix:suffix]" ⟼ IRI`
*/
func (iri *IRI) UnmarshalJSON(b []byte) error {
	var path string
	err := json.Unmarshal(b, &path)
	if err != nil {
		return err
	}

	*iri = New(path)
	return nil
}

//------------------------------------------------------------------------------
//
// Pair
//
//------------------------------------------------------------------------------

/*

ID is a product type of IRIs
*/
// type ID struct{ PKey, SKey IRI }

/*

MarshalJSON `ID ⟼ "[prefix:suffix][prefix:suffix]"`
*/
// func (id ID) MarshalJSON() ([]byte, error) {
// 	pkey := id.PKey.Safe()
// 	skey := id.SKey.Safe()

// 	if skey == "" {
// 		return json.Marshal(pkey)
// 	}

// 	return json.Marshal(pkey + skey)
// }

/*

UnmarshalJSON `"[prefix:suffix][prefix:suffix]" ⟼ ID`
*/
// func (id *ID) UnmarshalJSON(b []byte) error {
// 	var val string
// 	err := json.Unmarshal(b, &val)
// 	if err != nil {
// 		return err
// 	}

// 	seq := strings.Split(val, "][")
// 	id.PKey = New(seq[0])
// 	if len(seq) == 1 {
// 		return nil
// 	}

// 	id.SKey = New(seq[1])
// 	return nil
// }

//------------------------------------------------------------------------------
//
// String
//
//------------------------------------------------------------------------------

/*

String is safe representation of IRI
*/
type String string

/*

IRI convers Safe CURIE to IRI type
*/
func (s String) IRI() IRI {
	return New(string(s))
}

//------------------------------------------------------------------------------
//
// IRI algebra
//
//------------------------------------------------------------------------------

/*

New transform category of strings to IRI.
*/
func New(iri string, args ...interface{}) IRI {
	val := iri

	if len(args) > 0 {
		val = fmt.Sprintf(iri, args...)
	}

	return IRI{seq: split(strings.Trim(val, "[]"))}
}

/*

IsEmpty is an alias to curie.Rank(iri) == 0
*/
func IsEmpty(iri IRI) bool {
	return len(iri.seq) == 0
}

/*

Rank of CURIE, number of segments
Rank is an alias of len(curie.Seq(iri))
*/
func Rank(iri IRI) int {
	return len(iri.seq)
}

/*

Seq Returns CURIE segments
  a: ⟼ [ a ]
  a:b/c ⟼ [a, b, c]
  a:b/c/d ⟼ [a, b, c, d ]
*/
func Seq(iri IRI) []string {
	return iri.seq
}

/*

Safe transforms CURIE to safe string
*/
func Safe(iri IRI) *String {
	val := String(iri.Safe())
	return &val
}

/*

Path converts CURIE to relative file system path.
  a: ⟼ a
  a:b/c ⟼ a/b/c
  a:b/c/d ⟼ a/b/c/d
*/
func Path(iri IRI) string {
	return path.Join(iri.seq...)
}

/*

URI converts CURIE to fully qualified URL
  wikipedia:CURIE ⟼ http://en.wikipedia.org/wiki/CURIE
*/
func URI(prefixes map[string]string, iri IRI) string {
	if IsEmpty(iri) {
		return ""
	}
	seq := iri.seq

	if seq[0] == "" {
		return strings.Join(seq[1:], "/")
	}

	prefix, exists := prefixes[seq[0]]
	if !exists {
		return iri.String()
	}

	return prefix + strings.Join(seq[1:], "/")
}

/*

URL converts CURIE to fully qualified url.URL type
  wikipedia:CURIE ⟼ http://en.wikipedia.org/wiki/CURIE
*/
func URL(prefixes map[string]string, iri IRI) (*url.URL, error) {
	uri := URI(prefixes, iri)

	if uri == "" {
		return &url.URL{}, nil
	}

	return url.Parse(uri)
}

/*

Eq compares two CURIEs, return true if they are equal.
*/
func Eq(a, b IRI) bool {
	if Rank(a) != Rank(b) {
		return false
	}

	seq := b.seq
	for i, v := range a.seq {
		if seq[i] != v {
			return false
		}
	}

	return true
}

/*

Lt compares two CURIEs, return true if left element is less than supplied one.
*/
func Lt(a, b IRI) bool {
	if Rank(a) != Rank(b) {
		return Rank(a) < Rank(b)
	}

	seq := b.seq
	for i, v := range a.seq {
		if seq[i] != v {
			return v < seq[i]
		}
	}

	return false
}

/*

Prefix decomposes CURIE and return its prefix CURIE as string value.
*/
func Prefix(iri IRI) string {
	if len(iri.seq) == 0 {
		return ""
	}
	return iri.seq[0]
}

/*

Suffix decomposes CURIE and return its suffix.
*/
func Suffix(iri IRI) string {
	if len(iri.seq) < 2 {
		return ""
	}

	return strings.Join(iri.seq[1:], "/")
}

/*

Join composes segments into new descendant CURIE.
  a:b × [c, d, e] ⟼ a:b/c/d/e
*/
func Join(iri IRI, segments ...string) IRI {
	seq := append([]string{}, iri.seq...)
	if len(seq) == 0 {
		seq = append(seq, "")
	}

	for _, x := range segments {
		seq = append(seq,
			strings.FieldsFunc(x,
				func(x rune) bool {
					return x == '/' || x == ':'
				},
			)...,
		)
	}

	return IRI{seq: seq}
}

/*

Parent decomposes CURIE and return its parent as CURIE type.

  a:b/c/d/e ⟼¹ a:b/c/d  a:b/c/d/e ⟼⁻¹ a:
  a:b/c/d/e ⟼² a:b/c    a:b/c/d/e ⟼⁻² a:b
  a:b/c/d/e ⟼³ a:b      a:b/c/d/e ⟼⁻³ a:b/c
  ...
  a:b/c/d/e ⟼ⁿ a:       a:b/c/d/e ⟼⁻ⁿ a:b/c/d/e
*/
func Parent(iri IRI, rank ...int) IRI {
	r := 1
	if len(rank) > 0 {
		r = rank[0]
	}
	if r < 0 {
		r = len(iri.seq) + r
	}

	n := len(iri.seq) - r
	switch {
	case n < 0:
		return IRI{seq: []string{}}
	case n > len(iri.seq):
		return IRI{seq: append([]string{}, iri.seq...)}
	case n == 1 && iri.seq[0] == "":
		return IRI{seq: []string{}}
	default:
		return IRI{seq: append([]string{}, iri.seq[:n]...)}
	}
}

/*

Child decomposes CURIE and return its suffix as CURIE type.

  a:b/c/d/e ⟿¹ e          a:b/c/d/e ⟿⁻¹ b/c/d/e
  a:b/c/d/e ⟿² d/e        a:b/c/d/e ⟿⁻² c/d/e
  a:b/c/d/e ⟿³ c/d/e      a:b/c/d/e ⟿⁻³ d/e
  ...
  a:b/c/d/e ⟿ⁿ a:b/c/d/e  a:b/c/d/e ⟿⁻ⁿ e
*/
func Child(iri IRI, rank ...int) string {
	r := 1
	if len(rank) > 0 {
		r = rank[0]
	}

	if r < 0 {
		r = len(iri.seq) + r
	}

	n := Rank(iri) - r
	switch {
	case n >= len(iri.seq):
		return ""
	case n > 0:
		return strings.Join(iri.seq[n:len(iri.seq)], "/")
	default:
		return join(iri.seq)
	}
}

//------------------------------------------------------------------------------
//
// private
//
//------------------------------------------------------------------------------

/*

split parses tokens `prefix:suffix[/suffix/...]`
*/
func split(val string) []string {
	seq := strings.Split(val, ":")

	// zero
	if len(seq) == 1 && seq[0] == "" {
		return []string{}
	}

	// zeroPrefix
	if len(seq) == 1 {
		return append([]string{""}, strings.Split(seq[0], "/")...)
	}

	// zeroSuffix
	if seq[1] == "" {
		return []string{seq[0]}
	}

	return append([]string{seq[0]}, strings.Split(seq[1], "/")...)
}

/*

join sequence to string `prefix:suffix[/suffix/...]`
*/
func join(seq []string) string {
	// zero
	if len(seq) == 0 || (len(seq) == 1 && seq[0] == "") {
		return ""
	}

	// zeroSuffix
	if len(seq) == 1 && seq[0] != "" {
		return seq[0] + ":"
	}

	// zeroPrefix
	if seq[0] == "" {
		return strings.Join(seq[1:], "/")
	}

	return seq[0] + ":" + strings.Join(seq[1:], "/")
}
