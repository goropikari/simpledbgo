package record

import "github.com/goropikari/simpledbgo/backend/domain"

func (sch *Schema) SetFields(fields []domain.FieldName) {
	sch.fields = fields
}

func (sch *Schema) SetInfo(info map[domain.FieldName]*FieldInfo) {
	sch.info = info
}

func NewFieldInfo(typ FieldType, length int) *FieldInfo {
	return &FieldInfo{
		typ:    typ,
		length: length,
	}
}

func NewLayoutByElement(schema *Schema, offsets map[domain.FieldName]int64, slotsize int64) *Layout {
	return &Layout{
		schema:   schema,
		offsets:  offsets,
		slotsize: slotsize,
	}
}
