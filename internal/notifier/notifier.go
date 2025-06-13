package notifier

import (
	"fmt"
	"log"
)

// INotifier 通知接口
type INotifier interface {
	IsEnabled() bool // 是否启用通知
	Send(message AlertMessage) error
	Name() string // 返回通知器名称
}

type AlertMessage struct {
	Domain     string
	ExpiryDate string
	DaysLeft   int
}

func (m *AlertMessage) String() string {
	if m.DaysLeft < 0 {
		return fmt.Sprintf("域名 %s 的证书已过期 (过期时间: %s)", m.Domain, m.ExpiryDate)
	}
	return fmt.Sprintf("域名 %s 的证书将在 %d 天后过期 (过期时间: %s)",
		m.Domain, m.DaysLeft, m.ExpiryDate)
}

type Notifier struct {
	notifiers []INotifier
}

func NewNotifier(notifiers ...INotifier) *Notifier {
	return &Notifier{
		notifiers,
	}
}

func (n *Notifier) Send(message AlertMessage) error {
	for _, notifier := range n.notifiers {
		if notifier.IsEnabled() {
			log.Printf("%s 开始发送通知: %s", notifier.Name(), message.String())
		} else {
			log.Printf("%s 未启用, 跳过通知", notifier.Name())
		}
		if err := notifier.Send(message); err != nil {
			log.Printf("%s 发送失败: %s", notifier.Name(), err)
			return err
		}
	}

	return nil
}

func (n *Notifier) Name() string {
	return "Notifier"
}

func (n *Notifier) IsEnabled() bool {
	return true
}
