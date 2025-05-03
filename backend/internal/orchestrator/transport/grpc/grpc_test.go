package grpc

import (
	"fmt"
	pb "github.com/distributed-calc/v1/pkg/proto/orchestrator"
	"github.com/distributed-calc/v1/test/mock"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"io"
	"testing"
	"time"
)

func TestServer_ProcessTasks(t *testing.T) {
	log, _ := zap.NewDevelopment()

	stream := mock.NewMockBidiServerStream[pb.TaskResult, pb.Task]()

	service := &mock.ServiceMock{}

	server := grpc.NewServer()

	go func() {
		defer func() {
			stream.RecvClosed = true
			close(stream.RecvCh)
		}()

		for i := range 3 {
			stream.RecvCh <- pb.TaskResult{
				Id:     fmt.Sprintf("test:recv:%d", i),
				Result: 1,
				Status: "completed",
				Final:  false,
			}

			stream.SetRecvErr(io.EOF)
		}
	}()

	go func() {
		defer func() {
			stream.SendClosed = true
			close(stream.SendCh)
		}()

		for i := range 3 {
			res, ok := <-stream.SendCh
			if !ok {
				return
			}

			fmt.Printf("test:send:%d, got from send ch: %v\n", i, res)
		}
	}()

	app := NewServer(&Config{
		SendTaskBackoff: 100 * time.Millisecond,
	}, server, log, service)

	err := app.ProcessTasks(stream)
	if err != nil {
		t.Fatalf("error processing tasks: %v", err)
	}
}
