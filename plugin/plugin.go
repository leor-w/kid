package plugin

type Plugin interface {
	Provide() interface{}
}
