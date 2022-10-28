package util

import (
	"crypto/md5"
	"encoding/hex"
	"os"
)

func FileExists(path string) bool {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return true
}

// 方式一通过Write传参
func MD5(str string) string {
	b := []byte(str)
	h := md5.New()
	h.Write(b)
	return hex.EncodeToString(h.Sum(nil))
}

// 方式二通过Sum传参
func MD5_2(str string) string {
	b := []byte(str)
	h := md5.New()
	return hex.EncodeToString(h.Sum(b))
}

func MD5_SALT(str string, salt string) string {
	b := []byte(str)
	s := []byte(salt)
	h := md5.New()
	h.Write(s) // 先写盐值
	h.Write(b)
	return hex.EncodeToString(h.Sum(nil))
}

// 错误的加盐并多次加密
func MD5_SALT_MULT(str string, salt string, times int) string {
	b := []byte(str)
	s := []byte(salt)
	h := md5.New()
	var res []byte
	for i := 0; i < times; i++ {
		h.Write(s)
		h.Write(b)
		res = h.Sum(nil)
		b = res
	}
	return hex.EncodeToString(res)
}
