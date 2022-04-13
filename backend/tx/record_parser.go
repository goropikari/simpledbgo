package tx

import (
	"fmt"

	"github.com/goropikari/simpledb_go/backend/tx/logrecord"
	"github.com/goropikari/simpledb_go/lib/bytes"
	"github.com/goropikari/simpledb_go/meta"
)

// RecordParse parses b as log record.
func RecordParse(b []byte) (logrecord.LogRecorder, error) {
	bb := bytes.NewBufferBytes(b)
	typ, err := bb.GetInt32(0)
	if err != nil {
		return nil, err
	}

	data, err := bb.GetBytes(meta.Int32Length)
	if err != nil {
		return nil, err
	}

	var rec logrecord.LogRecorder
	switch typ {
	case Start:
		rec = &logrecord.StartRecord{}
	case Commit:
		rec = &logrecord.CommitRecord{}
	case SetInt32:
		rec = &logrecord.SetInt32Record{}
	case SetString:
		rec = &logrecord.SetStringRecord{}
	default:
		panic(fmt.Errorf("unexpected record type: %v", typ))
	}

	if err := rec.Unmarshal(data); err != nil {
		return nil, err
	}

	return rec, nil
}
