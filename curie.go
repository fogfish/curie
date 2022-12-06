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
	"strings"
)

//------------------------------------------------------------------------------
//
// IRI
//
//------------------------------------------------------------------------------

// IRI is compact URI, defined as superset of XML QNames, with the modification
// that the format of the strings after the colon is looser.
//
// safe_curie  :=   '[' curie ']'
// curie       :=   [ [ prefix ] ':' ] reference
// prefix      :=   NCName
// reference   :=   irelative-ref (as defined in IRI, RFC 3987)
type IRI string

// Safe transforms CURIE to safe string
func (iri IRI) Safe() string {
	if len(iri) == 0 {
		return ""
	}

	return "[" + string(iri) + "]"
}

// MarshalJSON `IRI ⟼ "[prefix:suffix]"`
func (iri IRI) MarshalJSON() ([]byte, error) {
	if len(iri) == 0 {
		return json.Marshal("")
	}

	return json.Marshal(iri.Safe())
}

// UnmarshalJSON `"[prefix:suffix]" ⟼ IRI`
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
// Prefixes
//
//------------------------------------------------------------------------------

// Prefixes is a collection of prefixes defined by the application
type Prefixes interface {
	Create(string) IRI
	Lookup(string) (string, bool)
}

// Namespaces is constant in-memory collection of prefixes defined by the application
type Namespaces map[string]string

// Create new URI using prefix table
func (ns Namespaces) Create(uri string) IRI {
	// Note: All non-ASCII code points in the IRI should next be encoded as UTF-8
	// https://en.wikipedia.org/wiki/Internationalized_Resource_Identifier
	// https://www.ietf.org/rfc/rfc3987.html#section-5.3.2.3
	for key, val := range ns {
		if strings.HasPrefix(uri, val) {
			ref := Decode(uri[len(val):])
			return IRI(key + ":" + string(ref))
		}
	}

	return IRI(uri)
}

// Lookup prefix in the map
func (ns Namespaces) Lookup(prefix string) (string, bool) {
	val, exists := ns[prefix]
	return val, exists
}

//------------------------------------------------------------------------------
//
// IRI algebra
//
//------------------------------------------------------------------------------

// New transform category of strings to IRI.
// It expects UTF-8 string according to RFC 3987.
func New(iri string, args ...interface{}) IRI {
	val := iri
	if len(val) > 0 && (val[0] == '[' && val[len(val)-1] == ']') {
		val = val[1 : len(val)-1]
	}

	if len(args) > 0 {
		val = fmt.Sprintf(iri, args...)
	}

	return IRI(val)
}

// IsEmpty is an alias to curie.Rank(iri) == 0
func IsEmpty(iri IRI) bool {
	return len(iri) == 0
}

// Built-in CURIE ranks
const (
	EMPTY = iota
	PREFIX
	REFERENCE
)

// Rank of CURIE, number of elements
// Rank is an alias of len(curie.Seq(iri))
func Rank(iri IRI) int {
	if len(iri) == 0 {
		return EMPTY
	}

	n := strings.IndexRune(string(iri), ':')
	if n == len(iri)-1 {
		return PREFIX
	}

	return REFERENCE
}

// Seq Returns CURIE segments
//
//	a: ⟼ [ a ]
//	b ⟼ [ , b ]
//	a:b ⟼ [a, b]
//	a:b/c/d ⟼ [a, b/c/d ]
func Seq(iri IRI) (string, string) {
	if len(iri) == 0 {
		return "", ""
	}

	n := strings.IndexRune(string(iri), ':')
	if n == -1 {
		return "", string(iri)
	}

	if n == len(iri)-1 {
		return string(iri)[:n], ""
	}

	return string(iri)[:n], string(iri)[n+1:]
}

// Prefix decomposes CURIE and return its prefix CURIE as string value.
func Prefix(iri IRI) string {
	if len(iri) == 0 {
		return ""
	}

	n := strings.IndexRune(string(iri), ':')
	if n == -1 {
		return ""
	}

	return string(iri)[:n]
}

// Reference decomposes CURIE and return its reference as string value.
func Reference(iri IRI) string {
	if len(iri) == 0 {
		return ""
	}

	n := strings.IndexRune(string(iri), ':')
	if n == -1 {
		return string(iri)
	}

	return string(iri)[n+1:]
}

// URI converts CURIE to fully qualified URL
//
//	wikipedia:CURIE ⟼ http://en.wikipedia.org/wiki/CURIE
func URI(prefixes Prefixes, iri IRI) string {
	uri, err := URL(prefixes, iri)
	if err != nil {
		return string(iri)
	}

	return uri.String()
}

// URI converts fully qualified URL to CURIE
//
//	http://en.wikipedia.org/wiki/CURIE ⟼ wikipedia:CURIE
func FromURI(prefixes Prefixes, uri string) IRI {
	return prefixes.Create(uri)
}

// URL converts CURIE to fully qualified url.URL type
//
//	wikipedia:CURIE ⟼ http://en.wikipedia.org/wiki/CURIE
func URL(prefixes Prefixes, iri IRI) (*url.URL, error) {
	if len(iri) == 0 {
		return new(url.URL), nil
	}

	//
	// A host language MAY declare a default prefix value, or
	// MAY provide a mechanism for defining a default prefix value.
	// In such a host language, when the prefix is omitted from a CURIE,
	// the default prefix value MUST be used.
	//
	uri := string(iri)
	if prefix, exists := prefixes.Lookup(Prefix(iri)); exists {
		uri = prefix + Reference(iri)
	}

	return url.Parse(uri)
}

// Eq compares two CURIEs, return true if they are equal.
func Eq(a, b IRI) bool {
	return a == b
}

// Lt compares two CURIEs, return true if left element is less than supplied one.
func Lt(a, b IRI) bool {
	if Rank(a) != Rank(b) {
		return Rank(a) < Rank(b)
	}

	return a < b
}

// Join composes segments into new descendant CURIE.
//
//	a:b × [c, d, e] ⟼ a:b/c/d/e
func Join(iri IRI, segments ...string) IRI {
	if len(segments) == 0 {
		return iri
	}

	// custom strings.Join that ignores empty segments
	n := (len(segments) - 1)
	for i := 0; i < len(segments); i++ {
		n += len(segments[i])
	}

	var b strings.Builder
	b.Grow(n)
	b.WriteString(segments[0])
	for _, s := range segments[1:] {
		if s != "" {
			b.WriteString("/")
			b.WriteString(s)
		}
	}
	path := b.String()

	if len(path) == 0 {
		return iri
	}

	//
	switch Rank(iri) {
	case EMPTY:
		return IRI(path)
	case PREFIX:
		return IRI(string(iri) + path)
	default:
		return IRI(string(iri) + "/" + path)
	}
}
