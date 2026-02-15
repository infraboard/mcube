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

	// 功能开关
	AutoRenew bool `toml:"auto_renew" json:"auto_renew" yaml:"auto_renew" env:"AUTO_RENEW"` // 自动续期 token
	Trace     bool `toml:"trace" json:"trace" yaml:"trace" env:"TRACE"`

	// 内部状态
	client      *vault.Client
	log         *zerolog.Logger
	stopRenewal chan struct{}
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
