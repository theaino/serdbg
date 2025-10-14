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

func (m model) View() string {
	instructions, startIdx := m.getInstructionSlice(m.height - 4 - 1)
	text := make([]string, len(instructions))
	for idx, instruction := range instructions {
		realIdx := idx + startIdx
		instructionStyle := sentStyle
		if realIdx == m.instructionPointer {
			instructionStyle = currentStyle
		} else if realIdx > m.instructionPointer {
			instructionStyle = futureStyle
		}
		line := fmt.Sprintf("%s%s", lineNumStyle.Render(strconv.Itoa(realIdx)), instructionStyle.Render(instruction))
		text = append(text, line)
	}
	instructionView := strings.Join(text, "\n")

	dbgTable := table.New().
		Border(lipgloss.NormalBorder()).
		Headers("x", "y", "pen").Rows([]string{
			strconv.Itoa(m.hpglContext.X),
			strconv.Itoa(m.hpglContext.Y),
			map[bool]string{true: "down", false: "up"}[m.hpglContext.PenDown],
		})
	dbgView := lipgloss.NewStyle().Width(m.width).Height(4).AlignHorizontal(lipgloss.Center).Render(dbgTable.Render())

	footerText := fmt.Sprintf("s - Step\t[0-9]+s - Step n times\tq - Quit\t\t%s", m.numBuffer)

	footerView := lipgloss.NewStyle().Width(m.width).
		Background(lipgloss.Color("7")).
		Foreground(lipgloss.Color("0")).
		Bold(true).
		Render(footerText)

	return lipgloss.JoinVertical(lipgloss.Left, instructionView, dbgView, footerView)
}

func (m model) getInstructionSlice(height int) ([]string, int) {
	padding := max(float64(height - 1) / 2, 0)
	startIdx := max(0, m.instructionPointer - int(math.Ceil(padding)))
	endIdx := min(len(m.instructions), startIdx + 2 * int(padding) + 1)

	return m.instructions[startIdx:endIdx], startIdx
}
