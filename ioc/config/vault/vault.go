package vault

import (
	"context"
	"fmt"
	"os"
	"time"

	vault "github.com/hashicorp/vault-client-go"
	"github.com/hashicorp/vault-client-go/schema"
	"github.com/infraboard/mcube/v2/ioc"
	"github.com/infraboard/mcube/v2/ioc/config/log"
	"github.com/infraboard/mcube/v2/ioc/config/trace"
	"github.com/rs/zerolog"
)

func init() {
	ioc.Config().Registry(defaultConfig)
}

var defaultConfig = &Vault{
	Address:       "http://127.0.0.1:8200",
	Timeout:       30,
	AuthMethod:    "token",
	K8sTokenPath:  "/var/run/secrets/kubernetes.io/serviceaccount/token",
	TLSSkipVerify: false,
	AutoRenew:     true,
	Trace:         false,
	// 默认挂载点（与 Vault CLI 开发模式保持一致）
	KvMountPath:       "secret",
	TransitMountPath:  "transit",
	DatabaseMountPath: "database",
	PkiMountPath:      "pki",
}

type Vault struct {
	ioc.ObjectImpl

	// 服务配置
	Address   string `toml:"address" json:"address" yaml:"address" env:"ADDRESS"`
	Namespace string `toml:"namespace" json:"namespace" yaml:"namespace" env:"NAMESPACE"` // 企业版
	Timeout   int    `toml:"timeout" json:"timeout" yaml:"timeout" env:"TIMEOUT"`         // 秒

	// 认证配置
	AuthMethod string `toml:"auth_method" json:"auth_method" yaml:"auth_method" env:"AUTH_METHOD"` // token, approle, kubernetes
	Token      string `toml:"token" json:"token" yaml:"token" env:"TOKEN"`

	// AppRole 认证
	RoleID   string `toml:"role_id" json:"role_id" yaml:"role_id" env:"ROLE_ID"`
	SecretID string `toml:"secret_id" json:"secret_id" yaml:"secret_id" env:"SECRET_ID"`

	// Kubernetes 认证
	K8sRole      string `toml:"k8s_role" json:"k8s_role" yaml:"k8s_role" env:"K8S_ROLE"`
	K8sTokenPath string `toml:"k8s_token_path" json:"k8s_token_path" yaml:"k8s_token_path" env:"K8S_TOKEN_PATH"`

	// TLS 配置
	TLSSkipVerify bool   `toml:"tls_skip_verify" json:"tls_skip_verify" yaml:"tls_skip_verify" env:"TLS_SKIP_VERIFY"`
	CACert        string `toml:"ca_cert" json:"ca_cert" yaml:"ca_cert" env:"CA_CERT"` // CA 证书路径

	// 秘密引擎挂载点配置（与 Vault CLI 默认值保持一致）
	KvMountPath       string `toml:"kv_mount_path" json:"kv_mount_path" yaml:"kv_mount_path" env:"KV_MOUNT_PATH"`                         // KV v2 引擎挂载点，默认 "secret"
	TransitMountPath  string `toml:"transit_mount_path" json:"transit_mount_path" yaml:"transit_mount_path" env:"TRANSIT_MOUNT_PATH"`     // Transit 引擎挂载点，默认 "transit"
	DatabaseMountPath string `toml:"database_mount_path" json:"database_mount_path" yaml:"database_mount_path" env:"DATABASE_MOUNT_PATH"` // Database 引擎挂载点，默认 "database"
	PkiMountPath      string `toml:"pki_mount_path" json:"pki_mount_path" yaml:"pki_mount_path" env:"PKI_MOUNT_PATH"`                     // PKI 引擎挂载点，默认 "pki"

	// 功能开关
	AutoRenew bool `toml:"auto_renew" json:"auto_renew" yaml:"auto_renew" env:"AUTO_RENEW"` // 自动续期 token
	Trace     bool `toml:"trace" json:"trace" yaml:"trace" env:"TRACE"`

	// 内部状态
	client      *vault.Client
	log         *zerolog.Logger
	stopRenewal chan struct{}
}

