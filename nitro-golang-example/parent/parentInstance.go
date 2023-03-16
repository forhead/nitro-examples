package main
import  (
        "fmt"
        "io"
        "log"
        "net/http"
        "encoding/json"
        "golang.org/x/sys/unix"
    )
	
type parentInstance struct{
}

func (pi parentInstance) generateWallet(kmsArn string){

}

func (pi parentInstance) sign(encryptWallet string){

}


type iamCredential struct {
	aws_access_key_id string
	aws_secret_access_key string
	aws_session_token string
}

// struct of response from metadata get function
type iamCredentialToken struct{
	Code string
	LastUpdated string
	Type string
	AccessKeyId string
	SecretAccessKey string
	Token string
	Expiration string
}
/**
* get the credential of the IAM Role attached on EC2
*/
func getIAMToken() iamCredential{
	var token iamCredential
	res, err := http.Get("http://169.254.169.254/latest/meta-data/iam/security-credentials/")

	if err != nil{
			log.Fatal(err)
	}

	body, err := io.ReadAll(res.Body)
	res.Body.Close()
	instanceProfileName := string(body)
	// get the credential of iam role
	profileUri :=fmt.Sprintf("http://169.254.169.254/latest/meta-data/iam/security-credentials/%s",instanceProfileName)
	res, err = http.Get(profileUri)

	if err != nil{
			log.Fatal(err)
	}

	body, err = io.ReadAll(res.Body)
	res.Body.Close()

	var result iamCredentialToken
	json.Unmarshal(body, &result)

	token.aws_access_key_id = result.AccessKeyId
	token.aws_secret_access_key = result.SecretAccessKey
	token.aws_session_token = result.Token

	return token
}

// http handler
type enclaveHandler struct{}

func (h *enclaveHandler) ServeHTTP(w http.ResponseWriter, r *http.Request){
    // TODO implement a function which get kmskey from request
    send_token_to_enclave(16,5000,"kmskey")

}

func send_token_to_enclave(cid uint32, port uint32, kms_key_id string){
    session_token := getIAMToken()

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
       	log.Fatal(err)
    }
    data, err := json.Marshal(&session_token)
    unix.Write(socket,[]byte(data))

}
