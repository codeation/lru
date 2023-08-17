# lru

A lru is an asynchronous LRU cache (generic version).

[![GoDoc](https://godoc.org/github.com/codeation/lru?status.svg)](https://godoc.org/github.com/codeation/lru)

- based on golang map structure
- prevents infinite cache growth
- designed for asynchronous use

# Usage

For example, caching the output of the ioutil.ReadFile function to reduce disk I/O.

```
package main

import (
	"log"
	"os"

	"github.com/codeation/lru"
)

func readFileContent(key string) ([]byte, error) {
	log.Println("read once")
	return os.ReadFile(key)
}

func main() {
	cache := lru.NewCache(1024, readFileContent)
	for i := 0; i < 10; i++ {
		var data []byte
		data, err := cache.Get("input.txt")
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("file size is %d\n", len(data))
	}
}
```

The lru.NewCache parameter is the number of cache items until the last used item is removed from the cache. The second parameter is a func to get the value for the specified key and error.

The parameter of cache.Get func is a key (filename in this case).

An error is returned when the function returns an error.
