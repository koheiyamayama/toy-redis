package main

import (
	"errors"
	"sync"
)

type KV struct {
	mu    *sync.RWMutex
	store map[string][]byte
}

var ErrNotExistsKey = errors.New("not found key")

func NewKV() *KV {
	return &KV{
		mu:    &sync.RWMutex{},
		store: map[string][]byte{},
	}
}

func (kv *KV) Get(key []byte) ([]byte, error) {
	kv.mu.RLock()
	// この変換、無駄が多そう
	value, ok := kv.store[string(key)]
	kv.mu.RUnlock()

	// この型の代入、無駄多そう
	if ok {
		value = append([]byte("+"), value...)
		return value, nil
	} else {
		return value, ErrNotExistsKey
	}
}

func (kv *KV) Set(key []byte, value []byte) (ok bool, err error) {
	kv.mu.Lock()
	defer kv.mu.Unlock()
	// この変換、無駄が多そう
	kv.store[string(key)] = value

	return true, nil
}
