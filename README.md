# curie

The library defines identity types `curie` ("Compact URI") as defined by the [W3C](https://www.w3.org/TR/2010/NOTE-curie-20101216/) and `urn` as defined by [RFC8141](https://www.rfc-editor.org/rfc/rfc8141). These datatype supports type safe identity within domain driven design.

[![Version](https://img.shields.io/github/v/tag/fogfish/curie?label=version)](https://github.com/fogfish/curie/releases)
[![Documentation](https://pkg.go.dev/badge/github.com/fogfish/curie)](https://pkg.go.dev/github.com/fogfish/curie)
[![Build Status](https://github.com/fogfish/curie/workflows/build/badge.svg)](https://github.com/fogfish/curie/actions/)
[![Git Hub](https://img.shields.io/github/last-commit/fogfish/curie.svg)](https://github.com/fogfish/curie)
[![Coverage Status](https://coveralls.io/repos/github/fogfish/curie/badge.svg?branch=main)](https://coveralls.io/github/fogfish/curie?branch=main)
[![Go Report Card](https://goreportcard.com/badge/github.com/fogfish/curie)](https://goreportcard.com/report/github.com/fogfish/curie)


## Inspiration 

Linked-Data are used widely by Semantic Web to publish structured data so that it can be interlinked by applications. Internationalized Resource Identifiers (IRIs) are key elements to cross-link data structure and establish global references (pointers) to data elements. These IRIs may be written as relative, absolute or compact IRIs. The `curie` type is just a formal definition of **compact IRI** (superset of XML QNames, with the modification that the format of the strings after the colon is looser). 

Another challenge solved by `curie` is a formal mechanism to permit the use of the concept of scoping, where identities are created within a unique scope, and that scope's collection is managed by the group that defines it. All-in-all CURIEs expand to any IRI.

`urn` is an alternative for `curie` following same principles.

Both types implements coherent interface allowing application establish hierarchical idenity. 

- [curie](#curie)
  - [Inspiration](#inspiration)
  - [Getting started](#getting-started)
    - [Quick Example](#quick-example)
    - [Identity formats](#identity-formats)
      - [CURIE](#curie-1)
      - [URN](#urn)
    - [Interface](#interface)
    - [URI compatibility](#uri-compatibility)
    - [Linked-data](#linked-data)
  - [How To Contribute](#how-to-contribute)
    - [commit message](#commit-message)
    - [bugs](#bugs)
  - [License](#license)


## Getting started

The latest version of the library is available at `main` branch of this repository. All development, including new features and bug fixes, take place on the `main` branch using forking and pull requests as described in contribution guidelines. The stable version is available via Golang modules.

```bash
go get github.com/fogfish/curie/v2
```

### Quick Example

```go
package main

import (
	"fmt"

	"github.com/fogfish/curie/v2"
)

func main() {
	//
	// creates compacts URI to wiki article about CURIE data type
	iri := curie.New("wiki", "CURIE")
	fmt.Println(iri.Safe())

	//
	// expands compact URI to absolute one
	//   ⟿ http://en.wikipedia.org/wiki/CURIE
	prefixes := curie.Namespaces{
		"wiki": "http://en.wikipedia.org/wiki/",
	}
	url := curie.URI(prefixes, iri)
	fmt.Println(url)
}
```

### Identity formats

The library generalizes  identity as composition of `schema` and `reference`. `schema` defines the namespace and `reference` is a hierarchical identifier.

#### CURIE

Compact URI is superset of XML QNames, with the modification that the format of the strings after the colon is looser. It is comprised of two components: a prefix and a reference, separated by `:`. Omit prefix to declare a relative IRI; omit suffix to declare namespace only. See [W3C CURIE Syntax 1.0](https://www.w3.org/TR/2010/NOTE-curie-20101216/)

```
safe_curie  :=   '[' curie ']'
curie       :=   [ [ prefix ] ':' ] reference
prefix      :=   NCName
reference   :=   irelative-ref (as defined in IRI, RFC 3987)
```

#### URN

A Uniform Resource Name (URN) is a Uniform Resource Identifier that is assigned under the "urn" URI scheme and a particular URN namespace.

```
URN := "urn" ":" NID ":" NSS
NID := (alphanum) 0*30(ldh) (alphanum)
ldh := alphanum / "-"
NSS := pchar *(pchar / "/")
```

### Interface

The library provides packages (`curie`, `urn`). Each implements the coherent api.

```go
// Unconditionally cast string to type
curie.IRI("wiki:CURIE")

// Creates IRI definiting schema & reference
iri := curie.New("wiki", "CURIE")

// Get schema & reference from identity
schema := curie.Schema(iri)
ref := curie.Reference(iri)
schema, ref := curie.Split(iri)
```

Hierarchical identity is a key feature. 

```go
// Composes "hierarchical" identity, adding "segments" to reference
// ⟿ wiki:CURIE/Example
iri = curie.Join(iri, "Example")

// Read base component of identity
// ⟿ Example
curie.Base(iri)

// Read all identity but base component
// ⟿ wiki:CURIE
curie.Path(iri)

// Read head component of identity
// ⟿ CURIE
curie.Head(iri)

// Read all identity but head component
// ⟿ wiki:Example
curie.Tail(iri)

// Discard last N segments
// ⟿ wiki:CURIE
curie.Cut(iri, 1)
```

See [go doc](https://pkg.go.dev/github.com/fogfish/curie/v2).


### URI compatibility

The datatype is compatible with traditional URIs

```go
// absolute URIs are parseable to curie if supplied with namespaces
prefixes := curie.Namespaces{
"ex": "https://example.com/",
}

compact := curie.FromURI(prefixes,"https://example.com/a/b/c")

fmt.Println(compact)

// expands compact URI to absolute one
//   ⟿ https://example.com/a/b/c
uri := curie.URI(prefixes,"ex:a/b/c")
fmt.Println(uri)
```

### Linked-data

Cross-linking of structured data is an essential part of type safe domain driven design. The library helps developers to model relations between data instances using familiar data type:

```go
type Person struct {
  ID      curie.IRI
  Social  *curie.IRI
  Father  *curie.IRI
  Mother  *curie.IRI
  Friends []curie.IRI
}
```

This example uses CURIE data type. `ID` is a primary key, all other `IRI` is a "pointer" to linked-data.


## How To Contribute

The library is [MIT](LICENSE) licensed and accepts contributions via GitHub pull requests:

1. Fork it
2. Create your feature branch (`git checkout -b my-new-feature`)
3. Commit your changes (`git commit -am 'Added some feature'`)
4. Push to the branch (`git push origin my-new-feature`)
5. Create new Pull Request

The build and testing process requires [Go](https://golang.org) version 1.13 or later.

**build** and **test** library.

```bash
git clone https://github.com/fogfish/curie
cd curie
go test
```

### commit message

The commit message helps us to write a good release note, speed-up review process. The message should address two question what changed and why. The project follows the template defined by chapter [Contributing to a Project](http://git-scm.com/book/ch5-2.html) of Git book.

### bugs

If you experience any issues with the library, please let us know via [GitHub issues](https://github.com/fogfish/curie/issue). We appreciate detailed and accurate reports that help us to identity and replicate the issue. 


## License

[![See LICENSE](https://img.shields.io/github/license/fogfish/curie.svg?style=for-the-badge)](LICENSE)
