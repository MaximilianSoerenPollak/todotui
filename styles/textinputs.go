package styles

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

var (
	TiFocusedStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	TiBlurredStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	TiCursorStyle         = TiFocusedStyle.Copy()
	TiNoStyle             = lipgloss.NewStyle()
	TiHelpStyle           = TiBlurredStyle.Copy()
	TiCursorModeHelpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("244"))

	TiFocusedButton = TiFocusedStyle.Copy().Render("[ Submit ]")
	TiBlurredButton = fmt.Sprintf("[ %s ]", TiBlurredStyle.Render("Submit"))
)
