package auth

import (
	"crypto/md5"
	"fmt"
	"io"
	"os"
)

var (
	SALT0, SALT1 string
)

func InitSalt() {
	SALT0 = os.Getenv("BBB_SALT0")
	SALT1 = os.Getenv("BBB_SALT1")
}

func PassEncrypt(pwd string) string {
	h := md5.New()
	io.WriteString(h, pwd)
	pwdmd5 := fmt.Sprintf("%x", h.Sum(nil))
	io.WriteString(h,SALT0)
	io.WriteString(h,pwdmd5)
	io.WriteString(h,SALT1)
	return fmt.Sprintf("%x",h.Sum(nil))
}
