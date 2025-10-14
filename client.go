package main

import (
	"context"
	"errors"
	pb "serdbg/proto"
	"time"

	"go.bug.st/serial"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"
)

type SerialConnection struct {
	conn *grpc.ClientConn
	client pb.SerialClient
}

func NewSerialConnection() (connection *SerialConnection, err error) {
	connection = new(SerialConnection)
	connection.conn, err = grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return
	}
	connection.client = pb.NewSerialClient(connection.conn)

	return
}

func (c *SerialConnection) Close() error {
	return c.conn.Close()
}

func (c *SerialConnection) GetMode() (mode *serial.Mode, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	rMode, err := c.client.GetMode(ctx, &emptypb.Empty{})
	if err != nil {
		return
	}
	mode = &serial.Mode{
		BaudRate: int(rMode.BaudRate),
		DataBits: int(rMode.DataBits),
		Parity: serial.Parity(rMode.Parity),
		StopBits: serial.StopBits(rMode.StopBits),
		InitialStatusBits: &serial.ModemOutputBits{
			RTS: rMode.Rts,
			DTR: rMode.Dtr,
		},
	}
	return
}

func (c *SerialConnection) SetMode(mode *serial.Mode) (err error) {
	rMode := &pb.Mode{
		BaudRate: int64(mode.BaudRate),
		DataBits: int64(mode.DataBits),
		Parity: int64(mode.Parity),
		StopBits: int64(mode.StopBits),
		Rts: mode.InitialStatusBits.RTS,
		Dtr: mode.InitialStatusBits.DTR,
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	_, err = c.client.SetMode(ctx, rMode)
	if err != nil {
		return
	}
	return
}

func (c *SerialConnection) GetPort() (port string, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	rPort, err := c.client.GetPort(ctx, &emptypb.Empty{})
	port = rPort.Port
	return
}

func (c *SerialConnection) SetPort(port string) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	rError, err := c.client.SetPort(ctx, &pb.Port{Port: port})
	if err != nil {
		return
	}
	if rError.Failed {
		err = errors.New(rError.Value)
	}
	return
}

func (c *SerialConnection) SendString(data string) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	rError, err := c.client.SendString(ctx, &pb.SendStringRequest{Data: data})
	if err != nil {
		return
	}
	if rError.Failed {
		err = errors.New(rError.Value)
	}
	return
}
