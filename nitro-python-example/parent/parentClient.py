import ecdsa
import socket
import requests
import json
import boto3
import hashlib

"""
dynamodb design 
table name: walletTable

colume:
walletName: wallet name for this wallet
encryptedPrivateKey: encrypted wallet private key
publicKey: the public key of the wallet
encryptedDatakey: the data key used to encrypt the private key
keyId: kms alias id which used for encryption for the private key
"""

class walletClient:

    """
    ddbTableName: the dynamodb table name which used to store and retrive the encrypted 
    keyID: the kms alais id, used to encrypt the plaintext and decrypt the ciphertext
    cid: cid for vsock client to connect
    port: port for vsock client to connect
    """
    def __init__(self, ddbTableName, keyId, cid, port):
        self.__ddbTableName = ddbTableName
        self.__keyId = keyId
        self.__cid = cid
        self.__port = port

    """
    generateWallet which will send credential and kms alais id to enclave server, will wait enclave server send the encrypted private key back, then save to dynamodb
    walletName: name for the wallet
    """
    def generateWallet(self, walletName):
        credential = self.__getIAMToken()
        # Create a vsock socket object
        s = socket.socket(socket.AF_VSOCK, socket.SOCK_STREAM)
        # Connect to the server
        s.connect((self.__cid, self.__port))

        playload = {}
        playload['apiCall'] = "generateWallet"
        playload['credential'] = credential
        playload['keyId'] = self.__keyId
        # Send AWS credential and KMS keyId to the server running in enclave
        s.send(str.encode(json.dumps(playload)))

        # receive data from the server and save to dynamodb with the walletName
        response = s.recv(65536).decode()
        # plaintext = base64.b64decode(response['Plaintext']).decode()
        self.__saveEncryptWalletToDDB(walletName, response, self.__keyId)
        s.close()

       
    def sign(self, keyId, walletName, message):
        # Get item from DynamoDB
        dynamodb = boto3.resource('dynamodb' , region_name = 'ap-southeast-1')
        table = dynamodb.Table(self.__ddbTableName)
        try:
            response = table.get_item(Key={
                'keyId': keyId,
                'walletName': walletName
                })
            if 'Item' not in response:
                print('walletName ' + walletName + ' not found in DynamoDB')
                return
        except Exception as error:
            print(error)
            return
        
        encryptedPrivateKey = response['Item']['encryptedPrivateKey']
        publicKey = response['Item']['publicKey']
        encryptedDatakey = response['Item']['encryptedDatakey']

   
        
        credential = self.__getIAMToken()

        playload = {}
        playload['apiCall'] = "sign"
        playload['credential'] = credential
        playload['encryptedPrivateKey'] = encryptedPrivateKey
        playload['encryptedDatakey'] = encryptedDatakey
        playload['message'] = message

        s = socket.socket(socket.AF_VSOCK, socket.SOCK_STREAM)
        # Connect to the server
        s.connect((self.__cid, self.__port))
        # Send AWS credential and KMS keyId to the server running in enclave
        s.send(str.encode(json.dumps(playload)))

        # receive data from the server and save to dynamodb with the walletName
        response = s.recv(65536)
        # plaintext = base64.b64decode(response['Plaintext']).decode()
        s.close()
        # Generate Verifying key from existing Public Key(string)
        vk = ecdsa.VerifyingKey.from_string(bytes.fromhex(publicKey), curve=ecdsa.SECP256k1, hashfunc=hashlib.sha256) # the default is sha1
        # Get signed message from Enclave, and verify the signing with Verifying key(public_key)
        bmessage = bytes(message, 'utf-8') 

        if vk.verify(response, bmessage): # True
            print('Signed message verified by public key: True')
        else:
            print('Signed message verified by public key: False')
        return response


    def __getIAMToken(self):
        """
        Get the AWS credential from EC2 instance metadata
        """
        r = requests.get("http://169.254.169.254/latest/meta-data/iam/security-credentials/")
        instance_profile_name = r.text

        r = requests.get("http://169.254.169.254/latest/meta-data/iam/security-credentials/%s" % instance_profile_name)
        response = r.json()

        credential = {
            'aws_access_key_id' : response['AccessKeyId'],
            'aws_secret_access_key' : response['SecretAccessKey'],
            'aws_session_token' : response['Token']
        }
        return credential

    def __saveEncryptWalletToDDB(self, walletName, response, keyId):
        dynamodb = boto3.resource('dynamodb','ap-southeast-1')
        table = dynamodb.Table(self.__ddbTableName)
        response_json = json.loads(response)
        print(response_json)
        table.put_item(Item={
            'walletName': walletName,
            'encryptedPrivateKey': response_json['encryptedPrivateKey'],
            'publicKey': response_json['publicKey'],
            'encryptedDatakey': response_json['encryptedDatakey'],
            'keyId': keyId
        })


def main():
 
    client = walletClient("demoWalletTable","your kms id",16,5000)
    client.generateWallet("wallet1")
    signedValue = client.sign('your kms id', "your message","hello")
    print("signed with value: ", signedValue)

if __name__ == '__main__':
    main()