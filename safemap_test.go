package main

import (
	"testing"
)

func TestBasic(t *testing.T) {
	t.Run("get non existing key", func(t *testing.T) {
		m, _ := NewSafeMap()
		key := Key("key")

		_, err := m.Get(key)

		if err == nil {
			t.Errorf("should have returned an error")
		}
	})

	t.Run("delete non existing key", func(t *testing.T) {
		m, _ := NewSafeMap()
		key := Key("key")

		err := m.Del(key)

		if err == nil {
			t.Errorf("should have returned an error")
		}
	})

	t.Run("add and get key", func(t *testing.T) {
		m, _ := NewSafeMap()
		key := Key("key")
		wanted := Value("value")

		m.Add(key, wanted)
		got, _ := m.Get(key)

		if wanted != got {
			t.Errorf("wanted: %p, got: %p", wanted, got)
		}
	})

	t.Run("overwrite existing key", func(t *testing.T) {
		m, _ := NewSafeMap()
		key := Key("key")
		v := Value("old")
		wanted := Value("new")

		m.Add(key, v)
		m.Add(key, wanted)
		got, _ := m.Get(key)

		if wanted != got {
			t.Errorf("wanted: %p, got: %p", wanted, got)
		}
	})

	t.Run("add and get many keys", func(t *testing.T) {
		m, _ := NewSafeMap()
		items := []struct {
			key   Key
			value Value
		}{
			{Key("ka"), Value("va")},
			{Key("kb"), Value("vb")},
			{Key("kc"), Value("vc")},
			{Key("kd"), Value("vd")},
			{Key("ke"), Value("ve")},
		}

		for _, item := range items {
			m.Add(item.key, item.value)
		}

		for _, item := range items {
			got, _ := m.Get(item.key)

			if item.value != got {
				t.Errorf("wanted: %p, got: %p", item.value, got)
			}
		}
	})

	t.Run("delete existing key", func(t *testing.T) {
		m, _ := NewSafeMap()
		key := Key("key")
		wanted := Value("value")

		m.Add(key, wanted)
		err := m.Del(key)

		if err != nil {
			t.Errorf("should not return error: %v", err)
		}
	})

	t.Run("delete existing key if value is the same as stored", func(t *testing.T) {
		m, _ := NewSafeMap()
		key := Key("key")
		v1 := Value("value")
		v2 := Value("value")

		m.Add(key, v1)
		m.Add(key, v2)
		err := m.DelIf(key, v2)

		if err != nil {
			t.Errorf("should not return error: %v", err)
		}
	})

	t.Run("don't delete existing key if value is not the same as stored", func(t *testing.T) {
		m, _ := NewSafeMap()
		key := Key("key")
		v1 := Value("value")
		v2 := Value("value")

		m.Add(key, v1)
		m.Add(key, v2)
		err := m.DelIf(key, v1)

		if err == nil {
			t.Errorf("should return error")
		}
	})
}
