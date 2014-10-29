package redis

import (
    "strings"
)


type HandlerFn func(r *Request) (ReplyWriter, error)


func (srv *Server) Apply(r *Request) (ReplyWriter, error) {
	if srv == nil || srv.methods == nil {
		Debugf("The method map is uninitialized")
		return ErrMethodNotSupported, nil
	}
	fn, exists := srv.methods[strings.ToLower(r.Name)]
	if !exists {
		return ErrMethodNotSupported, nil
	}
	return fn(r)
}


func (srv *Server) Register(name string, fn HandlerFn) {
	if srv.methods == nil {
		srv.methods = make(map[string]HandlerFn)
	}
	if fn != nil {
		Debugf("REGISTER: %s", strings.ToLower(name))
		srv.methods[strings.ToLower(name)] = fn
	}
}