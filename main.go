package main

import (
	"fmt"

	"github.com/Al0ha0e/SZHOJ_back/backserver"
)

func main() {
	server := backserver.GetBackServer()
	err := server.Init()
	fmt.Println("OK2", err)
	server.Run()
	defer server.Dispose()
}
