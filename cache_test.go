package ccache

import (
	"github.com/karlseguin/gspec"
	"testing"
	"strconv"
	"time"
)

func TestGCsTheOldestItems(t *testing.T) {
	spec := gspec.New(t)
	cache := New(Configure().ItemsToPrune(10))
	for i := 0; i < 500; i++ {
		cache.Set(strconv.Itoa(i), i, time.Minute)
	}
	cache.gc()
	spec.Expect(cache.Get("9")).ToBeNil()
	spec.Expect(cache.Get("10").(int)).ToEqual(10)
}

func TestPromotedItemsDontGetPruned(t *testing.T) {
	spec := gspec.New(t)
	cache := New(Configure().ItemsToPrune(10).GetsPerPromote(1))
	for i := 0; i < 500; i++ {
		cache.Set(strconv.Itoa(i), i, time.Minute)
	}
	cache.Get("9")
	time.Sleep(time.Millisecond * 10)
	cache.gc()
	spec.Expect(cache.Get("9").(int)).ToEqual(9)
	spec.Expect(cache.Get("10")).ToBeNil()
	spec.Expect(cache.Get("11").(int)).ToEqual(11)
}

func TestTrackerDoesNotCleanupHeldInstance(t *testing.T) {
	spec := gspec.New(t)
	cache := New(Configure().ItemsToPrune(10).Track())
	for i := 0; i < 10; i++ {
		cache.Set(strconv.Itoa(i), i, time.Minute)
	}
	item := cache.TrackingGet("0")
	time.Sleep(time.Millisecond * 10)
	cache.gc()
	spec.Expect(cache.Get("0").(int)).ToEqual(0)
	spec.Expect(cache.Get("1")).ToBeNil()
	item.Release()
	cache.gc()
	spec.Expect(cache.Get("0")).ToBeNil()
}