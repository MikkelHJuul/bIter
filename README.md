# bIter

[![Go Report Card](https://goreportcard.com/badge/github.com/MikkelHJuul/bIter)](https://goreportcard.com/report/github.com/MikkelHJuul/bIter)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/MikkelHJuul/bIter)](https://pkg.go.dev/github.com/MikkelHJuul/bIter)
[![Maintainability](https://api.codeclimate.com/v1/badges/96eda04e3207e18786c2/maintainability)](https://codeclimate.com/github/MikkelHJuul/bIter/maintainability)
[![codecov](https://codecov.io/gh/MikkelHJuul/bIter/branch/main/graph/badge.svg?token=1RFY7XASKC)](https://codecov.io/gh/MikkelHJuul/bIter)
![GitHub License](https://img.shields.io/github/license/MikkelHJuul/bIter)

Iterator for badgerDB

```go
// Iterator adds an implementable target for variations of different Iterator's
// for simplification of functional code, that you can then implement this reduced
// interface such that primarily methods Rewind and Valid and Next are overloaded.
// usage:
// 		for iterator.Rewind(); iterator.Valid(); iterator.Next() {
//			iterator.Item()
//			...
//		}
// this is a general snippet of code that using this interface may have overloaded
// only Rewind and Valid are changes
type Iterator interface {
	Rewind()
	Valid() bool
	Next()
	Item() *badger.Item
}
```
This library also add 5 concrete iterators: prefix, from-to, from, to, prefix-from.

The library does not support backward key-scan (whence only a prefix-from)

Use the method `biter.KeyRangeIterator` and throw any combination of prefix, from and to at it to get the corresponding `Iterator`.

for some examples of functionality refer to [queries_test.go](queries_test.go)