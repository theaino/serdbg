package main

import (
	"fmt"
	"maps"
	"slices"
	"sort"
	"strconv"
	"strings"

	"go.bug.st/serial"
)

type SerialOption uint
const (
	OptionPort SerialOption = iota
	OptionBaudRate
	OptionDataBits
	OptionParity
	OptionStopBits
)

var ParityMap = map[string]serial.Parity{
	"no": serial.NoParity,
	"odd": serial.OddParity,
	"even": serial.EvenParity,
	"mark": serial.MarkParity,
	"space": serial.SpaceParity,
}

var StopBitsMap = map[string]serial.StopBits{
	"1": serial.OneStopBit,
	"1.5": serial.OnePointFiveStopBits,
	"2": serial.TwoStopBits,
}

func joinKeysReadable[T any](stringMap map[string]T) string {
	keys := slices.Collect(maps.Keys(stringMap))
	sort.Strings(keys)
	return strings.Join(append(
		keys[:len(stringMap)-1],
		fmt.Sprintf("or %s", keys[len(stringMap)-1]),
	), ", ")
}

var SerialOptionDefinitions = map[SerialOption]struct{
	Key string
	Name string
	Placeholder string
}{
	OptionPort: {
		Key: "p",
		Name: "Port",
		Placeholder: "Path...",
	},
	OptionBaudRate: {
		Key: "b",
		Name: "Baudrate",
		Placeholder: "Baudrate...",
	},
	OptionDataBits: {
		Key: "d",
		Name: "Data bits",
		Placeholder: "Data bits...",
	},
	OptionParity: {
		Key: "r",
		Name: "Parity",
		Placeholder: fmt.Sprintf("Parity (%v)...", joinKeysReadable(ParityMap)),
	},
	OptionStopBits: {
		Key: "t",
		Name: "Stop bits",
		Placeholder: fmt.Sprintf("Stop bits (%v)...", joinKeysReadable(StopBitsMap)),
	},
}

var SerialOptions = []SerialOption{OptionPort, OptionBaudRate, OptionDataBits, OptionParity, OptionStopBits}

func SerialOptionKeyMap() map[string]SerialOption {
	keyMap := make(map[string]SerialOption)
	for key, value := range SerialOptionDefinitions {
		keyMap[value.Key] = key
	}
	return keyMap
}

func (m *Model) OpenPort() {
	var err error
	m.SerialPort, err = serial.Open(m.PortName, m.SerialMode)
	if err != nil {
		m.Error(err)
	}
}

func (m *Model) SetSerialOption(option SerialOption, value string) {
	switch option {
	case OptionPort:
		m.PortName = value
	case OptionBaudRate:
		baudrate, err := strconv.Atoi(value)
		if err == nil {
			m.SerialMode.BaudRate = baudrate
		} else {
			m.Error(err)
			return
		}
	case OptionDataBits:
		dataBits, err := strconv.Atoi(value)
		if err == nil {
			m.SerialMode.DataBits = dataBits
		} else {
			m.Error(err)
			return
		}
	case OptionParity:
		parity, ok := map[string]serial.Parity{
			"no": serial.NoParity,
			"odd": serial.OddParity,
			"even": serial.EvenParity,
			"mark": serial.MarkParity,
			"space": serial.SpaceParity,
		}[strings.ToLower(value)]
		if ok {
			m.SerialMode.Parity = parity
		} else {
			m.Error(fmt.Sprintf("Invalid parity: %s", value))
			return
		}
	case OptionStopBits:
		stopBits, ok := map[string]serial.StopBits{
			"1": serial.OneStopBit,
			"1.5": serial.OnePointFiveStopBits,
			"2": serial.TwoStopBits,
		}[strings.ToLower(value)]
		if ok {
			m.SerialMode.StopBits = stopBits
		} else {
			m.Error(fmt.Sprintf("Invalid stop bits: %s", value))
			return
		}
	}
	if m.PortName != "" {
		m.OpenPort()
	}
}

func (m *Model) GetSerialOption(option SerialOption) string {
	switch option {
	case OptionPort:
		return m.PortName
	case OptionBaudRate:
		return strconv.Itoa(m.SerialMode.BaudRate)
	case OptionDataBits:
		return strconv.Itoa(m.SerialMode.DataBits)
	case OptionParity:
		for key, value := range ParityMap {
			if value == m.SerialMode.Parity {
				return key
			}
		}
	case OptionStopBits:
		for key, value := range StopBitsMap {
			if value == m.SerialMode.StopBits {
				return key
			}
		}
	}
	return ""
}
