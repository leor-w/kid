package signature

import (
	"crypto"
	"crypto/hmac"
	"encoding/base64"
	"errors"
)

type HMAC struct {
	name string
	Hash crypto.Hash
}

func (h *HMAC) Verify(verify, signed string, key []byte) error {
	sig, err := base64.StdEncoding.DecodeString(signed)
	if err != nil {
		return err
	}
	if !h.Hash.Available() {
		return ErrInvalidKeyType
	}
	hasher := hmac.New(h.Hash.New, key)
	hasher.Write([]byte(verify))
	if !hmac.Equal(sig, hasher.Sum(nil)) {
		return ErrSignatureInvalid
	}
	return nil
}

func (h *HMAC) Sign(signingString string, key []byte) (string, error) {
	if !h.Hash.Available() {
		return "", ErrHashUnavailable
	}
	hasher := hmac.New(h.Hash.New, key)
	hasher.Write([]byte(signingString))
	return base64.StdEncoding.EncodeToString(hasher.Sum(nil)), nil
}

func (h *HMAC) Name() string {
	return h.name
}

var (
	SignHAMCHS256       *HMAC
	SignHMACHS384       *HMAC
	SignHMACHS512       *HMAC
	ErrSignatureInvalid = errors.New("signature is invalid")
)

const (
	HS256 = "HS256"
	HS384 = "HS384"
	HS512 = "HS512"
)

func init() {
	// HS256
	SignHAMCHS256 = &HMAC{HS256, crypto.SHA256}
	Registry(SignHAMCHS256.Name(), func() Signatory {
		return SignHAMCHS256
	})

	// HS384
	SignHMACHS384 = &HMAC{HS384, crypto.SHA384}
	Registry(SignHMACHS384.Name(), func() Signatory {
		return SignHMACHS384
	})

	// HS512
	SignHMACHS512 = &HMAC{HS512, crypto.SHA512}
	Registry(SignHMACHS512.Name(), func() Signatory {
		return SignHMACHS512
	})
}
