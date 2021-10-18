package biter

import (
	"bytes"
	"github.com/dgraph-io/badger/v3"
)

// Iterator adds an implementable target for variations of different Iterator(s)
// for simplification of functional code, that you can then implement this reduced
// interface such that primarily methods Rewind and Valid and Next are overloaded.
// usage:
// 		for iterator.Rewind(); iterator.Valid(); iterator.Next() {
//			iterator.Item()
//			...
//		}
// this is a general snippet of code that using this interface may have changed implementation
// only Rewind and Valid are changes
type Iterator interface {
	Rewind()
	Valid() bool
	Next()
	Item() *badger.Item
}

//badgerPrefixIterator seek's to the prefix and ends at the end of the prefix
type badgerPrefixIterator struct {
	*badger.Iterator
	prefix []byte
}

func (b *badgerPrefixIterator) Rewind() {
	b.Seek(b.prefix)
}

func (b *badgerPrefixIterator) Valid() bool {
	return b.ValidForPrefix(b.prefix)
}

//badgerFromIterator seeks to from, then keeps going as usual
type badgerFromIterator struct {
	*badger.Iterator
	from []byte
}

func (b *badgerFromIterator) Valid() bool {
	return b.Iterator.Valid()
}

func (b *badgerFromIterator) Rewind() {
	b.Seek(b.from)
}

//badgerToIterator starts at the top, and keeps going until to is less than the current key
type badgerToIterator struct {
	*badger.Iterator
	to []byte
}

func (b *badgerToIterator) Valid() bool {
	return b.Iterator.Valid() && keyLtEq(b.Item().Key(), b.to)
}

//badgerFromToIterator seeks to from, and stops at to
type badgerFromToIterator struct {
	*badger.Iterator
	from, to []byte
}

func (b *badgerFromToIterator) Valid() bool {
	return b.Iterator.Valid() && keyLtEq(b.Item().Key(), b.to)
}

func (b *badgerFromToIterator) Rewind() {
	b.Seek(b.from)
}

//badgerPrefixFromIterator is like a badgerPrefixIterator,
//but it rewinds to from in stead
type badgerPrefixFromIterator struct {
	*badgerPrefixIterator
	from []byte
}

func (b *badgerPrefixFromIterator) Rewind() {
	b.Seek(b.from)
}

//badgerPrefixToIterator does not exist as the method set
//is completely similar to that of badgerFromToIterator

func keyLtEq(key, upperBound []byte) bool {
	return 1 > bytes.Compare(key, upperBound)
}

// KeyRangeIterator returns an Iterator interface implementation, that wraps the badger.Iterator
// in order to simplify iterating with from-to and/or prefix.
func KeyRangeIterator(it *badger.Iterator, prefix, from, to []byte) Iterator {
	lPref, lFrom, lTo := len(prefix), len(from), len(to)
	switch {
	case lPref+lFrom+lTo == 0:
		return it
	case lPref != 0 && lFrom+lTo == 0:
		return &badgerPrefixIterator{it, prefix}
	case lTo != 0 && lPref+lFrom == 0:
		return &badgerToIterator{it, to}
	case lFrom != 0 && lPref+lTo == 0:
		return &badgerFromIterator{it, from}
	case lFrom != 0 && lTo != 0 && lPref == 0:
		return &badgerFromToIterator{it, from, to}
	case lFrom == 0:
		lastInPrefix := lastInPrefix(prefix, to)
		if keyLtEq(lastInPrefix, to) {
			return &badgerPrefixIterator{it, prefix}
		}
		return &badgerFromToIterator{it, prefix, to}
	case lTo == 0:
		if bytes.Compare(prefix, from) >= 0 {
			return &badgerPrefixIterator{it, prefix}
		}
		return &badgerPrefixFromIterator{&badgerPrefixIterator{it, prefix}, from}

	default: // all is set
		return iteratorOfAll(it, prefix, from, to)
	}
}

func lastInPrefix(prefix, to []byte) []byte {
	lastValueInPrefix := prefix
	if len(to)-len(prefix) > 0 {
		padding := make([]byte, len(to)-len(prefix))
		for i := range padding {
			padding[i] = uint8(255)
		}
		lastValueInPrefix = append(lastValueInPrefix, padding...)
	}
	return lastValueInPrefix
}

func iteratorOfAll(it *badger.Iterator, prefix, from, to []byte) Iterator {
	f, t := from, to
	if bytes.Compare(prefix, from) >= 0 {
		f = prefix
	}
	lastInPrefix := lastInPrefix(prefix, to)
	if keyLtEq(lastInPrefix, to) {
		if bytes.Equal(f, prefix) {
			return &badgerPrefixIterator{it, prefix}
		}
		return &badgerPrefixFromIterator{&badgerPrefixIterator{it, prefix}, f}
	}
	return &badgerFromToIterator{it, f, t}
}
