package notifier

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type WeComNotifier struct {
	Enabled          bool     `yaml:"enabled"`
	Webhook          string   `yaml:"webhook"`
	MentionedMobiles []string `yaml:"mentioned_mobiles"`
}

func (n *WeComNotifier) Send(msg AlertMessage) error {
	if !n.Enabled {
		return nil
	}
	payload := map[string]interface{}{
		"msgtype": "text",
		"text": map[string]interface{}{
			"content":               msg.String(),
			"mentioned_mobile_list": n.MentionedMobiles,
		},
	}
	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal wecom payload error: %v", err)
	}
	resp, err := http.Post(n.Webhook, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return fmt.Errorf("send wecom message error: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("wecom returned non-200 status: %d", resp.StatusCode)
	}
	return nil
}
func (n *WeComNotifier) Name() string {
	return "WeCom"
}

func (n *WeComNotifier) IsEnabled() bool {
	return n.Enabled
}
