package metadata

type (
	// ViewName is type of view name.
	ViewName = string

	// ViewDef is type of view definition.
	ViewDef = string
)

const (
	maxTableNameLength = 64
	maxFieldNameLength = 16

	tableCatalog = "table_catalog"
	fieldCatalog = "field_catalog"

	fldTableName = "table_name"
	fldFieldName = "field_name"
	fldSlotSize  = "slot_size"
	fldType      = "type"
	fldLength    = "length"
	fldOffset    = "offset"
)

const (
	maxViewDefLength = 100

	fldViewName    = "view_name"
	fldViewDef     = "view_def"
	fldViewCatalog = "view_catalog"
)
