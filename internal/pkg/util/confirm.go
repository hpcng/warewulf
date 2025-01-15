package util

import (
	"github.com/manifoldco/promptui"
)

// ConfirmationPrompt prompt is a blocking confirmation prompt.
// Returns true on y or yes user input.
func Confirm(label string) (yes bool) {

	prompt := promptui.Prompt{
		Label:     label,
		IsConfirm: true,
	}

	result, _ := prompt.Run()
	if result == "y" || result == "yes" {
		yes = true
	}
	return
}
