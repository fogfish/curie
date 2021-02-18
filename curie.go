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
  curie       :=   [ [ prefix ] ':' ] reference
  prefix      :=   NCName
  reference   :=   irelative-ref (as defined in IRI)
*/
type IRI struct {
	seq []string
}

/*

New transform category of strings to IRI.
*/
func New(iri string, args ...interface{}) IRI {
	val := iri

	if len(args) > 0 {
		val = fmt.Sprintf(iri, args...)
	}

	return IRI{
		seq: split(strings.Trim(val, "[]")),
	}
}

/*

This is helper function to lift IRI into pointer
*/
func (iri IRI) This() *IRI {
	return &iri
}

/*

IsEmpty is an alias to iri.Rank() == 0
*/
func (iri IRI) IsEmpty() bool {
	return len(iri.seq) == 0
}

/*

Rank of CURIE, number of segments
Rank is an alias of len(iri.Seq())
*/
func (iri IRI) Rank() int {
	return len(iri.seq)
}

/*

Seq Returns CURIE segments
  a:b/c/d/e ⟼ [a, b, c, d]
*/
func (iri IRI) Seq() []string {
	return iri.seq
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

Path converts CURIE to relative file system path.
*/
func (iri IRI) Path() string {
	return path.Join(iri.seq...)
}

/*

URI converts CURIE to fully qualified URL
  wikipedia:CURIE ⟼ http://en.wikipedia.org/wiki/CURIE
*/
func (iri IRI) URI(prefix string) (*url.URL, error) {
	if iri.IsEmpty() {
		return &url.URL{}, nil
	}
	seq := iri.seq

	if seq[0] == "" {
		return url.Parse(strings.Join(seq[1:], "/"))
	}

	return url.Parse(prefix + strings.Join(seq[1:], "/"))
}

/*

Origin decomposes CURIE and returns its source.
  a:b/c/d/e ⟼¹ a:       a:b/c/d/e ⟼⁻¹ a:b/c/d
  a:b/c/d/e ⟼² a:b      a:b/c/d/e ⟼⁻² a:b/c
  a:b/c/d/e ⟼³ a:b/c    a:b/c/d/e ⟼⁻³ a:b
    ...
  a:b/c/d/e ⟼ⁿ a:       a:b/c/d/e ⟼⁻ⁿ a:
*/
func (iri IRI) Origin(rank ...int) IRI {
	r := 1
	if len(rank) > 0 {
		r = rank[0]
	}

	if r < 0 {
		r = len(iri.seq) + r
	}

	switch {
	case r < 0:
		return IRI{seq: []string{}}
	case iri.Rank() <= r:
		return IRI{seq: append([]string{}, iri.seq...)}
	default:
		return IRI{seq: append([]string{}, iri.seq[:r]...)}
	}
}

/*

Parent decomposes CURIE and return its parent CURIE. It return immediate parent
compact URI by default.
  a:b/c/d/e ⟼¹ a:b/c/d  a:b/c/d/e ⟼⁻¹ a:
  a:b/c/d/e ⟼² a:b/c    a:b/c/d/e ⟼⁻² a:b
  a:b/c/d/e ⟼³ a:b      a:b/c/d/e ⟼⁻³ a:b/c
  ...
  a:b/c/d/e ⟼ⁿ a:       a:b/c/d/e ⟼⁻ⁿ a:b/c/d/e
*/
func (iri IRI) Parent(rank ...int) IRI {
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

Prefix decomposes CURIE and return its prefix CURIE (parent) as string value.
*/
func (iri IRI) Prefix(rank ...int) string {
	r := 1
	if len(rank) > 0 {
		r = rank[0]
	}

	n := iri.Rank() - r
	if n < 0 {
		return ""
	}

	return join(iri.seq[:n])
}

/*

Suffix decomposes CURIE and return its suffix.
  a:b/c/d/e ⟿¹ e          a:b/c/d/e ⟿⁻¹ b/c/d/e
  a:b/c/d/e ⟿² d/e        a:b/c/d/e ⟿⁻² c/d/e
  a:b/c/d/e ⟿³ c/d/e      a:b/c/d/e ⟿⁻³ d/e
  ...
  a:b/c/d/e ⟿ⁿ a:b/c/d/e  a:b/c/d/e ⟿⁻ⁿ e
*/
func (iri IRI) Suffix(rank ...int) string {
	r := 1
	if len(rank) > 0 {
		r = rank[0]
	}

	if r < 0 {
		r = len(iri.seq) + r
	}

	n := iri.Rank() - r
	switch {
	case n >= len(iri.seq):
		return ""
	case n > 0:
		return strings.Join(iri.seq[n:len(iri.seq)], "/")
	default:
		return join(iri.seq)
	}
}

/*

Join composes segments into new descendant CURIE.
  a:b × [c, d, e] ⟼ a:b/c/d/e
*/
func (iri IRI) Join(segments ...string) IRI {
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

	return IRI{seq: seq}
}

/*

Heir composes two CURIEs into new descendant CURIE.
	a:b × c/d/e ⟼ a:b/c/d/e
	a:b × c:d/e ⟼ a:b/c/d/e
*/
func (iri IRI) Heir(other IRI) IRI {
	return IRI{seq: append(append([]string{}, iri.seq...), other.Seq()...)}
}

/*

Eq compares two CURIEs, return true if they are equal.
*/
func (iri IRI) Eq(x IRI) bool {
	if iri.Rank() != x.Rank() {
		return false
	}

	seq := x.Seq()
	for i, v := range iri.seq {
		if seq[i] != v {
			return false
		}
	}

	return true
}

/*

Lt compares two CURIEs, return true if left element is less than supplied one.
*/
func (iri IRI) Lt(x IRI) bool {
	if iri.Rank() != x.Rank() {
		return iri.Rank() < x.Rank()
	}

	seq := x.Seq()
	for i, v := range iri.seq {
		if seq[i] != v {
			return v < seq[i]
		}
	}

	return false
}

/*

MarshalJSON `IRI ⟼ "[prefix:suffix]"`
*/
func (iri IRI) MarshalJSON() ([]byte, error) {
	if iri.Rank() == 0 {
		return json.Marshal("")
	}

	return json.Marshal("[" + iri.String() + "]")
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
