package mock

import (
	"context"
	"fmt"
	ma "github.com/distributed-calc/v1/internal/agent/models"
	mo "github.com/distributed-calc/v1/internal/orchestrator/models"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"io"
)

type OrchestratorMock struct {
	Err error
}

func (o *OrchestratorMock) GetTask(_ context.Context) (*ma.AgentTask, error) {
	if o.Err != nil {
		return nil, o.Err
	}

	return &ma.AgentTask{
		Id:       fmt.Sprint(10),
		Op:       "+",
		LeftArg:  2.0,
		RightArg: 3.0,
	}, nil
}

func (o *OrchestratorMock) PostResult(_ context.Context, _ *ma.TaskResult) error {
	if o.Err != nil {
		return o.Err
	}

	return nil
}

type CalculatorMock struct {
	Err error
}

func (c *CalculatorMock) Evaluate(task *ma.AgentTask) (*ma.TaskResult, error) {
	fmt.Println("evaluate called: ", task)

	if c.Err != nil {
		return nil, c.Err
	}

	return &ma.TaskResult{
		Id:     fmt.Sprint(10),
		Result: 5.0,
	}, nil
}

type ServiceMock struct {
	Err error
}

func (s ServiceMock) Register(ctx context.Context, creds *mo.UserCredentials) error {
	return s.Err
}

func (s ServiceMock) Login(_ context.Context, _ *mo.UserCredentials) (*mo.JWTTokens, error) {
	return &mo.JWTTokens{
		AccessToken:  "",
		RefreshToken: "",
	}, s.Err
}

func (s ServiceMock) GetUserID(_ context.Context, _ string) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (s ServiceMock) VerifyJWT(_ context.Context, _ string) error {
	//TODO implement me
	panic("implement me")
}

func (s ServiceMock) RefreshTokens(_ context.Context, _ string) (string, string, error) {
	//TODO implement me
	panic("implement me")
}

func (s ServiceMock) Evaluate(_ context.Context, _, _ string) (string, error) {
	if s.Err != nil {
		return "", s.Err
	}

	id, _ := uuid.NewV7()

	return id.String(), nil
}

func (s ServiceMock) Get(_ context.Context, _, _ string) (*mo.Expression, error) {
	if s.Err != nil {
		return nil, s.Err
	}

	return &mo.Expression{
		Id:     fmt.Sprint(10),
		Result: 0.0,
		Status: "pending",
	}, nil
}

func (s ServiceMock) GetAll(_ context.Context, _, _ string, _ int64) ([]*mo.Expression, error) {
	if s.Err != nil {
		return nil, s.Err
	}

	return []*mo.Expression{
		{
			Id:     fmt.Sprint(10),
			Result: 0.0,
			Status: "pending",
		},
	}, nil
}

func (s ServiceMock) GetTask(_ context.Context) (*mo.AgentTask, error) {
	if s.Err != nil {
		return nil, s.Err
	}

	return &mo.AgentTask{
		Id:            fmt.Sprint(10),
		LeftArg:       10,
		RightArg:      10,
		Op:            "+",
		OperationTime: 0,
		Final:         true,
	}, nil
}

func (s ServiceMock) FinishTask(_ context.Context, _ *mo.TaskResult) error {
	if s.Err != nil {
		return s.Err
	}

	return nil
}

func (s ServiceMock) Finalize(_ context.Context, _ string, _ float64) error {
	if s.Err != nil {
		return s.Err
	}

	return nil
}

type BidiServerStream[Req, Res any] struct {
	grpc.ServerStream
	RecvCh     chan Req
	SendCh     chan Res
	recvErr    error
	sendErr    error
	SendClosed bool
	RecvClosed bool
}

func NewMockBidiServerStream[Req, Res any]() *BidiServerStream[Req, Res] {
	return &BidiServerStream[Req, Res]{
		RecvCh: make(chan Req),
		SendCh: make(chan Res),
	}
}

func (b *BidiServerStream[Req, Res]) SetSendErr(err error) {
	b.sendErr = err
}

func (b *BidiServerStream[Req, Res]) SetRecvErr(err error) {
	b.recvErr = err
}

func (b *BidiServerStream[Req, Res]) Recv() (*Req, error) {
	req := <-b.RecvCh

	if b.recvErr != nil {
		return nil, b.recvErr
	}

	return &req, nil
}

func (b *BidiServerStream[Req, Res]) Send(res *Res) error {
	if b.SendClosed {
		return nil
	}

	b.SendCh <- *res

	if b.sendErr != nil {
		return b.sendErr
	}

	return nil
}

func (b *BidiServerStream[Req, Res]) Context() context.Context {
	return context.Background()
}

type BidiClientStream[Req, Res any] struct {
	grpc.ServerStream
	RecvCh     chan Res
	SendCh     chan Req
	recvErr    error
	sendErr    error
	SendClosed bool
	RecvClosed bool
}

func NewMockBidiClientStream[Req, Res any]() *BidiClientStream[Req, Res] {
	return &BidiClientStream[Req, Res]{
		RecvCh: make(chan Res),
		SendCh: make(chan Req),
	}
}

func (b *BidiClientStream[Req, Res]) SetSendErr(err error) {
	b.sendErr = err
}

func (b *BidiClientStream[Req, Res]) SetRecvErr(err error) {
	b.recvErr = err
}

func (b *BidiClientStream[Req, Res]) Send(req *Req) error {
	if b.SendClosed {
		return nil
	}

	b.SendCh <- *req

	if b.sendErr != nil {
		return b.sendErr
	}

	return nil
}

func (b *BidiClientStream[Req, Res]) Recv() (*Res, error) {
	req := <-b.RecvCh

	if b.recvErr != nil {
		return nil, b.recvErr
	}

	return &req, nil
}

func (b *BidiClientStream[Req, Res]) Header() (metadata.MD, error) {
	return nil, nil
}

func (b *BidiClientStream[Req, Res]) Trailer() metadata.MD {
	return nil
}

func (b *BidiClientStream[Req, Res]) CloseSend() error {
	b.SetRecvErr(io.EOF)
	return nil
}
