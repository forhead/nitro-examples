package main
import  (
        "fmt"
        "io"
        "log"
        "net/http"
        "encoding/json"
        "golang.org/x/sys/unix"

		// "github.com/aws/aws-sdk-go-v2/aws"
		// "github.com/aws/aws-sdk-go-v2/config"
		// "github.com/aws/aws-sdk-go-v2/service/dynamodb"
		// "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

    )

"""
dynamodb design
table name: demoWalletTable

colume:
walletName: wallet name for this wallet
encryptedPrivateKey: encrypted wallet private key
publicKey: the public key of the wallet
encryptedDatakey: the data key used to encrypt the private key
keyId: kms alias id which used for encryption for the private key
"""






// response
// generate encryptedDatakey
// 'encryptedPrivateKey': encrypted_privatekey,
//             'publicKey': public_key_hex,
//             'encryptedDatakey': encrypted_datakey

// sign
// content = sk.sign(bmessage)  # Signature in bytes

type parentClient struct{
	region string
	ddbTableName string
    keyId string
    cid uint32
    port uint32
}

type generateWalletResponse struct{
	encryptedPrivateKey string,
	publicKey string,
	encryptedDatakey string,
}

func (pc parentClient) generateWallet(walletName string){
	credential := getIAMToken()

	socket, err := unix.Socket(unix.AF_VSOCK, unix.SOCK_STREAM, 0)
    if err != nil{
        log.Fatal(err)
    }

    sockaddr := &unix.SockaddrVM{
        CID : pc.cid,
        Port : pc.port,
    }

    err = unix.Connect(socket, sockaddr) 
    if err != nil {
       	log.Fatal(err)
    }

    data, err := json.Marshal(&credential)
    unix.Write(socket,[]byte(data))

	playload := map[string]string{
		"apiCall":"generateWallet",
		"aws-access-key-id": credential.aws_access_key_id,
		"aws-secret-access-key" : credential.aws_secret_access_key,
		"aws-session-token" : credential.aws_session_token,
		"keyId":pc.keyId,
	}
	// Send AWS credential and KMS keyId to the server running in enclave
	b, err := json.Marshal(playload)
	unix.Write(socket,b)

	// receive data from the server and save to dynamodb with the walletName
	response := []byte{}
	unix.Read(socket,response)
	fmt.Println(string(response))

	// __saveEncryptWalletToDDB(walletName, response, self.__keyId)
	
}

func (pc parentClient) sign(keyId string, walletName string, message string){
	credential := getIAMToken()

	socket, err := unix.Socket(unix.AF_VSOCK, unix.SOCK_STREAM, 0)
    if err != nil{
        log.Fatal(err)
    }

    sockaddr := &unix.SockaddrVM{
        CID : pc.cid,
        Port : pc.port,
    }

    err = unix.Connect(socket, sockaddr) 
    if err != nil {
       	log.Fatal(err)
    }

    data, err := json.Marshal(&credential)
    unix.Write(socket,[]byte(data))

	var encryptoedPrivateKey = "" 
	var encryptedDatakey = ""

	playload := map[string]string{
		"apiCall":"generateWallet",
		"aws-access-key-id": credential.aws_access_key_id,
		"aws-secret-access-key" : credential.aws_secret_access_key,
		"aws-session-token" : credential.aws_session_token,
		"encryptedPrivateKey" : encryptoedPrivateKey,
		"encryptedDatakey" :encryptedDatakey,
		"keyId":pc.keyId,
		"message": message,
	}
	// Send AWS credential and KMS keyId to the server running in enclave
	b, err := json.Marshal(playload)
	unix.Write(socket,b)

	// receive data from the server and save to dynamodb with the walletName
	response := []byte{}
	unix.Read(socket,response)
	fmt.Println(string(response))
}

func (pc parentClient) saveEncryptWalletToDDB(walletName string, response []byte){
// 	cfg, err := config.LoadDefaultConfig(context.TODO(), func(o *config.LoadOptions) error {
//         o.Region = pc.region
//         return nil
//     })
//     if err != nil {
//         panic(err)
//     }
// // walletName: wallet name for this wallet
// // encryptedPrivateKey: encrypted wallet private key
// // publicKey: the public key of the wallet
// // encryptedDatakey: the data key used to encrypt the private key
// // keyId: kms alias id which used for encryption for the private key

// // 'encryptedPrivateKey': encrypted_privatekey,
// // 'publicKey': public_key_hex,
// // 'encryptedDatakey': encrypted_datakey
// 	var table ddbTable
// 	json.Unmarshal(body, &result)
// 	svc := dynamodb.NewFromConfig(cfg)
// 	out, err := svc.PutItem(context.TODO(), &dynamodb.PutItemInput{
//         TableName: pc.tableName,
//         Item: map[string]types.AttributeValue{
//             "walletName":    &types.AttributeValueMemberS{Value: "12346"},
//             "encryptedPrivateKey":  &types.AttributeValueMemberS{Value: "John Doe"},
//             "publicKey": &types.AttributeValueMemberS{Value: "john@doe.io"},
// 			"encryptedDatakey":,
// 			"keyId":,
//         },
//     })

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


func main(){
	client := parentClient{"ap-southeast-1", "demoWalletTable" ,"keyid",16,5000}
	client.generateWallet("wallet1")
}
