package main

import (
	"fmt"
	"sync"
)

type KV struct {
	mu    *sync.RWMutex
	store map[string][]byte
}

func NewKV() *KV {
	return &KV{
		mu:    &sync.RWMutex{},
		store: map[string][]byte{},
	}
}

func (kv *KV) GET(key []byte) ([]byte, error) {
	kv.mu.RLock()
	// この変換、無駄が多そう
	value, ok := kv.store[string(key)]
	kv.mu.RUnlock()

	// この型の代入、無駄多そう
	value = append([]byte("+"), value...)
	if ok {
		return value, nil
	} else {
		return value, fmt.Errorf("this key does not exist: key=%s", string(key))
	}
}

func (kv *KV) SET(key []byte, value []byte) (ok bool, err error) {
	kv.mu.Lock()
	defer kv.mu.Unlock()
	// この変換、無駄が多そう
	kv.store[string(key)] = value

	return true, nil
}
