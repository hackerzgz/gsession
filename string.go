package randstr

import (
	"bytes"
	"errors"
	"math"
)

type charType uint

const (
	CharDigit = 1 << iota
	CharLowerCase
	CharUpperCase

	digitSeq     = "0123456789"
	lowerCaseSeq = "abcdefghijklmnopqrstuvwxyz"
	upperCaseSeq = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

func NewString(typ charType, l int) (string, error) {
	var candidates bytes.Buffer

	for {
		if typ == 0 {
			break
		}

		switch {
		case typ&CharDigit > 0:
			candidates.WriteString(digitSeq)
			typ ^= CharDigit
		case typ&CharLowerCase > 0:
			candidates.WriteString(lowerCaseSeq)
			typ ^= CharLowerCase
		case typ&CharUpperCase > 0:
			candidates.WriteString(upperCaseSeq)
			typ ^= CharUpperCase
		}
	}

	if candidates.Len() == 0 {
		return "", errors.New("invalid char type")
	}

	c := candidates.Bytes()

	charBits := uint(bits(len(c) - 1))
	charMask := uint(1)<<charBits - 1
	charNum := 64 / charBits
	b := make([]byte, l)

	for i, cache, remain := l-1, _rand.Uint64(), charNum; i >= 0; {
		if remain == 0 {
			cache, remain = _rand.Uint64(), charNum
		}

		if idx := uint(cache) & charMask; idx < uint(len(c)) {
			b[i] = c[idx]
			i--
		}
		cache >>= charBits
		remain--
	}

	return string(b), nil
}

func bits(i int) int {
	return math.Ilogb(float64(i)) + 1
}
