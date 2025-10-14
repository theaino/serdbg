package main

import (
	"io"
	"os"
	"strconv"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	hpglContext *HPGLContext
	width, height int
	numBuffer string
	instructions []string
	instructionPointer int
}

func newModel(file *os.File) (model model, err error) {
	model.hpglContext = NewHPGLContext()

	source, err := io.ReadAll(file)
	if err != nil {
		return
	}
	model.instructions = ParseHPGL(string(source))
	model.instructionPointer = 0
	
	return
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "s":
			amount := 1
			if m.numBuffer != "" {
				var err error
				amount, err = strconv.Atoi(m.numBuffer)
				if err != nil {
					panic(err)
				}
				m.numBuffer = ""
			}
			m = m.step(amount)
		default:
			if strings.Contains("0123456789", msg.String()) {
				m.numBuffer += msg.String()
			} else {
				m.numBuffer = ""
			}
		}
	}
	return m, nil
}

func (m model) step(n int) model {
	for range n {
		instruction := m.instructions[m.instructionPointer]
		m.hpglContext.RunInstruction(instruction)
		m.instructionPointer++
	}
	return m
}
