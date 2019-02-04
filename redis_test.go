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

	t.Run("delete non existing key", func(t *testing.T) {
		r, _ := NewRedis()

		err := r.Del("key")
		if err == nil {
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

	t.Run("overwrite existing key", func(t *testing.T) {
		r, _ := NewRedis()

		key := "key"
		old := "old"
		wanted := "value"

		err := r.Set(key, old)
		if err != nil {
			t.Errorf("could not set key and value: %v", err)
			return
		}

		err = r.Set(key, wanted)
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

	t.Run("add and get many keys", func(t *testing.T) {
		r, _ := NewRedis()
		items := []struct {
			key   string
			value string
		}{
			{"ka", "va"},
			{"kb", "vb"},
			{"kc", "vc"},
			{"kd", "vd"},
			{"ke", "ve"},
		}

		for _, item := range items {
			r.Set(item.key, item.value)
		}

		for _, item := range items {
			got, _ := r.Get(item.key)

			if item.value != got {
				t.Errorf("wanted: %s, got: %s", item.value, got)
			}
		}
	})

	t.Run("delete existing key", func(t *testing.T) {
		r, _ := NewRedis()

		key := "key"
		value := "value"

		err := r.Set(key, value)
		if err != nil {
			t.Errorf("could not set key and value: %v", err)
			return
		}

		err = r.Del(key)
		if err != nil {
			t.Errorf("should not return error: %v", err)
			return
		}

		if _, err := r.Get("key"); err == nil {
			t.Errorf("should return error")
		}
	})

	t.Run("increment non existing key", func(t *testing.T) {
		r, _ := NewRedis()

		key := "key"

		i, err := r.Incr(key)
		if err != nil {
			t.Errorf("could not incr key: %v", err)
			return
		}

		if i != 1 {
			t.Errorf("wanted: %d, got: %d", 1, i)
		}
	})

	t.Run("increment valid existing key", func(t *testing.T) {
		r, _ := NewRedis()

		key := "key"
		value := "5"

		err := r.Set(key, value)
		if err != nil {
			t.Errorf("could not set key and value: %v", err)
			return
		}

		i, err := r.Incr(key)
		if err != nil {
			t.Errorf("could not incr key: %v", err)
			return
		}

		if i != 6 {
			t.Errorf("wanted: %d, got: %d", 1, i)
		}
	})

	t.Run("increment invalid existing key", func(t *testing.T) {
		r, _ := NewRedis()

		key := "key"
		value := "asd"

		err := r.Set(key, value)
		if err != nil {
			t.Errorf("could not set key and value: %v", err)
			return
		}

		_, err = r.Incr(key)
		if err == nil {
			t.Errorf("should return error")
		}
	})

	// Todo: set expires
}
