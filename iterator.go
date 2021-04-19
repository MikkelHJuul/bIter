package biter

import (
	"bytes"
	"github.com/dgraph-io/badger/v3"
)

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

type badgerPrefixIterator struct {
	*badger.Iterator
	prefix []byte
}

//prefix is used as from
type badgerFromToIterator struct {
	*badgerPrefixIterator
	to []byte
}

type badgerFromIterator struct {
	*badgerPrefixIterator
}

type badgerToIterator struct {
	*badgerFromToIterator
}

// a badgerPrefixToIterator is simply a FromToIterator,
//while badgerPrefixFromIterator is slightly faster than
//a badgerFromToIterator iterator, and the logic wouldn't work anyway
type badgerPrefixFromIterator struct {
	*badgerPrefixIterator
	from []byte
}

func (b *badgerPrefixFromIterator) Rewind() {
	b.Seek(b.from)
}

func (b *badgerToIterator) Rewind() {
	b.Iterator.Rewind()
}

func (b *badgerFromIterator) Valid() bool {
	return b.Iterator.Valid()
}

func (b *badgerPrefixIterator) Rewind() {
	b.Seek(b.prefix)
}

func (b *badgerPrefixIterator) Valid() bool {
	return b.ValidForPrefix(b.prefix)
}

func (b *badgerFromToIterator) Valid() bool {
	return b.Iterator.Valid() && keyLtEq(b.Item().Key(), b.to)
}

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
		return &badgerToIterator{&badgerFromToIterator{&badgerPrefixIterator{it, nil}, to}}
	case lFrom != 0 && lPref+lTo == 0:
		return &badgerFromIterator{&badgerPrefixIterator{it, from}}
	case lFrom != 0 && lTo != 0 && lPref == 0:
		return &badgerFromToIterator{&badgerPrefixIterator{it, from}, to}
	case lFrom == 0:
		lastInPrefix := lastInPrefix(prefix, to)
		if keyLtEq(lastInPrefix, to) {
			return &badgerPrefixIterator{it, prefix}
		}
		return &badgerFromToIterator{&badgerPrefixIterator{it, prefix}, to}
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
	return &badgerFromToIterator{&badgerPrefixIterator{it, f}, t}
}
