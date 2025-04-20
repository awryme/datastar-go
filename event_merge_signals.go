package datastar

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/awryme/sse-go/sseserver"
)

// MergeSignals is a shortcut to create a `datastar-merge-signals` event with a Signals value.
// Signals is a map[string]any that is transformed (unflattened) and set as Value field.
//
// Top level signal names are split by the "." separator.
// If you send a signal {"user.fullname.name": "john"} it will be turned into {"user": {"fullname": {"name": "john"}}}.
// You can use nested signals if you wish.
func MergeSignals(signals Signals) EventMergeSignals {
	return EventMergeSignals{
		Signals: signals,
	}
}

// MergeSignalsObj is a shortcut to create a `datastar-merge-signals` event with a struct value.
func MergeSignalsObj(signals any) EventMergeSignals {
	return EventMergeSignals{
		Value: signals,
	}
}

// EventMergeSignals is the implementation for `datastar-merge-signals` event.
//
// Only one of Signals or Value may be provided. If both are set, only Value will be used.
//
// Use respective MergeSignals and MergeSignalsObj functions to simplify creation of the event.
type EventMergeSignals struct {
	// Signals values, refer to MergeSignals for docs.
	Signals Signals

	// Value is marshalled to a json object and sent to frontend.
	Value any

	// OnlyIfMissing determines whether to update the signals with new values only if the key does not exist.
	OnlyIfMissing bool
}

func (event EventMergeSignals) Name() string {
	return "datastar-merge-signals"
}

func (event EventMergeSignals) WriteEvent(writer *sseserver.EventWriter, req *http.Request) error {
	// transform signals
	if len(event.Signals) > 0 {
		err := transformTopLevelSignals(event.Signals)
		if err != nil {
			return fmt.Errorf("transform signals: %w", err)
		}
		event.Value = event.Signals
	}

	if event.OnlyIfMissing {
		writer.Write("onlyIfMissing true")
	}

	data, err := json.Marshal(event.Value)
	if err != nil {
		return fmt.Errorf("marshal signals: %w", err)
	}

	writer.Format("signals %s", string(data))
	return nil
}
