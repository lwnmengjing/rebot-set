package lark

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"time"
)

// SendLarkMessage 发送消息到 Lark
func SendLarkMessage(webhook, secret, color, title, foot string, lines []string) error {
	timestamp := time.Now().Unix()
	sign, err := GenLarkMessageSign(secret, timestamp)
	if err != nil {
		slog.Error("Generate sign failed", "err", err)
		return err
	}
	elements := make([]map[string]interface{}, 0)
	for i := range lines {
		elements = append(elements, map[string]interface{}{
			"tag": "div",
			"text": map[string]string{
				"tag":     "lark_md",
				"content": lines[i],
			},
		}, map[string]interface{}{
			"tag": "hr",
		})
	}
	elements = append(elements, map[string]interface{}{
		"tag": "note",
		"elements": []map[string]string{
			{
				"tag":     "plain_text",
				"content": foot,
			},
		},
	})
	// 构建告警卡片消息
	m := map[string]interface{}{
		"timestamp": strconv.FormatInt(timestamp, 10),
		"sign":      sign,
		"msg_type":  "interactive",
		"card": map[string]interface{}{
			"config": map[string]interface{}{
				"wide_screen_mode": true,
			},
			"header": map[string]interface{}{
				"title": map[string]string{
					"tag":     "plain_text",
					"content": title,
				},
				"template": color, // 设置卡片头部颜色
			},
			"elements": elements,
		},
	}

	// 将消息结构体转换为 JSON
	rb, err := json.Marshal(m)
	if err != nil {
		slog.Error("Marshal message failed", "err", err)
		return err
	}

	// 发送卡片消息到 Lark
	resp, err := http.Post(webhook, "application/json", bytes.NewBuffer(rb))
	if err != nil {
		slog.Error("Send message to Lark failed", "err", err)
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		slog.Error("Send message to Lark failed", "status", resp.Status, "body", string(body))
		return fmt.Errorf("http status code: %d", resp.StatusCode)
	}
	return nil
}

// GenLarkMessageSign 生成 Lark 消息签名
func GenLarkMessageSign(secret string, timestamp int64) (string, error) {
	//timestamp + key 做sha256, 再进行base64 encode
	stringToSign := fmt.Sprintf("%v", timestamp) + "\n" + secret
	var data []byte
	h := hmac.New(sha256.New, []byte(stringToSign))
	_, err := h.Write(data)
	if err != nil {
		return "", err
	}
	signature := base64.StdEncoding.EncodeToString(h.Sum(nil))
	return signature, nil
}
