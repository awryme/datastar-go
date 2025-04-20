package datastar

import (
	"net/http"
	"strings"

	"github.com/awryme/sse-go/sseserver"
)

// Event is a datastar event that will be sent to the client.
// Current implementations: EventMergeFragments, EventMergeSignals, EventRemoveFragments, EventRemoveSignals, EventExecuteScript.
// Refer to individual events for details.
type Event interface {
	Name() string
	WriteEvent(writer *sseserver.EventWriter, req *http.Request) error
}

func splitLines(data string) []string {
	return strings.Split(data, "\n")
}
