package ml_expression

import (
	"reflect"
	"strings"
	"unsafe"
)

func (e *Expression) Bind(v any) {
	e.Value = v

	start := reflect.ValueOf(v).UnsafePointer()

	typeOf := reflect.TypeOf(v).Elem()
	valueOf := reflect.ValueOf(v).Elem()
	for i := 0; i < valueOf.NumField(); i++ {
		e.Map[typeOf.Field(i).Name] = unsafe.Add(start, typeOf.Field(i).Offset)
	}
}

func (e *Expression) Execute() any {
	stackCounter := 0
	anyStack := make([]any, 16)
	ops := e.TokenList

	for i := 0; i < len(ops); i++ {
		v := ops[i]

		switch v.Token {
		case TokenIdentifier:
			anyStack[stackCounter] = e.Map[v.Value.(string)]
			stackCounter += 1
			break
		case TokenString:
			anyStack[stackCounter] = v.Address
			stackCounter += 1
			break
		case TokenNumber:
			anyStack[stackCounter] = v.Address
			stackCounter += 1
			break
		case TokenOp:
			switch v.Value.(string) {
			case "not in":
				stackCounter -= 1
				b := anyStack[stackCounter].(unsafe.Pointer)
				stackCounter -= 1
				a := anyStack[stackCounter].(unsafe.Pointer)

				anyStack[stackCounter] = !strings.Contains(*(*string)(b), *(*string)(a))
				stackCounter += 1

				break
			case "in":
				stackCounter -= 1
				b := anyStack[stackCounter].(unsafe.Pointer)
				stackCounter -= 1
				a := anyStack[stackCounter].(unsafe.Pointer)

				anyStack[stackCounter] = strings.Contains(*(*string)(b), *(*string)(a))
				stackCounter += 1

				break
			case OperatorStringCompare:
				stackCounter -= 1
				b := anyStack[stackCounter].(unsafe.Pointer)
				stackCounter -= 1
				a := anyStack[stackCounter].(unsafe.Pointer)

				anyStack[stackCounter] = *(*string)(a) == *(*string)(b)
				stackCounter += 1
				break
			case "==":
				stackCounter -= 1
				b := anyStack[stackCounter].(unsafe.Pointer)
				stackCounter -= 1
				a := anyStack[stackCounter].(unsafe.Pointer)

				// fmt.Printf("B: %v\n", *(*string)(b))
				// fmt.Printf("A: %v\n", *(*string)(a))

				anyStack[stackCounter] = *(*int)(a) == *(*int)(b)
				stackCounter += 1
				break
			case "&&":
				stackCounter -= 1
				b := anyStack[stackCounter]
				stackCounter -= 1
				a := anyStack[stackCounter]

				anyStack[stackCounter] = a == b
				stackCounter += 1
				break
			case ">":
				stackCounter -= 1
				b := anyStack[stackCounter].(unsafe.Pointer)
				stackCounter -= 1
				a := anyStack[stackCounter].(unsafe.Pointer)

				anyStack[stackCounter] = *(*int)(a) > *(*int)(b)
				stackCounter += 1
				break
			case "+":
				stackCounter -= 1
				b := anyStack[stackCounter].(unsafe.Pointer)
				stackCounter -= 1
				a := anyStack[stackCounter].(unsafe.Pointer)

				anyStack[stackCounter] = *(*int)(a) + *(*int)(b)
				stackCounter += 1
				break
			}
			break
		}
	}

	return anyStack[0]
}
