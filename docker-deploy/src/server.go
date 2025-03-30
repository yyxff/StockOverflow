package main

import (
	"log"

	"github.com/panjf2000/gnet/v2"
)

type Server struct {
	gnet.BuiltinEventEngine
}

func (es *Server) OnBoot(eng gnet.Engine) gnet.Action {
	log.Printf("echo server with multi-core= is listening on\n")
	return gnet.None
}

func (es *Server) OnTraffic(c gnet.Conn) gnet.Action {
	buf, _ := c.Next(-1)
	c.Write(buf)
	return gnet.None
}
