package cxspec

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLocate(t *testing.T) {
	t.Run("split_loc_string", func(t *testing.T) {
		type TestCase struct {
			name string // test name
			in string   // test input

			// expected outputs
			prefix   LocPrefix
			suffix   string
			hasErr   bool
			hasPanic bool
		}

		cases := []TestCase{
			{
				name: "0_parts",
				in: "",
				hasErr: true,
			},
			{
				name: "1_parts",
				in: "9f86d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a08",
				prefix: TrackerLoc,
				suffix: "9f86d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a08",
			},
			{
				name: "2_parts",
				in: fmt.Sprintf("%s:%s", TrackerLoc, "9f86d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a08"),
				prefix: TrackerLoc,
				suffix: "9f86d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a08",
			},
			{
				name: "2_parts",
				in: fmt.Sprintf("%s:%s", FileLoc, `this\:/is/a\:/test/file.json`),
				prefix: FileLoc,
				suffix: `this\:/is/a\:/test/file.json`,
			},
		}

		for _, c := range cases {
			c := c

			t.Run(c.name, func(t *testing.T) {
				prefix, suffix, err := splitLocString(c.in)

				if c.hasErr {
					assert.Error(t, err)
					return
				}

				assert.Equal(t, c.prefix, prefix)
				assert.Equal(t, c.suffix, suffix)
				assert.NoError(t, err)
			})
		}
	})
}