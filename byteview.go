package geecache

type ByteView struct {
	b []byte
}

// 实现Value接口
func (bv ByteView) Len() int {
	return len(bv.b)
}

//这个下面还没有实现完，感觉现在不需要
