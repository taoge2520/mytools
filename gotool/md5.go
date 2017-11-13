package gotool

import (
	"crypto/md5"
	"encoding/hex"
)

func MD5(s string) string {
	m := md5.New()
	m.Write([]byte(s))
	return hex.EncodeToString(m.Sum(nil))
}

func MD5WithSalt(s string, salts ...string) string {
	m := md5.New()
	m.Write([]byte(s))
	for _, v := range salts {
		m.Write([]byte(v))
	}
	return hex.EncodeToString(m.Sum(nil))
}
