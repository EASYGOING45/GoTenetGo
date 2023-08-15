package lru

import "container/list"

// Cache is a LRU cache It is not safe for concurrent access.
// 缓存是一个LRU缓存，它不是并发访问安全的。
type Cache struct {
	maxBytes int64                    //允许使用的最大内存
	nbytes   int64                    //当前已使用的内存
	ll       *list.List               //双向链表 用于实现LRU
	cache    map[string]*list.Element //字典，key是字符串，值是双向链表中对应节点的指针
	// optional and executed when an entry is purged. 值被移除时的回调函数，可以为nil
	OnEvicted func(key string, value Value) //某条记录被移除时的回调函数，可以为nil
}

type entry struct {
	key   string
	value Value
}

// Value use Len to count how many bytes it takes 值需要实现Len方法来获取其所占的内存大小
type Value interface {
	Len() int //返回值所占用的内存大小
}

// New is the Constructor of Cache  New是Cache的构造函数
func New(maxBytes int64, onEvicted func(string, Value)) *Cache {
	return &Cache{
		maxBytes:  maxBytes,                       //允许使用的最大内存
		ll:        list.New(),                     //初始化双向链表
		cache:     make(map[string]*list.Element), //初始化字典
		OnEvicted: onEvicted,                      //回调函数
	}
}

// Add adds a value to the cache
func (c *Cache) Add(key string, value Value) {
	if ele, ok := c.cache[key]; ok {
		c.ll.MoveToFront(ele)                                  //如果键存在，则将对应节点移动到队尾
		kv := ele.Value.(*entry)                               //类型断言
		c.nbytes += int64(value.Len()) - int64(kv.value.Len()) //更新已使用的内存
		kv.value = value                                       //更新值
	} else {
		ele := c.ll.PushFront(&entry{key, value})
		c.cache[key] = ele
		c.nbytes += int64(len(key)) + int64(value.Len()) //更新已使用的内存
	}
	for c.maxBytes != 0 && c.maxBytes < c.nbytes {
		c.RemoveOldest() //如果超过了设定的最大内存，则移除最少访问的节点
	}
}

// Get look ups a key's value 查找功能
func (c *Cache) Get(key string) (value Value, ok bool) {
	if ele, ok := c.cache[key]; ok {
		c.ll.MoveToFront(ele)    //将对应节点移动到队尾
		kv := ele.Value.(*entry) //类型断言
		return kv.value, true
	}
	return
}

// RemoveOldest removes the oldest item 移除最少访问的节点（队首）
func (c *Cache) RemoveOldest() {
	ele := c.ll.Back() //获取队首节点
	if ele != nil {
		c.ll.Remove(ele)                                       //从链表中删除该节点
		kv := ele.Value.(*entry)                               //类型断言
		delete(c.cache, kv.key)                                //从字典中删除该节点的映射关系
		c.nbytes -= int64(len(kv.key)) + int64(kv.value.Len()) //更新已使用的内存
		if c.OnEvicted != nil {
			c.OnEvicted(kv.key, kv.value) //如果回调函数OnEvicted不为nil，则调用回调函数
		}
	}
}

// Len the number of cache entries
func (c *Cache) Len() int {
	return c.ll.Len()
}
