package impl

import (
	"github.com/dgraph-io/badger/v3"
	"reflect"
	"testing"
)

func TestKeyRangeIterator(t *testing.T) {
	bOpts := badger.DefaultOptions("").WithInMemory(true)
	aBadger, err := badger.Open(bOpts)
	if err != nil {
		t.Error("unable to start database!")
	}
	anIterator := aBadger.NewTransaction(false).NewIterator(badger.DefaultIteratorOptions)
	type args struct {
		prefix []byte
		from   []byte
		to     []byte
	}
	tests := []struct {
		name string
		args args
		want Iterator
	}{
		{
			name: "none-returns the wrapping iterator",
			want: anIterator,
			args: args{},
		}, {
			name: "prefix returns a prefix iterator",
			want: &badgerPrefixIterator{anIterator, []byte("prefix")},
			args: args{prefix: []byte("prefix")},
		}, {
			name: "from-to returns a from-to iterator",
			want: &badgerFromToIterator{&badgerPrefixIterator{anIterator, []byte("from")}, []byte("to")},
			args: args{from: []byte("from"), to: []byte("to")},
		}, {
			name: "from returns a from iterator",
			want: &badgerFromIterator{&badgerPrefixIterator{anIterator, []byte("from")}},
			args: args{from: []byte("from")},
		}, {
			name: "to returns a to iterator",
			want: &badgerToIterator{&badgerFromToIterator{&badgerPrefixIterator{anIterator, nil}, []byte("to")}},
			args: args{to: []byte("to")},
		}, {
			name: "to-prefix returns a FromTo iterator",
			want: &badgerFromToIterator{&badgerPrefixIterator{anIterator, []byte("d")}, []byte("dg")},
			args: args{to: []byte("dg"), prefix: []byte("d")},
		}, {
			name: "to-prefix returns a FromTo iterator, which is invalid",
			want: &badgerFromToIterator{&badgerPrefixIterator{anIterator, []byte("d")}, []byte("c")},
			args: args{to: []byte("c"), prefix: []byte("d")},
		}, {
			name: "to-prefix returns a Prefix iterator",
			want: &badgerPrefixIterator{anIterator, []byte("g")},
			args: args{to: []byte("z"), prefix: []byte("g")},
		}, {
			name: "from-prefix returns a FromTo iterator",
			want: &badgerPrefixFromIterator{&badgerPrefixIterator{anIterator, []byte("d")}, []byte("dg")},
			args: args{from: []byte("dg"), prefix: []byte("d")},
		}, {
			name: "from-prefix returns a Prefix iterator",
			want: &badgerPrefixIterator{anIterator, []byte("i")},
			args: args{from: []byte("g"), prefix: []byte("i")},
		}, {
			name: "from-prefix returns a PrefixFrom iterator, which will be invalid",
			want: &badgerPrefixFromIterator{&badgerPrefixIterator{anIterator, []byte("i")}, []byte("k")},
			args: args{from: []byte("k"), prefix: []byte("i")},
		}, {
			name: "from-to-prefix returns a FromTo Iterator",
			want: &badgerFromToIterator{&badgerPrefixIterator{anIterator, []byte("ka")}, []byte("kb")},
			args: args{from: []byte("ka"), to: []byte("kb"), prefix: []byte("k")},
		}, {
			name: "from-to-prefix returns a FromTo Iterator a out of bounds",
			want: &badgerFromToIterator{&badgerPrefixIterator{anIterator, []byte("k")}, []byte("kb")},
			args: args{from: []byte("a"), to: []byte("kb"), prefix: []byte("k")},
		}, {
			name: "from-to-prefix returns a Prefix Iterator",
			want: &badgerPrefixIterator{anIterator, []byte("b")},
			args: args{from: []byte("a"), to: []byte("c"), prefix: []byte("b")},
		}, {
			name: "from-to-prefix returns a PrefixFrom Iterator",
			want: &badgerPrefixFromIterator{&badgerPrefixIterator{anIterator, []byte("b")}, []byte("bb")},
			args: args{from: []byte("bb"), to: []byte("c"), prefix: []byte("b")},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := KeyRangeIterator(anIterator, tt.args.prefix, tt.args.from, tt.args.to)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("KeyRangeIterator() = %v, want %v", got, tt.want)
			}
		})
	}
}
