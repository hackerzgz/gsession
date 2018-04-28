package randstr

import (
	"math/rand"
	"time"
)

var _rand *rand.Rand

func init() {
	_rand = rand.New(rand.NewSource(time.Now().UnixNano()))
}
