package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/ethereum/go-ethereum/common/hexutil"

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
	KeyId               string
	Name                string
	Address             string
	EncryptedDataKey    string
	EncryptedPrivateKey string
}

type signedValueTable struct{
	Name string
	Transaction string
	SignedValue string
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
	response := make([]byte, 4096)
	n, err := unix.Read(socket, response)
	if err != nil {
		fmt.Println(err)
	}
	var responseStruct generateAccountResponse
	json.Unmarshal(response[:n], &responseStruct)

	ac.saveEncryptAccountToDDB(name, responseStruct, ac.keyId)

}

func (ac accountClient) saveEncryptAccountToDDB(name string, response generateAccountResponse, keyId string) {
	// Create Session
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(ac.region)},
	)

	if err != nil {
		panic(err)
	}

	svc := dynamodb.New(sess)

	at := accountTable{
		Name:                name,
		KeyId:               keyId,
		Address:             response.Address,
		EncryptedPrivateKey: response.EncryptedPrivateKey,
		EncryptedDataKey:    response.EncryptedDataKey,
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
	fmt.Println("account", name, "info is saved to dynamodb")
}

func (ac accountClient) sign(keyId string, name string, transaction string) string {
	credential := getIAMToken()
	// get item from dynamodb
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(ac.region)},
	)

	if err != nil {
		panic(err)
	}

	svc := dynamodb.New(sess)

	result, _ := svc.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(ac.ddbTableName),
		Key: map[string]*dynamodb.AttributeValue{
			"KeyId": {
				S: aws.String(keyId),
			},
			"Name": {
				S: aws.String(name),
			},
		},
	})

	if err != nil {
		fmt.Println("ddb query err:", err)
	}

	var at accountTable
	err = dynamodbattribute.UnmarshalMap(result.Item, &at)
	if err != nil {
		panic(fmt.Sprintf("Failed to unmarshal Record, %v", err))
	}
	var encryptedDataKey = at.EncryptedDataKey
	var encryptedPrivateKey = at.EncryptedPrivateKey

	fmt.Println(at)

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
		ApiCall:               "sign",
		Aws_access_key_id:     credential.aws_access_key_id,
		Aws_secret_access_key: credential.aws_secret_access_key,
		Aws_session_token:     credential.aws_session_token,
		KeyId:                 "",
		EncryptedPrivateKey:   encryptedPrivateKey,
		EncryptedDataKey:      encryptedDataKey,
		Transaction:           transaction,
	}

	// Send AWS credential and KMS keyId to the server running in enclave
	b, err := json.Marshal(playload)
	if err != nil {
		log.Fatal(err)
	}
	unix.Write(socket, b)
	// receive data from the server and save to dynamodb with the walletName
	response := make([]byte, 4096)
	n, err := unix.Read(socket, response)
	if err != nil {
		fmt.Println(err)
	}
	signedValue := hexutil.Encode(response[:n])
	return signedValue
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
	body, _ := io.ReadAll(res.Body)
	res.Body.Close()
	instanceProfileName := string(body)
	profileUri := fmt.Sprintf("http://169.254.169.254/latest/meta-data/iam/security-credentials/%s", instanceProfileName)
	res, err = http.Get(profileUri)
	if err != nil {
		log.Fatal(err)
	}
	body, _ = io.ReadAll(res.Body)
	res.Body.Close()
	var result iamCredentialToken
	json.Unmarshal(body, &result)
	token.aws_access_key_id = result.AccessKeyId
	token.aws_secret_access_key = result.SecretAccessKey
	token.aws_session_token = result.Token

	return token
}

