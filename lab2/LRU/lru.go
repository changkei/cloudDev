package lru

import (
	"errors"
	"fmt"
	"strconv"
)

type Cacher interface {
	Get(interface{}) (interface{}, error)
	Put(interface{}, interface{}) error
}

type lruCache struct {
	size      int
	remaining int
	cache     map[string]string
	queue     []string
}

func NewCache(size int) Cacher {
	return &lruCache{size: size, remaining: size, cache: make(map[string]string), queue: make([]string, size)}
}

func (lru *lruCache) Get(key interface{}) (interface{}, error) {
	// Your code here....

	//convert to string
	k := key.(string)

	//if there is no value, return error
	if lru.cache[k] == "" {
		return nil, errors.New("Not in cache")
	} else { //delete key from queue and move it to the tail
		lru.qDel(k)
		lru.queue = append(lru.queue, k)
		return lru.cache[k], nil
	}
}

func (lru *lruCache) Put(key, val interface{}) error {
	// Your code here....

	//convert to strings
	k := key.(string)
	v := val.(string)

	//error check, for empty key/val
	if len(k) == 0 || len(v) == 0 {
		return errors.New("Empty key/value")
	}

	//if cache is full and it is a new key/val
	//delete oldest key/value from cache/queue
	//and add new key to the tail of the queue
	if lru.remaining == 0 {
		delete(lru.cache, lru.queue[0])
		lru.qDel(lru.queue[0])
		lru.remaining++
		lru.queue = append(lru.queue, k)
	} else {
		//if cache is not full, shift keys in queue
		//to the left by storing all but the oldest
		//element and copying it to the queue
		//then give the last index the newest key
		temptail := lru.queue[1:]
		copy(lru.queue[0:], temptail)
		lru.queue[len(lru.queue)-1] = k
	}
	//if value is not in cache, put it in cache
	if lru.cache[k] != v {
		lru.cache[k] = v
		lru.remaining--
	}

	fmt.Println("Queue Length: " + strconv.Itoa(len(lru.queue)))
	fmt.Println(lru.queue)

	return nil
}

// Delete element from queue
func (lru *lruCache) qDel(ele string) {
	for i := 0; i < len(lru.queue); i++ {
		if lru.queue[i] == ele {
			oldlen := len(lru.queue)
			copy(lru.queue[i:], lru.queue[i+1:])
			lru.queue = lru.queue[:oldlen-1]
			break
		}
	}
}
