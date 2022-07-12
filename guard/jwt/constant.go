package jwt

type SigningMethod string

const (
	SigningMethodSH256 SigningMethod = "sha256"
	SigningMethodSH384 SigningMethod = "sha384"
	SigningMethodSH512 SigningMethod = "sha512"
)

const (
	BlacklistKey = "blacklist.key.%s"
)
