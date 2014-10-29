// go-redis-server is a helper library for building server software capable of speaking the redis protocol.
// This could be an alternate implementation of redis, a custom proxy to redis,
// or even a completely different backend capable of "masquerading" its API as a redis database.

package redis

import (
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"reflect"
)

type Server struct {
	Addr    string // TCP address to listen on, ":6389" if empty
	methods map[string]HandlerFn
}

func (srv *Server) ListenAndServe() error {
	addr := srv.Addr
	if addr == "" {
		addr = ":6389"
	}
	l, e := net.Listen("tcp", addr)
	if e != nil {
		return e
	}
	return srv.Serve(l)
}

// Serve accepts incoming connections on the Listener l, creating a
// new service goroutine for each.  The service goroutines read requests and
// then call srv.Handler to reply to them.
func (srv *Server) Serve(l net.Listener) error {
	defer l.Close()
	for {
		rw, err := l.Accept()
		if err != nil {
			return err
		}
		go srv.ServeClient(rw)
	}
}

// Serve starts a new redis session, using `conn` as a transport.
// It reads commands using the redis protocol, passes them to `handler`,
// and returns the result.
func (srv *Server) ServeClient(conn net.Conn) (err error) {
	defer func() {
		if err != nil {
			fmt.Fprintf(conn, "-%s\n", err)
		}
		conn.Close()
	}()

	clientChan := make(chan struct{})

	// Read on `conn` in order to detect client disconnect
	go func() {
		// Close chan in order to trigger eventual selects
		defer close(clientChan)
		defer Debugf("Client disconnected")
		// FIXME: move conn within the request.
		if false {
			io.Copy(ioutil.Discard, conn)
		}
	}()

	var clientAddr string

	switch co := conn.(type) {
	case *net.UnixConn:
		f, err := conn.(*net.UnixConn).File()
		if err != nil {
			return err
		}
		clientAddr = f.Name()
	default:
		clientAddr = co.RemoteAddr().String()
	}

	for {
		request, err := parseRequest(conn)
		if err != nil {
			return err
		}
		request.Host = clientAddr
		request.ClientChan = clientChan
		reply, err := srv.Apply(request)
		if err != nil {
			return err
		}
		if _, err = reply.WriteTo(conn); err != nil {
			return err
		}
	}
	return nil
}

func NewServer(addr string, handler interface{}) (*Server, error) {
	srv := &Server{}

	srv.Addr = addr

	rh := reflect.TypeOf(handler)
	for i := 0; i < rh.NumMethod(); i++ {
		method := rh.Method(i)
		if method.Name[0] > 'a' && method.Name[0] < 'z' {
			continue
		}
		handlerFn, err := srv.createHandlerFn(handler, &method.Func)
		if err != nil {
			return nil, err
		}
		srv.Register(method.Name, handlerFn)
	}
	return srv, nil
}
