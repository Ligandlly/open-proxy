package main

import (
	"fmt"
	"os"
)

func main() {
	var ip string
	var port = "7890"
	if len(os.Args) >= 2 {
		ip = os.Args[1]
	}
	if len(os.Args) == 3 {
		port = os.Args[2]
	}
	if ip == "" {
		fmt.Printf("The ip address is required")
		os.Exit(1)
	}

	fmt.Printf("export https_proxy=http://%s:%s;export http_proxy=http://%s:%s;export all_proxy=socks5://%s:%s\n", ip, port, ip, port, ip, port)
}
