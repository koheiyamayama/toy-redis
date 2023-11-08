package main

import (
	"context"
	"errors"
	"log/slog"
	"sync"
	"time"

	"github.com/koheiyamayama/toy-redis/logger"
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

func (kv *KV) Set(key []byte, value []byte, exp uint32) (ok bool, err error) {
	kv.mu.Lock()
	// この変換、無駄が多そう
	kv.store[string(key)] = value
	kv.mu.Unlock()

	go kv.expireEntry(key, exp)
	return true, nil
}

func (kv *KV) expireEntry(key []byte, exp uint32) (ok bool, err error) {
	logger.DebugCtx(context.Background(), "expireEntry",
		slog.String("key", string(key)),
		slog.Uint64("exp", uint64(exp)),
	)

	<-time.After(time.Second * time.Duration(exp))
	slog.Debug("delete entry by kv#expireEntry")
	return kv.Del(key)
}

func (kv *KV) Del(key []byte) (ok bool, err error) {
	kv.mu.Lock()
	delete(kv.store, string(key))
	kv.mu.Unlock()

	return true, nil
}
