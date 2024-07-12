package lru

import "container/list"

// 缓存
type Cache struct {
	//最大缓存
	maxBytes int
	//当前缓存
	nowBytes int

	//lru链表 方便内存不足时能淘汰元素
	//这里其实可以看出来，golang没泛型是真的坑，这样list里放的什么都不知道
	lruList *list.List
	//缓存本身 缓存存的是lru list的节点，而非value，一般来说list.element是在value的基础上加了pre.next等链表信息
	cache map[string]*list.Element
}

// value必须要实现len，这样才能判断出内存占用多少
type Value interface {
	Len() int
}

//list内的节点（其实也是map的value结构）
//这里之所以要包一层，里面加上key是为了删除缓存的时候，根据lruList拿到要删除的节点 可以通过key去map里再删除一遍

type Node struct {
	key   string
	value Value
}

func New(maxBytes int) *Cache {
	return &Cache{
		maxBytes: maxBytes,
		nowBytes: 0,
		lruList:  list.New(),
		cache:    make(map[string]*list.Element),
	}
}

func (cache *Cache) Add(key string, value Value) {
	//如果存在走update
	if element, ok := cache.cache[key]; ok {
		//按更新处理
		node := element.Value.(*Node)
		node.value = value
		//更新LRU认为是最新
		cache.lruList.MoveToFront(element)
		//更新占用字节
		cache.nowBytes += value.Len() - node.value.Len()
	} else {
		//先加入LRU缓存
		node := &Node{key, value}
		element := cache.lruList.PushFront(node)
		//再写入缓存
		cache.cache[key] = element
		//更新内存
		cache.nowBytes += value.Len() + len(key)
	}
	cache.lru()
}

func (cache *Cache) Delete(key string) Value {
	if element, ok := cache.cache[key]; ok {
		node := element.Value.(*Node)
		//删除缓存
		delete(cache.cache, key)
		//删除LRU列表
		cache.lruList.Remove(element)
		//更新字节占用
		cache.nowBytes -= node.value.Len()
		return node.value
	}
	return nil
}

func (cache *Cache) Get(key string) Value {
	if element, ok := cache.cache[key]; ok {
		cache.lruList.MoveToFront(element)
		node := element.Value.(*Node)
		return node.value
	}
	return nil
}

func (cache *Cache) lru() {
	for cache.maxBytes > 0 && cache.nowBytes > cache.maxBytes {
		cache.rmOldest()
	}
}

func (cache *Cache) rmOldest() {
	element := cache.lruList.Back()
	if element != nil {
		cache.lruList.Remove(element)
		node := element.Value.(*Node)
		delete(cache.cache, node.key)
		cache.nowBytes -= len(node.key) + node.value.Len()
	}
}
