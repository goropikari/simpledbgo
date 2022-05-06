package domain

// IndexName is a value object of index name.
type IndexName string

// NewIndexName constructs IndexName.
func NewIndexName(name string) (IndexName, error) {
	if len(name) > MaxFieldNameLength {
		return "", ErrExceedMaxFieldNameLength
	}

	return IndexName(name), nil
}

// String stringfy name.
func (name IndexName) String() string {
	return string(name)
}
