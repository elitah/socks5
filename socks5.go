package socks5

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"

	"github.com/armon/go-socks5"
)

var (
	EBadListener = fmt.Errorf("bad listener")
)

type cheatDNSResolver struct {
}

func (this *cheatDNSResolver) Resolve(ctx context.Context, name string) (context.Context, net.IP, error) {
	//
	return ctx, nil, nil
}

func (this *cheatDNSResolver) Rewrite(ctx context.Context, request *socks5.Request) (context.Context, *socks5.AddrSpec) {
	//
	return ctx, request.DestAddr
}

func dialWithContext(dial func(string, string) (net.Conn, error), logger *log.Logger) func(context.Context, string, string) (net.Conn, error) {
	//
	return func(ctx context.Context, network, address string) (net.Conn, error) {
		//
		var err error
		//
		ch0 := make(chan net.Conn)
		ch1 := make(chan error)
		//
		logger.Printf("dialWithContext: start connect to %s:%s", network, address)
		//
		go func(ch0 chan net.Conn, ch1 chan error, network, address string) {
			//
			if conn, err := dial(network, address); nil == err {
				//
				ch0 <- conn
			} else {
				//
				ch1 <- err
			}
			//
			close(ch0)
			close(ch1)
		}(ch0, ch1, network, address)
		//
		select {
		case <-ctx.Done():
			//
			err = context.DeadlineExceeded
		case conn, ok := <-ch0:
			//
			if ok {
				//
				logger.Printf("dialWithContext: success connect to %s:%s", network, address)
				//
				return conn, nil
			}
		case _err, ok := <-ch1:
			//
			if ok {
				//
				err = _err
			}
		}
		//
		logger.Printf("dialWithContext: failed connect to %s:%s, %v", network, address, err)
		//
		return nil, err
	}
}

func GetRawSocks5Server(args ...interface{}) error {
	//
	var l net.Listener
	//
	var logger *log.Logger
	//
	var handle0 func(context.Context, string, string) (net.Conn, error)
	//
	var handle1 func(string, string) (net.Conn, error)
	//
	for _, item := range args {
		//
		switch v := item.(type) {
		case net.Listener:
			//
			l = v
		case io.Writer:
			//
			logger = log.New(v, "", 0)
		case *log.Logger:
			//
			logger = v
		case func(context.Context, string, string) (net.Conn, error):
			//
			handle0 = v
		case func(string, string) (net.Conn, error):
			//
			handle1 = v
		}
	}
	//
	if nil != l {
		//
		if nil == logger {
			//
			logger = log.New(io.Discard, "", 0)
		}
		//
		if nil == handle0 && nil != handle1 {
			//
			handle0 = dialWithContext(handle1, logger)
		}
		//
		if nil == handle0 {
			//
			handle0 = (&net.Dialer{}).DialContext
		}
		//
		resolver := &cheatDNSResolver{}
		//
		if srv, err := socks5.New(&socks5.Config{
			Resolver: resolver,
			Rewriter: resolver,
			Logger:   logger,
			Dial:     handle0,
		}); nil == err {
			//
			return srv.Serve(l)
		} else {
			//
			return err
		}
	} else {
		//
		return EBadListener
	}
}
