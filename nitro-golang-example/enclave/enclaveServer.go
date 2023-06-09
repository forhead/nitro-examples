package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/ecdsa"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"

	"golang.org/x/sys/unix"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

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

type generateDataKeyResponse struct {
	datakey_plaintext_base64  string
	datakey_ciphertext_base64 string
}

type generateAccountResponse struct {
	EncryptedPrivateKey string
	Address             string
	EncryptedDataKey    string
}

func call_kms_generate_datakey(aws_access_key_id string, aws_secret_access_key string, aws_session_token string, keyId string) generateDataKeyResponse {
	var result generateDataKeyResponse
	cmd := exec.Command(
		"/app/kmstool_enclave_cli",
		"genkey",
		"--region", os.Getenv("REGION"),
		"--proxy-port", "8000",
		"--aws-access-key-id", aws_access_key_id,
		"--aws-secret-access-key", aws_secret_access_key,
		"--aws-session-token", aws_session_token,
		"--key-id", keyId,
		"--key-spec", "AES-256")

	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		fmt.Println("kms generate datakey error:", err)
	}

	datakey_split := strings.Split(out.String(), "\n")
	datakeyCiphertext_base64 := strings.TrimSpace(strings.Split(datakey_split[0], ":")[1])
	datakeyPlaintext_base64 := strings.TrimSpace(strings.Split(datakey_split[1], ":")[1])
	result.datakey_plaintext_base64 = datakeyPlaintext_base64
	result.datakey_ciphertext_base64 = datakeyCiphertext_base64

	return result
}

func call_kms_decrypt(aws_access_key_id string, aws_secret_access_key string, aws_session_token string, ciphertext string) string {
	cmd := exec.Command(
		"/app/kmstool_enclave_cli",
		"decrypt",
		"--region", os.Getenv("REGION"),
		"--proxy-port", "8000",
		"--aws-access-key-id", aws_access_key_id,
		"--aws-secret-access-key", aws_secret_access_key,
		"--aws-session-token", aws_session_token,
		"--ciphertext", ciphertext)

	fmt.Println("datakey:", ciphertext)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		log.Fatal("kms call", err)
	}
	result := out.String()
	fmt.Println("decrypt result:", result)
	return result
}

func generateAccount(aws_access_key_id string, aws_secret_access_key string, aws_session_token string, keyId string) generateAccountResponse {
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		log.Fatal(err)
	}

	privateKeyBytes := crypto.FromECDSA(privateKey)
	fmt.Println("SAVE BUT DO NOT SHARE THIS (Private Key):", hexutil.Encode(privateKeyBytes))

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}

	publicKeyBytes := crypto.FromECDSAPub(publicKeyECDSA)
	fmt.Println("Public Key:", hexutil.Encode(publicKeyBytes))

	address := crypto.PubkeyToAddress(*publicKeyECDSA).Hex()
	fmt.Println("Address:", address)

	datakeys := call_kms_generate_datakey(aws_access_key_id, aws_secret_access_key, aws_session_token, keyId)
	datakey_plaintext_base64 := datakeys.datakey_plaintext_base64
	datakey_ciphertext_base64 := datakeys.datakey_ciphertext_base64

	datakey_plaintext, _ := base64.StdEncoding.DecodeString(datakey_plaintext_base64)

	encryptedPrivateKey := encrypt([]byte(datakey_plaintext), string(privateKeyBytes))

	response := generateAccountResponse{
		EncryptedPrivateKey: encryptedPrivateKey,
		Address:             address,
		EncryptedDataKey:    datakey_ciphertext_base64,
	}
	return response
}

func sign(aws_access_key_id string, aws_secret_access_key string, aws_session_token string, encryptedDataKey string, encryptedPrivateKey string, transaction string) []byte {
	datakey_plaintext_base64 := call_kms_decrypt(aws_access_key_id, aws_secret_access_key, aws_session_token, encryptedDataKey)
	datakey_plaintext_base64_string := strings.TrimSpace(strings.Split(datakey_plaintext_base64, ":")[1])
	datakey_plaintext, err := base64.StdEncoding.DecodeString(datakey_plaintext_base64_string)
	if err != nil {
		log.Fatal("datakey", err)
	}
	private_key := decrypt(datakey_plaintext, encryptedPrivateKey)
	privateKey, err := crypto.ToECDSA([]byte(private_key))
	if err != nil {
		log.Fatal("privateKey error", err)
	}
	data := []byte(transaction)
	hash := crypto.Keccak256Hash(data)
	signature, err := crypto.Sign(hash.Bytes(), privateKey)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("hex:", hexutil.Encode(signature))
	fmt.Println("byte", signature)

	return signature
}

