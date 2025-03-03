package lru

import (
	"reflect"
	"testing"
)

type String string

func (d String) Len() int {
	return len(d)
}
func Test_Get(t *testing.T) {
	ca := NewCache(0, nil)
	ca.Add("wx", String("good"))
	//需要变成
	if v, ok := ca.Get("wx"); !ok || v.(String) != "good" {
		t.Fatal("get fail")
	}
	//下面这个肯定让不能通过
	if _, ok := ca.Get("w,,x"); !ok {
		t.Fatal("get fail,No key")
	}
}
func Test_DelectOldset(t *testing.T) {
	k1, k2, k3 := "key1", "key2", "k3"
	v1, v2, v3 := "value1", "value2", "v3"
	len := len(k1 + k2 + v1 + v2)
	ca := NewCache(len, nil)
	ca.Add(k1, String(v1))
	ca.Add(k2, String(v2))
	ca.Add(k3, String(v3))
	if _, ok := ca.Get(k1); ok || ca.Len() != 2 {
		t.Fatal("delect fail")
	}
}

// 测试回调函数是否可以被调用，我不知道要这个回调函数有什么用
func Test_Onevicted(t *testing.T) {
	keys := make([]string, 0)
	callback := func(key string, value Value) {
		keys = append(keys, key)
	}
	ca := NewCache(10, callback)
	ca.Add("key1", String("123456"))
	ca.Add("k2", String("v2"))
	ca.Add("k3", String("v3"))
	ca.Add("k4", String("v4"))
	expect := []string{"key1", "k2"}
	//感觉key2不会被淘汰
	if !reflect.DeepEqual(expect, keys) {
		t.Fatalf("keys not equal:%s", expect)
	}
}
