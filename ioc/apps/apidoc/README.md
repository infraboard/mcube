# apidoc

mcube 内置 API 文档：一份契约，双壳展示。

## 默认产品形态

| 角色 | UI | 路径 |
|------|-----|------|
| **阅读（默认）** | **Scalar** | `/ui.html` → 加载 `openapi.json` |
| **试调** | Swagger UI | `/swagger.html` → 同契约试调 |

不默认 Redoc（已弃用；`UI_ENGINE=redoc` 会映射为 Scalar）。

## 端点

| Path | 说明 |
|------|------|
| `/openapi.json` | OpenAPI 3.0.3（Scalar / Swagger UI 主契约） |
| `/swagger.json` | Swagger 2.0（兼容旧工具） |
| `/ui.html` | Scalar |
| `/swagger.html` | Swagger UI（Try it） |

`?open_to_api_key=true`：仅开放给 API Key 的 operation。

## 配置

| 字段 | 默认 | 说明 |
|------|------|------|
| `UI_ENGINE` | `both` | `scalar` / `swagger` / `both` |
| `CDN_SCALAR_JS` | jsDelivr `@scalar/api-reference` | 内网可改本地 |
| `CDN_SWAGGER_UI_BASE` | unpkg swagger-ui-dist | 试调壳 |
| `ENABLED` | true | false → 全部 404 |
| `CACHE_TTL_SECOND` | 0 | 0=进程缓存；-1=每次重建 |

## 为何默认 Scalar 而不是 Swagger UI？

- Swagger UI 太常见，当主文档缺少辨识度
- Scalar 阅读体验更好，且原生偏 OAS3
- 试调仍保留 Swagger UI（Authorize / Try it out 成熟）

## 深链（Deep Link）

与 `ScalarDeepLink(ui, tag, operationId)` 一致：

```text
/ui.html#tag/{OpenAPI标签}/{operationId}
```

例如：`#tag/蓝图编排/QueryBlueprint`

业务侧（如 API Key 开放接口列表）应生成该格式，勿再用 Redoc 的 `#operation/...`。

## x-* 扩展

`x-perm`、`x-open-to-api-key`、`x-required-auth`、`x-required-namespace`、`x-resource`、`x-action`

