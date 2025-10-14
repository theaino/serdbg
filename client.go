package main

import (
	pb "serdbg/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type SerialConnection struct {
	conn *grpc.ClientConn
}

func NewSerialConnection() (connection *SerialConnection, err error) {
	connection = new(SerialConnection)
	connection.conn, err = grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return
	}
	pb.NewSerialClient(connection.conn)
	return
}
