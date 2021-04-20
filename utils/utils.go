package utils

import (
	"crypto/md5"
	"encoding/hex"
	"net/http"
)

func MD5(s string) string {
	hash := md5.Sum([]byte(s))
	return hex.EncodeToString(hash[:])
}

func GetIP(r *http.Request) string {
	forwarded := r.Header.Get("X-FORWARDED-FOR")
	if forwarded != "" {
		return forwarded
	}
	return r.RemoteAddr
}

func GetAddrIndexFromNodes(target string, addrs []string) int {
	res := 0
	for i, addr := range addrs {
		if addr == target {
			res = i
			break
		}
	}
	return res
}

func IsAddrInNodes(target string, addrs []string) bool {
	for _, addr := range addrs {
		if addr == target {
			return true
		}
	}
	return false
}
