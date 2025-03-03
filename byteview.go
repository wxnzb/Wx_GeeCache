package geecache

type ByteView struct {
	b []byte
}

// 实现Value接口
func (bv ByteView) Len() int {
	return len(bv.b)
}
