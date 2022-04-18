package meta

const (
	// Int32Length is byte length of int32.
	Int32Length = 4

	// Uint32Length is byte length of uint32.
	Uint32Length = 4
)

// Constant is constant type of database.
type Constant struct {
	I32val int32
	Sval   string
}
