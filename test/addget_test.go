package test

import (
	"fmt"
	"testing"
	"time"

	"github.com/SzymonJaroslawski/pokedex-go/pokeapi"
)

func TestAddGet(t *testing.T) {
	const interval = 5 * time.Second
	cases := []struct {
		key string
		val []byte
	}{
		{
			key: "https://example.com",
			val: []byte("testdata"),
		},
		{
			key: "https://example.com/path",
			val: []byte("moretestdata"),
		},
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("Test case %v", i), func(t *testing.T) {
			cache := pokeapi.NewCache(interval)
			cache.Add(c.key, c.val)
			cache.Add("key2", []byte("val2"))
			val, ok := cache.Get(c.key)
			if !ok {
				t.Errorf("expected to find key")
				return
			}

			if string(val) != string(c.val) {
				t.Errorf("expected to find value")
				return
			}
		})
	}
}
