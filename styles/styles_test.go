package styles_test

import (
	"ollie/styles"
	"reflect"
	"testing"

	lipgloss "github.com/charmbracelet/lipgloss"
)

func TestHighlightStyle(t *testing.T) {
	expected := lipgloss.NewStyle().Foreground(lipgloss.Color("#003456")).Padding(1, 1, 1, 1).Align(lipgloss.Center)

	if !reflect.DeepEqual(styles.HighlightStyle, expected) {
		t.Errorf("Expected HighlightStyle to be %v, but got %v", expected, styles.HighlightStyle)
	}
}
