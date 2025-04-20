package datastar

import (
	"net/http"

	"github.com/awryme/sse-go/sseserver"
)

// RemoveSignals is a shortcut to create a `datastar-remove-signals` event.
func RemoveSignals(paths ...string) EventRemoveSignals {
	return EventRemoveSignals{paths}
}

// EventRemoveSignals is the implementation for `datastar-remove-signals` event.
type EventRemoveSignals struct {
	// Full paths to match the signals to remove.
	Paths []string
}

func (event EventRemoveSignals) Name() string {
	return "datastar-remove-signals"
}

func (event EventRemoveSignals) WriteEvent(writer *sseserver.EventWriter, req *http.Request) error {
	for _, path := range event.Paths {
		writer.Format("paths %s", path)
	}
	return nil
}
