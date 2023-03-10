package caiyun

import (
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"math/rand"
	"net/url"
	"sort"
	"strings"
	"time"
)

func formatCaiYunDate(date string) time.Time {
	t, _ := time.Parse("20060102150405", date)
	return t
}

func normalizePath(p string) string {
	return strings.TrimRight(p, "/")
}

func getRandomStr(n int) string {
	letters := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789")
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func getSign(ts, key, data string) string {
	data = strings.TrimSpace(data)
	data = encodeURIComponent(data)
	c := strings.Split(data, "")
	sort.Strings(c)
	data = strings.Join(c, "")
	s1 := genMD5(base64.StdEncoding.EncodeToString([]byte(data)))
	s2 := genMD5(ts + ":" + key)
	return strings.ToUpper(genMD5(s1 + s2))
}

func genMD5(data string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(data)))
}

func encodeURIComponent(str string) string {
	r := url.QueryEscape(str)
	r = strings.Replace(r, "+", "%20", -1)
	r = strings.Replace(r, "%21", "!", -1)
	r = strings.Replace(r, "%27", "'", -1)
	r = strings.Replace(r, "%28", "(", -1)
	r = strings.Replace(r, "%29", ")", -1)
	r = strings.Replace(r, "%2A", "*", -1)
	return r
}
