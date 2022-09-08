package cmd

import (
	// "time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"net"

	"google.golang.org/grpc"

	pb "github.com/LieAlbertTriAdrian/clean-arch-golang/domain/proto"
	todoGrpc "github.com/LieAlbertTriAdrian/clean-arch-golang/todo/delivery/grpc"
)

var grpcCommand = &cobra.Command{
	Use:   "grpc",
	Short: "Start gRPC server",
	Run:   grpcServer,
}

func init() {
	rootCmd.AddCommand(grpcCommand)
}

func grpcServer(cmd *cobra.Command, args []string) {
	lis, err := net.Listen("tcp", ":8081")
	if err != nil {
		logrus.Fatalf("failed to listen: %v", err)
	}

	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	todoGrpcServer, _ := todoGrpc.New()

	pb.RegisterTodoRpcServer(grpcServer, todoGrpcServer)

	logrus.Printf("server listening at %v", lis.Addr())
	if err := grpcServer.Serve(lis); err != nil {
		logrus.Fatalf("failed to serve: %v", err)
	}
}
