package consistenthash

import (
	"hash/crc32"
	"sort"
	"strconv"
)

type Hash func([]byte) uint32
type Map struct {
	replicas int
	hash     Hash
	keys     []int
	hashMap  map[int]string
}

func NewMap(replicas int, fn Hash) *Map {
	m := &Map{
		replicas: replicas,
		hash:     fn,
		hashMap:  make(map[int]string), //m.hashMap[hash] 存储的是对应的节点名称（如 "NodeA"），而不是key
	}
	if m.hash == nil {
		m.hash = crc32.ChecksumIEEE
	}
	return m
}
func (m *Map) Add(keys ...string) {
	for _, key := range keys {
		for i := 0; i < m.replicas; i++ {
			//func Itoa(i int) string
			hash := int(m.hash([]byte(strconv.Itoa(i) + key)))
			m.keys = append(m.keys, hash)
			m.hashMap[hash] = key
		}
	}
	sort.Ints(m.keys)
}
func (m *Map) Get(key string) string {
	if len(m.keys) == 0 {
		return ""
	}
	hash := int(m.hash([]byte(key)))
	idx := sort.Search(len(m.keys), func(i int) bool {
		return m.keys[i] >= hash
	})
	return m.hashMap[m.keys[idx%len(m.keys)]]
}
