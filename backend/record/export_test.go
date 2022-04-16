package record

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