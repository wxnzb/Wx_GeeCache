package lru

import (
	"container/list"
)

type Cache struct {
	maxSize     int
	alreadyUsed int
	cache       map[string]*list.Element
	list        *list.List
	Onevicted   func(key string, value Value)
}
type Entry struct {
	value Value //这里将int修改成interface{}接口实现通用性
	key   string
}
type Value interface {
	Len() int
}

// 新建缓存
func NewCache(maxSize int, Onevicted func(key string, value Value)) *Cache {
	return &Cache{
		maxSize:     maxSize,
		alreadyUsed: 0, //不需要？
		cache:       make(map[string]*list.Element),
		list:        list.New(),
		Onevicted:   Onevicted,
	}
}

// 在这个中，用的少的放在链表的后便，用的多的放在链表的前面
// 查找
func (c *Cache) Get(key string) (Value, bool) {
	if ele, ok := c.cache[key]; ok {
		c.list.MoveToFront(ele)
		kv := ele.Value.(*Entry)
		return kv.value, true
	}
	return nil, false
}

// 删除,这里删除实现的是删除最后一项不常用的
func (c *Cache) DeleteOldset() {
	ele := c.list.Back()
	if ele != nil {
		c.list.Remove(ele)
		kv := ele.Value.(*Entry)
		delete(c.cache, kv.key)
		c.alreadyUsed -= kv.value.Len() + len(kv.key)
		if c.Onevicted != nil {
			c.Onevicted(kv.key, kv.value)
		}
	}
}

// 新增或修改,找到就更新value,找不到就加入
// 自认为这个应该也可以
//
//	func (c *Cache)Add(key string,value Value){
//	     if ele,ok:=c.cache[key];ok{
//			if valueAgo,ok:=c.Get(key);ok{
//				kv:=ele.Value.(*Entry)
//				kv.value=value
//				c.alreadyUsed=value.Len()-valueAgo.Len()
//			}
//		 }else{
//			c.list.PushFront((&Entry{key:key,value:value}))
//			c.cache[key]=c.list.Front()
//		 }
//	}
func (c *Cache) Add(key string, value Value) {
	if ele, ok := c.cache[key]; ok {
		kv := ele.Value.(*Entry)
		c.alreadyUsed = value.Len() - kv.value.Len()
		kv.value = value
		c.list.MoveToFront(ele)
	} else {
		c.list.PushFront((&Entry{key: key, value: value}))
		c.cache[key] = c.list.Front()
	}
	//这里为啥要用for?
	//是为了确保 循环删除多余的元素，直到nbytes小于maxBytes
	for c.alreadyUsed > c.maxSize {
		c.DeleteOldset()
	}
}
func (c *Cache) Len() int {
	return c.list.Len()
}
