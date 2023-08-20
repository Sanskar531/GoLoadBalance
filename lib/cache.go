package lib

import (
	"crypto/sha256"
	"encoding/hex"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"
)

type Cache struct {
	values *sync.Map
}

func InitCache() *Cache {
	cache := Cache{
		values: &sync.Map{},
	}

	return &cache
}

// Basic Hash function to generate a hash based on
// Method + Auth + IP + Real IP + Path
func (cache *Cache) hash(request *http.Request) string {
	hasher := sha256.New()
	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		log.Print("Error while reading request body")
	}

	auth := request.Header.Get("Authorization")
	path := request.URL.Path
	method := request.Method
	realIp := request.Header.Get("X-Real-Ip")
	xForwardedFor := request.Header.Get("X-Forwarded-For")

	hasher.Write([]byte(method + path + realIp + xForwardedFor + auth + string(body)))
	return hex.EncodeToString(hasher.Sum(nil))
}

func (cache *Cache) check(request *http.Request) *map[string]any {
	if val, ok := cache.values.Load(cache.hash(request)); ok {
		return val.(*map[string]any)
	}

	return nil
}

func (cache *Cache) save(request *http.Request, body *string, response *http.Response, cacheTimeoutInSeconds int) {
	responseAndBody := make(map[string]any)
	responseAndBody["body"] = body
	responseAndBody["response"] = response
	hashedKey := cache.hash(request)
	cache.values.Store(hashedKey, &responseAndBody)
	go cache.keepHashedKeyAliveFor(hashedKey, time.Second*time.Duration(cacheTimeoutInSeconds))
}

func (cache *Cache) keepHashedKeyAliveFor(hashKey string, duration time.Duration) {
	time.Sleep(duration)
	cache.values.Delete(hashKey)
}
