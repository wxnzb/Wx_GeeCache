package main

import (
	"fmt"
	"geecache"
	"log"
	"net/http"
)

func main() {
	db := map[string]string{
		"Tom":  "630",
		"Jack": "589",
		"Sam":  "567",
	}
	geecache.NewGroup("scores", geecache.GetterFunc(func(key string) ([]byte, error) {
		fmt.Println("search key:", key)
		if v, ok := db[key]; ok {
			return []byte(v), nil
		}
		return nil, fmt.Errorf("%s not exist", key)
	}), 0)
	addr := "localhost:9999"
	log.Println("geecache is running at", addr)
	p := geecache.NewHttpPool(addr)
	http.ListenAndServe(addr, p)
}
