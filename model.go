package main

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"go.bug.st/serial"
)

type ModelState any

type StateNormal struct {}
type StateSerialInput struct {
	Option SerialOption
}

type Model struct {
	State ModelState
	Width, Height int
	NumBuffer string
	ErrorBuffer struct {
		Value string
		Time time.Time
	}

	HpglContext *HPGLContext
	Instructions []string
	InstructionPointer int

	TextInput textinput.Model
	InputValue string

	PortName string
	SerialPort serial.Port
	SerialMode *serial.Mode
}

func NewModel(file *os.File) (model *Model, err error) {
	source, err := io.ReadAll(file)
	if err != nil {
		return
	}

	model = &Model{
		State: StateNormal{},

		HpglContext: NewHPGLContext(),
		Instructions: ParseHPGL(string(source)),
		InstructionPointer: 0,

		TextInput: textinput.New(),

		SerialMode: &serial.Mode{
			BaudRate: 9600,
			DataBits: 7,
			Parity: serial.OddParity,
			StopBits: serial.OneStopBit,
		},
	}
	
	return
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmds := make([]tea.Cmd, 0)

	switch m.State.(type) {
	case StateSerialInput:
		var cmd tea.Cmd
		m.TextInput, cmd = m.TextInput.Update(msg)
		cmds = append(cmds, cmd)
	}

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
		switch s := m.State.(type) {
		case StateNormal:
			var cmd tea.Cmd
			cmd = m.HandleKey(msg.String())
			cmds = append(cmds, cmd)
		case StateSerialInput:
			switch msg.String() {
			case "enter":
				m.SetSerialOption(s.Option, m.TextInput.Value())
				fallthrough
			case "esc":
				m.State = StateNormal{}
			}
		}
	}
	return m, tea.Batch(cmds...)
}

func (m *Model) HandleKey(key string) tea.Cmd {
	var cmd tea.Cmd
	switch key {
	case "s":
		amount := 1
		if m.NumBuffer != "" {
			var err error
			amount, err = strconv.Atoi(m.NumBuffer)
			if err != nil {
				panic(err)
			}
			m.NumBuffer = ""
		}
		m.step(amount)
	case "o":
		m.OpenPort()
	default:
		option, ok := SerialOptionKeyMap()[key]
		if ok {
			m.State = StateSerialInput{
				Option: option,
			}
			cmd = m.TextInput.Focus()
			m.TextInput.Reset()
			m.TextInput.Placeholder = SerialOptionDefinitions[option].Placeholder
			break
		}
		if strings.Contains("0123456789", key) {
			m.NumBuffer += key
		} else {
			m.NumBuffer = ""
		}
	}
	return cmd
}

func (m *Model) step(n int) {
	for range n {
		instruction := m.Instructions[m.InstructionPointer]
		m.HpglContext.RunInstruction(instruction)
		m.InstructionPointer++
	}
}

func (m *Model) Error(value any) {
	m.ErrorBuffer.Value = fmt.Sprintf("%v", value)
	m.ErrorBuffer.Time = time.Now()
}
