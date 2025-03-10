package main

import (
	"flag"
	"fmt"
	"geecache"
	"log"
	"net/http"
	"strings"
)

// 这个他是直接监听9999端口，然后通过http协议访问
//
//	func main() {
//		db := map[string]string{
//			"Tom":  "630",
//			"Jack": "589",
//			"Sam":  "567",
//		}
//		geecache.NewGroup("scores", geecache.GetterFunc(func(key string) ([]byte, error) {
//			fmt.Println("search key:", key)
//			if v, ok := db[key]; ok {
//				return []byte(v), nil
//			}
//			return nil, fmt.Errorf("%s not exist", key)
//		}), 0)
//		addr := "localhost:9999"
//		log.Println("geecache is running at", addr)
//		p := geecache.NewHttpPool(addr)
//		http.ListenAndServe(addr, p)
//	}
var db = map[string]string{
	"Tom":  "630",
	"Jack": "589",
	"Sam":  "567",
}

func creatGroup() *geecache.Group {
	return geecache.NewGroup("scores", geecache.GetterFunc(func(key string) ([]byte, error) {
		fmt.Println("search key:", key)
		if v, ok := db[key]; ok {
			return []byte(v), nil
		}
		return nil, fmt.Errorf("%s not exist", key)
	}), 0)
}
func startCacheServer(addr string, addrs []string, gee *geecache.Group) {
	log.Println("geecache is running at", addr)
	//p是一个httpPool
	p := geecache.NewHttpPool(addr)
	p.Set(addrs...)
	//*httpPool实现了 PeerPicker接口
	gee.RegisterPeers(p)
	addr = strings.TrimPrefix(addr, "http://")
	err := http.ListenAndServe(addr, p)
	//http.ListenAndServe(addr[7:], p)
	if err != nil {
		log.Fatal(err)
	}
}
func startApi(apiAddr string, gee *geecache.Group) {
	http.Handle("/api", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//print("nononono")
		key := r.URL.Query().Get("key")
		//fmt.Printf("%v", gee)
		value, err := gee.Get(key)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/octet-stream")
		w.Write(value.ByteSlice())
	}))
	log.Printf("api server is running at %s", apiAddr)
	apiAddr = strings.TrimPrefix(apiAddr, "http://")
	http.ListenAndServe(apiAddr, nil)
}
func main() {
	var api bool
	var port int
	flag.BoolVar(&api, "api", false, "Start a api server")
	flag.IntVar(&port, "port", 8001, "Cache server port")
	flag.Parse()
	gee := creatGroup()
	//fmt.Printf("wuxiwuxi%vheihei", gee)
	apiaddr := "http://localhost:9999"
	addrmap := map[int]string{
		8001: "http://localhost:8001",
		8002: "http://localhost:8002",
		8003: "http://localhost:8003",
	}
	var addrs []string
	for _, v := range addrmap {
		addrs = append(addrs, v)
	}
	if api {
		go startApi(apiaddr, gee)
	}
	startCacheServer(addrmap[port], addrs, gee)
}
