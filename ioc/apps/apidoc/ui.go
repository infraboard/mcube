package apidoc

import (
	"fmt"
	"html"
	"strings"
)

// RenderScalarHTML 默认阅读页。specURL 建议指向 openapi.json（OAS3）。
// switchHref: Swagger UI 试调链接；空则不展示切换条。
func RenderScalarHTML(switchHref, specURL, scalarJS string) string {
	specURL = html.EscapeString(specURL)
	scalarJS = html.EscapeString(scalarJS)
	switchBar := ""
	if strings.TrimSpace(switchHref) != "" {
		switchBar = fmt.Sprintf(
			`<div class="doc-switch">
  <span class="doc-switch-label">文档视图</span>
  <span class="doc-switch-current">Scalar</span>
  <a class="doc-switch-link" href="%s">Swagger UI · 试调</a>
</div>`,
			html.EscapeString(switchHref),
		)
	}
	return fmt.Sprintf(`<!DOCTYPE html>
<html>
  <head>
    <title>API Doc</title>
    <meta charset="utf-8"/>
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <style>
      html, body { margin: 0; padding: 0; height: 100%%; }
      body { display: flex; flex-direction: column; min-height: 100vh; }
      .doc-switch {
        flex: 0 0 auto;
        display: flex;
        align-items: center;
        gap: 10px;
        padding: 8px 16px;
        border-bottom: 1px solid #e8e8e8;
        background: #fafafa;
        font: 13px/1.4 -apple-system, BlinkMacSystemFont, "Segoe UI", sans-serif;
        z-index: 1;
      }
      .doc-switch-label { color: #8c8c8c; }
      .doc-switch-current { color: #1d1d1d; font-weight: 600; }
      .doc-switch-link {
        margin-left: auto;
        color: #722ED1;
        text-decoration: none;
        padding: 4px 10px;
        border: 1px solid #d3adf7;
        background: #fff;
      }
      .doc-switch-link:hover { background: #f9f0ff; }
      #app { flex: 1 1 auto; min-height: 0; }
    </style>
  </head>
  <body>
    %s
    <div id="app"></div>
    <script src="%s"></script>
    <script>
      Scalar.createApiReference('#app', {
        url: '%s',
        hideModels: false,
        defaultHttpClient: { targetKey: 'shell', clientKey: 'curl' },
        // 与业务侧 DocURL 约定一致：#tag/{tag}/{operationId}
        generateTagSlug: (tag) => (tag && tag.name) ? String(tag.name).replace(/ /g, '-') : 'default',
        generateOperationSlug: (operation) => {
          if (operation && operation.operationId) return String(operation.operationId);
          const method = (operation && operation.method) ? String(operation.method).toLowerCase() : 'get';
          const path = (operation && operation.path) ? String(operation.path) : '';
          return method + path;
        }
      });
    </script>
  </body>
</html>`, switchBar, scalarJS, specURL)
}

// RenderSwaggerUIHTML 试调页。switchHref: Scalar 阅读页链接。
func RenderSwaggerUIHTML(switchHref, specURL, swaggerUIBase string) string {
	specURL = html.EscapeString(specURL)
	base := html.EscapeString(strings.TrimRight(swaggerUIBase, "/"))
	switchBar := ""
	if strings.TrimSpace(switchHref) != "" {
		switchBar = fmt.Sprintf(
			`<div class="doc-switch">
  <span class="doc-switch-label">文档视图</span>
  <span class="doc-switch-current">Swagger UI · 试调</span>
  <a class="doc-switch-link" href="%s">Scalar · 阅读</a>
</div>`,
			html.EscapeString(switchHref),
		)
	}
	return fmt.Sprintf(`<!DOCTYPE html>
<html>
  <head>
    <title>API Doc · Try it</title>
    <meta charset="utf-8"/>
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <link rel="stylesheet" href="%s/swagger-ui.css"/>
    <style>
      html, body { margin: 0; padding: 0; }
      .doc-switch {
        display: flex;
        align-items: center;
        gap: 10px;
        padding: 8px 16px;
        border-bottom: 1px solid #e8e8e8;
        background: #fafafa;
        font: 13px/1.4 -apple-system, BlinkMacSystemFont, "Segoe UI", sans-serif;
      }
      .doc-switch-label { color: #8c8c8c; }
      .doc-switch-current { color: #1d1d1d; font-weight: 600; }
      .doc-switch-link {
        margin-left: auto;
        color: #722ED1;
        text-decoration: none;
        padding: 4px 10px;
        border: 1px solid #d3adf7;
        background: #fff;
      }
      .doc-switch-link:hover { background: #f9f0ff; }
    </style>
  </head>
  <body>
    %s
    <div id="swagger-ui"></div>
    <script src="%s/swagger-ui-bundle.js"></script>
    <script>
      window.ui = SwaggerUIBundle({
        url: '%s',
        dom_id: '#swagger-ui',
        deepLinking: true,
        presets: [SwaggerUIBundle.presets.apis],
        layout: 'BaseLayout'
      });
    </script>
  </body>
</html>`, base, switchBar, base, specURL)
}

// RenderRedocHTML 已废弃，内部转 Scalar。
func RenderRedocHTML(switchHref, specURL, _ string) string {
	return RenderScalarHTML(switchHref, specURL, DefaultApiDoc().CDNScalarJS)
}
