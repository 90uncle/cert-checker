package main

import (
	"cert-checker/internal/checker"
	"cert-checker/internal/config"
	"cert-checker/internal/notifier"
	"flag"
	"log"
)

func main() {
	configFile := flag.String("config", "config/config.yaml", "Path to config file")
	flag.Parse()

	// 1. 加载配置
	cfg, err := config.Load(*configFile)
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}
	// 2. 初始化通知器
	ntf := notifier.NewNotifier(&cfg.Notifiers.Email, &cfg.Notifiers.DingTalk, &cfg.Notifiers.WeCom)

	// 3. 检查每个域名的证书
	for _, domain := range cfg.Domains {
		info, err := checker.CheckCert(domain, cfg.Alert.Threshold)
		if err != nil {
			continue
		}

		log.Printf("域名: %-30s 过期时间: %s 剩余天数: %d\n",
			domain,
			info.ExpiryDate.Format("2006-01-02"),
			info.ExpiresIn)

		if info.IsExpired || info.IsWarning {
			msg := notifier.AlertMessage{
				Domain:     domain,
				ExpiryDate: info.ExpiryDate.Format("2006-01-02 15:04:05"),
				DaysLeft:   info.ExpiresIn,
			}

			_ = ntf.Send(msg)
		}
	}
}
