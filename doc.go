//
// Copyright (C) 2020 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/fogfish/curie
//

/*

Package curie The type `curie` ("Compact URI") defines a generic syntax for
expressing URIs by abbreviated literal as defined by the W3C.
https://www.w3.org/TR/2010/NOTE-curie-20101216/.

The type supports type safe domain driven design using aspects of hierarchical
linked-data.


Inspiration

Linked-Data are used widely by Semantic Web to publish structured data so that
it can be interlinked by applications. Internationalized Resource Identifiers
(IRIs) are key elements to cross-link data structure and establish references
(pointers) to data elements. These IRIs may be written as relative, absolute or
compact IRIs. The `curie` type is just a formal definition of compact IRI
(superset of XML QNames). Another challenge solved by `curie` is a formal
mechanism to permit the use of hierarchical extensible name collections and its
serialization. All-in-all CURIEs expand to any IRI.


CURIE format

Compact URI is superset of XML QNames. It is comprised of two components:
a prefix and a suffix, separated by `:`. Omit prefix to declare a relative URI;
omit suffix to declare namespace only; omit both components to declare empty
URI. See W3C CURIE Syntax 1.0 https://www.w3.org/TR/2010/NOTE-curie-20101216/

  safe_curie  :=   '[' curie ']'
  curie       :=   [ [ prefix ] ':' ] reference
  prefix      :=   NCName
  reference   :=   irelative-ref (as defined in IRI)


CURIE algebra

The type defines a simple algebra for manipulating instances of compact URI

↣ zero: empty compact URI

↣ transform: string ⟼ CURIE

↣ binary compose: CURIE × CURIE ⟼ CURIE

↣ unary decompose: CURIE ⟼ CURIE

↣ rank: |CURIE| ⟼ Int

↣ binary ordering: CURIE ≼ CURIE ⟼ bool


Linked data

Cross-linking of structured data is an essential part of type safe domain
driven design. The library helps developers to model relations between data
instances using familiar data type:

  type Person struct {
    curie.ID
    Father  *curie.IRI
    Mother  *curie.IRI
    Friends []curie.IRI
  }

`curie.ID` and `curie.IRI` are sibling, equivalent CURIE data type.
`ID` is only used as primary key, `IRI` is a "pointer" to linked-data.

CURIE type is core type to organize hierarchies. An application declares
`A ⟼ B` hierarchical relation using path at suffix. For example, the root is
`curie.New("some:a")`, 2nd rank node `curie.New("some:a/b")` and so on
`curie.New("some:a/b/c/e/f")`.

*/
package curie
