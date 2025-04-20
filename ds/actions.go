package ds

import "fmt"

// Copies the provided evaluated expression to the clipboard.
func ActionClipboard(data string) string {
	return fmt.Sprintf("@clipboard('%s')", data)
}

// Sets the value of all matching signals to the expression provided in the second argument. The first argument accepts one or more space-separated signal paths. You can use * to match a single path segment and ** to match multiple path segments.
func ActionSetAll(paths string, value string) string {
	return fmt.Sprintf("@setAll('%s', %s)", paths, value)
}

// Toggles the value of all matching signals. The first argument accepts one or more space-separated signal paths. You can use * to match a single path segment and ** to match multiple path segments.
func ActionToggle(paths string) string {
	return fmt.Sprintf("@toggleAll('%s')", paths)
}
