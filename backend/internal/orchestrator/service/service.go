package service

import (
	"context"
	"fmt"
	"github.com/distributed-calc/v1/internal/orchestrator/errors"
	"github.com/distributed-calc/v1/internal/orchestrator/models"
	"github.com/google/uuid"
	"go/ast"
	"go/parser"
	"strconv"
	"unicode"
)

const (
	StatusPending   = "pending"
	StatusCompleted = "completed"
	StatusFailed    = "failed"

	number      = "NUMBER"
	operator    = "OPERATOR"
	parenthesis = "PARENTHESIS"
)

type ExpRepo interface {
	Add(ctx context.Context, exp *models.Expression) error
	Get(ctx context.Context, id string) (*models.Expression, error)
	GetAll(ctx context.Context) ([]*models.Expression, error)
	Update(ctx context.Context, exp *models.Expression) error
}

type TaskRepo interface {
	AddTasks(ctx context.Context, tasks []*models.Task) error
	GetTask(ctx context.Context) (*models.Task, error)
	UpdateTask(ctx context.Context, task *models.Task) error
	DeleteTasks(ctx context.Context, expID string) error
}

type Service struct {
	expRepo  ExpRepo
	taskRepo TaskRepo
}

func NewService(expRepo ExpRepo, taskRepo TaskRepo) *Service {
	return &Service{
		expRepo:  expRepo,
		taskRepo: taskRepo,
	}
}

func (s *Service) Evaluate(ctx context.Context, expression string) (string, error) {
	err := validate(expression)
	if err != nil {
		return "", err
	}

	expID, _ := uuid.NewV7()

	exp := &models.Expression{
		Id:     expID.String(),
		Status: StatusPending,
		Result: 0,
	}

	tasks, err := parseExpression(expression, expID.String())

	err = s.expRepo.Add(ctx, exp)
	if err != nil {
		return "", err
	}

	err = s.taskRepo.AddTasks(ctx, tasks)
	if err != nil {
		return "", err
	}

	return expID.String(), nil
}

func (s *Service) Get(ctx context.Context, id string) (*models.Expression, error) {
	return s.expRepo.Get(ctx, id)
}

func (s *Service) GetAll(ctx context.Context) ([]*models.Expression, error) {
	return s.expRepo.GetAll(ctx)
}

func (s *Service) GetTask(ctx context.Context) (*models.AgentTask, error) {
	task, err := s.taskRepo.GetTask(ctx)
	if err != nil {
		return nil, err
	}

	return &models.AgentTask{
		Id:            task.ID,
		LeftArg:       task.LeftArg,
		RightArg:      task.RightArg,
		Op:            task.Op,
		OperationTime: 0,
		Final:         task.Final,
	}, nil
}

func (s *Service) FinishTask(ctx context.Context, task *models.TaskResult) error {
	err := s.taskRepo.UpdateTask(ctx, &models.Task{
		ID:     task.Id,
		Result: task.Result,
		Status: task.Status,
		Final:  task.Final,
	})
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) Finalize(ctx context.Context, expID string, result float64) error {
	exp := &models.Expression{
		Id:     expID,
		Status: StatusCompleted,
		Result: result,
	}

	err := s.expRepo.Update(ctx, exp)
	if err != nil {
		return err
	}

	err = s.taskRepo.DeleteTasks(ctx, expID)
	if err != nil {
		return err
	}

	return nil
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

func validate(expression string) error {
	var currentToken token
	var dotEncountered bool
	var parenthesisCount int
	var previousRune rune

	if len(expression) < 1 {
		return fmt.Errorf("%w: expression is empty", errors.ErrInvalidExpression)
	}

	for i, r := range expression {
		if unicode.IsSpace(r) {
			continue
		}

		// Prevent '*2+3' or '2+3*'
		if i == 0 && isOperator(r) {
			if expression[0] != '-' && expression[0] != '+' {
				return fmt.Errorf("%w: unexpected operator at the beginning of expression", errors.ErrInvalidExpression)
			}
			currentToken.tokenType = number
			currentToken.value += string(expression[0])
			continue
		}

		if i == len(expression)-1 && isOperator(r) {
			return fmt.Errorf("%w: unexpected operator at the end of expression", errors.ErrInvalidExpression)
		}

		if unicode.IsDigit(r) {
			currentToken.tokenType = number
			currentToken.value += string(r)
		} else if r == '.' {
			// Prevent '3.14.2'
			if dotEncountered {
				return fmt.Errorf("%w: multiple decimal points in the same number", errors.ErrInvalidExpression)
			}
			currentToken.tokenType = number
			currentToken.value += string(r)
			dotEncountered = true
		} else {
			if currentToken.value != "" {
				currentToken = token{}
			}
			if isOperator(r) {
				if isOperator(previousRune) {
					return fmt.Errorf("%w: multiple sequent operators", errors.ErrInvalidExpression)
				}

				currentToken.tokenType = operator
				currentToken.value = string(r)

				currentToken.value = ""
			} else if isParenthesis(r) {
				// Prevent '2+()+3'
				if isParenthesis(r) && isParenthesis(previousRune) && r != previousRune {
					return fmt.Errorf("%w: empty parentheses", errors.ErrInvalidExpression)
				}

				currentToken.tokenType = parenthesis
				currentToken.value = string(r)

				currentToken.value = ""

				if r == '(' {
					parenthesisCount++
				} else {
					parenthesisCount--
				}
			} else {
				return fmt.Errorf("%w: invalid character: %c", errors.ErrInvalidExpression, r)
			}
		}

		previousRune = r
	}

	if parenthesisCount != 0 {
		return fmt.Errorf("%w: invalid parenthesis", errors.ErrInvalidExpression)
	}

	return nil
}

func parseExpression(exprStr, expID string) ([]*models.Task, error) {
	expr, err := parser.ParseExpr(exprStr)
	if err != nil {
		return nil, err
	}

	tasks := make([]*models.Task, 0)
	taskID := 1

	var visit func(expr ast.Expr) *models.Task
	visit = func(expr ast.Expr) *models.Task {
		var t *models.Task
		switch e := expr.(type) {
		case *ast.BasicLit:
			val, err := strconv.ParseFloat(e.Value, 64)
			if err != nil {
				panic(err)
			}
			t = &models.Task{
				ID:      fmt.Sprintf("%s:%d", expID, taskID),
				ExpID:   expID,
				LeftArg: val,
				Status:  "ready",
			}
			taskID++

		case *ast.BinaryExpr:
			leftTask := visit(e.X)
			rightTask := visit(e.Y)

			t = &models.Task{
				ID:      fmt.Sprintf("%s:%d", expID, taskID),
				ExpID:   expID,
				Op:      e.Op.String(),
				LeftID:  &leftTask.ID,
				RightID: &rightTask.ID,
			}
			taskID++

		case *ast.ParenExpr:
			return visit(e.X)
		}

		if t.LeftID == nil && t.RightID == nil {
			t.Status = "ready"
		}

		tasks = append(tasks, t)
		return t
	}

	visit(expr)

	tasks[len(tasks)-1].Final = true
	return tasks, nil
}
