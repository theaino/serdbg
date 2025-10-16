package main

import (
	"flag"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	port := flag.String("port", "", "Specify the serial port")
	baudrate := flag.Int("baudrate", 9600, "Specify the baudrate")
	dataBits := flag.Int("databits", 7, "Specify the data bits")
	parity := flag.String("parity", "odd", "Specify the parity")
	stopBits := flag.String("stopbits", "1", "Specify the stop bits")
	flag.Parse()

	model, err := NewModel()
	if err != nil {
		panic(err)
	}

	path := flag.Arg(0)
	if path != "" {
		model.LoadFile(path)
	}

	model.SerialMode.BaudRate = *baudrate
	model.SerialMode.DataBits = *dataBits
	model.SetSerialOption(OptionParity, *parity)
	model.SetSerialOption(OptionStopBits, *stopBits)
	model.PortName = *port
	if model.PortName != "" {
		model.OpenPort()
	}

	program := tea.NewProgram(model)
	go model.StepListener(program)
	if _, err := program.Run(); err != nil {
		panic(err)
	}
}
