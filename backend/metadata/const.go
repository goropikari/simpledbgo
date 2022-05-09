package metadata

type (
	// ViewDef is type of view definition.
	ViewDef = string
)

const (
	tableCatalog    = "table_catalog"
	fieldCatalog    = "field_catalog"
	fldIndexCatalog = "index_catalog"

	fldTableName = "table_name"
	fldFieldName = "field_name"
	fldSlotSize  = "slot_size"
	fldType      = "type"
	fldLength    = "length"
	fldOffset    = "offset"

	fldIndexName = "indexname"
)

const (
	fldViewName    = "view_name"
	fldViewDef     = "view_def"
	fldViewCatalog = "view_catalog"
)

const (
	updateTimes = 100
)
