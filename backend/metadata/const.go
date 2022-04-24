package metadata

type (
	// ViewName is type of view name.
	ViewName = string

	// ViewDef is type of view definition.
	ViewDef = string
)

const (
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
	fldViewName    = "view_name"
	fldViewDef     = "view_def"
	fldViewCatalog = "view_catalog"
)
