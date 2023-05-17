package main

import (
	"encoding/base64"
	"fmt"
)

// 将字母表中的字母替换为另一个字母表中的字母
func substituteAlphabet(input string) string {
	alphabet := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	substitution := "NOPQRSTUVWXYZABCDEFGHIJKLMnopqrstuvwxyzabcdefghijklm"
	findIndex := func(n byte) int {
		for index, b := range alphabet {
			if byte(b) == n {
				return index
			}
		}
		return -1
	}
	substitute := func(n byte) byte {
		if findIndex(n) > -1 {
			return substitution[findIndex(n)]
		}
		return n
	}
	result := ""
	for _, c := range input {
		result += string(substitute(byte(c)))
	}
	return result
}

// 从字符串中删除双字节字符
func removeDoubleByte(input string) string {
	bytes := []byte(input)
	for i := 0; i < len(bytes); i++ {
		if bytes[i] == 194 && i+1 < len(bytes) && bytes[i+1] >= 128 && bytes[i+1] <= 191 {
			bytes[i] = bytes[i+1]
			bytes[i+1] = 0
		}
	}
	return string(bytes)
}

// 对输入的字符串进行base64解码
func decodeBase64(input string) string {
	decoded, err := base64.StdEncoding.DecodeString(input)
	if err != nil {
		panic(err)
	}
	return removeDoubleByte(string(decoded))
}

func main() {
	code := substituteAlphabet("5Y2t5nJ977lZ5YvJ55JZ77lO")
	text := decodeBase64(code)
	fmt.Println(text)
}
