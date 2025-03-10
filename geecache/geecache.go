package geecache

import (
	"errors"
	"geecache/singleflight"
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
	loader *singleflight.Group
}

func (g *Group) RegisterPeers(peers PeerPicker) {
	g.peers = peers
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
		loader: &singleflight.Group{},
	}
	groups[name] = g
	return g
}
func GetGroup(name string) *Group {
	mu.RLock()
	defer mu.RUnlock()
	if groups[name] == nil {
		log.Print("你是最棒的")
	}
	return groups[name]
}
func (g *Group) Get(key string) (ByteView, error) {
	//fmt.Printf("key=%s", key)
	if key == "" {
		return ByteView{}, errors.New("key is required")
	}
	if bv, ok := (&g.mainca).get(key); ok {
		log.Println("[GeeCache] hit")
		return bv, nil
	}
	//没有命中，没有的话就要先在远程节点中找，要是找不到就调用回调函数
	return g.load(key)
}
func (g *Group) load(key string) (ByteView, error) {
	view, err := g.loader.Do(key, func() (interface{}, error) {
		if g.peers != nil {
			if peer, ok := g.peers.Pickpeer(key); ok {
				//fmt.Printf("ooo%s %v", g.name, key)
				if value, err := peer.Get(g.name, key); err == nil {
					return ByteView{b: value}, nil
				}
			}
		}
		return g.getLocally(key)
	})
	if err != nil {
		return ByteView{}, err
	}
	return view.(ByteView), nil
}

// 根据key获取本地数据
func (g *Group) getLocally(key string) (ByteView, error) {
	// 调用getter的Get方法获取key对应的value
	bytes, err := g.getter.Get(key)
	if err != nil {
		return ByteView{}, err
	}
	//这里记得先克隆
	value := ByteView{b: cloneBytes(bytes)}
	g.mainca.add(key, value)
	return value, nil
}

//现在就是还没有实现getter
