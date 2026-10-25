package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSetStasherRequest_Validate(t *testing.T) {

	testcases := []struct {
		testName string
		name     string
		email    string
		checkFn  func(t *testing.T, err error)
	}{
		{
			testName: "missing name",
			name:     "",
			email:    "john@example.com",
			checkFn: func(t *testing.T, err error) {
				require.Error(t, err)
			},
		},
		{
			testName: "missing email",
			name:     "John Doe",
			checkFn: func(t *testing.T, err error) {
				require.Error(t, err)
			},
		},
		{
			testName: "reject invalid email",
			email:    "invalid-email",
			name:     "Joe Bloggs",
			checkFn: func(t *testing.T, err error) {
				require.Error(t, err)
			},
		},
		{
			testName: "valid data",
			name:     "John Doe",
			email:    "john@example.com",
			checkFn: func(t *testing.T, err error) {
				require.NoError(t, err)
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.testName, func(t *testing.T) {

			req := &SetStasherRequest{
				Email: tc.email,
				Name:  tc.name,
			}
			err := req.Validate()
			tc.checkFn(t, err)
		})
	}
}
