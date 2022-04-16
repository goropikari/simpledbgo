package tx_test

import (
	"testing"

	"github.com/goropikari/simpledbgo/backend/tx"
	"github.com/goropikari/simpledbgo/backend/tx/logrecord"
	"github.com/goropikari/simpledbgo/lib/bytes"
	"github.com/goropikari/simpledbgo/meta"
	"github.com/stretchr/testify/require"
)

func TestParseRecord(t *testing.T) {
	var tests = []struct {
		name   string
		typ    logrecord.RecordType
		record logrecord.LogRecorder
	}{
		{
			name:   "start log",
			typ:    logrecord.Start,
			record: &logrecord.StartRecord{TxNum: 1},
		},
		{
			name:   "commit log",
			typ:    logrecord.Commit,
			record: &logrecord.CommitRecord{TxNum: 1},
		},
		{
			name: "set int32 log",
			typ:  logrecord.SetInt32,
			record: &logrecord.SetInt32Record{
				FileName:    "hoge",
				TxNum:       1,
				BlockNumber: 2,
				Offset:      3,
				Val:         4,
			},
		},
		{
			name: "set string log",
			typ:  logrecord.SetString,
			record: &logrecord.SetStringRecord{
				FileName:    "hoge",
				TxNum:       1,
				BlockNumber: 2,
				Offset:      3,
				Val:         "piyo",
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			bb := bytes.NewBuffer(100)

			var err error
			err = bb.SetInt32(0, int32(tt.typ))
			require.NoError(t, err)

			data, err := tt.record.Marshal()
			require.NoError(t, err)

			err = bb.SetBytes(meta.Int32Length, data)
			require.NoError(t, err)

			rec, err := tx.ParseRecord(bb.GetData())
			require.NoError(t, err)
			require.Equal(t, tt.record, rec)
		})
	}
}

func TestParseRecord_Error(t *testing.T) {
	tests := []struct {
		name string
		data []byte
	}{
		{
			name: "failed GetInt32",
			data: []byte{0},
		},
		{
			name: "faialed GetBytes",
			data: []byte{0, 0, 0, 0, 0},
		},
		{
			name: "invalid record type",
			data: []byte{0, 0, 0, 255, 0, 0, 0, 0},
		},
		{
			name: "invalid record",
			data: []byte{0, 0, 0, 1, 0, 0, 0, 1, 0},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			_, err := tx.ParseRecord(tt.data)
			require.Error(t, err)
		})
	}
}
