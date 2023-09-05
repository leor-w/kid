package rabbitmq

type Options struct {
	User   string
	Pwd    string
	Host   string
	Port   int
	Queues []*Queue
}

type Queue struct {
	Name     string
	Exchange string
	Key      string
}
