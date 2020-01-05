package main

import (
	"go-agar/internal/gateway"
)

func main() {
	g, e := gateway.NewGateway()
	if e != nil {
		println("gateway init error", e)
		return
	}
	g.Run()
}
