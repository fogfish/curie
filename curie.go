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
  curie       :=   [ [ scheme ] ':' ] reference
  scheme      :=   NCName
  reference   :=   prefix [ / suffix ]
  prefix      :=   irelative-ref (as defined in IRI)
  suffix      :=   irelative-ref (as defined in IRI)
*/
type IRI struct {
	// number of suffix segments
	rankSuffix int
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
	return "[" + iri.String() + "]"
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

/*

New transform category of strings to IRI.
*/
func New(iri string, args ...interface{}) IRI {
	val := iri

	if len(args) > 0 {
		val = fmt.Sprintf(iri, args...)
	}

	rankSuffix := 1
	seq := split(strings.Trim(val, "[]"))
	if len(seq) < 3 {
		rankSuffix = 0
	}

	return IRI{rankSuffix: rankSuffix, seq: seq}
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
  a:b/c/d/e ⟼ [a, b, c, d]
*/
func Seq(iri IRI) []string {
	return iri.seq
}

/*

Safe transforms CURIE to safe string
*/
func Safe(iri IRI) *String {
	val := String("[" + iri.String() + "]")
	return &val
}

/*

Path converts CURIE to relative file system path.

  a:b/c/d/e ⟼ a/b/c/d/e
*/
func Path(iri IRI) string {
	return path.Join(iri.seq...)
}

/*

URI converts CURIE to fully qualified URL
  wikipedia:CURIE ⟼ http://en.wikipedia.org/wiki/CURIE
*/
func URI(prefix string, iri IRI) (*url.URL, error) {
	if IsEmpty(iri) {
		return &url.URL{}, nil
	}
	seq := iri.seq

	if seq[0] == "" {
		return url.Parse(strings.Join(seq[1:], "/"))
	}

	return url.Parse(prefix + strings.Join(seq[1:], "/"))
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

Split IRI defining a new suffix rank, defining behavior of splitters
such as Parent, Prefix Suffix. Every new IRI has suffix rank 1

  a:b/c/d/e ⟼⁰   a:b/c/d/e |
  a:b/c/d/e ⟼¹   a:b/c/d | e
  a:b/c/d/e ⟼²   a:b/c | d/e
  a:b/c/d/e ⟼³   a:b | c/d/e
  ...
  a:b/c/d/e ⟼ⁿ⁻² a:b | c/d/e
  a:b/c/d/e ⟼ⁿ⁻¹ a:b | c/d/e
  a:b/c/d/e ⟼ⁿ   a:b | c/d/e
*/
func Split(iri IRI, rankSuffix int) IRI {
	seq := append([]string{}, iri.seq...)

	if len(iri.seq) <= 2 {
		return IRI{rankSuffix: 0, seq: seq}
	}

	r := rankSuffix
	if r > len(iri.seq)-2 {
		r = len(iri.seq) - 2
	}

	return IRI{rankSuffix: r, seq: seq}
}

/*

Scheme ...
*/
func Scheme(iri IRI) string {
	if len(iri.seq) == 0 {
		return ""
	}
	return iri.seq[0]
}

/*

Parent decomposes CURIE and return its parent CURIE.
It return immediate parent compact URI if rank is not defined.
*/
func Parent(iri IRI) IRI {
	seq := iri.seq[:len(iri.seq)-iri.rankSuffix]
	rankSuffix := 1
	if len(seq) < 3 {
		rankSuffix = 0
	}

	return IRI{rankSuffix: rankSuffix, seq: append([]string{}, seq...)}
}

/*

Prefix decomposes CURIE and return its prefix CURIE as string value.
*/
func Prefix(iri IRI) string {
	return join(iri.seq[:len(iri.seq)-iri.rankSuffix])
}

/*

Suffix decomposes CURIE and return its suffix.
*/
func Suffix(iri IRI) string {
	if iri.rankSuffix == 0 {
		return ""
	}

	return strings.Join(iri.seq[len(iri.seq)-iri.rankSuffix:], "/")
}

/*

Join composes segments into new descendant CURIE.
  a:b × [c, d, e] ⟼ a:b/c/d/e
*/
func Join(iri IRI, segments ...string) IRI {
	seq := append([]string{}, iri.seq...)
	for _, x := range segments {
		seq = append(seq,
			strings.FieldsFunc(x,
				func(x rune) bool {
					return x == '/' || x == ':'
				},
			)...,
		)
	}

	rankSuffix := 1
	if len(seq) < 3 {
		rankSuffix = 0
	}

	return IRI{rankSuffix: rankSuffix, seq: seq}
}

/*

Heir composes two CURIEs into new descendant CURIE.
	a:b × c/d/e ⟼ a:b/c/d/e
	a:b × c:d/e ⟼ a:b/c/d/e
*/
func Heir(prefix, suffix IRI) IRI {
	seq := append(append([]string{}, prefix.seq...), suffix.seq...)

	rankSuffix := 1
	if len(seq) < 3 {
		rankSuffix = 0
	}

	return IRI{rankSuffix: rankSuffix, seq: seq}
}

//------------------------------------------------------------------------------
//
// private
//
//------------------------------------------------------------------------------

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
