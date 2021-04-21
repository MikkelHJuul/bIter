import "testing"

func TestDifferentQueryScenarios(t *testing.T) {
	tests := []struct {
		name     string
		From, To, Prefix []byte
		response [][]byte
	}{
		{
			name:     "get range within",
			From: []byte("12"), 
			To: []byte("17"),
			response: [](
				[]byte("12"),
				[]byte("13"),
				[]byte("14"),
				[]byte("15"),
				[]byte("16"),
				[]byte("17")),
		},
		{
			name:     "get range overlap",
			From: "99", 
			To: "a",
			response: [][]byte("99"),
		},
		{
			name:     "get prefix",
			keyRange: &proto.KeyRange{Prefix: "9"},
			response: []*proto.KeyValue{
				{Key: "90", Value: []byte("90")},
				{Key: "91", Value: []byte("91")},
				{Key: "92", Value: []byte("92")},
				{Key: "93", Value: []byte("93")},
				{Key: "94", Value: []byte("94")},
				{Key: "95", Value: []byte("95")},
				{Key: "96", Value: []byte("96")},
				{Key: "97", Value: []byte("97")},
				{Key: "98", Value: []byte("98")},
				{Key: "99", Value: []byte("99")},
			},
		},
		{
			name:     "get prefix From",
			keyRange: &proto.KeyRange{Prefix: "9", From: "92"},
			response: []*proto.KeyValue{
				{Key: "92", Value: []byte("92")},
				{Key: "93", Value: []byte("93")},
				{Key: "94", Value: []byte("94")},
				{Key: "95", Value: []byte("95")},
				{Key: "96", Value: []byte("96")},
				{Key: "97", Value: []byte("97")},
				{Key: "98", Value: []byte("98")},
				{Key: "99", Value: []byte("99")},
			},
		},
		{
			name:     "get prefix To",
			keyRange: &proto.KeyRange{Prefix: "9", To: "92"},
			response: []*proto.KeyValue{
				{Key: "90", Value: []byte("90")},
				{Key: "91", Value: []byte("91")},
				{Key: "92", Value: []byte("92")},
			},
		},
		{
			name:     "get prefix FromTo",
			keyRange: &proto.KeyRange{Prefix: "9", From: "91", To: "92"},
			response: []*proto.KeyValue{
				{Key: "91", Value: []byte("91")},
				{Key: "92", Value: []byte("92")},
			},
		},
		{
			name:     "get prefix Pattern",
			keyRange: &proto.KeyRange{Prefix: "9", Pattern: ".2"},
			response: []*proto.KeyValue{
				{Key: "92", Value: []byte("92")},
			},
		},
		{
			name:     "get Pattern",
			keyRange: &proto.KeyRange{Pattern: ".3"},
			response: []*proto.KeyValue{
				{Key: "03", Value: []byte("03")},
				{Key: "13", Value: []byte("13")},
				{Key: "23", Value: []byte("23")},
				{Key: "33", Value: []byte("33")},
				{Key: "43", Value: []byte("43")},
				{Key: "53", Value: []byte("53")},
				{Key: "63", Value: []byte("63")},
				{Key: "73", Value: []byte("73")},
				{Key: "83", Value: []byte("83")},
				{Key: "93", Value: []byte("93")},
			},
		},
		{
			name:     "get prefix and from to mismatch",
			keyRange: &proto.KeyRange{Prefix: "9", From: "12", To: "78"},
			response: []*proto.KeyValue{},
		},
		{
			name:     "get prefix and from to mismatch",
			keyRange: &proto.KeyRange{Prefix: "5", From: "60"},
			response: []*proto.KeyValue{},
		},
		{
			name:     "get prefix and from to mismatch",
			keyRange: &proto.KeyRange{Prefix: "5", To: "49"},
			response: []*proto.KeyValue{},
		},
		{
			name:     "get range outside",
			keyRange: &proto.KeyRange{From: "a"},
			response: []*proto.KeyValue{},
		},
		{
			name:     "get range to zero",
			keyRange: &proto.KeyRange{To: "0"},
			response: []*proto.KeyValue{},
		},
		{
			name:     "get To one",
			keyRange: &proto.KeyRange{To: "1"},
			response: []*proto.KeyValue{
				{Key: "00", Value: []byte("00")},
				{Key: "01", Value: []byte("01")},
				{Key: "02", Value: []byte("02")},
				{Key: "03", Value: []byte("03")},
				{Key: "04", Value: []byte("04")},
				{Key: "05", Value: []byte("05")},
				{Key: "06", Value: []byte("06")},
				{Key: "07", Value: []byte("07")},
				{Key: "08", Value: []byte("08")},
				{Key: "09", Value: []byte("09")},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			txn := b.NewTransaction(false)
			defer txn.Commit()
			it := txn.NewIterator(badger.DefaultIteratorOptions)
			defer it.Close()
			iter := KeyRangeIterator(it, tt.Prefix, tt.From, tt.To)
			for iter.Rewind(); iter.Valid(); iter.Next() {
				item := iter.Item()
			validateReturn(t, tt.response, server.receive)
		})
	}
}
