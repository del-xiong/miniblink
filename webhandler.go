package miniblink

import (
	"strings"
)

// 黑白名单检查 请求发出之前检查 若未通过检测则不会发送请求
func (view *WebView) isRequestAllowed(url string) bool {
	var allow bool = false
	if len(view.whitelist) > 0 {
		// 白名单模式下 不匹配任意一白名单则禁止访问
		for _, v := range view.whitelist {
			if strings.Contains(url, v) {
				allow = true
			}
		}
		if !allow {
			return allow
		}
	}
	allow = true
	// 匹配任意一黑名单则禁止访问
	for _, v := range view.blacklist {
		if strings.Contains(url, v) {
			allow = false
		}
	}
	return allow
}