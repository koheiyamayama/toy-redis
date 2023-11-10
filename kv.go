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
	mu          *sync.RWMutex
	doneExpChan chan struct{}
	store       map[string][]byte
}

var ErrNotExistsKey = errors.New("not found key")

func NewKV() *KV {
	return &KV{
		mu:          &sync.RWMutex{},
		store:       map[string][]byte{},
		doneExpChan: make(chan struct{}),
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
	var exists bool
	kv.mu.RLock()
	_, exists = kv.store[string(key)]
	kv.mu.RUnlock()

	if exists {
		slog.Debug("exists key")
		kv.Expire(key, exp)
	} else {
		slog.Debug("not exists key")
		kv.mu.Lock()
		// この変換、無駄が多そう
		kv.store[string(key)] = value
		kv.mu.Unlock()
		kv.expire(key, exp)
	}

	return true, nil
}

func (kv *KV) Expire(key []byte, exp uint32) (ok bool, err error) {
	slog.Debug("before writing doneExpChan")
	kv.doneExpChan <- struct{}{}
	slog.Debug("after writing doneExpChan")
	return kv.expire(key, exp)
}

func (kv *KV) expire(key []byte, exp uint32) (ok bool, err error) {
	logger.DebugCtx(context.Background(), "expireEntry",
		slog.String("key", string(key)),
	)

	go func(k []byte, exp uint32) {
		for {
			select {
			case <-time.After(time.Duration(exp) * time.Second):
				slog.Debug("delete entry by kv#expireEntry")
				kv.Del(key)
				return
			case <-kv.doneExpChan:
				slog.Debug("kv.doneExpChan")
				return
			}
		}
	}(key, exp)

	return true, nil
}

func (kv *KV) Del(key []byte) (ok bool, err error) {
	kv.mu.Lock()
	delete(kv.store, string(key))
	kv.mu.Unlock()

	return true, nil
}
