package notifier

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type DingTalkNotifier struct {
	Enabled   bool     `yaml:"enabled"`
	Webhook   string   `yaml:"webhook"`
	AtMobiles []string `yaml:"at_mobiles"`
}

func (n *DingTalkNotifier) Send(msg AlertMessage) error {
	if !n.Enabled {
		return nil
	}
	payload := map[string]interface{}{
		"msgtype": "text",
		"text": map[string]string{
			"content": msg.String(),
		},
		"at": map[string]interface{}{
			"atMobiles": n.AtMobiles,
			"isAtAll":   false,
		},
	}
	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal dingtalk payload error: %v", err)
	}

	resp, err := http.Post(n.Webhook, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return fmt.Errorf("send dingtalk message error: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("dingtalk returned non-200 status: %d, body: %s", resp.StatusCode, body)
	}
	return nil
}

func (n *DingTalkNotifier) Name() string {
	return "DingTalk"
}

func (n *DingTalkNotifier) IsEnabled() bool {
	return n.Enabled
}
