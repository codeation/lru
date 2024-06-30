package lru

import (
	"testing"
)

func TestCacheHit(t *testing.T) {
	for _, tt := range []struct {
		title  string
		key    string
		cached bool
	}{{
		title:  "extra",
		key:    "F",
		cached: false,
	}, {
		title:  "latest",
		key:    "E",
		cached: true,
	}, {
		title:  "first",
		key:    "A",
		cached: true,
	}, {
		title:  "missed",
		key:    "B",
		cached: false,
	}} {
		t.Run(tt.title, func(t *testing.T) {
			funcKey := ""
			cache := NewSyncLRU(4, func(key string) (int, error) {
				funcKey = key
				return 42, nil
			})
			for _, key := range []string{"A", "B", "A", "C", "D", "E", "E"} { // expected cache keys are "A", "C", "D", "E"
				value, err := cache.Get(key)
				if err != nil {
					t.Fatalf("Error is: %v", err)
				}
				if value != 42 {
					t.Fatalf("Value is: %d", value)
				}
			}
			funcKey = ""

			value, err := cache.Get(tt.key)
			if err != nil {
				t.Fatalf("Error is: %v", err)
			}
			if value != 42 {
				t.Fatalf("Value is: %d", value)
			}
			if tt.cached != (funcKey != tt.key) {
				t.Fatalf("FuncKey is: %s", funcKey)
			}
		})
	}
}

func TestCacheOrder(t *testing.T) {
	for _, tt := range []struct {
		title    string
		keys     string
		expected string
	}{{
		title: "zero",
	}, {
		title:    "row",
		keys:     "ABCD",
		expected: "ABCD",
	}, {
		title:    "half",
		keys:     "ABA",
		expected: "BA",
	}, {
		title:    "first",
		keys:     "ABACADA",
		expected: "BCDA",
	}, {
		title:    "extra",
		keys:     "ABCDEFG",
		expected: "DEFG",
	}, {
		title:    "complex",
		keys:     "ABCDAEACCCFG",
		expected: "ACFG",
	}} {
		t.Run(tt.title, func(t *testing.T) {
			cache := NewSyncLRU(4, func(r rune) (int, error) { return 42, nil })
			for _, r := range tt.keys {
				_, _ = cache.Get(r)
			}

			output := ""
			for e := cache.queue.Back(); e != nil; e = e.Prev() {
				output += string(e.Value.(*oncePair[rune, int]).key)
			}

			if output != tt.expected {
				t.Fatalf("Cache keys unexpected: %s", output)
			}
		})
	}
}
