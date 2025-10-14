package server

import (
	"context"
	"fmt"
	"net"
	pb "serdbg/proto"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type server struct {
	pb.UnimplementedSerialServer
}

func (s *server) SendString(ctx context.Context, req *pb.SendStringRequest) (*emptypb.Empty, error) {
	fmt.Println("Test")
	return &emptypb.Empty{}, nil
}

func RunServer() {
	lis, _ := net.Listen("tcp", ":50051")
	s := grpc.NewServer()
	pb.RegisterSerialServer(s, &server{})
	s.Serve(lis)
}
