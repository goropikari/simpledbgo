package record_test

import (
	"testing"

	"github.com/goropikari/simpledbgo/backend/domain"
	"github.com/goropikari/simpledbgo/backend/record"
	"github.com/goropikari/simpledbgo/testing/fake"
	"github.com/stretchr/testify/require"
)

func TestSchema_AddField(t *testing.T) {
	fld1 := domain.FieldName(fake.RandString())
	intSch := record.NewSchema()
	intSch.SetFields([]domain.FieldName{fld1})
	mp1 := make(map[domain.FieldName]*record.FieldInfo)
	mp1[fld1] = record.NewFieldInfo(record.Int32, 0)
	intSch.SetInfo(mp1)

	fld2 := domain.FieldName(fake.RandString())
	strSch := record.NewSchema()
	strSch.SetFields([]domain.FieldName{fld2})
	mp2 := make(map[domain.FieldName]*record.FieldInfo)
	mp2[fld2] = record.NewFieldInfo(record.String, 8)
	strSch.SetInfo(mp2)

	tests := []struct {
		name     string
		fldname  domain.FieldName
		typ      record.FieldType
		length   int
		expected *record.Schema
	}{
		{
			name:     "add int field",
			fldname:  fld1,
			typ:      record.Int32,
			expected: intSch,
		},
		{
			name:     "add int field",
			fldname:  fld2,
			typ:      record.String,
			length:   8,
			expected: strSch,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			sch := record.NewSchema()
			sch.AddField(tt.fldname, tt.typ, tt.length)
			require.Equal(t, tt.expected, sch)
		})
	}
}

func TestSchema_AddIntField(t *testing.T) {
	fld1 := domain.FieldName(fake.RandString())
	intSch := record.NewSchema()
	intSch.SetFields([]domain.FieldName{fld1})
	mp1 := make(map[domain.FieldName]*record.FieldInfo)
	mp1[fld1] = record.NewFieldInfo(record.Int32, 0)
	intSch.SetInfo(mp1)

	tests := []struct {
		name     string
		fldname  domain.FieldName
		expected *record.Schema
	}{
		{
			name:     "add int field",
			fldname:  fld1,
			expected: intSch,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			sch := record.NewSchema()
			sch.AddInt32Field(tt.fldname)
			require.Equal(t, tt.expected, sch)
		})
	}
}

func TestSchema_AddStringField(t *testing.T) {
	fld1 := domain.FieldName(fake.RandString())
	strSch := record.NewSchema()
	strSch.SetFields([]domain.FieldName{fld1})
	mp1 := make(map[domain.FieldName]*record.FieldInfo)
	mp1[fld1] = record.NewFieldInfo(record.String, 8)
	strSch.SetInfo(mp1)

	tests := []struct {
		name     string
		fldname  domain.FieldName
		length   int
		expected *record.Schema
	}{
		{
			name:     "add string field",
			fldname:  fld1,
			length:   8,
			expected: strSch,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			sch := record.NewSchema()
			sch.AddStringField(tt.fldname, tt.length)
			require.Equal(t, tt.expected, sch)
		})
	}
}

func TestSchema_Add(t *testing.T) {
	fld1 := domain.FieldName(fake.RandString())
	fld2 := domain.FieldName(fake.RandString())
	bsch := record.NewSchema()
	bsch.SetFields([]domain.FieldName{fld1, fld2})
	mp1 := make(map[domain.FieldName]*record.FieldInfo)
	mp1[fld1] = record.NewFieldInfo(record.String, 8)
	mp1[fld2] = record.NewFieldInfo(record.Int32, 0)
	bsch.SetInfo(mp1)

	strSch := record.NewSchema()
	mp2 := make(map[domain.FieldName]*record.FieldInfo)
	mp2[fld1] = record.NewFieldInfo(record.String, 8)
	strSch.SetFields([]domain.FieldName{fld1})
	strSch.SetInfo(mp2)

	intSch := record.NewSchema()
	mp3 := make(map[domain.FieldName]*record.FieldInfo)
	mp3[fld2] = record.NewFieldInfo(record.Int32, 0)
	intSch.SetFields([]domain.FieldName{fld2})
	intSch.SetInfo(mp3)

	tests := []struct {
		name     string
		fldname  domain.FieldName
		base     *record.Schema
		expected *record.Schema
	}{
		{
			name:     "add string field",
			fldname:  fld1,
			base:     bsch,
			expected: strSch,
		},
		{
			name:     "add string field",
			fldname:  fld2,
			base:     bsch,
			expected: intSch,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			sch := record.NewSchema()
			sch.Add(tt.fldname, tt.base)
			require.Equal(t, tt.expected, sch)
		})
	}
}

func TestSchema_AddAllFields(t *testing.T) {
	fld1 := domain.FieldName(fake.RandString())
	fld2 := domain.FieldName(fake.RandString())
	bsch := record.NewSchema()
	bsch.SetFields([]domain.FieldName{fld1, fld2})
	mp1 := make(map[domain.FieldName]*record.FieldInfo)
	mp1[fld1] = record.NewFieldInfo(record.String, 8)
	mp1[fld2] = record.NewFieldInfo(record.Int32, 0)
	bsch.SetInfo(mp1)

	tests := []struct {
		name     string
		base     *record.Schema
		expected *record.Schema
	}{
		{
			name:     "add all fields",
			base:     bsch,
			expected: bsch,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			sch := record.NewSchema()
			sch.AddAllFields(tt.base)
			require.Equal(t, tt.expected, sch)
		})
	}
}
