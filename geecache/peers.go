package geecache

// PeerPicker：根据 key 选择合适的缓存节点。
// PeerGetter：向选中的缓存节点获取数据（可能是 HTTP 或 RPC 请求）。
// 假设我们有 3 台缓存服务器（A、B、C），其中：
// A 需要获取 key1，但它自己没有缓存这个 key。
// 通过 PickPeer(key1) 发现 key1 在 B 上。
// 然后，A 通过 PeerGetter.Get("group1", "key1") 从 B 获取数据。
type PeerPicker interface {
	Pickpeer(key string) (PeerGetter, bool)
}
type PeerGetter interface {
	Get(group, key string) ([]byte, error)
}
