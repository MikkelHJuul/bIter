# bIter
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

Use the method `bIter.KeyRangeIterator` and throw any combination of prefix, from and to at it to get the corresponding `Iterator`.
