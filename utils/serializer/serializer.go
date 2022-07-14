package serializer

type Serializer interface {
	Encode(interface{}) (string, error)
	Decode([]byte) (interface{}, error)
}
