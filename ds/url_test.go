package ds

import (
	"testing"

	"github.com/matryer/is"
)

func TestMapToJS(t *testing.T) {
	is := is.New(t)

	vals := map[string]string{
		"x":      "y",
		"asdQwe": "32",
		"inner": mapToJs(map[string]string{
			"k1": "'v1'",
			"k2": "'v2'",
		}),
		"qq": "123",
	}

	obj := mapToJs(vals)
	expected := "{asdQwe: 32, inner: {k1: 'v1', k2: 'v2'}, qq: 123, x: y}"
	is.Equal(obj, expected)
}
