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

/*

Thing is the most generic type of item. The interfaces declares anything with
unique identifier. Embedding CURIE ID into struct makes it Thing compatible.
*/
type Thing interface {
	// The identifier property represents any kind of identifier for
	// any kind of Thing
	Identity() ID
}

/*

IRI ...
*/
type IRI interface {
	IsEmpty() bool
	Rank() int
	Seq() []string
	String() string
	Safe() string
	Parent(rank ...int) IRI
	Prefix(rank ...int) string
	Suffix(rank ...int) string
	Join(segments ...string) IRI
	Heir(other IRI) IRI
	Eq(x IRI) bool
	Lt(x IRI) bool
}

//------------------------------------------------------------------------------
//
// ID
//
//------------------------------------------------------------------------------

/*

ID is compact URI (CURIE) type for struct tagging, It declares unique identity
of a thing. The tagged struct belongs to Thing category (implements Thing interface)

  type MyStruct struct {
    curie.ID
  }

*/
type ID struct {
	IRI `dynamodbav:"id" json:"id"`
}

/*

New transform category of strings to CURIE.
*/
func New(iri string, args ...interface{}) ID {
	return ID{IRI: NewIRI(iri, args...)}
}

/*

Identity makes CURIE compliant to Thing interface so that embedding ID makes any
struct to be Thing.
*/
func (iri ID) Identity() ID {
	return iri
}

/*

IsEmpty returns true if compact URI is not defined (empty)
*/
func (iri ID) IsEmpty() bool {
	return iri.IRI.IsEmpty()
}

/*

Rank value of CURIE
*/
func (iri ID) Rank() int {
	return iri.Rank()
}

/*

Safe transforms CURIE to safe string
*/
func (iri ID) Safe() string {
	return iri.IRI.Safe()
}

/*

Parent decomposes CURIE and return its parent CURIE. It return immediate parent
compact URI by default. Use optional rank param to extract "grant" parents,
non immediate value distant at rank.
	a:b/c/d/e ⟼¹ a:b/c/d
	a:b/c/d/e ⟼² a:b/c
	a:b/c/d/e ⟼³ a:b
*/
func (iri ID) Parent(rank ...int) ID {
	return ID{IRI: iri.IRI.Parent(rank...)}
}

/*

Prefix decomposes CURIE and return its prefix CURIE (parent) as string value.
*/
func (iri ID) Prefix(rank ...int) string {
	return iri.IRI.Prefix(rank...)
}

/*

Suffix decomposes CURIE and return its suffix.
	a:b/c/d/e ⟿¹ e
	a:b/c/d/e ⟿² d/e
	a:b/c/d/e ⟿³ c/d/e
	...
	a:b/c/d/e ⟿ⁿ a:b/c/d/e
*/
func (iri ID) Suffix(rank ...int) string {
	return iri.IRI.Suffix(rank...)
}

/*

Join composes segments into new descendant CURIE.
  a:b × [c, d, e] ⟼ a:b/c/d/e
*/
func (iri ID) Join(segments ...string) ID {
	return ID{IRI: iri.IRI.Join(segments...)}
}

/*

Heir composes two CURIEs into new descendant CURIE.
	a:b × c/d/e ⟼ a:b/c/d/e
	a:b × c:d/e ⟼ a:b/c/d/e
*/
func (iri ID) Heir(other ID) ID {
	return ID{IRI: iri.IRI.Heir(other.IRI)}
}

/*

URI converts CURIE to fully qualified URL
  wikipedia:CURIE ⟼ http://en.wikipedia.org/wiki/CURIE
*/
func (iri ID) URI(prefix string) (*url.URL, error) {
	if iri.IRI.IsEmpty() {
		return &url.URL{}, nil
	}
	seq := iri.IRI.Seq()

	if seq[0] == "" {
		return url.Parse(strings.Join(seq[1:], "/"))
	}

	return url.Parse(prefix + strings.Join(seq[1:], "/"))
}

/*

Path converts CURIE to relative file system path.
*/
func (iri ID) Path() string {
	return path.Join(iri.IRI.Seq()...)
}

/*

Seq converts CURIE to sequence of segments
*/
func (iri ID) Seq() []string {
	return iri.IRI.Seq()
}

/*

Eq compares two CURIEs, return true if they are equal.
*/
func (iri ID) Eq(x ID) bool {
	return iri.IRI.Eq(x.IRI)
}

/*

Lt compares two CURIEs, return true if left element is less than supplied one.
*/
func (iri ID) Lt(x ID) bool {
	return iri.IRI.Lt(x.IRI)
}

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
type CURIE struct {
	seq []string
}

