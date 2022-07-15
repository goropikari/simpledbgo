package domain

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"log"
)

// Constant is constant type of database.
type Constant struct {
	typ FieldType
	val any
}

// NewConstant construts a Constant.
func NewConstant(typ FieldType, val any) Constant {
	return Constant{typ: typ, val: val}
}

// AsInt32 returns a value as int32.
func (c Constant) AsInt32() int32 {
	v, ok := c.val.(int32)
	if !ok {
		log.Fatal(errors.New("ToInt32 cant't convert Constant to int32"))
	}

	return v
}

// AsString returns a value as string.
func (c Constant) AsString() string {
	v, ok := c.val.(string)
	if !ok {
		log.Fatal(errors.New("AsString cant't convert Constant to string"))
	}

	return v
}

// String stringfies constant.
func (c Constant) String() string {
	return fmt.Sprintf("%v", c.val)
}

// AsVal returns constant as any.
func (c Constant) AsVal() any {
	return c.val
}

// HashCode return hash value of c.
func (c Constant) HashCode() int {
	mod := 998244353

	b := sha256.Sum256([]byte(fmt.Sprintf("%v", c)))

	x := 0
	for _, v := range b {
		x += int(v)
		x %= mod
	}

	return x
}

// Equal checks the equality of Constant.
func (c Constant) Equal(other Constant) bool {
	if c.typ != other.typ {
		return false
	}

	if c.val != other.val {
		return false
	}

	return true
}

func (c Constant) Less(other Constant) bool {
	if c.typ != other.typ {
		panic(errors.New("compare different types"))
	}

	switch c.typ {
	case Int32FieldType:
		return c.val.(int32) < other.val.(int32)
	case StringFieldType:
		return c.val.(string) < other.val.(string)
	case UnknownFieldType:
		panic(errors.New("unsupported field type"))
	default:
		panic(errors.New("unsupported field type"))
	}
}
