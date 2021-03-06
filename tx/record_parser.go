package tx

import (
	"fmt"

	"github.com/goropikari/simpledbgo/common"
	"github.com/goropikari/simpledbgo/errors"
	"github.com/goropikari/simpledbgo/lib/bytes"
	"github.com/goropikari/simpledbgo/tx/logrecord"
)

const (
	recordLengthOffset = 0
	recordOffset       = common.Int32Length
)

// ParseRecord parses b as log record.
func ParseRecord(b []byte) (logrecord.LogRecorder, error) {
	bb := bytes.NewBufferBytes(b)
	typ, err := bb.GetInt32(recordLengthOffset)
	if err != nil {
		return nil, errors.Err(err, "GetInt32")
	}

	data, err := bb.GetBytes(recordOffset)
	if err != nil {
		return nil, errors.Err(err, "GetBytes")
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
		return nil, errors.Err(err, "Unmarshal")
	}

	return rec, nil
}
