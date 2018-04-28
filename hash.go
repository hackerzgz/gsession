package gsession

import (
	"crypto/md5"
)

// hash use sha256 to get unique hashing string
func hash(str string) []byte {
	h := md5.New()
	h.Write([]byte(str))
	return h.Sum(nil)
}
