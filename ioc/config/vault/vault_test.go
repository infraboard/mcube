package vault_test

import (
	"testing"

	"github.com/hashicorp/vault-client-go"
	"github.com/infraboard/mcube/v2/ioc"
	iocvault "github.com/infraboard/mcube/v2/ioc/config/vault"
	"github.com/infraboard/mcube/v2/tools/file"
)

func TestVaultClientKvV2Read(t *testing.T) {
	client := iocvault.Client()
	if client != nil {
		t.Log("vault client initialized")
		// 示例：读取秘密
		// 注意：Vault 开发模式下 KV v2 引擎挂载在 secret/ 路径
		// - 挂载点: secret
		// - 秘密路径: myapp/config
		// - 实际 API 路径: secret/data/myapp/config
		secret, err := client.Secrets.KvV2Read(
			t.Context(),
			"myapp/config",                // 秘密路径（不含挂载点）
			vault.WithMountPath("secret"), // 明确指定挂载点
		)
		if err != nil {
			t.Fatal(err)
		}
		t.Logf("Secret data: %+v", secret.Data.Data)
		t.Logf("Secret metadata: %+v", secret.Data.Metadata)

		// 访问具体字段
		if apiKey, ok := secret.Data.Data["api_key"].(string); ok {
			t.Logf("api_key: %s", apiKey)
		}
		if dbURL, ok := secret.Data.Data["database_url"].(string); ok {
			t.Logf("database_url: %s", dbURL)
		}
	}
}

// TestReadSecretHelper 测试辅助方法（自动使用配置的挂载点）
func TestReadSecretHelper(t *testing.T) {
	// 使用辅助方法，无需手动指定挂载点
	secret, err := iocvault.ReadSecret(t.Context(), "myapp/config")
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("使用辅助方法读取成功")
	t.Logf("Secret data: %+v", secret.Data.Data)

	// 访问具体字段
	if apiKey, ok := secret.Data.Data["api_key"].(string); ok {
		t.Logf("api_key: %s", apiKey)
	}
}

// TestWriteSecretHelper 测试写入辅助方法
func TestWriteSecretHelper(t *testing.T) {
	// 写入秘密
	data := map[string]interface{}{
		"test_key":   "test_value",
		"created_at": "2026-02-20",
	}

	_, err := iocvault.WriteSecret(t.Context(), "myapp/test", data)
	if err != nil {
		t.Fatal(err)
	}

	t.Log("写入秘密成功")

	// 读取验证
	secret, err := iocvault.ReadSecret(t.Context(), "myapp/test")
	if err != nil {
		t.Fatal(err)
	}

	if val, ok := secret.Data.Data["test_key"].(string); ok && val == "test_value" {
		t.Log("验证成功: test_key = test_value")
	} else {
		t.Error("验证失败")
	}
}

func TestDefaultConfig(t *testing.T) {
	file.MustToToml(
		iocvault.AppName,
		ioc.Config().Get(iocvault.AppName),
		"test/default.toml",
	)
}

// TestMountPathValidation 测试挂载点验证逻辑
func TestMountPathValidation(t *testing.T) {
	tests := []struct {
		name          string
		kvMount       string
		transitMount  string
		databaseMount string
		pkiMount      string
		wantErr       bool
		errContains   string
	}{
		{
			name:          "valid mount paths",
			kvMount:       "secret",
			transitMount:  "transit",
			databaseMount: "database",
			pkiMount:      "pki",
			wantErr:       false,
		},
		{
			name:          "valid custom mount paths",
			kvMount:       "prod-kv",
			transitMount:  "prod_transit",
			databaseMount: "db-engine",
			pkiMount:      "pki_v2",
			wantErr:       false,
		},
		{
			name:          "empty kv mount path",
			kvMount:       "",
			transitMount:  "transit",
			databaseMount: "database",
			pkiMount:      "pki",
			wantErr:       true,
			errContains:   "kv_mount_path cannot be empty",
		},
		{
			name:          "invalid character in mount path",
			kvMount:       "secret/data",
			transitMount:  "transit",
			databaseMount: "database",
			pkiMount:      "pki",
			wantErr:       true,
			errContains:   "contains invalid character",
		},
		{
			name:          "reserved path sys",
			kvMount:       "sys",
			transitMount:  "transit",
			databaseMount: "database",
			pkiMount:      "pki",
			wantErr:       true,
			errContains:   "cannot use reserved path 'sys'",
		},
		{
			name:          "reserved path identity",
			kvMount:       "secret",
			transitMount:  "identity",
			databaseMount: "database",
			pkiMount:      "pki",
			wantErr:       true,
			errContains:   "cannot use reserved path 'identity'",
		},
		{
			name:          "reserved path cubbyhole",
			kvMount:       "secret",
			transitMount:  "transit",
			databaseMount: "cubbyhole",
			pkiMount:      "pki",
			wantErr:       true,
			errContains:   "cannot use reserved path 'cubbyhole'",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &iocvault.Vault{}
			v.SetKvMountPath(tt.kvMount)
			v.SetTransitMountPath(tt.transitMount)
			v.SetDatabaseMountPath(tt.databaseMount)
			v.SetPkiMountPath(tt.pkiMount)

			err := v.ValidateMountPaths()

			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error containing '%s', got nil", tt.errContains)
					return
				}
				if tt.errContains != "" && !contains(err.Error(), tt.errContains) {
					t.Errorf("expected error containing '%s', got '%s'", tt.errContains, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			}
		})
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > 0 && len(substr) > 0 && findSubstring(s, substr)))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func init() {
	// 初始化 IOC 容器
	// 如果需要测试，请在 test/ 目录下创建 application.toml 配置文件
	req := ioc.NewLoadConfigRequest()
	req.ConfigFile.Enabled = true
	req.ConfigFile.Paths = []string{"test/application.toml"}
	req.ConfigFile.SkipIFNotExist = true

	err := ioc.ConfigIocObject(req)
	if err != nil {
		panic(err)
	}
}
