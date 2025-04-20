package ds

import "strings"

// FilterEnterKey filters keyboard inputs with that have Enter key pressed.
// This allows html text inputs to run actions on Enter.
const FilterEnterKey = "evt.key == 'Enter'"

// Join joins scripts/actions with "&&"
func Join(parts ...string) string {
	return strings.Join(parts, " && ")
}

// OnEnterKey applies FilterEnterKey to a set of actions.
func OnEnterKey(actions ...string) string {
	return Join(FilterEnterKey, Join(actions...))
}
