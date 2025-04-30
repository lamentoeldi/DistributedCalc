package service

import (
	"context"
	"fmt"
	errors2 "github.com/distributed-calc/v1/internal/orchestrator/errors"
	models2 "github.com/distributed-calc/v1/internal/orchestrator/models"
	"github.com/distributed-calc/v1/pkg/authenticator"
	"github.com/distributed-calc/v1/pkg/models"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"strconv"
	"unicode"
)

const (
	StatusPending   = "pending"
	StatusCompleted = "completed"
	StatusFailed    = "failed"
)

type task struct {
	id       string
	leftID   string
	rightID  string
	arg1     *float64
	arg2     *float64
	operator string
	result   *float64
	topLevel bool
}

type TaskPlanner interface {
	PlanTask(ctx context.Context, task *models.Task) (TaskPromise, error)
	FinishTask(ctx context.Context, res *models.TaskResult) error
}

type TaskPromise interface {
	WaitForResult(ctx context.Context) (float64, error)
}

type Repository interface {
	Add(ctx context.Context, exp *models.Expression) error
	Get(ctx context.Context, id string) (*models.Expression, error)
	GetAll(ctx context.Context) ([]*models.Expression, error)
	AddUser(ctx context.Context, user *models2.User) error
	GetUser(ctx context.Context, login string) (*models2.User, error)
}

type Queue[T any] interface {
	Enqueue(ctx context.Context, obj *T) error
	Dequeue(ctx context.Context) (*T, error)
}

type Service struct {
	repo Repository
	p    TaskPlanner
	q    Queue[models.Task]
	auth *authenticator.Authenticator
}

func NewService(r Repository, p TaskPlanner, q Queue[models.Task], auth *authenticator.Authenticator) *Service {
	return &Service{
		repo: r,
		p:    p,
		q:    q,
		auth: auth,
	}
}

func (s *Service) StartEvaluation(ctx context.Context, expression string) (string, error) {
	// In blocking part, initial validation is performed
	tokens, err := s.tokenize(expression)
	if err != nil {
		return "", err
	}

	ast, err := s.buildAST(tokens)
	if err != nil {
		return "", err
	}

	id := uuid.NewString()

	exp := &models.Expression{
		Id:     id,
		Status: StatusPending,
		Result: 0,
	}

	err = s.repo.Add(context.TODO(), exp)
	if err != nil {
		return "", err
	}

	// In non-blocking part, further calculation is initiated
	go func(exp *models.Expression) {
		res, err := s.taskAST(ast)
		if err != nil {
			exp.Status = StatusFailed
			_ = s.repo.Add(ctx, exp)
			return
		}

		exp.Status = StatusCompleted

		// Throw away excessive decimal places
		// It is ok ignoring this error since it may occur only if call above failed
		//res, _ = strconv.ParseFloat(fmt.Sprintf("%.8f", res), 64)
		//exp.Result = res
		//_ = s.repo.Add(ctx, exp)

		// TODO: Task creation
		_ = res
	}(exp)

	return id, nil
}

func (s *Service) Get(ctx context.Context, id string) (*models.Expression, error) {
	return s.repo.Get(ctx, id)
}

func (s *Service) GetAll(ctx context.Context) ([]*models.Expression, error) {
	return s.repo.GetAll(ctx)
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

func (s *Service) taskAST(n *node) ([]*task, error) {
	tasks := make([]*task, 0)

	var checkNode func(node *node) (*task, error)
	checkNode = func(in *node) (*task, error) {
		if in.left == nil && in.right == nil {
			taskID := uuid.NewString()

			res, err := strconv.ParseFloat(in.value.value, 64)
			if err != nil {
				return nil, err
			}

			t := &task{
				id:       taskID,
				operator: "none",
				result:   &res,
			}

			tasks = append(tasks, t)

			return t, nil
		}

		leftID, err := checkNode(in.left)
		if err != nil {
			return nil, err
		}

		rightID, err := checkNode(in.right)
		if err != nil {
			return nil, err
		}

		taskID := uuid.NewString()

		t := &task{
			id:       taskID,
			leftID:   leftID.id,
			rightID:  rightID.id,
			operator: in.value.value,
		}

		if leftID.operator == "none" {
			t.arg1 = leftID.result
		}

		if rightID.operator == "none" {
			t.arg2 = rightID.result
		}

		tasks = append(tasks, t)

		return t, nil
	}

	_, err := checkNode(n)
	if err != nil {
		return nil, err
	}

	// Mark the top level node (it is always the last node)
	last := len(tasks) - 1
	tasks[last].topLevel = true

	return tasks, nil
}

func (s *Service) Register(ctx context.Context, creds *models2.UserCredentials) error {
	id, _ := uuid.NewV7()

	hash, err := bcrypt.GenerateFromPassword([]byte(creds.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to register user: %w", err)
	}

	user := &models2.User{
		Id:             id.String(),
		Username:       creds.Username,
		HashedPassword: hash,
	}

	err = s.repo.AddUser(ctx, user)
	if err != nil {
		return fmt.Errorf("failed to register user: %w", err)
	}

	return nil
}

func (s *Service) Login(ctx context.Context, creds *models2.UserCredentials) (*models2.JWTTokens, error) {
	user, err := s.repo.GetUser(ctx, creds.Username)
	if err != nil {
		return nil, fmt.Errorf("failed to authorize user: %w: %w", errors2.ErrUnauthorized, err)
	}

	err = bcrypt.CompareHashAndPassword(user.HashedPassword, []byte(creds.Password))
	if err != nil {
		return nil, fmt.Errorf("failed to authorize user: %w: %w", errors2.ErrUnauthorized, err)
	}

	access, refresh, err := s.auth.SignTokens(s.auth.IssueTokens(user.Id))
	if err != nil {
		return nil, fmt.Errorf("failed to issue tokens: %w", err)
	}

	return &models2.JWTTokens{
		AccessToken:  access,
		RefreshToken: refresh,
	}, nil
}
