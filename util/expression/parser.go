package expression

import (
	"fmt"
	ml_slice "github.com/maldan/go-ml/util/slice"
	"strconv"
	"unsafe"
)

func precedence(op TokenType) int {
	/*if op.Type == TokenBinaryExpression {
		return 18
	}
	*/
	switch op.Token {
	case "-":
	case "+":
		return 11
	case "*":
	case "/":
		return 12
	case "==":
		return 8
	case "&&":
		return 4
	default:
		return -1
	}
	return -1
}

func infixToPostfix(tokens []TokenType) []TokenType {
	postfix := make([]TokenType, 0)
	stack := make([]TokenType, 0)

	for i := 0; i < len(tokens); i++ {
		if tokens[i].Type == TokenOp {
			for {
				if len(stack) == 0 {
					break
				}

				if precedence(top(&stack)) < precedence(tokens[i]) {
					break
				}

				postfix = append(postfix, top(&stack))
				pop(&stack)
			}
			stack = append(stack, tokens[i])
		} else {
			postfix = append(postfix, tokens[i])
		}
	}

	for {
		if len(stack) == 0 {
			break
		}
		postfix = append(postfix, top(&stack))
		pop(&stack)
	}

	/*for i := 0; i < len(postfix); i++ {
		fmt.Printf("%v\n", postfix[i].Token)
	}*/
	fmt.Printf("X %v\n", ml_slice.Map(postfix, func(t TokenType) string {
		return "[" + t.Token + "]"
	}))
	// fmt.Printf("0: %v\n", postfix[0])
	return postfix
}

func pop[T any](s *[]T) T {
	v := (*s)[len(*s)-1]
	*s = (*s)[0 : len(*s)-1]
	return v
}

func top[T any](s *[]T) T {
	return (*s)[len(*s)-1]
}

/*func parseWhere(queryInfo *QueryInfo, tokens []TokenType) {
	// Change tokes
	for i := 0; i < len(tokens); i++ {
		if tokens[i].Token == "AND" {
			tokens[i].Token = "&&"
			tokens[i].Type = core.TokenOp
		}
	}

	queryInfo.WhereCondition = infixToPostfix(tokens)
}*/

func tokenizer(str string) []TokenType {
	out := make([]TokenType, 0)
	tempStr := ""
	tempNumber := ""
	mode := ""
	previousMode := ""
	isQuoteMode := false
	tempQuote := ""
	str += " "
	tempOp := ""

	for i := 0; i < len(str); i++ {
		if isQuoteMode {
			if str[i] == '\'' {
				isQuoteMode = false
				mode = ""
				continue
			}
			tempQuote += string(str[i])
			continue
		}

		if str[i] == '\'' {
			mode = "quote"
			isQuoteMode = true
		} else if str[i] == ' ' {
			mode = "space"
		} else if str[i] >= '0' && str[i] <= '9' {
			tempNumber += string(str[i])
			mode = "number"
		} else if str[i] == '(' {
			mode = "("
		} else if str[i] == ')' {
			mode = ")"
		} else if str[i] == '*' {
			mode = "*"
		} else if str[i] == '=' {
			tempOp += "="
			mode = "="
		} else if str[i] == '&' {
			tempOp += "&"
			mode = "&"
		} else {
			tempStr += string(str[i])
			mode = "string"
		}

		if mode != previousMode {
			if previousMode == "&" {
				out = append(out, TokenType{Token: tempOp, Type: TokenOp})
				tempOp = ""
			} else if previousMode == "=" {
				out = append(out, TokenType{Token: tempOp, Type: TokenOp})
				tempOp = ""
			} else if previousMode == ")" {
				out = append(out, TokenType{Token: ")", Type: TokenOp})
			} else if previousMode == "(" {
				out = append(out, TokenType{Token: "(", Type: TokenOp})
			} else if previousMode == "*" {
				out = append(out, TokenType{Token: "*", Type: TokenOp})
			} else if previousMode == "space" {

			} else if previousMode == "string" {
				out = append(out, TokenType{Token: tempStr, Type: TokenIdentifier})
				tempStr = ""
			} else if previousMode == "number" {
				out = append(out, TokenType{Token: tempNumber, Type: TokenNumber})
				tempNumber = ""
			} else if previousMode == "quote" {
				out = append(out, TokenType{Token: tempQuote, Type: TokenString})
				tempQuote = ""
			}

			previousMode = mode
		}
	}
	return out
}

func toAst(list []TokenType) []TokenType {
	for i := 0; i < len(list); i++ {
		if list[i].Type == TokenIdentifier && list[i+1].Type == TokenOp && list[i+2].Type == TokenNumber {
			x := infixToPostfix([]TokenType{list[i], list[i+1], list[i+2]})

			list[i] = TokenType{Type: TokenBinaryExpression, List: x}

			list = ml_slice.RemoveAt(list, i+1)
			list = ml_slice.RemoveAt(list, i+1)

			i -= 1
		}
	}

	return list
}

func flatten(list []TokenType) []TokenType {
	out := make([]TokenType, 0)

	for i := 0; i < len(list); i++ {
		if len(list[i].List) > 0 {
			for j := 0; j < len(list[i].List); j++ {
				out = append(out, list[i].List[j])
			}
		} else {
			out = append(out, list[i])
		}
	}

	return out
}

func Parse(query string) (Expression, error) {
	expr := Expression{
		Map:  map[string]unsafe.Pointer{},
		Vars: make([]any, 0, 64),
	}

	tokens := tokenizer(query)

	ast := toAst(tokens)
	x := infixToPostfix(ast)
	f := flatten(x)
	ff := make([]ExpressionToken, 0)
	// ml_console.PrettyPrint(f)

	for i := 0; i < len(f); i++ {
		v := any(f[i].Token)
		if f[i].Type == TokenNumber {
			n, _ := strconv.Atoi(f[i].Token)
			v = n
		}

		if f[i].Type == TokenNumber {
			expr.Vars = append(expr.Vars, v)
			anyPtr := unsafe.Pointer(&(expr.Vars[len(expr.Vars)-1]))
			iface := (*emptyInterface)(anyPtr)

			ff = append(ff, ExpressionToken{
				Type:    f[i].Type,
				Value:   v,
				Address: iface.ptr,
			})
		} else {
			ff = append(ff, ExpressionToken{
				Type:  f[i].Type,
				Value: v,
			})
		}
	}

	expr.TokenList = ff
	return expr, nil
}
