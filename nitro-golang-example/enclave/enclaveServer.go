package main

import (
	"fmt"
	"log"
	// "net/http"
	// "encoding/json"
	"golang.org/x/sys/unix"
)


func generateWallet(){
	
}

func main(){
	fmt.Println("Start nitro enclave vsock server...")

	fd, err := unix.Socket(unix.AF_VSOCK, unix.SOCK_STREAM, 0)
    if err != nil{
        log.Fatal(err)
    }

	// Bind socket to cid 16, port 5000.
	sockaddr := &unix.SockaddrVM{
		CID : 16,
		Port : 5000,
	}
	err = unix.Bind(fd, sockaddr)
	if err != nil {
		log.Fatal("Bind ",err)
	}
	// Listen for up to 32 incoming connections.
	err = unix.Listen(fd, 32)
	if err != nil {
		log.Fatal("Listen ",err)
	}
	for {
		var data []byte
		unix.Read(fd, data)
		fmt.Println(string(data))
	}
	
}