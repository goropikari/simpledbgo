package domain

func (buf *Buffer) LSN() LSN {
	return buf.lsn
}

func (sch *Schema) SetFields(fields []FieldName) {
	sch.fields = fields
}

func (sch *Schema) SetInfo(info map[FieldName]*FieldInfo) {
	sch.info = info
}

func NewFieldInfo(typ FieldType, length int) *FieldInfo {
	return &FieldInfo{
		typ:    typ,
		length: length,
	}
}

func NewLayoutByElement(schema *Schema, offsets map[FieldName]int64, slotsize int64) *Layout {
	return &Layout{
		schema:   schema,
		offsets:  offsets,
		slotsize: slotsize,
	}
}