// SetKvMountPath 设置 KV 引擎挂载点（用于测试）
func (v *Vault) SetKvMountPath(path string) {
	v.KvMountPath = path
}

// SetTransitMountPath 设置 Transit 引擎挂载点（用于测试）
func (v *Vault) SetTransitMountPath(path string) {
	v.TransitMountPath = path
}

// SetDatabaseMountPath 设置 Database 引擎挂载点（用于测试）
func (v *Vault) SetDatabaseMountPath(path string) {
	v.DatabaseMountPath = path
}

// SetPkiMountPath 设置 PKI 引擎挂载点（用于测试）
func (v *Vault) SetPkiMountPath(path string) {
	v.PkiMountPath = path
}

// ValidateMountPaths 验证所有挂载点配置的合法性（用于测试）
func (v *Vault) ValidateMountPaths() error {
	return v.validateMountPaths()
}

func (v *Vault) Name() string {
	return AppName
}

func (v *Vault) Priority() int {
	return 696
}

func (v *Vault) Init() error {
	v.log = log.Sub(v.Name())
	v.stopRenewal = make(chan struct{})

	// 0. 验证挂载点配置
	if err := v.validateMountPaths(); err != nil {
		return fmt.Errorf("invalid mount path configuration: %w", err)
	}

	// 1. 创建客户端选项
	opts := []vault.ClientOption{
		vault.WithAddress(v.Address),
		vault.WithRequestTimeout(time.Duration(v.Timeout) * time.Second),
	}

	if v.Namespace != "" {
		v.log.Info().Msgf("using namespace: %s", v.Namespace)
		// Namespace 通过环境变量设置
		os.Setenv("VAULT_NAMESPACE", v.Namespace)
	}

	// 2. TLS 配置
	if v.TLSSkipVerify {
		v.log.Warn().Msg("TLS verification disabled - not recommended for production")
		// vault-client-go 会读取 VAULT_SKIP_VERIFY 环境变量
		os.Setenv("VAULT_SKIP_VERIFY", "true")
	}

	// 3. 创建客户端
	client, err := vault.New(opts...)
	if err != nil {
		return fmt.Errorf("create vault client: %w", err)
	}

	// 4. 认证（根据 AuthMethod）
	if err := v.authenticate(client); err != nil {
		return fmt.Errorf("vault authentication: %w", err)
	}

	// 5. 启用 trace（如果需要）
	if trace.Get().Enable && v.Trace {
		v.log.Info().Msg("vault trace enabled")
		// TODO: 可以在这里集成 OpenTelemetry HTTP transport
	}

	v.client = client
	v.log.Info().Msgf("vault client initialized (address=%s, auth=%s)", v.Address, v.AuthMethod)

	return nil
}

func (v *Vault) Close(ctx context.Context) {
	if v.stopRenewal != nil {
		close(v.stopRenewal)
	}

	if v.client != nil {
		v.log.Info().Msg("vault client closed")
	}
}

