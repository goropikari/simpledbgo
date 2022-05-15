package domain_test

import (
	"testing"

	"github.com/goropikari/simpledbgo/domain"
	"github.com/goropikari/simpledbgo/testing/fake"
	"github.com/stretchr/testify/require"
)

func TestSchema_AddField(t *testing.T) {
	fld1 := domain.FieldName(fake.RandString())
	intSch := domain.NewSchema()
	intSch.SetFields([]domain.FieldName{fld1})
	mp1 := make(map[domain.FieldName]*domain.FieldInfo)
	mp1[fld1] = domain.NewFieldInfo(domain.FInt32, 0)
	intSch.SetInfo(mp1)

	fld2 := domain.FieldName(fake.RandString())
	strSch := domain.NewSchema()
	strSch.SetFields([]domain.FieldName{fld2})
	mp2 := make(map[domain.FieldName]*domain.FieldInfo)
	mp2[fld2] = domain.NewFieldInfo(domain.FString, 8)
	strSch.SetInfo(mp2)

	tests := []struct {
		name     string
		fldname  domain.FieldName
		typ      domain.FieldType
		length   int
		expected *domain.Schema
	}{
		{
			name:     "add int field",
			fldname:  fld1,
			typ:      domain.FInt32,
			expected: intSch,
		},
		{
			name:     "add int field",
			fldname:  fld2,
			typ:      domain.FString,
			length:   8,
			expected: strSch,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			sch := domain.NewSchema()
			sch.AddField(tt.fldname, tt.typ, tt.length)
			require.Equal(t, tt.expected, sch)
		})
	}
}

func TestSchema_AddIntField(t *testing.T) {
	fld1 := domain.FieldName(fake.RandString())
	intSch := domain.NewSchema()
	intSch.SetFields([]domain.FieldName{fld1})
	mp1 := make(map[domain.FieldName]*domain.FieldInfo)
	mp1[fld1] = domain.NewFieldInfo(domain.FInt32, 0)
	intSch.SetInfo(mp1)

	tests := []struct {
		name     string
		fldname  domain.FieldName
		expected *domain.Schema
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
			sch := domain.NewSchema()
			sch.AddInt32Field(tt.fldname)
			require.Equal(t, tt.expected, sch)
		})
	}
}

func TestSchema_AddStringField(t *testing.T) {
	fld1 := domain.FieldName(fake.RandString())
	strSch := domain.NewSchema()
	strSch.SetFields([]domain.FieldName{fld1})
	mp1 := make(map[domain.FieldName]*domain.FieldInfo)
	mp1[fld1] = domain.NewFieldInfo(domain.FString, 8)
	strSch.SetInfo(mp1)

	tests := []struct {
		name     string
		fldname  domain.FieldName
		length   int
		expected *domain.Schema
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
			sch := domain.NewSchema()
			sch.AddStringField(tt.fldname, tt.length)
			require.Equal(t, tt.expected, sch)
		})
	}
}

func TestSchema_Add(t *testing.T) {
	fld1 := domain.FieldName(fake.RandString())
	fld2 := domain.FieldName(fake.RandString())
	bsch := domain.NewSchema()
	bsch.SetFields([]domain.FieldName{fld1, fld2})
	mp1 := make(map[domain.FieldName]*domain.FieldInfo)
	mp1[fld1] = domain.NewFieldInfo(domain.FString, 8)
	mp1[fld2] = domain.NewFieldInfo(domain.FInt32, 0)
	bsch.SetInfo(mp1)

	strSch := domain.NewSchema()
	mp2 := make(map[domain.FieldName]*domain.FieldInfo)
	mp2[fld1] = domain.NewFieldInfo(domain.FString, 8)
	strSch.SetFields([]domain.FieldName{fld1})
	strSch.SetInfo(mp2)

	intSch := domain.NewSchema()
	mp3 := make(map[domain.FieldName]*domain.FieldInfo)
	mp3[fld2] = domain.NewFieldInfo(domain.FInt32, 0)
	intSch.SetFields([]domain.FieldName{fld2})
	intSch.SetInfo(mp3)

	tests := []struct {
		name     string
		fldname  domain.FieldName
		base     *domain.Schema
		expected *domain.Schema
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
			sch := domain.NewSchema()
			sch.Add(tt.fldname, tt.base)
			require.Equal(t, tt.expected, sch)
		})
	}
}

func TestSchema_AddAllFields(t *testing.T) {
	fld1 := domain.FieldName(fake.RandString())
	fld2 := domain.FieldName(fake.RandString())
	bsch := domain.NewSchema()
	bsch.SetFields([]domain.FieldName{fld1, fld2})
	mp1 := make(map[domain.FieldName]*domain.FieldInfo)
	mp1[fld1] = domain.NewFieldInfo(domain.FString, 8)
	mp1[fld2] = domain.NewFieldInfo(domain.FInt32, 0)
	bsch.SetInfo(mp1)

	tests := []struct {
		name     string
		base     *domain.Schema
		expected *domain.Schema
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
			sch := domain.NewSchema()
			sch.AddAllFields(tt.base)
			require.Equal(t, tt.expected, sch)
		})
	}
}

func TestLayout(t *testing.T) {
	t.Run("constructs Layout", func(t *testing.T) {
		schema := domain.NewSchema()
		schema.AddField("hoge", domain.FInt32, 0)
		schema.AddField("piyo", domain.FString, 8)

		layout := domain.NewLayout(schema)

		mp := map[domain.FieldName]int64{
			"hoge": 4,
			"piyo": 8,
		}

		expected := domain.NewLayoutByElement(schema, mp, 20)

		require.Equal(t, expected, layout)
	})
}
