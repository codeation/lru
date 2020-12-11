# lru
A lru is an asynchronous LRU cache. 

[![GoDoc](https://godoc.org/github.com/codeation/lru?status.svg)](https://godoc.org/github.com/codeation/lru)

- based on golang map structure
- prevents infinite cache growth
- designed for asynchronous use

# Usage

For example, caching the output of the ioutil.ReadFile function to reduce disk I/O.

```
package main

import (
	"io/ioutil"
	"log"

	"github.com/codeation/lru"
)

func readFileContent(key string) (interface{}, error) {
	log.Printf("read once\n")
	return ioutil.ReadFile(key)
}

func main() {
	cache := lru.NewCache(1024)
	for i := 0; i < 10; i++ {
		var data []byte
		if err := cache.Get("input.txt", readFileContent, &data); err != nil {
			log.Fatal(err)
		}
		log.Printf("file size is %d\n", len(data))
	}
}
```

First parameter of cache.Get func is a key (filename in this case). The second parameter is a func to get the value for the specified key and error. The func type declaration must be the same as readFileContent in the example. The third parameter is a pointer to a variable to set the value. 

An error is returned when the function returns an error. An error is also returned if the return value of the function cannot be assigned to the variable.

The lru.NewCache parameter is the number of cache items until the last used item is removed from the cache.
