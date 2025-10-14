package main

import (
	"strconv"
	"strings"
)

type HPGLContext struct {
	X, Y int
	PenDown bool
}

func ParseHPGL(source string) (instructions []string) {
	instructions = strings.Split(strings.ReplaceAll(source, "\n", ""), ";")
	instructions = instructions[:len(instructions)-1]
	for idx, instruction := range instructions {
		instructions[idx] = instruction + ";"
	}
	return
}

func NewHPGLContext() *HPGLContext {
	return &HPGLContext{
		X: 0,
		Y: 0,
		PenDown: false,
	}
}

func (c *HPGLContext) RunInstruction(instruction string) error {
	instruction = strings.TrimSuffix(instruction, ";")
	parts := strings.Split(instruction, " ")
	switch strings.ToUpper(parts[0]) {
	case "PA":
		if len(parts) < 2 {
			c.X = 0
			c.Y = 0
		} else {
			coordinate := strings.Split(parts[1], ",")
			var err error
			c.X, err = strconv.Atoi(coordinate[0])
			if err != nil {
				return err
			}
			c.Y, err = strconv.Atoi(coordinate[1])
			if err != nil {
				return err
			}
		}
	case "PR":
		coordinate := strings.Split(parts[1], ",")
		x, err := strconv.Atoi(coordinate[0])
		if err != nil {
			return err
		}
		y, err := strconv.Atoi(coordinate[1])
		if err != nil {
			return err
		}
		c.X += x
		c.Y += y
	case "PU":
		c.PenDown = false
	case "PD":
		c.PenDown = true
	}
	return nil
}
