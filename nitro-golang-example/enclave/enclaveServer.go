package main

import (
	"fmt"
	"log"
	"encoding/json"
	"golang.org/x/sys/unix"
)


func generateWallet() string{
	return "generateWallet"
}

func sign() string{
	return "sign"
}

type requestContext struct{
	// all request will contains
	apiCall string  //generateWallet and sign
	_aws_access_key_id string
    _aws_secret_access_key string
    _aws_session_token string
	// contains only in generateWallet
	keyId string
	// contains only in sign
	encryptedPrivateKey string
	encryptedDatakey string
	message string
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
		peerFd, fromSockAdde, err := unix.Accept(fd)
		if err != nil {
			log.Fatal("Accept ",err)
		}
		fmt.Println("fromSockAdde: ",fromSockAdde)
		fmt.Println("peerFd is: ", peerFd)

		var requestData []byte
		var rc requestContext

		unix.Recvfrom(peerFd, requestData, 0)
		json.Unmarshal(requestData, &rc)

		apiCall := rc.apiCall

        if apiCall == "generateWallet"{
            fmt.Println("generateWallet request")
			fmt.Println("rc: ",rc)
			_aws_access_key_id := rc._aws_access_key_id
			_aws_secret_access_key := rc._aws_secret_access_key
			_aws_session_token := rc._aws_session_token
            keyId := rc.keyId
            result := generateWallet()
            //  send back to parent instance
            // unix.Write(byte[]("generatewallet finished"))
			fmt.Println(_aws_access_key_id,_aws_secret_access_key,_aws_session_token,keyId,result)
            fmt.Println("generateWallet finished")
		} else if apiCall == "sign"{
			fmt.Println("sign request")
            message := rc.message
            encryptedPrivateKey := rc.encryptedPrivateKey
            encryptedDatakey := rc.encryptedDatakey
			fmt.Println(encryptedPrivateKey,encryptedDatakey,message)
            // signedStr = server.sign(
            //     credential, encryptedPrivateKey, encryptedDatakey, message)
            // c.send(signedStr)

			fmt.Println("sign fihished")

		}else {
			fmt.Println("nothing to do")
		}
	}
	
}