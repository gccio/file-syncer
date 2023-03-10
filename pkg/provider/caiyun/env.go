package caiyun

import "os"

var (
	Telephone      = os.Getenv("CAIYUN_TELEPHONE")
	Token          = os.Getenv("CAIYUN_ORCHES-C-TOKEN")
	Account        = os.Getenv("CAIYUN_ORCHES-C-ACCOUNT")
	AccountEncrypt = os.Getenv("CAIYUN_ORCHES-I-ACCOUNT-ENCRYPT")
)
var cookie = ""
var basePath = "https://yun.139.com"

func init() {
	cookieMap := map[string]string{
		"ORCHES-C-TOKEN":           Token,
		"ORCHES-C-ACCOUNT":         Account,
		"ORCHES-I-ACCOUNT-ENCRYPT": AccountEncrypt,
	}
	for k, v := range cookieMap {
		cookie += k + "=" + v + ";"
	}
}
