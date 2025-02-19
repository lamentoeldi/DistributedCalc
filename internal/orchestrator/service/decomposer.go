package service

import (
	"DistributedCalc/pkg/models"
	"fmt"
	"math/rand"
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

func evaluateAST(node *node, ch chan *models.Task, rch chan *models.TaskResult) (float64, error) {
	if node.left == nil && node.right == nil {
		return strconv.ParseFloat(node.value.value, 64)
	}

	// All this goroutine bullshit may cause race, demands test

	chLeft := make(chan float64)
	chRight := make(chan float64)
	chErr := make(chan error)

	go func() {
		lRes, err := evaluateAST(node.left, ch, rch)
		if err != nil {
			chLeft <- 0
			chErr <- err
			return
		}

		chLeft <- lRes
		chErr <- nil
	}()

	go func() {
		rRes, err := evaluateAST(node.right, ch, rch)
		if err != nil {
			chRight <- 0
			chErr <- err
			return
		}

		chRight <- rRes
		chErr <- nil
	}()

	leftResult := <-chLeft
	rightResult := <-chRight

	select {
	case err := <-chErr:
		if err != nil {
			return 0, err
		}
	default:
	}

	// Pass this to channel, block until get result
	t := &models.Task{
		Id:            rand.Int(),
		Arg1:          leftResult,
		Arg2:          rightResult,
		Operation:     node.value.value,
		OperationTime: 0,
	}

	ch <- t
	i := <-rch
	return i.Result, nil
}

func walkAST(n *node) {
	if n.left.left == nil && n.left.right == nil || n.right.left == nil && n.right.right == nil {
		fmt.Println(n)
		return
	}

	walkAST(n.left)
	walkAST(n.right)
}
