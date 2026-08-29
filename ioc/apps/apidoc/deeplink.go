package apidoc

import (
	"fmt"
	"net/url"
	"strings"
)

// ScalarOperationHash 生成 Scalar 默认路由下的 operation 锚点。
// 需与 RenderScalarHTML 里 generateTagSlug / generateOperationSlug 一致：
//
//	#tag/{tag}/{operationId}
func ScalarOperationHash(tag, operationID string) string {
	tag = strings.TrimSpace(tag)
	operationID = strings.TrimSpace(operationID)
	if operationID == "" {
		return ""
	}
	if tag == "" {
		// 无 tag 时 Scalar 仍可能落到 introduction；仅带 operation 不可靠
		return ""
	}
	// 保留中文 tag；仅编码会破坏片段解析的字符
	return "#tag/" + scalarHashSegment(tag) + "/" + scalarHashSegment(operationID)
}

// ScalarDeepLink UI 相对路径 + Scalar operation 锚点。
func ScalarDeepLink(uiPath, tag, operationID string) string {
	uiPath = strings.TrimSpace(uiPath)
	hash := ScalarOperationHash(tag, operationID)
	if uiPath == "" {
		return hash
	}
	if hash == "" {
		return uiPath
	}
	return uiPath + hash
}

func scalarHashSegment(s string) string {
	// 空格 → 与 Scalar slug 常见行为接近；其余原样（中文 tag 已验证可用）
	s = strings.ReplaceAll(s, " ", "-")
	// 避免注入打断路径
	s = strings.ReplaceAll(s, "/", "-")
	s = strings.ReplaceAll(s, "#", "")
	s = strings.ReplaceAll(s, "?", "")
	// 不要 PathEscape 中文，浏览器地址栏与 Scalar 都按字面匹配
	if strings.ContainsAny(s, "%") {
		if dec, err := url.PathUnescape(s); err == nil {
			s = dec
		}
	}
	return s
}

// FormatScalarOperationSlug 供文档说明/测试，与前端 generateOperationSlug 对齐。
func FormatScalarOperationSlug(operationID, method, path string) string {
	if id := strings.TrimSpace(operationID); id != "" {
		return id
	}
	return fmt.Sprintf("%s%s", strings.ToLower(strings.TrimSpace(method)), strings.TrimSpace(path))
}