/*

NewIRI transform category of strings to CURIE.
TODO: rename Parse
*/
func NewIRI(iri string, args ...interface{}) CURIE {
	val := iri

	if len(args) > 0 {
		val = fmt.Sprintf(iri, args...)
	}

	return CURIE{
		seq: split(strings.Trim(val, "[]")),
	}
}

func (iri CURIE) IsEmpty() bool {
	return len(iri.seq) == 0
}

func (iri CURIE) Rank() int {
	return len(iri.seq)
}

func (iri CURIE) Seq() []string {
	return iri.seq
}

/*

This return pointer to itself
*/
func (iri CURIE) This() IRI {
	return &iri
}

/*

String transform CURIE to string
*/
func (iri CURIE) String() string {
	return join(iri.seq)
}

/*

Safe transforms CURIE to safe string
*/
func (iri CURIE) Safe() string {
	return "[" + iri.String() + "]"
}

/*

Parent decomposes CURIE and return its parent CURIE. It return immediate parent
compact URI by default. Use optional rank param to extract "grant" parents,
non immediate value distant at rank.
	a:b/c/d/e ⟼¹ a:b/c/d
	a:b/c/d/e ⟼² a:b/c
	a:b/c/d/e ⟼³ a:b
*/
func (iri CURIE) Parent(rank ...int) IRI {
	r := 1
	if len(rank) > 0 {
		r = rank[0]
	}

	n := len(iri.seq) - r
	if n < 0 {
		return CURIE{seq: []string{}}
	}

	if n == 1 && iri.seq[0] == "" {
		return CURIE{seq: []string{}}
	}

	return CURIE{seq: append([]string{}, iri.seq[:n]...)}
}

/*

Prefix decomposes CURIE and return its prefix CURIE (parent) as string value.
*/
func (iri CURIE) Prefix(rank ...int) string {
	r := 1
	if len(rank) > 0 {
		r = rank[0]
	}

	n := len(iri.seq) - r
	if n < 0 {
		return ""
	}

	return join(iri.seq[:n])
}

/*

Suffix decomposes CURIE and return its suffix.
	a:b/c/d/e ⟿¹ e
	a:b/c/d/e ⟿² d/e
	a:b/c/d/e ⟿³ c/d/e
	...
	a:b/c/d/e ⟿ⁿ a:b/c/d/e
*/
func (iri CURIE) Suffix(rank ...int) string {
	r := 1
	if len(rank) > 0 {
		r = rank[0]
	}

	n := len(iri.seq) - r
	if n > 0 {
		return strings.Join(iri.seq[n:len(iri.seq)], "/")
	}

	return join(iri.seq)
}

/*

Join composes segments into new descendant CURIE.
  a:b × [c, d, e] ⟼ a:b/c/d/e
*/
func (iri CURIE) Join(segments ...string) IRI {
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

	return CURIE{seq: seq}
}

/*

Heir composes two CURIEs into new descendant CURIE.
	a:b × c/d/e ⟼ a:b/c/d/e
	a:b × c:d/e ⟼ a:b/c/d/e
*/
func (iri CURIE) Heir(other IRI) IRI {
	return CURIE{seq: append(append([]string{}, iri.seq...), other.Seq()...)}
}

/*

Eq compares two CURIEs, return true if they are equal.
*/
func (iri CURIE) Eq(x IRI) bool {
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
func (iri CURIE) Lt(x IRI) bool {
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
func (iri CURIE) MarshalJSON() ([]byte, error) {
	if len(iri.seq) == 0 {
		return json.Marshal("")
	}

	return json.Marshal("[" + iri.String() + "]")
}

/*

UnmarshalJSON `"[prefix:suffix]" ⟼ IRI`
*/
func (iri *CURIE) UnmarshalJSON(b []byte) error {
	var path string
	err := json.Unmarshal(b, &path)
	if err != nil {
		return err
	}

	*iri = NewIRI(path)
	return nil
}

/*

MarshalDynamoDBAttributeValue `IRI ⟼ "prefix:suffix"`
*/
// func (iri IRI) MarshalDynamoDBAttributeValue(av *dynamodb.AttributeValue) error {
// 	if len(iri.Seq) == 0 {
// 		av.NULL = aws.Bool(true)
// 		return nil
// 	}

// 	// Note: we are using string representation to allow linked data in dynamo tables
// 	val, err := dynamodbattribute.Marshal(iri.String())
// 	if err != nil {
// 		return err
// 	}

// 	av.S = val.S
// 	return nil
// }

/*

UnmarshalDynamoDBAttributeValue `"prefix:suffix" ⟼ IRI`
*/
// func (iri *IRI) UnmarshalDynamoDBAttributeValue(av *dynamodb.AttributeValue) error {
// 	*iri = NewIRI(aws.StringValue(av.S))
// 	return nil
// }

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
