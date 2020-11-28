package main

import (
	"fmt"

	"github.com/Al0ha0e/SZHOJ_back/backserver"
)

func main() {
	server := backserver.GetBackServer()
	fmt.Println("OK1", server == nil)
	err := server.Init()
	fmt.Println("OK2", err)
	server.Run()
}
