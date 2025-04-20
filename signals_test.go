package datastar

import (
	"errors"
	"testing"

	"github.com/matryer/is"
)

func TestTransform(t *testing.T) {
	is := is.New(t)

	runTest := func(name string, signals Signals, expected Signals, expectedErr error) {
		t.Run(name, func(t *testing.T) {
			is := is.New(t)

			err := transformTopLevelSignals(signals)
			is.True(errors.Is(err, expectedErr)) // error should be expected
			if expectedErr != nil {
				return
			}

			is.Equal(signals, expected) // transformed signals should be equal to expected
		})
	}

	runTest("ok: no transform",
		Signals{
			"asd": 3,
			"qwe": 5,
		},
		Signals{
			"asd": 3,
			"qwe": 5,
		},
		nil,
	)

	runTest("ok: 1 transform",
		Signals{
			"asd.qwe": 3,
			"qwe":     5,
		},
		Signals{
			"asd": Signals{
				"qwe": 3,
			},
			"qwe": 5,
		},
		nil,
	)
	runTest("ok: 2 nested merge transforms",
		Signals{
			"asd.qwe.data1": 3,
			"asd.qwe.data2": 5,
		},
		Signals{
			"asd": Signals{
				"qwe": Signals{
					"data1": 3,
					"data2": 5,
				},
			},
		},
		nil,
	)

	runTest("fail: transform into existing flat value",
		Signals{
			"asd.qwe.data1": 3,
			"asd.qwe":       5,
		},
		nil,
		ErrTransform,
	)

	runTest("fail: transform into existing nested value",
		Signals{
			"asd.qwe.data1": 3,
			"asd": Signals{
				"qwe": 5,
			},
		},
		nil,
		ErrTransform,
	)

	runTest("fail: transform into existing nested object",
		Signals{
			"asd.qwe": 5,
			"asd": Signals{
				"qwe": Signals{
					"data1": 3,
				},
			},
		},
		nil,
		ErrTransform,
	)
}
