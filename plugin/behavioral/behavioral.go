package behavioral

type Behavioral interface {
	Verify(ip, data string) (string, error)
}
