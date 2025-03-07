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
	httpGetters map[string]*httpGetter
}

var defaultBasePath = "/_geecache/"
var defaultReplicas = 50

func NewHttpPool(self string) *HttpPool {
	return &HttpPool{
		self:     self,
		basePath: defaultBasePath,
	}
}

// func (p *HttpPool)Log()
func (p *HttpPool) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s %s", r.Method, r.URL.Path)
	//r.URL.Path[len(p.basePath)：]表示去掉basepath后的路径
	parts := strings.SplitN(r.URL.Path[len(p.basePath):], "/", 2)
	groupname := parts[0]
	key := parts[1]
	gee := GetGroup(groupname)
	_, _ = gee.Get(key)
}

// 实现客户端
type httpGetter struct {
	baseUrl string
}
type peerGetter interface {
	Get(group, key string) ([]byte, error)
}

func (h *httpGetter) Get(group, key string) ([]byte, error) {
	u := fmt.Sprintf("%v/%v/%v", h.baseUrl, url.QueryEscape(group), url.QueryEscape(key))
	resp, err := http.Get(u)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server returned:%v", resp.Status)
	}
	return ioutil.ReadAll(resp.Body)
}

var _peerGetter = (*httpGetter)(nil)

func (p *HttpPool) Set(peers ...string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.peers = consistenthash.NewMap(defaultReplicas, nil)
	p.peers.Add(peers...)
	p.httpGetters = make(map[string]*httpGetter)
	for _, peer := range peers {
		p.httpGetters[peer] = &httpGetter{baseUrl: peer + p.basePath}
	}
}

type PeerPicker interface {
	Pickpeer(key string) (*httpGetter, bool)
}

func (p *HttpPool) Pickpeer(key string) (*httpGetter, bool) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if peer := p.peers.Get(key); peer != "" && peer != p.self {
		return p.httpGetters[peer], true
	}
	return nil, false
}

var _PeerPicker = (*HttpPool)(nil)
