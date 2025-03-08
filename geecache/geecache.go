package geecache

import (
	"errors"
	"log"
	"sync"
)

// 咱们之前实现的都是结构体实现接口
// 而现在这个是函数类型实现接口
type Getter interface {
	Get(key string) ([]byte, error)
}
type GetterFunc func(kry string) ([]byte, error)

func (f GetterFunc) Get(key string) ([]byte, error) {
	return f(key)
}

// Group结构体
var mu sync.RWMutex
var groups = make(map[string]*Group)

type Group struct {
	name   string
	getter Getter
	mainca cache
	peers  PeerPicker
}

func (g *Group) RegisterPeers(peers PeerPicker) {
	g.peers = peers
}
func (g *Group) GetFromPeer(key string) (ByteView, error) {
	if g.peers != nil {
		//找到key对应的远程节点，在从远程节点中找到key对应的value
		if peer, ok := g.peers.Pickpeer(key); ok {
			if value, err := peer.Get(g.name, key); err == nil {
				return ByteView{b: value}, nil
			}
		}
	}
	//如果远程没有，就从本地获取
	return g.Get(key)
}
func NewGroup(name string, getter Getter, maxsize int) *Group {
	//这里为啥一定要有回调函数
	if getter == nil {
		panic("nil Getter")
	}
	mu.Lock()
	defer mu.Unlock()
	g := &Group{
		name:   name,
		getter: getter,
		mainca: cache{maxSize: maxsize},
	}
	groups[name] = g
	return g
}
func GetGroup(name string) *Group {
	mu.RLock()
	defer mu.RUnlock()
	return groups[name]
}
func (g *Group) Get(key string) (ByteView, error) {
	if key == "" {
		return ByteView{}, errors.New("key is required")
	}
	if bv, ok := (&g.mainca).get(key); ok {
		log.Println("[GeeCache] hit")
		return bv, nil
	}
	//没有命中，没有的话就要通过回调函数获取并加入数据库
	return g.load(key)
}
func (g *Group) load(key string) (ByteView, error) {
	bv, err := g.getter.Get(key)
	if err != nil {
		return ByteView{}, err
	}
	value := ByteView{b: cloneBytes(bv)}
	g.mainca.add(key, value)
	return value, nil
}

//现在就是还没有实现getter
