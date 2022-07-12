package server

type Server interface {
	// Init 初始化服务
	Init(...Option)
	// Options 获取服务的配置选项
	Options() Options
	// Handle 设置服务的 handler
	Handle(Handler) error
	// NewHandler 创建一个 handler 包装
	NewHandler(interface{}) Handler
	// Start 启动服务
	Start() error
	// Stop 实现服务优雅退出
	Stop() error
}

type Option func(*Options)

// Handler handler 的包装接口
type Handler interface {
	Handler() interface{}
}
