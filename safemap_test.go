package main

import (
	"testing"
)

func TestSafeMap(t *testing.T) {
	t.Run("get non existing key", func(t *testing.T) {
		m := NewSafeMap()
		key := "key"

		_, err := m.Get(key)

		if err == nil {
			t.Errorf("should have returned an error")
		}
	})

	t.Run("delete non existing key", func(t *testing.T) {
		m := NewSafeMap()
		key := "key"

		err := m.Del(key)

		if err == nil {
			t.Errorf("should have returned an error")
		}
	})

	t.Run("add and get key", func(t *testing.T) {
		m := NewSafeMap()
		key := "key"
		wanted := "value"

		m.Set(key, wanted)
		got, _ := m.Get(key)

		if wanted != got {
			t.Errorf("wanted: %s, got: %s", wanted, got)
		}
	})

	t.Run("overwrite existing key", func(t *testing.T) {
		m := NewSafeMap()
		key := "key"
		v := "old"
		wanted := "new"

		m.Set(key, v)
		m.Set(key, wanted)
		got, _ := m.Get(key)

		if wanted != got {
			t.Errorf("wanted: %s, got: %s", wanted, got)
		}
	})

	t.Run("add and get many keys", func(t *testing.T) {
		m := NewSafeMap()
		items := []struct {
			key   string
			value interface{}
		}{
			{"ka", "va"},
			{"kb", "vb"},
			{"kc", "vc"},
			{"kd", "vd"},
			{"ke", "ve"},
		}

		for _, item := range items {
			m.Set(item.key, item.value)
		}

		for _, item := range items {
			got, _ := m.Get(item.key)

			if item.value != got {
				t.Errorf("wanted: %p, got: %p", item.value, got)
			}
		}
	})

	t.Run("delete existing key", func(t *testing.T) {
		m := NewSafeMap()
		key := "key"
		wanted := "value"

		m.Set(key, wanted)
		err := m.Del(key)

		if err != nil {
			t.Errorf("should not return error: %v", err)
		}
	})

	t.Run("delete existing key if value is the same as stored", func(t *testing.T) {
		m := NewSafeMap()
		key := "key"
		v1 := "value"
		v2 := "value"

		m.Set(key, &v1)
		m.Set(key, &v2)
		err := m.DelIf(key, &v2)

		if err != nil {
			t.Errorf("should not return error: %v", err)
		}
	})

	t.Run("don't delete existing key if value is not the same as stored", func(t *testing.T) {
		m := NewSafeMap()
		key := "key"
		v1 := "value"
		v2 := "value"

		m.Set(key, &v1)
		m.Set(key, &v2)
		err := m.DelIf(key, &v1)

		if err == nil {
			t.Errorf("should return error")
		}
	})
}
