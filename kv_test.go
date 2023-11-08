package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetAndSet(t *testing.T) {
	kv := NewKV()
	// TODO: テストケース修正
	// - expireEntryのテストする
	ok, err := kv.Set([]byte("key"), []byte("value"), 0)
	assert.Equal(t, true, ok)
	assert.Nil(t, err)

	v, err := kv.Get([]byte("key"))
	assert.Equal(t, []byte("+value"), v)
	assert.Nil(t, err)

	v, err = kv.Get([]byte("not_exists_key"))
	assert.Equal(t, []byte(nil), v)
	assert.EqualError(t, err, ErrNotExistsKey.Error())

}

func TestDel(t *testing.T) {
	t.SkipNow()
}
