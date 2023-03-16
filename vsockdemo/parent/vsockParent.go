package main
import  (
        // "fmt"
        // "io"
        "log"
        // "net/http"
        "encoding/json"
        "golang.org/x/sys/unix"
    )


var KMS_KEY_id string = "KMS_KEY_ID"
var ENCRYPT_DATA string = "ENCRYPT_DATA" 
// http handler
// type enclave_handler struct{}

// func (h *enclave_handler) ServeHTTP(w http.ResponseWriter, r *http.Request){
//     // TODO implement a function which get kmskey from request
//     // send_token_to_enclave(16,5000,"kmskey")

// }

func sendEncryptDataToEnclave(cid uint32, port uint32, kmsKeyId string, encryptData string){

    socket, err := unix.Socket(unix.AF_VSOCK, unix.SOCK_STREAM, 0)
    if err != nil{
        log.Fatal(err)
    }

    sockaddr := &unix.SockaddrVM{
        CID : cid,
        Port : port,
    }

    err = unix.Connect(socket, sockaddr) 
    if err != nil {
       	log.Fatal("connect ",err)
    }
    
	dataMap := make(map[string]interface{})
	dataMap[KMS_KEY_id] = kmsKeyId
	dataMap[ENCRYPT_DATA] = encryptData
	dataMapBytes, err := json.Marshal(dataMap)
	if err != nil {
		log.Fatal(err.Error())
	}
    unix.Write(socket,dataMapBytes)
}

func main(){
    //send data to enclave via vsock
	sendEncryptDataToEnclave(16,5000,"kmsKeyId_id","encryptData_data")
	
	// start to listen
	// http.Handle("/", &enclave_handler{})

    // log.Println("Listening...")
    // err := http.ListenAndServe(":443", nil)
    // if err != nil {
    //     log.Fatal(err)
    // }
}
