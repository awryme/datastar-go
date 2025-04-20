package datastar

import (
	"net/http"

	"github.com/awryme/sse-go/sseserver"
)

// RemoveFragments is a shortcut to create a `datastar-remove-fragments` event.
func RemoveFragments(selector string) EventRemoveFragments {
	return EventRemoveFragments{selector}
}

// EventRemoveFragments is the implementation for `datastar-remove-fragments` event.
type EventRemoveFragments struct {
	// Selector is used to match the elements to remove from the DOM.
	// Only one selector can be provided in the event.
	Selector string
}

func (event EventRemoveFragments) Name() string {
	return "datastar-remove-fragments"
}

func (event EventRemoveFragments) WriteEvent(writer *sseserver.EventWriter, req *http.Request) error {
	writer.Format("selector %s", event.Selector)
	return nil
}
