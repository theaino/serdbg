package main

import (
	"fmt"
	"slices"
	"strconv"
	"strings"
)

var ETXCmds = []string{"LB", "BL"}

type HPGLInstruction struct {
	Command string
	Argument string
	Source string
}

func NewHPGLInstruction(command, argument, source string) HPGLInstruction {
	return HPGLInstruction{
		strings.ToUpper(command),
		argument,
		source,
	}
}

func (i HPGLInstruction) String() string {
	if i.Argument == "" {
		return fmt.Sprintf("%s;", i.Command)
	}
	if slices.Contains(ETXCmds, i.Command) {
		return fmt.Sprintf("%s %#v;", i.Command, i.Argument)
	}
	return fmt.Sprintf("%s %s;", i.Command, i.Argument)
}

type HPGLParsingState struct {
	Terminator string
}

func NewHPGLParsingState() *HPGLParsingState {
	return &HPGLParsingState{
		Terminator: "\x03",
	}
}

func (s *HPGLParsingState) ParseInstructions(source string) []HPGLInstruction {
	instructions := make([]HPGLInstruction, 0)
	inETX := false
	currentCommand := ""
	currentArgument := ""
	currentSource := ""
	for _, char := range []byte(source) {
		if !strings.Contains("\n", string(char)) {
			currentSource += string(char)
		}
		if inETX {
			if string(char) == s.Terminator {
				inETX = false
			} else {
				currentArgument += string(char)
			}
			continue
		}
		if strings.Contains(" \n", string(char)) {
			continue
		}
		if len(currentCommand) < 2 {
			currentCommand += string(char)
			inETX = slices.Contains(ETXCmds, currentCommand)
			continue
		}
		if char == ';' {
			instruction := NewHPGLInstruction(currentCommand, currentArgument, currentSource)
			s.HandleInstruction(instruction)
			instructions = append(instructions, instruction)
			currentCommand = ""
			currentArgument = ""
			currentSource = ""
			continue
		}
		currentArgument += string(char)
	}
	if len(currentCommand) > 0 {
		instruction := NewHPGLInstruction(currentCommand, currentArgument, currentSource)
		s.HandleInstruction(instruction)
		instructions = append(instructions, instruction)
	}
	return instructions
}

func (s *HPGLParsingState) HandleInstruction(instruction HPGLInstruction) {
	if instruction.Command == "DT" {
		s.Terminator = instruction.Argument
	}
}

type HPGLInterpretingState struct {
	X, Y int
	PenDown bool
}

func NewHPGLInterpretingState() *HPGLInterpretingState {
	return &HPGLInterpretingState{
		X: 0,
		Y: 0,
		PenDown: false,
	}
}

func (s *HPGLInterpretingState) InterpretInstruction(instruction HPGLInstruction) (err error) {
	switch instruction.Command {
	case "PA":
		if instruction.Argument == "" {
			s.X = 0
			s.Y = 0
		} else {
			coordinate := strings.Split(instruction.Argument, ",")
			s.X, err = strconv.Atoi(coordinate[0])
			if err != nil {
				return
			}
			s.Y, err = strconv.Atoi(coordinate[1])
			if err != nil {
				return
			}
		}
	case "PR":
		if instruction.Argument == "" {
		} else {
			coordinate := strings.Split(instruction.Argument, ",")
			var x, y int
			x, err = strconv.Atoi(coordinate[0])
			if err != nil {
				return
			}
			y, err = strconv.Atoi(coordinate[1])
			if err != nil {
				return
			}
			s.X += x
			s.Y += y
		}
	case "PU":
		s.PenDown = false
	case "PD":
		s.PenDown = true
	}
	return
}

type HPGLState struct {
	*HPGLParsingState
	*HPGLInterpretingState
}

func NewHPGLState() *HPGLState {
	return &HPGLState{
		NewHPGLParsingState(),
		NewHPGLInterpretingState(),
	}
}

func (s *HPGLState) RunInstruction(instruction HPGLInstruction) error {
	s.HandleInstruction(instruction)
	return s.InterpretInstruction(instruction)
}

