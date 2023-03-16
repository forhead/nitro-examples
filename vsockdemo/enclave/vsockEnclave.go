package main

import (
	"fmt"
	"log"
	// "net/http"
	// "encoding/json"
	"golang.org/x/sys/unix"
)


func main(){
	
	fmt.Println("Start nitro enclave vsock server...")

	fd, err := unix.Socket(unix.AF_VSOCK, unix.SOCK_STREAM, 0)
    if err != nil{
        log.Fatal(err)
    }
	fmt.Println("fd is: ", fd)

	// Bind socket to cid 16, port 5000.
	sockaddr := &unix.SockaddrVM{
		CID : unix.VMADDR_CID_ANY,
		Port : 5000,
	}

	err = unix.Bind(fd, sockaddr)
	if err != nil {
		log.Fatal("Bind ",err.Error())
	}
	// Listen for up to 32 incoming connections.
	err = unix.Listen(fd, 4)
	if err != nil {
		log.Fatal("Listen ",err)
	}

	for {
		peerFd, fromSockAdde, err := unix.Accept(fd)
		if err != nil {
			log.Fatal("Accept ",err)
		}
		fmt.Println("fromSockAdde: ",fromSockAdde)
		fmt.Println("peerFd is: ", peerFd)
		// nfd int, sa Sockaddr, err error
		// unix.Accept(fd)
		var data []byte
		unix.Recvfrom(peerFd, data, 0)
		fmt.Println(string(data))
	}
}