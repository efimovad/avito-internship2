package app

import (
	"errors"
	"math"
	"net"
)

func GetNetIP(cidr int) (net.IPMask, error){
	if cidr > 24 || cidr < 0 {
		return nil, errors.New("wrong cidr value")
	}

	var maskByte[4] byte
	for i := 0; i < 4 || cidr > 0; i++ {
		if cidr >= 8 {
			maskByte[i] = 255
		} else {
			maskByte[i] = byte(256 - int(math.Pow(2, float64(8 - cidr))))
		}
		cidr = cidr - 8
	}
	mask := net.IPv4Mask(maskByte[0], maskByte[1], maskByte[2], maskByte[3])
	return mask, nil
}
