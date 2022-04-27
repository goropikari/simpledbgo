package meta

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"log"
)

const (
	// Int32Length is byte length of int32.
	Int32Length = 4

	// Uint32Length is byte length of uint32.
	Uint32Length = 4
)

// ValueType is a type of value.
type ValueType uint

const (
	// Undefined means undef type.
	Undefined ValueType = iota

	// Int32 means int32 type.
	Int32

	// String means string type.
	String
)

// Constant is constant type of database.
type Constant struct {
	typ ValueType
	val any
}

// NewConstant construts a Constant.
func NewConstant(typ ValueType, val any) Constant {
	return Constant{typ: typ, val: val}
}

// ToInt32 returns a value as int32.
func (c Constant) ToInt32() int32 {
	v, ok := c.val.(int32)
	if !ok {
		log.Fatal(errors.New("Cant't convert Constant to int32"))
	}

	return v
}

// ToString returns a value as string.
func (c Constant) ToString() string {
	v, ok := c.val.(string)
	if !ok {
		log.Fatal(errors.New("Cant't convert Constant to string"))
	}

	return v
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