func main() {

	region := "ap-northeast-1"
	keyId := "0f360b0f-1ad4-4c6b-b405-932d2f606779"
	walletAccountName := "account1"
	tableName :="AccountTable"
	signedTableName :="SignedValueTable"
	
	// check dynamodb AccountTable exist or not, create it if not exists
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region)},
	)

	if err != nil {
		panic(err)
	}

	svc := dynamodb.New(sess)
	describe_input := &dynamodb.DescribeTableInput{
		TableName: aws.String(tableName),
	}

	result, err := svc.DescribeTable(describe_input)
	if err != nil {
		fmt.Println(err)
		fmt.Println(result)
		fmt.Println("create the table",tableName)
		create_input := &dynamodb.CreateTableInput{
			AttributeDefinitions: []*dynamodb.AttributeDefinition{
				{
					AttributeName: aws.String("KeyId"),
					AttributeType: aws.String("S"),
				},
				{
					AttributeName: aws.String("Name"),
					AttributeType: aws.String("S"),
				},
			},
			KeySchema: []*dynamodb.KeySchemaElement{
				{
					AttributeName: aws.String("KeyId"),
					KeyType:       aws.String("HASH"),
				},
				{
					AttributeName: aws.String("Name"),
					KeyType:       aws.String("RANGE"),
				},
			},
			ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
				ReadCapacityUnits:  aws.Int64(5),
				WriteCapacityUnits: aws.Int64(5),
			},
			TableName: aws.String(tableName),
		}
		
		result, err := svc.CreateTable(create_input)
		if err != nil {
			fmt.Println(err.Error())
			fmt.Println(result)
		}
	}

	client := accountClient{region, tableName, keyId, 16, 5000}
	client.generateAccount(walletAccountName)

	//test sign
	transaction := map[string]interface{}{
		"value":    1000000000,
		"to":       "0xF0109fC8DF283027b6285cc889F5aA624EaC1F55",
		"nonce":    0,
		"chainId":  4,
		"gas":      100000,
		"gasPrice": 234567897654321,
	}

	b := new(bytes.Buffer)
	for key, value := range transaction {
		fmt.Fprintf(b, "%s=\"%s\"\n", key, value)
	}

	signedValue := client.sign(keyId, walletAccountName, b.String())
	// for demo, automatically create dynamodb table and save the signed value to it
	describe_input = &dynamodb.DescribeTableInput{
		TableName: aws.String(signedTableName),
	}

	result, err = svc.DescribeTable(describe_input)
	if err != nil {
		fmt.Println(err)
		fmt.Println(result)
		fmt.Println("create the table",signedTableName)
		create_input := &dynamodb.CreateTableInput{
			AttributeDefinitions: []*dynamodb.AttributeDefinition{
				{
					AttributeName: aws.String("Name"),
					AttributeType: aws.String("S"),
				},
				{
					AttributeName: aws.String("Transaction"),
					AttributeType: aws.String("S"),
				},
			},
			KeySchema: []*dynamodb.KeySchemaElement{
				{
					AttributeName: aws.String("Name"),
					KeyType:       aws.String("HASH"),
				},
				{
					AttributeName: aws.String("Transaction"),
					KeyType:       aws.String("RANGE"),
				},
			},
			ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
				ReadCapacityUnits:  aws.Int64(5),
				WriteCapacityUnits: aws.Int64(5),
			},
			TableName: aws.String(signedTableName),
		}
		
		result, err := svc.CreateTable(create_input)
		if err != nil {
			fmt.Println(err.Error())
			fmt.Println(result)
		}
	}

	// save to table
	svt := signedValueTable{
		Name:           walletAccountName,
		Transaction: 	b.String(),
		SignedValue: 	signedValue,
	}

	sv, err := dynamodbattribute.MarshalMap(svt)

	if err != nil {
		fmt.Println("Got error marshalling map:")
		fmt.Println(err.Error())
	}

	input := &dynamodb.PutItemInput{
		Item:      sv,
		TableName: aws.String(signedTableName),
	}

	_, err = svc.PutItem(input)
	if err != nil {
		fmt.Println("Got error calling PutItem:")
		fmt.Println(err.Error())
	}
}