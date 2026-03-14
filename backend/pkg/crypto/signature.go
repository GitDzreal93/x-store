package crypto

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sort"
	"strings"
)

// GenerateSignature 生成防重放签名
// 将请求参数按 key 排序后拼接，加上 timestamp + nonce + secret 进行 HMAC-SHA256
func GenerateSignature(params map[string]string, timestamp, nonce, secret string) string {
	// 按 key 排序
	keys := make([]string, 0, len(params))
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// 拼接参数
	var parts []string
	for _, k := range keys {
		parts = append(parts, fmt.Sprintf("%s=%s", k, params[k]))
	}
	payload := strings.Join(parts, "&")
	payload = fmt.Sprintf("%s&timestamp=%s&nonce=%s", payload, timestamp, nonce)

	// HMAC-SHA256
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(payload))
	return hex.EncodeToString(h.Sum(nil))
}

// VerifySignature 验证签名是否匹配
func VerifySignature(params map[string]string, timestamp, nonce, secret, signature string) bool {
	expected := GenerateSignature(params, timestamp, nonce, secret)
	return hmac.Equal([]byte(expected), []byte(signature))
}
