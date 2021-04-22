package biter

import (
	"bytes"
	"fmt"
	"github.com/dgraph-io/badger/v3"
	"testing"
)

//TestDifferentQueryScenarios is adapted from github.com/MikkelHJuul/ld/impl/get_test.go
func TestDifferentQueryScenarios(t *testing.T) {
	bOpts := badger.DefaultOptions("").WithInMemory(true)
	badgerDB, err := badger.Open(bOpts)
	if err != nil {
		t.Fatal("could not initiate database!", err)
	}
	txn := badgerDB.NewTransaction(true)
	for i := 0; i < 100; i++ {
		key := fmt.Sprintf("%02d", i)
		if err = txn.Set([]byte(key), []byte(key)); err != nil {
			t.Fatal("error setting DB-values", err)
		}
	}
	if err = txn.Commit(); err != nil {
		t.Fatal("could not commit initialisation", err)
	}

	tests := []struct {
		name             string
		From, To, Prefix []byte
		response         [][]byte
	}{
		{
			name: "get range within",
			From: []byte("12"),
			To:   []byte("17"),
			response: [][]byte{
				{'1', '2'},
				{'1', '3'},
				{'1', '4'},
				{'1', '5'},
				{'1', '6'},
				{'1', '7'}},
		},
		{
			name:     "get range overlap",
			From:     []byte("99"),
			To:       []byte("a"),
			response: [][]byte{{'9', '9'}},
		},
		{
			name:   "get prefix",
			Prefix: []byte{'9'},
			response: [][]byte{
				{'9', '0'},
				{'9', '1'},
				{'9', '2'},
				{'9', '3'},
				{'9', '4'},
				{'9', '5'},
				{'9', '6'},
				{'9', '7'},
				{'9', '8'},
				{'9', '9'}},
		},
		{
			name:   "get prefix From",
			Prefix: []byte("9"),
			From:   []byte("92"),
			response: [][]byte{
				{'9', '2'},
				{'9', '3'},
				{'9', '4'},
				{'9', '5'},
				{'9', '6'},
				{'9', '7'},
				{'9', '8'},
				{'9', '9'}},
		},
		{
			name:   "get prefix To",
			Prefix: []byte("9"),
			To:     []byte("92"),
			response: [][]byte{
				{'9', '0'},
				{'9', '1'},
				{'9', '2'}},
		},
		{
			name:   "get prefix FromTo",
			Prefix: []byte("9"),
			From:   []byte("91"),
			To:     []byte("92"),
			response: [][]byte{
				{'9', '1'},
				{'9', '2'}},
		},
		{
			name:   "get prefix and from to mismatch",
			Prefix: []byte("9"),
			From:   []byte("12"),
			To:     []byte("78"),
		},
		{
			name:   "get prefix and from to mismatch",
			Prefix: []byte("5"),
			From:   []byte("60"),
		},
		{
			name:   "get prefix and from to mismatch",
			Prefix: []byte("5"),
			To:     []byte("49"),
		},
		{
			name: "get range outside",
			From: []byte("a"),
		},
		{
			name: "get range to zero",
			To:   []byte("0"),
		},
		{
			name: "get To one",
			To:   []byte("1"),
			response: [][]byte{
				{'0', '0'},
				{'0', '1'},
				{'0', '2'},
				{'0', '3'},
				{'0', '4'},
				{'0', '5'},
				{'0', '6'},
				{'0', '7'},
				{'0', '8'},
				{'0', '9'}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			txn := badgerDB.NewTransaction(false)
			defer txn.Commit()
			it := txn.NewIterator(badger.DefaultIteratorOptions)
			defer it.Close()
			iter := KeyRangeIterator(it, tt.Prefix, tt.From, tt.To)
			var keys [][]byte
			for iter.Rewind(); iter.Valid(); iter.Next() {
				keys = append(keys, iter.Item().Key())
			}
			validateReturn(t, tt.response, keys)
		})
	}
}

func validateReturn(t *testing.T, expected, got [][]byte) {
	if len(got) != len(expected) {
		t.Errorf("not the same amount of results, %d =|= %d", len(expected), len(got))
	}
	for _, aVal := range got {
		isThere := false
		for _, res := range expected {
			if bytes.Equal(aVal, res) {
				isThere = true
				break
			}
		}
		if !isThere {
			t.Errorf("results are not alike! %v != %v", got, expected)
		}
	}
}
