package server

import (
	"encoding/binary"
	"fmt"

	"github.com/goropikari/simpledbgo/domain"
	"github.com/goropikari/simpledbgo/errors"
)

type Result struct {
	typ     ResultType
	records [][]any
	fields  []string
}

func makeDataRow(rec []any) []byte {
	dataRow := make([]byte, 0)
	nc := len(rec)
	ncb := make([]byte, 2)
	binary.BigEndian.PutUint16(ncb, uint16(nc))
	dataRow = append(dataRow, ncb...)
	for _, val := range rec {
		if val == nil {
			dataRow = append(dataRow, []byte{0xff, 0xff, 0xff, 0xff}...)
		} else {
			s := fmt.Sprintf("%v", val)
			sb := []byte(s)
			slen := len(sb)
			lenByte := make([]byte, payloadBytesLength)
			binary.BigEndian.PutUint32(lenByte, uint32(slen))
			dataRow = append(dataRow, lenByte...)
			dataRow = append(dataRow, sb...)
		}
	}

	payload := make([]byte, 0)
	payload = append(payload, 'D') // 0x44 -> D: DataRow
	lenByte := make([]byte, payloadBytesLength)
	binary.BigEndian.PutUint32(lenByte, uint32(len(dataRow)+payloadBytesLength))
	payload = append(payload, lenByte...)
	payload = append(payload, dataRow...)

	return payload
}

func makeColDesc(cols []string) []byte {
	payload := make([]byte, 0)
	n := len(cols)
	numCols := make([]byte, 2)
	binary.BigEndian.PutUint16(numCols, uint16(n))
	payload = append(payload, numCols...)

	for k, col := range cols {
		payload = append(payload, []byte(col)...)
		payload = append(payload, nullEnd)
		payload = append(payload, []byte{0x00, 0x00, 0x40, 0x06}...) // object id
		idx := make([]byte, 2)
		binary.BigEndian.PutUint16(idx, uint16(k+1))
		payload = append(payload, idx...)                            // col id
		payload = append(payload, []byte{0x00, 0x00, 0x04, 0x13}...) // data type
		payload = append(payload, []byte{0xff, 0xff}...)             // data type size
		payload = append(payload, []byte{0xff, 0xff, 0xff, 0xff}...) // type modifier
		payload = append(payload, []byte{0x00, 0x00}...)             // format code
	}

	length := make([]byte, payloadBytesLength)
	binary.BigEndian.PutUint32(length, uint32(len(payload)+payloadBytesLength))
	packet := make([]byte, 0)
	packet = append(packet, 'T') // 0x54 -> T: RowDescription
	packet = append(packet, length...)
	packet = append(packet, payload...)

	return packet
}

func makeDataRows(recs [][]any) []byte {
	dataRows := make([]byte, 0)
	for _, rec := range recs {
		dataRows = append(dataRows, makeDataRow(rec)...)
	}

	return dataRows
}

type Rows struct {
	scan   domain.Scanner
	fields []domain.FieldName
}

func (cn *Connection) makeResult(rows *Rows) (Result, error) {
	recs := make([][]any, 0)
	for rows.scan.HasNext() {
		rec := make([]any, 0)
		for _, fld := range rows.fields {
			v, err := rows.scan.GetVal(fld)
			if err != nil {
				return Result{}, errors.Err(err, "GetVal")
			}
			rec = append(rec, v)
		}
		recs = append(recs, rec)
	}
	if err := rows.scan.Err(); err != nil {
		return Result{}, errors.Err(err, "HasNext")
	}

	fields := make([]string, 0)
	for _, fld := range rows.fields {
		fields = append(fields, string(fld))
	}

	return Result{
		typ:     queryResult,
		records: recs,
		fields:  fields,
	}, nil
}
