package application

import (
	"encoding/base64"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	aesgcm "github.com/infraboard/mcube/v2/crypto/aes_gcm"
	"github.com/infraboard/mcube/v2/crypto/cbc"
	"github.com/infraboard/mcube/v2/ioc"
)

func init() {
	ioc.Config().Registry(defaultConfig)
}

var defaultConfig = &Application{
	AppName:          "",
	AppGroup:         "default",
	EncryptKey:       "defualt app encrypt key",
	CipherPrefix:     "@ciphered@",
	EncryptAlgorithm: ENCRYPT_ALGORITHM_AES_GCM,
	KeyLength:        aesgcm.AES256,
	Extras:           map[string]string{},
}

type Application struct {
	ioc.ObjectImpl

	AppGroup         string            `json:"group" yaml:"group" toml:"group" env:"GROUP"`
	AppName          string            `json:"name" yaml:"name" toml:"name" env:"NAME"`
	AppDescription   string            `json:"description" yaml:"description" toml:"description" env:"DESCRIPTION"`
	AppAddress       string            `json:"address" yaml:"address" toml:"address" env:"ADDRESS"`
	Debug            bool              `json:"debug" yaml:"debug" toml:"debug" env:"DEBUG"`
	InternalAddress  string            `json:"internal_address" yaml:"internal_address" toml:"internal_address" env:"INTERNAL_ADDRESS"`
	InternalToken    string            `json:"internal_token" yaml:"internal_token" toml:"internal_token" env:"INTERNAL_TOKEN"`
	EncryptAlgorithm EncryptAlgorithm  `json:"encrypt_algorithm" yaml:"encrypt_algorithm" toml:"encrypt_algorithm" env:"ENCRYPT_ALGORITHM"`
	KeyLength        aesgcm.KeySize    `json:"key_length" yaml:"key_length" toml:"key_length" env:"KEY_LENGTH"`
	EncryptKey       string            `json:"encrypt_key" yaml:"encrypt_key" toml:"encrypt_key" env:"ENCRYPT_KEY"`
	CipherPrefix     string            `json:"cipher_prefix" yaml:"cipher_prefix" toml:"cipher_prefix" env:"CIPHER_PREFIX"`
	Extras           map[string]string `json:"extras" yaml:"extras" toml:"extras" env:"EXTRAS"`

	appURL *url.URL
}

func (i *Application) GetExtras(key string) string {
	if i.Extras == nil {
		i.Extras = map[string]string{}
	}
	return i.Extras[key]
}

func (i *Application) Domain() string {
	if i.appURL != nil {
		return strings.Split(i.appURL.Host, ":")[0]
	}
	return "localhost"
}

func (i *Application) GenerateRandomEncryptKey() string {
	var key []byte
	switch i.EncryptAlgorithm {
	case ENCRYPT_ALGORITHM_AES_CBC:
		key = cbc.GenKeyBySHA1Hash2(i.EncryptKey, cbc.AES_KEY_LEN(i.KeyLength))
	default:
		key = aesgcm.MustGenerateKey(i.KeyLength)
	}
	return base64.StdEncoding.EncodeToString(key)
}

func (i *Application) EncryptByte(plaintext []byte) ([]byte, error) {
	if i.EncryptKey == "" {
		return nil, aesgcm.ErrInvalidKeySize
	}

	key, err := base64.StdEncoding.DecodeString(i.EncryptKey)
	if err != nil {
		return nil, err
	}

	switch i.EncryptAlgorithm {
	case ENCRYPT_ALGORITHM_AES_CBC:
		cbc := cbc.MustNewAESCBCCihper(key)
		return cbc.Encrypt(plaintext)
	default:
		gcm, err := aesgcm.NewAESGCM(key)
		if err != nil {
			return nil, err
		}
		return gcm.Encrypt(plaintext)
	}
}

func (i *Application) EncryptString(plaintext string) (string, error) {
	if i.EncryptKey == "" {
		return "", aesgcm.ErrInvalidKeySize
	}

	key, err := base64.StdEncoding.DecodeString(i.EncryptKey)
	if err != nil {
		return "", err
	}

	switch i.EncryptAlgorithm {
	case ENCRYPT_ALGORITHM_AES_CBC:
		cbc := cbc.MustNewAESCBCCihper(key)
		return cbc.EncryptToString(plaintext)
	default:
		gcm, err := aesgcm.NewAESGCM(key)
		if err != nil {
			return "", err
		}
		return gcm.EncryptToString(plaintext)
	}
}

func (i *Application) GetKeyInfo() (map[string]any, error) {
	if i.EncryptKey == "" {
		return nil, nil
	}

	key, err := base64.StdEncoding.DecodeString(i.EncryptKey)
	if err != nil {
		return nil, err
	}

	switch i.EncryptAlgorithm {
	case ENCRYPT_ALGORITHM_AES_CBC:
		return nil, fmt.Errorf("cbc not key info")
	default:
		gcm, err := aesgcm.NewAESGCM(key)
		if err != nil {
			return nil, err
		}
		return gcm.GetKeyInfoDetailed(), nil
	}
}

func (i *Application) DecryptByte(ciphertext []byte) ([]byte, error) {
	if i.EncryptKey == "" {
		return nil, aesgcm.ErrInvalidKeySize
	}

	key, err := base64.StdEncoding.DecodeString(i.EncryptKey)
	if err != nil {
		return nil, err
	}

	switch i.EncryptAlgorithm {
	case ENCRYPT_ALGORITHM_AES_CBC:
		cbc := cbc.MustNewAESCBCCihper(key)
		return cbc.Decrypt(ciphertext)
	default:
		gcm, err := aesgcm.NewAESGCM(key)
		if err != nil {
			return nil, err
		}
		return gcm.Decrypt(ciphertext)
	}
}

func (i *Application) DecryptString(ciphertext string) (string, error) {
	if i.EncryptKey == "" {
		return "", aesgcm.ErrInvalidKeySize
	}

	key, err := base64.StdEncoding.DecodeString(i.EncryptKey)
	if err != nil {
		return "", err
	}

	switch i.EncryptAlgorithm {
	case ENCRYPT_ALGORITHM_AES_CBC:
		cbc := cbc.MustNewAESCBCCihper(key)
		return cbc.DecryptFromCipherText(ciphertext)
	default:
		gcm, err := aesgcm.NewAESGCM(key)
		if err != nil {
			return "", err
		}
		return gcm.DecryptFromString(ciphertext)
	}
}

func (i *Application) Host() string {
	if i.appURL != nil {
		return i.appURL.Host
	}
	return ""
}

func (i *Application) GetAppName() string {
	if i.AppName != "" {
		return i.AppName
	}

	// 获取当前可执行文件的路径
	executablePath, err := os.Executable()
	if err != nil {
		return ""
	}

	// 获取可执行文件的名称
	return filepath.Base(executablePath)
}

func (i *Application) IsInternalIP() bool {
	return i.Domain() != "localhost" && i.Domain() != "" && i.Domain() != "127.0.0.1"
}

func (i *Application) Init() error {
	sn := os.Getenv("OTEL_SERVICE_NAME")
	if sn == "" {
		os.Setenv("OTEL_SERVICE_NAME", i.AppName)
	}

	if i.AppAddress != "" {
		v, err := url.Parse(i.AppAddress)
		if err != nil {
			return err
		}
		i.appURL = v
	}
	return nil
}

func (i *Application) Name() string {
	return AppName
}

// 优先初始化, 以供后面的组件使用
func (i *Application) Priority() int {
	return 999
}
