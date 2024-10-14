package crypt

type Crypt interface {
	GenKeyPair() ([]byte, []byte, error) // 生成公钥和私钥，返回公钥和私钥的字节数组
	Init() error                         // 初始化
	Sign(raw string) (string, error)     // 签名
	Verify(plaintext, sign string) error // 验证签名
}
