package grpc

import (
	"context"

	pb "github.com/LieAlbertTriAdrian/clean-arch-golang/domain/proto"
)

// handler is the interface which exposes the User Server methods
type server struct {
	pb.UnimplementedTodoRpcServer
}

// New returns the object for the RPC handler
func New() (pb.TodoRpcServer, error) {
	return &server{}, nil
}

// FetchTodos function returns the list of users
func (h *server) FetchTodos(ctx context.Context, r *pb.EmptyReq) (*pb.FetchTodosResponse, error) {
	return &pb.FetchTodosResponse{
		Todos: []*pb.Todo{
			{
				Id:      "id",
				Message: "message",
				Status:  "status",
			},
		},
	}, nil
}
