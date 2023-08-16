package crypt

type Crypt interface {
	Init() error
	Sign(raw string) (string, error)
	Verify(plaintext, sign string) error
}
