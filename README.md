# bIter
Iterator for badgerDB

```
// Iterator adds an implementable target for variations of different Iterator's
// for simplification of functional code, that you can then implement this reduced
// interface such that primarily methods Rewind and Valid and be overloaded.
// usage:
// 		for iterator.Rewind(); iterator.Valid(); iterator.Next() {
//			iterator.Item()
//			...
//		}
// is a general snippet of code that using this interface may have overloaded
// ... Rewind to seek to a prefix or a value
// ... Valid to validate the key is still within bounds using fx bytes.Compare
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
