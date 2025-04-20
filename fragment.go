package datastar

import (
	"context"
	"io"
)

// MergeMode sets the mode to merge fragments with.
// Refer to individual MergeMode constants for details.
type MergeMode string

const (
	ModeMorph            MergeMode = "morph"            //Merges the fragment using Idiomorph. This is the default merge strategy.
	ModeInner            MergeMode = "inner"            //Replaces the target’s innerHTML with the fragment.
	ModeOuter            MergeMode = "outer"            //Replaces the target’s outerHTML with the fragment.
	ModePrepend          MergeMode = "prepend"          //Prepends the fragment to the target’s children.
	ModeAppend           MergeMode = "append"           //Appends the fragment to the target’s children.
	ModeBefore           MergeMode = "before"           //Inserts the fragment before the target as a sibling.
	ModeAfter            MergeMode = "after"            //Inserts the fragment after the target as a sibling.
	ModeUpsertAttributes MergeMode = "upsertAttributes" //Merges attributes from the fragment into the target – useful for updating a signals.
)

// Fragment is a standard fragment renderer.
// Examples are: gomponents, gostar.
type Fragment interface {
	Render(w io.Writer) error
}

// CtxFragment is a ctx aware fragment renderer.
// Examples are: templ.
type CtxFragment interface {
	Render(ctx context.Context, w io.Writer) error
}
