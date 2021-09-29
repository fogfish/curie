# curie

The type `curie` ("Compact URI") defines a generic syntax for expressing URIs by abbreviated literal as defined by the [W3C](https://www.w3.org/TR/2010/NOTE-curie-20101216/). This datatype supports type safe identity within domain driven design. 


[![Documentation](https://pkg.go.dev/badge/github.com/fogfish/curie)](https://pkg.go.dev/github.com/fogfish/curie)
[![Build Status](https://github.com/fogfish/curie/workflows/build/badge.svg)](https://github.com/fogfish/curie/actions/)
[![Git Hub](https://img.shields.io/github/last-commit/fogfish/curie.svg)](https://github.com/fogfish/curie)
[![Coverage Status](https://coveralls.io/repos/github/fogfish/curie/badge.svg?branch=main)](https://coveralls.io/github/fogfish/curie?branch=main)
[![Go Report Card](https://goreportcard.com/badge/github.com/fogfish/curie)](https://goreportcard.com/report/github.com/fogfish/curie)
[![Maintainability](https://api.codeclimate.com/v1/badges/bdad0e2fd29d488217fd/maintainability)](https://codeclimate.com/github/fogfish/curie/maintainability)


## Inspiration 

Linked-Data are used widely by Semantic Web to publish structured data so that it can be interlinked by applications. Internationalized Resource Identifiers (IRIs) are key elements to cross-link data structure and establish global references (pointers) to data elements. These IRIs may be written as relative, absolute or compact IRIs. The `curie` type is just a formal definition of **compact IRI** (superset of XML QNames). 

Another challenge solved by `curie` is a formal mechanism to permit the use of hierarchical extensible name collections and its serialization. All-in-all CURIEs expand to any IRI.


## Getting started

The latest version of the library is available at `main` branch of this repository. All development, including new features and bug fixes, take place on the `main` branch using forking and pull requests as described in contribution guidelines.

```go
import "github.com/fogfish/curie"

//
// creates compacts URI to wiki article about CURIE data type
compact := curie.New("wikipedia:CURIE")

//
// expands compact URI to absolute one
//   ⟿ http://en.wikipedia.org/wiki/CURIE
prefix := map[string]string{"wiki": "http://en.wikipedia.org/wiki/"}
url := curie.URI(prefix, compact)
```

The type specification is available at [go doc](https://pkg.go.dev/github.com/fogfish/curie).


### CURIE format

Compact URI is superset of XML QNames. It is comprised of two components: a prefix and a suffix, separated by `:`. Omit prefix to declare a relative URI; omit suffix to declare namespace only; omit both components to declare empty URI. See [W3C CURIE Syntax 1.0](https://www.w3.org/TR/2010/NOTE-curie-20101216/)

```
safe_curie  :=   '[' curie ']'
curie       :=   [ [ prefix ] ':' ] suffix
prefix      :=   NCName
suffix      :=   NCName [ / suffix ]
```

### CURIE "algebra"

The type defines a simple algebra for manipulating instances of compact URI:

```go
// zero: empty compact URI
z := curie.New("")

// transform: string ⟼ CURIE
a := curie.New(/* ... */)
b := curie.New(/* ... */)

// rank: |CURIE| ⟼ Int
curie.Rank(a)

// unary decompose: CURIE ⟼ string
curie.Prefix(c)
curie.Suffix(c)

// binary ordering: CURIE ≼ CURIE ⟼ bool 
curie.Eq(a, b)
curie.Lt(a, b)


// binary compose: CURIE × CURIE ⟼ CURIE
curie.Join(a, b)
```

### URI compatibility

The datatype is compatible with traditional URIs

```go
// any absolute URIs are parsable to CURIE
compact := curie.New("https://example.com/a/b/c")

// String is an identity function
//   ⟿ https://example.com/a/b/c
compact.String()

//
// expands compact URI to absolute one
//   ⟿ https://example.com/a/b/c
url, err := compact.URI("https://example.com/a/b/c")
```

### Hierarchy

CURIE type is core type to organize hierarchies. An application declares `A ⟼ B` hierarchical relation using paths, prefixes and suffixes. 

```go
root := curie.New("some:a")

// construct 2nd rank curie using one of those functions
rank2 := curie.New("some:a/b")
rank2 := curie.Join(root, "b")

//
// parent and prefix of rank2 node is root
//  ⟿ some:a
curie.Parent(rank2)

//
// suffix of rank2 node is 
//  ⟿ b
curie.Child(rank2)

//
// and so on
curie.New("some:a/b/c/d/e")
```

### Linked-data

Cross-linking of structured data is an essential part of type safe domain driven design. The library helps developers to model relations between data instances using familiar data type:

```go
type Person struct {
  ID      curie.IRI
  Social  *curie.String
  Father  *curie.IRI
  Mother  *curie.IRI
  Friends []curie.IRI
}
```

This example uses CURIE data type. `ID` is only used as primary key, `IRI` is a "pointer" to linked-data. `curie.String` is an alternative approach to defined IRI using safe notation.


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
