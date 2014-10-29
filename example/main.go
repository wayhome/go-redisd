package main

import (
    "fmt"
    . "github.com/youngking/go-redisd"
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


func main() {
    srv, err := NewServer(":6389", &MyHandler{values: make(map[string][]byte)})
    if err != nil {
        panic(err)
    }
    srv.ListenAndServe()
}
