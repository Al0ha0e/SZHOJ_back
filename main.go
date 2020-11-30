package main

import (
	"fmt"

	"github.com/Al0ha0e/SZHOJ_back/backserver"
)

func main() {
	server := backserver.GetBackServer()
	err := server.Init()
	if err != nil {
		fmt.Println(err)
		return
	}

	server.Run()
	defer server.Dispose()
}
