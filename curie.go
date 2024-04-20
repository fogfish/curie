//
// Copyright (C) 2020 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/fogfish/curie
//

// Package curie implements the type for compact URI. It defines a generic syntax
// for expressing URIs by abbreviated literal as defined by the W3C.
// https://www.w3.org/TR/2010/NOTE-curie-20101216/
package curie

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"github.com/fogfish/curie/v2/internal/reference"
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

// Empty IRI
const Empty = IRI("")

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
	var val string
	err := json.Unmarshal(b, &val)
	if err != nil {
		return err
	}

	if len(val) == 0 {
		*iri = IRI("")
		return nil
	}

	if val[0] != '[' && val[len(val)-1] != ']' {
		return fmt.Errorf("invalid CURIE %s", val)
	}

	val = val[1 : len(val)-1]

	*iri = IRI(val)
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
	if prefix, exists := prefixes.Lookup(Schema(iri)); exists {
		uri = prefix + Reference(iri)
	}

	return url.Parse(uri)
}

//------------------------------------------------------------------------------
//
// IRI algebra
//
//------------------------------------------------------------------------------

// New transform category of strings to IRI.
// It expects UTF-8 string according to RFC 3987.
func New(schema, ref string, args ...interface{}) IRI {
	curie := schema
	if len(curie) > 0 && !strings.HasSuffix(curie, ":") {
		curie += ":"
	}

	if len(ref) > 0 {
		curie += fmt.Sprintf(ref, args...)
	}

	return IRI(curie)
}

// Split URN into NID and NSS
func (iri IRI) Split() (string, string) { return Split(iri) }

// Split CURIE
//
//	a: ⟼ [ a ]
//	b ⟼ [ , b ]
//	a:b ⟼ [a, b]
//	a:b/c/d ⟼ [a, b/c/d ]
func Split(iri IRI) (string, string) {
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

// Return CURIE Schema (prefix)
func (iri IRI) Schema() string { return Schema(iri) }

// Return CURIE Schema (prefix)
func Schema(iri IRI) string {
	schema, _ := Split(iri)
	return schema
}

// Return CURIE Reference
func (iri IRI) Reference() string { return Reference(iri) }

// Return URN Reference
func Reference(iri IRI) string {
	_, ref := Split(iri)
	return ref
}

// IsEmpty is an alias to len(iri) == 0
func (iri IRI) IsEmpty() bool { return len(iri) == 0 }

// IsEmpty is an alias to len(iri) == 0
func IsEmpty(iri IRI) bool {
	return len(iri) == 0
}

// Join composes segments into new descendant CURIE.
func (iri IRI) Join(segments ...string) IRI { return Join(iri, segments...) }

// Join composes segments into new descendant URN.
//
// a:b × [c, d, e] ⟼ a:b/c/d/e
func Join(iri IRI, segments ...string) IRI {
	schema, ref := Split(iri)
	return New(schema, reference.Join(ref, '/', segments...))
}

// Disjoin decomposes CURIE.
func (iri IRI) Disjoin(n int) IRI { return Disjoin(iri, n) }

// Disjoin decomposes CURIE.
func Disjoin(iri IRI, n int) IRI {
	schema, ref := Split(iri)
	return New(schema, reference.Split(ref, '/', n))
}
