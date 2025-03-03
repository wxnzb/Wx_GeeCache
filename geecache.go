package geecache

// 咱们之前实现的都是结构体实现接口
// 而现在这个是函数类型实现接口
type Getter interface {
	Get(key string) ([]byte, error)
}
type GetterFunc func(kry string) ([]byte, error)

func (f GetterFunc) Get(key string) ([]byte, error) {
	return f(key)
}
