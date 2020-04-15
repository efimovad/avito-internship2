package app

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_getNetIP(t *testing.T) {
	testCases := []struct {
		name	string
		cidr 	int
		mask	[]byte
		err		error
	}{
		{
			name: 	"valid: /23",
			cidr: 	23,
			mask:	[]byte{255, 255, 254, 0},
			err: 	nil,
		},
		{
			name: 	"valid: /16",
			cidr: 	16,
			mask:	[]byte{255, 255, 0, 0},
			err: 	nil,
		},
		{
			name: 	"valid: /8",
			cidr: 	8,
			mask:	[]byte{255, 0, 0, 0},
			err: 	nil,
		},
		{
			name:	"invalid: negative cidr",
			cidr:	-10,
			mask:	nil,
			err:	errors.New("wrong cidr value"),
		},
		{
			name:	"invalid: too big value",
			cidr:	100,
			mask:	nil,
			err:	errors.New("wrong cidr value"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mask, err := GetNetIP(tc.cidr)
			assert.Equal(t, []byte(mask), tc.mask)
			assert.Equal(t, err, tc.err)
		})
	}
}