func encrypt(key []byte, message string) string {
	//Create byte array from the input string
	plainText := []byte(message)

	//Create a new AES cipher using the key
	block, err := aes.NewCipher(key)

	//IF NewCipher failed, exit:
	if err != nil {
		log.Fatal(err)
	}

	//Make the cipher text a byte array of size BlockSize + the length of the message
	cipherText := make([]byte, aes.BlockSize+len(plainText))

	//iv is the ciphertext up to the blocksize (16)
	iv := cipherText[:aes.BlockSize]
	if _, err = io.ReadFull(rand.Reader, iv); err != nil {
		log.Fatal(err)
	}

	//Encrypt the data:
	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(cipherText[aes.BlockSize:], plainText)

	//Return string encoded in base64
	return base64.RawStdEncoding.EncodeToString(cipherText)
}

func decrypt(key []byte, secure string) string {
	//Remove base64 encoding:
	cipherText, err := base64.RawStdEncoding.DecodeString(secure)

	//IF DecodeString failed, exit:
	if err != nil {
		log.Fatal(err)
	}

	//Create a new AES cipher with the key and encrypted message
	block, err := aes.NewCipher(key)

	//IF NewCipher failed, exit:
	if err != nil {
		log.Fatal(err)
	}

	//IF the length of the cipherText is less than 16 Bytes:
	if len(cipherText) < aes.BlockSize {
		fmt.Println("Ciphertext block size is too short!")
	}

	iv := cipherText[:aes.BlockSize]
	cipherText = cipherText[aes.BlockSize:]

	//Decrypt the message
	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(cipherText, cipherText)

	return string(cipherText)
}

func main() {
	fmt.Println("Start nitro enclave vsock server...")

	fd, err := unix.Socket(unix.AF_VSOCK, unix.SOCK_STREAM, 0)
	if err != nil {
		log.Fatal(err)
	}
	// Bind socket to cid 16, port 5000.
	sockaddr := &unix.SockaddrVM{
		CID:  unix.VMADDR_CID_ANY,
		Port: 5000,
	}
	err = unix.Bind(fd, sockaddr)
	if err != nil {
		log.Fatal("Bind ", err)
	}
	// Listen for up to 32 incoming connections.
	err = unix.Listen(fd, 32)
	if err != nil {
		log.Fatal("Listen ", err)
	}

	for {
		nfd, fromSockAdde, err := unix.Accept(fd)
		if err != nil {
			log.Fatal("Accept ", err)
		}
		fmt.Println("fromSockAdde: ", fromSockAdde)
		fmt.Println("conn is: ", nfd)

		requestData := make([]byte, 4096)
		var playload requestPlayload

		n, err := unix.Read(nfd, requestData)
		if err != nil {
			log.Fatal("Accept ", err)
		}

		err = json.Unmarshal(requestData[:n], &playload)
		if err != nil {
			fmt.Println(err.Error())
		}

		fmt.Println("apicall:", playload.ApiCall)

		apiCall := playload.ApiCall
		fmt.Println(apiCall)

		if apiCall == "generateAccount" {
			result := generateAccount(playload.Aws_access_key_id, playload.Aws_secret_access_key,
				playload.Aws_session_token, playload.KeyId)

			b, err := json.Marshal(result)
			if err != nil {
				fmt.Println(err)
			}
			//  send back to parent instance
			unix.Write(nfd, b)
			fmt.Println("generateAccount finished")
		} else if apiCall == "sign" {
			fmt.Println("sign request")
			result := sign(playload.Aws_access_key_id, playload.Aws_secret_access_key, playload.Aws_session_token,
				playload.EncryptedDataKey, playload.EncryptedPrivateKey, playload.Transaction)
			fmt.Println("result is:", result)
			unix.Write(nfd, result)
			fmt.Println("sign fihished")
		} else {
			fmt.Println("nothing to do")
		}
		unix.Close(nfd)
	}
}