func (v *Vault) authenticate(client *vault.Client) error {
	ctx := context.Background()

	switch v.AuthMethod {
	case "token":
		if v.Token == "" {
			return fmt.Errorf("token is required for token auth")
		}
		if err := client.SetToken(v.Token); err != nil {
			return fmt.Errorf("set token: %w", err)
		}
		v.log.Info().Msg("using token authentication")

	case "approle":
		if v.RoleID == "" || v.SecretID == "" {
			return fmt.Errorf("role_id and secret_id are required for approle auth")
		}

		resp, err := client.Auth.AppRoleLogin(
			ctx,
			schema.AppRoleLoginRequest{
				RoleId:   v.RoleID,
				SecretId: v.SecretID,
			},
			vault.WithMountPath("approle"),
		)
		if err != nil {
			return fmt.Errorf("approle login: %w", err)
		}

		token, ok := resp.Data["client_token"].(string)
		if !ok || token == "" {
			return fmt.Errorf("approle login failed: no token returned")
		}

		if err := client.SetToken(token); err != nil {
			return fmt.Errorf("set token from approle: %w", err)
		}

		v.log.Info().Msg("approle authentication successful")

		// 启动自动续期
		if v.AutoRenew {
			// 尝试从响应中获取 lease_duration，如果失败则使用默认值 3600 秒
			leaseDuration := int32(3600)
			if ld, ok := resp.Data["lease_duration"].(float64); ok {
				leaseDuration = int32(ld)
			}
			go v.startTokenRenewal(leaseDuration)
		}

	case "kubernetes":
		if v.K8sRole == "" {
			return fmt.Errorf("k8s_role is required for kubernetes auth")
		}

		// 读取 ServiceAccount token
		jwt, err := os.ReadFile(v.K8sTokenPath)
		if err != nil {
			return fmt.Errorf("read k8s token from %s: %w", v.K8sTokenPath, err)
		}

		resp, err := client.Auth.KubernetesLogin(
			ctx,
			schema.KubernetesLoginRequest{
				Role: v.K8sRole,
				Jwt:  string(jwt),
			},
			vault.WithMountPath("kubernetes"),
		)
		if err != nil {
			return fmt.Errorf("kubernetes login: %w", err)
		}

		token, ok := resp.Data["client_token"].(string)
		if !ok || token == "" {
			return fmt.Errorf("kubernetes login failed: no token returned")
		}

		if err := client.SetToken(token); err != nil {
			return fmt.Errorf("set token from kubernetes: %w", err)
		}

		v.log.Info().Msg("kubernetes authentication successful")

		// 启动自动续期
		if v.AutoRenew {
			// 尝试从响应中获取 lease_duration，如果失败则使用默认值 3600 秒
			leaseDuration := int32(3600)
			if ld, ok := resp.Data["lease_duration"].(float64); ok {
				leaseDuration = int32(ld)
			}
			go v.startTokenRenewal(leaseDuration)
		}

	default:
		return fmt.Errorf("unsupported auth method: %s (supported: token, approle, kubernetes)", v.AuthMethod)
	}

	return nil
}

// validateMountPaths 验证所有挂载点配置的合法性
func (v *Vault) validateMountPaths() error {
	paths := map[string]string{
		"kv_mount_path":       v.KvMountPath,
		"transit_mount_path":  v.TransitMountPath,
		"database_mount_path": v.DatabaseMountPath,
		"pki_mount_path":      v.PkiMountPath,
	}

	for name, path := range paths {
		if path == "" {
			return fmt.Errorf("%s cannot be empty", name)
		}

		// 检查是否包含非法字符（Vault 路径只允许字母、数字、连字符、下划线）
		for _, char := range path {
			if !((char >= 'a' && char <= 'z') ||
				(char >= 'A' && char <= 'Z') ||
				(char >= '0' && char <= '9') ||
				char == '-' || char == '_') {
				return fmt.Errorf("%s contains invalid character '%c', only alphanumeric, hyphen and underscore allowed", name, char)
			}
		}

		// 检查是否使用了 Vault 保留路径
		reservedPaths := []string{"sys", "identity", "cubbyhole"}
		for _, reserved := range reservedPaths {
			if path == reserved {
				return fmt.Errorf("%s cannot use reserved path '%s'", name, reserved)
			}
		}
	}

	return nil
}

// startTokenRenewal 暂时简化实现，因为 vault-client-go v0.4.3 的 API 返回值较复杂
// TODO: 完善自动续期逻辑
func (v *Vault) startTokenRenewal(leaseDuration int32) {
	if leaseDuration <= 0 {
		v.log.Warn().Msg("invalid lease duration, auto-renewal disabled")
		return
	}

	// 在 token 生命周期的 80% 时进行续期
	renewInterval := time.Duration(leaseDuration) * time.Second * 80 / 100

	v.log.Info().Msgf("token auto-renewal enabled (interval=%s)", renewInterval)

	ticker := time.NewTicker(renewInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			ctx := context.Background()
			_, err := v.client.Auth.TokenRenewSelf(
				ctx,
				schema.TokenRenewSelfRequest{},
			)
			if err != nil {
				v.log.Error().Err(err).Msg("failed to renew token")
				// 续期失败，可以选择重新认证或告警
			} else {
				v.log.Debug().Msg("token renewed successfully")
			}

		case <-v.stopRenewal:
			v.log.Info().Msg("token renewal stopped")
			return
		}
	}
}
