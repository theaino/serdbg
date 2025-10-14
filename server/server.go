package server

import (
	"context"
	"fmt"
	"net"
	pb "serdbg/proto"

	"go.bug.st/serial"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type server struct {
	pb.UnimplementedSerialServer
	SerialPort string
	SerialMode *serial.Mode
}

func newServer() *server {
	return &server{
		SerialMode: &serial.Mode{
			InitialStatusBits: &serial.ModemOutputBits{},
		},
	}
}

func (s *server) GetMode(ctx context.Context, req *emptypb.Empty) (*pb.Mode, error) {
	return &pb.Mode{
		BaudRate: int64(s.SerialMode.BaudRate),
		DataBits: int64(s.SerialMode.DataBits),
		Parity: int64(s.SerialMode.Parity),
		StopBits: int64(s.SerialMode.StopBits),
		Rts: s.SerialMode.InitialStatusBits.RTS,
		Dtr: s.SerialMode.InitialStatusBits.DTR,
	}, nil
}

func (s *server) SetMode(ctx context.Context, req *pb.Mode) (*emptypb.Empty, error) {
	s.SerialMode = &serial.Mode{
		BaudRate: int(req.BaudRate),
		DataBits: int(req.DataBits),
		Parity: serial.Parity(req.Parity),
		StopBits: serial.StopBits(req.StopBits),
		InitialStatusBits: &serial.ModemOutputBits{
			RTS: req.Rts,
			DTR: req.Dtr,
		},
	}
	return &emptypb.Empty{}, nil
}

func (s *server) GetPort(ctx context.Context, req *emptypb.Empty) (*pb.Port, error) {
	return &pb.Port{Port: s.SerialPort}, nil
}

func (s *server) SetPort(ctx context.Context, req *pb.Port) (*pb.Error, error) {
	// TODO: implement checking
	s.SerialPort = req.Port
	return &pb.Error{Failed: false}, nil
}

func (s *server) SendString(ctx context.Context, req *pb.SendStringRequest) (*pb.Error, error) {
	fmt.Println(req.Data)
	return &pb.Error{Failed: false}, nil
}

func RunServer() {
	lis, _ := net.Listen("tcp", ":50051")
	s := grpc.NewServer()
	pb.RegisterSerialServer(s, newServer())
	s.Serve(lis)
}
