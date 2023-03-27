package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"

	"golang.org/x/sys/unix"
)

// dynamodb design
// table name: AccountTable

// colume:
// keyId: kms alias id which used for encryption for the private key
// Name: account name for this Account
// encryptedPrivateKey: encrypted Account private key
// address: the address of the Account
// encryptedDataKey: the data key used to encrypt the private key

type accountTable struct {
	keyId               string
	name                string
	address             string
	encryptedDataKey    string
	encryptedPrivateKey string
}

type accountClient struct {
	region       string
	ddbTableName string
	keyId        string
	cid          uint32
	port         uint32
}

type generateAccountResponse struct {
	EncryptedPrivateKey string
	Address             string
	EncryptedDataKey    string
}

type requestPlayload struct {
	ApiCall               string
	Aws_access_key_id     string
	Aws_secret_access_key string
	Aws_session_token     string
	KeyId                 string // this is for generateAccount
	//this 3 is for sign
	EncryptedPrivateKey string
	EncryptedDataKey    string
	Transaction         string
}

func (ac accountClient) saveEncryptWalletToDDB(name string, response generateAccountResponse, keyId string) {
	// Create Session
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(ac.region)},
	)

	if err != nil {
		panic(err)
	}

	svc := dynamodb.New(sess)

	at := accountTable{
		name:                name,
		keyId:               keyId,
		address:             response.Address,
		encryptedPrivateKey: response.EncryptedPrivateKey,
		encryptedDataKey:    response.EncryptedDataKey,
	}

	av, err := dynamodbattribute.MarshalMap(at)

	if err != nil {
		fmt.Println("Got error marshalling map:")
		fmt.Println(err.Error())
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(ac.ddbTableName),
	}

	_, err = svc.PutItem(input)
	if err != nil {
		fmt.Println("Got error calling PutItem:")
		fmt.Println(err.Error())
	}
}

func (ac accountClient) generateAccount(name string) {
	credential := getIAMToken()

	socket, err := unix.Socket(unix.AF_VSOCK, unix.SOCK_STREAM, 0)
	if err != nil {
		log.Fatal(err)
	}

	sockaddr := &unix.SockaddrVM{
		CID:  ac.cid,
		Port: ac.port,
	}

	err = unix.Connect(socket, sockaddr)
	if err != nil {
		log.Fatal(err)
	}

	playload := requestPlayload{
		ApiCall:               "generateAccount",
		Aws_access_key_id:     credential.aws_access_key_id,
		Aws_secret_access_key: credential.aws_secret_access_key,
		Aws_session_token:     credential.aws_session_token,
		KeyId:                 ac.keyId,
		EncryptedPrivateKey:   "",
		EncryptedDataKey:      "",
		Transaction:           "",
	}

	// Send AWS credential and KMS keyId to the server running in enclave
	b, err := json.Marshal(playload)
	if err != nil {
		fmt.Println(err)
	}
	unix.Write(socket, b)

	// receive data from the server and save to dynamodb with the walletName
	response := []byte{}
	unix.Read(socket, response)

	var responseStruct generateAccountResponse
	json.Unmarshal(response, &responseStruct)
	fmt.Println(string(response))

	// ac.saveEncryptWalletToDDB(name, responseStruct, ac.keyId)

}

func (ac accountClient) sign(keyId string, name string, transaction string) {
	credential := getIAMToken()

	// get item from dynamodb
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(ac.region)},
	)

	if err != nil {
		panic(err)
	}

	svc := dynamodb.New(sess)

	result, err := svc.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(ac.ddbTableName),
		Key: map[string]*dynamodb.AttributeValue{
			"keyId": {
				N: aws.String(keyId),
			},
			"name": {
				S: aws.String(name),
			},
		},
	})

	var at accountTable
	err = dynamodbattribute.UnmarshalMap(result.Item, &at)
	if err != nil {
		panic(fmt.Sprintf("Failed to unmarshal Record, %v", err))
	}
	var encryptedDataKey = at.encryptedDataKey
	var encryptedPrivateKey = at.encryptedPrivateKey

	socket, err := unix.Socket(unix.AF_VSOCK, unix.SOCK_STREAM, 0)
	if err != nil {
		log.Fatal(err)
	}

	sockaddr := &unix.SockaddrVM{
		CID:  ac.cid,
		Port: ac.port,
	}

	err = unix.Connect(socket, sockaddr)
	if err != nil {
		log.Fatal(err)
	}

	playload := map[string]string{
		"apiCall":               "generateWallet",
		"aws-access-key-id":     credential.aws_access_key_id,
		"aws-secret-access-key": credential.aws_secret_access_key,
		"aws-session-token":     credential.aws_session_token,
		"encryptedPrivateKey":   encryptedPrivateKey,
		"encryptedDataKey":      encryptedDataKey,
		"keyId":                 keyId,
		"transaction":           transaction,
	}
	fmt.Println(playload)
	// Send AWS credential and KMS keyId to the server running in enclave
	b, err := json.Marshal(playload)

	fmt.Println(b)

	unix.Write(socket, b)

	// receive data from the server and save to dynamodb with the walletName
	response := []byte{}
	unix.Read(socket, response)
	fmt.Println(string(response))
}

type iamCredentialResponse struct {
	aws_access_key_id     string
	aws_secret_access_key string
	aws_session_token     string
}

// struct of response from metadata get function
type iamCredentialToken struct {
	Code            string
	LastUpdated     string
	Type            string
	AccessKeyId     string
	SecretAccessKey string
	Token           string
	Expiration      string
}

/**
* get the credential of the IAM Role attached on EC2
 */
func getIAMToken() iamCredentialResponse {
	var token iamCredentialResponse
	res, err := http.Get("http://169.254.169.254/latest/meta-data/iam/security-credentials/")
	if err != nil {
		log.Fatal(err)
	}
	body, err := io.ReadAll(res.Body)
	res.Body.Close()
	instanceProfileName := string(body)
	profileUri := fmt.Sprintf("http://169.254.169.254/latest/meta-data/iam/security-credentials/%s", instanceProfileName)
	res, err = http.Get(profileUri)
	if err != nil {
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

func main() {
	client := accountClient{"ap-northeast-1", "AccountTable", "0f360b0f-1ad4-4c6b-b405-932d2f606779", 16, 5000}
	client.generateAccount("wallet1")
}
