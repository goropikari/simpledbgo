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
	case logrecord.Start:
		rec = &logrecord.StartRecord{}
	case logrecord.Commit:
		rec = &logrecord.CommitRecord{}
	case logrecord.Checkpoint:
		rec = &logrecord.CheckpointRecord{}
	case logrecord.SetInt32:
		rec = &logrecord.SetInt32Record{}
	case logrecord.SetString:
		rec = &logrecord.SetStringRecord{}
	case logrecord.Rollback:
		rec = &logrecord.RollbackRecord{}
	default:
		return nil, fmt.Errorf("unexpected record type: %v", typ)
	}

	if err := rec.Unmarshal(data); err != nil {
		return nil, err
	}

	return rec, nil
}
