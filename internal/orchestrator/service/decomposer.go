package service

import (
	"fmt"
	"strconv"
	"unicode"
)

const (
	number      = "NUMBER"
	operator    = "OPERATOR"
	parenthesis = "PARENTHESIS"
)

type node struct {
	left  *node
	right *node
	value token
}

type AST node

type token struct {
	tokenType string
	value     string
}

func isOperator(r rune) bool {
	return r == '+' || r == '-' || r == '*' || r == '/'
}

func isParenthesis(r rune) bool {
	return r == '(' || r == ')'
}

func tokenize(expression string) ([]token, error) {
	var tokens []token
	var currentToken token
	var dotEncountered bool
	var parenthesisCount int
	var previousRune rune

	for i, r := range expression {
		// Prevent '*2+3' or '2+3*'
		if (i == 0 || i == len(expression)-1) && isOperator(r) {
			return nil, fmt.Errorf("unexpected operator at the beginning or end of expression: %c", r)
		}

		if unicode.IsDigit(r) {
			currentToken.tokenType = number
			currentToken.value += string(r)
		} else if r == '.' {
			// Prevent '3.14.2'
			if dotEncountered {
				return nil, fmt.Errorf("multiple decimal points in the same number: %s", expression)
			}
			currentToken.tokenType = number
			currentToken.value += string(r)
			dotEncountered = true
		} else {
			if currentToken.value != "" {
				tokens = append(tokens, currentToken)
				currentToken = token{}
			}
			if isOperator(r) {
				if isOperator(previousRune) {
					return nil, fmt.Errorf("multiple operators in a row: %c", previousRune)
				}

				currentToken.tokenType = operator
				currentToken.value = string(r)

				tokens = append(tokens, currentToken)
				currentToken.value = ""
			} else if isParenthesis(r) {
				// Prevent '2+()+3'
				if isParenthesis(r) && isParenthesis(previousRune) && r != previousRune {
					return nil, fmt.Errorf("empty parentheses: %c", previousRune)
				}

				currentToken.tokenType = parenthesis
				currentToken.value = string(r)

				tokens = append(tokens, currentToken)
				currentToken.value = ""

				if r == '(' {
					parenthesisCount++
				} else {
					parenthesisCount--
				}
			} else {
				return nil, fmt.Errorf("invalid character: %c", r)
			}
		}

		previousRune = r
	}

	if currentToken.value != "" {
		tokens = append(tokens, currentToken)
	}

	if parenthesisCount != 0 {
		return nil, fmt.Errorf("parenthesis count mismatch")
	}

	return tokens, nil
}

func hasHigherPrecedence(op1, op2 token) bool {
	precedence := map[string]int{
		"+": 1,
		"-": 1,
		"*": 2,
		"/": 2,
	}

	return precedence[op1.value] > precedence[op2.value]
}

func buildAST(tokens []token) (*node, error) {
	var operatorStack []token
	var operandStack []*node

	for _, t := range tokens {
		switch t.tokenType {
		case number:
			operandStack = append(operandStack, &node{value: t})

		case operator:
			for len(operatorStack) > 0 && hasHigherPrecedence(operatorStack[len(operatorStack)-1], t) {
				op := operatorStack[len(operatorStack)-1]
				operatorStack = operatorStack[:len(operatorStack)-1]

				right := operandStack[len(operandStack)-1]
				operandStack = operandStack[:len(operandStack)-1]
				left := operandStack[len(operandStack)-1]
				operandStack = operandStack[:len(operandStack)-1]

				operandStack = append(operandStack, &node{
					left:  left,
					right: right,
					value: op,
				})
			}
			operatorStack = append(operatorStack, t)

		case parenthesis:
			if t.value == "(" {
				operatorStack = append(operatorStack, t)
			} else {
				for len(operatorStack) > 0 && operatorStack[len(operatorStack)-1].value != "(" {
					op := operatorStack[len(operatorStack)-1]
					operatorStack = operatorStack[:len(operatorStack)-1]

					right := operandStack[len(operandStack)-1]
					operandStack = operandStack[:len(operandStack)-1]
					left := operandStack[len(operandStack)-1]
					operandStack = operandStack[:len(operandStack)-1]

					operandStack = append(operandStack, &node{
						left:  left,
						right: right,
						value: op,
					})
				}
				operatorStack = operatorStack[:len(operatorStack)-1]
			}
		}
	}

	for len(operatorStack) > 0 {
		op := operatorStack[len(operatorStack)-1]
		operatorStack = operatorStack[:len(operatorStack)-1]

		right := operandStack[len(operandStack)-1]
		operandStack = operandStack[:len(operandStack)-1]
		left := operandStack[len(operandStack)-1]
		operandStack = operandStack[:len(operandStack)-1]

		operandStack = append(operandStack, &node{
			left:  left,
			right: right,
			value: op,
		})
	}

	if len(operandStack) == 1 {
		return operandStack[0], nil
	}

	return nil, fmt.Errorf("invalid expression")
}

// This implementation is used for tests only, needed to check that AST is built correctly
func evaluateAST(node *node) (float64, error) {
	if node.left == nil && node.right == nil {
		return strconv.ParseFloat(node.value.value, 64)
	}

	var leftResult, rightResult float64
	var leftErr, rightErr error

	leftResult, leftErr = evaluateAST(node.left)
	if leftErr != nil {
		return 0, leftErr
	}

	rightResult, rightErr = evaluateAST(node.right)
	if rightErr != nil {
		return 0, rightErr
	}

	switch node.value.value {
	case "+":
		return leftResult + rightResult, nil
	case "-":
		return leftResult - rightResult, nil
	case "*":
		return leftResult * rightResult, nil
	case "/":
		if rightResult == 0 {
			return 0, fmt.Errorf("division by zero")
		}
		return leftResult / rightResult, nil
	default:
		return 0, fmt.Errorf("unknown operator %s", node.value.value)
	}
}
