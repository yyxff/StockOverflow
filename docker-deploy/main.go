package main

import (
	"log"

	"github.com/panjf2000/gnet/v2"
)

func main() {

	echo := &Server{}
	log.Fatal(gnet.Run(echo, "tcp://localhost:12345", gnet.WithMulticore(true)))
}
