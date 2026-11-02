package data

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_IsLocked(t *testing.T) {

	testcases := []struct {
		name  string
		stash *Stash
		exp   bool
	}{
		{
			name: "should return false when stash time in the past and details are set",
			stash: &Stash{
				Date:        "2026-11-01",
				ReleaseTime: "12:00",
				Details:     &Details{},
			},
			exp: false,
		},
		{
			name: "should return true when stash time in the past and details NOT set",
			stash: &Stash{
				Date:        "2026-11-01",
				ReleaseTime: "12:00",
				Details:     nil,
			},
			exp: true,
		},
		{
			name: "should return true when stash time in the future",
			stash: &Stash{
				Date:        "2026-11-02",
				ReleaseTime: "21:13",
			},
			exp: true,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {

			tm := time.Date(2026, 11, 2, 21, 0, 34, 0, time.UTC)
			res := tc.stash.IsLocked(tm)
			assert.Equal(t, tc.exp, res)
		})
	}
}
