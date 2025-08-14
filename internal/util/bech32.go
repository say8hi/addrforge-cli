package util

import (
	"errors"
	"strings"
)

const charset = "qpzry9x8gf2tvdw0s3jn54khce6mua7l"

var gen = [5]uint32{0x3b6a57b2, 0x26508e6d, 0x1ea119fa, 0x3d4233dd, 0x2a1462b3}

// bech32Polymod computes the bech32 polymod checksum
func bech32Polymod(values []byte) uint32 {
	chk := uint32(1)
	for _, value := range values {
		top := chk >> 25
		chk = (chk&0x1ffffff)<<5 ^ uint32(value)
		for i := 0; i < 5; i++ {
			if (top>>uint(i))&1 == 1 {
				chk ^= gen[i]
			}
		}
	}
	return chk
}

// bech32HrpExpand expands the human readable part for checksum computation
func bech32HrpExpand(hrp string) []byte {
	h := []byte(strings.ToLower(hrp))
	ret := make([]byte, len(h)*2+1)
	for i, c := range h {
		ret[i] = c >> 5
	}
	ret[len(h)] = 0
	for i, c := range h {
		ret[len(h)+1+i] = c & 31
	}
	return ret
}

// bech32CreateChecksum creates the checksum for bech32 encoding
func bech32CreateChecksum(hrp string, data []byte) []byte {
	values := append(bech32HrpExpand(hrp), data...)
	values = append(values, []byte{0, 0, 0, 0, 0, 0}...)
	polymod := bech32Polymod(values) ^ 1
	ret := make([]byte, 6)
	for i := 0; i < 6; i++ {
		ret[i] = byte((polymod >> uint(5*(5-i))) & 31)
	}
	return ret
}

// convertBits converts between bit groups
func convertBits(data []byte, fromBits, toBits uint, pad bool) ([]byte, error) {
	acc := uint(0)
	bits := uint(0)
	ret := []byte{}
	maxv := uint((1 << toBits) - 1)
	maxAcc := uint((1 << (fromBits + toBits - 1)) - 1)

	for _, value := range data {
		if uint(value) >= (1 << fromBits) {
			return nil, errors.New("invalid data for base conversion")
		}
		acc = ((acc << fromBits) | uint(value)) & maxAcc
		bits += fromBits
		for bits >= toBits {
			bits -= toBits
			ret = append(ret, byte((acc>>bits)&maxv))
		}
	}

	if pad {
		if bits > 0 {
			ret = append(ret, byte((acc<<(toBits-bits))&maxv))
		}
	} else if bits >= fromBits || ((acc<<(toBits-bits))&maxv) != 0 {
		return nil, errors.New("invalid padding bits")
	}

	return ret, nil
}

// EncodeBech32 encodes data to bech32 format
func EncodeBech32(hrp string, data []byte) (string, error) {
	// Convert 8-bit data to 5-bit
	conv, err := convertBits(data, 8, 5, true)
	if err != nil {
		return "", err
	}

	// Create checksum
	checksum := bech32CreateChecksum(hrp, conv)
	combined := append(conv, checksum...)

	// Build the final string
	ret := hrp + "1"
	for _, d := range combined {
		ret += string(charset[d])
	}

	return ret, nil
}

// DecodeBech32 decodes bech32 encoded string
func DecodeBech32(bech string) (string, []byte, error) {
	// Check for mixed case
	if strings.ToLower(bech) != bech && strings.ToUpper(bech) != bech {
		return "", nil, errors.New("mixed case")
	}

	bech = strings.ToLower(bech)

	// Find separator
	pos := strings.LastIndex(bech, "1")
	if pos < 1 || pos > 83 || pos+7 > len(bech) {
		return "", nil, errors.New("invalid separator position")
	}

	hrp := bech[:pos]
	data := bech[pos+1:]

	// Decode data part
	decoded := make([]byte, len(data))
	for i, c := range data {
		d := strings.Index(charset, string(c))
		if d == -1 {
			return "", nil, errors.New("invalid character")
		}
		decoded[i] = byte(d)
	}

	// Verify checksum
	if bech32Polymod(append(bech32HrpExpand(hrp), decoded...)) != 1 {
		return "", nil, errors.New("invalid checksum")
	}

	// Convert 5-bit to 8-bit
	conv, err := convertBits(decoded[:len(decoded)-6], 5, 8, false)
	if err != nil {
		return "", nil, err
	}

	return hrp, conv, nil
}
