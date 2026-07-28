package crypt

import (
	"context"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/lsjhtang/kid/injector"

	"github.com/lsjhtang/kid/config"
	"github.com/lsjhtang/kid/logger"
	"github.com/lsjhtang/kid/utils"
)

type Rsa struct {
	publicKey  *rsa.PublicKey
	privateKey *rsa.PrivateKey
	options    *RsaOptions
}

func (r *Rsa) Provide(ctx context.Context) interface{} {
	var confName string
	if name, ok := ctx.Value(injector.NameKey{}).(string); ok && len(name) > 0 {
		confName = "." + name
	}
	confPrefix := fmt.Sprintf("rsa%s", confName)
	if !config.Exist(confPrefix) {
		panic(fmt.Sprintf("配置文件未找到配置项 [%s]", confPrefix))
	}
	return NewRsa(
		WithPrivateKey(config.GetString(utils.GetConfigurationItem(confPrefix, "privateKey"))),
		WithPrivateKeyFile(config.GetString(utils.GetConfigurationItem(confPrefix, "privateKeyFile"))),
		WithPublicKey(config.GetString(utils.GetConfigurationItem(confPrefix, "publicKey"))),
		WithPublicKeyFile(config.GetString(utils.GetConfigurationItem(confPrefix, "publicKeyFile"))),
	)
}

func (r *Rsa) GenKeyPair() ([]byte, []byte, error) {
	return nil, nil, nil
}

func (r *Rsa) Init() error {
	pubKey, err := r.LoadPublicKey()
	if err != nil {
		return err
	}
	r.publicKey = pubKey
	priKey, err := r.LoadPrivateKey()
	if err != nil {
		return err
	}
	r.privateKey = priKey
	return nil
}

func (r *Rsa) LoadPublicKey() (*rsa.PublicKey, error) {
	var (
		publicKey []byte
		err       error
	)
	if len(r.options.PublicKey) > 0 {
		publicKey = []byte(r.options.PublicKey)
	}
	if len(r.options.PublicKeyFile) > 0 {
		publicKey, err = ioutil.ReadFile(r.options.PublicKeyFile)
		if err != nil {
			return nil, err
		}
	}
	block, _ := pem.Decode(publicKey)
	if block == nil {
		return nil, errors.New("load public key failed")
	}
	pubKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return pubKey.(*rsa.PublicKey), nil
}

func (r *Rsa) LoadPrivateKey() (*rsa.PrivateKey, error) {
	var (
		privateKey []byte
		err        error
	)
	if len(r.options.PrivateKey) != 0 {
		privateKey = []byte(r.options.PrivateKey)
	}
	if len(r.options.PrivateKeyFile) != 0 {
		privateKey, err = os.ReadFile(r.options.PrivateKeyFile)
		if err != nil {
			return nil, err
		}
	}
	if len(privateKey) == 0 {
		return nil, errors.New("private key is empty")
	}
	block, _ := pem.Decode(privateKey)
	if block == nil {
		return nil, errors.New("private key parse failed")
	}
	priKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return priKey, nil
}

func (r *Rsa) Sign(raw string) (string, error) {
	hasher := sha256.New()
	hasher.Write([]byte(raw))
	hashed := hasher.Sum(nil)
	signed, err := rsa.SignPKCS1v15(rand.Reader, r.privateKey, crypto.SHA256, hashed)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(signed), nil
}

func (r *Rsa) Verify(plaintext, sign string) error {
	plaintextHashed := sha256.Sum256([]byte(plaintext))
	decodeSign, err := base64.StdEncoding.DecodeString(sign)
	if err != nil {
		return err
	}
	if err = rsa.VerifyPKCS1v15(r.publicKey, crypto.SHA256, plaintextHashed[:], decodeSign); err != nil {
		return err
	}
	return nil
}

func NewRsa(opts ...RsaOption) *Rsa {
	var options RsaOptions
	for _, opt := range opts {
		opt(&options)
	}
	rsaCrypt := &Rsa{
		options: &options,
	}
	if err := rsaCrypt.Init(); err != nil {
		logger.Errorf("init rsa crypt failed: %s", err.Error())
		return nil
	}
	return rsaCrypt
}

type RsaOption func(*RsaOptions)

type RsaOptions struct {
	PublicKey      string
	PublicKeyFile  string
	PrivateKey     string
	PrivateKeyFile string
}

func WithPublicKey(publicKey string) RsaOption {
	return func(o *RsaOptions) {
		o.PublicKey = publicKey
	}
}

func WithPublicKeyFile(publicKeyFile string) RsaOption {
	return func(o *RsaOptions) {
		o.PublicKeyFile = publicKeyFile
	}
}
func WithPrivateKey(privateKey string) RsaOption {
	return func(o *RsaOptions) {
		o.PrivateKey = privateKey
	}
}

func WithPrivateKeyFile(privateKeyFile string) RsaOption {
	return func(o *RsaOptions) {
		o.PrivateKeyFile = privateKeyFile
	}
}
