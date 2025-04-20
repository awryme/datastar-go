package datastar

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/awryme/datastar-go/bufpool"
	"github.com/awryme/sse-go/sseserver"
)

// MergeFragments is a shortcut to create a `datastar-merge-fragments` event with a set of fragments.
func MergeFragments(fragments ...Fragment) EventMergeFragments {
	return EventMergeFragments{
		Fragments: fragments,
	}
}

// MergeCtxFragments is a shortcut to create a `datastar-merge-fragments` event with a set of ctx aware fragments.
func MergeCtxFragments(fragments ...CtxFragment) EventMergeFragments {
	return EventMergeFragments{
		CtxFragments: fragments,
	}
}

// EventMergeFragments is the implementation for `datastar-merge-fragments` event.
// You can set both Fragment and CtxFragments, as many as you need in a single request.
// Event options can be set with respective fields.
type EventMergeFragments struct {
	Fragments    []Fragment
	CtxFragments []CtxFragment

	// Selects the target element of the merge process using a CSS selector.
	Selector string

	// Sets the mode to merge fragments with.
	// Refer to individual MergeMode constants for details.
	MergeMode MergeMode

	// Determines whether to use view transitions when merging into the DOM.
	UseViewTransition bool
}

func (event EventMergeFragments) Name() string {
	return "datastar-merge-fragments"
}

func (event EventMergeFragments) WriteEvent(writer *sseserver.EventWriter, req *http.Request) error {
	if event.Selector != "" {
		writer.Format("selector %s", event.Selector)
	}
	if event.MergeMode != "" {
		writer.Format("mergeMode %s", event.MergeMode)
	}
	if event.UseViewTransition {
		writer.Write("useViewTransition true")
	}

	ctx := req.Context()

	buf := bufpool.GetBuffer()
	defer bufpool.PutBuffer(buf)

	renderFragment := func(render func(context.Context, io.Writer) error) error {
		buf.Reset()

		if err := render(ctx, buf); err != nil {
			return fmt.Errorf("render fragment: %w", err)
		}

		for _, line := range splitLines(buf.String()) {
			writer.Format("fragments %s", line)
		}
		return nil
	}

	for _, fragment := range event.CtxFragments {
		err := renderFragment(fragment.Render)
		if err != nil {
			return err
		}
	}

	for _, fragment := range event.Fragments {
		err := renderFragment(func(ctx context.Context, w io.Writer) error {
			return fragment.Render(w)
		})
		if err != nil {
			return err
		}
	}

	return nil
}
