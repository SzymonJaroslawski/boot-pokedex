package test

import (
	"fmt"
	"testing"
	"time"

	"github.com/SzymonJaroslawski/pokedex-go/pokeapi"
)

func TestReap(t *testing.T) {
	interval := time.Millisecond * 10
	cache := pokeapi.NewCache(interval)

	keyOne := "key1"
	cache.Add(keyOne, []byte("val1"))

	time.Sleep(interval / 2)

	_, ok := cache.Get(keyOne)
	if !ok {
		vals := cache.GetAll()
		for _, val := range vals {
			fmt.Printf("%s\n", string(val))
		}
		t.Errorf("%s should not have been reapd", keyOne)
	}
}
