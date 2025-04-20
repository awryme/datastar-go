package datastar

import (
	"net/http"

	"github.com/awryme/sse-go/sseserver"
)

// ExecuteScript is a shortcut to create a `datastar-execute-script` event.
func ExecuteScript(scrpit ...string) EventExecuteScript {
	return EventExecuteScript{
		Script: scrpit,
	}
}

// EventExecuteScript is the implementation for `datastar-execute-script` event.
// Event options can be set with respective fields.
type EventExecuteScript struct {
	// Script to send to the client.
	// It can have multiple lines that will be split by newline at event creation.
	Script []string

	// AutoRemove determines whether to remove the script after execution.
	AutoRemove bool

	// Attributes to add to the script element.
	Attributes map[string]string
}

func (event EventExecuteScript) Name() string {
	return "datastar-execute-script"
}

func (event EventExecuteScript) WriteEvent(writer *sseserver.EventWriter, req *http.Request) error {
	if event.AutoRemove {
		writer.Write("autoRemove true")
	}

	if len(event.Attributes) > 0 {
		for name, value := range event.Attributes {
			writer.Format("attributes %s %s", name, value)
		}
	}

	for _, script := range event.Script {
		for _, line := range splitLines(script) {
			writer.Format("script %s", line)
		}
	}
	return nil
}
