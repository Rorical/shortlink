package main

import (
	"shortlink/internal/server"
)

func main() {
	ser := server.NewServer()
	ser.Init()
	ser.Run()
}
