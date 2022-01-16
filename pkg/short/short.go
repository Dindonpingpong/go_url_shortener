package short

import (
	"crypto/md5"
	"encoding/hex"
)

func GenereteShortString(s string) (string, error) {
	h := md5.New()

	h.Write([]byte(s))

	return hex.EncodeToString(h.Sum(nil)), nil
}
