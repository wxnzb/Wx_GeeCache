package geecache

import (
	"fmt"
	"geecache/consistenthash"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
)

type HttpPool struct {
	self        string
	basePath    string
	mu          sync.Mutex
	peers       *consistenthash.Map
	httpGetters map[string]*httpGetter //现在又忘了这个是干啥的了
}

// 等于是有两个环，一个环上面存的是远程节点，另外的环上面存的是每个远程节点对应的k,v
var defaultBasePath = "/_geecache/"
var defaultReplicas = 50

func NewHttpPool(self string) *HttpPool {
	return &HttpPool{
		self:     self,
		basePath: defaultBasePath,
	}
}

func (p *HttpPool) Log(format string, v ...interface{}) {
	log.Printf("[Server %s] %s", p.self, fmt.Sprintf(format, v...))
}
func (p *HttpPool) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//log.Printf("%s %s", r.Method, r.URL.Path)
	//r.URL.Path[len(p.basePath)：]表示去掉basepath后的路径
	p.Log("%s %s", r.Method, r.URL.Path)
	parts := strings.SplitN(r.URL.Path[len(p.basePath):], "/", 2)
	groupname := parts[0]
	key := parts[1]
	gee := GetGroup(groupname)
	value, err := gee.Get(key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Write(value.ByteSlice())
}

// 实现客户端
type httpGetter struct {
	baseUrl string
}

// // PeerPicker：根据 key 选择合适的缓存节点。
// // PeerGetter：向选中的缓存节点获取数据（可能是 HTTP 或 RPC 请求）。
// // 假设我们有 3 台缓存服务器（A、B、C），其中：
// // A 需要获取 key1，但它自己没有缓存这个 key。
// // 通过 PickPeer(key1) 发现 key1 在 B 上。
// // 然后，A 通过 PeerGetter.Get("group1", "key1") 从 B 获取数据。
// type PeerPicker interface {
// 	Pickpeer(key string) (PeerGetter, bool)
// }
// type PeerGetter interface {
// 	Get(group, key string) ([]byte, error)
// }

func (p *HttpPool) Pickpeer(key string) (PeerGetter, bool) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if peer := p.peers.Get(key); peer != "" && peer != p.self {
		log.Printf("Pick peer %s", peer)
		return p.httpGetters[peer], true
	}
	return nil, false
}

var _PeerPicker = (*HttpPool)(nil)

// 这个函数就是找到远程接口后那个接口上的环中找
func (h *httpGetter) Get(groupname, key string) ([]byte, error) {
	// print("00000")
	// fmt.Printf("11111%s", h.baseUrl)
	// fmt.Printf("22222%s", groupname)
	// fmt.Printf("33333%s", key)
	//我找你这个bug找了这么长时间
	u := fmt.Sprintf("%v%v/%v", h.baseUrl, url.QueryEscape(groupname), url.QueryEscape(key))
	//这里是关键，你必须先开启对这个端口的监听才行
	resp, err := http.Get(u)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server returned:%v", resp.Status)
	}
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response body:%v", err)
	}
	return bytes, nil
}

var _peerGetter = (*httpGetter)(nil)

func (p *HttpPool) Set(peers ...string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	//将每个远程节点存进环中
	p.peers = consistenthash.NewMap(defaultReplicas, nil)
	p.peers.Add(peers...)
	p.httpGetters = make(map[string]*httpGetter)
	for _, peer := range peers {
		//每个远程节点 创建一个 httpGetter，用于 通过 HTTP 访问该远程节点，实现分布式缓存的远程数据获取
		p.httpGetters[peer] = &httpGetter{baseUrl: peer + p.basePath}
	}
}
