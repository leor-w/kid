package signature

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"sync"
)

var (
	signatoryMap = SyncMap{data: make(map[string]func() Signatory)}
)

var (
	ErrInvalidKeyType  = errors.New("key is of invalid type")
	ErrHashUnavailable = errors.New("the requested hash function is unavailable")
)

type SyncMap struct {
	sync.RWMutex
	data map[string]func() Signatory
}

type Signatory interface {
	Verify(verify, signed string, key []byte) error
	Sign(signingString string, key []byte) (string, error)
	Name() string
}

func Registry(name string, f func() Signatory) {
	signatoryMap.Lock()
	defer signatoryMap.Unlock()
	signatoryMap.data[name] = f
}

func GetSignatory(name string) Signatory {
	signatoryMap.RLock()
	defer signatoryMap.RUnlock()
	f, exist := signatoryMap.data[name]
	if exist {
		return f()
	}
	return nil
}

func EncodeToURLBase64(data interface{}) (string, error) {
	dBytes, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(dBytes), nil
}

func DecodeURLBase64(raw string, receive interface{}) error {
	bytes, err := base64.URLEncoding.DecodeString(raw)
	if err != nil {
		return err
	}
	return json.Unmarshal(bytes, &receive)
}
