package structs

import (
	"crypto/sha256"
	"encoding/hex"
	"hash"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
)

type Cache struct {
	hasher hash.Hash
	values *sync.Map
}

func InitCache() *Cache {
	cache := Cache{
		hasher: sha256.New(),
		values: &sync.Map{},
	}
	return &cache
}

// Bootleg hash function to generate a hash based on
// Method + Auth + IP + Real IP + Path
func (cache *Cache) hash(request *http.Request) string {
	// Hasher needs to be reset as we write to
	cache.hasher.Reset()
	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		log.Print("Error while reading request body")
	}

	auth := request.Header.Get("Authorization")
	path := request.URL.Path
	method := request.Method
	realIp := request.Header.Get("X-Real-Ip")
	xForwardedFor := request.Header.Get("X-Forwarded-For")

	cache.hasher.Write([]byte(method + path + realIp + xForwardedFor + auth + string(body)))
	return hex.EncodeToString(cache.hasher.Sum(nil))
}

func (cache *Cache) check(request *http.Request) *http.Response {
	if val, ok := cache.values.Load(cache.hash(request)); ok {
		return val.(*http.Response)
	}

	return nil
}

func (cache *Cache) save(request *http.Request, response *http.Response) {
	cache.values.Store(cache.hash(request), response)
}

func (cache *Cache) removeInvalidEntires() {

}
