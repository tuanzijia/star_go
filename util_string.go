package starGo

import (
	"crypto/md5"
	"errors"
	"fmt"
	"github.com/satori/go.uuid"
	"runtime"
	"strings"
	"unicode"
)

// 获取随机长度的字符串
func GetRandomString(l int) string {
	bytes := []byte(baseString)
	var result []byte
	for i := 0; i < l; i++ {
		result = append(result, bytes[RandIntN(len(bytes))])
	}
	return string(result)
}

// 检查一个字符串是否是空字符串
func IsEmpty(content string) bool {
	if len(content) <= 0 {
		return true
	}

	return strings.IndexFunc(content, func(item rune) bool {
		return unicode.IsSpace(item) == false
	}) < 0
}

// 根据不同平台获取换行符
func GetNewLineString() string {
	switch os := runtime.GOOS; os {
	case "windows":
		return "\r\n"
	default:
		return "\n"
	}
}

// 获取新的UUID字符串
func GetNewUUID() string {
	return fmt.Sprintf("%v", uuid.Must(uuid.NewV4(), nil))
}

// 判断UUID是否为空
func IsUUIDEmpty(uuid string) bool {
	if uuid == "" || uuid == "00000000-0000-0000-0000-000000000000" {
		return true
	}

	return false
}

// 比较UUID是否相等
func IsUUIDEqual(uuid1, uuid2 string) bool {
	u1, err1 := uuid.FromString(uuid1)
	u2, err2 := uuid.FromString(uuid2)
	if err1 != nil || err2 != nil {
		return false
	}
	return uuid.Equal(u1, u2)
}

// 对字符数组进行MD5加密，并且可以选择返回大、小写
func Md5Bytes(b []byte, ifUpper bool) string {
	if len(b) == 0 {
		panic(errors.New("input []byte can't be empty"))
	}

	md5Instance := md5.New()
	md5Instance.Write(b)
	result := md5Instance.Sum([]byte(""))
	if ifUpper {
		return fmt.Sprintf("%X", result)
	} else {
		return fmt.Sprintf("%x", result)
	}
}

// 对字符串进行MD5加密，并且可以选择返回大、小写
func Md5String(s string, ifUpper bool) string {
	if len(s) == 0 {
		panic(errors.New("input string can't be empty"))
	}

	return Md5Bytes([]byte(s), ifUpper)
}
