//
// Copyright (C) 2020 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/fogfish/curie
//

// Package urn implements the type for URN. It is defined URN as identity aspect only.
// https://www.rfc-editor.org/rfc/rfc8141
package urn

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/fogfish/curie/v2"
	"github.com/fogfish/curie/v2/internal/reference"
)

// URN is a Uniform Resource Identifier (URI) that is assigned under the
// "urn" URI scheme and a particular URN namespace, with the intent that
// the URN will be a persistent, location-independent resource identifier.
//
//	namestring   = "urn" ":" NID ":" NSS
//	NID          = (alphanum) 0*30(ldh) (alphanum)
//	ldh          = alphanum / "-"
//	NSS          = pchar *(pchar / "/" / ":")
type URN string

// Empty URN
const Empty = URN("")

// Create new URN
func New(schema, ref string) URN {
	urn := "urn:" + schema
	if len(ref) > 0 {
		if len(ref) > 0 {
			urn += ":" + ref
		}
	}

	return URN(urn)
}

// MarshalJSON `URN ⟼ "urn:schema:reference"`
func (urn URN) MarshalJSON() ([]byte, error) {
	if len(urn) == 0 || (len(urn) > 5 && strings.HasPrefix(string(urn), "urn:")) {
		return json.Marshal(string(urn))
	}

	return nil, fmt.Errorf("invalid URN %s", urn)
}

// UnmarshalJSON `"urn:schema:reference" ⟼ URN`
func (urn *URN) UnmarshalJSON(b []byte) error {
	var val string
	err := json.Unmarshal(b, &val)
	if err != nil {
		return err
	}

	if len(val) == 0 || (len(val) > 5 && strings.HasPrefix(val, "urn:")) {
		*urn = URN(val)
		return nil
	}

	return fmt.Errorf("invalid URN %s", val)
}

// Return URN Schema
func (urn URN) Schema() string { return Schema(urn) }

// Return URN Schema
func Schema(urn URN) string {
	schema, _ := Split(urn)
	return schema
}

// Return URN Reference
func (urn URN) Reference() string { return Reference(urn) }

// Return URN Reference
func Reference(urn URN) string {
	_, ref := Split(urn)
	return ref
}

// Split URN into NID and NSS
func (urn URN) Split() (string, string) { return Split(urn) }

// Split URN into NID and NSS
func Split(urn URN) (string, string) {
	if len(urn) < 5 {
		return "", ""
	}

	s := urn[4:]
	n := strings.IndexRune(string(s), ':')

	if n == -1 {
		return string(s), ""
	}

	return string(s[:n]), string(s[n+1:])
}

// Base returns the last element of CURIE reference
func Base(iri URN) string {
	ref := Reference(iri)

	if len(ref) == 0 {
		return ""
	}

	n := strings.LastIndexByte(ref, ':')
	if n == -1 {
		return ref
	}

	return ref[n+1:]
}

// Path returns all but the last element of CURIE reference
func Path(iri URN) URN {
	schema, ref := Split(iri)
	if len(ref) == 0 {
		return iri
	}

	n := strings.LastIndexByte(ref, ':')
	if n == -1 {
		ref = ""
	} else {
		ref = ref[:n]
	}

	return New(schema, ref)
}

// Head returns the head element of CURIE reference
func Head(iri URN) string {
	ref := Reference(iri)

	if len(ref) == 0 {
		return ""
	}

	n := strings.IndexRune(string(ref), ':')
	if n == -1 {
		return ref
	}

	return ref[:n]
}

// Path returns all but the fiirst element of CURIE reference
func Tail(iri URN) URN {
	schema, ref := Split(iri)
	if len(ref) == 0 {
		return iri
	}

	n := strings.IndexRune(string(ref), ':')
	if n == -1 {
		ref = ""
	} else {
		ref = ref[n+1:]
	}

	return New(schema, ref)
}

// Join composes segments into new descendant URN.
func (urn URN) Join(segments ...string) URN { return Join(urn, segments...) }

// Join composes segments into new descendant URN.
//
// urn:a:b:c × [d, e, f] ⟼ a:b:c:d:e:f
func Join(urn URN, segments ...string) URN {
	schema, ref := Split(urn)
	return New(schema, reference.Join(ref, ':', segments...))
}

// Cut N components from URN NSS
func (urn URN) Cut(n int) URN { return Cut(urn, n) }

// Cut N components from URN NSS
func Cut(urn URN, n int) URN {
	schema, ref := Split(urn)
	return New(schema, reference.Split(ref, ':', n))
}

// Conver URN to IRI
func ToIRI(urn URN) curie.IRI {
	schema, ref := Split(urn)
	ref = filepath.Join(strings.Split(ref, ":")...)
	return curie.New(schema, ref)
}

func ToURN(iri curie.IRI) URN {
	schema, ref := curie.Split(iri)
	ref = strings.Join(strings.Split(ref, "/"), ":")
	return New(schema, ref)
}
