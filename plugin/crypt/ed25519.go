package crypt

import (
	"context"
	"crypto/ed25519"
	"crypto/rand"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"fmt"
	"os"

	"github.com/lsjhtang/kid/logger"

	"github.com/lsjhtang/kid/injector"

	"github.com/lsjhtang/kid/config"
	"github.com/lsjhtang/kid/utils"
)

type Ed25519 struct {
	privateKey *ed25519.PrivateKey
	publicKey  *ed25519.PublicKey
	options    *Ed25519Options
}

func (e *Ed25519) Provide(ctx context.Context) interface{} {
	var confName string
	if name, ok := ctx.Value(injector.NameKey{}).(string); ok && len(name) > 0 {
		confName = "." + name
	}
	confPrefix := fmt.Sprintf("crypt.ed25519%s", confName)
	if !config.Exist(confPrefix) {
		panic(fmt.Sprintf("config.yaml file not found configuration item [%s]", confPrefix))
	}
	return NewEd25519(
		WithEd25519PrivateKey(config.GetString(utils.GetConfigurationItem(confPrefix, "privateKey"))),
		WithEd25519PrivateKeyFile(config.GetString(utils.GetConfigurationItem(confPrefix, "privateKeyFile"))),
		WithEd25519PublicKey(config.GetString(utils.GetConfigurationItem(confPrefix, "publicKey"))),
		WithEd25519PublicKeyFile(config.GetString(utils.GetConfigurationItem(confPrefix, "publicKeyFile"))),
	)
}

func (e *Ed25519) GenKeyPair() ([]byte, []byte, error) {
	privateKey, publicKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return nil, nil, fmt.Errorf("生成密钥对失败: %s", err.Error())
	}
	return []byte(hex.EncodeToString(privateKey)), []byte(hex.EncodeToString(publicKey)), nil
}

func (e *Ed25519) Init() error {
	pubKey, err := e.LoadPublicKey()
	if err != nil {
		return err
	}
	e.publicKey = pubKey
	priKey, err := e.LoadPrivateKey()
	if err != nil {
		return err
	}
	e.privateKey = priKey
	return nil
}

func (e *Ed25519) LoadPublicKey() (*ed25519.PublicKey, error) {
	if len(e.options.PublicKey) > 0 {
		publicKey, err := x509.ParsePKIXPublicKey([]byte(e.options.PublicKey))
		if err != nil {
			return nil, err
		}
		ed25519PublicKey, ok := publicKey.(ed25519.PublicKey)
		if !ok {
			return nil, fmt.Errorf("公钥不是ed25519类型")
		}
		return &ed25519PublicKey, nil
	}
	if len(e.options.PublicKeyFile) == 0 {
		return nil, fmt.Errorf("公钥文件未配置")
	}
	pemData, err := os.ReadFile(e.options.PublicKeyFile)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(pemData)
	if block == nil || block.Type != "PUBLIC KEY" {
		return nil, fmt.Errorf("解析公钥文件失败")
	}

	publicKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	ed25519PublicKey, ok := publicKey.(ed25519.PublicKey)
	if !ok {
		return nil, fmt.Errorf("公钥不是ed25519类型")
	}

	return &ed25519PublicKey, nil
}

func (e *Ed25519) LoadPrivateKey() (*ed25519.PrivateKey, error) {
	if len(e.options.PrivateKey) > 0 {
		privateKey, err := x509.ParsePKCS8PrivateKey([]byte(e.options.PrivateKey))
		if err != nil {
			return nil, err
		}
		ed25519PrivateKey, ok := privateKey.(ed25519.PrivateKey)
		if !ok {
			return nil, fmt.Errorf("私钥不是ed25519类型")
		}
		return &ed25519PrivateKey, nil
	}
	pemData, err := os.ReadFile(e.options.PrivateKeyFile)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(pemData)
	if block == nil || block.Type != "PRIVATE KEY" {
		return nil, fmt.Errorf("解析私钥文件失败")
	}

	privateKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	ed25519PrivateKey, ok := privateKey.(ed25519.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("私钥不是ed25519类型")
	}

	return &ed25519PrivateKey, nil
}

func (e *Ed25519) Sign(raw string) (string, error) {
	signature := ed25519.Sign(*e.privateKey, []byte(raw))
	return string(signature), nil
}

func (e *Ed25519) Verify(plaintext, sign string) error {
	if !ed25519.Verify(*e.publicKey, []byte(plaintext), []byte(sign)) {
		return fmt.Errorf("验证签名失败")
	}
	return nil
}

func NewEd25519(opts ...Ed25519Option) *Ed25519 {
	var options Ed25519Options
	for _, opt := range opts {
		opt(&options)
	}
	e := &Ed25519{
		options: &options,
	}
	if err := e.Init(); err != nil {
		logger.Errorf("init rsa crypt failed: %s", err.Error())
		return nil
	}
	return e
}

type Ed25519Option func(*Ed25519Options)

type Ed25519Options struct {
	PrivateKey     string
	PrivateKeyFile string
	PublicKey      string
	PublicKeyFile  string
}

func WithEd25519PublicKey(publicKey string) Ed25519Option {
	return func(o *Ed25519Options) {
		o.PublicKey = publicKey
	}
}

func WithEd25519PublicKeyFile(publicKeyFile string) Ed25519Option {
	return func(o *Ed25519Options) {
		o.PublicKeyFile = publicKeyFile
	}
}

func WithEd25519PrivateKey(privateKey string) Ed25519Option {
	return func(o *Ed25519Options) {
		o.PrivateKey = privateKey
	}
}

func WithEd25519PrivateKeyFile(privateKeyFile string) Ed25519Option {
	return func(o *Ed25519Options) {
		o.PrivateKeyFile = privateKeyFile
	}
}
