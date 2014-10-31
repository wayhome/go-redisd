package main

import (
	"fmt"
	. "github.com/youngking/go-redisd"
	"strconv"
)

type MyHandler struct {
	values map[string][]byte
}

func (h *MyHandler) GET(key string) ([]byte, error) {
	v := h.values[key]
	return v, nil
}

func (h *MyHandler) SET(key string, value []byte, expire string) error {
	ts, err := strconv.Atoi(expire)
	if err != nil {
		return err
	}
	fmt.Printf("expire in: %d \n", ts)
	fmt.Printf("SET value: %v \n", value)
	h.values[key] = value
	return nil
}

func main() {
	srv, err := NewServer("unix:///tmp/redis.sock", &MyHandler{values: make(map[string][]byte)})
	if err != nil {
		panic(err)
	}
	srv.ListenAndServe()
}
