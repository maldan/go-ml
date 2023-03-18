package ml_expression

import "unsafe"

const TokenString = 0
const TokenOp = 1
const TokenIdentifier = 2
const TokenNumber = 3
const TokenBinaryExpression = 4

const KindBool = uint8(1)
const KindNumber = uint8(2)
const KindString = uint8(3)

const OperatorStringCompare = "strcmp"

type TokenType struct {
	Token string
	Type  uint8
	// Value any
	// Value any
	// B     []byte
	// Int   int
	List []TokenType
}

type ExpressionToken struct {
	Token   uint8
	Kind    uint8
	Value   any
	Address unsafe.Pointer
}

type Expression struct {
	TokenList []ExpressionToken
	Value     any
	Map       map[string]unsafe.Pointer
	Vars      []any
}

type emptyInterface struct {
	typ unsafe.Pointer
	ptr unsafe.Pointer
}
