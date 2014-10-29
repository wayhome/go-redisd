## go-redisd
A go library help you build server which speack the redis protocol. Forked from [go-redis-server](https://github.com/docker/go-redis-server)

## Usage

```go	
type MyHandler struct {
	values map[string][]byte
}

func (h *MyHandler) GET(key string) ([]byte, error) {
    v := h.values[key]
    return v, nil
    }

func (h *MyHandler) SET(key string, value []byte) error {
     h.values[key] = value
     return nil
}


srv, err := NewServer(":6389", &MyHandler{})
if err != nil {
	panic(err)
}
srv.ListenAndServe()
```
