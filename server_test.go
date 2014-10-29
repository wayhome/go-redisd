package redis

import (
	"fmt"
	. "github.com/garyburd/redigo/redis"
	"testing"
	"time"
)

type MyHandler struct {
	values map[string][]byte
}

func (h *MyHandler) GET(key string) ([]byte, error) {
	v := h.values[key]
	return v, nil
}

func (h *MyHandler) SET(key string, value []byte, expire int) error {
	fmt.Printf("SET value: %v \n", value)
	fmt.Printf("expire in: %d \n", expire)
	h.values[key] = value
	return nil
}

func TestServer(t *testing.T) {
	srv, err := NewServer(":6389", &MyHandler{values: make(map[string][]byte)})
	if err != nil {
		panic(err)
	}
	go srv.ListenAndServe()
	time.Sleep(100 * time.Millisecond)
	c, _ := Dial("tcp", ":6389")
	defer c.Close()
	c.Do("SET", "hello", "world", 110)
	n, _ := c.Do("GET", "hello")
	fmt.Printf("GET value: %v \n", n)
}
