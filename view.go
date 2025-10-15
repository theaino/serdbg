package main

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
)

var lineNumStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("8")).Width(4)
var sentStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("4"))
var currentStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("2")).Bold(true)
var futureStyle = lipgloss.NewStyle().Italic(true)

func (m *Model) View() string { if m.Width == 0 {
		return ""
	}

	wrappedErrorText, errorViewHeight := WrapString(fmt.Sprintf(" %v", m.ErrorBuffer.Value), m.Width)
	errorView := lipgloss.NewStyle().Foreground(lipgloss.Color("1")).
		Bold(true).
		Height(errorViewHeight).
		Render(wrappedErrorText)

	dbgViewHeight := 5
	dbgTable := table.New().
		Border(lipgloss.NormalBorder()).
		Headers("X", "Y", "PEN", "ETX").Rows([]string{
			strconv.Itoa(m.HpglState.X),
			strconv.Itoa(m.HpglState.Y),
			map[bool]string{true: "down", false: "up"}[m.HpglState.PenDown],
			fmt.Sprintf("%#v", m.HpglState.Terminator),
		})
	dbgView := lipgloss.NewStyle().Width(m.Width).Height(dbgViewHeight).Render(dbgTable.Render())

	serialViewEntries := make([]string, len(SerialOptions))
	for idx, option := range SerialOptions {
		serialViewEntries[idx] = fmt.Sprintf("%s: %s", SerialOptionDefinitions[option].Name, m.GetSerialOption(option))
	}
	if m.SerialWrittenBuffer != "" {
		serialViewEntries = append(serialViewEntries,
			lipgloss.NewStyle().Italic(true).Render(fmt.Sprintf("  %s b written", m.SerialWrittenBuffer)),
		)
	}
	serialText, serialViewHeight := JoinStringWrapped(serialViewEntries, "  ", m.Width)
	serialView := lipgloss.NewStyle().Foreground(lipgloss.Color(map[bool]string{
		true: "3",
		false: "2",
	}[m.SerialPort == nil])).Height(serialViewHeight).Render(serialText)

	m.TextInput.PlaceholderStyle = m.TextInput.PlaceholderStyle.
		Foreground(lipgloss.Color("8")).
		Background(lipgloss.Color("7")).
		Bold(true)

	m.TextInput.Width = m.Width
	var footerValue string
	footerViewHeight := 1
	switch m.State.(type) {
	case StateNormal:
		footerEntries := []string{
			"s - Step",
			"[0-9]+s - Step n times",
			"e - Send whole buffer",
			"x - Stop sending",
		}
		for _, option := range SerialOptions {
			definition := SerialOptionDefinitions[option]
			footerEntries = append(footerEntries, fmt.Sprintf("%s - Set %s", definition.Key, definition.Name))
		}
		footerEntries = append(footerEntries, "o - Reopen port", "^C - Quit")
		if m.NumBuffer != "" {
			footerEntries = append(footerEntries, fmt.Sprintf("  %s",
				lipgloss.NewStyle().Italic(true).Render(m.NumBuffer),
			))
		}
		footerValue, footerViewHeight = JoinStringWrapped(footerEntries, "    ", m.Width)
	case StateSerialInput:
		footerValue = m.TextInput.View()
	}

	footerView := lipgloss.NewStyle().Width(m.Width).
		Background(lipgloss.Color("7")).
		Foreground(lipgloss.Color("0")).
		Bold(true).
		Height(footerViewHeight).
		Render(footerValue)

	instructionViewHeight := m.Height - errorViewHeight - dbgViewHeight - serialViewHeight - footerViewHeight
	instructions, startIdx := m.getInstructionSlice(instructionViewHeight)
	text := make([]string, len(instructions))
	for idx, instruction := range instructions {
		realIdx := int64(idx) + startIdx
		instructionStyle := sentStyle
		if realIdx == m.InstructionPointer.Load() {
			instructionStyle = currentStyle
		} else if realIdx > m.InstructionPointer.Load() {
			instructionStyle = futureStyle
		}
		lineNumStyle := lineNumStyle
		if realIdx == m.InstructionPointerTarget {
			lineNumStyle = lineNumStyle.Foreground(lipgloss.Color("1"))
		}
		line := fmt.Sprintf("%s%s", lineNumStyle.Render(strconv.FormatInt(realIdx, 10)), instructionStyle.Render(instruction))
		text[idx] = line
	}
	instructionView := strings.Join(text, "\n")

	return lipgloss.JoinVertical(lipgloss.Left, instructionView, errorView, dbgView, serialView, footerView)
}

func (m *Model) getInstructionSlice(height int) ([]string, int64) {
	padding := max(float64(height - 1) / 2, 0)
	startIdx := max(0, m.InstructionPointer.Load() - int64(math.Ceil(padding)))
	endIdx := startIdx + int64(height)

	stringInstructions := make([]string, endIdx - startIdx)
	for idx, instruction := range m.Instructions[startIdx:endIdx] {
		stringInstructions[idx] = instruction.String()
	}
	return stringInstructions, startIdx
}
