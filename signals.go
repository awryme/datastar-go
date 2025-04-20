package datastar

import (
	"fmt"
	"strings"
)

// Signals is a map[string]any that is transformed (unflattened) and set as Value field.
//
// Top level signal names are split by the "." separator.
// If you send a signal {"user.fullname.name": "john"} it will be turned into {"user": {"fullname": {"name": "john"}}}.
// You can use nested signals if you wish.
type Signals map[string]any

// ErrTransform is error in case we coudn't transform signals into a nested object. Basically indicates if two fields are conflicting.
//
// Example: signal {"user.fullname.name": "user name"} is used along with signal {"user.fullname" = "users full name"}
var ErrTransform = fmt.Errorf("cannot transform signals")

const signalSeparator = "."

func transformTopLevelSignals(signals Signals) error {
	transformed := make([]string, 0, 1)
	for name, value := range signals {
		if strings.Contains(name, signalSeparator) {
			nameParts := strings.Split(name, signalSeparator)
			err := upsertSignal(signals, "", nameParts, value)
			if err != nil {
				return err
			}
			transformed = append(transformed, name)
		}
	}
	for _, name := range transformed {
		delete(signals, name)
	}
	return nil
}

func addName(prefix, name string) string {
	if prefix == "" {
		return name
	}
	return prefix + signalSeparator + name
}

func upsertSignal(signals Signals, nameprefix string, nameParts []string, value any) error {
	if len(nameParts) == 0 {
		return nil
	}
	name := nameParts[0]
	nameprefix = addName(nameprefix, name)
	if len(nameParts) == 1 {
		obj, ok := signals[name]
		if !ok {
			signals[name] = value
		}
		if _, ok := obj.(Signals); ok {
			return fmt.Errorf("%w: cannot place value to %s: signal object already exists", ErrTransform, nameprefix)
		}
		signals[name] = value
		return nil
	}

	obj, ok := signals[name]
	if !ok {
		// object doesn't exist, create new one
		partSignals := make(Signals)
		err := upsertSignal(partSignals, nameprefix, nameParts[1:], value)
		if err != nil {
			return err
		}
		signals[name] = partSignals
		return nil
	}
	if partSignals, ok := obj.(Signals); ok {
		// Signals object exists
		return upsertSignal(partSignals, nameprefix, nameParts[1:], value)
	}
	return fmt.Errorf("%w: cannot place signals object to %s: signal value already exists", ErrTransform, nameprefix)
}
