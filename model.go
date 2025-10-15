package main

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
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
	SerialWrittenBuffer string
	ErrorBuffer struct {
		Value string
		Time time.Time
	}

	HpglState *HPGLState
	Instructions []HPGLInstruction
	InstructionPointer atomic.Int64
	InstructionPointerTarget int64

	TextInput textinput.Model
	InputValue string

	KillStepListener chan bool
	SerialMutex sync.Mutex
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

		HpglState: NewHPGLState(),
		Instructions: NewHPGLParsingState().ParseInstructions(string(source)),

		TextInput: textinput.New(),

		KillStepListener:make(chan bool),
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
			m.KillStepListener <- true
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
		amount := int64(1)
		if m.NumBuffer != "" {
			var err error
			amount, err = strconv.ParseInt(m.NumBuffer, 10, 64)
			if err != nil {
				panic(err)
			}
			m.NumBuffer = ""
		}
		m.InstructionPointerTarget += amount
	case "e":
		m.InstructionPointerTarget = int64(len(m.Instructions)) - 1
	case "x":
		m.InstructionPointerTarget = m.InstructionPointer.Load()
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

func (m *Model) StepListener(program *tea.Program) {
	for {
		select {
		case <-m.KillStepListener:
			return
		default:
		}
		if m.InstructionPointer.Load() >= m.InstructionPointerTarget {
			continue
		}
		m.SerialMutex.Lock()
		instruction := m.Instructions[m.InstructionPointer.Load()]
		m.HpglState.RunInstruction(instruction)
		if m.SerialPort != nil {
			n, err := m.SerialPort.Write([]byte(instruction.Source))
			if err != nil {
				m.Error(err)
			} else if err := m.SerialPort.Drain(); err != nil {
				m.Error(err)
			} else {
				m.SerialWrittenBuffer = strconv.Itoa(n)
			}
		}
		m.InstructionPointer.Add(1)
		m.SerialMutex.Unlock()
		program.Send(struct{}{})
	}
}

func (m *Model) Error(value any) {
	m.ErrorBuffer.Value = fmt.Sprintf("%v", value)
	m.ErrorBuffer.Time = time.Now()
}
