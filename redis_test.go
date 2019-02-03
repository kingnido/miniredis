package main

import (
	"testing"
)

func TestRedis(t *testing.T) {
	t.Run("get non existing key", func(t *testing.T) {
		r, _ := NewRedis()

		if _, err := r.Get("key"); err == nil {
			t.Errorf("should return error")
		}
	})

	t.Run("add and get key", func(t *testing.T) {
		r, _ := NewRedis()

		key := "key"
		wanted := "value"

		err := r.Set(key, wanted)
		if err != nil {
			t.Errorf("could not set key and value: %v", err)
			return
		}

		got, err := r.Get(key)
		if err != nil {
			t.Errorf("could not get key and value: %v", err)
			return
		}

		if wanted != got {
			t.Errorf("wanted: %s, got: %s", wanted, got)
			return
		}
	})
}
