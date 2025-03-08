package main

import (
	"fmt"
	"geecache"
	"log"
	"net/http"
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
func main() {
	db := map[string]string{
		"Tom":  "630",
		"Jack": "589",
		"Sam":  "567",
	}
	gee := geecache.NewGroup("scores", geecache.GetterFunc(func(key string) ([]byte, error) {
		fmt.Println("search key:", key)
		if v, ok := db[key]; ok {
			return []byte(v), nil
		}
		return nil, fmt.Errorf("%s not exist", key)
	}), 0)
	addr := "localhost:9999"
	log.Println("geecache is running at", addr)
	//p是一个httpPool
	p := geecache.NewHttpPool(addr)
	addrmap := map[int]string{
		8001: "localhost:8001",
		8002: "localhost:8002",
		8003: "localhost:8003",
	}
	var addrs []string
	for _, v := range addrmap {
		addrs = append(addrs, v)
	}
	p.Set(addrs...)
	//*httpPool实现了 PeerPicker接口
	gee.RegisterPeers(p)
	http.ListenAndServe(addr, p)
}
