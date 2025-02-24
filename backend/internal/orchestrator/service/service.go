package service

import (
	"DistributedCalc/pkg/models"
	"context"
	"fmt"
	"math/rand"
	"strconv"
	"unicode"
)

const (
	StatusPending   = "pending"
	StatusCompleted = "completed"
	StatusFailed    = "failed"
)

type TaskPlanner interface {
	PlanTask(ctx context.Context, task *models.Task) (TaskPromise, error)
	FinishTask(ctx context.Context, res *models.TaskResult) error
}

type TaskPromise interface {
	WaitForResult(ctx context.Context) (float64, error)
}

type Repository interface {
	Add(ctx context.Context, exp *models.Expression) error
	Get(ctx context.Context, id int) (*models.Expression, error)
	GetAll(ctx context.Context) ([]*models.Expression, error)
}

type Queue[T any] interface {
	Enqueue(ctx context.Context, obj *T) error
	Dequeue(ctx context.Context) (*T, error)
}

type Service struct {
	r Repository
	p TaskPlanner
	q Queue[models.Task]
}

func NewService(r Repository, p TaskPlanner, q Queue[models.Task]) *Service {
	return &Service{
		r: r,
		p: p,
		q: q,
	}
}

func (s *Service) StartEvaluation(ctx context.Context, expression string) (int, error) {
	// In blocking part, initial validation is performed
	tokens, err := s.tokenize(expression)
	if err != nil {
		return 0, err
	}

	ast, err := s.buildAST(tokens)
	if err != nil {
		return 0, err
	}

	id := rand.Int()

	exp := &models.Expression{
		Id:     id,
		Status: StatusPending,
		Result: 0,
	}

	err = s.r.Add(context.TODO(), exp)
	if err != nil {
		return 0, err
	}

	// In non-blocking part, further calculation is initiated
	go func(exp *models.Expression) {
		res, err := s.evaluateAST(ast)
		if err != nil {
			exp.Status = StatusFailed
			_ = s.r.Add(ctx, exp)
			return
		}

		exp.Status = StatusCompleted

		// Throw away excessive decimal places
		// It is ok ignoring this error since it may occur only if call above failed
		res, _ = strconv.ParseFloat(fmt.Sprintf("%.8f", res), 64)
		exp.Result = res
		_ = s.r.Add(ctx, exp)
	}(exp)

	return id, nil
}

func (s *Service) Get(ctx context.Context, id int) (*models.Expression, error) {
	return s.r.Get(ctx, id)
}

func (s *Service) GetAll(ctx context.Context) ([]*models.Expression, error) {
	return s.r.GetAll(ctx)
}

func (s *Service) FinishTask(ctx context.Context, result *models.TaskResult) error {
	return s.p.FinishTask(ctx, result)
}

func (s *Service) Enqueue(ctx context.Context, task *models.Task) error {
	return s.q.Enqueue(ctx, task)
}

func (s *Service) Dequeue(ctx context.Context) (*models.Task, error) {
	return s.q.Dequeue(ctx)
}

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

func (s *Service) tokenize(expression string) ([]token, error) {
	var tokens []token
	var currentToken token
	var dotEncountered bool
	var parenthesisCount int
	var previousRune rune

	if len(expression) < 1 {
		return nil, fmt.Errorf("%w: expression is empty", models.ErrInvalidExpression)
	}

	for i, r := range expression {
		if unicode.IsSpace(r) {
			continue
		}

		// Prevent '*2+3' or '2+3*'
		if i == 0 && isOperator(r) {
			if expression[0] != '-' && expression[0] != '+' {
				return nil, fmt.Errorf("%w: unexpected operator at the beginning of expression", models.ErrInvalidExpression)
			}
			currentToken.tokenType = number
			currentToken.value += string(expression[0])
			continue
		}

		if i == len(expression)-1 && isOperator(r) {
			return nil, fmt.Errorf("%w: unexpected operator at the end of expression", models.ErrInvalidExpression)
		}

		if unicode.IsDigit(r) {
			currentToken.tokenType = number
			currentToken.value += string(r)
		} else if r == '.' {
			// Prevent '3.14.2'
			if dotEncountered {
				return nil, fmt.Errorf("%w: multiple decimal points in the same number", models.ErrInvalidExpression)
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
					return nil, fmt.Errorf("%w: multiple sequent operators", models.ErrInvalidExpression)
				}

				currentToken.tokenType = operator
				currentToken.value = string(r)

				tokens = append(tokens, currentToken)
				currentToken.value = ""
			} else if isParenthesis(r) {
				// Prevent '2+()+3'
				if isParenthesis(r) && isParenthesis(previousRune) && r != previousRune {
					return nil, fmt.Errorf("%w: empty parentheses", models.ErrInvalidExpression)
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
				return nil, fmt.Errorf("%w: invalid character: %c", models.ErrInvalidExpression, r)
			}
		}

		previousRune = r
	}

	if currentToken.value != "" {
		tokens = append(tokens, currentToken)
	}

	if parenthesisCount != 0 {
		return nil, fmt.Errorf("%w: invalid parenthesis", models.ErrInvalidExpression)
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

func (s *Service) buildAST(tokens []token) (*node, error) {
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

	return nil, models.ErrInvalidExpression
}

func (s *Service) evaluateAST(node *node) (float64, error) {
	if node.left == nil && node.right == nil {
		return strconv.ParseFloat(node.value.value, 64)
	}

	// All this goroutine bullshit may cause race, demands test

	chLeft := make(chan float64)
	chRight := make(chan float64)
	chErr := make(chan error)

	go func() {
		lRes, err := s.evaluateAST(node.left)
		if err != nil {
			chLeft <- 0
			chErr <- err
			return
		}

		chLeft <- lRes
		chErr <- nil
	}()

	go func() {
		rRes, err := s.evaluateAST(node.right)
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
		Arg1:      leftResult,
		Arg2:      rightResult,
		Operation: node.value.value,
	}

	p, err := s.p.PlanTask(context.TODO(), t)
	if err != nil {
		return 0, err
	}

	return p.WaitForResult(context.TODO())
}
