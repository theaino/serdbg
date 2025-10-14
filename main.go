package main

import (
	"flag"
	"os"
	"serdbg/server"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	server := flag.Bool("server", false, "Run as serial server")
	
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
	file, err := os.Open(os.Args[1])
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
