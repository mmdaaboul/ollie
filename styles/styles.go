package styles

import lipgloss "github.com/charmbracelet/lipgloss"

var HighlightStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#668599")).Padding(1, 1, 1, 1).Align(lipgloss.Center)
var ErrorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF0000")).Padding(1, 1, 1, 1).Align(lipgloss.Center)
