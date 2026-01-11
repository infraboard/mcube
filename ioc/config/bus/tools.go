package bus

import (
	"regexp"
	"strings"
)

// sanitizeQueueName 清理队列名称，使其符合 RabbitMQ 命名规范
// 只允许字母、数字、连字符、点号、下划线
func SanitizeQueueName(name string) string {
	// 替换非法字符为下划线
	reg := regexp.MustCompile(`[^a-zA-Z0-9._-]`)
	sanitized := reg.ReplaceAllString(name, "_")

	// 移除开头和结尾的点号
	sanitized = strings.Trim(sanitized, ".")

	// 确保不为空
	if sanitized == "" {
		sanitized = "default"
	}

	// 限制长度（RabbitMQ 队列名最多 255 字节）
	if len(sanitized) > 255 {
		sanitized = sanitized[:255]
		// 再次移除结尾的点号
		sanitized = strings.TrimRight(sanitized, ".")
	}

	return sanitized
}
