package main

import (
	"flag"
	"fmt"
	"os"
	"serdbg/server"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	server := flag.Bool("server", false, "Run as serial server")
	flag.Parse()
	
	if *server {
		runServer()
	} else {
		runClient()
	}
}

func runServer() {
	server.RunServer()
}

func runClient() {
	serial, err := NewSerialConnection()
	if err != nil {
		panic(err)
	}
	err = serial.SendString("Hello")
	if err != nil {
		panic(err)
	}
	port, err := serial.GetPort()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Port: %s\n", port)
	err = serial.SetPort("/dev/ttyUSB0")
	if err != nil {
		panic(err)
	}
	port, err = serial.GetPort()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Port: %s\n", port)
	serial.Close()

	file, err := os.Open(flag.Arg(0))
	if err != nil {
		panic(err)
	}
	model, err := newModel(file)
	if err != nil {
		panic(err)
	}
	program := tea.NewProgram(model)
	if _, err := program.Run(); err != nil {
		panic(err)
	}
}
