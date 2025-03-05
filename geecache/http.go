package geecache

import (
	"log"
	"net/http"
	"strings"
)

type HttpPool struct {
	self     string
	basePath string
}

var defaultBasePath = "/_geecache/"

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
