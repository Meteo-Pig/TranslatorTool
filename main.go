package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

var AUTHORIZATION string

// 字典请求体
type DictPayload struct {
	Source    string `json:"source"`
	TransType string `json:"trans_type"`
}

// 字典响应体
type DictResponse struct {
	Rc   int `json:"rc"`
	Wiki struct {
	} `json:"wiki"`
	Dictionary struct {
		Prons struct {
			EnUs string `json:"en-us"`
			En   string `json:"en"`
		} `json:"prons"`
		Explanations []string      `json:"explanations"`
		Synonym      []string      `json:"synonym"`
		Antonym      []interface{} `json:"antonym"`
		WqxExample   [][]string    `json:"wqx_example"`
		Entry        string        `json:"entry"`
		Type         string        `json:"type"`
		Related      []interface{} `json:"related"`
		Source       string        `json:"source"`
	} `json:"dictionary"`
}

// 翻译请求体
type TranslatorPayload struct {
	Source    string `json:"source"`
	TransType string `json:"trans_type"`
	RequestID string `json:"request_id"`
	Media     string `json:"media"`
	OsType    string `json:"os_type"`
	Dict      bool   `json:"dict"`
	Cached    bool   `json:"cached"`
	Replaced  bool   `json:"replaced"`
	Detect    bool   `json:"detect"`
	BrowserID string `json:"browser_id"`
}

// 翻译响应体
type TranslatorResponse struct {
	Isdict     int     `json:"isdict"`
	Confidence float64 `json:"confidence"`
	Target     string  `json:"target"`
	Rc         int     `json:"rc"`
	Jwt        string  `json:"jwt"`
}

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

/**
 * 网络请求
 * @param {string}  url - 请求地址
 * @param {any} payload - 请求体
 * @returns {[]byte} 响应体
 */
func request(url string, payload any) []byte {
	client := &http.Client{}
	buf, err := json.Marshal(payload)
	if err != nil {
		log.Fatal(err)
	}
	data := bytes.NewReader(buf)
	req, err := http.NewRequest("POST", url, data)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("authority", "api.interpreter.caiyunai.com")
	req.Header.Set("accept", "application/json, text/plain, */*")
	req.Header.Set("content-type", "application/json;charset=UTF-8")
	req.Header.Set("device-id", "0c6b3518ed767c4bba8921fb6b2f24bb")
	req.Header.Set("origin", "https://fanyi.caiyunapp.com")
	req.Header.Set("referer", "https://fanyi.caiyunapp.com/")
	req.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/113.0.0.0 Safari/537.36 Edg/113.0.1774.42")
	req.Header.Set("x-authorization", "token:qgemv4jr1y38jyq6vhvi")
	// 判断是否是翻译请求
	switch payload.(type) {
	case TranslatorPayload:
		// 判断是否是翻译请求
		if payload.(TranslatorPayload).Detect {
			req.Header.Set("t-authorization", AUTHORIZATION)
		}
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	bodyText, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	// fmt.Printf("%s\n", bodyText)
	return bodyText
}

/**
 * 字典查询
 * @param {string} word - 需要查询的单词
 * @returns {DictResponse} 查询结果
 */
func dictQuery(word string) DictResponse {
	payload := DictPayload{
		Source:    word,
		TransType: "en2zh",
	}
	res := request("https://api.interpreter.caiyunai.com/v1/dict", payload)
	var dictRes DictResponse
	json.Unmarshal(res, &dictRes)
	return dictRes
}

/**
 * 翻译
 * @param {string} word - 需要翻译的内容
 * @returns {TranslatorResponse} 翻译结果
 */
func translator(word string) TranslatorResponse {
	payload := TranslatorPayload{
		BrowserID: "0c6b3518ed767c4bba8921fb6b2f24bb",
		Cached:    true,
		Detect:    true,
		Dict:      true,
		Media:     "text",
		OsType:    "web",
		Replaced:  true,
		RequestID: "web_fanyi",
		Source:    word,
		TransType: "auto2zh",
	}
	res := request("https://api.interpreter.caiyunai.com/v1/translator", payload)
	var translatorRes TranslatorResponse
	json.Unmarshal(res, &translatorRes)
	return translatorRes
}

/**
 * 生成JWT
 * @param {string} input - 输入的字符串
 * @returns {string} JWT
 */
func generateJwt() string {
	payload := TranslatorPayload{BrowserID: "0c6b3518ed767c4bba8921fb6b2f24bb"}
	res := request("https://api.interpreter.caiyunai.com/v1/user/jwt/generate", payload)
	var translatorRes TranslatorResponse
	json.Unmarshal(res, &translatorRes)
	// fmt.Printf("%s\n", res)
	return translatorRes.Jwt
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("请输入需要翻译的内容")
		return
	}
	word := os.Args[1]
	AUTHORIZATION = generateJwt()
	translatorRes := translator(word)
	code := substituteAlphabet(translatorRes.Target)
	text := decodeBase64(code)
	fmt.Println(text)

	if translatorRes.Isdict == 1 {
		dictResponse := dictQuery(word)
		for _, v := range dictResponse.Dictionary.Explanations {
			fmt.Printf("%s\n", v)
		}
	}
}
