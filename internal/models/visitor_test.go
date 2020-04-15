package models

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestVisitors_GetVisitor(t *testing.T) {
	testCases := []struct {
		name	string
		ip		string
		err		error
	}{
		{
			name: 	"valid",
			ip:		"128.69.0.1",
			err: 	nil,
		},
	}

	visitors := NewVisitors(10, time.Second * 60, time.Second * 120)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			v := visitors.GetVisitor(tc.ip)
			assert.NotNil(t, v)
		})
	}
}
