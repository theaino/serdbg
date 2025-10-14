package main

import (
	"flag"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	flag.Parse()
	
	runClient()
}

func runClient() {
	file, err := os.Open(flag.Arg(0))
	if err != nil {
		panic(err)
	}
	model, err := NewModel(file)
	if err != nil {
		panic(err)
	}
	program := tea.NewProgram(model)
	if _, err := program.Run(); err != nil {
		panic(err)
	}
}
