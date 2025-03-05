package geecache

import (
	"fmt"
	"log"
	"reflect"
	"testing"
)

func Test_Getter(t *testing.T) {
	var f = GetterFunc(func(key string) ([]byte, error) {
		return []byte(key), nil
	})
	expect := []byte("key")
	if v, _ := f.Get("key"); !reflect.DeepEqual(v, expect) {
		t.Errorf("got %v, expect %v", v, expect)
	}
}

var db = map[string]string{
	"wwxx": "666",
	"wx":   "580",
	"w大x":  "588",
}

func Test_GetterFunc(t *testing.T) {
	var dbcounts = make(map[string]int, len(db))
	g := NewGroup("scores", GetterFunc(func(key string) ([]byte, error) {
		if v, ok := db[key]; ok {
			log.Println("get from db", key)
			//这里要是db里面没有比如jj对应的value,在第一次没命中想从db中找到加进去时，也找不到，所以一直都命中不了
			if _, ok := dbcounts[key]; !ok {
				dbcounts[key] = 0
			} else {
				dbcounts[key]++
			}
			// if _, ok := dbcounts[key]; !ok {
			// 	dbcounts[key] = 1
			// }
			return []byte(v), nil
		}
		return nil, fmt.Errorf("no value for %s", key)
	}), 2<<10) //为啥设置2<<10
	for k, v := range db {
		//这里返回的val时ByteView类型的，要转成string，需要自己实现
		if val, err := g.Get(k); err != nil || v != val.String() {
			t.Fatal("failed to get value of wwxx")
		}
		if _, err := g.Get(k); err != nil {
			t.Fatal("failed to get value of wwxx")
		}
	}
}
