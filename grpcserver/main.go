package grpcserver

import (
	"fmt"
	"net"

	"google.golang.org/grpc"
)

func StartGRPCServer(server *grpc.Server, port string) error {
	listener, err := net.Listen("tcp", port)
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}

	fmt.Println("gRPC server started on", port)

	return server.Serve(listener)
}
